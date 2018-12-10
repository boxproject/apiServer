// Copyright 2018. bolaxy.org authors.
// 
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// 
// 		 http://www.apache.org/licenses/LICENSE-2.0
// 
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package timer

import (
	"github.com/boxproject/apiServer/db"
	voucher "github.com/boxproject/apiServer/rpc"
	log "github.com/alecthomas/log4go"
	"github.com/boxproject/apiServer/service/template"
	"time"
	"github.com/boxproject/apiServer/common"
	"github.com/boxproject/apiServer/service/message"
	"github.com/boxproject/apiServer/service/logger"
)

// 上链
func uploadTemplate(tem db.Template) {
	var tmp db.TemplateDBService
	var tso db.TemplateOperDBService

	err := tmp.FindTemlateById(&tem)
	if err!=nil {
		log.Error("审批流查询错误",err)
	}
	arrOpers := tso.FindTemOperById(tem.ID)
	var data_type int
	if tem.Status == 3 {
		// 作废
		data_type = 1
	} else if tem.Status == 1 {
		// 同意
		data_type = 2
	}
	// 请求签名机上链
	template.FlowHashAddToVoucher(tem.Content, tem.Hash, arrOpers, data_type)
	temTimer := time.NewTicker(5 * time.Second)
	ret := -1
	for ret == -1{
		select{
			case <- temTimer.C:
				ret = getFlowHashStatusFromVoucher(&tem,arrOpers)
		}
	}
	temTimer.Stop()
	tmp.UpdateTemplateStatus(tem.Hash, ret)
}


// 获取审批流模板哈希状态
//func parseHashStatus() {
//	// 获取中间状态的审批流模板列表
//	var ts db.TemplateDBService
//	//var flow_now_status int
//	list, err := ts.IntermediateStateHashList()
//	if err != nil {
//		log.Error("Get IntermediateStateHashList Error", err)
//	}
//	for i := 0; i < len(list); i++ {
//		// 从签名机获取对应哈希的状态
//		fmt.Println("模板哈希状态", list[i].Hash, list[i].Status)
//		template_status := getFlowHashStatusFromVoucher(&list[i])
//		// 更新审批流模板状态
//		if template_status == db.Template_Disabled || template_status == db.Template_Avaliable {
//			ts.UpdateTemplateStatus(list[i].Hash, template_status)
//			if template_status == db.Template_Avaliable {
//				// 审批流生效日志
//				logErr := logger.AddLog("template", "", "", "templateActive", list[i].ID)
//				if logErr != nil {
//					log.Error("定时器模板日志添加失败", list[i].ID)
//				}
//			}
//		}
//	}
//}

// 从签名机获取审批流哈希状态
func getFlowHashStatusFromVoucher(tempInfo *db.Template,arr_temp_opt []db.TemplateOper) (status int) {
	data_type := 2
	if tempInfo.Status == common.TemplateToDisable {
		data_type = 1
	}
	_, voucherRes, voucherRet := template.FlowHashAddToVoucher(tempInfo.Content, tempInfo.Hash, arr_temp_opt, data_type)
	if voucherRet == voucher.VRET_CLIENT {
		log.Debug("voucher return status", voucherRes.Status)
		if voucherRes.Status == voucher.STATUS_OK {
			if tempInfo.Status == common.TemplateToDisable {
				return common.TemplateDisabled
			} else if tempInfo.Status == common.TemplateApprovaled {
				// 添加日志
				logger.AddLog("template", "", "", common.LoggerTemplateActive, tempInfo.ID)
				return common.TemplateAvaliable
			} else {
				return tempInfo.Status
			}
		} else if voucherRes.Status == voucher.STATUS_INSUFFICIENT_FUNDS || voucherRes.Status == voucher.STATUS_FLOW_VERIFY_ERROR || voucherRes.Status == voucher.STATUS_HASH_LOSEBASEINFO {
			if voucherRes.Status == voucher.STATUS_FLOW_VERIFY_ERROR || voucherRes.Status == voucher.STATUS_HASH_LOSEBASEINFO {
				// 发送站内信
				err := message.VoucherFail(message.VoucherRetTempFailMsgType, voucherRes.Status, tempInfo, nil)
				log.Error("发送站内信失败", err)
				return common.TemplateFailed
			}
			log.Error("voucher return error", voucherRes.Status)
			return common.TemplateFailed
		}
	}
	log.Debug("apply voucher", voucherRet)
	return -1
}


