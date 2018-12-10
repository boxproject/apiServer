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
package config

type Config struct {
	Server   serverInfo
	Database DBInfo
	Voucher  VoucherInfo
	LogPath  string
	LangPath string
	Confirm  ConfirmInfo
}

type serverInfo struct {
	Port          string
	APIVersion    string
	Mode          string
	AppPath       string
	BTCHeightDiff uint64
	ETHHeightDiff uint64
	LTCHeightDiff uint64
}

type DBInfo struct {
	User      string
	Password  string
	DbName    string
	Host      string
	MaxOpen   int
	MaxIdle   int
	DebugMode bool
}

type VoucherInfo struct {
	Port      string
	ServerPem string
	ServerKey string
	ClientPem string
}

type ConfirmInfo struct {
	ETH int
	BTC int
	LTC int
}
