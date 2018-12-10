// Copyright 2018. bolaxy.org authors.
// 
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// 
// 		 http://www.apache.org/licenses/LICENSE-2.0
// 
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package timer

import (
	"encoding/json"
	"strconv"
	"time"
	"strings"
	log "github.com/alecthomas/log4go"
	"github.com/boxproject/apiServer/db"
	voucher "github.com/boxproject/apiServer/rpc"
	msg "github.com/boxproject/apiServer/service/message"
	"github.com/boxproject/apiServer/utils"
	"github.com/boxproject/boxwallet/bccoin"
	"github.com/boxproject/boxwallet/bccore"
	walletError "github.com/boxproject/boxwallet/errors"
	txutil "github.com/boxproject/boxwallet/signature"
	"github.com/shopspring/decimal"
	"github.com/boxproject/apiServer/config"
	"github.com/boxproject/apiServer/common"
)

//所有btc和LTC的地址  如AllBTCAddress【"BTC"】
var AllBTCAddress = make(map[string][]string)

//所有BTC 地址的索引
var AllBTCAddressMap = make(map[string]string)

//所有LTC 地址的索引
var AllLTCAddressMap = make(map[string]string)

//是否占用
var isUsing = utils.MyMap{
	IsUsing: make(map[string]int),
}

func getAllmasterAddress() {
	timerLock.Lock()
	defer timerLock.Unlock()
	if db.AllMasterAddress == nil {
		addressService := &db.AccManageService{}
		addressDB := addressService.GetMasterAddress()
		var BTCAddresses []string
		var LTCAddresses []string
		if len(addressDB) > 0 {
			db.AllMasterAddress = make(map[string]string)
			for i := 0; i < len(addressDB); i++ {
				if addressDB[i].Type == 0 {
					db.AllMasterAddress[addressDB[i].CoinName] = addressDB[i].Address
				}
				if addressDB[i].CoinName == "BTC" {
					BTCAddresses = append(BTCAddresses, addressDB[i].Address)
					AllBTCAddressMap[addressDB[i].Address] = addressDB[i].Deep
				}
				if addressDB[i].CoinName == "LTC" {
					LTCAddresses = append(LTCAddresses, addressDB[i].Address)
					AllLTCAddressMap[addressDB[i].Address] = addressDB[i].Deep
				}

			}
			AllBTCAddress["BTC"] = BTCAddresses
			AllBTCAddress["LTC"] = LTCAddresses

			log.Debug("同步公钥结束")
			return
		} else {
			if voucher.GetVoucherStatus().ServerStatus < 2 {
				log.Debug("状态小于2", voucher.GetVoucherStatus().ServerStatus)
				return
			}
			oper := &voucher.GrpcServer{
				Type: voucher.VOUCHER_OPERATE_MASTER_PUBKEY,
			}
			_, res, ret := voucher.SendVoucherData(oper)
			if ret != 2 {
				log.Error("签名机获取主公钥错误")
				return
			} else {
				coinsService := &db.CoinService{}
				log.Debug("voucherPUblickey", res.PublicKeys)
				walletErr := Wallet.SaveMasterKey(res.PublicKeys)
				if walletErr != nil {
					log.Error("保存主公钥错误", walletErr)
					return
				}
				coins, _ := coinsService.GetCoins()
				addressEntity := &db.Address{}
				//生成的主地址
				masterAddress, err2 := Wallet.GetMasterAddress(bccore.BC_ETH)
				if err2 != nil {
					log.Error("保存主公钥错误", err2)
					return
				}

				db.AllMasterAddress = make(map[string]string)
				for i := 0; i < len(coins); i++ {
					if coins[i].Name == "ETH" {
						masterAddress, err2 = Wallet.GetMasterAddress(bccore.BC_ETH)
						if err2 != nil {
							log.Error("获取ETH地址错误", err2)
							continue
						}
						addressEntity = &db.Address{
							Address:  masterAddress.Address(),
							Type:     0,
							CoinName: coins[i].Name,
							CoinId:   coins[i].ID,
						}

					} else if coins[i].Name == "USDT" {
						masterAddress, err2 = Wallet.GetMasterAddress(bccore.BC_USDT)
						if err2 != nil {
							log.Error("获取USDT地址错误", err2)
							continue
						}
						addressEntity = &db.Address{
							Address:  masterAddress.Address(),
							Type:     0,
							CoinName: coins[i].Name,
							CoinId:   coins[i].ID,
						}
					} else if coins[i].Name == "BTC" {
						masterAddress, err2 = Wallet.GetMasterAddress(bccore.BC_BTC)

						if err2 != nil {
							log.Error("获取BTC地址错误", err2)
							continue
						}
						BTCAddresses = append(BTCAddresses, masterAddress.Address())
						addressEntity = &db.Address{
							Address:  masterAddress.Address(),
							Type:     0,
							CoinName: coins[i].Name,
							CoinId:   coins[i].ID,
						}
					} else if coins[i].Name == "LTC" {
						masterAddress, err2 = Wallet.GetMasterAddress(bccore.BC_LTC)
						if err2 != nil {
							log.Error("获取LTC地址错误", err2)
							continue
						}
						LTCAddresses = append(LTCAddresses, masterAddress.Address())
						addressEntity = &db.Address{
							Address:  masterAddress.Address(),
							Type:     0,
							CoinName: coins[i].Name,
							CoinId:   coins[i].ID,
						}
					} else {
						continue
					}
					addressService.SaveMasterAddress(addressEntity)
					db.AllMasterAddress[coins[i].Name] = masterAddress.Address()

					AllBTCAddress["BTC"] = BTCAddresses
					AllBTCAddress["LTC"] = LTCAddresses

				}
			}
		}
	} else {
		Time2.Stop()
		masterRet = false
	}

}

func transferChan() {
	go func() {
		for x := range db.BTCTransferCh {
			time.Sleep(time.Second)
			btcusdtBusiness(&x)
		}
	}()
	go func() {
		for x := range db.LTCTransferCh {
			time.Sleep(time.Second)
			btcusdtBusiness(&x)
		}
	}()
	go func() {
		for x := range db.ETHTransferCh {
			time.Sleep(time.Second)
			hashChain()
			ethBusiness(&x)
		}
	}()

	go func() {
		for x := range db.MergeETHTransferCh {
			time.Sleep(time.Second)
			mergeBusiness(&x)
		}
	}()

}

func backChannel(transfer *db.TransferSend) {
	//go func(){
	log.Info("重新塞回去", transfer.Amount)

	if transfer.Currency == "ETH" || (transfer.Currency == "ERC20") {
		db.ETHTransferCh <- *transfer
	} else if (transfer.CoinName == "BTC") || (transfer.CoinName == "USDT") {
		db.BTCTransferCh <- *transfer
	} else if transfer.CoinName == "LTC" {
		db.LTCTransferCh <- *transfer
	}
	//}()

}

func hashChain() {
	//模版上链操作
	temchain := db.TemplateChains.GetTemplateChain()
	if len(temchain) > 0 {
		uploadTemplate(temchain[0])
		db.TemplateChains.DelTemplateChain()
	}
}

func btcusdtBusiness(transfer *db.TransferSend) {

	if transfer.CoinName == "USDT" {
		usdtBusiness(transfer)
	} else if transfer.CoinName == "LTC" || (transfer.CoinName == "BTC") {
		btcBusiness(transfer)
	}
}

//btc或者ltc
func btcBusiness(transfer *db.TransferSend) {
	ts := &db.TransferDBService{}
	var err error
	var fromAddress = AllBTCAddress[transfer.CoinName]
	AllMasterAddress := db.AllMasterAddress
	//组装map
	addMap := make(map[string][]uint32)
	addMap[AllMasterAddress[transfer.CoinName]] = nil
	for i := 0; i < len(AllBTCAddress[transfer.CoinName]); i++ {
		maps, mapsError := transferMap(AllLTCAddressMap[AllBTCAddress[transfer.CoinName][i]])
		if transfer.CoinName == "BTC" {
			maps, mapsError = transferMap(AllBTCAddressMap[AllBTCAddress[transfer.CoinName][i]])
		}
		if mapsError != nil {
			log.Error("索引解析错误", mapsError)
			ts.UpdateOrderExpire(&db.Transfer{OrderId: transfer.OrderId, Status: common.TransferAddressIndexError})
			return
		}
		addMap[AllBTCAddress[transfer.CoinName][i]] = maps

	}

	miner, err := strconv.ParseFloat(transfer.Miner, 64)
	if err != nil {
		log.Error("矿工费转化float错误", err)
		ts.UpdateOrderExpire(&db.Transfer{OrderId: transfer.OrderId, Status: common.TransferMinerToFloatError})
		return
	}
	//余额不足，重新转账
	if transfer.Status == common.TransferInsufficientBalance {
		if transfer.Deadline != "" {
			now := time.Now()
			formattime, _ := time.ParseInLocation("2006-01-02 15:04:05", transfer.Deadline, time.Local)
			if now.After(formattime) {
				//转账过期
				ts.UpdateOrderExpire(&db.Transfer{OrderId: transfer.OrderId, Status: common.TransferApproveExpire})
				//插入站内信
				msg.TransferFailOutOfDate(transfer.CoinName, transfer.Amount, transfer.OrderId, transfer.ApplyerId)
				return
			}
		}
	}
	var transferMsg ApplyMsg

	err = json.Unmarshal([]byte(transfer.TransferMsg), &transferMsg)
	if err != nil {
		log.Error("msg解析错误")
		ts.UpdateOrderExpire(&db.Transfer{OrderId: transfer.OrderId, Status: common.TransferDecodeMsgError})
		return
	}
	bc_bc := bccore.BC_BTC
	bc_str := bccore.STR_BTC
	if transfer.CoinName == "LTC" {
		bc_bc = bccore.BC_LTC
		bc_str = bccore.STR_LTC
	}

	var amounts []*bccoin.AddressAmount
	for i := 0; i < len(transferMsg.ApplyVos); i++ {
		trans := transferMsg.ApplyVos[i]
		//decimalAmounts,_ := decimal.NewFromString("0")
		for j := 0; j < len(trans.Amount); j++ {
			decimalAmount, err := decimal.NewFromString(trans.Amount[j])
			if err != nil {
				log.Error("数据解析错误", err)
				ts.UpdateOrderExpire(&db.Transfer{OrderId: transfer.OrderId, Status: common.TransferDecodeDataError})
				return
			}
			coinAmount, err := bccoin.NewCoinAmount(bc_bc, "", decimalAmount.String())
			amount := &bccoin.AddressAmount{
				Address: trans.ToAddress,
				Amount:  coinAmount,
			}
			amounts = append(amounts, amount)
			//decimalAmounts = decimal.Sum(decimalAmounts,decimalAmount)
		}
		//coinAmount, err := bccoin.NewCoinAmount(bc_bc, "", decimalAmounts.String())
		if err != nil {
			log.Error("精度不够", err)
			ts.UpdateOrderExpire(&db.Transfer{OrderId: transfer.OrderId, Status: common.TransferDicemalError})
			return
		}
		//amount := &bccoin.AddressAmount{
		//	Address: trans.ToAddress,
		//	Amount:  coinAmount,
		//}
		//amounts = append(amounts, amount)
		fromAddress = removeArr(fromAddress, trans.ToAddress)
	}

	//构造交易
	uuid, txu := createTxBTC(bc_str, fromAddress, "", amounts, miner, transfer)
	if txu == nil {
		return
	}
	//签名
	txu, voucherRet := voucherSign(transfer, txu, addMap, string(bc_str))
	if !voucherRet {
		ts.UpdateOrderExpire(&db.Transfer{OrderId: transfer.OrderId, Status: common.TransferVoucherReturnError})
		return
	}
	//更新转账状态为中间状态3
	ts.UpdateOrderExpire(&db.Transfer{OrderId: transfer.OrderId, Status: common.TransferSomeSuccess,TxId:txu.TxId()})
	// 插入站内信，部分成功
	orderInfo, _ := ts.FindOrderById(transfer.OrderId)
	var ms db.MessageService
	ms.TransferPartialSuccessful(orderInfo.CoinName, orderInfo.Amount, orderInfo.ID, orderInfo.ApplyerId)
	//上链操作
	sendTxBTC(bc_str, txu, uuid, transfer)

}

//usdt
func usdtBusiness(transfer *db.TransferSend) {
	defer func() {
		go asycOrderStatus(transfer.OrderId)
	}()
	ts := &db.TransferDBService{}
	var err error
	var fromAddress = []string{}
	AllMasterAddress := db.AllMasterAddress
	//组装map
	addMap := make(map[string][]uint32)
	addMap[AllMasterAddress["USDT"]] = nil
	miner, err := strconv.ParseFloat(transfer.Miner, 64)
	if err != nil {
		log.Error("矿工费转化float错误", err)
		ts.UpdateTransfer(&db.Transfer{ID: transfer.TransferId, Status: common.TransferMinerToFloatError})
		return
	}
	//余额不足，重新转账
	if transfer.Status == common.TransferInsufficientBalance {
		if transfer.Deadline != "" {
			now := time.Now()
			formattime, _ := time.ParseInLocation("2006-01-02 15:04:05", transfer.Deadline, time.Local)
			if now.After(formattime) {
				//转账过期
				ts.UpdateTransfer(&db.Transfer{ID: transfer.TransferId, Status: common.TransferApproveExpire})
				////插入站内信
				//msg.TransferFailOutOfDate(transfer.CoinName, transfer.Amount, transfer.TransferId, transfer.ApplyerId)
				return
			}
		}
	}
	fromAddress = append(fromAddress, AllMasterAddress["USDT"])
	var amounts []*bccoin.AddressAmount
	bc_bc := bccore.BC_USDT
	bc_str := bccore.STR_USDT
	coinAmount, _ := bccoin.NewCoinAmount(bc_bc, bccore.Token(USDT_PPID), transfer.Amount)
	amount := &bccoin.AddressAmount{
		Address: transfer.ToAddress,
		Amount:  coinAmount,
	}
	if coinAmount == nil {
		//精度不够
		ts.UpdateTransfer(&db.Transfer{ID: transfer.TransferId, Status: common.TransferDicemalError})
		return
	}
	amounts = append(amounts, amount)
	//构造交易
	uuid, txu := createTx(bc_str, fromAddress, transfer.TokenAddress, amounts, miner, transfer)
	if txu == nil {
		return
	}
	//签名
	txu, voucherRet := voucherSign(transfer, txu, addMap, string(bc_str))
	if !voucherRet {
		ts.UpdateTransfer(&db.Transfer{ID: transfer.TransferId, Status: common.TransferVoucherReturnError})
		return
	}
	//更新转账状态为中间状态3
	ts.UpdateTransfer(&db.Transfer{ID: transfer.TransferId, Status: common.TransferSomeSuccess,TxId:txu.TxId()})
	//上链操作
	sendTx(bc_str, txu, uuid, transfer)
}

//eth
func ethBusiness(transfer *db.TransferSend) {
	//为了促发模板上链
	if transfer.Types == 3 {
		return
	}
	ts := &db.TransferDBService{}
	var err error
	var fromAddress = []string{}
	AllMasterAddress := db.AllMasterAddress
	coinType := transfer.CoinName
	//正在使用的地址或者币种 ，如果是eth地址则填地址
	if transfer.TokenAddress != "" {
		coinType = "ERC20"
	}
	defer func() {
		go asycOrderStatus(transfer.OrderId)
	}()
	//组装map
	addMap := make(map[string][]uint32)
	addMap[strings.ToLower(AllMasterAddress["ETH"])] = nil
	miner, err := strconv.ParseFloat(transfer.Miner, 64)
	if err != nil {
		log.Error("矿工费转化float错误", err)
		ts.UpdateTransfer(&db.Transfer{ID: transfer.TransferId, Status: common.TransferMinerToFloatError})
		return
	}

	//余额不足，重新转账
	if transfer.Status == common.TransferInsufficientBalance {
		if transfer.Deadline != "" {
			now := time.Now()
			formattime, _ := time.ParseInLocation("2006-01-02 15:04:05", transfer.Deadline, time.Local)
			if now.After(formattime) {
				//转账过期
				ts.UpdateTransfer(&db.Transfer{ID: transfer.TransferId, Status: common.TransferApproveExpire})
				////插入站内信
				//msg.TransferFailOutOfDate(transfer.CoinName, transfer.Amount, transfer.TransferId, transfer.ApplyerId)
				return
			}
		}
	}
	fromAddress = append(fromAddress, AllMasterAddress["ETH"])
	var amounts []*bccoin.AddressAmount
	bc_bc := bccore.BC_ETH
	bc_str := bccore.STR_ETH
	coinAmount, _ := bccoin.NewCoinAmount(bc_bc, bccore.Token(transfer.TokenAddress), transfer.Amount)

	if coinType == "ERC20" {
		bc_bc = bccore.BC_ERC20
		bc_str = bccore.STR_ERC20
		coinAmount, _ = bccoin.NewCoinAmount(bc_bc, bccore.Token(transfer.TokenAddress), transfer.Amount)
	}

	amount := &bccoin.AddressAmount{
		Address: transfer.ToAddress,
		Amount:  coinAmount,
	}

	if coinAmount == nil {
		//精度不够
		ts.UpdateTransfer(&db.Transfer{ID: transfer.TransferId, Status: common.TransferDicemalError})
		return
	}
	amounts = append(amounts, amount)

	//构造交易
	uuid, txu := createTx(bc_str, fromAddress, transfer.TokenAddress, amounts, miner, transfer)
	if txu == nil {
		return
	}
	//签名
	txu, voucherRet := voucherSign(transfer, txu, addMap, string(bc_str))
	if !voucherRet {
		ts.UpdateTransfer(&db.Transfer{ID: transfer.TransferId, Status: common.TransferVoucherReturnError})
		return
	}
	//更新转账状态为中间状态3
	ts.UpdateTransfer(&db.Transfer{ID: transfer.TransferId, Status: common.TransferSomeSuccess,TxId:txu.TxId()})
	//上链操作
	sendTx(bc_str, txu, uuid, transfer)

}

//eth合并
func mergeBusiness(transfer *db.TransferSend) {

	defer func() {
		go asycCombineOrder(transfer.OrderId)
	}()
	miner := 2.0
	ts := &db.TransferDBService{}
	var err error
	var fromAddress = []string{}

	//正在使用的地址或者币种 ，如果是eth地址则填地址

	//组装索引
	addMap := make(map[string][]uint32)
	addresses := &TransferMsgJson{}

	err = json.Unmarshal([]byte(transfer.AddressMsg), addresses)
	if err != nil {
		log.Error("地址索引解析错误", err)
		ts.UpdateTransfer(&db.Transfer{ID: transfer.TransferId, Status: common.TransferAddressIndexError})
		return
	}
	transfer.TokenAddress = addresses.Token
	addMap[strings.ToLower(addresses.FromAddress.Address)] = addresses.FromAddress.Deep
	addMap[strings.ToLower(addresses.ToAddress.Address)] = addresses.ToAddress.Deep
	fromAddress = append(fromAddress, addresses.FromAddress.Address)
	var amounts []*bccoin.AddressAmount
	bc_bc := bccore.BC_ETH
	bc_str := bccore.STR_ETH
	coinAmount, _ := bccoin.NewCoinAmount(bc_bc, "", transfer.Amount)

	if transfer.TokenAddress != "" {
		bc_bc = bccore.BC_ERC20
		bc_str = bccore.STR_ERC20
		coinAmount, _ = bccoin.NewCoinAmount(bc_bc, bccore.Token(transfer.TokenAddress), transfer.Amount)
	}

	amount := &bccoin.AddressAmount{
		Address: addresses.ToAddress.Address,
		Amount:  coinAmount,
	}

	if coinAmount == nil {
		//精度不够
		ts.UpdateTransfer(&db.Transfer{ID: transfer.TransferId, Status: common.TransferDicemalError})
		return
	}
	amounts = append(amounts, amount)

	//构造交易
	uuid, txu := createTx(bc_str, fromAddress, transfer.TokenAddress, amounts, miner, transfer)
	if txu == nil {
		return
	}
	//签名
	txu, voucherRet := voucherSign(transfer, txu, addMap, string(bc_str))
	if !voucherRet {
		ts.UpdateTransfer(&db.Transfer{ID: transfer.TransferId, Status: common.TransferVoucherReturnError})
		return
	}
	//更新转账状态为中间状态3
	ts.UpdateTransfer(&db.Transfer{ID: transfer.TransferId, Status: common.TransferSomeSuccess})
	//上链操作
	sendTx(bc_str, txu, uuid, transfer)


}

func transferBusiness(transfer *db.TransferSend) {
	log.Debug("定时转账，开始------", transfer.Amount)

	//为了促发模板上链
	if transfer.Types == 3 {
		return
	}
	//默认矿工费
	miner := 2.0
	var fromAddress = []string{}
	AllMasterAddress := db.AllMasterAddress
	ts := &db.TransferDBService{}
	var err error

	coinType := transfer.CoinName
	//正在使用的地址或者币种 ，如果是eth地址则填地址，如果是其他币种就是币种名称
	using := transfer.CoinName
	if transfer.TokenAddress != "" {
		coinType = "ERC20"
	}
	if transfer.CoinName == "USDT" {
		coinType = "USDT"
	}

	//组装索引
	addMap := make(map[string][]uint32)
	addresses := &TransferMsgJson{}
	if transfer.Types == 2 {

		errs := json.Unmarshal([]byte(transfer.AddressMsg), addresses)
		if errs != nil {
			log.Error("地址索引解析错误", errs)
			ts.UpdateTransfer(&db.Transfer{ID: transfer.TransferId, Status: common.TransferAddressIndexError})
			return
		}
		transfer.TokenAddress = addresses.Token
		addMap[strings.ToLower(addresses.FromAddress.Address)] = addresses.FromAddress.Deep
		addMap[strings.ToLower(addresses.ToAddress.Address)] = addresses.ToAddress.Deep
		using = addresses.FromAddress.Address
		if addresses.Token != "" {
			coinType = "ERC20"
		}
	} else {
		if coinType == "ERC20" || (coinType == "ETH") {
			addMap[strings.ToLower(AllMasterAddress["ETH"])] = nil
			using = AllMasterAddress["ETH"]
		} else if coinType == "BTC" {
			for i := 0; i < len(AllBTCAddress["BTC"]); i++ {
				maps, mapsError := transferMap(AllBTCAddressMap[AllBTCAddress["BTC"][i]])
				if mapsError != nil {
					log.Error("索引解析错误", mapsError)
					ts.UpdateTransfer(&db.Transfer{ID: transfer.TransferId, Status: common.TransferAddressIndexError})
					return
				}
				addMap[AllBTCAddress["BTC"][i]] = maps
			}
		} else if coinType == "LTC" {
			for i := 0; i < len(AllBTCAddress["LTC"]); i++ {
				maps, mapsError := transferMap(AllLTCAddressMap[AllBTCAddress["LTC"][i]])
				if mapsError != nil {
					log.Error("索引解析错误")
					ts.UpdateTransfer(&db.Transfer{ID: transfer.TransferId, Status: common.TransferAddressIndexError})
					return
				}
				addMap[AllBTCAddress["LTC"][i]] = maps
			}
		} else {
			addMap[AllMasterAddress[transfer.CoinName]] = nil
		}

		miner, err = strconv.ParseFloat(transfer.Miner, 64)
		if err != nil {
			log.Error("矿工费转化float错误", err)
			ts.UpdateTransfer(&db.Transfer{ID: transfer.TransferId, Status: common.TransferMinerToFloatError})
			return
		}
	}

	if isUsing.Get(using) == 1 {
		log.Info("该地址正在占用")
		ts.UpdateTransfer(&db.Transfer{ID: transfer.TransferId, Status: common.TransferApprovaled})
		backChannel(transfer)
		return
	}
	isUsing.Set(using, 1)
	defer func() {
		isUsing.Set(using, 0)
		if transfer.Types == 1 {
			go asycOrderStatus(transfer.OrderId)
		} else if transfer.Types == 2 {
			go asycCombineOrder(transfer.OrderId)
		}

	}()

	//余额不足，重新转账
	if transfer.Status == common.TransferInsufficientBalance {
		if transfer.Deadline != "" {
			now := time.Now()
			formattime, _ := time.ParseInLocation("2006-01-02 15:04:05", transfer.Deadline, time.Local)
			if now.After(formattime) {
				//转账过期
				ts.UpdateTransfer(&db.Transfer{ID: transfer.TransferId, Status: common.TransferApproveExpire})
				////插入站内信
				//msg.TransferFailOutOfDate(transfer.CoinName, transfer.Amount, transfer.TransferId, transfer.ApplyerId)
				return
			}
		}
	}
	var amounts []*bccoin.AddressAmount
	bc_bc := bccore.BC_ETH
	bc_str := bccore.STR_ETH
	coinAmount, _ := bccoin.NewCoinAmount(bccore.BC_ETH, "", transfer.Amount)
	switch coinType {
	case "ETH":
		{
			bc_bc = bccore.BC_ETH
			bc_str = bccore.STR_ETH
			fromAddress = append(fromAddress, AllMasterAddress["ETH"])
			coinAmount, _ = bccoin.NewCoinAmount(bc_bc, "", transfer.Amount)
			break
		}
	case "BTC":
		{
			bc_bc = bccore.BC_BTC
			bc_str = bccore.STR_BTC
			fromAddress = removeArr(AllBTCAddress["BTC"], transfer.ToAddress)
			coinAmount, _ = bccoin.NewCoinAmount(bc_bc, "", transfer.Amount)
			break
		}
	case "ERC20":
		{
			bc_bc = bccore.BC_ERC20
			bc_str = bccore.STR_ERC20
			fromAddress = append(fromAddress, AllMasterAddress["ETH"])
			log.Error("haohaoxianghxiang", transfer.TokenAddress, transfer.Amount)
			coinAmount, err = bccoin.NewCoinAmount(bc_bc, bccore.Token(transfer.TokenAddress), transfer.Amount)
			if err != nil {
				log.Error("hohaoxiangxiang", err)
			}
			break
		}
	case "LTC":
		{
			bc_bc = bccore.BC_LTC
			bc_str = bccore.STR_LTC
			fromAddress = removeArr(AllBTCAddress["LTC"], transfer.ToAddress)
			coinAmount, _ = bccoin.NewCoinAmount(bc_bc, "", transfer.Amount)
			break
		}
	case "USDT":
		{
			bc_bc = bccore.BC_USDT
			bc_str = bccore.STR_USDT
			fromAddress = append(fromAddress, AllMasterAddress["USDT"])
			coinAmount, _ = bccoin.NewCoinAmount(bc_bc, bccore.Token(USDT_PPID), transfer.Amount)
			break
		}
	default:
		{
			ts.UpdateTransfer(&db.Transfer{ID: transfer.TransferId, Status: common.TransferCoinNotFound})
			log.Error("无此币种类别")
			return
		}
	}
	if transfer.Types == 2 {
		fromAddress = nil
		fromAddress = append(fromAddress, addresses.FromAddress.Address)
	}

	amount := &bccoin.AddressAmount{
		Address: transfer.ToAddress,
		Amount:  coinAmount,
	}

	if coinAmount == nil {
		//精度不够
		ts.UpdateTransfer(&db.Transfer{ID: transfer.TransferId, Status: common.TransferDicemalError})
		return
	}
	amounts = append(amounts, amount)

	//构造交易
	uuid, txu := createTx(bc_str, fromAddress, transfer.TokenAddress, amounts, miner, transfer)
	if txu == nil {
		return
	}
	//签名
	txu, voucherRet := voucherSign(transfer, txu, addMap, string(bc_str))
	if !voucherRet {
		ts.UpdateTransfer(&db.Transfer{ID: transfer.TransferId, Status: common.TransferVoucherReturnError})
		return
	}
	//更新转账状态为中间状态3
	ts.UpdateTransfer(&db.Transfer{ID: transfer.TransferId, Status: common.TransferSomeSuccess,TxId:txu.TxId()})
	//上链操作
	sendTx(bc_str, txu, uuid, transfer)

}

func createTx(bc_str bccore.BlockChainSign, fromAddress []string, token string, amounts []*bccoin.AddressAmount, miner float64, transfer *db.TransferSend) (string, txutil.TxUtil) {
	log.Info("构造交易的各个参数", bc_str, fromAddress, transfer.TokenAddress, miner)
	log.Info("amounts", amounts[0].Amount, amounts[0].Address)
	ts := &db.TransferDBService{}
	uuid, txu, err := Wallet.CreateTx(bc_str, fromAddress, transfer.TokenAddress, amounts, miner)
	if err != nil {
		if err == walletError.ERR_ADDRESS_QUEUE_BLOCKED {
			log.Error("阻塞", err)
			ts.UpdateTransfer(&db.Transfer{ID: transfer.TransferId, Status: common.TransferApprovaled})
			backChannel(transfer)
		} else if err == walletError.ERR_NOT_ENOUGH_COIN {
			log.Error("余额不足", err)
			if transfer.Types == 1 {
				ts.UpdateTransfer(&db.Transfer{ID: transfer.TransferId, Status: common.TransferInsufficientBalance})
				transfer.Status = common.TransferInsufficientBalance
				backChannel(transfer)
			} else {
				ts.UpdateTransfer(&db.Transfer{ID: transfer.TransferId, Status: common.TransferFailed})
			}
		} else {
			ts.UpdateTransfer(&db.Transfer{ID: transfer.TransferId, Status: common.TransferGenTxError})
			log.Error("构造交易失败", err)
		}
		return "", nil
	} else {
		return uuid, txu
	}
}

func sendTx(bc_str bccore.BlockChainSign, txu txutil.TxUtil, uuid string, transfer *db.TransferSend) {
	log.Debug("正在上链")
	err := Wallet.SendTx(bc_str, txu, uuid)
	ts := &db.TransferDBService{}
	if err != nil {
		log.Error("上链错误", err)

		//余额不足
		if err == walletError.ERR_NOT_ENOUGH_COIN {
			if transfer.Types == 1 {
				ts.UpdateTransfer(&db.Transfer{ID: transfer.TransferId, Status: common.TransferInsufficientBalance, TxId: txu.TxId()})
				transfer.Status = common.TransferInsufficientBalance
				backChannel(transfer)
			} else {
				ts.UpdateTransfer(&db.Transfer{ID: transfer.TransferId, Status: common.TransferFailed, TxId: txu.TxId()})
				//插入站内信
				//msg.TransferFail(transfer.CoinName, transfer.Amount, transfer.TransferId, transfer.ApplyerId)
			}

		} else if err == walletError.ERR_PIPELINE_DATA_ILLEGAL {
			ts.UpdateTransfer(&db.Transfer{ID: transfer.TransferId, Status: common.TransferSuccess, TxId: txu.TxId()})
			//if transfer.Currency == "ERC20" {
			//	updateBalances(transfer.CoinName,AllBTCAddress["ETH"],transfer.coi)
			//
			//}

			//插入站内信
			//err = msg.TransferSuccess(transfer.CoinName, transfer.Amount, transfer.OrderId, transfer.ApplyerId)
			//if err != nil {
			//	log.Error("插入站内信失败", err)
			//}
			log.Debug("上链成功")

			return

		} else {
			ts.UpdateTransfer(&db.Transfer{ID: transfer.TransferId, Status: common.TransferFailed, TxId: txu.TxId()})
			//插入站内信
			//msg.TransferFail(transfer.CoinName, transfer.Amount, transfer.TransferId, transfer.ApplyerId)
		}
		return
	}
	ts.UpdateTransfer(&db.Transfer{ID: transfer.TransferId, Status: common.TransferSuccess, TxId: txu.TxId()})

	//插入站内信
	//err = msg.TransferSuccess(transfer.CoinName, transfer.Amount, transfer.OrderId, transfer.ApplyerId)
	//if err != nil {
	//	log.Error("插入站内信失败", err)
	//}
	log.Debug("上链成功")
}

func createTxBTC(bc_str bccore.BlockChainSign, fromAddress []string, token string, amounts []*bccoin.AddressAmount, miner float64, transfer *db.TransferSend) (string, txutil.TxUtil) {
	log.Info("构造交易的各个参数", bc_str, fromAddress, transfer.TokenAddress, miner)

	ts := &db.TransferDBService{}
	uuid, txu, err := Wallet.CreateTx(bc_str, fromAddress, transfer.TokenAddress, amounts, miner)
	if err != nil {
		if err == walletError.ERR_ADDRESS_QUEUE_BLOCKED {
			log.Error("阻塞", err)
			ts.UpdateOrderExpire(&db.Transfer{OrderId: transfer.OrderId, Status: common.TransferApprovaled})
			backChannel(transfer)
		} else if err == walletError.ERR_NOT_ENOUGH_COIN {
			log.Error("余额不足", err)
			if transfer.Types == 1 {
				ts.UpdateOrderExpire(&db.Transfer{OrderId: transfer.OrderId, Status: common.TransferInsufficientBalance})
				transfer.Status = common.TransferInsufficientBalance
				backChannel(transfer)
			} else {
				ts.UpdateOrderExpire(&db.Transfer{OrderId: transfer.OrderId, Status: common.TransferFailed})
				// 插入站内信
				orderInfo, _ := ts.FindOrderById(transfer.OrderId)
				err = msg.TransferFail(orderInfo.CoinName, orderInfo.Amount, orderInfo.ID, orderInfo.ApplyerId)
				if err != nil {
					log.Error("插入站内信失败", err)
				}
			}
		} else {
			ts.UpdateOrderExpire(&db.Transfer{OrderId: transfer.OrderId, Status: common.TransferGenTxError})
			log.Error("构造交易失败", err)
		}
		return "", nil
	} else {
		return uuid, txu
	}
}

func sendTxBTC(bc_str bccore.BlockChainSign, txu txutil.TxUtil, uuid string, transfer *db.TransferSend) {
	log.Debug("正在上链")
	err := Wallet.SendTx(bc_str, txu, uuid)
	ts := &db.TransferDBService{}
	if err != nil {
		log.Error("上链错误", err)

		//余额不足
		if err == walletError.ERR_NOT_ENOUGH_COIN {
			if transfer.Types == 1 {
				ts.UpdateOrderExpire(&db.Transfer{OrderId: transfer.OrderId, Status: common.TransferInsufficientBalance, TxId: txu.TxId()})
				transfer.Status = common.TransferInsufficientBalance
				backChannel(transfer)
			} else {
				ts.UpdateOrderExpire(&db.Transfer{OrderId: transfer.OrderId, Status: common.TransferFailed, TxId: txu.TxId()})
				//插入站内信
				msg.TransferFail(transfer.CoinName, transfer.Amount, transfer.OrderId, transfer.ApplyerId)
			}

		} else if err == walletError.ERR_PIPELINE_DATA_ILLEGAL {
			ts.UpdateOrderExpire(&db.Transfer{OrderId: transfer.OrderId, Status: common.TransferSuccess, TxId: txu.TxId()})

			//插入站内信
			err = msg.TransferSuccess(transfer.CoinName, transfer.Amount, transfer.OrderId, transfer.ApplyerId)
			if err != nil {
				log.Error("插入站内信失败", err)
			}
			log.Debug("上链成功")

			return
		} else {
			ts.UpdateOrderExpire(&db.Transfer{OrderId: transfer.OrderId, Status: common.TransferFailed, TxId: txu.TxId()})
			//插入站内信
			msg.TransferFail(transfer.CoinName, transfer.Amount, transfer.OrderId, transfer.ApplyerId)
		}
		return
	}
	ts.UpdateOrderExpire(&db.Transfer{OrderId: transfer.OrderId, Status: common.TransferSuccess, TxId: txu.TxId()})

	//插入站内信
	err = msg.TransferSuccess(transfer.CoinName, transfer.Amount, transfer.OrderId, transfer.ApplyerId)
	if err != nil {
		log.Error("插入站内信失败", err)
	}
	log.Debug("上链成功")
}

//签名机签名
func voucherSign(transfer *db.TransferSend, txu txutil.TxUtil, addMap map[string][]uint32, coinType string) (txutil.TxUtil, bool) {
	transfers := TransferInfo{
		TransferMsgs:   transfer.TransferMsg,
		ApplyAccount:   transfer.ApplyAccount,
		ApplyPublickey: transfer.ApplyPublickey,
		ApplySign:      transfer.ApplySign,
		ApproversSign:  transfer.ApproversSign,
		ToAddress:      transfer.ToAddress,
		AmountIndex:    transfer.AmountIndex,
		TemInfo:        transfer.TemInfo,
	}
	log.Info("addMap", addMap)

	jsonAddMap, _ := json.Marshal(addMap)
	jsonTransfers, _ := json.Marshal(transfers)
	txinfo, _ := txu.Marshal()
	voucherVo := &voucher.GrpcServer{
		Currency: coinType,
		AddrMap:  jsonAddMap,
		TxInfo:   txinfo,
		TransMsg: jsonTransfers,
		Type:     voucher.VOUCHER_OPERATE_SIGN,
	}

	if transfer.Types == 2 {
		voucherVo.Type = voucher.VOUCHER_OPERATE_TRANSSIGN_INTERIOR
	}
	_, res, ret := voucher.SendVoucherData(voucherVo)
	voucherTimer := time.NewTicker(3 * time.Second)
	for ret == voucher.VRET_TIMEOUT || (res.Status == voucher.VOUCHER_STATUS_PAUSED) {
		select {
		case <-voucherTimer.C:
			_, res, ret = voucher.SendVoucherData(voucherVo)
		}
	}
	if ret != voucher.VRET_CLIENT {
		log.Error("签名机错误")
		return nil, false
	}
	if res.Status != voucher.STATUS_OK {
		err := msg.VoucherFail(msg.VoucherRetTransFailMsgType, res.Status, nil, transfer)
		if err != nil {
			log.Error("发送站内信失败", err)
		}
		voucherResMsg, _ := json.Marshal(res)
		log.Error("签名机返回错误", string(voucherResMsg))
		return nil, false
	}
	err := txu.UnMarshal(res.SignTx)
	if err != nil {
		log.Error("签名后解析交易错误", err)
		return nil, false
	}
	return txu, true
}

//同步订单状态
func asycOrderStatus(orderId string) {
	ts := &db.TransferDBService{}

	list := ts.GetTransferByOrderId(orderId)
	if len(list) == 0 {
		return
	}
	success := 0
	transfering := 0
	for i := 0; i < len(list); i++ {
		// 1和3， 转账中  12 余额不足等待转账
		if list[i].Status == common.TransferSomeSuccess || (list[i].Status == common.TransferApprovaled) || (list[i].Status == common.TransferInsufficientBalance) {
			transfering++
		}
		if list[i].Status == common.TransferSuccess {
			success++
		}
	}
	if transfering > 0 {
		return
	}
	if success == len(list) {
		ts.UpdateOrderStatus(orderId, common.TransferSuccess)
	}
	if (success > 0) && (success < len(list)) && (transfering == 0) {
		ts.UpdateOrderStatus(orderId, common.TransferSomeSuccess)
	}
	if (success == 0) && (transfering == 0) {
		ts.UpdateOrderStatus(orderId, common.TransferFailed)
	}

}

//更新合并账户状态
func asycCombineOrder(orderId string) {
	ts := &db.TransferDBService{}
	tx := &db.TxinfosService{}

	list := ts.GetTransferByOrderId(orderId)
	if len(list) == 0 {
		return
	}

	transfering := 0
	for i := 0; i < len(list); i++ {
		// 1和3， 转账中
		if list[i].Status == common.TransferSomeSuccess || (list[i].Status == common.TransferApprovaled) {
			transfering = 1
		}
	}
	if transfering > 0 {
		return
	}
	tx.UpadateConfigs("combine_account", "0")
	db.ConfigMap["combine_account"] = "0"

}

func updateConnectStatus() {
	if voucher.IsConnect() {
		voucher.HeartStatus = append(voucher.HeartStatus, 1)
	} else {
		voucher.HeartStatus = append(voucher.HeartStatus, 0)
	}
	if len(voucher.HeartStatus) > 5 {
		voucher.HeartStatus = voucher.HeartStatus[len(voucher.HeartStatus)-5 : len(voucher.HeartStatus)]
	}

}

func updateAsyncStatus(){

	result := Wallet.GetHeights()
	coinType := map[bccore.BloclChainType]string{1: "BTC", 2: "ETH", 3: "LTC"}
	cfg := config.Conf
	heightDiff := map[bccore.BloclChainType]uint64{1: cfg.Server.BTCHeightDiff, 2: cfg.Server.ETHHeightDiff, 3: cfg.Server.LTCHeightDiff}
	// returned data
	var list []interface{}
	aggregatedStatus := 0 //正常
	for i, v := range result {
		status := 0
		if (v.PubHeight - v.CurHeight) > heightDiff[i] {
			status = 1
			aggregatedStatus = 1
		}
		list = append(list, map[string]interface{}{
			"PubHeight" : v.PubHeight,
			"CurHeight" : v.CurHeight,
			"Name"     : coinType[i],
			"Status"    : status,
		})
	}
	if db.AsyncBlockChains.AggregatedStatus != aggregatedStatus {
		if aggregatedStatus == 0 {
			msg.AysncBlockChainStatus(1)
		} else {
			msg.AysncBlockChainStatus(2)
		}
	}
	db.AsyncBlockChains = &db.AsyncBlockChain{
		List:list,
		AggregatedStatus : aggregatedStatus,
	}

}

func removeArr(s []string, value string) []string {
	index := -1
	for i := 0; i < len(s); i++ {
		if s[i] == value {
			index = i
			break
		}
	}
	if index == -1 {
		return s
	}
	return append(s[:index], s[index+1:]...)
}

func transferMap(maps string) ([]uint32, error) {
	var addMap []uint32
	if maps == "" {
		return nil, nil
	}
	err := json.Unmarshal([]byte(maps), &addMap)

	return addMap, err
}
