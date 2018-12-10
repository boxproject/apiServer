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

import "github.com/jinzhu/gorm"

// AuthService db
type AuthService struct {
}

// GetAuthList 获取权限列表
func (*AuthService) GetAuthList() ([]Auth, error) {
	var authList []Auth
	err := rdb.Find(&authList).Error
	return authList, err
}

// GetAuthAccounts 获取单个权限下成员列表
func (*AuthService) GetAuthAccounts(id int) ([]Account, error) {
	var authMap []AuthMap
	var account []Account
	// 获取指定权限的关联表内容
	var accountIds []int

	authErr := rdb.Where(&AuthMap{AuthId: id}).Find(&authMap).Error
	if authErr != nil {
		return account, authErr
	}
	for _, v := range authMap {
		accountIds = append(accountIds, v.AccountId)
	}
	err := rdb.Where("id in (?) and userType=0", accountIds).Find(&account).Error
	return account, err
}

// GetAuthAccountsCount 获取指定权限下成员数
func (*AuthService) GetAuthAccountsCount(id int) int {
	var count int
	rdb.Model(&AuthMap{}).Where(&AuthMap{AuthId: id}).Count(&count)
	return count
}

// AddAuthToAccount 账户添加权限
func (*AuthService) AddAuthToAccount(add []int, id int) int {
	var authMap []AuthMap
	var account []Account
	var auth []Auth
	// 判断权限是否存在
	authMapRow := rdb.First(&auth, id)
	if authMapRow.RowsAffected == 0 {
		return 2
	}

	for _, v := range add {
		// 判断账户是否存在
		accountRow := rdb.First(&account, v)
		if accountRow.RowsAffected != 1 {
			return 3
		}
		// 判断是否已设置权限
		authMapRow := rdb.Where(&AuthMap{AccountId: v, AuthId: id}).First(&authMap)
		if authMapRow.RowsAffected == 1 {
			continue
		}
		authmap := AuthMap{
			AccountId: v,
			AuthId:    id,
		}
		err := rdb.Create(&authmap).Error
		if err != nil {
			return 4
		}
	}
	return 1
}

// DelAuthFromAccount 移除账户权限
func (*AuthService) DelAuthFromAccount(authID int, accountID int) error {
	authMapErr := rdb.Where(&AuthMap{AuthId: authID, AccountId: accountID}).Delete(&AuthMap{}).Error
	return authMapErr
}

// AuthCheck 权限检查
func (*AuthService) AuthCheck(authType string, accountID int) error {
	var authMap []AuthMap
	var account []Account
	authtype := authType + "_auth"
	accErr := rdb.Where(&Account{ID: accountID}).First(&account).Error
	if accErr != nil {
		return accErr
	}
	if len(account) > 0 && account[0].UserType > 0 {
		return nil
	}
	err := rdb.Joins("INNER JOIN auth on auth.id = authmap.authId and authmap.accountId=?", accountID).Where("authType=?", authtype).Find(&authMap).Error
	if err != nil {
		return err
	}
	if len(authMap) == 0 {
		return gorm.ErrRecordNotFound
	}
	return err
}

// GetAuthByID 获取单个权限
func (*AuthService) GetAuthByID(authID int) ([]Auth, error) {
	var auth []Auth
	err := rdb.Where(&Auth{ID: authID}).Find(&auth).Error
	return auth, err
}

// GetAuthMapByAccID 获取用户所有权限
func (*AuthService) GetAuthMapByAccID(accountID int) ([]AuthMap, error) {
	var authMap []AuthMap
	err := rdb.Where(&AuthMap{AccountId: accountID}).Find(&authMap).Error
	return authMap, err
}
