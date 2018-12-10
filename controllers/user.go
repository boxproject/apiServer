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
	"strconv"

	log "github.com/alecthomas/log4go"
	"github.com/boxproject/apiServer/errors"
	middle "github.com/boxproject/apiServer/middleware"
	"github.com/boxproject/apiServer/service/logger"
	userService "github.com/boxproject/apiServer/service/user"
	"github.com/boxproject/apiServer/service/verify"
	"github.com/gin-gonic/gin"
	"github.com/boxproject/apiServer/common"
)

// Signup registrat a new account
func Signup(ctx *gin.Context) {
	var user userService.UserVo
	var msg errors.ErrModel
	if err := ctx.ShouldBind(&user); err != nil {
		log.Error("Bind Param Error", err)
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	} else {
		j_user, _ := json.Marshal(user)
		log.Debug("注册参数", string(j_user))
		msg = userService.Signup(&user)
	}
	msg.RetErr(ctx)
}

// Login
func Login(ctx *gin.Context) {
	var user userService.UserVo
	var msg errors.ErrModel
	if err := ctx.ShouldBind(&user); err != nil {
		log.Error("Bind Param Error", err)
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	} else {
		msg = userService.Login(&user)
	}
	msg.RetErr(ctx)
}

// ModifyPassword change the password to log in
func ModifyPassword(ctx *gin.Context) {
	var obj userService.ModifyPSD
	var msg errors.ErrModel
	claims := ctx.MustGet("claims").(*middle.CustomClaims)
	if err := ctx.ShouldBind(&obj); err != nil {
		log.Error("Bind Param Error", err)
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	} else {
		msg = userService.ModifyPassword(&obj, claims.Account)
	}

	msg.RetErr(ctx)
}

// InsideLetter Station letter
func InsideLetter(ctx *gin.Context) {
	var msg errors.ErrModel
	var obj userService.ModifyPSD
	claims := ctx.MustGet("claims").(*middle.CustomClaims)
	if err := ctx.ShouldBind(&obj); err != nil {
		log.Error("Bind Param Error", err)
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	} else {
		msg = userService.InsideLetter(ctx.GetHeader(common.HeaderLangKey), claims.ID, claims.UserType, obj.Page)
	}
	msg.RetErr(ctx)
}

// ReadLetter Save reading history
func ReadLetter(ctx *gin.Context) {
	var msg errors.ErrModel
	var user userService.ModifyPSD
	claims := ctx.MustGet("claims").(*middle.CustomClaims)
	if err := ctx.ShouldBind(&user); err != nil {
		log.Error("Bind Param Error", err)
		msg = errors.ErrModel{Code: errors.ParamsNil, Err: errors.PARAMS_NULL}
	} else {
		msg = userService.ReadLetter(claims.ID, user.ID)
	}

	msg.RetErr(ctx)
}

// Scan Superior submit the registration
func Scan(ctx *gin.Context) {
	tokenUser := ctx.MustGet("claims").(*middle.CustomClaims)
	var user userService.UserVo
	var msg errors.ErrModel
	if err := ctx.ShouldBind(&user); err != nil {
		log.Error("Bind Param Error", err)
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	} else {
		user.AppId = tokenUser.AppID
		user.Name = tokenUser.Account
		user.Level = tokenUser.Level
		msg = userService.Scan(&user)
	}

	msg.RetErr(ctx)
}

// FindUserStatus get the account status
func FindUserStatus(ctx *gin.Context) {
	var user userService.UserVo
	var msg errors.ErrModel
	if err := ctx.ShouldBind(&user); err != nil {
		log.Error("Bind Param Error", err)
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	} else {
		msg = userService.FindUserStatus(&user)
	}
	msg.RetErr(ctx)
}

// VerifyUser Superior approval the registration
func VerifyUser(ctx *gin.Context) {
	var user userService.UserVo
	log.Debug("审核参数", user.DepID)
	var msg errors.ErrModel
	if err := ctx.ShouldBind(&user); err != nil {
		log.Error("Bind Param Error", err)
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	} else {
		msg = userService.VerifyUser(&user)
	}
	msg.RetErr(ctx)
}

// UserTree get the Organizational structure
func UserTree(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*middle.CustomClaims)
	msg := userService.UserTree(claims.PubKey)
	msg.RetErr(ctx)
}

// GetAccountsByType get the accounts info by user type
func GetAccountsByType(ctx *gin.Context) {
	var user userService.UserVo
	var msg errors.ErrModel
	if err := ctx.ShouldBind(&user); err != nil {
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	} else {
		msg = userService.GetAccountsByType(user.UserType)
	}
	msg.RetErr(ctx)
}

// AddAdmin add a new manager
func AddAdmin(ctx *gin.Context) {
	var user userService.UserVo
	var msg errors.ErrModel
	session := ctx.MustGet("claims").(*middle.CustomClaims)
	if err := ctx.ShouldBind(&user); err != nil {
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	} else {
		msg = userService.AddAdmin(user.Add, session.ID)
	}
	msg.RetErr(ctx)
}

// DelAdmin deprivation the account's admin right
func DelAdmin(ctx *gin.Context) {
	var user userService.UserVo
	var msg errors.ErrModel
	session := ctx.MustGet("claims").(*middle.CustomClaims)
	if err := ctx.ShouldBind(&user); err != nil {
		log.Error("删除管理员错误", err)
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	} else {
		msg = userService.DelAdmin(user.ID, session.ID)
	}
	msg.RetErr(ctx)
}

// ReSignUp re-registration
func ReSignUp(ctx *gin.Context) {
	var user userService.UserVo
	var msg errors.ErrModel
	if err := ctx.ShouldBind(&user); err != nil {
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	} else {
		msg = userService.ReSignUp(user.AppId)
	}
	msg.RetErr(ctx)
}

// GetAllUsers get all members
func GetAllUsers(ctx *gin.Context) {
	var msg errors.ErrModel
	msg = userService.GetAllUsers()
	msg.RetErr(ctx)
}

// DelayTaskNum count To-do list
func DelayTaskNum(ctx *gin.Context) {
	var msg errors.ErrModel
	session := ctx.MustGet("claims").(*middle.CustomClaims)
	msg = userService.DelayTaskNum(session.Account, session.UserType, session.ID)
	msg.RetErr(ctx)
}

// GetUserByID get account info by id
func GetUserByID(ctx *gin.Context) {
	var user userService.UserVo
	var msg errors.ErrModel
	if err := ctx.ShouldBind(&user); err != nil {
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	} else {
		msg = userService.GetUserByID(user.ID)
	}
	msg.RetErr(ctx)
}

// SetUser update account info
func SetUser(ctx *gin.Context) {
	var user userService.PAccount
	var msg errors.ErrModel
	claims := ctx.MustGet("claims").(*middle.CustomClaims)
	if err := ctx.ShouldBind(&user); err != nil {
		log.Error("setUser", err)
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	} else {
		msg = userService.SetUser(&user, claims.ID)
	}
	msg.RetErr(ctx)
}

// DisableAcc froze an account
func DisableAcc(ctx *gin.Context) {
	var user userService.UserVo
	var msg errors.ErrModel
	claims := ctx.MustGet("claims").(*middle.CustomClaims)
	if err := ctx.ShouldBind(&user); err != nil {
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	} else {
		msg = userService.DisableAcc(user.ID, claims.Account)
	}
	msg.RetErr(ctx)
}

// GetRegList get the registration list
func GetRegList(ctx *gin.Context) {
	var msg errors.ErrModel
	msg = userService.GetRegList()
	msg.RetErr(ctx)
}

// HasRecovery check the recover apply form the owner
func HasRecovery(ctx *gin.Context) {
	var msg errors.ErrModel
	var user userService.UserVo
	if err := ctx.ShouldBind(&user); err != nil {
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	} else {
		msg = userService.HasRecovery(user.Name)
	}
	msg.RetErr(ctx)
}

// RecoveryOwner recover a owner account
func RecoveryOwner(ctx *gin.Context) {
	log.Debug("RecoveryOwner...")
	var msg errors.ErrModel
	var user userService.UserVo
	if err := ctx.ShouldBind(&user); err != nil {
		log.Debug("恢复股东参数绑定错误", err)
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	} else {
		msg = userService.RecoveryOwner(&user)
		// 日志
		if msg.Code == errors.Success {
			regID := ""
			if data, ok := msg.Data.(map[string]interface{}); ok {
				if id, ok := data["RegId"].(int); ok {
					regID = strconv.Itoa(id)
				}
				logErr := logger.AddLog("recovery", "", user.Name, common.LoggerOwnerRecovery, regID)
				if logErr != nil {
					log.Error("日志未记录ownerRecovery", user.Name)
				}
			}
		}
	}
	msg.RetErr(ctx)
}

// RecoveryList list of accounts to recover
func RecoveryList(ctx *gin.Context) {
	var msg errors.ErrModel
	claims := ctx.MustGet("claims").(*middle.CustomClaims)
	msg = userService.RecoveryList(claims.Account, claims.PubKey)
	msg.RetErr(ctx)
}

// SubRecovery submit the recover apply
func SubRecovery(ctx *gin.Context) {
	var msg errors.ErrModel
	var user userService.UserVo
	claims := ctx.MustGet("claims").(*middle.CustomClaims)
	if err := ctx.ShouldBind(&user); err != nil {
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	} else {
		if claims.Account == user.Name {
			msg = errors.ErrModel{Code: errors.User_2014, Err: errors.MSG_2014}
		} else {
			msg = userService.SubRecovery(user.AppId, user.Name, user.RegID, claims.Account)
		}
	}
	msg.RetErr(ctx)
}

// VerifyPwd verify the account
func VerifyPwd(ctx *gin.Context) {
	var msg errors.ErrModel
	var user userService.UserVo
	claims := ctx.MustGet("claims").(*middle.CustomClaims)
	if err := ctx.ShouldBind(&user); err != nil {
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	} else {
		msg = userService.VerifyPwd(user.Pwd, claims.Account, claims.ID)
	}
	msg.RetErr(ctx)
}

// RecoveryResult get the result of verifaction
func RecoveryResult(ctx *gin.Context) {
	var msg errors.ErrModel
	var user userService.UserVo
	if err := ctx.ShouldBind(&user); err != nil {
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	} else {
		msg = userService.RecoveryResult(user.RegID, user.AppId)
	}
	msg.RetErr(ctx)
}

// VerifyRecovery check the owner account recovered
func VerifyRecovery(ctx *gin.Context) {
	var msg errors.ErrModel
	var voucher userService.VoucherVo
	claims := ctx.MustGet("claims").(*middle.CustomClaims)
	if err := ctx.ShouldBind(&voucher); err != nil {
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	} else {
		msg = userService.VerifyRecovery(claims.Account, &voucher)
	}
	msg.RetErr(ctx)
}

// ResetRecovery re-registration
func ResetRecovery(ctx *gin.Context) {
	var msg errors.ErrModel
	var user userService.UserVo
	if err := ctx.ShouldBind(&user); err != nil {
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	} else {
		msg = userService.ResetRecovery(user.Name)
	}
	msg.RetErr(ctx)
}

// ActiveRecovery Activate owner
func ActiveRecovery(ctx *gin.Context) {
	var msg errors.ErrModel
	var user userService.VoucherVo
	if err := ctx.ShouldBind(&user); err != nil {
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	} else {
		msg = userService.ActiveRecovery(&user)
	}
	msg.RetErr(ctx)
}

// GetPubKeys get the subordinates' signed pubkey info
func GetPubKeys(ctx *gin.Context) {
	var msg errors.ErrModel
	var user userService.UserVo
	if err := ctx.ShouldBind(&user); err != nil {
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	} else {
		msg = userService.GetPubkeys(user.Name)
	}
	msg.RetErr(ctx)
}

// VerifyPassword Check password
func VerifyPassword() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session := ctx.MustGet("claims").(*middle.CustomClaims)
		password := ctx.PostForm("password")
		result, _, errorCode := verify.VerifyPSW(session.Account, password)
		log.Debug("very password", session.Account, password)
		if !result.Result {
			msg := errors.ErrModel{Code: errorCode, Err: errors.New(result.Reason), Data: map[string]interface{}{"data": result.Data}}
			ctx.Abort()
			msg.RetErr(ctx)
		}
		return
	}
}

// GetVersion get the api version
func GetVersion(ctx *gin.Context) {
	var msg errors.ErrModel
	msg = userService.GetVersion()
	msg.RetErr(ctx)
}

// GetBlockHeight get the current block height info of the system
func GetBlockHeight(ctx *gin.Context) {
	var msg errors.ErrModel
	msg = userService.GetBlockHeight()
	msg.RetErr(ctx)
}
