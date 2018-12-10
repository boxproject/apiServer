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
package transfer

import (
	"time"
)

// 返回网页转账状态
const (
	WebTransWaitingForQr int = iota + 1
	WebTransSubmitted
	WebTransExpired
)

type TransferVo struct {
	AppId            string `form:"app_id"`
	Account          string `form:"account"`
	CoinId           int    `form:"coin_id"`
	TemplateId       string `form:"template_id"`
	TemplateHash     string
	Password         string         `form:"password"`
	OrderInfo        string         `form:"order_info"`
	Deadline         time.Time      `form:"deadline"`
	Fee              string         `form:"fee"`
	TransferInfoList []TransferInfo `form:"transferinfo_list"`
	TransferId       string         `form:"transfer_id"`
	ApplyMsg         string         `form:"apply_msg"`
	ApplySign        string         `form:"apply_sign"`
	ApplyerId        int
	AppName          string
	OrderIds         string `form:"order_ids"`
	ArrOrderIds      []string
	OrderId          string `form:"order_id"`
}

type TransferOrder struct {
	Id   string `form:"id"`
	Sort int    `form:"sort"`
	Hash string `form:"hash"`
}

type TransferInfo struct {
	Amount  string `form:"amount"`
	Address string `form:"address"`
	mark    string `form:"mark"`
}

type ApplyContent struct {
	Amount      string `form:"amount"`
	Coin        string `form:"coin_name"`
	Info        string `form:"info"`
	Miner       string `form:"miner"`
	Destination string `form:"destination"`
	Timestamp   int64  `form:"timestamp"`
	FlowHash    string `form:"flow_hash"`
	Sign        string `form:"sign"`
}

type ApplyMsg struct {
	CoinName  string    `json:"coin_name"`
	TemHash   string    `json:"t_hash"`
	Reason    string    `json:"reason"`
	Deadline  string    `json:"deadline"`
	Miner     string    `json:"miner"`
	Timestamp int       `json:"timestamp"`
	Amount    string    `json:"amount"`
	ApplyVos  []ApplyVo `json:"applys"`
	CoinId    int       `json:"coin_id"`
	Token     string    `json:"token"`
	Currency  string    `json:"currency"`
}

type ApplyVo struct {
	ToAddress string   `json:"to_address"`
	Tag       string   `json:"tag"`
	Amount    []string `json:"amount"`
}

type VerifyApplyVo struct {
	Status       int    `form:"status"`
	OrderId      string `form:"order_id"`
	AppName      string
	TransferSign string `form:"transfer_sign"`
	Reason       string `form:"reason"`
	AccountId    int
	OrderIds     string `form:"order_ids"`
	ArrOrderIds  []string
	TemHash      string `form:"tem_hash"`
	TemInfo      string
}

//对当前层数的汇总
type CountVo struct {
	AgreeNum   int
	RefuseNum  int
	TotalNum   int
	RequireNum int
	NowLevel   int
	TotalLevel int
}

//通过币种查询模板
type TemplateVo struct {
	CoinId int `form:"coin_id"`
}

//转账数据库里msg 的json数据
type TransferMsgJson struct {
	FromAddress AddressJson
	ToAddress   AddressJson
	Token       string
}
type AddressJson struct {
	Address string
	deep    string
}

//批量审批专用
type BatchVerify struct {
	OrderIds string `form:"order_ids"`
	Reason   string `form:"reason"`
}

//批量审批专用
type BatchOrder struct {
	OrderId      string `json:"order_id"`
	Status       int    `json:"status"`
	TransferSign string `json:"transfer_sign"`
	AppName      string
}

type WebTransfer struct {
	Id          string `form:"id"`
	Msg         string `form:"msg"`
	CreatedTime int64  `form:"createdtime"`
}
