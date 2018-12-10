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
package template

import (
	"encoding/json"
	"errors"
	"sync"
	log "github.com/alecthomas/log4go"
	"github.com/boxproject/apiServer/db"
	reterror "github.com/boxproject/apiServer/errors"
	"github.com/boxproject/apiServer/service/logger"
	msg "github.com/boxproject/apiServer/service/message"
	"github.com/boxproject/apiServer/middleware"
	voucher "github.com/boxproject/apiServer/rpc"
	"github.com/boxproject/apiServer/utils"
	"github.com/jinzhu/gorm"
	"github.com/boxproject/apiServer/service/verify"
	"github.com/boxproject/apiServer/common"
)

var templateLock = new(sync.Mutex)

func CreateNewTemplate(t *AddTmp, claims *middleware.CustomClaims) reterror.ErrModel {
	log.Debug("CreateTemplate...")
	//required parameters
	if t.Content == "" || t.TemplateSign == "" {
		log.Error("参数不能为空")
		return reterror.ErrModel{Code: reterror.ParamsNil, Err: reterror.PARAMS_NULL}
	}
	var templateModel db.TemplateContent
	// parse template content
	err := json.Unmarshal([]byte(t.Content), &templateModel)
	if err != nil {
		log.Error("解析模板内容", err)
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}

	if templateModel.Period < 0 || templateModel.Period > common.TemplateExpireUpperLimit {
		log.Error("冻结周期非法")
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	ts := &db.TemplateDBService{}
	// check if name exists
	templateExists, err := ts.NameExists(templateModel.Name)
	if err != nil {
		log.Error("审批流模板重名", err)
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	if templateExists {
		log.Error("审批流已存在")
		return reterror.ErrModel{Code: reterror.Template_9001, Err: reterror.MSG_9001}
	}

	// validate template rsa sign
	err = utils.RsaVerify([]byte(claims.PubKey), []byte(t.Content), []byte(t.TemplateSign))
	if err != nil {
		log.Error("模板签名错误")
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	// create template
	templateID, err := ts.Create(&templateModel, t.Content, claims.ID)
	if err != nil {
		log.Error("创建审批流模板落库", err)
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	return reterror.ErrModel{Code: reterror.Success, Data: templateID}
}

// TemplateList get template list
func TemplateList(types int, userType int, accountId int) reterror.ErrModel {
	ts := &db.TemplateDBService{}
	templateList, err := ts.TemplateList(types, userType, accountId)
	if err != nil {
		log.Error("模板列表错误", err)
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	} else {
		return reterror.ErrModel{Code: reterror.Success, Data: templateList}
	}

}

// FindTemplateById get template by id
func FindTemplateById(tempVo *TemplateVo, accountId int) reterror.ErrModel {
	if tempVo.Id == "" {
		log.Error("参数不能为空")
		return reterror.ErrModel{Code: reterror.ParamsNil, Err: reterror.PARAMS_NULL}
	}
	ts := &db.TemplateDBService{}
	tems := &db.Template{
		ID: tempVo.Id,
	}
	tem, err := ts.FindTemById(tems, accountId)
	if err != nil {
		log.Error("Find Template Error", err)
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	} else {
		if tem == nil {
			log.Error("template not found")
			return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
		}
		return reterror.ErrModel{Code: reterror.Success, Data: tem}
	}
}

// VerifyTemplate verify template 1.passed 2.rejected
func VerifyTemplate(tempVo *TemplateVo, claims *middleware.CustomClaims) reterror.ErrModel {
	log.Debug("VerifyTemplate...")
	// add lock
	defer templateLock.Unlock()
	templateLock.Lock()
	var temOperApprovalModel = &db.TemplateOper{
		TemplateId: tempVo.Id,
		AccountId:  claims.ID,
		Sign:       tempVo.TemplateSign,
		AppId:      claims.AppID,
		AppName:    claims.Account,
	}
	if tempVo.AdminStatus == common.TemplateRejected {
		if tempVo.Reason == "" {
			log.Error("Reject 需要填原因")
			return reterror.ErrModel{Code: reterror.ParamsNil, Err: reterror.PARAMS_NULL}
		}
		temOperApprovalModel.Reason = tempVo.Reason
	}
	if tempVo.Sign == "" || (tempVo.AdminStatus != common.TemplateApprovaled && tempVo.AdminStatus != common.TemplateRejected) || tempVo.Timestamp == 0 || tempVo.Id == "" || tempVo.TemplateSign == "" || tempVo.Msg == "" {
		log.Error("参数不能为空")
		return reterror.ErrModel{Code: reterror.ParamsNil, Err: reterror.PARAMS_NULL}
	}

	// get template info
	ts := db.TemplateDBService{}
	tem := &db.Template{
		ID: tempVo.Id,
	}
	temp, err := ts.FindTemplateById(tem)

	if err != nil {
		log.Error("Find Template Error", err)
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}

	if temp == nil {
		log.Error("Template not found", tem.ID)
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}

	// validate hash
	t_temp_hash := utils.GenHashStr(tem.Content)
	if t_temp_hash != tem.Hash {
		log.Error("template hash", tempVo.Id)
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}

	// check paramete
	if (tempVo.AdminStatus != common.TemplateApprovaled) && (tempVo.AdminStatus != common.TemplateRejected) {
		log.Error("Not the correct Approval Info", tempVo.AdminStatus)
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	//如果已经拒绝，返回
	if tem.Status == common.TemplateRejected {
		log.Error("已经有股东拒绝")
		return reterror.ErrModel{Code: reterror.Template_9010, Err: reterror.MSG_9010}
	}
	// 可以审批
	if tem.Status != common.TemplatePending {
		log.Error("Can't Apprlval cause the current status is ", tem.Status)
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}

	// avoid double-submit
	tso := &db.TemplateOperDBService{}
	acc_approval_info, err := tso.TemApprovalInfoByAccId(tempVo.Id, claims.ID)
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error("Dumplicate", err)
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	if acc_approval_info.Status != common.TemplatePending {
		log.Debug("Dumplicate")
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}

	// validate rsa sign
	err = utils.RsaVerify([]byte(claims.PubKey), []byte(tem.Content), []byte(tempVo.TemplateSign))
	if err != nil {
		log.Error("Verify Sign", err)
		return reterror.ErrModel{Code: reterror.Code_2, Err: reterror.MSG_2}
	}
	log.Debug("校验签名", err)

	// owner rejected
	if tempVo.AdminStatus == common.TemplateRejected {
		temModel := &db.Template{
			ID:     tempVo.Id,
			Status: common.TemplateRejected,
		}
		temOperApprovalModel.Status = common.TemplateRejected
		rej_result := tso.VerifyTem(temModel, temOperApprovalModel)
		if !rej_result {
			log.Error("reject")
			return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
		}
		// add log
		logErr := logger.AddLog("template", temOperApprovalModel.Reason, claims.Account, common.LoggerTemplateReject, tempVo.Id)
		if logErr != nil {
			log.Error("模板日志添加失败", claims.Account)
		}
		// send message
		msg.FailVerifyTemp(claims.Account, temp.Name, tem.ID, tem.CreatorID)
	} else {
		// 股东同意
		agree := 0
		log.Debug("股东同意")
		temOperApprovalModel.Status = common.TemplateApprovaled
		approval_result := tso.VerifyTem(nil, temOperApprovalModel)
		if !approval_result {
			log.Error("approval")
			return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
		}
		// add log
		logErr := logger.AddLog("template", "", claims.Account, common.LoggerTemplateAccept, tempVo.Id)
		if logErr != nil {
			log.Error("模板日志添加失败", claims.Account)
		}

		// re-acquire template status
		arrOpers := tso.FindTemOperById(tempVo.Id)
		// 所有人的审批意见
		for i := 0; i < len(arrOpers); i++ {
			if arrOpers[i].Status == common.TemplateApprovaled {
				agree++
			}
		}

		// 全部同意，先本地记录状态，不上链
		if agree == len(arrOpers) {
			// send message
			err = msg.SuccessVerifyTemp(temp.Name, temp.ID, temp.CreatorID)
			if err != nil {
				log.Error("插入站内信失败")
			}
			// template progress 0ready to validate 1.passed 2.rejected 3.disale apply 4.add to chain 5.disable success 6.voucher connect failed, re-add to chain
			temFinalModel := &db.Template{
				ID:     tempVo.Id,
				Status: common.TemplateApprovaled,
			}
			// update template status
			final_result := tso.VerifyTem(temFinalModel, nil)
			if !final_result {
				log.Debug("approval template error")
				return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
			}
			db.TemplateChains.SetTemplateChain(*temFinalModel)
			db.ETHTransferCh <- db.TransferSend{Types: 3}

		}
	}
	return reterror.ErrModel{Code: reterror.Success}
}

// TemplateListTransfer get template list
func TemplateListTransfer(coinId int) reterror.ErrModel {
	ts := &db.TemplateDBService{}
	if coinId <= 0 {
		log.Error("参数不能为空")
		return reterror.ErrModel{Code: reterror.ParamsNil, Err: reterror.PARAMS_NULL}
	}
	//update template transfer limit
	ts.UpdateTemplateLimit()
	templateList, err := ts.TemplateListTransfer(coinId)
	if err != nil {
		log.Error("模板列表错误", err)
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	return reterror.ErrModel{Code: reterror.Success, Data: templateList}
}

type FlowSign struct {
	AppId string
	Sign  string
}

type VoucherSigns struct {
	AppID   string
	AppName string
	Status  int
	Sign    string
}

// FlowHashAddToVoucher request template hash add to chain(include availaible and unavailable)
func FlowHashAddToVoucher(flow_content, flow_hash string, arr_temp_opt []db.TemplateOper, data_type int) (*voucher.VoucherStatus, *voucher.GrpcClient, int) {
	// data_type 1-作废 2-同意
	flow_signs := []VoucherSigns{}
	for i := 0; i < len(arr_temp_opt); i++ {
		flow_signs = append(flow_signs, VoucherSigns{arr_temp_opt[i].AppId, arr_temp_opt[i].AppName, arr_temp_opt[i].Status, arr_temp_opt[i].Sign})
	}
	flow_signs_byte, _ := json.Marshal(flow_signs)
	oper_addhash := &voucher.GrpcServer{
		Type:      voucher.VOUCHER_OPERATE_ADDHASH,
		Hash:      flow_hash,
		FlowInfo:  []byte(flow_content),
		FlowSigns: flow_signs_byte,
	}
	if data_type == 1 {
		// disabled
		oper_addhash.Type = voucher.VOUCHER_OPERATE_DISHASH
	}
	return voucher.SendVoucherData(oper_addhash)
}

// CancelTemp disabled template
func CancelTemp(c_temp *ParamCancelTemplate, claims *middleware.CustomClaims) reterror.ErrModel {
	log.Debug("CancelTemplate...")
	// add lock
	defer templateLock.Unlock()
	templateLock.Lock()
	// required parameters
	if c_temp.Sign == "" || c_temp.TemplateId == "" || c_temp.TemplateSign == "" || c_temp.Timestamp == 0 || c_temp.Pwd == "" {
		log.Error("参数不能为空")
		return reterror.ErrModel{Code: reterror.ParamsNil, Err: reterror.PARAMS_NULL}
	}

	// validate password
	verifyResult, _, errorCode := verify.VerifyPSW(claims.Account, c_temp.Pwd)
	if verifyResult.Result == false {
		return reterror.ErrModel{Code: errorCode, Err: errors.New(verifyResult.Reason), Data: map[string]interface{}{"data": verifyResult.Data}}
	}

	// get template content
	ts := db.TemplateDBService{}
	tempp := &db.Template{
		ID: c_temp.TemplateId,
	}
	temp, err := ts.FindTemplateById(tempp)
	if err != nil {
		log.Error("Get template info", err)
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}

	if temp == nil {
		log.Debug("template not found")
		return reterror.ErrModel{Code: reterror.Template_9002, Err: reterror.MSG_9002}
	}

	// validate template content
	if temp.Hash != utils.GenHashStr(temp.Content) {
		log.Error("审批流被篡改")
		return reterror.ErrModel{Code: reterror.Template_9009, Err: reterror.MSG_9009}
	}
	// check if template can be disabled
	if temp.Status != common.TemplateAvaliable {
		log.Error("Template Status Error", temp.Hash, temp.Status)
		return reterror.ErrModel{Code: reterror.Template_9003, Err: reterror.MSG_9003}
	}

	// validate rsa sign
	err = utils.RsaVerify([]byte(claims.PubKey), []byte(temp.Content), []byte(c_temp.TemplateSign))
	if err != nil {
		log.Error("verify flow sign", err)
		return reterror.ErrModel{Code: reterror.Code_2, Err: reterror.MSG_2}
	}

	// disable template and cancel related transfer apply
	var tids, thashs []string
	tids = append(tids, temp.ID)
	thashs = append(thashs, temp.Hash)
	_, err = ts.DisableTemplate(tids, thashs)
	if err != nil {
		log.Error("Disable flow", err)
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	db.TemplateChains.SetTemplateChain(db.Template{ID: temp.ID, Status: common.TemplateToDisable})
	db.ETHTransferCh <- db.TransferSend{Types: 3}

	// add log
	logErr := logger.AddLog("template", "", claims.Account, common.LoggerTemplateCancel, temp.ID)
	if logErr != nil {
		log.Error("模板日志添加失败", claims.Account)
	}
	// send message
	msg.CancelTemplate(claims.Account, "", temp.Name, temp.ID)
	return reterror.ErrModel{Code: reterror.Success}
}

// DisableTemplateByUserName get template by user name
func DisableTemplateByUserName(user *db.Account, operator string) error {
	log.Debug("DisableTemplateByUserName...")
	// get available templates
	var ts db.TemplateDBService
	list, err := ts.TemplistListAvailableByApprover(user.Name)
	if err != nil {
		log.Debug("获取可用的审批流模板", err)
		return err
	}
	log.Debug("template to be disabled", list)
	// cancel all rejected template
	var to_disable_list, disabled_list []db.Template
	for i := 0; i < len(list); i++ {
		if list[i].Status == common.TemplateAvaliable {
			// 已经上链的向签名机申请作废
			to_disable_list = append(to_disable_list, list[i])
			db.TemplateChains.SetTemplateChain(db.Template{ID: list[i].ID, Status: common.TemplateToDisable})
			db.ETHTransferCh <- db.TransferSend{Types: 3}
		} else {
			disabled_list = append(disabled_list, list[i])
		}
	}
	// get hash of templates to be disabled
	to_disabled_hash, err := todisableTemplateList(to_disable_list, user)
	if err != nil {
		log.Debug("作废审批流失败", err)
		return err
	}
	// 作废对应的转账申请
	txToDisableOrders, err := ts.CancelTemplateAdnTxByIDs(to_disabled_hash, common.TemplateToDisable)
	if err != nil {
		log.Error("批量直接作废审批流", err)
		return err
	}

	// get hash of template can be disabled
	disabled_hash, err := disableTemplateList(disabled_list, user)
	if err != nil {
		log.Debug("作废审批流失败", err)
		return err
	}
	txDisablteOrder, err := ts.CancelTemplateAdnTxByIDs(disabled_hash, common.TemplateDisabled)
	if err != nil {
		log.Error("批量预作废审批流", err)
		return err
	}
	//tempsDisables := append(to_disable_list, disabled_list...)
	log.Debug("被作废审批流模板", list)
	// 插入作废审批流站内信
	for _, tem := range list {
		err = msg.CancelTemplate(operator, user.Name, tem.Name, tem.ID)
		if err != nil {
			log.Error("作废审批流失败", tem.ID, err)
		}
	}
	// 插入作废转账站内信
	txOrders := append(txToDisableOrders, txDisablteOrder...)
	trs := db.TransferDBService{}
	for _, orderId := range txOrders {
		txInfo, _ := trs.FindApply(orderId)
		err = msg.CancelTransfer(txInfo.ApplyerId, txInfo.CoinName, txInfo.Amount, orderId, true)
	}
	return nil
}

// userBlongToTemplate return true if approval include user
func userBlongToTemplate(user *db.Account, template_content string) (bool, error) {
	// parse template content
	var content db.TemplateContent
	err := json.Unmarshal([]byte(template_content), &content)
	if err != nil {
		log.Error("Unmarshal Template Content Error", err)
		return false, err
	}

	for i := 0; i < len(content.ApprovalInfo); i++ {
		for j := 0; j < len(content.ApprovalInfo[i].Approvers); j++ {
			if content.ApprovalInfo[i].Approvers[j].Account == user.Name {
				return true, nil
			}
		}
	}
	return false, nil
}

// disableTemplateToVoucher send request to voucher to disable template
func disableTemplateToVoucher(list []db.Template, user *db.Account) ([]string, error) {
	var data []string
	for i := 0; i < len(list); i++ {
		if list[i].Hash != utils.GenHashStr(list[i].Content) {
			log.Error("审批流内容被篡改", list[i].ID)
			return nil, reterror.MSG_9009
		}
		user_belong_to_temp, err := userBlongToTemplate(user, list[i].Content)
		if err != nil {
			log.Error("校验账号是否在审批流中", err)
			return nil, err
		}

		if user_belong_to_temp {
			// send request to voucher
			oper := &voucher.GrpcServer{
				Type: voucher.VOUCHER_OPERATE_DISHASH,
				Hash: list[i].Hash,
			}
			_, _, voucherRet := voucher.SendVoucherData(oper)
			if voucherRet == voucher.VRET_ERR {
				// reset data
				log.Debug("Disable Flow Voucher Return Error", list[i].Hash, voucherRet)
				return nil, errors.New("Voucher Ret Error")
			}
			data = append(data, list[i].Hash)
		}
	}
	return data, nil
}

// todisableTemplateList get templates to be disabled
func todisableTemplateList(list []db.Template, user *db.Account) ([]string, error) {
	var data []string
	for i := 0; i < len(list); i++ {
		if list[i].Hash != utils.GenHashStr(list[i].Content) {
			log.Error("审批流内容被篡改", list[i].ID)
			return nil, reterror.MSG_9009
		}
		user_belong_to_temp, err := userBlongToTemplate(user, list[i].Content)
		if err != nil {
			log.Error("校验账号是否在审批流中", err)
			return nil, err
		}

		if user_belong_to_temp {
			data = append(data, list[i].Hash)
		}
	}
	return data, nil
}

// disableTemplateList get templates can be disabled
func disableTemplateList(list []db.Template, user *db.Account) ([]string, error) {
	var data []string
	for i := 0; i < len(list); i++ {
		if list[i].Hash != utils.GenHashStr(list[i].Content) {
			log.Error("审批流内容被篡改", list[i].ID)
			return nil, reterror.MSG_9009
		}
		user_belong_to_temp, err := userBlongToTemplate(user, list[i].Content)
		if err != nil {
			log.Error("校验账号是否在审批流中", err)
			return nil, err
		}

		if user_belong_to_temp {
			data = append(data, list[i].Hash)
		}
	}
	return data, nil
}

// TempStatistics get template include account
func TempStatistics(account string) reterror.ErrModel {
	if account == "" {
		return reterror.ErrModel{Code: reterror.ParamsNil, Err: reterror.PARAMS_NULL}
	}
	dt := &db.TemplateDBService{}
	tempCount, transCount, err := dt.TempStatistics(account)
	if err != nil {
		log.Error("TempStatistics error", err)
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	return reterror.ErrModel{Err: err, Code: reterror.Success, Data: map[string]interface{}{"TemplateCount": tempCount, "TransferCount": transCount}}
}

// TxNumByTempID get transfer amount by template id
func TxNumByTempID(tempID string) reterror.ErrModel {
	if tempID == "" {
		return reterror.ErrModel{Code: reterror.ParamsNil, Err: reterror.PARAMS_NULL}
	}
	dt := &db.TemplateDBService{}
	transCount, err := dt.CountTxByTempID(tempID)
	if err != nil {
		log.Debug("统计转账数目", err)
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	return reterror.ErrModel{Err: err, Code: reterror.Success, Data: map[string]interface{}{"TransferCount": transCount}}
}
