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
	"strconv"
	log "github.com/alecthomas/log4go"
	"github.com/jinzhu/gorm"
	rerr "github.com/boxproject/apiServer/errors"
)

type DepService struct {
}

// DepartmentNameExists check if the department exists by given name and its id
func (*DepService) DepartmentNameExists(departmentName string, id int) (bool, error) {
	log.Debug("verify duplicate department name", departmentName)
	var dep Department
	var err error
	if id == 0 {
		err = rdb.Where("name = ? and available = 0", departmentName).First(&dep).Error
	} else {
		err = rdb.First(&dep, "id != ? and name = ? and available = 0", id, departmentName).Error
	}

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, rerr.DEP_EXISTS
}

// AddDepartment add new department
func (d *DepService) AddDepartment(depName string, creatorID int) error {
	order := d.GetDepCount() + 1
	dep := Department{
		Name:         depName,
		CreatorAccId: creatorID,
		Order:        order,
	}
	err := rdb.Where("name=? and available=0", depName).FirstOrCreate(&dep).Error
	return err
}

// EditDepartment edit a existed department
func (*DepService) EditDepartment(dep *Department) error {
	//err := d.DB.Where("id=?",dep.ID).Model(&Department{}).Update("name",dep.Name).Error
	err := rdb.Model(dep).Update(dep).Error
	return err
}

// DelDepartment delete a existed department
func (*DepService) DelDepartment(id int) error {
	var dep []Department
	err := rdb.Model(&dep).Where("id=? and available=0", id).Update("available", 1).Error
	return err
}

// DepartmentList get the existed department list
func (*DepService) DepartmentList() ([]Department, error) {
	var deps []Department
	err := rdb.Order("`order` asc").Where("available = 0").Find(&deps).Error
	return deps, err
}

// GetAccountByDepartment get all the members belong to the given department
func (*DepService) GetAccountByDepartment(id int, authID int) ([]Account, []AuthMap, error) {
	var account []Account
	var authMap []AuthMap
	err := rdb.Where("departmentId=? and isDeleted = 0", id).Find(&account).Error
	if authID != 0 {
		rdb.Where(&AuthMap{AuthId: authID}).Find(&authMap)
	}
	return account, authMap, err
}

// GetDepByID get the department info by given account
func (*DepService) GetDepByID(id int) ([]Department, error) {
	var dep []Department
	err := rdb.Where("id=?", id).Find(&dep).Error
	return dep, err
}

// GetDepAccCount count the member of given department
func (*DepService) GetDepAccCount(depID int) int {
	var count int
	rdb.Model(&Account{}).Where("departmentId=? and isDeleted=0", depID).Count(&count)
	return count
}

// GetDepCount count the existed department
func (*DepService) GetDepCount() int {
	var count int
	rdb.Model(&Department{}).Where("available = 0").Count(&count)
	return count
}

// SortDep sort the department order by app
func (d *DepService) SortDep(depID []string) error {
	if len(depID) != 0 {
		// 校验排序部门是否有效
		var depCount int
		rdb.Model(&Department{}).Where("id in (?) and available=0", depID).Count(&depCount)
		available := d.GetDepCount()
		if depCount != available {
			return gorm.ErrRecordNotFound
		}
		// 排序
		query := ""
		for i, v := range depID {
			query += " WHEN " + v + " THEN " + strconv.Itoa(i+1)
		}
		sql := "CASE id" + query + " END"
		err := rdb.Table("department").Where("id IN (?) and available=0", depID).Updates(map[string]interface{}{"order": gorm.Expr(sql)}).Error
		if err != nil {
			return err
		}
	}
	return nil
}
