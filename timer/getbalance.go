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
	"time"
	"github.com/boxproject/apiServer/db"
	"github.com/boxproject/boxwallet/bccore"
	"github.com/boxproject/boxwallet/bccoin"
	log "github.com/alecthomas/log4go"
	"github.com/boxproject/apiServer/config"
)

func updateDepositBalance(){
	confirm :=0
	//cfg,_ := config.LoadConfig()
	cfg := config.Conf
	ts := &db.TxinfosService{}
	as := &db.AccManageService{}
	//拿取充值列表的数据
	depositList := ts.FindDeposit()
	if len(depositList) == 0{
		return
	} else{
		deposit := depositList[0]
		var amount  bccoin.CoinAmounter
		if deposit.CoinName == "ETH"{
			confirm = cfg.Confirm.ETH
		} else if deposit.CoinName == "BTC"{
			confirm = cfg.Confirm.BTC
		} else if deposit.CoinName == "USDT"{
			confirm = cfg.Confirm.BTC
		} else if deposit.CoinName == "LTC"{
			confirm = cfg.Confirm.LTC
		}else if deposit.TokenType!=0{
			confirm = cfg.Confirm.ETH
		}
		var err error
		if deposit.Nowconfirm < confirm{
			status := &db.Deposit{
				ID:deposit.DepositId,
				IsUpdate:3,
				UpdatedAt:time.Now(),
				Confirm:deposit.Nowconfirm,
			}
			ts.UpdateDeposit(status)
			return
		} else {
			//if deposit.ExtValid == 0 {
			//	status := &db.Deposit{
			//		ID:deposit.DepositId,
			//		IsUpdate:4,
			//		UpdatedAt:time.Now(),
			//		Confirm:deposit.Nowconfirm,
			//	}
			//	ts.UpdateDeposit(status)
			//	return
			//}
			if deposit.TokenType !=0 {
				amount,err = Wallet.GetBalance(bccore.STR_ERC20,deposit.ToAddr,deposit.TokenAddress)
			} else {
				if deposit.CoinName == "USDT"{
					amount,err = Wallet.GetBalance(bccore.STR_USDT,deposit.ToAddr,USDT_PPID)
				} else{
					amount,err = Wallet.GetBalance(bccore.BlockChainSign(deposit.CoinName),deposit.ToAddr,"")
				}

			}
			if err != nil{
				log.Error("wallet getbalance cuowu",err)
				return
			}
			capital := &db.Capital{
				AddressID:deposit.AddressId,
				CoinName:deposit.CoinName,
				Balance:amount.String(),
				CoinId:deposit.CoinId,
				Address:deposit.ToAddr,

			}
			ret := as.UpdateBalance(capital)
			if ret==true {
				status := &db.Deposit{
					ID:deposit.DepositId,
					Confirm:deposit.Nowconfirm,
					IsUpdate:2,
				}
				ts.UpdateDeposit(status)
			}

			if deposit.CoinName == "BTC"{
				AllBTCAddress["BTC"] = addArrValue(AllBTCAddress["BTC"],deposit.ToAddr)
				AllBTCAddressMap[deposit.ToAddr] = deposit.Deep
			} else if deposit.CoinName == "LTC" {
				AllBTCAddress["LTC"] = addArrValue(AllBTCAddress["LTC"],deposit.ToAddr)
				AllLTCAddressMap[deposit.ToAddr] = deposit.Deep
			}
		}

	}
}




func updateBalance(){
	ts := &db.TxinfosService{}
	balances := ts.GetBalanceList()
	if len(balances)==0 {
		return
	} else {
		balance := balances[0]
		var amount  bccoin.CoinAmounter
		var err error
		if balance.TokenType !=0 {
			amount,err = Wallet.GetBalance(bccore.STR_ERC20,balance.Address,balance.TokenAddress)
		} else {
			if balance.CoinName == "USDT"{
				amount,err = Wallet.GetBalance(bccore.STR_USDT,balance.Address,USDT_PPID)
			} else{
				amount,err = Wallet.GetBalance(bccore.BlockChainSign(balance.CoinName),balance.Address,"")
			}
		}
		if err != nil{
			log.Error("wallet getbalance cuowu",err)
			return
		}
		if amount == nil {
			log.Error("wallet getbalance amount",amount)
			return
		}
		ts.UpdateBalance(&db.Capital{ID:balance.CapitalId,Balance:amount.String(),CreatedAt:time.Now()})
	}
}


func updateBalances(coinName string,address string,coinId int,tokenAddress string){
	ts := &db.TxinfosService{}
	var amount  bccoin.CoinAmounter
	var err error
	if tokenAddress != "" {
		amount,err = Wallet.GetBalance(bccore.STR_ERC20,address,tokenAddress)
	} else {
		if coinName == "USDT" {
			amount,err = Wallet.GetBalance(bccore.STR_USDT,address,USDT_PPID)
		} else {
			amount,err = Wallet.GetBalance(bccore.BlockChainSign(coinName),address,"")
		}
	}
	if err !=nil {
		log.Error("更新余额出问题",coinName,address,tokenAddress)
		return
	}
	ts.UpdateBalances(coinName,address,amount.String())
}


func sumBalance() {
	dt := &db.TxinfosService{}
	err := dt.SumBalance()
	if err != nil {
		log.Error("余额更新失败", err)
		return
	}
}



func addArrValue(arr []string,value string) []string{
	index := 0
	for i:=0;i<len(arr);i++ {
		if arr[i] == value{
			index = 1
		}
	}
	if index == 0 {
		return append(arr,value)
	} else {
		return arr
	}

}


