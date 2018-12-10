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
package logger

var (
	OperLoggerTempMap = make(map[string]map[string]interface{})
)

const (
	ZhLoggerFilePath = "./static/lang/logger/logger_zh.toml"
	EnLoggerFilePath = "./static/lang/logger/logger_en.toml"
)

type OperateLogger struct {
	Operator string
}
