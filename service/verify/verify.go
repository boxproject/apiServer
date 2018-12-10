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
package verify

import (
	"fmt"
	"strconv"
	"time"

	log "github.com/alecthomas/log4go"
	"github.com/boxproject/apiServer/db"
	"github.com/boxproject/apiServer/errors"
	"github.com/boxproject/apiServer/utils"
)

// VerifyPSW validate password record wrong time and freeze account
func VerifyPSW(name, password string) (*VerifyResult, *db.Account, errors.CodeType) {
	var result VerifyResult
	if name == "" || password == "" {
		result.Result = false
		result.Reason = "缺少参数"
		return &result, nil, errors.ParamsNil
	}
	as := &db.AccDBService{}
	var salt string
	// validate account
	valideAccInfo, err := as.ExistAccount(name)
	if err != nil {
		result.Result = false
		result.Reason = "系统出错了"
		return &result, nil, errors.Failed
	}
	if valideAccInfo == nil {
		result.Result = false
		result.Reason = "账户不存在"
		return &result, nil, errors.User_3001
	}
	if valideAccInfo.IsDeleted == 1 {
		result.Result = false
		result.Reason = "账户已被停用"
		return &result, nil, errors.User_3005
	} else {
		salt = valideAccInfo.Salt
	}
	frozen, frozenTo, err := as.AccoutIsFrozen(name)
	if err != nil {
		result.Result = false
		result.Reason = "系统出错了"
		return &result, nil, errors.Failed
	}
	timesp := frozenTo.Unix() //frozen time
	newTime := time.Now().Unix()
	if err != nil {
		result.Result = false
		result.Reason = "系统出错了"
		return &result, nil, errors.Failed
	}
	if frozen == 1 && timesp-newTime > 0 { // if still frozen
		result.Result = false
		result.Data = frozenTo.Local().Format("2006-01-02 03:04:05 PM")
		result.Reason = "输入错误次数过多,账户被锁定到:" + result.Data
		return &result, nil, errors.User_3007
	}
	pwd := utils.PwdSalt(password, salt)
	validePassword, account, err := as.Login(name, pwd)
	if err != nil {
		result.Result = false
		result.Reason = "系统出错了"
		return &result, nil, errors.Failed
	}
	if validePassword == false { // wrong password and record times
		num, FrozenToTime := as.RecordAttempts(name)
		if num == 0 {
			result.Result = false
			result.Reason = "系统出错了"
			return &result, nil, errors.Failed
		} else if num == 5 { // 5 time incorrect password
			result.Result = false
			result.Data = FrozenToTime.Local().Format("2006-01-02 03:04:05 PM")
			result.Reason = "输入错误次数过多,账户被锁定到:" + result.Data
			return &result, nil, errors.User_3007
		} else if 5 > num && num > 0 {
			result.Result = false
			result.Data = strconv.Itoa(num)
			result.Reason = "连续输错5次将锁定账户8小时,您已输错" + result.Data + "次"
			return &result, nil, errors.User_3002
		}
	}
	// reset attemps
	err = as.ResetAttempts(name)
	if err != nil {
		log.Error("登录成功重置尝试次数", err)
		result.Result = false
		result.Reason = "系统出错了"
		return &result, nil, errors.Failed
	}
	result.Result = true
	return &result, account, errors.Success
}

// IsFrozen check if account is frozen
func IsFrozen(appName string) errors.ErrModel {
	ds := &db.AccDBService{}
	frozen, frozenTo, accErr := ds.AccoutIsFrozen(appName)
	if accErr != nil {
		return errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	}
	if time.Now().Before(frozenTo) && frozen == 1 {
		message := "输入错误次数过多,账户被锁定到:" + frozenTo.Local().Format("2006-01-02 03:04:05 AM")
		return errors.ErrModel{Code: errors.User_3007, Err: errors.New(message), Data: map[string]interface{}{
			"data": frozenTo.Local().Format("2006-01-02 03:04:05 AM"),
		}}
	}
	return errors.ErrModel{Code: errors.Success, Err: nil}
}

// Keyword record keyword failed times
func Keyword(appID, appName string) errors.ErrModel {
	var code errors.CodeType
	ds := db.AccDBService{}
	num, frozenTime := ds.RecordAttempts(appName)
	message := fmt.Sprintf("连续输错5次将锁定账户8小时,您已输错%v次", num)
	data := strconv.Itoa(num)
	code = errors.User_3002
	if num == 5 {
		code = errors.User_3007
		message = "输入错误次数过多,账户被锁定到:" + frozenTime.Local().Format("2006-01-02 03:04:05 PM")
		data = frozenTime.Local().Format("2006-01-02 03:04:05 PM")
	}
	log.Info(data)
	return errors.ErrModel{Code: code, Err: errors.New(message), Data: map[string]interface{}{"data": data}}
}
