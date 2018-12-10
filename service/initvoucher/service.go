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
package initvoucher

import (
	"encoding/json"
	"strconv"
	"sync"
	"time"

	log "github.com/alecthomas/log4go"
	db "github.com/boxproject/apiServer/db"
	reterror "github.com/boxproject/apiServer/errors"
	JWT "github.com/boxproject/apiServer/middleware"
	voucher "github.com/boxproject/apiServer/rpc"
	"github.com/boxproject/apiServer/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/boxproject/apiServer/service/logger"
	"github.com/boxproject/apiServer/service/verify"
	"github.com/boxproject/apiServer/common"
)

const TOKEN_EXP = 24

var initVoucherLock = new(sync.Mutex)

// QrCommit get code from
func QrCommit(voucherVo *VoucherVo) reterror.ErrModel {
	log.Debug("QrCommit...")
	initVoucherLock.Lock()
	defer initVoucherLock.Unlock()
	// 判断注册表是否有数据并添加数据
	da := &db.RegDBService{}
	duplicate, err := da.DuplicateName(voucherVo.AppName)
	if err != nil {
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	if duplicate {
		return reterror.ErrModel{Code: reterror.User_2002, Err: reterror.MSG_2002}
	}
	oper := &voucher.GrpcServer{
		Type:      voucher.VOUCHER_OPERATE_ADDKEY,
		AppId:     voucherVo.AppId,
		AppName:   voucherVo.AppName,
		Sign:      []byte(voucherVo.Sign),
		Timestamp: voucherVo.Timestamp,
		AesKey:    []byte(voucherVo.AesKey),
		Msg:       []byte(voucherVo.Msg),
	}
	_, res, ret := voucher.SendVoucherData(oper)
	if ret == voucher.VRET_ERR {
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	if res.Status != voucher.STATUS_OK {
		log.Debug("签名机状态错误", res.Status)
		return reterror.ErrModel{Code: reterror.VoucherFail, Err: reterror.MSG_1006}
	}
	reg := &db.Registration{
		Name:  voucherVo.AppName,
		AppId: voucherVo.AppId,
	}
	regErr := da.AddUser(reg)
	if regErr != nil {
		log.Debug("添加注册表失败", regErr)
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	return reterror.ErrModel{Err: nil, Code: reterror.Success, Data: res}
}

// CommandCommit send keywords and backup password
func CommandCommit(voucherVo *VoucherVo) reterror.ErrModel {
	log.Debug("CommandCommit...")
	initVoucherLock.Lock()
	defer initVoucherLock.Unlock()
	oper := &voucher.GrpcServer{
		Type:      voucher.VOUCHER_OPERATE_CREATE,
		AppId:     voucherVo.AppId,
		AppName:   voucherVo.AppName,
		AesKey:    []byte(voucherVo.AesKey),
		Msg:       []byte(voucherVo.Msg),
		Sign:      []byte(voucherVo.Sign),
		Timestamp: voucherVo.Timestamp,
	}
	log.Error("commit", voucherVo.AesKey)
	_, res, ret := voucher.SendVoucherData(oper)
	if ret == voucher.VRET_ERR {
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	return reterror.ErrModel{Err: nil, Code: reterror.Success, Data: res}
}

// GetOtherPubKey get other owners public key
func GetOtherPubKey(voucherVo *VoucherVo) reterror.ErrModel {
	log.Debug("GetOtherPubKey...")
	initVoucherLock.Lock()
	defer initVoucherLock.Unlock()
	acc := &db.AccDBService{}
	account, accErr := acc.AccountByName(voucherVo.AppName)
	if accErr != nil {
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	if len(account) != 0 {
		return reterror.ErrModel{Code: reterror.User_2002, Err: reterror.MSG_2002}
	}
	oper := &voucher.GrpcServer{
		Type:      voucher.VOUCHER_OPERATE_APP_PUBKEY,
		AppId:     voucherVo.AppId,
		AppName:   voucherVo.AppName,
		AesKey:    []byte(voucherVo.AesKey),
		Msg:       []byte(voucherVo.Msg),
		Sign:      []byte(voucherVo.Sign),
		Timestamp: voucherVo.Timestamp,
	}
	_, res, ret := voucher.SendVoucherData(oper)
	if ret == voucher.VRET_ERR {
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	} else {
		if res.Status != voucher.STATUS_OK {
			log.Debug("send to voucher failed", res.Status)
			return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR, Data: res}
		}
		salt := strconv.FormatInt((time.Now().UnixNano() / 1e6), 10)
		pwdSalt := utils.PwdSalt(voucherVo.Password, salt)
		var otherMapInfo = make(map[string][]byte)
		jsonErr := json.Unmarshal(res.Other, &otherMapInfo)
		if jsonErr != nil {
			log.Error("签名机app公钥反解错误", jsonErr)
			return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
		}

		account := &db.Account{
			Name:         voucherVo.AppName,
			Pwd:          pwdSalt,
			Salt:         salt,
			AppId:        voucherVo.AppId,
			PubKey:       string(otherMapInfo["appPubKey"][:]),
			Msg:          string(otherMapInfo["voucherSign"][:]),
			Level:        1,
			UserType:     common.OwnerAccType,
			DepartmentId: 1,
			FrozenTo:     time.Now(),
			SourceAppId:  "0",
		}
		acc := &db.AccDBService{}
		err := acc.AddAccount(account)
		if err != nil {
			log.Error("保存股东账号错误", err)
			return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
		}
		res.Other = nil
		return reterror.ErrModel{Err: nil, Code: reterror.Success, Data: res}
	}
}

// BackupKey backup key
func BackupKey(voucherVo *VoucherVo) reterror.ErrModel {
	log.Debug("BackupKey...")
	initVoucherLock.Lock()
	defer initVoucherLock.Unlock()
	oper := &voucher.GrpcServer{
		Type:      voucher.VOUCHER_OPERATE_HDBAKUP,
		AppId:     voucherVo.AppId,
		AppName:   voucherVo.AppName,
		AesKey:    []byte(voucherVo.AesKey),
		Msg:       []byte(voucherVo.Msg),
		Sign:      []byte(voucherVo.Sign),
		Timestamp: voucherVo.Timestamp,
		BakAction: voucherVo.BakAction,
	}
	_, res, ret := voucher.SendVoucherData(oper)
	if ret == voucher.VRET_ERR {
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	return reterror.ErrModel{Err: nil, Code: reterror.Success, Data: res}
}

// StartVoucher start voucher
func StartVoucher(voucherVo *VoucherVo) reterror.ErrModel {
	log.Debug("StartVoucher...")
	initVoucherLock.Lock()
	defer initVoucherLock.Unlock()
	// check if account isFrozen
	verifyErr := verify.IsFrozen(voucherVo.AppName)
	if verifyErr.Code != 0 {
		return verifyErr
	}
	oper := &voucher.GrpcServer{
		Type:      voucher.VOUCHER_OPERATE_START,
		AppId:     voucherVo.AppId,
		AppName:   voucherVo.AppName,
		AesKey:    []byte(voucherVo.AesKey),
		Msg:       []byte(voucherVo.Msg),
		Sign:      []byte(voucherVo.Sign),
		Timestamp: voucherVo.Timestamp,
	}
	_, res, ret := voucher.SendVoucherData(oper)
	// logs
	var logEntity []voucher.OperateLog
	err := json.Unmarshal(res.Logs, &logEntity)
	if err != nil {
		log.Error("log json转化失败", err)
	}
	logErr := logger.AddLog("voucher", "", voucherVo.AppName, common.LoggerVoucherKeyWord, "")
	if logErr != nil {
		log.Info("启动签名机日志", logErr)
	}

	if ret == voucher.VRET_ERR {
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	// logs
	if logEntity[0].Other == "voucher start success" {
		logErr := logger.AddLog("voucher", "", "", common.LoggerVoucherStart, "")
		if logErr != nil {
			log.Info("启动签名机日志（最后）:", logErr)
		}
	}

	switch res.Status {
	case voucher.STATUS_APP_VERIFY_ERROR:
		return reterror.ErrModel{Code: reterror.Code_102, Err: reterror.MSG_102}
	case voucher.STATUS_APP_PASSWORD_ERROR:
		keywordErr := verify.Keyword(voucherVo.AppId, voucherVo.AppName)
		return keywordErr
	case voucher.STATUS_APP_NAME_NOTMATCH:
		return reterror.ErrModel{Code: reterror.Code_106, Err: reterror.MSG_106}
	case voucher.STATUS_APP_SAMENAME:
		return reterror.ErrModel{Code: reterror.Code_107, Err: reterror.MSG_107}
	case voucher.VRET_TIMEOUT:
		voucherLog := logger.AddLog("voucher", "", "", common.LoggerVoucherTimeout, "")
		if voucherLog != nil {
			log.Info("签名机超时日志", voucherLog)
		}
		return reterror.ErrModel{Code: reterror.VoucherTimeout, Err: reterror.MSG_Voucher_Timeout}
	}
	ds := db.AccDBService{}
	_, accErr := ds.ResetAccount(voucherVo.AppId)
	if accErr != nil {
		log.Info("取消账户冻结错误", accErr)
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	return reterror.ErrModel{Err: nil, Code: reterror.Success, Data: res}
}

// StopVoucher stop voucher
func StopVoucher(voucherVo *VoucherVo) reterror.ErrModel {
	initVoucherLock.Lock()
	defer initVoucherLock.Unlock()
	// check if account isFrozen
	verifyErr := verify.IsFrozen(voucherVo.AppName)
	if verifyErr.Code != 0 {
		return verifyErr
	}

	oper := &voucher.GrpcServer{
		Type:      voucher.VOUCHER_OPERATE_STOP,
		AppId:     voucherVo.AppId,
		AppName:   voucherVo.AppName,
		AesKey:    []byte(voucherVo.AesKey),
		Msg:       []byte(voucherVo.Msg),
		Sign:      []byte(voucherVo.Sign),
		Timestamp: voucherVo.Timestamp,
	}
	_, res, ret := voucher.SendVoucherData(oper)
	if ret == voucher.VRET_ERR {
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}

	switch res.Status {
	case voucher.STATUS_APP_VERIFY_ERROR:
		return reterror.ErrModel{Code: reterror.Code_102, Err: reterror.MSG_102}
	case voucher.STATUS_APP_PASSWORD_ERROR:
		keywordErr := verify.Keyword(voucherVo.AppId, voucherVo.AppName)
		return keywordErr
	case voucher.STATUS_APP_NAME_NOTMATCH:
		return reterror.ErrModel{Code: reterror.Code_106, Err: reterror.MSG_106}
	case voucher.STATUS_APP_SAMENAME:
		return reterror.ErrModel{Code: reterror.Code_107, Err: reterror.MSG_107}
	case voucher.VRET_TIMEOUT:
		voucherLog := logger.AddLog("voucher", "", "", common.LoggerVoucherTimeout, "")
		if voucherLog != nil {
			log.Info("签名机超时日志", voucherLog)
		}
		return reterror.ErrModel{Code: reterror.VoucherTimeout, Err: reterror.MSG_Voucher_Timeout}
	}

	logErr := logger.AddLog("voucher", "", voucherVo.AppName, common.LoggerVoucherStop, "")
	if logErr != nil {
		log.Info("关停签名机日志", logErr)
	}
	ds := db.AccDBService{}
	_, accErr := ds.ResetAccount(voucherVo.AppId)
	if accErr != nil {
		log.Info("取消账户冻结错误", accErr)
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	return reterror.ErrModel{Err: nil, Code: reterror.Success, Data: res}
}

// GetStatus get voucher status
func GetStatus(voucherVo *VoucherVo) reterror.ErrModel {
	oper := &voucher.GrpcServer{
		Type:      voucher.VOUCHER_STATUS,
		Timestamp: time.Now().UnixNano(),
	}
	voucherStatus, _, ret := voucher.SendVoucherData(oper)
	if ret == voucher.VRET_ERR {
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	} else if ret == voucher.VRET_TIMEOUT {
		return reterror.ErrModel{Code: reterror.VoucherTimeout, Err: reterror.MSG_Voucher_Timeout}
	} else {
		return reterror.ErrModel{Err: nil, Code: reterror.Success, Data: voucherStatus}
	}
}

// GetConnectStatus get voucher connect status
// 0.未连接 1.未创建 2.已创建 3.已备份 4.已启动 5.已停止
func GetConnectStatus() reterror.ErrModel {
	fail := 0
	voucherHeart := voucher.GetVoucherStatus()
	voucherstatus := voucherHeart.ServerStatus
	startPerson := voucherHeart.NodesAuthorized["start"]
	if voucherstatus == voucher.VOUCHER_STATUS_PAUSED {
		//如果有人点启动了
		if len(startPerson) > 0 {
			return reterror.ErrModel{Err: nil, Code: reterror.Success, Data: map[string]interface{}{"status": voucher.VOUCHER_STATUS_BAKUP, "start": startPerson}}
		} else {
			return reterror.ErrModel{Err: nil, Code: reterror.Success, Data: map[string]interface{}{"status": voucher.VOUCHER_STATUS_STARTED}}
		}
	}
	arr := voucher.HeartStatus
	if len(arr) < 5 {
		return reterror.ErrModel{Err: nil, Code: reterror.Success, Data: map[string]interface{}{"status": voucher.VOUCHER_STATUS_UNCONNETED, "start": startPerson}}
	}
	arr = arr[len(arr)-5 : len(arr)]
	for i := 0; i < len(arr); i++ {
		if arr[i] == 0 {
			fail++
		}
	}
	if fail >= 3 {
		return reterror.ErrModel{Err: nil, Code: reterror.Success, Data: map[string]interface{}{"status": voucher.VOUCHER_STATUS_CREATED, "start": startPerson}}
	}
	if (fail >= 1) && (fail < 3) {
		return reterror.ErrModel{Err: nil, Code: reterror.Success, Data: map[string]interface{}{"status": voucher.VOUCHER_STATUS_UNCREATED, "start": startPerson}}
	}
	return reterror.ErrModel{Err: nil, Code: reterror.Success, Data: map[string]interface{}{"status": voucher.VOUCHER_STATUS_UNCONNETED, "start": startPerson}}

}

// SavePubkeySign save other public key
func SavePubkeySign(voucherVo *VoucherVo) reterror.ErrModel {
	log.Debug("SavePubkeySign...")
	initVoucherLock.Lock()
	defer initVoucherLock.Unlock()
	//TODO 没有验证
	acc := &db.AccDBService{}
	account := &db.Account{
		AppId: voucherVo.AppId,
		Msg:   voucherVo.Sign,
	}
	err := acc.SavePubkeySign(account)
	if err == nil {
		log.Debug("SavePubkeySign error", err)
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	return reterror.ErrModel{Err: nil, Code: reterror.Success}

}

// GetToken owner get token
func GetToken(voucherVo *VoucherVo) reterror.ErrModel {
	//TODO 没有验证
	// acc := &db.AccDBService{}
	// account := &db.Account{
	// 	AppId:  voucherVo.AppId,
	// }
	// err := acc.GetToken(account)
	// 生成token
	var j *JWT.JWT = &JWT.JWT{
		[]byte(JWT.GetSignKey()),
	}
	claims := JWT.CustomClaims{AppID: voucherVo.AppId, Account: voucherVo.AppName, UserType: common.OwnerAccType, StandardClaims: jwt.StandardClaims{ExpiresAt: time.Now().Add((TOKEN_EXP) * time.Hour).Unix()}}
	token, err := j.CreateToken(claims)
	if err != nil {
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	return reterror.ErrModel{Err: nil, Code: reterror.Success, Data: map[string]interface{}{
		"token":    token,
		"name":     voucherVo.AppName,
		"userType": common.OwnerAccType}}
}

//GetMasterAddress get master address
func GetMasterAddress() reterror.ErrModel {
	log.Debug("GetMasterAddress...")
	masterAddress := db.AllMasterAddress
	var maps []map[string]string
	if masterAddress == nil {
		log.Error("未获取到主地址")
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	for k, v := range masterAddress {
		value := make(map[string]string)
		value["CoinName"] = k
		value["Address"] = v
		maps = append(maps, value)
	}
	return reterror.ErrModel{Err: nil, Code: reterror.Success, Data: maps}
}

// SaveMasterAddress save master address
func SaveMasterAddress(value string) reterror.ErrModel {
	txinfosService := &db.TxinfosService{}
	if db.ConfigMap["address_sign"] != "" {
		log.Debug("sign data cant bee null.")
		return reterror.ErrModel{Code: reterror.Success, Err: nil}
	}
	err := txinfosService.UpadateConfigs("address_sign", value)
	if err != nil {
		log.Error("保存错误")
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	db.ConfigMap["address_sign"] = value
	return reterror.ErrModel{Err: nil, Code: reterror.Success}
}

// OperateKeyExchange exchange public key
func OperateKeyExchange(voucherVo *VoucherVo) reterror.ErrModel {
	log.Debug("OperateKeyExchange...")
	initVoucherLock.Lock()
	defer initVoucherLock.Unlock()
	voucherStatus := voucher.GetVoucherStatus()
	applyer := voucherStatus.KeyStoreStatus[voucherVo.AppId].ApplyerId
	if applyer != "" && voucherVo.IsRecover == 1 {
		return reterror.ErrModel{Code: reterror.User_3018, Err: reterror.MSG_3018}
	}
	oper := &voucher.GrpcServer{
		Type:      voucher.VOUCHER_OPERATE_KEY_EXCHANGE,
		AppId:     voucherVo.AppId,
		AppName:   voucherVo.AppName,
		AesKey:    []byte(voucherVo.AesKey),
		Msg:       []byte(voucherVo.Msg),
		Sign:      []byte(voucherVo.Sign),
		Timestamp: voucherVo.Timestamp,
	}
	_, res, ret := voucher.SendVoucherData(oper)
	if ret == voucher.VRET_ERR {
		log.Error("voucher return error", ret)
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	return reterror.ErrModel{Err: nil, Code: reterror.Success, Data: res}
}
