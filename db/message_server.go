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
	"encoding/json"
	"fmt"
	"strings"
	log "github.com/alecthomas/log4go"
	"github.com/boxproject/apiServer/common"
	"time"
)

type MessageService struct {
}

// AddAdmin 添加管理员消息通知
func (*MessageService) AddAdmin(ids []int, accountId int) error {
	//组织发送给股东/管理员消息
	valueStrings := make([]string, 0, len(ids))
	valueArgs := make([]interface{}, 0, len(ids)*3)
	var operator Account //操作员
	err := rdb.First(&operator, "id=? and isDeleted =0", accountId).Error
	if err != nil {
		return err
	}
	var accounts []Account
	err = rdb.Find(&accounts, "userType in (1,2) and isDeleted =0 and id not in (?)", ids).Error //查询所有股东和管理员
	if err != nil {
		return err
	}
	revicer := make([]interface{}, 0, len(accounts))
	for _, acc := range accounts {
		revicer = append(revicer, acc.ID)
	}
	json_revicer, _ := json.Marshal(revicer)

	for _, id := range ids {
		var account Account
		err := rdb.First(&account, "id=? and isDeleted =0", id).Error
		if err != nil {
			return err
		}
		content := common.MSgContentAppointAdminNoticeOthers
		paddingOther := &common.MsgTemplatePadding{AccountType: common.MsgOwnerAccountType, OperAccountName: operator.Name, PromoterAccountName: account.Name}
		valueStrings = append(valueStrings, "(?,?,?,?)")
		valueArgs = append(valueArgs, common.MsgTitleAppointAdmin)
		paddingOtherByte, _ := json.Marshal(paddingOther)
		valueArgs = append(valueArgs, content)
		valueArgs = append(valueArgs, string(paddingOtherByte))
		valueArgs = append(valueArgs, string(json_revicer))
	}
	// 组织发送给被操作人消息
	paramJson, _ := json.Marshal(ids)
	valueStrings = append(valueStrings, "(?,?,?,?)")
	valueArgs = append(valueArgs, common.MsgTitleAppointAdmin)
	valueArgs = append(valueArgs, common.MsgContentAppointAdmin)
	padding := &common.MsgTemplatePadding{AccountType: common.MsgOwnerAccountType, OperAccountName: operator.Name}
	paddingByte, _ := json.Marshal(padding)
	valueArgs = append(valueArgs, string(paddingByte))
	valueArgs = append(valueArgs, string(paramJson))
	tx := rdb.Begin()
	stmt := fmt.Sprintf("insert into message (title,content,padding,receiver) VALUES %s", strings.Join(valueStrings, ","))
	err = tx.Exec(stmt, valueArgs...).Error
	if err != nil {
		tx.Rollback()
		log.Error("添加消息失败", err)
		return err
	}
	tx.Commit()
	return nil
}

// DelAdmin 取消管理员消息通知
func (*MessageService) DelAdmin(accountID, id int) (err error) {
	var operator Account   //操作人
	var account Account    //被操作人
	var accounts []Account //所有股东和管理员
	err = rdb.First(&operator, "id=? and isDeleted =0", id).Error
	err = rdb.First(&account, "id=? and isDeleted =0", accountID).Error
	err = rdb.Find(&accounts, "userType in (1,2) and isDeleted =0").Error
	if err != nil {
		return err
	}
	revicer := make([]interface{}, 0, len(accounts))
	for _, acc := range accounts {
		revicer = append(revicer, acc.ID)
	}
	json_revicer, _ := json.Marshal(revicer)
	padding := &common.MsgTemplatePadding{AccountType: common.MsgOwnerAccountType, OperAccountName: operator.Name, PromoterAccountName: account.Name}
	paddingByte, _ := json.Marshal(padding)
	tx := rdb.Begin()
	err = tx.Exec("insert into message (title,content,padding,receiver) VALUES (?,?,?,?)", common.MsgTitleDeleteAdmin, common.MsgContentDeleteAdminNoticeOthers, string(paddingByte), string(json_revicer)).Error
	if err != nil {
		tx.Rollback()
		log.Error("添加消息失败", err)
		return
	}
	revicer = make([]interface{}, 0, 1)
	revicer = append(revicer, accountID)
	json_revicer, _ = json.Marshal(revicer)
	err = tx.Exec("insert into message (title,content,padding,receiver) VALUES (?,?,?,?)", common.MsgTitleDeleteAdmin, common.MsgContentDeleteAdmin, string(paddingByte), string(json_revicer)).Error
	if err != nil {
		tx.Rollback()
		log.Error("添加消息失败", err)
		return
	}
	tx.Commit()
	return
}

// ChangeNode 切换节点
func (*MessageService) ChangeNode() (err error) {
	var accounts []Account //所有股东和管理员
	adminUserTypes := []int{common.AdminAccType, common.OwnerAccType}
	err = rdb.Find(&accounts, "userType in (?) and isDeleted =0", adminUserTypes).Error
	if err != nil {
		return err
	}
	revicer := make([]interface{}, 0, len(accounts))
	for _, acc := range accounts {
		revicer = append(revicer, acc.ID)
	}
	json_revicer, _ := json.Marshal(revicer)
	err = rdb.Exec("insert into message (title,content,receiver) VALUES (?,?,?)", common.MsgTitleNodeSwitch, common.MsgContentNodeSwitch, string(json_revicer)).Error
	return
}

// SuccessVerifyTemp  审批流 审批成功
func (*MessageService) SuccessVerifyTemp(name, id string, creatorId int) (err error) {
	var accounts []Account //所有股东
	err = rdb.Find(&accounts, "userType =2 and isDeleted =0").Error
	if err != nil {
		return err
	}
	revicer := make([]interface{}, 0, len(accounts)+1)
	for _, acc := range accounts {
		if acc.ID != creatorId {
			revicer = append(revicer, acc.ID)
		}
	}
	revicer = append(revicer, creatorId) //发送给股东和发起人
	json_revicer, _ := json.Marshal(revicer)

	param := make(map[string]interface{})
	param["id"] = id
	param_json, _ := json.Marshal(param)
	padding := common.MsgTemplatePadding{TemplateName: name}
	paddingByte, _ := json.Marshal(padding)
	err = rdb.Exec("insert into message (title,content,padding,receiver,type,param) VALUES (?,?,?,?,?,?)", common.MsgTitleApprovalTemplateSuccess, common.MsgContentApprovalTemplateSuccess, string(paddingByte), json_revicer, 1, string(param_json)).Error
	return
}

// FailVerifyTemp 审批流 审批失败
func (*MessageService) FailVerifyTemp(opername, name, id string, creatorId int) (err error) {
	var accounts []Account //所有股东
	adminUserTypes := []int{common.AdminAccType, common.OwnerAccType}
	err = rdb.Find(&accounts, "userType in (?) and isDeleted =0", adminUserTypes).Error
	if err != nil {
		return err
	}
	revicer := make([]interface{}, 0, len(accounts)+1)
	for _, acc := range accounts {
		if acc.ID != creatorId && acc.UserType == common.OwnerAccType {
			revicer = append(revicer, acc.ID)
		}
	}
	revicer = append(revicer, creatorId) //发送给股东和发起人
	json_revicer, _ := json.Marshal(revicer)

	param := make(map[string]interface{})
	param["id"] = id
	param_json, _ := json.Marshal(param)

	title := common.MsgTitleApprovalTemplateFail
	var content string
	var padding common.MsgTemplatePadding
	var warnType = common.NormalMessage
	if opername == "" {
		warnType = common.WarnMessage
		// 发给股东和管理员
		adminReceiver := []int{}
		for _, acc := range accounts {
			if acc.ID != creatorId {
				adminReceiver = append(adminReceiver, acc.ID)
			}
		}
		json_revicer, _ = json.Marshal(adminReceiver)
		content = common.MsgContentApprovalTemplateFailByVoucher
		padding = common.MsgTemplatePadding{TemplateName: name}
	} else {
		content = common.MsgContentApprovalTemplateFailByOwner
		padding = common.MsgTemplatePadding{TemplateName: name, OperAccountName: opername}
	}
	paddingByte, _ := json.Marshal(padding)
	err = rdb.Exec("insert into message (title,content,padding,receiver,type,param, warnType) VALUES (?,?,?,?,?,?,?)", title, content, string(paddingByte), json_revicer, 1, string(param_json), warnType).Error
	return
}

// SuccessVerifyTrans 转账审批成功消息
func (*MessageService) SuccessVerifyTrans(CoinName, Amount string, ApplyerId int, id string) (err error) {
	revicer := make([]interface{}, 0, 1)
	revicer = append(revicer, ApplyerId)
	json_revicer, _ := json.Marshal(revicer)

	param := make(map[string]interface{})
	param["id"] = id
	param_json, _ := json.Marshal(param)
	padding := &common.MsgTemplatePadding{TransOrderCoinName: CoinName, TransOrderAmount: Amount}
	paddingByte, _ := json.Marshal(padding)
	err = rdb.Exec("insert into message (title,content,padding,receiver,type,param) VALUES (?,?,?,?,?,?)", common.MsgTitleTransferApproved, common.MsgContentTransferApproved, string(paddingByte), string(json_revicer), 2, string(param_json)).Error
	return err
}

// FailVerifyTrans 转账审批失败消息
func (*MessageService) FailVerifyTrans(CoinName, Amount, reason string, ApplyerId int, id string) (err error) {
	revicer := make([]interface{}, 0, 1)
	revicer = append(revicer, ApplyerId)
	json_revicer, _ := json.Marshal(revicer)

	param := make(map[string]interface{})
	param["id"] = id
	param_json, _ := json.Marshal(param)
	padding := &common.MsgTemplatePadding{TransOrderCoinName: CoinName, TransOrderAmount: Amount}
	paddingByte, _ := json.Marshal(padding)
	err = rdb.Exec("insert into message (title,content,padding,receiver,type,param) VALUES (?,?,?,?,?,?)", common.MsgTitleTransferReject, common.MsgContentTransferReject, string(paddingByte), string(json_revicer), 2, string(param_json)).Error
	return err
}

// TransferSuccess 转账成功消息
func (*MessageService) TransferSuccess(coinName, amount, orderId string, applyerId int) (err error) {
	revicer := make([]interface{}, 0, 1)
	revicer = append(revicer, applyerId)
	json_revicer, _ := json.Marshal(revicer)

	param := make(map[string]interface{})
	param["id"] = orderId
	param_json, _ := json.Marshal(param)
	padding := &common.MsgTemplatePadding{TransOrderCoinName: coinName, TransOrderAmount: amount}
	paddingByte, _ := json.Marshal(padding)
	err = rdb.Exec("insert into message (title,content,padding,receiver,type,param) VALUES (?,?,?,?,?,?)", common.MsgTitleTransferSuccess, common.MsgContentTransferSuccess, string(paddingByte), string(json_revicer), 2, string(param_json)).Error
	return err
}

// TransferFail 转账失败消息
func (*MessageService) TransferFail(coinName, amount, transferId string, applyerId, warnType int) (err error) {
	revicer := []int{applyerId}
	if warnType == common.WarnMessage {
		var msgs []Message
		var transferInfo Transfer
		rdb.Where("id = ?", transferId).First(&transferInfo)
		transferId = transferInfo.OrderId
		// 过滤重复的orderid
		rdb.Where("warnType = ?", common.WarnMessage).Find(&msgs)
		for _, v := range msgs {
			var orderidParam map[string]string
			json.Unmarshal([]byte(v.Param), &orderidParam)
			if orderidParam["id"] == transferInfo.OrderId {
				return nil
			}
		}
	}
	json_revicer, _ := json.Marshal(revicer)
	param := make(map[string]interface{})
	param["id"] = transferId
	param_json, _ := json.Marshal(param)
	padding := &common.MsgTemplatePadding{TransOrderCoinName: coinName, TransOrderAmount: amount}
	paddingByte, _ := json.Marshal(padding)
	err = rdb.Exec("insert into message (title,content,padding,receiver,type,param,warnType) VALUES (?,?,?,?,?,?,?)", common.MsgTitleTransferFail, common.MsgContentTransferFail, string(paddingByte), string(json_revicer), 2, string(param_json), warnType).Error
	return err
}

// TransferPartialSuccessful 部分转账成功
func (*MessageService) TransferPartialSuccessful(coinName, amount, orderId string, applyerId int) (err error) {
	revicer := []int{applyerId}
	json_revicer, _ := json.Marshal(revicer)
	param := make(map[string]interface{})
	param["id"] = orderId
	param_json, _ := json.Marshal(param)
	padding := &common.MsgTemplatePadding{TransOrderCoinName: coinName, TransOrderAmount: amount}
	paddingByte, _ := json.Marshal(padding)
	err = rdb.Exec("insert into message (title,content,padding,receiver,type,param) VALUES (?,?,?,?,?,?)", common.MsgTitleTransferFail, common.MsgContentTransferPartiallySuceeded, string(paddingByte), string(json_revicer), 2, string(param_json)).Error
	return err
}

// TransferFailOutOfDate 由于过期而转账失败消息
func (*MessageService) TransferFailOutOfDate(coinName, amount, transferId string, applyerId, warnType int) (err error) {
	revicer := []int{applyerId}
	json_revicer, _ := json.Marshal(revicer)

	param := make(map[string]interface{})
	param["id"] = transferId
	param_json, _ := json.Marshal(param)
	padding := &common.MsgTemplatePadding{TransOrderCoinName: coinName, TransOrderAmount: amount}
	paddingByte, _ := json.Marshal(padding)
	err = rdb.Exec("insert into message (title,content,padding,receiver,type,param,warnType) VALUES (?,?,?,?,?,?,?)", common.MsgTitleTransferFail, common.MsgContentTransferFailCauseTimeout, string(paddingByte), string(json_revicer), 2, string(param_json), warnType).Error
	return err
}

func (*MessageService) TransferApprovalTimeOut(coinName, amount, transferId string, applyerId, warnType int) (err error) {
	revicer := []int{applyerId}
	json_revicer, _ := json.Marshal(revicer)

	param := make(map[string]interface{})
	param["id"] = transferId
	param_json, _ := json.Marshal(param)
	padding := &common.MsgTemplatePadding{TransOrderCoinName: coinName, TransOrderAmount: amount}
	paddingByte, _ := json.Marshal(padding)
	err = rdb.Exec("insert into message (title,content,padding,receiver,type,param,warnType) VALUES (?,?,?,?,?,?,?)", common.MsgTitleTransferFail, common.MsgContentTransferApprovalFaliCauseTimeOut, string(paddingByte), string(json_revicer), 2, string(param_json), warnType).Error
	return err
}

// CancelTemplate 作废审批流消息
func (*MessageService) CancelTemplate(opername, userAcc, name, id string) (err error) {
	content := common.MsgContentCancelTempByOwner
	var padding common.MsgTemplatePadding
	padding = common.MsgTemplatePadding{TemplateName: name, OperAccountName: opername, DisabledAccount: userAcc}
	paddingByte, _ := json.Marshal(padding)
	if opername == "" {
		content = common.MsgContentCancelTemp
	}

	if userAcc != "" {
		content = common.MsgContentCancelTempByDisableUser
	}
	var accounts []Account
	ownerAndAdminAccIds := []int{common.AdminAccType, common.OwnerAccType}
	rdb.Find(&accounts, "userType in (?) and isDeleted =0", ownerAndAdminAccIds) //查询所有股东和管理员
	revicer := make([]interface{}, 0, len(accounts))
	for _, acc := range accounts {
		revicer = append(revicer, acc.ID)
	}
	json_revicer, _ := json.Marshal(revicer)

	param := make(map[string]interface{})
	param["id"] = id
	param_json, _ := json.Marshal(param)
	err = rdb.Exec("insert into message (title,content,padding,receiver,type,param) VALUES (?,?,?,?,?,?)", common.MsgTitleCancelTemplate, content, string(paddingByte), string(json_revicer), 1, string(param_json)).Error
	return err
}

//CancelTransfer 作废转账消息
func (*MessageService) CancelTransfer(ApplyerId int, CoinName, Amount, id string, isCalcelUser bool) (err error) {
	content := common.MsgContentTransferCanceled
	if isCalcelUser == true {
		content = common.MsgContentTransferCanceledCauseDisableUser
	}

	revicer := make([]interface{}, 0, 1)
	revicer = append(revicer, ApplyerId)
	json_revicer, _ := json.Marshal(revicer)

	param := make(map[string]interface{})
	param["id"] = id
	param_json, _ := json.Marshal(param)
	padding := &common.MsgTemplatePadding{TransOrderCoinName: CoinName, TransOrderAmount: Amount}
	paddingByte, _ := json.Marshal(padding)
	err = rdb.Exec("insert into message (title,content,padding,receiver,type,param) VALUES (?,?,?,?,?,?)", common.MsgTitleTransferCanceled, content, string(paddingByte), string(json_revicer), 2, string(param_json)).Error
	return err
}

// AddAuth 添加权限消息
func (*MessageService) AddAuth(add, authIds []int, userType int, operator string) (err error) {
	var auths []Auth
	var padding common.MsgTemplatePadding
	if userType == common.AdminAccType {
		padding = common.MsgTemplatePadding{AccountType: common.MsgAdminAccountType, OperAccountName: operator}
	}
	if userType == common.OwnerAccType {
		padding = common.MsgTemplatePadding{AccountType: common.MsgOwnerAccountType, OperAccountName: operator}
	}
	rdb.Find(&auths, "id in (?)", authIds)
	for i := 0; i < len(auths); i++ {
		padding.AuthNames = append(padding.AuthNames, auths[i].AuthType)
	}
	paddingByte, _ := json.Marshal(padding)
	addJson, _ := json.Marshal(add)
	err = rdb.Exec("insert into message (title,content,padding,receiver) VALUES (?,?,?,?)", common.MsgTitleAuthorization, common.MsgContentAuthGrant, string(paddingByte), string(addJson)).Error
	return err
}

// DelAuth 取消权限消息
func (*MessageService) DelAuth(authId, accountId, userType int, account string) (err error) {
	var auth Auth
	var padding common.MsgTemplatePadding
	if userType == common.AdminAccType {
		padding = common.MsgTemplatePadding{AccountType: common.MsgAdminAccountType, OperAccountName: account}
	}
	if userType == common.OwnerAccType {
		padding = common.MsgTemplatePadding{AccountType: common.MsgOwnerAccountType, OperAccountName: account}
	}
	rdb.First(&auth, "id=?", authId)
	padding.AuthNames = append(padding.AuthNames, auth.AuthType)
	revicer := make([]interface{}, 0, 1)
	revicer = append(revicer, accountId)
	ReceiverJson, _ := json.Marshal(revicer)
	paddingByte, _ := json.Marshal(padding)
	err = rdb.Exec("insert into message (title,content,padding,receiver) VALUES (?,?,?,?)", common.MsgTitleAuthorization, common.MsgContentAuthCancel, string(paddingByte), string(ReceiverJson)).Error
	return err
}

// VoucherWarn 签名机警告消息
func (*MessageService) VoucherWarn(accountId, warnNum, pageType int, param_json string) (*Message, error) {
	receiver := make([]interface{}, 0, 1)
	receiver = append(receiver, accountId)
	json_revicer, _ := json.Marshal(receiver)
	padding := &common.MsgTemplatePadding{VoucherWarnCount: warnNum}
	paddingByte, _ := json.Marshal(padding)
	msg := &Message{Title: common.MsgTitleVoucherWarn, Content: common.MsgContentVoucherWarn, Padding: string(paddingByte), Receiver: string(json_revicer), Type: pageType, Param: param_json, CreatedAt: time.Now(), WarnType: common.ErrorMessage}
	err := rdb.Create(msg).Error
	if err != nil {
		return nil, err
	}
	return msg, nil
}

// LetterList 站内信列表
func (*MessageService) LetterList() (msgListAsc []Message, msgListDesc []Message, err error) {
	err = rdb.Order("createdAt desc").Find(&msgListDesc).Error
	err = rdb.Find(&msgListAsc).Error
	return
}

// GetMessageById get message info by id
func (*MessageService) GetMessageById(id int) (Message, error) {
	var message Message
	err := rdb.First(&message, "id=?", id).Error
	return message, err
}

// ReadLetter app read message
func (*MessageService) ReadLetter(id int, readers string) error {
	var message Message
	err := rdb.Model(&message).Where("id = ?", id).Update("reader", readers).Error
	return err
}

func (*MessageService) AysncBlockChainStatus(types int) (err error) {
	content := common.MsgContentNodeSyncFail
	if types == 1 {
		content = common.MsgContentNodeSyncSuccess
	}
	var accounts []Account
	rdb.Find(&accounts, "userType in (?) and isDeleted =0", []int{common.AdminAccType, common.OwnerAccType}) //查询所有股东和管理员
	revicer := make([]interface{}, 0, len(accounts))
	for _, acc := range accounts {
		revicer = append(revicer, acc.ID)
	}
	json_revicer, _ := json.Marshal(revicer)

	err = rdb.Exec("insert into message (title,content,receiver,type) VALUES (?,?,?,?)", common.MsgTitleNodeSync, content, string(json_revicer), 3).Error
	return err
}
