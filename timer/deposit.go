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
	"github.com/boxproject/apiServer/db"
	"encoding/json"
	log "github.com/alecthomas/log4go"
	"strconv"
	"github.com/satori/go.uuid"
)


type TxinfoEntity struct {
	In 	 	[]addr
	Out  	[]addr
	ExtIn	[]addr
	ExtOut	[]addr
}

type addr struct{
	Amt		string
	Addr	string
}

func depositBusiness() {
	ts := &db.TxinfosService{}
	txinfos := ts.FindTxinfos()
	if len(txinfos) == 0 {
		return
	}
	for i:=0;i<len(txinfos);i++{
		depositTx(txinfos[i])
		db.ConfigMap["getbalance_time"] = txinfos[i].CreateAt.Format("2006-01-02 15:04:05")
	}
	ts.UpadateConfigs("getbalance_time",db.ConfigMap["getbalance_time"])
}

func depositTx (tx db.TxInfos){
	ts := &db.TxinfosService{}
	as := &db.AccManageService{}
	var txinfo TxinfoEntity
	coinName := ""
	fromAddress:=""
	toAddress:=""
	amount :=""
	err:=json.Unmarshal([]byte(tx.TxObj),&txinfo)
	if err!=nil {
		log.Error("txinfo json失败",err)
		return
	}
	switch tx.Type {
		case 1:{

			if tx.Token !="" {
				coinName = "USDT"
				fromAddress = txinfo.ExtIn[0].Addr
				toAddress = txinfo.ExtOut[0].Addr
				amount = txinfo.ExtOut[0].Amt
				b := as.FindAddressByAddress(toAddress,coinName)
				//如果地址不存在，返回
				if !b {
					return
				}
			}else {
				coinName = "BTC"
				getBTCAddress(&txinfo,coinName,tx)
				return
			}
		}
		case 2:{
			if tx.Token !=""{
				coinName = "ETH"
				fromAddress = txinfo.ExtIn[0].Addr
				toAddress = txinfo.ExtOut[0].Addr
				amount = txinfo.ExtOut[0].Amt
			}else{
				coinName = "ETH"
				if txinfo.Out == nil {
					return 
				}
				fromAddress = txinfo.In[0].Addr
				toAddress = txinfo.Out[0].Addr
				amount = txinfo.Out[0].Amt
			}
			b := as.FindAddressByAddress(toAddress,coinName)
			//如果地址不存在，返回
			if !b {
				return
			}
		}
		case 3:{
			coinName = "LTC"
			getBTCAddress(&txinfo,coinName,tx)
			return

		}
	}
	coin ,ret:=ts.FindCoin(coinName,tx.Token)
	if ret == 0{
		return
	}

	deposit := &db.Deposit{
		FromAddr:fromAddress,
		ToAddr:toAddress,
		Amount:amount,
		TxId:tx.Txid,
		BlockNum:strconv.Itoa(tx.Height),
		Confirm:tx.Confirm,
		ID:uuid.Must(uuid.NewV4()).String(),
		TokenAddress:tx.Token,

	}
	if ret == 2 {
		deposit.CoinId = coin.ID
		deposit.CoinName = coin.Name
		deposit.IsUpdate = 3
		deposit.Precise = coin.Precise
	}

	ts.SaveDeposit(deposit)
}

func getBTCAddress (tx *TxinfoEntity,coinName string,txinfos db.TxInfos) (){
	ts := &db.TxinfosService{}
	coin ,ret:=ts.FindCoin(coinName,"")
	if ret == 0{
		return
	}
	addr := make(map[string]string)
	as := &db.AccManageService{}
	fromAddress := tx.In[0].Addr
	toAddresses := []string{}
	amount := ""
	mapAmount := make(map[string][]string)
	for i:=0;i<len(tx.In);i++{
		addr[tx.In[i].Addr] = "1"
	}

	for j:=0;j<len(tx.Out);j++{
		if addr[tx.Out[j].Addr] !="1" {
			toAddresses = append(toAddresses,tx.Out[j].Addr)
			amount = tx.Out[j].Amt
			mapAmount[tx.Out[j].Addr] = append(mapAmount[tx.Out[j].Addr],amount)
		}
	}
	log.Error("mapAmount----",mapAmount)
	toAddress,_ := as.FindAddressByBTCAddress(toAddresses,coin.ID)
	for i:=0;i<len(toAddress);i++{
		for j:=0;j<len(mapAmount[toAddress[i].Address]);j++{
			deposit := &db.Deposit{
				FromAddr:fromAddress,
				ToAddr:toAddress[i].Address,
				Amount:mapAmount[toAddress[i].Address][j],
				TxId:txinfos.Txid,
				BlockNum:strconv.Itoa(txinfos.Height),
				Confirm:txinfos.Confirm,
				ID:uuid.Must(uuid.NewV4()).String(),
				CoinId : coin.ID,
				CoinName : coin.Name,
				IsUpdate : 3,
				Precise : coin.Precise,
			}
			ts.SaveDeposit(deposit)
		}

	}
}



