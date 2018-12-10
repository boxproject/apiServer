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
	log "github.com/alecthomas/log4go"
	"github.com/boxproject/apiServer/service/coin"
	"github.com/boxproject/apiServer/errors"
	"github.com/gin-gonic/gin"
)

// AddCoin used to adds one coin to the list
func AddCoin(ctx *gin.Context) {
	var pcoin coin.Pcoin
	var msg errors.ErrModel
	if err := ctx.ShouldBind(&pcoin); err != nil {
		log.Error("Bind Param Error", err)
		msg = errors.ErrModel{Code: errors.ParamsNil, Err: errors.PARAMS_NULL}
	} else {
		msg = coin.AddCoin(&pcoin)
	}

	msg.RetErr(ctx)
}

// VerifyAddress used to verify the given address.
func VerifyAddress(ctx *gin.Context) {
	var pcoin coin.Pcoin
	var msg errors.ErrModel
	if err := ctx.ShouldBind(&pcoin); err != nil {
		log.Error("Bind Param Error", err)
		msg = errors.ErrModel{Code: errors.ParamsNil, Err: errors.PARAMS_NULL}
	} else {
		msg = coin.VerifyAddress(pcoin.TokenAddress)
	}

	msg.RetErr(ctx)
}

// CoinList used to get the coin list which the system now supports
func CoinList(ctx *gin.Context) {
	var msg errors.ErrModel
	var pcoin coin.Pcoin
	if err := ctx.ShouldBind(&pcoin); err != nil {
		log.Error("Bind Param Error", err)
		msg = errors.ErrModel{Code: errors.ParamsNil, Err: errors.PARAMS_NULL}
	} else {
		msg = coin.CoinList(pcoin.CoinType)
	}
	msg.RetErr(ctx)
}

// CoinBalance used to get the balance of each coin
func CoinBalance(ctx *gin.Context) {
	var msg errors.ErrModel
	msg = coin.CoinBalance()
	msg.RetErr(ctx)
}

// CoinStauts used to enable/disable one coin
func CoinStauts(ctx *gin.Context) {
	var msg errors.ErrModel
	var pcoin coin.Pcoin
	if err := ctx.ShouldBind(&pcoin); err != nil {
		log.Error("Bind Param Error", err)
		msg = errors.ErrModel{Code: errors.ParamsNil, Err: errors.PARAMS_NULL}
	} else {
		msg = coin.CoinStauts(pcoin.ID, pcoin.Status)
	}
	msg.RetErr(ctx)
}

// QRcode used to get the deposit address
func QRcode(ctx *gin.Context) {
	var msg errors.ErrModel
	var pcoin coin.Pcoin
	if err := ctx.ShouldBind(&pcoin); err != nil {
		log.Error("Bind Param Error", err)
		msg = errors.ErrModel{Code: errors.ParamsNil, Err: errors.PARAMS_NULL}
	} else {
		msg = coin.QRcode(pcoin.ID)
	}
	msg.RetErr(ctx)
}
