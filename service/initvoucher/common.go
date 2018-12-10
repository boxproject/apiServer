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
package initvoucher

type VoucherVo struct {
	AppId     string `form:"app_id"`
	AppName   string `form:"name"`
	AesKey    string `form:"aeskey"`
	Msg       string `form:"msg"`
	Sign      string `form:"sign"`
	Timestamp int64  `form:"timestamp"`
	Password  string `form:"password"`
	BakAction string `form:"bak_action"`
	IsRecover int    `form:"isRecover"`
}
