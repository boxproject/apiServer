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
	"github.com/boxproject/apiServer/service/auth"
	"github.com/boxproject/apiServer/errors"
	middle "github.com/boxproject/apiServer/middleware"
	"github.com/gin-gonic/gin"
	"github.com/boxproject/apiServer/common"
)

// AuthList gets the auth rights list.
func AuthList(ctx *gin.Context) {
	var msg errors.ErrModel
	msg = auth.GetAuthList(ctx.GetHeader(common.HeaderLangKey))
	msg.RetErr(ctx)
}

// AuthAccounts gets the member list by auth right
func AuthAccounts(ctx *gin.Context) {
	var msg errors.ErrModel
	var authMap auth.PAuth
	if err := ctx.ShouldBind(&authMap); err != nil {
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	} else {
		msg = auth.GetAuthAccounts(authMap.ID)
	}
	msg.RetErr(ctx)
}

// AddAuthToAccount used to adds a new auth right
func AddAuthToAccount(ctx *gin.Context) {
	var msg errors.ErrModel
	var authParams auth.PAuth
	claims := ctx.MustGet("claims").(*middle.CustomClaims)
	if err := ctx.ShouldBind(&authParams); err != nil {
		msg = errors.ErrModel{Code: errors.Auth_5002, Err: errors.MSG_5002}
	} else {
		msg = auth.AddAuthToAccount(authParams, claims.UserType, claims.Account)
	}
	msg.RetErr(ctx)
}

// DelAuthFromAccount used to cancel the auth right from the account
func DelAuthFromAccount(ctx *gin.Context) {
	var msg errors.ErrModel
	var pAuth auth.PAuth
	claims := ctx.MustGet("claims").(*middle.CustomClaims)
	if err := ctx.ShouldBind(&pAuth); err != nil {
		msg = errors.ErrModel{Code: errors.Auth_5002, Err: errors.MSG_5002}
	} else {
		msg = auth.DelAuthFromAccount(pAuth, claims.UserType, claims.Account)
	}
	msg.RetErr(ctx)
}
