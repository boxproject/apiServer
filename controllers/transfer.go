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
	"encoding/json"
	log "github.com/alecthomas/log4go"
	coinService "github.com/boxproject/apiServer/service/coin"
	"github.com/boxproject/apiServer/errors"
	middle "github.com/boxproject/apiServer/middleware"
	templateService "github.com/boxproject/apiServer/service/template"
	transferService "github.com/boxproject/apiServer/service/transfer"
	"github.com/gin-gonic/gin"
)

// GetCoins get the list of coins
func GetCoins(ctx *gin.Context) {
	var msg errors.ErrModel
	msg = coinService.CoinListByTransfer()
	msg.RetErr(ctx)
}

// GetTemplate get the list of template flow
func GetTemplate(ctx *gin.Context) {
	var msg errors.ErrModel
	var order transferService.TemplateVo
	if err := ctx.ShouldBind(&order); err != nil {
		log.Error("Bind Param Error", err)
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	} else {
		msg = templateService.TemplateListTransfer(order.CoinId)
	}
	msg.RetErr(ctx)
}

// FindApplyList get the list of transfer
func FindApplyList(ctx *gin.Context) {
	var msg errors.ErrModel
	var order transferService.TransferOrder
	session := ctx.MustGet("claims").(*middle.CustomClaims)
	if err := ctx.ShouldBind(&order); err != nil {
		log.Error("Bind Param Error", err)
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	} else {
		msg = transferService.Findapplylist(order.Sort, session.ID, session.Account)
	}
	msg.RetErr(ctx)
}

// FindApplyLog get the log of transfer
func FindApplyLog(ctx *gin.Context) {
	var msg errors.ErrModel
	var order transferService.TransferOrder
	if err := ctx.ShouldBind(&order); err != nil {
		log.Error("Bind Param Error", err)
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	} else {
		msg = transferService.FindApplyLog(order.Id)
	}
	msg.RetErr(ctx)
}

// FindApplyById get the detail info of a transfer order by id
func FindApplyById(ctx *gin.Context) {
	var msg errors.ErrModel
	var order transferService.TransferOrder
	if err := ctx.ShouldBind(&order); err != nil {
		log.Error("Bind Param Error", err)
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	} else {
		msg = transferService.Findapplybyid(order.Id)
	}
	msg.RetErr(ctx)
}


// Apply employee apply for transfer
func Apply(ctx *gin.Context) {

	var trans transferService.TransferVo
	var msg errors.ErrModel
	if err := ctx.ShouldBind(&trans); err != nil {
		log.Error("Bind Param Error", err)
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	} else {
		session := ctx.MustGet("claims").(*middle.CustomClaims)
		trans.ApplyerId = session.ID
		trans.AppName = session.Account
		msg = transferService.Apply(&trans)
	}

	msg.RetErr(ctx)
}

// Verify approval the transfer
func Verify(ctx *gin.Context) {
	log.Debug("verify transfer...")
	var trans transferService.VerifyApplyVo
	var msg errors.ErrModel
	if err := ctx.ShouldBind(&trans); err != nil {
		log.Error("Bind Param Error", err)
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	} else {
		session := ctx.MustGet("claims").(*middle.CustomClaims)
		trans.AppName = session.Account
		msg = transferService.Verify(&trans)
	}
	msg.RetErr(ctx)
}

// Cancel cancel the transfer by applyer
func Cancel(ctx *gin.Context) {
	var trans transferService.VerifyApplyVo
	var msg errors.ErrModel
	if err := ctx.ShouldBind(&trans); err != nil {
		log.Error("Bind Param Error", err)
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	} else {
		session := ctx.MustGet("claims").(*middle.CustomClaims)
		trans.AccountId = session.ID
		msg = transferService.Cancel(&trans)
	}
	msg.RetErr(ctx)
}

// BatchList list the transfer to approval
func BatchList(ctx *gin.Context) {
	var trans transferService.TransferVo
	var msg errors.ErrModel
	if err := ctx.ShouldBind(&trans); err != nil {
		log.Error("Bind Param Error", err)
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	} else {
		log.Debug("BatchList orderids", trans.OrderIds)
		var arr []string
		if err = json.Unmarshal([]byte(trans.OrderIds), &arr); err != nil {
			log.Error("json解析错误", err)
			msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
		} else {
			trans.ArrOrderIds = arr
			session := ctx.MustGet("claims").(*middle.CustomClaims)
			trans.ApplyerId = session.ID
			msg = transferService.BatchList(&trans)
		}
	}
	msg.RetErr(ctx)
}

// BatchVerify batch approval the transfer
func BatchVerify(ctx *gin.Context) {
	log.Debug("BatchVerify...")
	var Batch transferService.BatchVerify
	var msg errors.ErrModel
	log.Debug("batch", Batch)
	if err := ctx.ShouldBind(&Batch); err != nil {
		log.Error("Bind Param Error", err)
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	} else {
		var arr [](transferService.BatchOrder)
		if err = json.Unmarshal([]byte(Batch.OrderIds), &arr); err != nil {
			log.Error("json解析错误", err)
			msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
		} else {
			var trans transferService.VerifyApplyVo
			session := ctx.MustGet("claims").(*middle.CustomClaims)
			trans.AppName = session.Account
			var rets = make(map[string]interface{})
			for i := 0; i < len(arr); i++ {
				trans.OrderId = arr[i].OrderId
				trans.Status = arr[i].Status
				trans.TransferSign = arr[i].TransferSign
				trans.Reason = Batch.Reason
				ret := transferService.Verify(&trans)
				rets[trans.OrderId] = map[string]interface{}{
					"transStatus": trans.Status,
					"retCode":     ret.Code,
					"retErr":      ret.RetErr,
				}
			}
			log.Info("批量审批结果", rets)
		}
		msg = errors.ErrModel{Code: errors.Success, Err: nil}
	}
	msg.RetErr(ctx)
}

// GetTemplatebyHash get the template info by flow hash
func GetTemplatebyHash(ctx *gin.Context) {
	var trans transferService.VerifyApplyVo
	var msg errors.ErrModel
	if err := ctx.ShouldBind(&trans); err != nil {
		log.Error("Bind Param Error", err)
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	} else {
		msg = transferService.GetTemplatebyHash(trans.TemHash)
	}
	msg.RetErr(ctx)
}


// TransferCommit web提交转账
func TransferCommit(ctx *gin.Context) {
	var trans transferService.WebTransfer
	var msg errors.ErrModel
	if err := ctx.ShouldBind(&trans); err != nil {
		log.Error("Bind Param Error", err)
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	} else {
		msg = transferService.TransferCommit(&trans)
	}
	msg.RetErr(ctx)

}

// FindTranfersById app获取转账信息
func FindTranfersById(ctx *gin.Context) {
	var trans transferService.WebTransfer
	var msg errors.ErrModel
	if err := ctx.ShouldBind(&trans); err != nil {
		log.Error("Bind Param Error", err)
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	} else {
		msg = transferService.FindTranfersById(&trans)
	}
	msg.RetErr(ctx)
}

// GetWebRouter 获取web转账网页路径
func GetWebRouter(ctx *gin.Context) {
	var msg errors.ErrModel
	msg = transferService.GetWebRouter()
	msg.RetErr(ctx)

}

func GetCommitStatus(ctx *gin.Context) {
	var trans transferService.WebTransfer
	var msg errors.ErrModel
	if err := ctx.ShouldBind(&trans); err != nil {
		log.Error("Bind Param Error", err)
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	} else {
		msg = transferService.GetCommitStatus(&trans)
	}

	msg.RetErr(ctx)

}
