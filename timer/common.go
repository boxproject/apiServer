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

type ApplyMsg struct {
	CoinName  string     `json:"coin_name"`
	TemHash   string     `json:"t_hash"`
	Reason    string     `json:"reason"`
	Deadline  string     `json:"deadline"`
	Miner     string     `json:"miner"`
	Timestamp int        `json:"timestamp"`
	Amount    string     `json:"amount"`
	ApplyVos  []ApplyVos `json:"applys"`
	CoinId    int        `json:"coin_id"`
	Token     string     `json:"token"`
	Currency  string     `json:"currency"`
}

type ApplyVos struct {
	ToAddress string   `json:"to_address"`
	Tag       string   `json:"tag"`
	Amount    []string `json:"amount"`
}

type TransferInfo struct {
	TransferMsgs   string
	ApplyAccount   string
	ApplyPublickey string
	ApplySign      string
	ApproversSign  string
	AmountIndex    int
	ToAddress      string
	TemInfo        string
}

type TransferMsgJson struct {
	FromAddress AddressJson
	ToAddress   AddressJson
	Token       string
}
type AddressJson struct {
	Address string
	Deep    []uint32
}
