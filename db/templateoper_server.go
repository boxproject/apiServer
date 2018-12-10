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
	"github.com/jinzhu/gorm"
	"github.com/boxproject/apiServer/common"
)

type TemplateOperDBService struct {
}

// VerifyTem 审核模板
func (*TemplateOperDBService) VerifyTem(template *Template, templateOper *TemplateOper) bool {
	conn := rdb.Begin()
	if templateOper != nil {
		var err error
		if templateOper.Status == common.TemplateRejected {
			err = conn.Model(&TemplateOper{}).Where("templateId = ? and accountId = ?", templateOper.TemplateId, templateOper.AccountId).Updates(map[string]interface{}{"status": templateOper.Status, "sign": templateOper.Sign, "appId": templateOper.AppId, "appName": templateOper.AppName, "reason": templateOper.Reason}).Error
		} else {
			err = conn.Model(&TemplateOper{}).Where("templateId = ? and accountId = ?", templateOper.TemplateId, templateOper.AccountId).Updates(map[string]interface{}{"status": templateOper.Status, "sign": templateOper.Sign, "appId": templateOper.AppId, "appName": templateOper.AppName}).Error
		}
		if err != nil {
			log.Error("股东审核错误", err)
			conn.Rollback()
			return false
		}
	}
	if template != nil {
		err := conn.Model(template).Update(template).Error
		if err != nil {
			log.Error("股东审核错误", err)
			conn.Rollback()
			return false
		}
	}
	conn.Commit()
	return true
}

// FindTemOperById 通过模板id查询所有股东的操作
func (*TemplateOperDBService) FindTemOperById(TemplateId string) []TemplateOper {
	var tempOpers []TemplateOper
	err := rdb.Where(&TemplateOper{TemplateId: TemplateId}).Find(&tempOpers).Error
	if err != nil {
		log.Error("查询列表错误", err)
		return nil
	} else {
		return tempOpers
	}
}

// TemApprovalInfoByAccId 是否重复提交审批模板意见
func (*TemplateOperDBService) TemApprovalInfoByAccId(templateID string, accID int) (*TemplateOper, error) {
	var optemp TemplateOper
	err := rdb.Where(&TemplateOper{TemplateId: templateID, AccountId: accID}).First(&optemp).Error
	return &optemp, err
}

//// DisableTemplateWhenDisableUser
//func(*TemplateDBService)DisableTemplateWhenDisableUser(tIds, tHashs []string ) ([]string, error) {
//	// 获取审批中的转账申请
//	var txOrders []TransferOrder
//	var txOrderIds []string
//	err := rdb.Where(&Template{Status:common.TransferToApproval}).Where("hash in (?)", tIds).Where("status").Find(&txOrders).Pluck("id", &txOrderIds).Error
//	if err != nil && err != gorm.ErrRecordNotFound {
//		log.Error("获取审批中的转账申请")
//		return nil, err
//	}
//	conn := rdb.Begin()
//	// 作废审批流模板
//	err = conn.Model(&Template{}).Where("id in (?)", tIds).Update(&Template{Status: common.TemplateToDisable}).Error
//	if err != nil {
//		log.Error("作废审批流模板", err)
//		conn.Rollback()
//		return nil, err
//	}
//
//	// 作废大订单
//	err = conn.Model(&TransferOrder{}).Where("hash in (?)", tHashs).Updates(&TransferOrder{Status:common.TransferDisableEmployee}).Error
//	if err != nil {
//		log.Error("作废相关转账申请", err)
//		conn.Rollback()
//		return nil, err
//	}
//
//	// 作废子订单
//	err = conn.Model(&Transfer{}).Where("orderId in (?)", txOrderIds).Updates(&Transfer{Status:common.TransferDisableEmployee}).Error
//	if err != nil {
//		log.Error("作废子订单", err)
//		conn.Rollback()
//		return nil, err
//	}
//	conn.Commit()
//	return txOrderIds, nil
//}

// DisableTemplate disable template by given ids
func (*TemplateDBService) DisableTemplate(tIds, tHashs []string) ([]string, error) {
	// 获取审批中的转账申请
	var txOrders []TransferOrder
	var txOrderIds []string
	err := rdb.Where("status = ? and hash in (?)", common.TransferToApproval, tHashs).Find(&txOrders).Pluck("id", &txOrderIds).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error("获取审批中的转账申请")
		return nil, err
	}
	conn := rdb.Begin()
	// 作废审批流模板
	err = conn.Model(&Template{}).Where("id in (?)", tIds).Update(&Template{Status: common.TemplateToDisable}).Error
	if err != nil {
		log.Error("作废审批流模板", err)
		conn.Rollback()
		return nil, err
	}

	// 作废大订单
	// 0.审批中 1.审批通过，转账中 2.拒绝 3.部分转账成功 4.转账失败 5.全部成功 6.撤回 7.非法 8.审批过期 9.转账过期 10.员工作废 11.模板停用作废
	err = conn.Model(&TransferOrder{}).Where("status = ? and hash in (?)", common.TransferToApproval, tHashs).Updates(&TransferOrder{Status: common.TransferCancelTemp}).Error
	if err != nil {
		log.Error("作废相关转账申请", err)
		conn.Rollback()
		return nil, err
	}

	// 作废子订单
	err = conn.Model(&Transfer{}).Where("orderId in (?) and status = ? ", txOrderIds, common.TransferToApproval).Updates(&Transfer{Status: common.TransferCancelTemp}).Error
	if err != nil {
		log.Error("作废子订单", err)
		conn.Rollback()
		return nil, err
	}
	conn.Commit()
	return txOrderIds, nil
}
