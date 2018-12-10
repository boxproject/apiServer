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
package template

type AddTmp struct {
	Content      string `form:"content"`
	TemplateSign string `form:"template_sign"`
	Types        int    `form:"types"` // 获取列表的类型 0.全部  1.等待审批

}

type TemplateVo struct {
	Id          string `form:"template_id"`
	AdminStatus int    `form:"status"`
	AccountId   int

	AppId        string `form:"app_id"`
	AppName      string `form:"app_name"`
	KeyLine      string `form:"key_kine"`
	Timestamp    int64  `form:"timestamp"`
	Sign         string `form:"sign"`
	TemplateSign string `form:"template_sign"`
	Aeskey       string `form:"aeskey"`
	Msg          string `form:"msg"`
	Reason       string `form:"reason"`
}

// 作废审批流模板
type ParamCancelTemplate struct {
	TemplateId   string `form:"template_id"`
	TemplateSign string `form:"template_sign"`
	Timestamp    int64  `form:"timestamp"`
	Sign         string `form:"sign"`
	Pwd          string `form:"pwd"`
}
