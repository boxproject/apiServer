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
package accmanage

import (
	"encoding/json"
	"errors"
	"fmt"

	"sync"

	log "github.com/alecthomas/log4go"
	"github.com/boxproject/apiServer/db"
	reterror "github.com/boxproject/apiServer/errors"
	"github.com/boxproject/apiServer/service/verify"
	"github.com/boxproject/boxwallet/bccore"
	walletClient "github.com/boxproject/boxwallet/cli"
	"github.com/shopspring/decimal"
	"github.com/boxproject/apiServer/common"
	"strings"
)

var accLock = new(sync.Mutex)

// Statistical 主链币多账户统计主账户和子账户的个数
func Statistical() reterror.ErrModel {
	am := &db.AccManageService{}
	list := am.Statistical()
	accounts := make([]interface{}, 0, len(list))
	for _, account := range list {
		//TODO:目前只有ETH有代币，这个地方以后可能会修改
		if account.Name == "ETH" {
			account.HasToken = true
			account.HasAddChild = true
			account.HasMerge = true
		}
		if account.Name == string(bccore.STR_BTC) || account.Name == string(bccore.STR_LTC) {
			account.HasAddChild = true
		}
		accounts = append(accounts, account)
	}
	return reterror.ErrModel{Err: nil, Code: reterror.Success, Data: accounts}
}

// ChildStatistical 子账户统计
// id 币种ID
func ChildStatistical(id int) reterror.ErrModel {
	if id == 0 {
		return reterror.ErrModel{Code: reterror.ParamsNil, Err: reterror.PARAMS_NULL}
	}
	am := &db.AccManageService{}
	childCount, otherCount := am.ChildStatistical(id)
	data := make(map[string]interface{})
	data["ChildCount"] = childCount
	data["OtherCount"] = otherCount
	data["UsableCount"] = MaxCount - childCount - otherCount
	data["MaxCount"] = MaxCount
	return reterror.ErrModel{Err: nil, Code: reterror.Success, Data: data}
}

// CountAllChild 所有子账户数量和总额
func CountAllChild(id int) reterror.ErrModel {
	if id == 0 {
		return reterror.ErrModel{Code: reterror.ParamsNil, Err: reterror.PARAMS_NULL}
	}
	am := &db.AccManageService{}
	balance, amount, err := am.CountAllChild(id)
	if err != nil {
		log.Error("统计子账户失败", err)
		if err == reterror.MSG_6017 {
			return reterror.ErrModel{Err: reterror.MSG_6017, Code: reterror.Acc_Manage_6017}
		}
		return reterror.ErrModel{Err: reterror.SYSTEM_ERROR, Code: reterror.Failed}
	}
	data := make(map[string]interface{})
	data["Balance"] = balance
	data["Amount"] = amount
	return reterror.ErrModel{Err: nil, Code: reterror.Success, Data: data}
}

// AddressList 主链币对应的账户列表
// coinId 币种id （tag 别名 address 地址）模糊查询 （start ,end） 分页
func AddressList(qt *Qaddress) reterror.ErrModel {
	coinId := qt.CoinID
	condition := qt.Condition
	page := qt.Page
	tp := qt.Type
	if coinId == 0 {
		return reterror.ErrModel{Code: reterror.ParamsNil, Err: reterror.PARAMS_NULL}
	}
	am := &db.AccManageService{}
	list, total, err := am.AddressList(coinId, condition, page, tp)
	if err != nil {
		log.Error("get address list error", err)
		return reterror.ErrModel{Err: reterror.SYSTEM_ERROR, Code: reterror.Failed}
	}
	data := make(map[string]interface{})
	data["rows"] = list
	data["total"] = total
	return reterror.ErrModel{Err: nil, Code: reterror.Success, Data: data}
}

// AccountDetail 查询账户详细信息
func AccountDetail(id, coinId int) reterror.ErrModel {
	if id == 0 || coinId == 0 {
		return reterror.ErrModel{Code: reterror.ParamsNil, Err: reterror.PARAMS_NULL}
	}
	am := &db.AccManageService{}
	cd := &db.CoinService{}
	coin, err := cd.GetCoinById(coinId)
	if err != nil {
		return reterror.ErrModel{Err: reterror.SYSTEM_ERROR, Code: reterror.Failed}
	}
	detail, err := am.FindAddressById(id)
	if err != nil {
		return reterror.ErrModel{Err: reterror.MSG_6016, Code: reterror.Acc_Manage_6016}
	}
	coinName := coin.Name
	var sign bccore.BlockChainSign
	switch strings.ToUpper(coinName) {
	case "ETH":
		sign = bccore.STR_ETH
		break
	case "BTC":
		sign = bccore.STR_BTC
		break
	case "ERC20":
		sign = bccore.STR_ERC20
		break
	case "USDT":
		sign = bccore.STR_USDT
		break
	}
	if sign == "" {
		return reterror.ErrModel{Code: reterror.ParamsNil, Err: reterror.PARAMS_NULL}
	}
	walletService := walletClient.NewAppServer()
	balance, err := walletService.GetBalance(sign, detail.Address, coin.TokenAddress)
	if err != nil {
		return reterror.ErrModel{Code: reterror.WallertFailed, Err: reterror.MSG_1009}
	}
	data := make(map[string]interface{})
	data["Total"] = balance.String()
	data["Tag"] = detail.Tag
	data["Address"] = detail.Address
	return reterror.ErrModel{Err: nil, Code: reterror.Success, Data: data}
}

// SetTag 为账户设置别名
func SetTag(id int, tag string) reterror.ErrModel {
	if id == 0 {
		return reterror.ErrModel{Code: reterror.ParamsNil, Err: reterror.PARAMS_NULL}
	}
	am := &db.AccManageService{}
	// 获取地址信息
	addressInfo, err := am.FindAddressById(id)
	if err != nil {
		log.Error("获取地址信息失败", id, err)
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}

	if addressInfo == nil {
		log.Error("没有获取到地址信息", id)
		return reterror.ErrModel{Code: reterror.Acc_Manage_6016, Err: reterror.MSG_6016}
	}

	// 主账户地址不能修改
	if addressInfo.Type == common.MainAddressType {
		log.Debug("主账户地址无法修改tag")
		return reterror.ErrModel{Code: reterror.Acc_Manage_6018, Err: reterror.MSG_6018}
	}
	am.SetTag(id, tag)
	return reterror.ErrModel{Err: nil, Code: reterror.Success}
}

// TransferRecord  查询账户的转账记录列表
func TransferRecord(id int) reterror.ErrModel {
	if id == 0 {
		return reterror.ErrModel{Code: reterror.ParamsNil, Err: reterror.PARAMS_NULL}
	}
	am := &db.AccManageService{}
	list, err := am.TransferRecord(id)
	if err != nil {
		log.Error("get transfer record error", err)
		return reterror.ErrModel{Err: reterror.SYSTEM_ERROR, Code: reterror.Failed}
	}
	return reterror.ErrModel{Err: nil, Code: reterror.Success, Data: list}
}

// AddressTopTen  查询指定币种的账户金额top10
func AddressTopTen(qt *Qaddress) reterror.ErrModel {
	coinId := qt.CoinID
	if coinId == 0 {
		return reterror.ErrModel{Code: reterror.ParamsNil, Err: reterror.PARAMS_NULL}
	}
	am := &db.AccManageService{}
	list, total, err := am.AddressTopTen(coinId)
	if err != nil {
		log.Error("get assets top 10 error", err)
		return reterror.ErrModel{Err: reterror.SYSTEM_ERROR, Code: reterror.Failed}
	}
	data := make(map[string]interface{})
	data["rows"] = list
	data["total"] = total
	return reterror.ErrModel{Err: nil, Code: reterror.Success, Data: data}
}

// TokenRecord 代币转账记录
func TokenRecord(id, coinId int) reterror.ErrModel {
	if coinId == 0 {
		return reterror.ErrModel{Code: reterror.ParamsNil, Err: reterror.PARAMS_NULL}
	}
	am := &db.AccManageService{}
	list, err := am.TokenRecord(id, coinId)
	if err != nil {
		log.Error("get token record error", err)
		return reterror.ErrModel{Err: reterror.SYSTEM_ERROR, Code: reterror.Failed}
	}
	return reterror.ErrModel{Err: nil, Code: reterror.Success, Data: list}
}

// TokenList 代币明细
func TokenList(id int) reterror.ErrModel {
	if id == 0 {
		return reterror.ErrModel{Code: reterror.ParamsNil, Err: reterror.PARAMS_NULL}
	}
	am := &db.AccManageService{}
	list, err := am.TokenList(id)
	if err != nil {
		log.Error("get token list error", err)
		return reterror.ErrModel{Err: reterror.SYSTEM_ERROR, Code: reterror.Failed}
	}
	return reterror.ErrModel{Err: nil, Code: reterror.Success, Data: list}
}

// AddChildAccount 添加子账户
// coinId 币种id amount 数量 tag别名
func AddChildAccount(ct *Qaddress, account string) reterror.ErrModel {
	accLock.Lock()
	defer accLock.Unlock()
	coinId := ct.CoinID
	amount := ct.Amount
	tag := ct.Tag
	psw := ct.PSW
	if coinId == 0 || amount == 0 || psw == "" {
		return reterror.ErrModel{Code: reterror.ParamsNil, Err: reterror.PARAMS_NULL}
	}
	am := &db.AccManageService{}
	//校验数量
	childCount, otherCount := am.ChildStatistical(coinId)
	count := MaxCount - childCount - otherCount
	if amount > count {
		log.Error("数量超过可添加子账号数量", count)
		return reterror.ErrModel{Code: reterror.Acc_Manage_6014, Err: reterror.MSG_6014}
	}
	//校验密码
	verifyResult, _, errorCode := verify.VerifyPSW(account, psw)
	if verifyResult.Result == false {
		return reterror.ErrModel{Code: errorCode, Err: errors.New(verifyResult.Reason), Data: map[string]interface{}{"data": verifyResult.Data}}
	}
	var sign bccore.BlockChainSign
	var kt bccore.BloclChainType
	cs := &db.CoinService{}
	coin, err := cs.GetCoinById(coinId)
	if err != nil {
		log.Error("get coin info error", err)
		return reterror.ErrModel{Err: reterror.SYSTEM_ERROR, Code: reterror.Failed}
	}
	coinName := coin.Name

	switch strings.ToUpper(coinName) {
	case "ETH":
		sign = bccore.STR_ETH
		kt = bccore.BC_ETH
		break
	case "BTC":
		sign = bccore.STR_BTC
		kt = bccore.BC_BTC
		break
	case "ERC20":
		sign = bccore.STR_ETH
		kt = bccore.BC_ETH
		break
	case "LTC":
		sign = bccore.STR_LTC
		kt = bccore.BC_LTC
		break
	}

	if sign == "" {
		return reterror.ErrModel{Code: reterror.Acc_Manage_6021, Err: reterror.MSG_6021}
	}
	walletService := walletClient.NewAppServer()
	keys, err := walletService.GeneraterKeys(amount, sign)
	if err != nil {
		log.Error("请求wallet获取地址出错", err)
		return reterror.ErrModel{Err: reterror.SYSTEM_ERROR, Code: reterror.Failed}
	}

	adds := make([]string, len(keys))
	deeps := make([]string, len(keys))
	for i, key := range keys {
		adds[i] = key.Address()
		deep := GetBtcKey(uint32(kt), key.CurrentNum(), key.CustomDeep())
		json, _ := json.Marshal(deep)
		deeps[i] = string(json)
	}
	if coinId == 0 || amount == 0 || psw == "" {
		return reterror.ErrModel{Code: reterror.ParamsNil, Err: reterror.PARAMS_NULL}
	}
	num := am.AddChildAccount(coinId, amount, tag, adds, deeps)
	if num == 1 {
		return reterror.ErrModel{Err: nil, Code: reterror.Success}
	}
	return reterror.ErrModel{Err: reterror.MSG_6008, Code: reterror.Acc_Manage_6008}
}

//组织deep
func GetBtcKey(kt uint32, curNum uint32, customDeep []uint32) []uint32 {

	var deeps []uint32
	deeps = append(deeps, uint32(kt))
	deeps = append(deeps, customDeep...)
	deeps = append(deeps, curNum)
	fmt.Println("address11-----", deeps)
	return deeps
}

//合并账户
func MergeAccount(coinId, mold int, ids []int, psw, account string) reterror.ErrModel {
	accLock.Lock()
	defer accLock.Unlock()
	if db.ConfigMap["combine_account"] == "1" {
		log.Error("账户正在合并中")
		return reterror.ErrModel{Code: reterror.Acc_Manage_6020, Err: reterror.MSG_6020}
	}
	if coinId == 0 {
		return reterror.ErrModel{Code: reterror.ParamsNil, Err: reterror.PARAMS_NULL}
	}
	//TODO:账户合并目前仅支持ETH及其代币币种
	cs := &db.CoinService{}
	tx := &db.TxinfosService{}
	coin, err := cs.GetCoinById(coinId)
	ETH, err := cs.GetETH()
	if err != nil {
		log.Error("合并子账户", err)
		return reterror.ErrModel{Err: reterror.SYSTEM_ERROR, Code: reterror.Failed}
	}
	coinName := coin.Name

	fmt.Println(coin.Name, coin.TokenType)
	if strings.ToUpper(coinName) != "ETH" && coin.TokenType != ETH.ID {
		return reterror.ErrModel{Err: reterror.MSG_6017, Code: reterror.Acc_Manage_6017}
	}
	if mold == 1 && len(ids) == 0 {
		return reterror.ErrModel{Code: reterror.ParamsNil, Err: reterror.PARAMS_NULL}
	}
	//verifyResult, _, errorCode:= verify.VerifyPSW(account, psw)
	//if verifyResult.Result == false {
	//	return reterror.ErrModel{Code: errorCode, Err: errors.New(verifyResult.Reason), Data: map[string]interface{}{
	//	"data": verifyResult.Data}}
	//}
	am := &db.AccManageService{}
	result := am.MergeAccount(coinId, mold, ids)
	if result == common.MergeAccountAddressNotFound {
		return reterror.ErrModel{Code: reterror.Acc_Manage_6009, Err: reterror.MSG_6009}
	}
	if result == common.MergeAccountSuccess {
		tx.UpadateConfigs("combine_account", "1")
		db.ConfigMap["combine_account"] = "1"
		return reterror.ErrModel{Code: reterror.Success}
	}
	if result == common.MergeAccountTransfering {
		return reterror.ErrModel{Code: reterror.Acc_Manage_6012, Err: reterror.MSG_6012}
	}
	if result == common.MergeAccountNoAddrToMegge {
		return reterror.ErrModel{Code: reterror.Acc_Manage_6011, Err: reterror.MSG_6011}
	}
	if result == common.MergeAccountSqlErr {
		return reterror.ErrModel{Code: reterror.Auth_5006, Err: reterror.MSG_5006}
	}
	return reterror.ErrModel{Code: reterror.Acc_Manage_6010, Err: reterror.MSG_6010}
}

// GetContractAddress 获取合约账户地址
func GetContractAddress() reterror.ErrModel {
	am := &db.AccManageService{}
	msg, err := am.ContractAddr()
	if err != nil {
		if err.Error() == "Voucher Return Fail" {
			return reterror.ErrModel{Code: reterror.VoucherFail, Err: reterror.MSG_1006}
		} else {
			log.Error("Get Contract Address Error", err)
			return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
		}
	}

	// 获取对地址的签名信息
	addressSign := db.ConfigMap["address_sign"]
	if addressSign == "" {
		log.Error("没有对地址签名")
		return reterror.ErrModel{Code: reterror.Acc_Manage_6016, Err: reterror.MSG_6016}
	}

	signs := &db.AddressSigns{}
	err = json.Unmarshal([]byte(addressSign), &signs)
	if err != nil {
		log.Error("地址签名解析json错误", err)
		return reterror.ErrModel{Code: reterror.ParamsNil, Err: reterror.PARAMS_NULL}
	}
	for i := 0; i < len(signs.AddressInfo); i++ {
		if strings.ToUpper(signs.AddressInfo[i].CoinName) == "ETH" {
			msg.Sign = signs.AddressInfo[i].Sign
			break
		}
	}
	msg.Account = signs.Account

	return reterror.ErrModel{Code: reterror.Success, Data: msg}
}

// GenContractAddress 创建合约账户
func GenContractAddress() reterror.ErrModel {
	accLock.Lock()
	defer accLock.Unlock()
	am := &db.AccManageService{}
	// 主账户余额小于0.01则无法创建合约地址
	asset, err := am.MainAccAssets()
	if err != nil {
		log.Error("Get Main Address Aset Error", err)
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	asset_decimal, _ := decimal.NewFromString(asset)
	lowest_asset := decimal.NewFromFloat(0.01)
	if asset_decimal.LessThan(lowest_asset) == true {
		log.Error("Asset less than 0.01E", asset_decimal.String())
		return reterror.ErrModel{Code: reterror.Acc_Manage_6013, Err: reterror.MSG_6013}
	}
	// 向签名机申请创建合约地址
	err = am.ApplyForContract()
	if err != nil {
		if err.Error() == "Voucher Return Error" {
			return reterror.ErrModel{Code: reterror.VoucherFail, Err: reterror.MSG_1006}
		} else if err.Error() == "Duplicate Init Contract" {
			// 重复创建合约
			return reterror.ErrModel{Code: reterror.Template_9008, Err: reterror.MSG_9008}
		} else {
			return reterror.ErrModel{Code: reterror.Template_9007, Err: reterror.MSG_9007}
		}
	}
	return reterror.ErrModel{Code: reterror.Success}
}
