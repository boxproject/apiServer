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
package db

import (
	"fmt"
	"encoding/json"
	"strconv"
	"strings"
	"time"
	log "github.com/alecthomas/log4go"
	"github.com/jinzhu/gorm"
	"github.com/boxproject/apiServer/common"
)

type AccDBService struct {
}

type AsyncBlockChain struct {
	List             []interface{}
	AggregatedStatus int
}

var AsyncBlockChains = &AsyncBlockChain{}

// AddAccount add account info to db
func (*AccDBService) AddAccount(account *Account) error {
	err := rdb.Create(account).Error
	return err
}

// AccoutIsFrozen check account is forzen
func (*AccDBService) AccoutIsFrozen(name string) (int, time.Time, error) {
	var account Account
	err := rdb.First(&account, "name = ?", name).Error
	fmt.Println(account.FrozenTo)
	return account.Frozen, account.FrozenTo, err
}

// RecordAttempts record incorrect password attempts
func (*AccDBService) RecordAttempts(name string) (result int, frozentime time.Time) {
	var account Account
	err := rdb.First(&account, "name = ?", name).Error
	if err != nil {
		result = 0
		return
	}
	now := time.Now()
	num := account.Attempts
	if account.Frozen == 1 && now.After(account.FrozenTo) {
		num = 1
		number := rdb.Model(&account).Where("name = ?", name).Updates(map[string]interface{}{"attempts": num, "frozen": 0}).RowsAffected
		if number == 1 { //小于5
			result = num
			return
		}
		result = 0
		return
	}
	if num == 4 {
		frozenTo := now.Add(time.Duration(8) * time.Hour)
		num++
		number := rdb.Model(&account).Where("name = ?", name).Updates(map[string]interface{}{"attempts": num, "frozen": 1, "frozenTo": frozenTo}).RowsAffected
		if number == 1 {
			result = 5
			frozentime = frozenTo
			return
		}
		result = 0
		return
	}
	if num == 5 {
		return 5, account.FrozenTo
	}
	num++
	number := rdb.Model(&account).Where("name = ?", name).Update("attempts", num).RowsAffected
	if number == 1 {
		result = num
		return
	}
	result = 0
	return
}

// ExistAccount check if account is existed and return salt
func (*AccDBService) ExistAccount(name string) (*Account, error) {
	var account Account
	err := rdb.First(&account, "name = ?", name).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			//不存在账号
			return nil, nil
		}
		//查询错误
		return nil, err
	}
	//存在,将盐取出
	return &account, nil
}

// AccountByName find account by account name
func (*AccDBService) AccountByName(name string) ([]Account, error) {
	var account []Account
	err := rdb.Where("name=?", name).Find(&account).Error
	if err != nil {
		return account, err
	}
	return account, nil
}

// AccountByID find account by id
func (*AccDBService) AccountByID(acc_id int) (*Account, error) {
	var account Account
	err := rdb.Where(&Account{ID: acc_id}).First(&account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Debug("用户未找到，accID = ", acc_id)
			return nil, nil
		}
		log.Debug("获取用户账号失败", err)
		return nil, err
	}
	return &account, nil
}

// Login check name and password and login
func (*AccDBService) Login(name string, pwd string) (bool, *Account, error) {
	var account Account
	err := rdb.First(&account, "name = ? and pwd = ?", name, pwd).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, &account, nil
		}
		return false, &account, err
	}
	return true, &account, nil
}

//GetSalt get salt password
func (*AccDBService) GetSalt(name string) (string, string, error) {
	var account Account
	err := rdb.First(&account, "name = ?", name).Error
	return account.Salt, account.Pwd, err
}

// ModifyPassword modify password
func (*AccDBService) ModifyPassword(name, new_pwd string) error {
	var account Account
	err := rdb.Model(&account).Where("name = ?", name).Update("pwd", new_pwd).Error
	return err
}

// ResetAttempts resume frozen account by name
func (*AccDBService) ResetAttempts(name string) error {
	var account Account
	err := rdb.Model(&account).Where("name = ?", name).Updates(map[string]interface{}{"attempts": 0, "frozen": 0}).Error
	return err
}

// DelayTaskNum count todos（other account to be verified）
func (*AccDBService) DelayTaskNum() (task DelayTask) {
	var total int
	var name string
	row := rdb.Table("registration").Select("name").Where("status = 1").Order("createdAt desc").Limit(1).Row()
	row.Scan(&name)
	row = rdb.Table("registration").Select("count(0) as total").Where("status = 1").Row()
	row.Scan(&total)
	task.Number = total
	task.Reason = name
	return
}

// ResetAccount resume frozen account by app id
func (*AccDBService) ResetAccount(appid string) (bool, error) {
	err := rdb.Model(&Account{}).Where(&Account{AppId: appid}).Updates(map[string]interface{}{"frozen": 0, "attempts": 0}).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

// VerifyUser verify account and add user info to table account
func (*AccDBService) VerifyUser(account *Account, reg *Registration) bool {
	tx := rdb.Begin()
	if reg.Status == 2 {
		if err := tx.Model(reg).Update(reg).Error; err != nil {
			tx.Rollback()
			log.Error("审核错误", err)
			return false
		}
		if err := tx.Create(account).Error; err != nil {
			tx.Rollback()
			log.Error("审核错误", err)
			return false
		}

	} else {
		err := tx.Model(reg).Where("id=?", reg.ID).Updates(map[string]interface{}{"status": reg.Status, "refuseReason": reg.RefuseReason}).Error
		if err != nil {
			tx.Rollback()
			log.Error("审核错误", err)
			return false
		}
	}
	tx.Commit()
	return true
}

// FindAccountByAppId find account by app id
func (*AccDBService) FindAccountByAppId(appId string) (*Account, error) {
	var acc Account
	err := rdb.First(&acc, "appId = ?", appId).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &acc, nil

}

// FindAllUsers get all account and order by level
func (*AccDBService) FindAllUsers() ([]Account, error) {
	var accounts []Account
	err := rdb.Order("level desc").Find(&accounts).Error
	return accounts, err
}

// GetAccountsByType get account info by user type  0: member, 1: admin, 2: owner
func (*AccDBService) GetAccountsByType(userType int) ([]Account, error) {
	var accounts []Account
	err := rdb.Where("userType=? and isDeleted=0", userType).Find(&accounts).Error
	return accounts, err
}

// GetAdminCount get admin count
func (*AccDBService) GetAdminCount() (count int, err error) {
	err = rdb.Table("account").Where("userType=1 and isDeleted=0").Count(&count).Error
	return count, err
}

// AddAdmin set account as admin
func (*AccDBService) AddAdmin(accountID int) error {
	var account []Account
	var authMap []AuthMap
	// update user type
	conn := rdb.Begin()
	accountErr := conn.Model(&account).Where("id=? and userType = 0", accountID).Update("userType", 1).Error
	if accountErr != nil {
		conn.Rollback()
		return accountErr
	}
	// remove authmap data
	err := conn.Model(&authMap).Where("accountId=?", accountID).Delete(&AuthMap{}).Error
	if err != nil {
		conn.Rollback()
		return err
	}
	conn.Commit()
	return err
}

// DelAdmin remove admin
func (*AccDBService) DelAdmin(accountID int) error {
	var account []Account
	err := rdb.Model(&account).Where("id=? and userType = 1", accountID).Update("userType", 0).Error
	return err
}

// GetAllUsers get all available accounts
func (*AccDBService) GetAllUsers() ([]Account, error) {
	var account []Account
	err := rdb.Where("isDeleted = 0").Find(&account).Error
	return account, err
}

// GetUserByID get user info by id
func (*AccDBService) GetUserByID(id int) ([]Account, error) {
	var account []Account
	err := rdb.Where("id=?", id).Find(&account).Error
	return account, err
}

// SavePubkeySign update msg(sign)
func (*AccDBService) SavePubkeySign(account *Account) error {
	err := rdb.Model(&account).Where("appId=?", account.AppId).Update("msg", account.Msg).Error
	return err
}

// SetUser update user info (auth, user type, department)
func (*AccDBService) SetUser(authID string, depID int, accID int, userType int) error {
	var account []Account
	var authMap []AuthMap
	var auth []Auth
	tx := rdb.Begin()
	accErr := tx.Model(&account).Where("id=?", accID).Find(&account).Update("userType", userType).Error
	if accErr != nil {
		tx.Rollback()
		return accErr
	}
	if len(account) == 0 {
		tx.Rollback()
		return nil
	}
	if err := tx.Model(&account).Where("id=?", accID).Update("departmentId", depID).Error; err != nil {
		tx.Rollback()
		return err
	}
	if authID != "" {
		authIds := strings.Split(authID, ",")
		parsedIds := make([]int, len(authIds))
		for i, v := range authIds {
			parsedIds[i], _ = strconv.Atoi(v)
		}
		delErr := tx.Model(&authMap).Where("accountId = ?", accID).Delete(&AuthMap{}).Error
		if delErr != nil {
			tx.Rollback()
			return delErr
		}
		for _, v := range parsedIds {
			authExist := tx.Model(&auth).Where("id=?", v).Find(&auth).RowsAffected
			if authExist == 0 {
				continue
			}
			addErr := tx.Model(&authMap).Create(&AuthMap{AccountId: accID, AuthId: v}).Error
			if addErr != nil {
				tx.Rollback()
				return addErr
			}
		}
	} else {
		delErr := tx.Model(&authMap).Where("accountId = ?", accID).Delete(&AuthMap{}).Error
		if delErr != nil {
			tx.Rollback()
			return delErr
		}
	}
	if userType != 0 {
		authErr := tx.Model(&authMap).Where("accountId=?", accID).Delete(&AuthMap{}).Error
		if authErr != nil {
			tx.Rollback()
			return authErr
		}
	}

	tx.Commit()
	return nil
}

// DisableAcc disable account by id
func (*AccDBService) DisableAcc(accID int) error {
	// var reg []Registration
	var authMap []AuthMap
	tx := rdb.Begin()
	err := tx.Model(&Account{ID: accID, IsDeleted: 0}).Where("userType != 2").Update("isDeleted", 1).Error
	if err != nil {
		log.Debug("删除用户账号", err)
		tx.Rollback()
		return err
	}

	authMapErr := tx.Model(&authMap).Where("accountId=?", accID).Delete(&AuthMap{}).Error
	if authMapErr != nil {
		tx.Rollback()
		return authMapErr
	}
	tx.Commit()
	return nil
}

// IsOwner check if user is owner by name
func (*AccDBService) IsOwner(name string) (bool, error) {
	var account []Account
	err := rdb.Where("name=? and userType=?", name, common.OwnerAccType).Find(&account).Error
	if err != nil {
		return false, err
	}
	if len(account) == 0 {
		return false, err
	}
	return true, err
}

// GetPubkeys get public keys by name
func (*AccDBService) GetPubkeys(name string) ([]Account, error) {
	var account []Account
	err := rdb.Joins("INNER JOIN registration on registration.name = account.name").Where("sourceAccount=?", name).Find(&account).Error
	return account, err
}

// Content represent content
type Content struct {
	Name string `json:"name"`
	Msg  string `json:"msg"`
}

// Replaced include
type Replaced struct {
	Data []Content `json:"data"`
}

// ReplaceMsg 替换 msg
func (*AccDBService) ReplaceMsg(rep string, appId string) error {
	tmp := Replaced{}
	jsonErr := json.Unmarshal([]byte(rep), &tmp)
	if jsonErr != nil { // TODO 错误码
		return jsonErr
	}
	var account []Account
	replacedArr := tmp.Data
	conn := rdb.Begin()
	for i := 0; i < len(replacedArr); i++ {
		err := conn.Model(&account).Where("name=?", replacedArr[i].Name).Updates(Account{Msg: replacedArr[i].Msg, SourceAppId: appId}).Error
		if err != nil {
			conn.Rollback()
			return err
		}
	}
	var config []Configs
	values := ConfigMap["tree_version"]
	val, _ := strconv.Atoi(values)
	val = val + 1
	updateErr := conn.Model(&config).Where("con_key=?", "tree_version").Update("con_value", strconv.Itoa(val)).Error
	if updateErr != nil {
		conn.Rollback()
		return updateErr
	}
	ConfigMap["tree_version"] = strconv.Itoa(val)
	conn.Commit()
	return nil
}
