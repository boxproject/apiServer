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
package auth

import (
	"strconv"
	"strings"
	"github.com/boxproject/apiServer/db"
	reterror "github.com/boxproject/apiServer/errors"
	msg "github.com/boxproject/apiServer/service/message"
	"github.com/boxproject/apiServer/utils"
	log "github.com/alecthomas/log4go"
)

// GetAuthList get auth list
func GetAuthList(lang string) reterror.ErrModel {
	ds := &db.AuthService{}
	list, err := ds.GetAuthList()
	langTag := utils.ConvertLangTag(lang)
	if err != nil {
		return reterror.ErrModel{Code: reterror.Failed, Err: err}
	}
	authMap := make([]interface{}, len(list))
	for i, v := range list {
		authMap[i] = map[string]interface{}{
			"ID":    v.ID,
			"Name":  db.AuthMaps[langTag][strings.ToLower(v.AuthType)],
			"Count": ds.GetAuthAccountsCount(v.ID),
		}
	}
	return reterror.ErrModel{Err: nil, Code: reterror.Success, Data: authMap}
}

// GetAuthAccounts get accounts has authed
func GetAuthAccounts(authID int) reterror.ErrModel {
	if authID == 0 {
		return reterror.ErrModel{Code: reterror.ParamsNil, Err: reterror.PARAMS_NULL}
	}
	ds := &db.AuthService{}
	list, err := ds.GetAuthAccounts(authID)
	if err != nil {
		log.Error("get account auth by authid error", err)
		return reterror.ErrModel{Err: reterror.SYSTEM_ERROR, Code: reterror.Failed}
	}
	authAccount := make([]interface{}, len(list))
	for i, v := range list {
		authAccount[i] = map[string]interface{}{"ID": v.ID, "Name": v.Name, "AppId": v.AppId}
	}
	return reterror.ErrModel{Err: nil, Code: reterror.Success, Data: authAccount}
}

// AddAuthToAccount add auth to account
func AddAuthToAccount(auth PAuth, userType int, account string) reterror.ErrModel {
	log.Debug("AddAuthToAccount...")
	if auth.Add == "" || auth.ID == 0 {
		return reterror.ErrModel{Code: reterror.ParamsNil, Err: reterror.PARAMS_NULL}
	}
	ds := &db.AuthService{}
	parsedAdd := strings.Split(auth.Add, ",")
	add := make([]int, len(parsedAdd))
	for i, v := range parsedAdd {
		add[i], _ = strconv.Atoi(v)
	}
	result := ds.AddAuthToAccount(add, auth.ID)
	switch result {
	case 2:
		return reterror.ErrModel{Err: reterror.MSG_5003, Code: reterror.Auth_5003}
	case 3:
		return reterror.ErrModel{Err: reterror.MSG_5004, Code: reterror.Auth_5004}
	case 4:
		return reterror.ErrModel{Err: reterror.MSG_5005, Code: reterror.Auth_5005}
	}
	// send message
	err := msg.AddAuth(add, []int{auth.ID}, userType, account)
	if err != nil {
		log.Error("添加权限站内信通知错误:", err)
	}
	return reterror.ErrModel{Err: nil, Code: reterror.Success}
}

// DelAuthFromAccount remove auth from account
func DelAuthFromAccount(authParams PAuth, userType int, account string) reterror.ErrModel {
	log.Debug("DelAuthFromAccount...")
	if authParams.ID == 0 || authParams.AccountID == 0 {
		return reterror.ErrModel{Code: reterror.ParamsNil, Err: reterror.PARAMS_NULL}
	}
	ds := &db.AuthService{}
	err := ds.DelAuthFromAccount(authParams.ID, authParams.AccountID)
	if err != nil {
		log.Error("delete auth from account error", err)
		return reterror.ErrModel{Code: reterror.Auth_5006, Err: reterror.MSG_5006}
	}
	// send message
	msg.DelAuth(authParams.ID, authParams.AccountID, userType, account)
	return reterror.ErrModel{Err: nil, Code: reterror.Success}
}
