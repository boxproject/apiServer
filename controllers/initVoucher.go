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
package controllers

import (
	initVoucher "github.com/boxproject/apiServer/service/initvoucher"
	"github.com/boxproject/apiServer/errors"
	log "github.com/alecthomas/log4go"
	"github.com/gin-gonic/gin"
	"encoding/json"
)

// GetStatus shows the status of the voucher
func GetStatus(ctx *gin.Context) {
	var vo initVoucher.VoucherVo
	var msg errors.ErrModel
	if err := ctx.ShouldBind(&vo); err != nil {
		log.Error("Bind Param Error", err)
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	} else {
		j_vo, _ := json.Marshal(vo)
		log.Debug("GetStatus...", string(j_vo))
		msg = initVoucher.GetStatus(&vo)
	}
	msg.RetErr(ctx)
}

// QrCommit used to confirm the registration
func QrCommit(ctx *gin.Context) {
	var vo initVoucher.VoucherVo
	var msg errors.ErrModel
	if err := ctx.ShouldBind(&vo); err != nil {
		log.Error("Bind Param Error", err)
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	} else {
		j_vo, _ := json.Marshal(vo)
		log.Debug("QrCommit...", string(j_vo))
		msg = initVoucher.QrCommit(&vo)
	}
	msg.RetErr(ctx)
}

// CommandCommit used to submit the key words to voucher
func CommandCommit(ctx *gin.Context) {
	var vo initVoucher.VoucherVo
	var msg errors.ErrModel
	if err := ctx.ShouldBind(&vo); err != nil {
		log.Error("Bind Param Error", err)
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	} else {
		j_vo, _ := json.Marshal(vo)
		log.Debug("提交备份口令", string(j_vo))
		msg = initVoucher.CommandCommit(&vo)
	}
	msg.RetErr(ctx)
}

// GetOtherPubKey used to get other owner's pubkeys
func GetOtherPubKey(ctx *gin.Context) {
	var vo initVoucher.VoucherVo
	var msg errors.ErrModel
	if err := ctx.ShouldBind(&vo); err != nil {
		log.Error("Bind Param Error", err)
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	} else {
		j_vo, _ := json.Marshal(vo)
		log.Debug("各个私钥app同步其他公钥", string(j_vo))
		msg = initVoucher.GetOtherPubKey(&vo)
	}
	msg.RetErr(ctx)
}

// BackupKey used to backup the key words to BOX KEY
func BackupKey(ctx *gin.Context) {
	var vo initVoucher.VoucherVo
	var msg errors.ErrModel
	if err := ctx.ShouldBind(&vo); err != nil {
		log.Error("Bind Param Error", err)
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	} else {
		j_vo, _ := json.Marshal(vo)
		log.Debug("备份关键句", string(j_vo))
		msg = initVoucher.BackupKey(&vo)
	}
	msg.RetErr(ctx)
}

// StartVoucher used to start the service of voucher
func StartVoucher(ctx *gin.Context) {
	var vo initVoucher.VoucherVo
	var msg errors.ErrModel
	if err := ctx.ShouldBind(&vo); err != nil {
		log.Error("Bind Param Error", err)
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	} else {
		j_vo, _ := json.Marshal(vo)
		log.Debug("启动签名机", string(j_vo))
		msg = initVoucher.StartVoucher(&vo)
	}
	msg.RetErr(ctx)
}

// StopVoucher used to stop the service of voucher
func StopVoucher(ctx *gin.Context) {
	var vo initVoucher.VoucherVo
	var msg errors.ErrModel

	if err := ctx.ShouldBind(&vo); err != nil {
		log.Error("Bind Param Error", err)
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	} else {
		j_vo, _ := json.Marshal(vo)
		log.Debug("关闭签名机", string(j_vo))
		msg = initVoucher.StopVoucher(&vo)
	}
	msg.RetErr(ctx)
}

// SavePubkeySign save the sign info of pubkey
func SavePubkeySign(ctx *gin.Context) {
	var vo initVoucher.VoucherVo
	var msg errors.ErrModel
	if err := ctx.ShouldBind(&vo); err != nil {
		log.Error("Bind Param Error", err)
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	} else {
		j_vo, _ := json.Marshal(vo)
		log.Debug("对其他股东app签名保存", string(j_vo))
		msg = initVoucher.SavePubkeySign(&vo)
	}
	msg.RetErr(ctx)
}

// GetToken owner get token
func GetToken(ctx *gin.Context) {
	var vo initVoucher.VoucherVo
	var msg errors.ErrModel
	if err := ctx.ShouldBind(&vo); err != nil {
		log.Error("Bind Param Error", err)
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	} else {
		j_vo, _ := json.Marshal(vo)
		log.Debug("获取token", string(j_vo))
		msg = initVoucher.GetToken(&vo)
	}
	msg.RetErr(ctx)
}

// GetConnectStatus get the connection status of voucher
func GetConnectStatus(ctx *gin.Context) {
	var msg errors.ErrModel
	msg = initVoucher.GetConnectStatus()
	msg.RetErr(ctx)
}

// GetMasterAddress get the master address
func GetMasterAddress(ctx *gin.Context) {
	var msg errors.ErrModel
	msg = initVoucher.GetMasterAddress()
	msg.RetErr(ctx)
}

// SaveMasterAddress save the signed master address form owner
func SaveMasterAddress(ctx *gin.Context) {
	var vo initVoucher.VoucherVo
	var msg errors.ErrModel
	if err := ctx.ShouldBind(&vo); err != nil {
		log.Error("Bind Param Error", err)
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	} else {
		j_vo, _ := json.Marshal(vo)
		log.Debug("保存对主地址的签名", string(j_vo))
		msg = initVoucher.SaveMasterAddress(vo.Msg)
	}
	msg.RetErr(ctx)
}

// OperateKeyExchange exchange the key to voucher
func OperateKeyExchange(ctx *gin.Context) {
	var vo initVoucher.VoucherVo
	var msg errors.ErrModel
	if err := ctx.ShouldBind(&vo); err != nil {
		log.Error("Bind Param Error", err)
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	} else {
		j_vo, _ := json.Marshal(vo)
		log.Debug("交换秘钥", string(j_vo))
		msg = initVoucher.OperateKeyExchange(&vo)
	}
	msg.RetErr(ctx)
}
