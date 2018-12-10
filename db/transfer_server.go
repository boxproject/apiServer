// Copyright 2018. bolaxy.org authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package db

import (
	"strings"
	"encoding/json"
	log "github.com/alecthomas/log4go"
	"github.com/jinzhu/gorm"
	"time"
	"github.com/boxproject/apiServer/common"
	"github.com/theplant/batchputs"
	"os"
	"fmt"
)

type TransferDBService struct {
}

var BTCTransferCh = make(chan TransferSend, 10000)
var ETHTransferCh = make(chan TransferSend, 20000)
var LTCTransferCh = make(chan TransferSend, 10000)
var MergeETHTransferCh = make(chan TransferSend, 1000)

func init() {
	go transferSends()
}

func (*TransferDBService) UpdateTimeExpire() error {
	var orders []TransferOrder
	var ordersIds []string
	tx := rdb.Begin()
	err := tx.Table("transferOrder").Where("status = 0 and  deadline <= now()").Find(&orders).Pluck("id", &ordersIds).Error

	if err != nil || (len(orders) == 0) {
		tx.Rollback()
		return err
	}

	err = tx.Exec("update transferOrder set status = 8 where id in (?)", ordersIds).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Exec("update transfer set status = 8 where orderId in (?)", ordersIds).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	tepViewColumns := []string{"title", "content", "type", "receiver", "param"}
	var tpvs [][]interface{}

	for i := 0; i < len(orders); i++ {
		param := make(map[string]interface{})
		param["id"] = orders[i].ID
		param_json, _ := json.Marshal(param)
		playerIds := []int{orders[i].ApplyerId}
		strPlayerIds, _ := json.Marshal(playerIds)
		content := fmt.Sprintf("您发起的转账申请：%s %s，转账失败。\n失败原因：超过截止时间", orders[i].CoinName, orders[i].Amount)
		tpvs = append(tpvs, []interface{}{common.MsgTitleTransferFail, content, 2, strPlayerIds, param_json})
	}
	err = batchputs.Put(rdb.DB(), os.Getenv("DB_DIALECT"), "message", "content", tepViewColumns, tpvs)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// FindApplyList 获取转账列表
// 1.全部  2.发起的 3.参与的 4.批量
func (*TransferDBService) FindApplyList(types, applyId int, account string) ([]bachInfo, error) {
	var tfo []bachInfo
	//var transorder []TransferOrder
	var err error
	if types == 1 {
		err = rdb.Table("transferOrder").Order("createdAt desc").Scan(&tfo).Error
		log.Error("types1", err)
	}
	if types == 2 {
		err = rdb.Table("transferOrder").Where("applyerId=? ", applyId).Order("createdAt desc").Scan(&tfo).Error
	}
	if types == 3 {
		rows, er := rdb.Table("transferOrder").Select("transferOrder.coinName,transferOrder.id,transferOrder.createdAt,transferOrder.amount,transferOrder.status,transferOrder.applyReason").
			Joins("join transferReview tr on transferOrder.id = tr.orderNum").
			Where(" tr.accountName = ?", account).Order("transferOrder.createdAt desc").
			Rows()
		err = er
		defer rows.Close()
		for rows.Next() {
			var transferOrder bachInfo
			rdb.ScanRows(rows, &transferOrder)
			tfo = append(tfo, transferOrder)
		}
	}

	if types == 4 {
		// 批量转账
		rows, er := rdb.Table("transferOrder").Select(`ac.name as applyName, transferOrder.coinName,transferOrder.id,transferOrder.amount,transferOrder.status,transferOrder.applyReason,transferOrder.createdAt`).
			Joins("left join transferReview as tr on transferOrder.id = tr.orderNum and tr.accountName = ?", account).
			Joins("left join account as ac on ac.id = transferOrder.applyerId").
			Where("transferOrder.status = 0 and tr.status = 0").Order("transferOrder.createdAt desc").
			Rows()
		err = er
		defer rows.Close()
		for rows.Next() {
			var transferOrder bachInfo
			rdb.ScanRows(rows, &transferOrder)
			tfo = append(tfo, transferOrder)
		}
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	} else {
		return tfo, nil
	}
}

// GetTemHashByOrderId 获取模版hash
func (*TransferDBService) GetTemHashByOrderId(id string) (string, error) {
	var order TransferOrder
	err := rdb.First(&order, "id = ?", id).Error
	return order.Hash, err
}

// DelayTaskNum 统计待审批数量和最新的待审批
func (*TransferDBService) DelayTaskNum(account string) (task DelayTask) {
	var total int
	var applyReason string
	row := rdb.Table("transferOrder").Select("count(0) as total").Joins("join transferReview tr on transferOrder.id = tr.orderNum").Where("transferOrder.status = 0 and tr.status = 0 and tr.accountName = ?", account).Row()
	row.Scan(&total)
	row = rdb.Table("transferOrder").Select("transferOrder.applyReason").Joins("join transferReview tr on transferOrder.id = tr.orderNum").Where("transferOrder.status = 0 and tr.status = 0 and tr.accountName = ?", account).Order("transferOrder.createdAt desc").Row()
	row.Scan(&applyReason)
	task.Number = total
	task.Reason = applyReason
	return
}

// FindAllApprove 获取所有待审批审批人
func (*TransferDBService) FindAllApprove(orderIds []string) ([]TransferReview, error) {
	var reviews []TransferReview
	err := rdb.Find(&reviews, "orderNum in (?)", orderIds).Error
	return reviews, err
}

// FindApplyLog 查询转账log
func (*TransferDBService) FindApplyLog(orderId string) ([]TransferReview, error) {
	var trs []TransferReview
	err := rdb.Where("orderNum = ? and status>0 and status < 3", orderId).Find(&trs).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	} else {
		return trs, nil
	}
}

// FindApply 通过id查找转账订单
func (*TransferDBService) FindApply(orderIds string) (*TransferOrder, error) {
	var reviews = &TransferOrder{
		ID: orderIds,
	}
	err := rdb.First(reviews).Error
	return reviews, err
}

// Findapplybyid 转账申请明细
func (*TransferDBService) Findapplybyid(id string) (*OrderDetail, error) {
	var order OrderDetail

	err := rdb.Table("transferOrder").Select("transferOrder.sign as applySign,transferOrder.id as orderId,transferOrder.content as content,transferOrder.nowLevel as nowLevel,"+
		"transferOrder.amount as amount,transferOrder.status as status,tem.content as strtemplate,transferOrder.coinName,a.name applyerName,"+
		"transferOrder.applyReason,transferOrder.createdAt,transferOrder.deadline,JSON_UNQUOTE(json_extract(transferOrder.content,'$.coin_fullname')) AS fullName").
		Joins("join account as a on a.id = transferOrder.applyerId and transferOrder.id = ?", id).
		Joins("join template as tem on tem.hash = transferOrder.hash").
		Scan(&order).Error
	if err != nil { //没有查询到订单信息
		return nil, err
	}
	var transfers []Transfer
	err = rdb.Table("transfer").Select("toAddress,amount,tag,status").Where("orderId=?", id).Order("toAddress", true).Scan(&transfers).Error
	if err != nil {
		return nil, err
	}
	count := len(transfers)
	vals := []interface{}{}
	var amountMap map[string]interface{}
	amountsMap := []interface{}{}
	var addressMap map[string]interface{}
	lastAddress := ""
	lastTag := ""
	for i := 0; i < count; i++ {
		transfer := transfers[i]
		if (i != 0) && (lastAddress != transfer.ToAddr) {
			addressMap = map[string]interface{}{"Address": lastAddress, "Amount": amountsMap, "Tag": lastTag}
			amountsMap = []interface{}{}
			vals = append(vals, addressMap)
		}
		amountMap = map[string]interface{}{"amount": transfer.Amount, "status": transfer.Status}
		amountsMap = append(amountsMap, amountMap)
		lastAddress = transfer.ToAddr
		lastTag = transfer.Tag

		if (i + 1) == count {
			addressMap = map[string]interface{}{"Address": lastAddress, "Amount": amountsMap, "Tag": lastTag}
			amountsMap = []interface{}{}
			vals = append(vals, addressMap)
		}
	}
	order.OrderAddress = vals
	objTemplate := &TemplateContent{}
	err = json.Unmarshal([]byte(order.StrTemplate), &objTemplate)

	if err != nil {
		log.Error("解析错误", err)
		return nil, err
	}
	order.Template = *objTemplate
	return &order, err
}

// GetApprovalStatus 查询审批人的审批状态
func (*TransferDBService) GetApprovalStatus(id, account string) int {
	var transferReview TransferReview
	err := rdb.First(&transferReview, "orderNum =? and accountName = ?", id, account).Error
	if err != nil {
		return 0
	}
	return transferReview.Status
}

// Apply 申请转账
func (*TransferDBService) Apply(order *TransferOrder, transfers []*Transfer, approvers []*TransferReview, temView *TemplateView) bool {
	tx := rdb.Begin()
	//保存order
	if err := tx.Create(order).Error; err != nil {
		tx.Rollback()
		log.Error("创建订单错误", err)
		return false
	}
	//保存转账列表
	sqlStr := "insert into transfer(id,coinName,amount,toAddress,orderId,tag,applyReason,amountIndex,fromAddress) values "
	vals := []interface{}{}
	rowSQL := "(?,?,?,?,?,?,?,?,?)"
	var inserts []string
	for _, elem := range transfers {
		inserts = append(inserts, rowSQL)
		vals = append(vals, elem.ID, elem.CoinName, elem.Amount, elem.ToAddr, elem.OrderId, elem.Tag, elem.ApplyReason, elem.AmountIndex, elem.FromAddress)
	}
	sqlStr = sqlStr + strings.Join(inserts, ",")
	err := tx.Exec(sqlStr, vals...).Error
	if err != nil {
		tx.Rollback()
		log.Error("创建转账错误", err)
		return false
	}
	//保存第一层审批人列表
	sqlStr = "insert into transferReview(orderNum,accountName,level) values "
	vals = []interface{}{}
	rowSQL = "(?,?,?)"
	inserts = []string{}

	for _, elem := range approvers {
		inserts = append(inserts, rowSQL)
		vals = append(vals, elem.OrderNum, elem.AccountName, 1)
	}
	sqlStr = sqlStr + strings.Join(inserts, ",")
	err = tx.Exec(sqlStr, vals...).Error
	if err != nil {
		tx.Rollback()
		log.Error("创建审批人错误", err)
		return false
	}
	if temView.Referspan != 0 {
		var errview error
		//如果该时间已经过期 ，重置解冻时间
		if temView.FrozenTo.Before(time.Now()) {
			addTime := time.Duration(time.Hour.Nanoseconds() * int64(temView.Referspan))
			temView.FrozenTo = time.Now().Add(addTime)
			errview = tx.Exec("update templateView set referAmount = ?,frozenTo = ? where coinId = ? and templateId = ?",
				temView.ReferAmount, temView.FrozenTo, temView.CoinId, temView.TemplateId).Error

		} else {
			errview = tx.Exec("update templateView set referAmount = ? where coinId = ? and templateId = ?",
				temView.ReferAmount, temView.CoinId, temView.TemplateId).Error
		}
		if errview != nil {
			tx.Rollback()
			log.Error("更新模板额度错误", err)
			return false
		}
	}
	tx.Commit()
	return true
}

// Verify verify transfer
// params :
//	order ：  更新订单状态  （可选）
//	transfer ：  更新转账列表的状态	（可选）
//	transferReview ：  更新审批人的操作	（必选）
//	approvers  ：如果当前操作的人是该层的最后一个审批人，且还有下一层，则将下一层的审批人插入数据库 （可选）
//	isOver :如果当前操作是最终操作，需要把其他没操作的审批人状态更新为不可用 (必选)
func (*TransferDBService) Verify(order *TransferOrder, transfer *Transfer, transferReview *TransferReview, approvers []*TransferReview, isOver bool) bool {
	tx := rdb.Begin()

	if order != nil {
		if err := tx.Model(order).Update(order).Error; err != nil {
			tx.Rollback()
			log.Error("更新订单错误", err)
			return false
		}
	}
	if transfer != nil {
		err := tx.Model(&Transfer{}).Where("orderId = ?", order.ID).Update("status", transfer.Status).Error
		if err != nil {
			tx.Rollback()
			log.Error("更新转账列表错误", err)
			return false
		}
	}
	if isOver {
		err := tx.Model(&TransferReview{}).Where("orderNum = ? and status = 0", order.ID).Update("status", 3).Error
		if err != nil {
			tx.Rollback()
			log.Error("更新其他转账列表错误", err)
			return false
		}
	}
	if err := tx.Model(transferReview).Update(transferReview).Error; err != nil {

		tx.Rollback()
		log.Error("更新审批人的操作", err)
		return false
	}

	//保存下一层审批人列表
	if len(approvers) > 0 {
		sqlStr := "insert into transferReview(orderNum,accountName,level) values "
		vals := []interface{}{}
		rowSQL := "(?,?,?)"
		var inserts []string
		for _, elem := range approvers {
			inserts = append(inserts, rowSQL)
			vals = append(vals, elem.OrderNum, elem.AccountName, order.NowLevel)
		}
		sqlStr = sqlStr + strings.Join(inserts, ",")
		err := tx.Exec(sqlStr, vals...).Error
		if err != nil {
			tx.Rollback()
			log.Error("创建审批人错误", err)
		}
	}
	tx.Commit()
	return true
}

// VerifyFailed 审批验证签名失败
func (*TransferDBService) VerifyFailed(order *TransferOrder) bool {
	tx := rdb.Begin()
	if order != nil {
		if err := tx.Model(&TransferOrder{}).Where("id = ?", order.ID).Update("status", common.TransferInvalid).Error; err != nil {
			tx.Rollback()
			log.Error("VerifyFailed更新订单错误", err)
			return false
		}
	}

	if order != nil {
		if err := tx.Model(&Transfer{}).Where("orderId = ?", order.ID).Update("status", common.TransferInvalid).Error; err != nil {
			tx.Rollback()
			log.Error("VerifyFailed更新订单转账错误", err)
			return false
		}
	}
	tx.Commit()
	return true
}

// FindApproveByOrderId 通过orderid查询该层所有审批人
// tran.NowLevel:
// 	0 所有  其他第几层
func (*TransferDBService) FindApproveByOrderId(tran *TransferOrder) ([]*TransferReview, error) {
	var trans []*TransferReview
	var err error
	if tran.NowLevel == 0 {
		err = rdb.Where("orderNum = ? and status = 1", tran.ID).Find(&trans).Error
	} else {
		err = rdb.Where("orderNum = ? and level = ?", tran.ID, tran.NowLevel).Find(&trans).Error
	}
	return trans, err
}

// FindAllApproveByOrderId 通过orderid查询所有层所有审批人
func (*TransferDBService) FindAllApproveByOrderId(orderId string) ([]*TransferReview, error) {
	var trans []*TransferReview
	err := rdb.Where("orderNum = ?", orderId).Find(&trans).Error
	return trans, err
}

// FindOrderById 通过id查找订单
func (*TransferDBService) FindOrderById(orderId string) (*TransferOrder, error) {
	var orderInfo TransferOrder
	err := rdb.First(&orderInfo, "id = ?", orderId).Error
	return &orderInfo, err
}

// RecoverAmount 恢复额度
func (*TransferDBService) RecoverAmount(order *TransferOrder) error {
	var tv []TemplateView
	err := rdb.Raw("select tv.* from templateView tv left join template t on tv.templateId = t.id "+
		"left join coin c on c.id = tv.coinId where t.hash = ? and c.name = ?", order.Hash, order.CoinName).Scan(&tv).Error
	if err != nil {
		log.Error("查询错误", err)
		return err
	}
	temView := tv[0]
	if temView.Referspan != 0 {
		if temView.FrozenTo.After(time.Now()) {
			err = rdb.Exec("update templateView set referAmount = referAmount + ? where templateId = ? and coinId = ?", order.Amount, temView.TemplateId, temView.CoinId).Error
			if err != nil {
				rdb.Rollback()
				log.Error("更新余额错误", err)
				return err
			}
		}
	}
	return nil
}

// Cancel
func (*TransferDBService) Cancel(order *TransferOrder) bool {
	tx := rdb.Begin()
	var tv []TemplateView
	err := tx.Raw("select tv.* from templateView tv left join template t on tv.templateId = t.id "+
		"left join coin c on c.id = tv.coinId where t.hash = ? and c.name = ?", order.Hash, order.CoinName).Scan(&tv).Error
	if err != nil || len(tv) == 0 {
		tx.Rollback()
		log.Error("查询错误", err)
		return false
	}
	temView := tv[0]
	if temView.Referspan != 0 {
		if temView.FrozenTo.After(time.Now()) {
			err = tx.Exec("update templateView set referAmount = referAmount + ? where templateId = ? and coinId = ?", order.Amount, temView.TemplateId, temView.CoinId).Error
			if err != nil {
				tx.Rollback()
				log.Error("更新余额错误", err)
				return false
			}
		}
	}

	if err := tx.Model(&TransferOrder{}).Where("id = ?", order.ID).Update("status", common.TransferCancel).Error; err != nil {
		tx.Rollback()
		log.Error("更新订单错误", err)
		return false
	}
	err = tx.Model(&Transfer{}).Where(" orderId = ?", order.ID).Update("status", common.TransferCancel).Error
	if err != nil {
		tx.Rollback()
		log.Error("更新订单错误", err)
		return false
	}

	err = tx.Model(&TransferReview{}).Where("orderNum = ? and status = ?", order.ID, common.TransferToApproval).Update("status", common.TransferSomeSuccess).Error
	if err != nil {
		tx.Rollback()
		log.Error("更新其他转账列表错误", err)
		return false
	}
	tx.Commit()
	return true

}

type TransferSend struct {
	TransferId     string `gorm:"column:transferId"`
	CoinName       string `gorm:"column:coinName"`
	TransferMsg    string `gorm:"column:transferMsg"`
	ApplyAccount   string `gorm:"column:applyAccount"`
	ApplyerId      int    `gorm:"column:applyerId"`
	ApplyPublickey string `gorm:"column:applyPublickey"`
	ApplySign      string `gorm:"column:applySign"`
	ApproversSign  string `gorm:"column:approversSign"`
	Miner          string `gorm:"column:miner"`
	Amount         string `gorm:"column:amount"`
	ToAddress      string `gorm:"column:toAddress"`
	AmountIndex    int    `gorm:"column:amountIndex"`
	AddressMsg     string `gorm:"column:addressMsg"`
	Types          int    `gorm:"column:types"` //  1.外部转账 2。内部转账
	TokenAddress   string `gorm:"column:tokenAddress"`
	Deadline       string `gorm:"column:deadline"`
	Status         int    `gorm:"column:status"`
	OrderId        string `gorm:"column:orderId"`
	Currency       string `gorm:"column:currency"` //属于哪个主链币
	TemInfo        string `gorm:"column:temInfo"`
}

// timer专用
func transferSends() {
	var transfer []TransferSend
	err := rdb.Raw("select tem.content as temInfo,ts.status,ts.types as types, ts.addressMsg as addressMsg,ts.amountIndex as amountIndex,JSON_UNQUOTE(json_extract(tso.content,'$.token')) AS tokenAddress,JSON_UNQUOTE(json_extract(tso.content,'$.currency')) AS currency, " +
		" ts.id as transferId,ts.amount,ts.toAddress,ts.orderId, tso.miner as miner,ac.name as applyAccount,ac.pubKey as applyPublickey," +
		" ts.coinName as coinName,tso.sign as applySign,tso.approversSign,tso.content as transferMsg,JSON_UNQUOTE(json_extract(tso.content,'$.deadline')) AS deadline " +
		" from (select status,amountIndex,coinName,id,amount,toAddress,orderId,types,msg as addressMsg from transfer where  status = 1 or status = 12) ts " +
		" left join transferOrder tso on ts.orderId = tso.id left join account ac on tso.applyerId = ac.id " +
		" left join template tem on tem.hash = tso.hash ").Scan(&transfer).Error
	if err != nil {
		log.Error("查询错误")
	}
	for i := 0; i < len(transfer); i++ {
		if transfer[i].Currency == "ETH" || (transfer[i].Currency == "ERC20") {
			ETHTransferCh <- transfer[i]
		} else if transfer[i].Types == 2 {
			MergeETHTransferCh <- transfer[i]
		} else if transfer[i].Currency == "USDT" {
			BTCTransferCh <- transfer[i]
		}
	}
	err = rdb.Raw("select 1 as types,tso.id as orderId,tso.coinName,tem.content as temInfo,tso.content as transferMsg,ac.name as applyAccount," +
		" ac.id as applyerId,ac.pubKey as applyPublickey,tso.sign as applySign,tso.approversSign as approversSign,tso.miner as miner," +
		" tso.deadline as deadline,tso.status as status from transferOrder tso" +
		" left join template tem on tem.hash = tso.hash " +
		" left join account ac on tso.applyerId = ac.id " +
		" where tso.status = 1 or tso.status = 12").Scan(&transfer).Error
	if err != nil {
		log.Error("查询错误")
	}
	for i := 0; i < len(transfer); i++ {
		if transfer[i].CoinName == "BTC" {
			BTCTransferCh <- transfer[i]
		} else if transfer[i].CoinName == "LTC" {
			LTCTransferCh <- transfer[i]
		}
	}

}

func (*TransferDBService) BatchList(orderIds []string, account string) ([]TransferOrder, error) {
	var tfo []TransferOrder
	err := rdb.Table("transferOrder").Select("transferOrder.coinName,transferOrder.id,transferOrder.createdAt,transferOrder.amount,transferOrder.status,transferOrder.applyReason").
		Joins("join transferReview tr on transferOrder.id = tr.orderNum and tr.accountName = ? and transferOrder.id in (?)", account, orderIds).Scan(&tfo).Error
	return tfo, err
}

// 转账更新 timer专用
func (*TransferDBService) UpdateTransfer(transfer *Transfer) bool {
	tx := rdb.Begin()

	err := tx.Model(&Transfer{}).Where(" id = ?", transfer.ID).Update(map[string]interface{}{"status": transfer.Status, "txId": transfer.TxId, "updatedAt": time.Now()}).Error
	if err != nil {
		tx.Rollback()
		log.Error("timer更新转账错误", err)
		return false
	}
	tx.Commit()
	return true

}

// UpdateBalanceStatus 更新转账余额
func (*TransferDBService) UpdateBalanceStatus(id string) bool {
	err := rdb.Model(&Transfer{}).Where(" id = ?", id).Update("isUpdate", 2).Error
	if err != nil {
		log.Error("timer更新转账错误", err)
		return false
	}
	return true
}

// GetBalanceStatus 获取应该更新余额的转账
func (*TransferDBService) GetBalanceStatus() []TransferOrder {
	var transferOrder []TransferOrder
	err := rdb.Where(" isUpdate = 1").Find(&transferOrder).Error
	if err != nil {
		log.Error("获取更新余额转账错误", err)
		return nil
	}
	return transferOrder
}

// GetFirstOrder 获取一条不是最终状态的订单
func (*TransferDBService) GetFirstOrder() *TransferOrder {
	var transferOrder []TransferOrder
	err := rdb.Raw("select * from transferOrder where  status = 1").Scan(&transferOrder).Error
	if err != nil {
		log.Error("获取订单错误", err)
		return nil
	}
	if len(transferOrder) == 0 {
		return nil
	}
	return &transferOrder[0]
}

// GetTransferByOrderId 获取一个订单里的所有转账列表
func (*TransferDBService) GetTransferByOrderId(orderId string) []Transfer {
	var transfer []Transfer
	err := rdb.Raw("select * from transfer where  orderId = ?", orderId).Scan(&transfer).Error
	if err != nil {
		log.Error("获取订单错误", err)
		return transfer
	}
	return transfer
}

// UpdateOrderStatus 更新订单
func (*TransferDBService) UpdateOrderStatus(orderId string, status int) {
	var orderInfo TransferOrder
	ms := MessageService{}
	err := rdb.Table("transferOrder").Where("id = ?", orderId).First(&orderInfo).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Error("订单未找到", orderId)
		}
		log.Error("查询订单失败", err)
	}
	err = rdb.Model(&TransferOrder{}).Where("id = ?", orderId).Update("status", status).Error
	if err != nil {
		log.Error("更新订单错误", err)
	}

	if status == common.TransferSuccess {
		//err := msg.TransferSuccess(orderInfo.CoinName, orderInfo.Amount, orderId, orderInfo.ApplyerId)
		err = ms.TransferSuccess(orderInfo.CoinName, orderInfo.Amount, orderId, orderInfo.ApplyerId)
		if err != nil {
			log.Error("插入成功站内信error", err)
		}
	}

	if status == common.TransferFailed {
		//err := msg.TransferFail(orderInfo.CoinName, orderInfo.Amount, orderId, orderInfo.ApplyerId)
		err = ms.TransferFail(orderInfo.CoinName, orderInfo.Amount, orderId, orderInfo.ApplyerId, common.NormalMessage)
		if err != nil {
			log.Error("插入失败站内信error", err)
		}
	}

	if status == common.TransferSomeSuccess {
		err = ms.TransferPartialSuccessful(orderInfo.CoinName, orderInfo.Amount, orderInfo.ID, orderInfo.ApplyerId)
		if err != nil {
			log.Error("插入部分成功站内信error", err)
		}
	}
}

// UpdateOrderExpire 更新
func (*TransferDBService) UpdateOrderExpire(transfer *Transfer) {
	tx := rdb.Begin()

	err := tx.Model(&Transfer{}).Where("orderId = ?", transfer.OrderId).Update(transfer).Error

	if err != nil {
		tx.Rollback()
		log.Error("更新转账错误", err)
		return
	}
	if transfer.Status == common.TransferSomeSuccess {
		tx.Commit()
		return
	}

	if transfer.Status > common.TransferInsufficientBalance {
		transfer.Status = common.TransferFailed
	}

	err = tx.Model(&TransferOrder{}).Where("id = ?", transfer.OrderId).Update(transfer).Error
	if err != nil {
		tx.Rollback()
		log.Error("更新订单错误", err)
		return
	}
	tx.Commit()
}

// UpdateOrder
func (*TransferDBService) UpdateOrder(orderId string, status int) {
	tx := rdb.Begin()

	err := tx.Model(&TransferOrder{}).Where("id = ?", orderId).Update("status", status).Error
	if err != nil {
		tx.Rollback()
		log.Error("更新订单错误", err)
		return
	}
	err = tx.Model(&Transfer{}).Where("orderId = ?", orderId).Update("status", status).Error

	if err != nil {
		tx.Rollback()
		log.Error("更新转账错误", err)
		return
	}
	tx.Commit()

}

// CreateWebTransfer
func (*TransferDBService) CreateWebTransfer(transfers *WebTransfers) error {
	err := rdb.Create(transfers).Error
	return err
}
func (*TransferDBService) FindWebTransfer(transfers *WebTransfers) error {
	err := rdb.Where("transferId = ?", transfers.TransferId).First(transfers).Error
	return err
}
func (*TransferDBService) UpdateWebTransfer(transfers *WebTransfers) error {
	err := rdb.Model(transfers).Update(*transfers).Error
	return err
}
