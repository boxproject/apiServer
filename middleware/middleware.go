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
package middleware

import (
	"time"
	"strconv"
	log "github.com/alecthomas/log4go"
	"github.com/boxproject/apiServer/db"
	"github.com/boxproject/apiServer/errors"
	reterror "github.com/boxproject/apiServer/errors"
	voucher "github.com/boxproject/apiServer/rpc"
	"github.com/boxproject/apiServer/utils"
	"github.com/boxproject/apiServer/service/verify"
	"github.com/gin-gonic/gin"
	"github.com/boxproject/apiServer/common"
)

const ReqTimeOut = 20 * 1000

// AdminAuth middleware to validate userType,  owner: "owner", admin:"admin"
func AdminAuth(authType string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		claims := ctx.MustGet("claims").(*CustomClaims)
		switch {
		case authType == "owner" && claims.UserType != common.OwnerAccType:
			msg := errors.ErrModel{Code: errors.Code_3, Err: reterror.MSG_3}
			msg.RetErr(ctx)
			ctx.Abort()
		case authType == "admin" && claims.UserType != common.AdminAccType && claims.UserType != common.OwnerAccType:
			msg := errors.ErrModel{Code: errors.Code_3, Err: reterror.MSG_3}
			msg.RetErr(ctx)
			ctx.Abort()
		}
	}
}

func VerifyParamSignMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		var sign string
		var timestamp string
		pubkey := c.MustGet("claims").(*CustomClaims).PubKey
		if c.Request.Form == nil {
			c.Request.ParseMultipartForm(32 << 20)
		}
		for k, v := range c.Request.Form {
			if k == "sign" {
				sign = v[0]
			} else if k == "timestamp" {
				timestamp = v[0]
			}
		}
		// 校验时间戳
		timeNow := time.Now().UnixNano() / 1e6
		timestamp = utils.Format(timestamp, 13)
		timestampInt64, err := strconv.ParseInt(timestamp, 10, 64)
		if err != nil {
			log.Debug("请求者：", c.MustGet("claims").(*CustomClaims).Account)
			log.Debug("解析时间戳出错", err)
			msg := errors.ErrModel{Code: errors.Failed, Err: reterror.SYSTEM_ERROR}
			c.Abort()
			msg.RetErr(c)
			return
		}

		if utils.AbsInt64(timeNow-timestampInt64) > int64(ReqTimeOut) {
			log.Error("时间差超过20秒", c.MustGet("claims").(*CustomClaims).Account)
			msg := errors.ErrModel{Code: errors.Code_5, Err: reterror.MSG_5}
			c.Abort()
			msg.RetErr(c)
			return
		}
		err = utils.RsaVerify([]byte(pubkey), []byte(timestamp), []byte(sign))
		if err != nil {
			log.Error("参数签名错误", c.MustGet("claims").(*CustomClaims).Account, err)
			msg := errors.ErrModel{Code: errors.VerifyParamFail, Err: reterror.MSG_4}
			c.Abort()
			msg.RetErr(c)
			return
		}
		return
	}
}

// AuthCheck 权限校验
func AuthCheck(authType string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accountID := ctx.MustGet("claims").(*CustomClaims).ID
		ds := &db.AuthService{}
		err := ds.AuthCheck(authType, accountID)
		if err != nil {
			msg := errors.ErrModel{Code: errors.Auth_5007, Err: reterror.MSG_5007}
			ctx.Abort()
			msg.RetErr(ctx)
		}
	}
}

// IsFrozen check account is frozen
func IsFrozen() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		account := ctx.MustGet("claims").(*CustomClaims).Account
		ds := &db.AccDBService{}
		frozen, frozenTo, accErr := ds.AccoutIsFrozen(account)
		if accErr != nil {
			msg := errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
			ctx.Abort()
			msg.RetErr(ctx)
		}
		if time.Now().Before(frozenTo) && frozen == 1 {
			message := "输入错误次数过多,账户被锁定到:" + frozenTo.Local().Format("2006-01-02 03:04:05 AM")
			msg := errors.ErrModel{Code: errors.User_3007, Err: errors.New(message), Data: map[string]interface{}{
				"data": frozenTo.Local().Format("2006-01-02 03:04:05 AM"),
			}}
			ctx.Abort()
			msg.RetErr(ctx)
		}
	}
}

type ReqBody struct {
	AppID     string `form:"app_id"`
	AppName   string `form:"name"`
	AesKey    string `form:"aeskey"`
	Msg       string `form:"msg"`
	Sign      string `form:"sign"`
	Timestamp int64  `form:"timestamp"`
}

// VerifyKeyWord validate key word
func VerifyKeyWord() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var reqBody ReqBody
		if err := ctx.ShouldBind(&reqBody); err != nil {
			msg := reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
			ctx.Abort()
			msg.RetErr(ctx)
			return
		}
		values, exists := ctx.Get("claims")
		appID := reqBody.AppID
		appName := reqBody.AppName
		if exists {
			claims := values.(*CustomClaims)
			appID = claims.AppID
			appName = claims.Account
		}
		// check if account is frozen
		accountErr := verify.IsFrozen(appName)
		if accountErr.Code != 0 {
			ctx.Abort()
			accountErr.RetErr(ctx)
		}
		ds := db.AccDBService{}
		// validate key word
		oper := &voucher.GrpcServer{
			Type:      voucher.VOUCHER_OPERATE_CHECK_PASS,
			AppId:     appID,
			AppName:   appName,
			AesKey:    []byte(reqBody.AesKey),
			Sign:      []byte(reqBody.Sign),
			Msg:       []byte(reqBody.Msg),
			Timestamp: reqBody.Timestamp,
		}
		_, voucherRes, voucherRet := voucher.SendVoucherData(oper)
		var msg reterror.ErrModel
		if voucherRet == voucher.VRET_ERR {
			msg = reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
		} else {
			switch voucherRes.Status {
			case voucher.STATUS_APP_VERIFY_ERROR:
				msg = reterror.ErrModel{Code: reterror.Code_102, Err: reterror.MSG_102}
				ctx.Abort()
				msg.RetErr(ctx)
				return
			case voucher.STATUS_APP_PASSWORD_ERROR:
				msg := verify.Keyword(appID, appName)
				ctx.Abort()
				msg.RetErr(ctx)
				return
			case 0:
				_, accErr := ds.ResetAccount(appID)
				if accErr != nil {
					msg = reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
					ctx.Abort()
					msg.RetErr(ctx)
				}
			default:
				msg = reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
				ctx.Abort()
				msg.RetErr(ctx)
				return
			}
		}
	}
}
