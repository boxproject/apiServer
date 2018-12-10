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
package department

import (
	"strings"
	log "github.com/alecthomas/log4go"
	"github.com/boxproject/apiServer/db"
	reterror "github.com/boxproject/apiServer/errors"
)

// AddDepartment add department
func AddDepartment(department *Pdepartment) reterror.ErrModel {
	log.Debug("AddDepartment...")
	if department.Name == "" {
		return reterror.ErrModel{Code: reterror.ParamsNil, Err: reterror.PARAMS_NULL}
	}

	ds := &db.DepService{}
	// check if department exists
	department_exists, err := ds.DepartmentNameExists(department.Name, 0)
	if err != nil {
		err := reterror.ErrModel{Err: err}
		if department_exists {
			err.Code = reterror.Duplicate
		} else {
			err.Code = reterror.Failed
		}
		return err
	}
	// add department
	err = ds.AddDepartment(department.Name, department.CreatorID)

	if err != nil {
		log.Error("Add department failed case:", err)
		return reterror.ErrModel{Err: reterror.MSG_1003, Code: reterror.AddDepFailed}
	}

	return reterror.ErrModel{Err: nil, Code: reterror.Success}
}

// EditDepartment edit department
func EditDepartment(department *Pdepartment) reterror.ErrModel {
	log.Debug("EditDepartment...")
	if department.Name == "" || (department.ID == 0) {
		return reterror.ErrModel{Code: reterror.ParamsNil, Err: reterror.PARAMS_NULL}
	}
	ds := &db.DepService{}
	// check if department exists
	department_exists, err := ds.DepartmentNameExists(department.Name, department.ID)
	if err != nil {
		log.Error("edit department error", err)
		err := reterror.ErrModel{Err: err}
		if department_exists {
			err.Code = reterror.Duplicate
		} else {
			err.Code = reterror.Failed
		}
		return err
	}
	//var dep db.Department
	dep := &db.Department{
		ID:   department.ID,
		Name: department.Name,
	}
	err = ds.EditDepartment(dep)

	if err != nil {
		log.Error("edit department failed case:", err)
		return reterror.ErrModel{Err: reterror.MSG_1004, Code: reterror.EditDepFailed}
	}
	return reterror.ErrModel{Err: nil, Code: reterror.Success}
}

// DepartmentList get department list
func DepartmentList() reterror.ErrModel {
	ds := &db.DepService{}
	list, _ := ds.DepartmentList()
	departments := make([]interface{}, len(list))
	for i, v := range list {
		departments[i] = map[string]interface{}{
			"ID":    v.ID,
			"Name":  v.Name,
			"Count": ds.GetDepAccCount(v.ID),
		}
	}
	return reterror.ErrModel{Err: nil, Code: reterror.Success, Data: departments}
}

// DelDepartment delete department
func DelDepartment(department *Pdepartment) reterror.ErrModel {
	log.Debug("DelDepartment...")
	if department.ID == 0 || department.ID == 1 {
		return reterror.ErrModel{Code: reterror.ParamsNil, Err: reterror.PARAMS_NULL}
	}
	ds := &db.DepService{}
	count := ds.GetDepAccCount(department.ID)
	if count != 0 {
		return reterror.ErrModel{Code: reterror.DelDepFailed, Err: reterror.MSG_1008}
	}
	err := ds.DelDepartment(department.ID)
	if err != nil {
		log.Error("delete department failed case:", err)
		return reterror.ErrModel{Err: reterror.MSG_1004, Code: reterror.EditDepFailed}
	}
	return reterror.ErrModel{Err: nil, Code: reterror.Success}
}

// GetAccountsByDepartment get accounts by department id or auth id
func GetAccountsByDepartment(depID int, authID int) reterror.ErrModel {
	if depID == 0 {
		return reterror.ErrModel{Code: reterror.ParamsNil, Err: reterror.PARAMS_NULL}
	}
	ds := &db.DepService{}
	list, authMap, err := ds.GetAccountByDepartment(depID, authID)
	authAccount := make([]interface{}, len(list))
	for i, ac := range list {
		authAccount[i] = map[string]interface{}{
			"ID":       ac.ID,
			"Name":     ac.Name,
			"HasAuth":  false,
			"UserType": ac.UserType,
			"AppId":    ac.AppId,
		}
		for _, auth := range authMap {
			if auth.AccountId == ac.ID {
				assertion := authAccount[i].(map[string]interface{})
				assertion["HasAuth"] = true
			}
		}
	}
	if err != nil {
		log.Error("GetAccountsByDepartment error", err)
		return reterror.ErrModel{Err: reterror.MSG_5006, Code: reterror.Auth_5006}
	}
	return reterror.ErrModel{Err: nil, Code: reterror.Success, Data: authAccount}
}

// SortDep sort department
func SortDep(depID string) reterror.ErrModel {
	log.Debug("SortDep...")
	parsedDepID := strings.Split(depID, ",")
	ds := &db.DepService{}
	count := ds.GetDepCount()
	if len(parsedDepID) != count {
		return reterror.ErrModel{Err: reterror.MSG_2013, Code: reterror.User_2013}
	}
	err := ds.SortDep(parsedDepID)
	if err != nil {
		log.Error("sort department order error", err)
		return reterror.ErrModel{Err: reterror.MSG_1005, Code: reterror.SetDepFailed}
	}
	return reterror.ErrModel{Err: nil, Code: reterror.Success}
}
