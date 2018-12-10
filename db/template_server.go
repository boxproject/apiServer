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
	log "github.com/alecthomas/log4go"
	"github.com/boxproject/apiServer/utils"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"os"
	"database/sql"
	"encoding/json"
	"time"
	"github.com/go-errors/errors"
	"github.com/theplant/batchputs"
	"github.com/boxproject/apiServer/common"
)

type TemplateDBService struct {
}

func init(){
	//初始化
	initTemplateChain()
}

// NameExists check if the template name exists
func (*TemplateDBService) NameExists(name string) (bool, error) {
	var tp Template
	err := rdb.Where(&Template{Name: name}).First(&tp).Error
	if err != nil {

		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Create create a new template
func (*TemplateDBService) Create(tc *TemplateContent, content string, creatorID int) (string, error) {
	templateID := uuid.Must(uuid.NewV4()).String()
	templateHash := utils.GenHashStr(content)
	tepViewColumns := []string{"templateId", "coinId", "referAmount", "amountLimit", "referspan"}
	conn := rdb.Begin()
	var tp_db_struct TemplateView
	var tpovs [][]interface{}
	var users []Account
	err := conn.Where(&Account{UserType: common.OwnerAccType}).Find(&users).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Error("admin not found")
			return "", nil
		}
		return "", err
	}

	for i := 0; i < len(users); i++ {
		tpovs = append(tpovs, []interface{}{templateID, users[i].ID, 0, sql.NullString{}, users[i].AppId, users[i].Name})
	}

	err = conn.Create(&Template{ID: templateID, Hash: templateHash, Name: tc.Name, CreatorID: creatorID, Content: content, Period: tc.Period}).Error

	if err != nil {
		log.Error("新建审批流", err)
		conn.Rollback()
		return "", err
	}
	err = conn.Exec("insert into templateOper(templateId,accountId,status) (select ?,id,0 from account where userType =2)", templateID).Error // (*sql.Rows, error)

	if err != nil {
		log.Error("新建审批流", err)
		conn.Rollback()
		return "", err
	}

	// 获取审批流对应支持的币种信息
	var coins []Coin
	var coinNames []string
	var tpvs [][]interface{}
	// 获取该审批流支持的币种信息
	for i := 0; i < len(tc.LimitInfo); i++ {
		coinNames = append(coinNames, tc.LimitInfo[i].Symbol)
	}
	err = conn.Where(&Coin{Available: 0}).Where("name in (?)", coinNames).Find(&coins).Error
	if err != nil {
		log.Error("获取审批流支持的币种信息", err)
		conn.Rollback()
		return "", err
	}
	//
	for i := 0; i < len(coins); i++ {
		for j := 0; j < len(tc.LimitInfo); j++ {
			if tc.LimitInfo[j].Symbol == coins[i].Name && tc.LimitInfo[j].FullName == coins[i].FullName {
				tpvs = append(tpvs, []interface{}{templateID, coins[i].ID, tc.LimitInfo[j].Limit, tc.LimitInfo[j].Limit, tc.Period})
				continue
			}
		}
	}
	if len(tpvs) != len(tc.LimitInfo) {
		log.Error("存在不支持的币种")
		conn.Rollback()
		return "", errors.New("Invalidate Template Limit Info.")
	}
	if err != nil {
		log.Error("插入额度信息", err)
		conn.Rollback()
		return "", err
	}
	// 插入额度信息
	err = batchputs.Put(rdb.DB(), os.Getenv("DB_DIALECT"), tp_db_struct.TableName(), "templateId", tepViewColumns, tpvs)
	if err != nil {
		log.Error("初始化额度信息", err)
		conn.Rollback()
		return "", err
	}

	conn.Commit()
	return templateID, nil
}

type TempList struct {
	Id          string
	Name        string
	Hash        string
	Status      int
	ReferAmount string `gorm:"column:referAmount"`
	AmountLimit string `gorm:"column:amountLimit"`
}

// DelayTaskNum 统计待审批数量
func (*TemplateDBService) DelayTaskNum(approverID int) (task DelayTask) {
	var total int
	var name string
	row := rdb.Table("template").Select("count(0) as total").Joins("left join templateOper as tpo on tpo.templateId = template.id").Where("template.status = 0 and tpo.status = 0 and tpo.accountId = ?", approverID).Row()
	row.Scan(&total)
	row = rdb.Table("template").Select("template.name").Joins("left join templateOper as tpo on tpo.templateId = template.id").Where("template.status = 0 and tpo.status = 0 and tpo.accountId = ?", approverID).Order("createdAt desc").Limit(1).Row()
	row.Scan(&name)
	task.Number = total
	task.Reason = name
	return
}

// TemplateList 审批流列表
// types :0.全部  1.等待审批
func (*TemplateDBService) TemplateList(types int, userType int, accountId int) ([]TempList, error) {
	var tems []TempList
	var err error
	if types > 0 {
		err = rdb.Table("template").Joins("inner join templateOper on templateOper.status=0 and (template.id = templateOper.templateId)").Where("templateOper.accountId=? and template.status=0", accountId).Order("createdAt desc").Find(&tems).Error
	} else {
		if userType == 0 {
			err = rdb.Table("template").Select([]string{"id", "hash", "name", "status"}).Where("status in (?)", []int{common.TemplateApprovaled, common.TemplateToDisable, common.TemplateAvaliable, common.TemplateDisabled}).Order("createdAt desc").Find(&tems).Error
		} else {
			err = rdb.Table("template").Select([]string{"id", "hash", "name", "status"}).Find(&tems).Error
		}
	}
	return tems, err
}

// TemplistListAvailable 可用的审批流模板列表
func (*TemplateDBService) TemplistListAvailableByApprover(approver string) ([]Template, error) {
	var list []Template
	var result []Template
	// 获取审批流模板列表 0待审批 1审批通过 2审批拒绝 3.申请禁用4.上链 5禁用成功
	err := rdb.Where("status in (?)", []int{common.TemplatePending, common.TemplateApprovaled, common.TemplateAvaliable}).Find(&list).Error
	for _, v := range list {
		var tempContent TemplateContent
		json.Unmarshal([]byte(v.Content), &tempContent)
		for i := 0; i < len(tempContent.ApprovalInfo); i++ {
			for j := 0; j < len(tempContent.ApprovalInfo[i].Approvers); j++ {
				if tempContent.ApprovalInfo[i].Approvers[j].Account == approver {
					result = append(result, v)
					j = len(tempContent.ApprovalInfo)
				}
			}

		}
	}
	return result, err
}

// FindTemplateById 根据id查询审批
func (*TemplateDBService) FindTemplateById(tem *Template) (*Template, error) {
	err := rdb.First(&tem).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Error("find template by id", err)
		return nil, err
	}
	return tem, nil
}

// TemplateListTransfer 审批流列表--转账专用
func (*TemplateDBService) TemplateListTransfer(coinId int) ([]TempList, error) {
	var tems []TempList
	rdb.Exec("update templateView set  referAmount  = amountLimit where frozenTo < now()")
	err := rdb.Raw("select t.id,t.hash,t.name,tv.referAmount,tv.amountLimit from templateView tv left join template t on tv.templateId = t.id where t.status = 4 and tv.coinId = ?", coinId).Scan(&tems).Error
	return tems, err
}

// FindTemplateByOrderId 根据orderid查询审批
func (*TemplateDBService) FindTemplateByOrderId(orderId string) (*Template, bool) {
	var tem []Template
	rdb.Raw("select t.content as content from transferOrder tos left join template t on tos.hash = t.hash where tos.id = ? and t.status = 4 ", orderId).Scan(&tem)
	//db.Raw("SELECT name, age FROM users WHERE name = ?", 3).Scan(&result)

	if len(tem) == 0 {
		log.Error("查询错误")
		return nil, false
	}
	return &(tem[0]), true
}

// FindTemplateByHash 根据temid查询审批
func (*TemplateDBService) FindTemplateByHash(tem *Template) (*Template, bool) {
	err := rdb.Where(tem).First(tem).Error
	if err != nil {
		return nil, false
	}
	return tem, true
}

// FindTemplateView 根据temid 和 coinid查询审批额度
func (*TemplateDBService) FindTemplateView(tem *TemplateView) (*TemplateView, bool) {
	err := rdb.Where(tem).First(tem).Error
	if err != nil {
		return nil, false
	}
	return tem, true
}

// IntermediateStateHashList 获取中间状态的审批流模板列表
func (*TemplateDBService) IntermediateStateHashList() ([]Template, error) {
	var list []Template
	err := rdb.Model(&Template{}).Where("status in (?)", []int{common.TemplateApprovaled, common.TemplateToDisable}).Find(&list).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Error("Get Template list Error", err)
		return nil, err
	}
	return list, nil
}

// UpdateTemplateStatus 更新审批流模板状态
func (*TemplateDBService) UpdateTemplateStatus(hash string, status int) error {
	err := rdb.Model(&Template{}).Where(&Template{Hash: hash}).Update("status", status).Error
	if err != nil {
		log.Error("update template status", err)
		return err
	}
	return nil
}

// CancelTemplateAdnTxByIDs 根据ID批量作废审批流模板
func (*TemplateDBService) CancelTemplateAdnTxByIDs(hashs []string, template_statu int) ([]string, error) {
	var txOrderIds []string
	var txOrders []TransferOrder
	err := rdb.Table("transferOrder").Where("hash in (?) and status = 0", hashs).Find(&txOrders).Pluck("id", &txOrderIds).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error("获取审批中的转账申请")
		return nil, err
	}
	conn := rdb.Begin()
	// 作废审批流模板
	err = conn.Table("template").Where("hash in (?)", hashs).Update("status", template_statu).Error
	if err != nil {
		conn.Rollback()
		log.Error("批量作废审批流模板", err)
		return nil, err
	}
	// 作废对应的转账信息
	// 作废大订单
	err = conn.Model(&TransferOrder{Status: common.TransferToApproval}).Where("hash in (?) and status=0", hashs).Updates(&TransferOrder{Status: common.TransferDisableEmployee}).Error
	if err != nil {
		log.Error("作废相关转账申请", err)
		conn.Rollback()
		return nil, err
	}

	// 作废子订单
	err = conn.Model(&Transfer{}).Where("orderId in (?) and status=0", txOrderIds).Updates(&Transfer{Status: common.TransferDisableEmployee}).Error
	if err != nil {
		log.Error("作废子订单", err)
		conn.Rollback()
		return nil, err
	}
	conn.Commit()
	return txOrderIds, nil
}

// template
type TemplateVos struct {
	ID            string
	Hash          string
	Name          string
	CreatorID     int
	Content       string
	Status        int
	Period        int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	ApproveStatus int
}

// FindTemById 根据id查询审批
func (*TemplateDBService) FindTemById(tem *Template, accountId int) (*TemplateVos, error) {
	err := rdb.First(&tem).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Error("find template by id", err)
		return nil, err
	}
	var temo TemplateOper
	err = rdb.Where(&TemplateOper{AccountId: accountId, TemplateId: tem.ID}).First(&temo).Error
	tt := &TemplateVos{
		ID:            tem.ID,
		Hash:          tem.Hash,
		Name:          tem.Name,
		CreatorID:     tem.CreatorID,
		Content:       tem.Content,
		Status:        tem.Status,
		Period:        tem.Period,
		CreatedAt:     tem.CreatedAt,
		UpdatedAt:     tem.UpdatedAt,
		ApproveStatus: temo.Status,
	}
	return tt, nil
}

// FindTemByHashs 根据hash查询审批
func (*TemplateDBService) FindTemByHashs(hash []string) ([]Template, error) {
	var tems []Template
	err := rdb.Find(&tems, "hash in (?)", hash).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Error("find template by hash", err)
		return nil, err
	}
	return tems, nil
}

// TempStatistics 审批流统计
func (*TemplateDBService) TempStatistics(account string) (int, int, error) {
	var template []Template
	var content TemplateContent
	tempCount := 0
	var transferCount int
	tempErr := rdb.Find(&template).Error
	if tempErr != nil {
		return 0, 0, tempErr
	}
	if len(template) == 0 {
		return 0, 0, nil
	}
	var hashMap []string
	// 循环查找结果
	for tempi := 0; tempi < len(template); tempi++ {
		err := json.Unmarshal([]byte(template[tempi].Content), &content)
		if err != nil {
			return 0, 0, err
		}
		for infoi := 0; infoi < len(content.ApprovalInfo); infoi++ {
			for approveri := 0; approveri < len(content.ApprovalInfo[infoi].Approvers); approveri++ {
				if content.ApprovalInfo[infoi].Approvers[approveri].Account == account {
					tempCount++
					hashMap = append(hashMap, template[tempi].Hash)
				}
			}
		}
	}
	countErr := rdb.Model(&TransferOrder{}).Where("hash in (?) and status=0", hashMap).Count(&transferCount).Error
	if countErr != nil {
		return 0, 0, countErr
	}
	return tempCount, transferCount, tempErr
}

// CountTxByTempID 根据审批流统计转账
func (*TemplateDBService) CountTxByTempID(tempID string) (int, error) {
	var transferCount int
	err := rdb.Table("transferOrder").Joins("join template on template.hash = transferOrder.hash where template.id = ? and transferOrder.status = 0", tempID).Count(&transferCount).Error
	return transferCount, err
}

// 获取待上链的审批流模板
func (*TemplateDBService) TemplateListToUpload() ([]Template, error) {
	template_list := []Template{}
	err := rdb.Where("status in (?)", []int{common.TemplateToDisable, common.TemplateApprovaled}).Find(&template_list).Limit(5).Error
	return template_list, err
}

// UpdateTemplateLimit 通过截止时间更新币种额度
func (*TemplateDBService) UpdateTemplateLimit() {
	err := rdb.Exec("update templateView set referAmount = amountLimit where referspan !=0 and frozenTo < now()").Error
	if err != nil {
		log.Error("通过截止时间更新币种额度", err)
	}
}

// FindTemlateById grep template by id
func (*TemplateDBService) FindTemlateById(template *Template) error {
	err := rdb.First(template).Error
	return err
}

func initTemplateChain() {
	var tems []Template
	err := rdb.Where("status in (?)", []int{common.TemplateApprovaled, common.TemplateToDisable}).Find(&tems).Error
	if err != nil {
		log.Error("initTemplateChain 错误", err)
	}

	for i := 0; i < len(tems); i++ {
		TemplateChains.SetTemplateChain(tems[i])
		ETHTransferCh <- TransferSend{Types: common.TemplateToDisable}
	}
}
