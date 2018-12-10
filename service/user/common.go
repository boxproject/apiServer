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
package user

import "time"

// 认证股东
const (
	VerifyOwnerApprove int = iota + 1
	VerifyOwnerReject
)

type UserVo struct {
	ID           int    `form:"account_id"`
	Name         string `form:"name"`
	AppId        string `form:"app_id"`
	PubKey       string `form:"pub_key"`
	Pwd          string `form:"pwd"`
	Msg          string `form:"msg"`
	NewAppId     string `form:"new_app_id"`
	Status       int    `form:"status"`				// 2.approval 3.reject
	RefuseReason string `form:"refuse_reason"`
	UserType     int    `form:"user_type"`
	AuthID       string `form:"auth_id"`
	Add          string `form:"add"`
	DepID        int    `form:"dep_id"`
	RegID        int    `form:"reg_id"`
	Level        int
}

type AccountTree struct {
	Name        string
	AppId       string
	PubKey      string
	Msg         string
	SourceAppId string
	Children    []AccountTree
}

type ModifyPSD struct {
	ID          int    `form:"id"`
	Page        int    `form:"page"`
	OldPassword string `form:"old_psd"`
	NewPassword string `form:"new_psd"`
}

type PAccount struct {
	ID       *int    `form:"account_id" binding:"exists"`
	AuthID   *string `form:"auth_id" binding:"exists"`
	DepID    *int    `form:"dep_id" binding:"exists"`
	UserType *int    `form:"user_type" binding:"exists"`
}

type AccountTreeWithCoreInfo struct {
	AccTree []AccountTree     `json:"trees"`
	Voucher VoucherPubkeyInfo `json:"voucher"`
}

type VoucherPubkeyInfo struct {
	ENPublickey string
	ENKey       string
}

type Letter struct {
	ID      int
	Title   string
	Status  int
	Time    time.Time
	Content string
	Type    int
	Param   string
	Warn 	bool
}

type VoucherVo struct {
	AppId     string `form:"app_id"`
	AppName   string `form:"name"`
	AesKey    string `form:"aeskey"`
	PubKey    string `form:"pub_key"`
	Pwd       string `form:"pwd"`
	Msg       string `form:"msg"`
	Sign      string `form:"sign"`
	Timestamp int64  `form:"timestamp"`
	Password  string `form:"string"`
	BakAction string `form:"bak_action"`
	Status    int    `form:"status"`
	RegId     int    `form:"reg_id"`
	Replaced  string `form:"replaced"`
}
