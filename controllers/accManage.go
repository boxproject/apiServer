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
	"strconv"
	"strings"
	log "github.com/alecthomas/log4go"
	"github.com/boxproject/apiServer/service/accmanage"
	"github.com/boxproject/apiServer/errors"
	middle "github.com/boxproject/apiServer/middleware"
	"github.com/gin-gonic/gin"
)

// Statistical statistics the number of address of the system
func Statistical(ctx *gin.Context) {
	var msg errors.ErrModel
	msg = accmanage.Statistical()
	msg.RetErr(ctx)
}

// ChildStatistical statistics the number of the child address of each currency
func ChildStatistical(ctx *gin.Context) {
	var msg errors.ErrModel
	var qt accmanage.Qaddress
	if err := ctx.ShouldBind(&qt); err != nil {
		log.Error("Bind Param Error", err)
		msg = errors.ErrModel{Code: errors.ParamsNil, Err: errors.PARAMS_NULL}
	} else {
		msg = accmanage.ChildStatistical(qt.CoinID)
	}
	msg.RetErr(ctx)
}

//CountAllChild counts the balance and number of each child of each currency
func CountAllChild(ctx *gin.Context) {
	var msg errors.ErrModel
	var qt accmanage.Qaddress
	if err := ctx.ShouldBind(&qt); err != nil {
		log.Error("Bind Param Error", err)
		msg = errors.ErrModel{Code: errors.ParamsNil, Err: errors.PARAMS_NULL}
	} else {
		msg = accmanage.CountAllChild(qt.CoinID)
	}
	msg.RetErr(ctx)
}

// AddressList lists the information of all the address
func AddressList(ctx *gin.Context) {
	var msg errors.ErrModel
	var qt accmanage.Qaddress
	if err := ctx.ShouldBind(&qt); err != nil {
		log.Error("Bind Param Error", err)
		msg = errors.ErrModel{Code: errors.ParamsNil, Err: errors.PARAMS_NULL}
	} else {
		msg = accmanage.AddressList(&qt)
	}
	msg.RetErr(ctx)
}

// AccountDetail shows the detail of the address information by given currency id
func AccountDetail(ctx *gin.Context) {
	var msg errors.ErrModel
	var qt accmanage.Qaddress
	if err := ctx.ShouldBind(&qt); err != nil {
		log.Error("Bind Param Error", err)
		msg = errors.ErrModel{Code: errors.ParamsNil, Err: errors.PARAMS_NULL}
	} else {
		msg = accmanage.AccountDetail(qt.ID, qt.CoinID)
	}
	msg.RetErr(ctx)
}

// SetTag used to set tag for the given address
func SetTag(ctx *gin.Context) {
	var msg errors.ErrModel
	var qt accmanage.Qaddress
	if err := ctx.ShouldBind(&qt); err != nil {
		log.Error("Bind Param Error", err)
		msg = errors.ErrModel{Code: errors.ParamsNil, Err: errors.PARAMS_NULL}
	} else {
		msg = accmanage.SetTag(qt.ID, qt.Tag)
	}

	msg.RetErr(ctx)
}

// AddressTopTen shows the top 10 rich address
func AddressTopTen(ctx *gin.Context) {
	var msg errors.ErrModel
	var qt accmanage.Qaddress
	if err := ctx.ShouldBind(&qt); err != nil {
		log.Error("Bind Param Error", err)
		msg = errors.ErrModel{Code: errors.ParamsNil, Err: errors.PARAMS_NULL}
	} else {
		msg = accmanage.AddressTopTen(&qt)
	}
	msg.RetErr(ctx)
}

// TransferRecord shows the transfer record of the given account
func TransferRecord(ctx *gin.Context) {
	var msg errors.ErrModel
	var qt accmanage.Qaddress
	if err := ctx.ShouldBind(&qt); err != nil {
		log.Error("Bind Param Error", err)
		msg = errors.ErrModel{Code: errors.ParamsNil, Err: errors.PARAMS_NULL}
	} else {
		msg = accmanage.TransferRecord(qt.ID)
	}
	msg.RetErr(ctx)
}

// TokenRecord shows the token transfer record of the given address
func TokenRecord(ctx *gin.Context) {
	var msg errors.ErrModel
	var qt accmanage.Qaddress
	if err := ctx.ShouldBind(&qt); err != nil {
		log.Error("Bind Param Error", err)
		msg = errors.ErrModel{Code: errors.ParamsNil, Err: errors.PARAMS_NULL}
	} else {
		msg = accmanage.TokenRecord(qt.ID, qt.CoinID)
	}
	msg.RetErr(ctx)
}

// TokenList lists the token info of the given address
func TokenList(ctx *gin.Context) {
	var msg errors.ErrModel
	var qt accmanage.Qaddress
	if err := ctx.ShouldBind(&qt); err != nil {
		log.Error("Bind Param Error", err)
		msg = errors.ErrModel{Code: errors.ParamsNil, Err: errors.PARAMS_NULL}
	} else {
		msg = accmanage.TokenList(qt.ID)
	}
	msg.RetErr(ctx)
}

// AddChildAccount used to create a new child account of the given account
func AddChildAccount(ctx *gin.Context) {
	var msg errors.ErrModel
	var ct accmanage.Qaddress
	claims := ctx.MustGet("claims").(*middle.CustomClaims)
	if err := ctx.ShouldBind(&ct); err != nil {
		log.Error("Bind Param Error", err)
		msg = errors.ErrModel{Code: errors.ParamsNil, Err: errors.PARAMS_NULL}
	} else {
		msg = accmanage.AddChildAccount(&ct, claims.Account)
	}
	msg.RetErr(ctx)
}

// MergeAccount used to consolidate account balance to master address
func MergeAccount(ctx *gin.Context) {
	var msg errors.ErrModel
	var qt accmanage.Qaddress
	claims := ctx.MustGet("claims").(*middle.CustomClaims)
	if err := ctx.ShouldBind(&qt); err != nil {
		log.Error("Bind Param Error", err)
		msg = errors.ErrModel{Code: errors.ParamsNil, Err: errors.PARAMS_NULL}
	} else {
		var ids []int
		sl := strings.Split(qt.Ids, ",")
		for _, val := range sl {
			an, _ := strconv.Atoi(val)
			ids = append(ids, an)
		}
		msg = accmanage.MergeAccount(qt.CoinID, qt.Type, ids, qt.PSW, claims.Account)
	}
	msg.RetErr(ctx)
}

// ContractAddr used to get the contract address from voucher
func ContractAddr(ctx *gin.Context) {
	msg := accmanage.GetContractAddress()
	msg.RetErr(ctx)
}

// GenContractAddr used to apply for the contract from contract
func GenContractAddr(ctx *gin.Context) {
	msg := accmanage.GenContractAddress()
	msg.RetErr(ctx)
}
