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
package utils

type RPCRsp struct {
	RspNo string
}

// 签名机状态
type VoucherStatus struct {
	RPCRsp
	Status       ServerStatus
	ApprovalInfo FlowStatus
	CoinStatus   []CoinStatu
	TokenInfos   []TokenInfo
}

// 查询审批流上链状态返回值结构
type ResHashStatus struct {
	RspNo        string
	Status       ServerStatus
	ApprovalInfo FlowStatus
}

// 审批流操作日志
type ResFlowOpLog struct {
	RPCRsp
	HashOperates []FlowOpLogInfo
}

// 审批流操作日志详情
type FlowOpLogInfo struct {
	ApplyerAccount string
	CaptainId      string
	Option         string
	Opinion        string
	CreateTime     string
}

type ServerStatus struct {
	ServerStatus    int64
	Status          int64
	Total           int64
	HashCount       int64
	TokenCount      int64
	Address         string
	ContractAddress string
	BtcAddress      string
	CoinStatus      []CoinStatu
}

type FlowStatus struct {
	Hash      string
	Name      string
	AppId     string
	CaptainId string
	Flow      string
	Sign      string
	Status    string
}

type CoinStatu struct {
	Name     string
	Category int64
	Decimals int64
	Used     bool
}

type TokenInfo struct {
	TokenName    string
	Decimals     int64
	ContractAddr string
	Category     int64
}
