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

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
	"sync"
	log "github.com/alecthomas/log4go"
	"github.com/boxproject/apiServer/db"
	reterror "github.com/boxproject/apiServer/errors"
	"github.com/boxproject/apiServer/service/logger"
	msg "github.com/boxproject/apiServer/service/message"
	JWT "github.com/boxproject/apiServer/middleware"
	voucher "github.com/boxproject/apiServer/rpc"
	"github.com/boxproject/apiServer/service/template"
	"github.com/boxproject/apiServer/utils"
	"github.com/boxproject/apiServer/service/verify"
	"github.com/boxproject/apiServer/versionlog"
	"github.com/dgrijalva/jwt-go"
	walletCli "github.com/boxproject/boxwallet/cli"
	"github.com/boxproject/boxwallet/bccore"
	"github.com/boxproject/apiServer/config"
	"path"
	"github.com/boxproject/apiServer/common"
	"reflect"
	tmpl "text/template"
	"bytes"
)

const TOKEN_EXP = 24

var accLock = new(sync.Mutex)

// Signup sign up
func Signup(user *UserVo) reterror.ErrModel {
	log.Debug("Signup...", user)
	if user.AppId == "" || (user.Name == "") || (user.PubKey == "") || (user.Pwd == "") {
		log.Error("参数不能为空")
		return reterror.ErrModel{Code: reterror.ParamsNil, Err: reterror.PARAMS_NULL}
	}
	name := strings.TrimSpace(user.Name)
	rs := &db.RegDBService{}
	// if name exists
	nameExists, err := rs.NameExists(name)
	if err != nil {
		log.Error("注册错误", err)
		if nameExists {
			return reterror.ErrModel{Code: reterror.User_2002, Err: reterror.MSG_2002}
		} else {
			return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
		}
	}
	appIDExists, err := rs.AppIdExists(user.AppId)
	if err != nil {
		log.Error("AppID exists", err)
		if appIDExists {
			return reterror.ErrModel{Code: reterror.User_2004, Err: reterror.MSG_2004}
		} else {
			return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
		}
	}

	//generate salt
	salt := strconv.FormatInt(time.Now().Unix(), 10)
	pwdSalt := utils.PwdSalt(user.Pwd, salt)
	dep := &db.Registration{
		Name:   name,
		Pwd:    pwdSalt,
		Salt:   salt,
		AppId:  user.AppId,
		PubKey: user.PubKey,
	}
	err = rs.AddUser(dep)

	if err != nil {
		log.Error("Add user failed case:", err)
		return reterror.ErrModel{Err: reterror.SYSTEM_ERROR, Code: reterror.User_2001}
	}

	return reterror.ErrModel{Err: nil, Code: reterror.Success}
}

// Login login
func Login(user *UserVo) reterror.ErrModel {
	log.Debug("Login...", user)
	name := strings.TrimSpace(user.Name)
	password := user.Pwd
	// parameters shouldn't been nil
	if name == "" || password == "" {
		return reterror.ErrModel{Code: reterror.ParamsNil, Err: reterror.PARAMS_NULL}
	}
	// validate password
	verifyResult, account, errorCode := verify.VerifyPSW(name, password)

	if verifyResult.Result == false {
		log.Error("verify password error")
		return reterror.ErrModel{Code: errorCode, Err: errors.New(verifyResult.Reason), Data: map[string]interface{}{"data": verifyResult.Data}}
	}
	// generate token
	var j *JWT.JWT = &JWT.JWT{
		[]byte(JWT.GetSignKey()),
	}

	claims := JWT.CustomClaims{AppID: account.AppId, Account: name, UserType: account.UserType, StandardClaims: jwt.StandardClaims{ExpiresAt: time.Now().Add((TOKEN_EXP) * time.Hour).Unix()}}
	token, err := j.CreateToken(claims)

	if err != nil {
		log.Error("登录生成TOKEN", err)
		return reterror.ErrModel{Code: reterror.User_3003, Err: reterror.MSG_3003}
	}

	da := &db.AuthService{}
	dd := &db.DepService{}
	authMap, authMapErr := da.GetAuthMapByAccID(account.ID)
	if authMapErr != nil {
		log.Error("get auth map error", authMapErr)
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}

	authList := make([]interface{}, len(authMap))
	for i, v := range authMap {
		auth, authErr := da.GetAuthByID(v.AuthId)
		if authErr != nil {
			log.Error("get auth by id error", v.AuthId)
			return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
		}
		if len(auth) == 0 {
			authList[i] = map[string]interface{}{}
		}
		authList[i] = map[string]interface{}{
			"AuthName": auth[0].Name,
			"AuthId":   v.AuthId,
		}
	}
	// get department of user by department id
	dep, depErr := dd.GetDepByID(account.DepartmentId)
	if depErr != nil {
		log.Error("get dep by account id error")
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}

	departmentName, departmentID := "", 0
	if len(dep) != 0 {
		departmentName = dep[0].Name
		departmentID = dep[0].ID
	}

	return reterror.ErrModel{
		Err:  nil,
		Code: reterror.Success,
		Data: map[string]interface{}{
			"token":          token,
			"name":           name,
			"userType":       account.UserType,
			"auths":          authList,
			"departmentName": departmentName,
			"departmentId":   departmentID,
		},
	}
}

// ModifyPassword modify password
func ModifyPassword(mp *ModifyPSD, name string) reterror.ErrModel {
	log.Debug("ModifyPassword...")
	as := &db.AccDBService{}
	salt, pwd, _ := as.GetSalt(name)
	oldPwd := utils.PwdSalt(mp.OldPassword, salt)
	if pwd != oldPwd {
		return reterror.ErrModel{Code: reterror.User_3002, Err: reterror.MSG_3002}
	}
	newPwd := utils.PwdSalt(mp.NewPassword, salt)
	err := as.ModifyPassword(name, newPwd)
	if err == nil {
		return reterror.ErrModel{Err: nil, Code: reterror.Success}
	}
	log.Error("ModifyPassword error", err)
	return reterror.ErrModel{Code: reterror.User_3010, Err: reterror.MSG_3010}
}

// InsideLetter send message
func InsideLetter(lang string, accountId, userType, page int) reterror.ErrModel {
	langTag := utils.ConvertLangTag(lang)
	if accountId == 0 {
		return reterror.ErrModel{Code: reterror.ParamsNil, Err: reterror.PARAMS_NULL}
	}
	var letters []Letter
	var unreadnumber int // number of unread message
	var total int        // number of message
	var unReadWarn bool
	ms := &db.MessageService{}
	msgListAsc, msgListDesc, err := ms.LetterList()
	if err != nil {
		log.Error("获取站内信列表失败", err)
		return reterror.ErrModel{Code: reterror.ParamsNil, Err: reterror.PARAMS_NULL}
	}
	//组织list

	for _, v := range msgListDesc {
		var receivers []int
		json.Unmarshal([]byte(v.Receiver), &receivers)
		for _, vR := range receivers {
			if vR == accountId {
				var letter Letter
				var paddings common.MsgTemplatePadding
				json.Unmarshal([]byte(v.Padding), &paddings)
				if paddings.AccountType != "" {
					paddings.AccountType = msg.MsgAccountTypeMap[langTag][strings.ToLower(paddings.AccountType)].(string)
				}
				if len(paddings.AuthNames) > 0 {
					paddings.AuthName = paddings.AuthNames[0]
					if len(paddings.AuthNames) > 1 {
						for i := 1; i < len(paddings.AuthNames); i++ {
							paddings.AuthName = paddings.AuthName + ", " + paddings.AuthNames[i]
						}
					}
				}
				letter.ID = v.ID
				//letter.Title = titlebuf.String()
				letter.Time = v.CreatedAt
				letter.Content = v.Content
				letter.Type = v.Type
				letter.Param = v.Param
				// title
				titleValue := reflect.ValueOf(v.Title)
				titleKey := strings.ToLower(titleValue.String())
				titleTmpl, _ := tmpl.New("title").Parse(msg.MsgTitleTempMap[langTag][titleKey].(string))
				titlebuf := new(bytes.Buffer)
				titleTmpl.Execute(titlebuf, paddings)
				letter.Title = titlebuf.String()
				// content
				cValue := reflect.ValueOf(v.Content)
				cKey := strings.ToLower(cValue.String())
				cTmpl, _ := tmpl.New("con").Parse(msg.MsgContentTempMap[langTag][cKey].(string))
				cbuf := new(bytes.Buffer)
				cTmpl.Execute(cbuf, paddings)
				letter.Content = cbuf.String()
				if v.WarnType == common.ErrorMessage {
					letter.Warn = true
				}
				var readers []int
				json.Unmarshal([]byte(v.Reader), &readers)
				for _, id := range readers {
					if id == accountId {
						letter.Status = 1
					}
				}
				letters = append(letters, letter)
				break
			}
		}
	}

	// number of unread message
	var warnCount int
	var warnMsg []db.Message
	for i := 0; i < len(msgListAsc); i++ {
		var receivers []int
		json.Unmarshal([]byte(msgListAsc[i].Receiver), &receivers)
		for _, v := range receivers {
			if v == accountId {
				var readers []int
				var number int
				total++
				json.Unmarshal([]byte(msgListAsc[i].Reader), &readers)
				for _, id := range readers {
					if id == accountId {
						number = 1
					}
				}
				if number == 0 {
					if msgListAsc[i].WarnType == common.ErrorMessage {
						unReadWarn = true
					}
					unreadnumber++
				}

			}
			if msgListAsc[i].WarnType == common.ErrorMessage {
				receivers := []int{}
				json.Unmarshal([]byte(msgListAsc[i].Receiver), &receivers)
				if receivers[0] == accountId {
					warnCount = 0
					warnMsg = nil
				}
			}

			if msgListAsc[i].WarnType == common.WarnMessage {
				warnCount++
				warnMsg = append(warnMsg, msgListAsc[i])
			}
			break
		}
	}
	if warnCount >= 2 {
		for i := len(warnMsg) - 1; i > 0; i-- {
			if warnMsg[i].WarnType == common.WarnMessage && (userType == common.AdminAccType || userType == common.OwnerAccType) {
				// 插入报警站内信
				wmsg, err := msg.VoucherWarn(warnCount, accountId, warnMsg[i].Type, warnMsg[i].Param)
				if err != nil {
					log.Error("插入报警失败", err)
					return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
				}
				var warnLetter Letter
				var paddings common.MsgTemplatePadding
				json.Unmarshal([]byte(wmsg.Padding), &paddings)
				paddings.AuthName = paddings.AuthNames[0]
				if len(paddings.AuthNames) > 1 {
					for i := 1; i < len(paddings.AuthNames); i++ {
						paddings.AuthName = paddings.AuthName + ", " + paddings.AuthNames[i]
					}
				}
				if paddings.AccountType != "" {
					paddings.AccountType = msg.MsgAccountTypeMap[langTag][strings.ToLower(paddings.AccountType)].(string)
				}
				warnLetter.ID = wmsg.ID
				//warnLetter.Title = wmsg.Title
				warnLetter.Time = wmsg.CreatedAt
				//warnLetter.Content = wmsg.Content
				warnLetter.Type = wmsg.Type
				warnLetter.Param = wmsg.Param
				warnLetter.Warn = true
				// title
				titleWarnValue := reflect.ValueOf(wmsg.Title)
				titleWarnKey := strings.ToLower(titleWarnValue.String())
				titleWarnTmpl, _ := tmpl.New("title").Parse(msg.MsgTitleTempMap[langTag][titleWarnKey].(string))
				titleWarnBuf := new(bytes.Buffer)
				titleWarnTmpl.Execute(titleWarnBuf, paddings)
				warnLetter.Title = titleWarnBuf.String()
				// content
				cWarnValue := reflect.ValueOf(wmsg.Content)
				cWarnKey := strings.ToLower(cWarnValue.String())
				cWarnTmpl, _ := tmpl.New("con").Parse(msg.MsgContentTempMap[langTag][cWarnKey].(string))
				cWarnbuf := new(bytes.Buffer)
				cWarnTmpl.Execute(cWarnbuf, paddings)
				warnLetter.Content = cWarnbuf.String()
				letters = append(letters, warnLetter)
				break
			}
		}
		warnCount = 0
	}
	// add paging
	var slice []Letter
	length := len(letters)
	if length-page*10 > 0 && length-page*10 >= 10 {
		slice = letters[page*10 : page*10+10]
	}
	if length-page*10 > 0 && length-page*10 < 10 {
		slice = letters[page*10 : length]
	}
	if length-page*10 < 0 {
		slice = nil
	}
	data := make(map[string]interface{})
	data["List"] = slice
	data["Total"] = total
	data["UnReadNumber"] = unreadnumber
	data["UnReadWarn"] = unReadWarn
	// get admin count ??
	da := db.AccDBService{}
	count, accErr := da.GetAdminCount()
	if accErr != nil {
		log.Error("ger admin account info error", accErr)
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	data["AdminCount"] = count
	return reterror.ErrModel{Err: nil, Code: reterror.Success, Data: data}
}

// ReadLetter mark message as readed
func ReadLetter(accountId, id int) reterror.ErrModel {
	if accountId == 0 || id == 0 {
		return reterror.ErrModel{Code: reterror.ParamsNil, Err: reterror.PARAMS_NULL}
	}
	ms := &db.MessageService{}
	message, err := ms.GetMessageById(id)
	if err != nil {
		log.Error("ger message by id error", err)
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	var readers []int
	json.Unmarshal([]byte(message.Reader), &readers)
	for _, reader := range readers {
		if reader == accountId {
			return reterror.ErrModel{Code: reterror.MessageRead, Err: reterror.MSG_1011}
		}
	}
	readers = append(readers, accountId)

	json, _ := json.Marshal(readers)
	err = ms.ReadLetter(id, string(json))
	if err != nil {
		log.Error("read letter error", err)
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	return reterror.ErrModel{Err: nil, Code: reterror.Success}
}

// Scan submit member info for qrcode scaned
func Scan(userVo *UserVo) reterror.ErrModel {
	log.Debug("Scan...")
	rs := &db.RegDBService{}
	// app id of new member
	if userVo.NewAppId == "" || userVo.Msg == "" {
		log.Error("参数错误")
		return reterror.ErrModel{Code: reterror.ParamsNil, Err: reterror.PARAMS_NULL}
	}
	//get public key of new member
	user, err := rs.FindUserByAppId(userVo.NewAppId)
	if err != nil {
		log.Error("查询公钥错误", err)
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}

	if user == nil {
		return reterror.ErrModel{Code: reterror.User_3004, Err: reterror.MSG_3004}
	}
	if user.Status == common.RegToApproval {
		return reterror.ErrModel{Code: reterror.User_3013, Err: reterror.MSG_3013}
	}
	if user.Status == common.RegApprovaled {
		return reterror.ErrModel{Code: reterror.User_3016, Err: reterror.MSG_3016}
	}

	reg := &db.Registration{
		ID:            user.ID,
		Status:        common.RegToApproval,
		Msg:           userVo.Msg,
		Level:         userVo.Level + 1,
		SourceAppId:   userVo.AppId,
		SourceAccount: userVo.Name,
	}
	rs.EditUser(reg)

	return reterror.ErrModel{Err: nil, Code: reterror.Success}

}

// FindUserStatus check user registered status
func FindUserStatus(userVo *UserVo) reterror.ErrModel {
	rs := &db.RegDBService{}
	// required parameters
	if userVo.AppId == "" {
		return reterror.ErrModel{Code: reterror.ParamsNil, Err: reterror.PARAMS_NULL}
	}
	user, err := rs.FindUserByAppId(userVo.AppId)
	if err != nil {
		log.Error("查询错误", err)
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	if user == nil {
		return reterror.ErrModel{Code: reterror.User_3004, Err: reterror.MSG_3004}
	}
	token := ""
	if user.Status == common.RegApprovaled {
		var j *JWT.JWT = &JWT.JWT{
			[]byte(JWT.GetSignKey()),
		}
		claims := JWT.CustomClaims{AppID: user.AppId, Account: user.Name, UserType: 0, StandardClaims: jwt.StandardClaims{ExpiresAt: time.Now().Add((TOKEN_EXP) * time.Hour).Unix()}}
		token, _ = j.CreateToken(claims)
		log.Debug("用户登录成功", user.Name, token)
	}
	return reterror.ErrModel{Err: nil, Code: reterror.Success, Data: map[string]interface{}{"status": user.Status, "refuse_reason": user.RefuseReason, "token": token}}
}

// VerifyUser verify user
func VerifyUser(userVo *UserVo) reterror.ErrModel {
	log.Debug("Verify user...")
	accLock.Lock()
	defer accLock.Unlock()
	rs := &db.RegDBService{}
	as := &db.AccDBService{}
	ds := &db.DepService{}
	// required parameters
	if userVo.NewAppId == "" || (userVo.Status == common.RegApprovaled && userVo.DepID == 0) {
		return reterror.ErrModel{Code: reterror.ParamsNil, Err: reterror.PARAMS_NULL}
	}
	// validate input params
	if (userVo.Status != common.RegApprovaled) && (userVo.Status != common.RegRejected) {
		return reterror.ErrModel{Code: reterror.ParamsNil, Err: reterror.PARAMS_NULL}
	}
	user, err := rs.FindUserByAppId(userVo.NewAppId)
	if err != nil {
		log.Error("查询错误", err)
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}

	if user == nil {
		return reterror.ErrModel{Code: reterror.User_3004, Err: reterror.MSG_3004}
	}
	if user.Status == common.RegWaitingToApproval { // not submit yet
		return reterror.ErrModel{Code: reterror.User_3008, Err: reterror.MSG_3008}
	}
	if user.Status == common.RegApprovaled || user.Status == common.RegRejected { // has been verified
		return reterror.ErrModel{Code: reterror.User_3017, Err: reterror.MSG_3017}
	}

	dbaccount := &db.Account{}
	dbuser := &db.Registration{}
	if userVo.Status == common.RegApprovaled {
		dep, err := ds.GetDepByID(userVo.DepID)
		if err != nil || len(dep) == 0 {
			return reterror.ErrModel{Code: reterror.SetDepFailed, Err: reterror.MSG_1005}
		}
		dbaccount = &db.Account{
			AppId:        user.AppId,
			Name:         user.Name,
			Pwd:          user.Pwd,
			Salt:         user.Salt,
			PubKey:       user.PubKey,
			Msg:          user.Msg,
			FrozenTo:     time.Now(),
			Level:        user.Level,
			DepartmentId: userVo.DepID,
			SourceAppId:  user.SourceAppId,
		}
		dbuser = &db.Registration{
			ID:     user.ID,
			Status: common.RegApprovaled,
		}
	} else {
		dbuser = &db.Registration{
			ID:           user.ID,
			Status:       common.RegRejected,
			RefuseReason: userVo.RefuseReason,
		}
	}
	ret := as.VerifyUser(dbaccount, dbuser)
	if ret {
		return reterror.ErrModel{Err: nil, Code: reterror.Success}
	} else {
		return reterror.ErrModel{Code: reterror.User_3008, Err: reterror.MSG_3008}
	}

}

// RegList get users to be verified
func RegList(userVo *UserVo) reterror.ErrModel {
	rs := &db.RegDBService{}
	dbuser := &db.Registration{
		Status: userVo.Status,
	}
	List, error := rs.RegList(dbuser)
	if error == nil {
		return reterror.ErrModel{Err: nil, Code: reterror.Success, Data: List}
	} else {
		return reterror.ErrModel{Code: reterror.User_3008, Err: reterror.MSG_3008}
	}

}

// UserTree generate user tree structure by level
func UserTree(userPubKey string) reterror.ErrModel {
	log.Debug("userTree...")
	length := 0
	as := &db.AccDBService{}
	accounts, err := as.FindAllUsers()
	if err != nil {
		log.Error("查找用户列表错误", err)
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	var mapAccount map[int][]AccountTree
	mapAccount = make(map[int][]AccountTree)
	for i := 0; i < len(accounts); i++ {
		if i == 0 {
			length = accounts[0].Level
		}
		account := &AccountTree{
			Name:        accounts[i].Name,
			AppId:       accounts[i].AppId,
			SourceAppId: accounts[i].SourceAppId,
			Msg:         accounts[i].Msg,
			PubKey:      accounts[i].PubKey,
		}
		createMaps(accounts[i].Level, account, mapAccount)
	}
	for i := length - 1; i > 0; i-- {
		// create user data reverse
		createTree(i, mapAccount)

	}
	var mapAccountTree []AccountTree
	mapAccountTree = mapAccount[1]
	// add voucher public key
	voucherPubKey, err := getCorePubkeyInfo()
	if err != nil || len(voucherPubKey) == 0 {
		log.Error("获取签名机公钥", err)
		return reterror.ErrModel{Code: reterror.Failed}
	}
	enPubkey, enAesKey, err := encodeVoucherPubkey([]byte(userPubKey), voucherPubKey)
	if err != nil {
		log.Error("加密公钥", err)
		return reterror.ErrModel{Code: reterror.Failed}
	}
	data := &AccountTreeWithCoreInfo{}
	data.AccTree = mapAccountTree
	data.Voucher.ENPublickey = string(enPubkey)
	data.Voucher.ENKey = string(enAesKey)
	log.Debug("userTree", data)
	return reterror.ErrModel{Err: nil, Code: reterror.Success, Data: data}

}

// generate tree struct
func createTree(key int, mapAccount map[int][]AccountTree) {
	parentAppId := ""
	for j := 0; j < len(mapAccount[key]); j++ {
		parentAppId = mapAccount[key][j].AppId
		for k := 0; k < len(mapAccount[key+1]); k++ {
			if parentAppId == mapAccount[key+1][k].SourceAppId {
				if len(mapAccount[key][j].Children) > 0 {
					mapAccount[key][j].Children = append(mapAccount[key][j].Children, mapAccount[key+1][k])
				} else {
					mapAccount[key][j].Children = []AccountTree{mapAccount[key+1][k]}
				}
			}
		}
	}
}

func createMaps(key int, account *AccountTree, mapAccount map[int][]AccountTree) {
	accounts, ok := mapAccount[key]
	if ok {
		accounts := append(accounts, *account)
		mapAccount[key] = accounts
	} else {
		var newAccounts = []AccountTree{*account}
		mapAccount[key] = newAccounts
	}
}

// GetAccountsByType get users by user type
func GetAccountsByType(userType int) reterror.ErrModel {
	ds := &db.AccDBService{}
	accounts, err := ds.GetAccountsByType(userType)
	if err != nil {
		log.Error("get accounts info by type error", err)
		return reterror.ErrModel{Err: reterror.SYSTEM_ERROR, Code: reterror.Failed}
	}
	rAccounts := make([]interface{}, len(accounts))
	for i, v := range accounts {
		rAccounts[i] = map[string]interface{}{
			"ID":       v.ID,
			"Name":     v.Name,
			"UserType": v.UserType,
			"AppId":    v.AppId,
		}
	}
	return reterror.ErrModel{Err: nil, Code: reterror.Success, Data: rAccounts}
}

// AddAdmin add
func AddAdmin(add string, accountID int) reterror.ErrModel {
	log.Debug("AddAdmin...")
	if add == "" {
		return reterror.ErrModel{Code: reterror.ParamsNil, Err: reterror.PARAMS_NULL}
	}
	parsedAdd := strings.Split(add, ",")
	ds := &db.AccDBService{}
	for _, v := range parsedAdd {
		accountID, _ := strconv.Atoi(v)
		err := ds.AddAdmin(accountID)
		if err != nil {
			return reterror.ErrModel{Err: reterror.MSG_3011, Code: reterror.User_3011}
		}
	}
	// send message
	err := msg.AddAdmin(add, accountID)
	if err != nil {
		return reterror.ErrModel{Err: reterror.MSG_1010, Code: reterror.MessageFailed}
	}
	return reterror.ErrModel{Err: nil, Code: reterror.Success}
}

// DelAdmin delete admin
func DelAdmin(accountID, id int) reterror.ErrModel {
	log.Debug("DelAdmin...")
	if accountID == 0 {
		return reterror.ErrModel{Code: reterror.ParamsNil, Err: reterror.PARAMS_NULL}
	}
	ds := &db.AccDBService{}
	err := ds.DelAdmin(accountID)
	if err != nil {
		log.Error("delete admin error", err)
		return reterror.ErrModel{Err: reterror.MSG_3012, Code: reterror.User_3012}
	}
	// send message
	err = msg.DelAdmin(accountID, id)
	if err != nil {
		log.Error("send message error", err)
		return reterror.ErrModel{Err: reterror.MSG_1010, Code: reterror.MessageFailed}
	}
	return reterror.ErrModel{Err: err, Code: reterror.Success}
}

// ReSignUp re-signup
func ReSignUp(appID string) reterror.ErrModel {
	log.Debug("ReSignUp...")
	if appID == "" {
		return reterror.ErrModel{Code: reterror.ParamsNil, Err: reterror.PARAMS_NULL}
	}
	ds := &db.RegDBService{}
	err := ds.ReSignUp(appID)
	if err != nil {
		log.Error("get resign info error", err)
		return reterror.ErrModel{Err: reterror.MSG_2005, Code: reterror.User_2005}
	}
	return reterror.ErrModel{Err: err, Code: reterror.Success}
}

// GetAllUsers get available users
func GetAllUsers() reterror.ErrModel {
	ds := &db.AccDBService{}
	list, err := ds.GetAllUsers()
	users := make([]interface{}, len(list))
	for i, v := range list {
		users[i] = map[string]interface{}{
			"Name":     v.Name,
			"ID":       v.ID,
			"AppID":    v.AppId,
			"UserType": v.UserType,
			"Frozen":   v.Frozen,
		}
	}
	if err != nil {
		log.Debug("get all users error", err)
		return reterror.ErrModel{Err: reterror.MSG_3014, Code: reterror.User_3014}
	}
	return reterror.ErrModel{Err: err, Code: reterror.Success, Data: users}
}

// DelayTaskNum get to-do tasks
func DelayTaskNum(account string, usertype, accountID int) reterror.ErrModel {
	data := make(map[string]interface{})
	ts := &db.TransferDBService{}
	ds := &db.AccDBService{}
	tt := &db.TemplateDBService{}
	if usertype == common.NormalAccType { // member
		trTask := ts.DelayTaskNum(account)
		data["Transfer"] = trTask
	}
	if usertype == common.AdminAccType { // admin
		trTask := ts.DelayTaskNum(account)
		data["Transfer"] = trTask
		acTask := ds.DelayTaskNum()
		data["Account"] = acTask
	}
	if usertype == common.OwnerAccType { // owner
		trTask := ts.DelayTaskNum(account)
		data["Transfer"] = trTask
		acTask := ds.DelayTaskNum()
		data["Account"] = acTask
		tpTask := tt.DelayTaskNum(accountID)
		data["Template"] = tpTask
	}
	return reterror.ErrModel{Err: nil, Code: reterror.Success, Data: data}
}

// GetUserByID get user info by id
func GetUserByID(id int) reterror.ErrModel {
	if id == 0 {
		return reterror.ErrModel{Code: reterror.ParamsNil, Err: reterror.PARAMS_NULL}
	}
	ds := &db.AccDBService{}
	da := &db.AuthService{}
	dd := &db.DepService{}

	userInfo, err := ds.GetUserByID(id)
	if err != nil {
		log.Error("get user info by id error", err)
		return reterror.ErrModel{Err: reterror.MSG_3014, Code: reterror.User_3014}
	}
	if len(userInfo) == 0 {
		return reterror.ErrModel{Err: reterror.MSG_3001, Code: reterror.User_3001}
	}
	// get all auth account has
	authMap, authMapErr := da.GetAuthMapByAccID(userInfo[0].ID)
	if authMapErr != nil {
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	// 获取权限
	authList := make([]interface{}, len(authMap))
	for i, v := range authMap {
		auth, authErr := da.GetAuthByID(v.AuthId)
		if authErr != nil {
			log.Error("get auth error", err)
			return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
		}
		if len(auth) == 0 {
			authList[i] = map[string]interface{}{}
		}
		authList[i] = map[string]interface{}{
			"Name": auth[0].Name,
			"ID":   v.AuthId,
		}
	}
	// get department
	dep, depErr := dd.GetDepByID(userInfo[0].DepartmentId)
	if depErr != nil {
		log.Error("ger user department error", err)
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	departmentName, departmentID := "", 0
	if len(dep) != 0 {
		departmentName = dep[0].Name
		departmentID = dep[0].ID
	}

	return reterror.ErrModel{Err: err, Code: reterror.Success, Data: map[string]interface{}{
		"Name":           userInfo[0].Name,
		"AppId":          userInfo[0].AppId,
		"UserType":       userInfo[0].UserType,
		"Auths":          authList,
		"DepartmentName": departmentName,
		"DepartmentId":   departmentID,
		"ID":             userInfo[0].ID,
	},
	}
}

// SetUser set user info
func SetUser(user *PAccount, operatorId int) reterror.ErrModel {
	log.Debug("SetUser...")
	ds := &db.AccDBService{}
	accInfo, err := ds.AccountByID(*user.ID)
	err = ds.SetUser(*user.AuthID, *user.DepID, *user.ID, *user.UserType)
	if err != nil {
		log.Error("set user fail", err)
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	// 发送站内信
	if accInfo.UserType != *user.UserType {
		// 设置或删除管理员权限
		if accInfo.UserType == common.NormalAccType {
			// 设置为管理员
			msg.AddAdmin(utils.ToStr(*user.ID), operatorId)
		} else {
			// 取消管理员
			msg.DelAdmin(*user.ID, operatorId)
		}
	}

	authAdd := strings.Split(*user.AuthID, ",")
	authIds := make([]int, len(authAdd))
	for i, v := range authAdd {
		authIds[i], _ = strconv.Atoi(v)
	}
	operatorInfo, err := ds.AccountByID(operatorId)
	if err != nil {
		log.Error("find account info error", err)
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	if len(authIds) > 0 && authIds[0] != 0 {
		err = msg.AddAuth([]int{*user.ID}, authIds, operatorInfo.UserType, operatorInfo.Name)
		if err != nil {
			log.Error("send message error", err)
			return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
		}
	}
	return reterror.ErrModel{Err: err, Code: reterror.Success}
}

// DisableAcc disable account
func DisableAcc(accID int, adminAcc string) reterror.ErrModel {
	log.Debug("DisableAcc...", accID, adminAcc)
	if accID == 0 {
		return reterror.ErrModel{Code: reterror.ParamsNil, Err: reterror.PARAMS_NULL}
	}
	ds := &db.AccDBService{}

	// get user info by id
	accInfo, err := ds.AccountByID(accID)
	if err != nil {
		log.Error("根据ID获取账户信息", accID, err)
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	// disable related template and transfer
	err = template.DisableTemplateByUserName(accInfo, adminAcc)
	if err != nil {
		log.Error("根据用户名作废所涉及的审批流模板", err)
		if err == reterror.MSG_9009 {
			return reterror.ErrModel{Code: reterror.Template_9009, Err: reterror.MSG_9009}
		}
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}

	err = ds.DisableAcc(accID)
	if err != nil {
		log.Error("停用账户", err)
		return reterror.ErrModel{Code: reterror.User_3015, Err: reterror.MSG_3015}
	}
	return reterror.ErrModel{Err: err, Code: reterror.Success}
}

// GetRegList account list to be verified
func GetRegList() reterror.ErrModel {
	ds := &db.RegDBService{}
	list, err := ds.GetRegList()
	users := make([]interface{}, len(list))
	for i, v := range list {
		users[i] = map[string]interface{}{
			"Name":          v.Name,
			"CreatedAt":     v.UpdatedAt, // submit times
			"Status":        v.Status,
			"SourceAccount": v.SourceAccount,
			"AppId":         v.AppId,
		}
	}
	if err != nil {
		log.Error("get reg list error", err)
		return reterror.ErrModel{Code: reterror.User_2006, Err: reterror.MSG_2006}
	}
	return reterror.ErrModel{Err: err, Code: reterror.Success, Data: users}
}

// getCorePubkeyInfo get and save voucher public key
func getCorePubkeyInfo() ([]byte, error) {
	// get local public key
	var voucherPubKey []byte
	//var voucherPubKey_bs64 string
	//parentPath := utils.GetParentDirectory(currentPath)
	currentPath := utils.GetCurrentDirectory()
	filename := path.Join(currentPath, common.CONF_PATH, common.VOUCHER_KEY_FILE)
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Error("读文件", err)
		return nil, err
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		log.Error("read core.key file", err)
		return nil, err
	}

	if len(buf) == 0 {
		// request voucher public key
		oper := &voucher.GrpcServer{
			Type:      voucher.VOUCHER_OPERATE_MASTER_PUBKEY,
			Timestamp: time.Now().UnixNano(),
		}
		_, voucherRes, voucherRet := voucher.SendVoucherData(oper)
		if voucherRet == voucher.VRET_ERR {
			log.Error("voucher return error", voucherRet)
			return nil, errors.New("Voucher Return Error")
		} else {
			if voucherRes.Status != voucher.VRET_ERR {
				log.Error("voucher status error", voucherRes.Status)
				return nil, errors.New("Voucher Status Error")
			}
		}
		// voucher public key
		voucherPubKey = voucherRes.Other
		// write file
		_, err = f.Write(voucherPubKey)
		defer f.Close()
		if err != nil {
			log.Error("write file error", err)
			return nil, err
		}

	} else {
		voucherPubKey = buf

	}
	//voucherPubKey_bs64 = base64.StdEncoding.EncodeToString(voucherPubKey)
	return voucherPubKey, nil
}

// encodeVoucherPubkey encode voucher public key
func encodeVoucherPubkey(userPubkey, voucherPubKey []byte) (enPubkey, enAesKey []byte, err error) {
	aesKeyBytes := utils.GetAesKeyRandom(utils.AESKEY_LEN_32)
	enAesKey, err = utils.RsaEncrypt(userPubkey, aesKeyBytes)
	if err != nil {
		log.Error("加密aes Key error", err)
		return
	}
	enPubkey, err = utils.AESCBCEncrypter(voucherPubKey, aesKeyBytes)
	if err != nil {
		log.Error("加密签名机公钥", err)
		return
	}
	return
}

// RecoveryOwner owner account recovery
func RecoveryOwner(user *UserVo) reterror.ErrModel {
	log.Debug("recover owner account...")
	if user.AppId == "" || user.Name == "" || user.PubKey == "" || user.Pwd == "" {
		return reterror.ErrModel{Code: reterror.ParamsNil, Err: reterror.PARAMS_NULL}
	}
	ds := &db.RegDBService{}
	da := &db.AccDBService{}
	// isOwner
	isOwner, existErr := da.IsOwner(user.Name)
	if existErr != nil {
		log.Error("判断是否为股东", existErr)
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	if !isOwner {
		return reterror.ErrModel{Code: reterror.User_2008, Err: reterror.MSG_2008}
	}
	// check recovery status
	regStatus, regErr := ds.CheckRecovStatus(user.AppId)
	if regStatus == common.OwnerRegError || regErr != nil {
		log.Error("查看恢复账号审批状态", regErr)
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	// rejected or unsubmit
	if regStatus != common.RegRejected && regStatus != common.OwnerRegNotFound {
		log.Debug("需重新注册或扫码提交认证", regStatus)
		return reterror.ErrModel{Code: reterror.User_2011, Err: reterror.MSG_2011}
	}

	salt := strconv.FormatInt(time.Now().Unix(), 10)
	pwdSalt := utils.PwdSalt(user.Pwd, salt)
	reg := &db.OwnerReg{
		Name:   user.Name,
		Pwd:    pwdSalt,
		Salt:   salt,
		AppId:  user.AppId,
		PubKey: user.PubKey,
	}
	id, err := ds.AddReg(reg)
	if err != nil {
		log.Error("添加恢复股东注册信息", err)
		return reterror.ErrModel{Code: reterror.User_2007, Err: reterror.MSG_2007}
	}

	// encode voucher public key
	mainPubkey, err := getCorePubkeyInfo()
	if err != nil {
		log.Debug("获取签名机公钥错误", err)
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	// encode voucher public key
	enPubkey, enAeskey, err := encodeVoucherPubkey([]byte(user.PubKey), mainPubkey)
	if err != nil {
		log.Debug("加密签名机公钥", err)
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	return reterror.ErrModel{Err: err, Code: reterror.Success, Data: map[string]interface{}{"RegId": id, "EnPubkey": string(enPubkey), "EnAesKey": string(enAeskey)}}
}

// RecoveryList list of recovery account
func RecoveryList(account string, pubkey string) reterror.ErrModel {
	ds := &db.RegDBService{}
	list, err := ds.RecoveryList()
	if err != nil {
		log.Error("获取待恢复股东列表", err)
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	var recoverySatus int
	accounts := make([]interface{}, len(list))
	for i, v := range list {
		aesKeyBytes := utils.GetAesKeyRandom(utils.AESKEY_LEN_32)
		enAesKey, err := utils.RsaEncrypt([]byte(pubkey), aesKeyBytes)
		if err != nil {
			log.Error("加密aes Key error", err)
			return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
		}
		enPubKey, err := utils.AESCBCEncrypter([]byte(v.PubKey), aesKeyBytes)
		// 恢复总状态
		recoverySatus = v.Status
		// 操作者操作状态
		status, statErr := ds.OperateStatus(account, v.ID)
		if statErr != nil {
			log.Debug("当前股东审批状态", statErr)
			return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
		}
		switch {
		case status == common.OwnerAccApprovaling && recoverySatus != common.RecoveryInvalid:
			recoverySatus = common.RegToApproval
		case status == common.OwnerAccApprovaled && recoverySatus != common.RecoveryReject:
			recoverySatus = common.HasApproved
		}
		accounts[i] = map[string]interface{}{
			"ID":        v.ID,
			"Name":      v.Name,
			"UpdatedAt": v.UpdatedAt,
			"Status":    recoverySatus,
			"AppId":     v.AppId,
			"EnAesKey":  string(enAesKey),
			"EnPubKey":  string(enPubKey),
		}
	}
	if err != nil {
		log.Debug("恢复股东列表", err)
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	return reterror.ErrModel{Err: err, Code: reterror.Success, Data: accounts}
}

// SubRecovery submit recovery account info
func SubRecovery(appID string, account string, regID int, operator string) reterror.ErrModel {
	log.Debug("SubRecovery...")
	if appID == "" || account == "" || regID == 0 {
		return reterror.ErrModel{Code: reterror.ParamsNil, Err: reterror.PARAMS_NULL}
	}
	ds := db.RegDBService{}
	da := db.AccDBService{}
	regStatus, _, regErr := ds.RegStatus(regID)
	if regErr != nil {
		log.Debug("审批状态", regErr)
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	existAccInfo, err := da.ExistAccount(account)
	if err != nil {
		log.Error("verify account exist error", err)
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	if existAccInfo == nil { // account doesn't exist
		log.Error("股东账号不存在", account)
		return reterror.ErrModel{Code: reterror.User_2008, Err: reterror.MSG_2008}
	}
	if existAccInfo.IsDeleted == 1 {
		log.Error("股东账号已被删除", account)
		return reterror.ErrModel{Code: reterror.User_2008, Err: reterror.MSG_2008}
	}
	if regStatus != 0 { // submit repeatly
		return reterror.ErrModel{Code: reterror.User_2010, Err: reterror.MSG_2010}
	}
	err = ds.SubRecovery(appID, account, regID)
	if err != nil {
		log.Debug("提交确认信息(关系表初始化数据)", err)
		return reterror.ErrModel{Code: reterror.User_2009, Err: reterror.MSG_2009}
	}
	logErr := logger.AddLog("recovery", "", operator, common.LoggerOwnerSubmit, strconv.Itoa(regID))
	if logErr != nil {
		log.Debug("股东恢复日志错误", logErr)
	}
	return reterror.ErrModel{Err: nil, Code: reterror.Success}
}

// VerifyPwd verify password and return user info
func VerifyPwd(pwd string, account string, id int) reterror.ErrModel {
	if pwd == "" || account == "" {
		return reterror.ErrModel{Code: reterror.ParamsNil, Err: reterror.PARAMS_NULL}
	}
	// 密码校验
	verified, _, errorCode := verify.VerifyPSW(account, pwd)
	if verified.Result == false {
		return reterror.ErrModel{Code: errorCode, Err: errors.New(verified.Reason), Data: map[string]interface{}{"data": verified.Data}}
	}

	ds := &db.AccDBService{}
	da := &db.AuthService{}
	dd := &db.DepService{}

	userInfo, err := ds.GetUserByID(id)
	if err != nil {
		log.Error("get user info by id error", err)
		return reterror.ErrModel{Err: reterror.MSG_3014, Code: reterror.User_3014}
	}
	if len(userInfo) == 0 {
		return reterror.ErrModel{Err: reterror.MSG_3001, Code: reterror.User_3001}
	}
	// 获取用户所有权限
	authMap, authMapErr := da.GetAuthMapByAccID(userInfo[0].ID)
	if authMapErr != nil {
		log.Error("get auth error", authMapErr)
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	// 获取权限
	authList := make([]interface{}, len(authMap))
	for i, v := range authMap {
		auth, authErr := da.GetAuthByID(v.AuthId)
		if authErr != nil {
			log.Error("get auth by id error", authErr)
			return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
		}
		if len(auth) == 0 {
			authList[i] = map[string]interface{}{}
		}
		authList[i] = map[string]interface{}{
			"AuthName": auth[0].Name,
			"AuthId":   v.AuthId,
		}
	}
	// get department
	dep, depErr := dd.GetDepByID(userInfo[0].DepartmentId)
	department := make([]interface{}, 1)
	if depErr != nil {
		log.Error("get dep error", depErr)
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	if len(dep) != 0 {
		department[0] = map[string]interface{}{
			"Department": dep[0].Name,
			"ID":         dep[0].ID,
		}
	} else {
		department = []interface{}{}
	}

	return reterror.ErrModel{Err: err, Code: reterror.Success, Data: map[string]interface{}{
		"Name":       userInfo[0].Name,
		"AppId":      userInfo[0].AppId,
		"UserType":   userInfo[0].UserType,
		"Auths":      authList,
		"Department": department,
		"ID":         userInfo[0].ID,
	},
	}
}

// RecoveryResult get recovery result
func RecoveryResult(regID int, appID string) reterror.ErrModel {
	if regID == 0 || appID == "" {
		return reterror.ErrModel{Code: reterror.ParamsNil, Err: reterror.PARAMS_NULL}
	}
	ds := db.RegDBService{}
	status, _, statusErr := ds.RegStatus(regID)
	if statusErr != nil {
		log.Debug("审批状态", statusErr)
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	list, err := ds.RecoveryResult(regID, appID)
	if err != nil {
		log.Debug("获取恢复股东认证结果", err)
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	return reterror.ErrModel{Err: err, Code: reterror.Success, Data: map[string]interface{}{
		"RecoveryStatus": status,
		"List":           list,
	}}
}

// VerifyRecovery verify owner
func VerifyRecovery(account string, voucherVo *VoucherVo) reterror.ErrModel {
	log.Debug("VerifyRecovery...")
	accLock.Lock()
	defer accLock.Unlock()
	// check if account isFrozen
	accErr := verify.IsFrozen(voucherVo.AppName)
	if accErr.Code != reterror.Success {
		return accErr
	}
	if voucherVo.RegId == 0 || (voucherVo.Status != VerifyOwnerApprove && voucherVo.Status != VerifyOwnerReject) || account == "" {
		return reterror.ErrModel{Code: reterror.ParamsNil, Err: reterror.PARAMS_NULL}
	}
	da := db.RegDBService{}
	status, operator, regErr := da.RegStatus(voucherVo.RegId)

	if operator == voucherVo.AppName {
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	if regErr != nil {
		log.Error("get reg status error", regErr)
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	if status == common.RegRejected {
		return reterror.ErrModel{Err: reterror.MSG_2015, Code: reterror.User_2015}
	}
	// verify owner
	verifyErr := VoucherVerifyRecovery(account, voucherVo)
	if verifyErr.Err != nil {
		return verifyErr
	}
	// resetting wrong password record
	ds := db.AccDBService{}
	_, err := ds.ResetAccount(voucherVo.AppId)
	if err != nil {
		log.Info("取消账户冻结错误", err)
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	// add log
	loggerType := common.LoggerOwnerAccept
	if voucherVo.Status == VerifyOwnerReject {
		loggerType = common.LoggerOwnerReject
	}
	logErr := logger.AddLog("recovery", "", account, loggerType, strconv.Itoa(voucherVo.RegId))
	if logErr != nil {
		log.Error("日志未记录ownerScan", logErr)
	}
	return reterror.ErrModel{Err: nil, Code: reterror.Success}
}

// VoucherVerifyRecovery recovery owner (voucher part)
func VoucherVerifyRecovery(account string, voucherVo *VoucherVo) reterror.ErrModel {
	// other owner verify
	oper := &voucher.GrpcServer{
		Type:      voucher.VOUCHER_OPERATE_RECOVER,
		AppId:     voucherVo.AppId,
		AppName:   voucherVo.AppName,
		AesKey:    []byte(voucherVo.AesKey),
		Msg:       []byte(voucherVo.Msg),
		Sign:      []byte(voucherVo.Sign),
		Timestamp: voucherVo.Timestamp,
	}
	_, result, ret := voucher.SendVoucherData(oper)
	if ret == 0 || ret == 4 {
		log.Debug("请求签名机", ret)
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	ds := db.RegDBService{}

	switch result.Status {
	case voucher.STATUS_APP_VERIFY_ERROR:
		return reterror.ErrModel{Code: reterror.Code_102, Err: reterror.MSG_102}
	case voucher.STATUS_APP_PASSWORD_ERROR:
		keywordErr := verify.Keyword(voucherVo.AppId, voucherVo.AppName)
		return keywordErr
	case voucher.STATUS_APP_NAME_NOTMATCH:
		return reterror.ErrModel{Code: reterror.Code_106, Err: reterror.MSG_106}
	case voucher.STATUS_APP_SAMENAME:
		return reterror.ErrModel{Code: reterror.Code_107, Err: reterror.MSG_107}
	}
	if result.Status != voucher.STATUS_OK { // fail to recovery
		log.Debug("签名机恢复失败", result.Status)
		// change recovery status to failed
		regErr := ds.UpdateRecoveryStatus(voucherVo.RegId)
		if regErr != nil {
			log.Error("修改恢复申请状态失败, id: ", voucherVo.RegId, result.Status)
		}
	}
	// get voucher status and check if all owner has passed recovery apply
	option := &voucher.GrpcServer{
		Type:      voucher.VOUCHER_STATUS,
		Timestamp: time.Now().UnixNano(),
	}
	voucherStatus, _, voucherRet := voucher.SendVoucherData(option)
	if voucherRet == 0 {
		log.Error("获取签名机状态失败")
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	// get recovery status
	verifyStatus := common.OwnerAccApprovaling
	other := int(binary.BigEndian.Uint64(result.Other))
	if (voucherStatus.Total - 1) == other {
		verifyStatus = common.OwnerAccApprovaled
	}
	err := ds.VerifyRecovery(account, voucherVo.Status, voucherVo.RegId, verifyStatus)
	if err != nil {
		return reterror.ErrModel{Code: reterror.User_2012, Err: reterror.MSG_2012}
	}
	return reterror.ErrModel{Code: reterror.Success, Err: nil}
}

// ResetRecovery reset owner recovery
func ResetRecovery(account string) reterror.ErrModel {
	log.Debug("ResetRecovery...")
	if account == "" {
		return reterror.ErrModel{Code: reterror.ParamsNil, Err: reterror.PARAMS_NULL}
	}
	ds := db.RegDBService{}
	err := ds.ResetRecovery(account)
	if err != nil {
		log.Error("restr recovery error", err)
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	return reterror.ErrModel{Err: err, Code: reterror.Success}
}

// ActiveRecovery active new owner account
func ActiveRecovery(voucherVo *VoucherVo) reterror.ErrModel {
	log.Debug("ActiveRecovery...")
	if voucherVo.AppId == "" || voucherVo.AppName == "" || voucherVo.RegId == 0 {
		return reterror.ErrModel{Code: reterror.ParamsNil, Err: reterror.PARAMS_NULL}
	}
	ds := db.RegDBService{}
	as := db.AccDBService{}
	voucherStatus := voucher.GetVoucherStatus()
	applyer := voucherStatus.KeyStoreStatus[voucherVo.AppId].ApplyerId
	if applyer != voucherVo.AppId {
		// send data to voucher for recovery
		oper := &voucher.GrpcServer{
			Type:      voucher.VOUCHER_OPERATE_RECOVER,
			AppId:     voucherVo.AppId,
			AppName:   voucherVo.AppName,
			AesKey:    []byte(voucherVo.AesKey),
			Msg:       []byte(voucherVo.Msg),
			Sign:      []byte(voucherVo.Sign),
			Timestamp: voucherVo.Timestamp,
		}
		_, result, ret := voucher.SendVoucherData(oper)

		if ret == voucher.VRET_ERR || ret == voucher.VRET_TIMEOUT {
			log.Error("connect to voucher error", ret)
			return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
		}
		// check if account is frozen
		frozen, frozenTo, accErr := as.AccoutIsFrozen(voucherVo.AppName)
		if accErr == nil {
			if time.Now().Before(frozenTo) && frozen == 1 {
				message := "输入错误次数过多,账户被锁定到:" + frozenTo.Local().Format("2006-01-02 03:04:05 AM")
				return reterror.ErrModel{Code: reterror.User_3007, Err: errors.New(message), Data: map[string]interface{}{
					"data": frozenTo.Local().Format("2006-01-02 03:04:05 AM"),
				}}
			}
		} else {
			log.Error("数据库查询错误", accErr)
		}
		if result.Status == voucher.STATUS_APP_PASSWORD_ERROR {
			num, frozenTime := as.RecordAttempts(voucherVo.AppName)
			message := fmt.Sprintf("连续输错5次将锁定账户8小时,您已输错%v次", num)
			data := strconv.Itoa(num)
			code := reterror.User_3002
			if num == common.MAX_PASWD_ATTEMPTS {
				code = reterror.User_3007
				message = "输入错误次数过多,账户被锁定到:" + frozenTime.Local().Format("2006-01-02 03:04:05 PM")
				data = frozenTime.Local().Format("2006-01-02 03:04:05 PM")
			}
			return reterror.ErrModel{Code: code, Err: errors.New(message), Data: map[string]interface{}{"data": data}}
		}
		// request voucher failed
		if result.Status != voucher.STATUS_APP_RECOVERY_SUCCESSS {
			log.Error("签名机未全部通过审核", result.Status)
			return reterror.ErrModel{Code: reterror.Code_6, Err: reterror.MSG_6}
		}

		// get owner msg
		var otherMapInfo = make(map[string][]byte)
		jsonErr := json.Unmarshal(result.Other, &otherMapInfo)
		if jsonErr != nil {
			log.Error("签名机app公钥反解错误", jsonErr)
			return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
		}
		msg := string(otherMapInfo["voucherSign"][:])
		// add owner msg to db
		updateErr := ds.UpdateMsg(msg, voucherVo.RegId)
		if updateErr != nil {
			log.Error("更新股东签名失败", msg)
		}
	}

	// get owner account info
	adminInfo, err := as.AccountByName(voucherVo.AppName)
	if err != nil {
		log.Error("获取股东账号信息失败", err)
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}

	if len(adminInfo) == 0 {
		log.Error("股东不存在", voucherVo.AppName)
		return reterror.ErrModel{Code: reterror.User_2008, Err: reterror.MSG_2008}
	}
	//// 更新审批流
	//templateErr := template.DisableTemplateByUserName(&adminInfo[0], voucherVo.AppName)
	//if templateErr != nil {
	//	log.Error("作废审批流失败", templateErr)
	//	if err == reterror.MSG_9009 {
	//		return reterror.ErrModel{Code: reterror.Template_9009, Err: reterror.MSG_9009}
	//	}
	//	return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	//}
	da := db.AccDBService{}
	// 替换用户公钥
	if voucherVo.Replaced != "" {
		accErr := da.ReplaceMsg(voucherVo.Replaced, voucherVo.AppId)
		if accErr != nil {
			log.Error("替换msg失败", accErr)
			return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
		}
	}
	// 激活用户
	err = ds.ActiveRecovery(voucherVo.AppName, voucherVo.RegId)
	if err != nil {
		log.Error("激活失败", err)
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}

	// 更新审批流
	templateErr := template.DisableTemplateByUserName(&adminInfo[0], voucherVo.AppName)
	log.Debug("恢复股东作废审批流...")
	if templateErr != nil {
		log.Error("作废审批流失败", templateErr)
		if err == reterror.MSG_9009 {
			return reterror.ErrModel{Code: reterror.Template_9009, Err: reterror.MSG_9009}
		}
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}

	return reterror.ErrModel{Err: err, Code: reterror.Success}
}

// GetPubkeys 获取被恢复股东签名的用户公钥
func GetPubkeys(account string) reterror.ErrModel {
	log.Debug("GetPubkeys...")
	if account == "" {
		return reterror.ErrModel{Code: reterror.ParamsNil, Err: reterror.PARAMS_NULL}
	}
	ds := &db.AccDBService{}
	pubkeys, err := ds.GetPubkeys(account)
	list := make([]interface{}, len(pubkeys))
	for i, v := range pubkeys {
		list[i] = map[string]interface{}{
			"Pubkey": v.PubKey,
			"Name":   v.Name,
		}
	}
	if err != nil {
		log.Error("get recovered owners pubkey error", err)
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	return reterror.ErrModel{Err: err, Code: reterror.Success, Data: list}
}

// HasRecovery check if owner has recovery apply
func HasRecovery(account string) reterror.ErrModel {
	if account == "" {
		return reterror.ErrModel{Code: reterror.ParamsNil, Err: reterror.PARAMS_NULL}
	}
	ds := &db.RegDBService{}
	count, err := ds.HasRecovery(account)
	if err != nil {
		log.Error("check has recovery apply error", err)
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	return reterror.ErrModel{Err: err, Code: reterror.Success, Data: map[string]interface{}{"Count": count}}
}

// GetVersion get api version
func GetVersion() reterror.ErrModel {
	treeVersion := db.ConfigMap["tree_version"]
	maps := versionlog.LogMap
	//cfg, _ := config.LoadConfig()
	maps["TreeVersion"] = treeVersion
	maps["DownLoadUrl"] = config.Conf.Server.AppPath
	return reterror.ErrModel{Err: nil, Code: reterror.Success, Data: maps}
}

// GetBlockHeight get block height
func GetBlockHeight() reterror.ErrModel {
	wallet := walletCli.NewAppServer()
	result := wallet.GetHeights()
	coinType := map[bccore.BloclChainType]bccore.BlockChainSign{bccore.BC_BTC: bccore.STR_BTC, bccore.BC_ETH: bccore.STR_ETH, bccore.BC_LTC: bccore.STR_LTC}
	cfg := config.Conf
	heightDiff := map[bccore.BloclChainType]uint64{bccore.BC_BTC: cfg.Server.BTCHeightDiff, bccore.BC_ETH: cfg.Server.ETHHeightDiff, bccore.BC_LTC: cfg.Server.LTCHeightDiff}
	// returned data
	var list []interface{}
	aggregatedStatus := 0
	for i, v := range result {
		status := common.NodeStable
		if (v.PubHeight - v.CurHeight) > heightDiff[i] {
			status = common.NodeSyncing
			aggregatedStatus = common.NodeSyncing
		}
		list = append(list, map[string]interface{}{
			"PubHeight": v.PubHeight,
			"CurHeight": v.CurHeight,
			"Name":      coinType[i],
			"Status":    status,
		})
	}
	return reterror.ErrModel{Err: nil, Code: reterror.Success, Data: map[string]interface{}{"List": list, "AgrStatus": aggregatedStatus}}
}
