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
package coin

import (
	"github.com/boxproject/apiServer/db"
	reterror "github.com/boxproject/apiServer/errors"
	"github.com/boxproject/boxwallet/bccore"
	walletClient "github.com/boxproject/boxwallet/cli"

	//"github.com/ethereum/go-ethereum/log"
	log "github.com/alecthomas/log4go"
	"encoding/json"
	voucher "github.com/boxproject/apiServer/rpc"
)

// AddCoin add coin
func AddCoin(coinInfo *Pcoin) reterror.ErrModel {
	ds := &db.CoinService{}
	// check if coin exists
	bool, err := ds.CoinIsExist(coinInfo.TokenAddress, coinInfo.Symbol)
	if bool == true { //已存在
		return reterror.ErrModel{Err: reterror.MSG_6001, Code: reterror.Coin_6001}
	}
	// add info to db
	err = ds.AddCoin(coinInfo.Symbol, coinInfo.Name, coinInfo.Decimals, coinInfo.TokenAddress)
	if err != nil {
		log.Error("Add Coin failed case:", err)
		return reterror.ErrModel{Err: reterror.MSG_1003, Code: reterror.AddDepFailed}
	}
	return reterror.ErrModel{Err: nil, Code: reterror.Success}
}

// VerifyAddress validate address
func VerifyAddress(address string) reterror.ErrModel {
	walletService := walletClient.NewAppServer()
	coinInfo, err := walletService.GetCoinInfo(bccore.BC_ETH, address)
	if err != nil {
		log.Error("校验地址失败", err)
		return reterror.ErrModel{Err: reterror.MSG_6003, Code: reterror.Coin_6003}
	}
	data := make(map[string]interface{})
	data["Name"] = coinInfo.Name
	data["Symbol"] = coinInfo.Symbol
	data["Decimals"] = coinInfo.Decimals
	return reterror.ErrModel{Err: nil, Code: reterror.Success, Data: data}
}

// CoinList get coin list
func CoinList(typeCoin int) reterror.ErrModel {
	ds := &db.CoinService{}
	list, err := ds.CoinList(typeCoin)
	if err != nil {
		log.Error("get coin list error", err)
		return reterror.ErrModel{Err: reterror.SYSTEM_ERROR, Code: reterror.Failed}
	}
	data := make([]interface{}, len(list))
	for i, v := range list {
		data[i] = map[string]interface{}{
			"ID":           v.ID,
			"Name":         v.Name,
			"FullName":     v.FullName,
			"TokenType":    v.TokenType,
			"TokenAddress": v.TokenAddress,
			"Status":       v.Available,
			"Precise":      v.Precise,
		}
	}
	return reterror.ErrModel{Err: nil, Code: reterror.Success, Data: data}
}

// CoinBalance get balance
func CoinBalance() reterror.ErrModel {
	ds := &db.CoinService{}
	list, err := ds.CoinBalance()
	if err != nil {
		log.Error("get coin balance error", err)
		return reterror.ErrModel{Err: reterror.SYSTEM_ERROR, Code: reterror.Failed}
	}
	return reterror.ErrModel{Err: nil, Code: reterror.Success, Data: list}
}

// QRcode address info for generate qrcode
func QRcode(id int) reterror.ErrModel {
	// required parameters
	if id == 0 {
		return reterror.ErrModel{Code: reterror.ParamsNil, Err: reterror.PARAMS_NULL}
	}

	// get signed address
	addressSign := db.ConfigMap["address_sign"]
	if addressSign == "" {
		log.Error("没有对地址签名")
		return reterror.ErrModel{Code: reterror.Acc_Manage_6016, Err: reterror.MSG_6016}
	}

	signs := &db.AddressSigns{}
	err := json.Unmarshal([]byte(addressSign), &signs)
	if err != nil {
		log.Error("地址签名解析json错误")
		return reterror.ErrModel{Code: reterror.ParamsNil, Err: reterror.PARAMS_NULL}
	}

	ds := &db.CoinService{}
	addresses, err := ds.QRcode(id)
	if err != nil {
		log.Error("get coin qrcode error", err)
		return reterror.ErrModel{Err: reterror.SYSTEM_ERROR, Code: reterror.Failed}
	}

	data := make(map[string]interface{})
	for _, v := range signs.AddressInfo {
		if v.MasterAddress == addresses["MainAddress"] {
			data["Sign"] = v.Sign
			break
		}
	}

	data["MainAddress"] = addresses["MainAddress"]
	data["RandomAddress"] = addresses["ChildAddress"]

	data["Account"] = signs.Account
	data["RandomIndex"] = addresses["Index"]
	// get voucher public key
	oper := &voucher.GrpcServer{
		Type: voucher.VOUCHER_OPERATE_MASTER_PUBKEY,
	}
	_, res, ret := voucher.SendVoucherData(oper)
	if ret != voucher.VRET_CLIENT {
		log.Error("签名机获取主公钥错误")
		return reterror.ErrModel{Err: reterror.SYSTEM_ERROR, Code: reterror.Failed}
	}
	data["MasterKey"] = res.PublicKeys
	return reterror.ErrModel{Err: nil, Code: reterror.Success, Data: data}
}

// CoinStauts available/unavailable coin
func CoinStauts(id, status int) reterror.ErrModel {
	ds := &db.CoinService{}
	coinInfo, err := ds.CoinStauts(id, status)
	if err != nil {
		log.Error("update coin status error", err)
		return reterror.ErrModel{Err: reterror.MSG_6006, Code: reterror.Coin_6006}
	}

	if coinInfo == nil {
		return reterror.ErrModel{Err: reterror.MSG_6005, Code: reterror.Coin_6005}
	}
	return reterror.ErrModel{Err: nil, Code: reterror.Success}

}

// CoinListByTransfer get coin list(for transfer)
func CoinListByTransfer() reterror.ErrModel {
	ds := &db.CoinService{}
	list, err := ds.CoinListByTransfer()
	if err != nil {
		log.Error("get coin list by transfer error", err)
		return reterror.ErrModel{Err: reterror.SYSTEM_ERROR, Code: reterror.Failed}
	}
	return reterror.ErrModel{Err: nil, Code: reterror.Success, Data: list}
}
