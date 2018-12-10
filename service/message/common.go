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
package message

const (
	VoucherRetTempFailMsgType int = iota
	VoucherRetTransFailMsgType
)

var (
	MsgContentTempMap = make(map[string]map[string]interface{})
	MsgAccountTypeMap = make(map[string]map[string]interface{})
	MsgTitleTempMap   = make(map[string]map[string]interface{})
)

const (
	ZhMsgFilePath     = "./static/lang/message/message_zh.toml"
	EnMsgFilePath     = "./static/lang/message/message_en.toml"
	MsgAccountTypeTag = "accountType"
	MsgTitleTag       = "title"
	MsgContentTag     = "content"
)

type AccountTypeName struct {
	AdminType string
	OwnerType string
}

type MessageTitleAndContent struct {
	OperAccountType  string
	OperAccountName  string
	AdminAccount     string
	AuthName         string
	TemplateName     string
	DisabledAccount  string
	TransOrderName   string
	VoucherWarnCount int
}
