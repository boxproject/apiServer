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
package db

import "time"

// blockchain explorer url
const (
	ETH_LINK       = "https://etherscan.io/tx/"
	BTC_LINK       = "https://www.blockchain.com/btc/tx/"
	LTC_LINK       = "https://chainz.cryptoid.info/ltc/tx.dws?"
	AuthTypeZhPath = "./static/lang/auth/auth_zh.toml"
	AuthTypeEnPath = "./static/lang/auth/auth_en.toml"
)

var AuthMaps = make(map[string]map[string]interface{})

// TemplateContent represents the template content struct
type TemplateContent struct {
	Name         string          `json:"name"`
	Period       int             `json:"period"`
	ApprovalInfo []ApprovalInfo  `json:"approvalInfo"`
	LimitInfo    []TemplateLimit `json:"limitInfo"`
}

// ApprovalInfo represents the approval info defined for a template
type ApprovalInfo struct {
	Require   int         `json:"require"`
	Approvers []Approvers `json:"approvers"`
}

// Approvers represents approvers info of a template
type Approvers struct {
	Account string `json:"account"`
	PubKey  string `json:"pubkey"`
}

// TemplateLimit represents the transfer limit amount information for every coin of the template
type TemplateLimit struct {
	TokenAddress string `json:"tokenAddress"`
	Symbol       string `json:"symbol"`
	FullName     string `json:"name"`
	Decimal      int    `json:"precise"`
	Limit        string `json:"limit"`
}

// TxApply represents the struct of transfer apply content
type TxApply struct {
	Amount       string `json:"amount"`
	CoinName     string `json:"coinName"`
	Info         string `json:"info"`
	Miner        string `json:"miner"`
	Destination  string `json:"destination"`
	Timestamp    int64  `json:"timestamp"`
	TemplateHash string `json:"templateHash"`
}

// AddressJson represents the address of its deep info
type AddressJson struct {
	Address string
	Deep    []uint32
}

// DelayTasks represents the struct of to-do task
type DelayTasks struct {
	Transfer DelayTask
	Account  DelayTask
	Template DelayTask
}

type DelayTask struct {
	Number int
	Reason string
}


type bachInfo struct {
	ApplyName   string    `gorm:"column:applyName"`
	CoinName    string    `gorm:"column:coinName"`
	ID          string    `json:"OrderId" gorm:"column:id"`
	Amount      string    `gorm:"column:amount"`
	Status      int       `gorm:"column:status"`
	ApplyReason string    `gorm:"column:applyReason"`
	CreatedAt   time.Time `gorm:"column:createdAt"`
}


// MessageTitle station letter
type MessageTitle string
