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
	"os"
	"github.com/jinzhu/gorm"
	rerr "github.com/boxproject/apiServer/errors"
	"github.com/theplant/batchputs"
	"github.com/boxproject/apiServer/common"
)

type RegDBService struct {
}

// NameExists check if account existed by name
func (*RegDBService) NameExists(name string) (bool, error) {
	var reg []Registration
	err := rdb.Raw("select * from ((select name from account) union (select name from registration)) user where name = ?", name).Find(&reg).Error
	if err != nil {
		return false, err
	}
	if len(reg) != 0 {
		return true, rerr.MSG_2002
	}
	return false, err
}

// AppIdExists check if account existed by appid
func (*RegDBService) AppIdExists(appId string) (bool, error) {
	var reg Registration
	err := rdb.First(&reg, "appId = ? and status !=3", appId).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, rerr.DEP_EXISTS
}

// AddUser add new user account
func (*RegDBService) AddUser(reg *Registration) error {
	err := rdb.Create(&reg).Error
	return err
}

// EditUser edit a user account
func (*RegDBService) EditUser(reg *Registration) error {
	err := rdb.Model(reg).Update(reg).Error
	return err
}

// Sign 登陆
func (*RegDBService) Sign(reg *Registration) error {
	err := rdb.First(&reg, "appId = ? and Pwd = ?", reg.AppId, reg.Pwd).Error
	return err
}

// FindUserByAppId 通过appid查询注册状态
func (*RegDBService) FindUserByAppId(appId string) (*Registration, error) {
	var reg Registration
	err := rdb.First(&reg, "appId = ?", appId).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &reg, nil
}

// RegList 查询注册列表
func (*RegDBService) RegList(reg *Registration) ([]Registration, error) {
	var regs []Registration
	err := rdb.Where(&regs).Find(&regs).Error
	return regs, err
}

// ReSignUp 重新注册用户（删除注册表中数据）
func (*RegDBService) ReSignUp(appID string) error {
	err := rdb.Where("appId = ? and status = 0", appID).Or("appId = ? and status = 3", appID).Delete(&Registration{}).Error
	return err
}

// GetRegList 获取扫码后注册用户列表
func (*RegDBService) GetRegList() ([]Registration, error) {
	var reg []Registration
	err := rdb.Where("status=?", common.RegToApproval).Order("updatedAt desc").Find(&reg).Error
	return reg, err
}

// AddReg 添加恢复股东注册信息
func (*RegDBService) AddReg(ownerReg *OwnerReg) (int, error) {
	err := rdb.Create(&ownerReg).Error
	return ownerReg.ID, err
}

// RecoveryList 获取待恢复股东列表
func (*RegDBService) RecoveryList() ([]OwnerReg, error) {
	var reg []OwnerReg
	err := rdb.Where("status != 0").Order("createdAt desc").Find(&reg).Error
	return reg, err
}

// CheckRecovStatus 查看恢复账号审批状态
func (*RegDBService) CheckRecovStatus(appID string) (int, error) {
	var ownerReg []OwnerReg
	err := rdb.Where("appId=?", appID).Find(&ownerReg).Error
	if err != nil {
		return common.OwnerRegError, err
	}
	if len(ownerReg) == 0 {
		return common.OwnerRegNotFound, nil
	}
	for _, v := range ownerReg {
		if v.Status == common.RegWaitingToApproval {
			return common.RegWaitingToApproval, nil
		}
	}
	return ownerReg[0].Status, err
}

// SubRecovery 提交确认信息 （关系表初始化数据）
func (*RegDBService) SubRecovery(appID string, acc string, regID int) error {
	var account []Account
	// 获取所有股东
	accountErr := rdb.Where("userType=? and name != ?", 2, acc).Find(&account).Error
	if accountErr != nil || len(account) == 0 {
		return accountErr
	}
	// 获取操作表最后一条的id?
	var owner []OwnerOper
	var count int
	ownerErr := rdb.Select("id").Last(&owner).Error
	if ownerErr != nil {
		return ownerErr
	}
	if len(owner) == 0 {
		count = 0
	} else {
		count = owner[0].ID
	}

	var ownerOper OwnerOper
	var ownerReg []OwnerReg
	conn := rdb.Begin()
	// 修改注册表状态
	regErr := conn.Model(&ownerReg).Where("appId = ? and name=? and status = 0", appID, acc).Find(&ownerReg).Update("status", 1).Error
	if regErr != nil || len(ownerReg) == 0 {
		conn.Rollback()
		return regErr
	}
	// 构造插入数据
	var vals [][]interface{}
	for _, v := range account {
		if acc != v.Name {
			count++
			vals = append(vals, []interface{}{count, v.Name, acc, 0, appID, ownerReg[0].ID})
		}
	}

	regColumns := []string{"id", "account", "operatedAccount", "status", "operatedAppid", "regId"}
	err := batchputs.Put(rdb.DB(), os.Getenv("DB_DIALECT"), ownerOper.TableName(), "id", regColumns, vals)
	if err != nil {
		conn.Rollback()
		return err
	}
	// 作废之前的待提交申请
	subErr := conn.Model(&ownerReg).Where("name=? and status=1 and id !=?", acc, regID).Update("status", 5).Error
	if subErr != nil {
		conn.Rollback()
		return subErr
	}
	conn.Commit()
	return err
}

// RecoveryResult 获取恢复股东认证结果
func (*RegDBService) RecoveryResult(regID int, appID string) ([]OwnerOper, error) {
	var ownerOper []OwnerOper
	err := rdb.Where("regID=? and operatedAppid=?", regID, appID).Find(&ownerOper).Error
	return ownerOper, err
}

// RegStatus 审批状态
func (*RegDBService) RegStatus(regID int) (int, string, error) {
	var ownerReg []OwnerReg
	err := rdb.Where("id=?", regID).Find(&ownerReg).Error
	if err != nil || len(ownerReg) == 0 {
		return -1, "", err
	}
	return ownerReg[0].Status, ownerReg[0].Name, err
}

// CancelRecovery 取消之前的申请
func (*RegDBService) CancelRecovery(account string) error {
	var ownerReg []OwnerReg
	err := rdb.Model(&ownerReg).Where("account=? and (status=0 or status=1)", account).Update("status", 5).Error
	return err
}

// VerifyRecovery 认证股东
func (*RegDBService) VerifyRecovery(account string, status int, regID int, voucherStatus int) error {
	var ownerOper []OwnerOper
	var ownerReg []OwnerReg
	conn := rdb.Begin()
	err := conn.Model(&ownerOper).Where("account=? and regId=? and status = 0", account, regID).Find(&ownerOper).Update("status", status).Error
	if err != nil {
		conn.Rollback()
		return err
	}
	if len(ownerOper) == 0 {
		conn.Rollback()
		return gorm.ErrRecordNotFound
	}
	// 拒绝时 修改注册表状态
	if status == common.OwnerAccRejected {
		regErr := conn.Model(&ownerReg).Where("id=?", regID).Update("status", 3).Error
		if regErr != nil {
			conn.Rollback()
			return regErr
		}
	}
	if status == common.OwnerAccApprovaled && voucherStatus == common.OwnerAccApprovaled {
		regErr := conn.Model(&ownerReg).Where("id=?", regID).Update("status", common.RegApprovaled).Error
		if regErr != nil {
			conn.Rollback()
			return regErr
		}
	}
	conn.Commit()
	return err
}

// ResetRecovery 重新注册恢复股东账户（状态为0的删除数据，其余不修改）
func (*RegDBService) ResetRecovery(account string) error {
	err := rdb.Where("name=? and status = 0", account).Delete(&OwnerReg{}).Error
	return err
}

// OperateStatus 当前股东审批状态
func (*RegDBService) OperateStatus(account string, regID int) (int, error) {
	var ownerOper []OwnerOper
	err := rdb.Where("account=? and regId=?", account, regID).Find(&ownerOper).Error
	if err != nil {
		return common.OwnerRegError, err
	}
	if len(ownerOper) == 0 {
		return common.OwnerRegNotFound, nil
	}
	return ownerOper[0].Status, err
}

// ActiveRecovery 激活股东
func (*RegDBService) ActiveRecovery(name string, regID int) error {
	var ownerReg []OwnerReg
	var account []Account
	err := rdb.Where("name=? and status=? and id=?", name, common.RegApprovaled, regID).Find(&ownerReg).Error
	if err != nil {
		return err
	}
	if len(ownerReg) == 0 {
		return gorm.ErrRecordNotFound
	}
	conn := rdb.Begin()
	accErr := conn.Model(&account).Where("name=?", name).Update(&Account{Pwd: ownerReg[0].Pwd, PubKey: ownerReg[0].PubKey, Salt: ownerReg[0].Salt, AppId: ownerReg[0].AppId, Msg: ownerReg[0].Msg}).Error
	if accErr != nil {
		conn.Rollback()
		return accErr
	}
	conn.Commit()
	return accErr
}

// UpdateRecoveryStatus 更新恢复股东申请状态
func (*RegDBService) UpdateRecoveryStatus(regID int) error {
	var ownerReg []OwnerReg
	conn := rdb.Begin()
	err := conn.Model(&ownerReg).Where("id=?", regID).Update("status", common.OwnerRegNotFound).Error
	if err != nil {
		conn.Rollback()
		return err
	}
	conn.Commit()
	return err
}

// HasRecovery 是否有进行中的恢复申请
func (*RegDBService) HasRecovery(account string) (int, error) {
	var ownerReg []OwnerReg
	var count int
	err := rdb.Where("name=? and status in (0, 1)", account).Find(&ownerReg).Count(&count).Error
	return count, err
}

// Duplicate 判断注册用户名是否重复
func (*RegDBService) DuplicateName(account string) (bool, error) {
	var reg []Registration
	regErr := rdb.Where("name=?", account).Find(&reg).Error
	if regErr != nil {
		return false, regErr
	}
	if len(reg) != 0 {
		return true, regErr
	}
	return false, regErr
}

// UpdateMsg update owner message
func (*RegDBService) UpdateMsg(msg string, regID int) error {
	err := rdb.Table("ownerReg").Where("id=?", regID).Update("msg", msg).Error
	return err
}
