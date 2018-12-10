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
package transfer

import (
	"encoding/json"

	log "github.com/alecthomas/log4go"
	"github.com/boxproject/apiServer/db"
	reterror "github.com/boxproject/apiServer/errors"
	msg "github.com/boxproject/apiServer/service/message"
	"github.com/satori/go.uuid"
	"database/sql"
	"sync"
	"time"
	"github.com/shopspring/decimal"
	"github.com/boxproject/apiServer/common"
	"github.com/jinzhu/gorm"
)

const WebPath = "webpage"

var transferLock = new(sync.Mutex)

// Findapplylist 查询转账申请列表 types 1.全部  2.发起的 3.参与的 ,4.批量审批列表
func Findapplylist(types, accountId int, account string) reterror.ErrModel {
	if types != 1 && types != 2 && types != 3 && types != 4 {
		return reterror.ErrModel{Code: reterror.Transfer_10001, Err: reterror.MSG_10001}
	}
	if accountId == 0 {
		return reterror.ErrModel{Code: reterror.Transfer_10001, Err: reterror.MSG_10001}
	}
	ts := &db.TransferDBService{}
	err := ts.UpdateTimeExpire()
	if err != nil {
		log.Error("更新状态错误 ", err)
		return reterror.ErrModel{Code: reterror.Transfer_10002, Err: reterror.MSG_10002}
	}
	list, err := ts.FindApplyList(types, accountId, account)
	//approveList, _ := ts.FindAllApprove()
	if err != nil {
		log.Error("列表错误 ", err)
		return reterror.ErrModel{Code: reterror.Transfer_10002, Err: reterror.MSG_10002}
	}
	return reterror.ErrModel{Code: reterror.Success, Data: list}
}

// FindApplyLog 查询转账log
func FindApplyLog(orderId string) reterror.ErrModel {
	as := &db.AccDBService{}
	if orderId == "" {
		return reterror.ErrModel{Code: reterror.Transfer_10001, Err: reterror.MSG_10001}
	}
	ts := &db.TransferDBService{}
	list, err := ts.FindApplyLog(orderId)

	if err != nil {
		log.Error("FindApplyLog查询log错误1", err)
		return reterror.ErrModel{Code: reterror.Transfer_10003, Err: reterror.MSG_10003}
	}
	data := make([]interface{}, 0)

	dd, err := ts.FindApply(orderId)

	if err != nil {
		log.Error("FindApplyLog查询log错误", err)
		return reterror.ErrModel{Code: reterror.Transfer_10003, Err: reterror.MSG_10003}
	}
	user, err := as.GetUserByID(dd.ApplyerId)
	if err != nil {
		log.Error("FindApplyLog查询loguser错误", err)
		return reterror.ErrModel{Code: reterror.Transfer_10003, Err: reterror.MSG_10003}
	}
	data = append(data, map[string]interface{}{
		"OrderNum":    orderId,
		"Status":      0,
		"UpdatedAt":   dd.CreatedAt,
		"AccountName": user[0].Name,
	})

	for _, v := range list {
		data = append(data, map[string]interface{}{
			"OrderNum":    v.OrderNum,
			"Status":      v.Status,
			"Reason":      v.Reason,
			"UpdatedAt":   v.UpdatedAt,
			"AccountName": v.AccountName,
			"Encode":      v.Encode,
		})
	}
	last := map[string]interface{}{
		"Status":    dd.Status,
		"UpdatedAt": dd.UpdatedAt,
	}
	datas := map[string]interface{}{
		"ApproversOper": data,
		"LastStatus":    last,
	}
	return reterror.ErrModel{Code: reterror.Success, Data: datas}
}

// Findapplybyid 通过订单id 查询明细
func Findapplybyid(id string) reterror.ErrModel {
	if id == "" {
		return reterror.ErrModel{Code: reterror.Transfer_10001, Err: reterror.MSG_10001}
	}
	ts := &db.TransferDBService{}
	detail, err := ts.Findapplybyid(id)
	if err != nil {
		log.Error("查找转账明细错误", err)
		return reterror.ErrModel{Code: reterror.Transfer_10005, Err: reterror.MSG_10005}
	}
	approvers, err := ts.FindAllApproveByOrderId(id)
	if err != nil {
		log.Error("查找审批人错误", err)
		return reterror.ErrModel{Code: reterror.Transfer_10005, Err: reterror.MSG_10005}
	}
	approversMap := make(map[string]interface{})
	for i := 0; i < len(approvers); i++ {
		approversMap[approvers[i].AccountName] = map[string]interface{}{
			"Status": approvers[i].Status,
			"sign":   approvers[i].Encode,
		}
	}
	detail.StrTemplate = ""
	detail.ApprovalContent = approversMap
	return reterror.ErrModel{Code: reterror.Success, Data: detail}
}

// Apply 申请转账
func Apply(transferVo *TransferVo) reterror.ErrModel {
	defer transferLock.Unlock()
	transferLock.Lock()
	ts := &db.TransferDBService{}
	temdb := &db.TemplateDBService{}
	orderId := uuid.Must(uuid.NewV4()).String()
	if transferVo.OrderId != "" {
		orderId = transferVo.OrderId
	}

	if transferVo.ApplySign == "" {
		log.Error("参数不能为空")
		return reterror.ErrModel{Code: reterror.Transfer_10001, Err: reterror.MSG_10001}
	}

	// 解析转账内容
	var applyMsg ApplyMsg
	err := json.Unmarshal([]byte(transferVo.ApplyMsg), &applyMsg)
	if err != nil {
		log.Error("解析转账内容", err)
		return reterror.ErrModel{Code: reterror.Transfer_10014, Err: reterror.MSG_10014}
	}
	if len(applyMsg.ApplyVos) == 0 {
		log.Error("参数不能为空applyvo")
		return reterror.ErrModel{Code: reterror.Transfer_10001, Err: reterror.MSG_10001}
	}

	total, err := decimal.NewFromString(applyMsg.Amount)
	if err != nil {
		log.Error("总金额必须为数字")
		return reterror.ErrModel{Code: reterror.Transfer_10001, Err: reterror.MSG_10001}
	}
	totalDecimal, _ := decimal.NewFromString("0")
	for j := 0; j < len(applyMsg.ApplyVos); j++ {
		num, b := countTotal(applyMsg.ApplyVos[j].Amount)
		if !b {
			log.Error("金额必须为数字")
			return reterror.ErrModel{Code: reterror.Transfer_10001, Err: reterror.MSG_10001}
		}
		totalDecimal = decimal.Sum(totalDecimal, *num)
	}
	if !total.Equal(totalDecimal) {
		log.Error("总金额不相等")
		return reterror.ErrModel{Code: reterror.Transfer_10001, Err: reterror.MSG_10001}
	}

	temParams := &db.Template{
		Status: common.TemplateAvaliable,
		Hash:   applyMsg.TemHash,
	}
	tem, _ := temdb.FindTemplateByHash(temParams)
	if tem == nil {
		log.Error("查找模板hash错误", applyMsg.TemHash)
		return reterror.ErrModel{Code: reterror.Transfer_10010, Err: reterror.MSG_10010}
	}
	templateView, _ := temdb.FindTemplateView(&db.TemplateView{TemplateId: tem.ID, CoinId: applyMsg.CoinId})
	if templateView == nil {
		log.Error("模板额度查找错误")
		return reterror.ErrModel{Code: reterror.Transfer_10033, Err: reterror.MSG_10033}
	}
	applyAmount, _ := decimal.NewFromString(applyMsg.Amount)
	referAmount, _ := decimal.NewFromString(templateView.ReferAmount)
	limitAmount, _ := decimal.NewFromString(templateView.AmountLimit)

	if limitAmount.LessThan(applyAmount) {
		log.Error("全部额度小于当前要申请的额度")
		return reterror.ErrModel{Code: reterror.Transfer_10034, Err: reterror.MSG_10034}
	}
	if templateView.Referspan != 0 {
		log.Debug("applyAmount", applyAmount, "referAmount", referAmount)
		if templateView.FrozenTo.After(time.Now()) && referAmount.LessThan(applyAmount) {
			log.Error("剩余额度小于当前要申请的额度")
			return reterror.ErrModel{Code: reterror.Transfer_10034, Err: reterror.MSG_10034}
		}
	}

	//将content转化为结构体
	var templateModel db.TemplateContent
	// 解析模板内容
	err = json.Unmarshal([]byte(tem.Content), &templateModel)
	if err != nil {
		log.Error("解析模板内容", err)
		return reterror.ErrModel{Code: reterror.Transfer_10011, Err: reterror.MSG_10011}
	}

	transferOrder := &db.TransferOrder{
		ID:          orderId,
		CoinName:    applyMsg.CoinName,
		Hash:        applyMsg.TemHash,
		ApplyReason: applyMsg.Reason,
		Miner:       applyMsg.Miner,
		Sign:        transferVo.ApplySign,
		Content:     transferVo.ApplyMsg,
		ApplyerId:   transferVo.ApplyerId,
		Amount:      applyMsg.Amount,
		NowLevel:    1,
	}
	if applyMsg.Deadline != "" {
		transferOrder.Deadline = sql.NullString{applyMsg.Deadline, true}
	} else {
		transferOrder.Deadline = sql.NullString{applyMsg.Deadline, false}
	}

	var arrApplys []*db.Transfer
	applys := applyMsg.ApplyVos
	fromAddress := ""
	if applyMsg.CoinName == "BTC" {
		fromAddress = db.AllMasterAddress["BTC"]
	} else if applyMsg.CoinName == "ETH" {
		fromAddress = db.AllMasterAddress["ETH"]
	} else if applyMsg.CoinName == "USDT" {
		fromAddress = db.AllMasterAddress["USDT"]
	} else if applyMsg.CoinName == "LTC" {
		fromAddress = db.AllMasterAddress["LTC"]
	} else {
		fromAddress = db.AllMasterAddress["ETH"]
	}
	for i := 0; i < len(applys); i++ {
		for j := 0; j < len(applys[i].Amount); j++ {
			transferId := uuid.Must(uuid.NewV4()).String()
			transfer := &db.Transfer{
				ID:          transferId,
				CoinName:    applyMsg.CoinName,
				Amount:      (applys[i].Amount)[j],
				ToAddr:      applys[i].ToAddress,
				OrderId:     orderId,
				Tag:         applys[i].Tag,
				ApplyReason: applyMsg.Reason,
				AmountIndex: (j + 1),
				FromAddress: fromAddress,
			}

			arrApplys = append(arrApplys, transfer)
		}
	}
	approvers := templateModel.ApprovalInfo[0].Approvers

	var arrApprover []*db.TransferReview
	for i := 0; i < len(approvers); i++ {
		approver := &db.TransferReview{
			OrderNum:    orderId,
			AccountName: approvers[i].Account,
		}
		arrApprover = append(arrApprover, approver)
	}
	temView := &db.TemplateView{
		TemplateId:  templateView.TemplateId,
		CoinId:      templateView.CoinId,
		Referspan:   templateView.Referspan,
		FrozenTo:    templateView.FrozenTo,
		AmountLimit: templateView.AmountLimit,
	}

	if templateView.Referspan != 0 {
		amountLimit, _ := decimal.NewFromString(temView.AmountLimit)
		if templateView.FrozenTo.Before(time.Now()) {
			temView.ReferAmount = amountLimit.Sub(applyAmount).String()
		} else {
			temView.ReferAmount = referAmount.Sub(applyAmount).String()
		}
	}

	log.Debug("approver", len(arrApprover))
	ret := ts.Apply(transferOrder, arrApplys, arrApprover, temView)
	if !ret {
		log.Error("申请保存数据库错误")
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	return reterror.ErrModel{Err: nil, Code: reterror.Success}
}

func countTotal(num []string) (*decimal.Decimal, bool) {
	numDecimal, _ := decimal.NewFromString("0")
	for i := 0; i < len(num); i++ {
		num, err := decimal.NewFromString(num[i])
		if err != nil {
			return nil, false
		} else {
			numDecimal = decimal.Sum(numDecimal, num)
		}
	}
	return &numDecimal, true

}

// Verify 审批转账 transferVo 转账内容 Status：1.同意 2.拒绝  7.app验证签名错误
func Verify(verifyApplyVo *VerifyApplyVo) reterror.ErrModel {
	log.Debug("verify transfer...")
	defer transferLock.Unlock()
	transferLock.Lock()
	temdb := &db.TemplateDBService{}
	ts := &db.TransferDBService{}
	log.Debug("approver account", verifyApplyVo.AppName)
	//判断参数
	if verifyApplyVo.OrderId == "" {
		log.Error("参数错误")
		return reterror.ErrModel{Code: reterror.Transfer_10001, Err: reterror.MSG_10001}
	}
	if !((verifyApplyVo.Status == common.TransferApprovaled) || (verifyApplyVo.Status == common.TransferReject) || (verifyApplyVo.Status == common.TransferSomeSuccess)) {
		log.Error("参数状态错误")
		return reterror.ErrModel{Code: reterror.ParamsNil, Err: reterror.PARAMS_NULL}
	}

	//判断是否该模板存在且状态正确
	tem, _ := temdb.FindTemplateByOrderId(verifyApplyVo.OrderId)
	if tem == nil {
		log.Error("模板不存在或者状态不正确")
		return reterror.ErrModel{Code: reterror.Transfer_10013, Err: reterror.MSG_10013}
	}

	//将content转化为结构体
	var templateModel db.TemplateContent
	// 解析模板内容
	err := json.Unmarshal([]byte(tem.Content), &templateModel)
	if err != nil {
		log.Error("解析模板内容", err)
		return reterror.ErrModel{Code: reterror.Transfer_10014, Err: reterror.MSG_10014}
	}

	//通过orderid查找转账信息
	transferOrder, err := ts.FindOrderById(verifyApplyVo.OrderId)
	//查找订单错误
	if err != nil {
		log.Error("查找订单错误", err)
		return reterror.ErrModel{Code: reterror.Transfer_10015, Err: reterror.MSG_10015}
	}

	//数据库订单状态不对
	if transferOrder.Status != 0 {
		log.Error("数据库订单状态不对", err)
		return reterror.ErrModel{Code: reterror.Transfer_10016, Err: reterror.MSG_10016}
	}
	strDeadline := transferOrder.Deadline.String
	if strDeadline != "" {

		t, err := time.Parse(time.RFC3339, strDeadline)
		if err != nil {
			log.Error("时间转换错误", strDeadline)
			return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
		}
		if t.Before(time.Now()) {
			log.Error("该订单过期")
			ts.UpdateOrderExpire(&db.Transfer{OrderId: transferOrder.ID, Status: common.TransferDisuse})
			// 插入站内信
			msg.TransferApprovalTimeOut(transferOrder.CoinName, transferOrder.Amount, transferOrder.ID, transferOrder.ApplyerId)
			return reterror.ErrModel{Code: reterror.Transfer_10032, Err: reterror.MSG_10032}
		}
	}

	//当前第几层
	nowLevel := transferOrder.NowLevel
	//一共多少层
	totalLevel := len(templateModel.ApprovalInfo)

	//取出当前需要的审核人
	approvers := templateModel.ApprovalInfo[nowLevel-1]
	//取出数据库里该层的审批人
	dbApprovers, approvers_err := ts.FindApproveByOrderId(transferOrder)
	if approvers_err != nil {
		log.Error("取出数据库里该层的查询错误", err)
		return reterror.ErrModel{Code: reterror.Transfer_10017, Err: reterror.MSG_10017}
	}
	//判断是否存在
	var existApprove *db.TransferReview = nil
	//同意人数
	agreeNum := 0
	//拒绝人数
	refuseNum := 0
	//该层所有人数
	totalNum := len(approvers.Approvers)
	//该层需要几个人同意
	requireNum := approvers.Require

	for i, v := range dbApprovers {
		if dbApprovers[i].AccountName == verifyApplyVo.AppName {
			existApprove = v
			if v.Status != 0 {
				//该审批状态不对
				log.Error("该审批状态不对", err)
				return reterror.ErrModel{Code: reterror.Transfer_10018, Err: reterror.MSG_10018}
			}
		}
		if dbApprovers[i].Status == 1 {
			agreeNum++
		}
		if dbApprovers[i].Status == 2 {
			refuseNum++
		}
	}
	if existApprove == nil {
		//该用户不能操作
		log.Error("该用户不能操作")
		return reterror.ErrModel{Code: reterror.Transfer_10019, Err: reterror.MSG_10019}
	}
	if agreeNum >= approvers.Require {
		//该审批流已经通过，无需审批
		log.Error("该审批流已经通过，无需审批")
		return reterror.ErrModel{Code: reterror.Transfer_10020, Err: reterror.MSG_10020}
	}
	//拒绝人数已经超过，结束
	if (totalNum - refuseNum) < requireNum {
		log.Error("拒绝人数已经超过，结束")
		return reterror.ErrModel{Code: reterror.Transfer_10021, Err: reterror.MSG_10021}
	}

	// 如果status为3，则审批不通过
	if verifyApplyVo.Status == 3 {
		transferOrder.Status = 7
		reet := ts.VerifyFailed(transferOrder)
		if !reet {
			//操作数据库失败
			log.Error("操作数据库失败")
			return reterror.ErrModel{Code: reterror.Transfer_10022, Err: reterror.MSG_10022}
		} else {
			//审批不通过
			return reterror.ErrModel{Code: reterror.Transfer_10023, Err: reterror.MSG_10023}
		}
	}

	countVo := &CountVo{
		AgreeNum:   agreeNum,
		RefuseNum:  refuseNum,
		TotalNum:   totalNum,
		RequireNum: requireNum,
		NowLevel:   nowLevel,
		TotalLevel: totalLevel,
	}
	rett := false
	//如果同意
	if verifyApplyVo.Status == 1 {
		verifyApplyVo.TemInfo = tem.Content
		rett = verifyAgree(countVo, verifyApplyVo, existApprove, &templateModel, transferOrder)
	}
	//如果拒绝
	if verifyApplyVo.Status == 2 {
		rett = verifyRefuse(countVo, verifyApplyVo, existApprove, transferOrder)
	}
	if !rett {
		log.Error("数据库错误")
		return reterror.ErrModel{Code: reterror.ParamsNil, Err: reterror.PARAMS_NULL}
	}
	return reterror.ErrModel{Err: err, Code: reterror.Success}
}

//agree
func verifyAgree(countVo *CountVo, verifyApplyVo *VerifyApplyVo, existApprove *db.TransferReview, templateModel *db.TemplateContent, transferOrder *db.TransferOrder) bool {
	ts := &db.TransferDBService{}

	transferReviewEntity := &db.TransferReview{
		ID:     existApprove.ID,
		Status: verifyApplyVo.Status,
		Reason: verifyApplyVo.Reason,
		Encode: verifyApplyVo.TransferSign,
	}
	//如果还差至少两个人同意
	if (countVo.RequireNum - countVo.AgreeNum) > 1 {
		ret := ts.Verify(nil, nil, transferReviewEntity, nil, false)
		return ret
	}
	//如果还剩最后一个人
	if (countVo.RequireNum - countVo.AgreeNum) == 1 {
		//如果不是最后一层
		if countVo.TotalLevel > countVo.NowLevel {
			order := &db.TransferOrder{
				ID:       verifyApplyVo.OrderId,
				NowLevel: countVo.NowLevel + 1,
			}
			//创建下一层审批人
			approvers := templateModel.ApprovalInfo[countVo.NowLevel].Approvers
			var arrApprover []*db.TransferReview
			for i := 0; i < len(approvers); i++ {
				approver := &db.TransferReview{
					OrderNum:    verifyApplyVo.OrderId,
					AccountName: approvers[i].Account,
				}
				arrApprover = append(arrApprover, approver)
			}
			ret := ts.Verify(order, nil, transferReviewEntity, arrApprover, true)
			return ret
		}
		//如果是最后一层---------最终操作
		if countVo.TotalLevel == countVo.NowLevel {
			mapApprovers := getAllApprovers(verifyApplyVo)
			mapApprovers[verifyApplyVo.AppName] = verifyApplyVo.TransferSign
			jsonApprovers, _ := json.Marshal(mapApprovers)
			strApprovers := string(jsonApprovers)
			order := &db.TransferOrder{
				ID:            verifyApplyVo.OrderId,
				Status:        1,
				ApproversSign: strApprovers,
			}
			transfers := &db.Transfer{
				OrderId: verifyApplyVo.OrderId,
				Status:  1,
			}
			ret := ts.Verify(order, transfers, transferReviewEntity, nil, true)
			//插入站内信
			err := msg.SuccessVerifyTrans(verifyApplyVo.OrderId)
			if err != nil {
				log.Error("站内信错误", err)
			}
			if ret {
				transferOrder.ApproversSign = strApprovers
				go insertChannel(transferOrder, verifyApplyVo.TemInfo)
			}

			return ret
		}
	}
	//其他错误情况返回false
	return false
}

//将审批成功的订单塞进channel
func insertChannel(order *db.TransferOrder, temInfo string) {
	acc := &db.AccDBService{}
	trans := &db.TransferDBService{}
	account, err := acc.AccountByID(order.ApplyerId)
	if err != nil {
		log.Error("channel account查询错误", err)
	}
	transfers := trans.GetTransferByOrderId(order.ID)
	if len(transfers) == 0 {
		log.Error("channel transfers查询错误")
	}
	// 解析转账内容
	var applyMsg ApplyMsg
	err = json.Unmarshal([]byte(order.Content), &applyMsg)
	if err != nil {
		log.Error("解析转账内容", err)
	}

	transferSendOrder := db.TransferSend{
		CoinName:       order.CoinName,
		TransferMsg:    order.Content,
		ApplyAccount:   account.Name,
		ApplyerId:      account.ID,
		ApplyPublickey: account.PubKey,
		ApplySign:      order.Sign,
		ApproversSign:  order.ApproversSign,
		Miner:          order.Miner,
		Types:          1,
		TokenAddress:   applyMsg.Token,
		Deadline:       applyMsg.Deadline,
		Status:         1, //1审批成功，转账中
		OrderId:        order.ID,
		Currency:       applyMsg.Currency,
		TemInfo:        temInfo,
	}
	if order.CoinName == "BTC" {
		db.BTCTransferCh <- transferSendOrder
	} else if order.CoinName == "LTC" {
		db.LTCTransferCh <- transferSendOrder
	} else {
		for i := 0; i < len(transfers); i++ {
			transferSend := db.TransferSend{
				TransferId:     transfers[i].ID,
				CoinName:       order.CoinName,
				TransferMsg:    order.Content,
				ApplyAccount:   account.Name,
				ApplyerId:      account.ID,
				ApplyPublickey: account.PubKey,
				ApplySign:      order.Sign,
				ApproversSign:  order.ApproversSign,
				Miner:          order.Miner,
				Amount:         transfers[i].Amount,
				ToAddress:      transfers[i].ToAddr,
				AmountIndex:    transfers[i].AmountIndex,
				Types:          1,
				TokenAddress:   applyMsg.Token,
				Deadline:       applyMsg.Deadline,
				Status:         1, //1审批成功，转账中
				OrderId:        order.ID,
				Currency:       applyMsg.Currency,
				TemInfo:        temInfo,
			}
			if applyMsg.Currency == "ETH" || (applyMsg.Currency == "ERC20") {
				db.ETHTransferCh <- transferSend
			} else if applyMsg.Currency == "USDT" {
				db.BTCTransferCh <- transferSend
			}
		}
	}

}

//refuse
func verifyRefuse(countVo *CountVo, verifyApplyVo *VerifyApplyVo, existApprove *db.TransferReview, orders *db.TransferOrder) bool {
	ts := &db.TransferDBService{}
	transferReviewEntity := &db.TransferReview{
		ID:     existApprove.ID,
		Status: verifyApplyVo.Status,
		Reason: verifyApplyVo.Reason,
		Encode: verifyApplyVo.TransferSign,
	}
	//如果再拒绝一次，该转账失效  ----------最终操作
	if countVo.TotalNum-countVo.RefuseNum == countVo.RequireNum {
		order := &db.TransferOrder{
			ID:     verifyApplyVo.OrderId,
			Status: 2,
		}
		transfers := &db.Transfer{
			OrderId: verifyApplyVo.OrderId,
			Status:  2,
		}
		ret := ts.Verify(order, transfers, transferReviewEntity, nil, true)
		//插入站内信
		err := msg.FailVerifyTrans(verifyApplyVo.OrderId, verifyApplyVo.Reason)
		if err != nil {
			log.Error("站内信错误", err)
		}
		if ret {
			err := ts.RecoverAmount(orders)
			if err != nil {
				log.Error("更新额度错误", err)
			}
		}
		return ret

	}

	//如果再拒绝不失效
	if countVo.TotalNum-countVo.RefuseNum > countVo.RequireNum {
		ret := ts.Verify(nil, nil, transferReviewEntity, nil, false)
		return ret
	}
	return false
}

// Cancel 撤销转账
func Cancel(verifyApplyVo *VerifyApplyVo) reterror.ErrModel {
	log.Debug("cancel transfer...")
	ts := &db.TransferDBService{}
	//判断参数
	if verifyApplyVo.OrderId == "" {
		log.Error("参数错误")
		return reterror.ErrModel{Code: reterror.Transfer_10001, Err: reterror.MSG_10001}
	}
	//通过orderid查找转账信息
	transferOrder, err := ts.FindOrderById(verifyApplyVo.OrderId)
	if err != nil {
		return reterror.ErrModel{Code: reterror.Transfer_10026, Err: reterror.MSG_10026}
	}
	if transferOrder.Status != common.TransferToApproval {
		//状态不对
		log.Error("状态不对", transferOrder.Status)
		return reterror.ErrModel{Code: reterror.Transfer_10027, Err: reterror.MSG_10027}
	}
	if transferOrder.ApplyerId != verifyApplyVo.AccountId {
		//该用户不能操作
		log.Error("该用户不能操作", transferOrder.Status)
		return reterror.ErrModel{Code: reterror.Transfer_10028, Err: reterror.MSG_10028}
	}

	ret := ts.Cancel(transferOrder)
	if !ret {
		return reterror.ErrModel{Code: reterror.Transfer_10029, Err: reterror.MSG_10029}
	}
	// 插入站内信
	err = msg.CancelTransfer(transferOrder.ApplyerId, transferOrder.CoinName, transferOrder.Amount, transferOrder.ID, false)
	if err != nil {
		log.Error("站内信err", err)
	}

	return reterror.ErrModel{Err: err, Code: reterror.Success}
}

// getAllApprovers 取出所有的审批人
func getAllApprovers(verifyApplyVo *VerifyApplyVo) map[string]string {
	//获取所有审批人
	ts := &db.TransferDBService{}
	approvers, _ := ts.FindApproveByOrderId(&db.TransferOrder{ID: verifyApplyVo.OrderId, NowLevel: 0})
	mapApprovers := make(map[string]string)
	for i := 0; i < len(approvers); i++ {
		mapApprovers[approvers[i].AccountName] = approvers[i].Encode
	}
	return mapApprovers

}

// BatchList 审批列表
func BatchList(transferVo *TransferVo) reterror.ErrModel {
	//获取所有审批人
	log.Error("订单的ids", transferVo.ArrOrderIds)
	//拉去所有审批的列表
	var batchs = []interface{}{}
	arrIds := transferVo.ArrOrderIds
	for i := 0; i < len(arrIds); i++ {
		ret := Findapplybyid(arrIds[i])
		if ret.Code == reterror.Success {
			batchs = append(batchs, ret.Data)
		}
	}
	return reterror.ErrModel{Code: reterror.Success, Data: batchs}
}

// GetTemplatebyHash 通过模板hash查找模板
func GetTemplatebyHash(temHash string) reterror.ErrModel {
	if temHash == "" {
		return reterror.ErrModel{Code: reterror.Transfer_10001, Err: reterror.MSG_10001}
	}
	ts := &db.TemplateDBService{}
	tem, _ := ts.FindTemplateByHash(&db.Template{Hash: temHash})
	if tem == nil {
		log.Error("FindTemplateByHash error")
		return reterror.ErrModel{Code: reterror.Transfer_10001, Err: reterror.MSG_10001}
	}
	objTemplate := &db.TemplateContent{}
	err := json.Unmarshal([]byte(tem.Content), &objTemplate)

	if err != nil {
		return reterror.ErrModel{Code: reterror.Transfer_10001, Err: reterror.MSG_10001}
	}
	return reterror.ErrModel{Code: reterror.Success, Data: objTemplate}
}

// TransferCommit web提交转账
func TransferCommit(transfer *WebTransfer) reterror.ErrModel {
	ts := &db.TransferDBService{}
	log.Error("orderid", transfer.Id)
	if transfer.Id == "" || (transfer.Msg == "") {
		log.Error(transfer)
		return reterror.ErrModel{Code: reterror.Transfer_10001, Err: reterror.MSG_10001}
	}
	webTrans := &db.WebTransfers{
		TransferId:  transfer.Id,
		Msg:         transfer.Msg,
		CreatedTime: time.Now().UnixNano() / 1e6,
	}
	err := ts.CreateWebTransfer(webTrans)
	if err != nil {
		log.Error("TransferCommit错误", err)
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	return reterror.ErrModel{Err: nil, Code: reterror.Success}
}

// FindTranfersById app获取转账信息
func FindTranfersById(transfer *WebTransfer) reterror.ErrModel {
	ts := &db.TransferDBService{}
	if transfer.Id == "" {
		return reterror.ErrModel{Code: reterror.Transfer_10001, Err: reterror.MSG_10001}
	}

	webTrans := &db.WebTransfers{
		TransferId: transfer.Id,
	}
	err := ts.FindWebTransfer(webTrans)
	if err != nil {
		log.Error("FindTranfersById错误", err)
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	timeNow := time.Now().UnixNano() / 1e6
	log.Error(timeNow, webTrans.CreatedTime)
	if ((timeNow - webTrans.CreatedTime) / 1000) > 60 {
		return reterror.ErrModel{Code: reterror.Transfer_10035, Err: reterror.MSG_10035}
	}
	return reterror.ErrModel{Err: nil, Code: reterror.Success, Data: webTrans}
}

// GetWebRouter
func GetWebRouter() reterror.ErrModel {
	return reterror.ErrModel{Err: nil, Code: reterror.Success, Data: map[string]string{"router": WebPath}}
}

// GetCommitStatus web获取提交状态
func GetCommitStatus(transfer *WebTransfer) reterror.ErrModel {
	ts := &db.TransferDBService{}
	if transfer.Id == "" {
		return reterror.ErrModel{Code: reterror.Transfer_10001, Err: reterror.MSG_10001}
	}

	webTrans := &db.WebTransfers{
		TransferId: transfer.Id,
	}
	err := ts.FindWebTransfer(webTrans)
	if err != nil {
		log.Error("FindTranfersById error", err)
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	_, err = ts.FindOrderById(transfer.Id)
	timeNow := time.Now().UnixNano() / 1e6
	log.Error(timeNow, webTrans.CreatedTime)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			if ((timeNow - webTrans.CreatedTime) / 1000) > 60 {
				// 过期
				return reterror.ErrModel{Code: reterror.Success, Data: map[string]int{"status": WebTransExpired}}
			}
			// 未扫码
			return reterror.ErrModel{Code: reterror.Success, Data: map[string]int{"status": WebTransWaitingForQr}}
		}
		log.Error("FindTranfersById错误", err)
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	//  扫码后提交
	return reterror.ErrModel{Err: nil, Code: reterror.Success, Data: map[string]int{"status": WebTransSubmitted}}
}
