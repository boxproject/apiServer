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
package controllers

import (
	log "github.com/alecthomas/log4go"
	"github.com/boxproject/apiServer/service/department"
	"github.com/boxproject/apiServer/errors"
	"github.com/gin-gonic/gin"
)

// AddDepartment used to add a new department
func AddDepartment(ctx *gin.Context) {
	var dep department.Pdepartment
	var msg errors.ErrModel
	if err := ctx.ShouldBind(&dep); err != nil {
		log.Error("Bind Param Error", err)
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	} else {
		msg = department.AddDepartment(&dep)
	}

	msg.RetErr(ctx)
}

// DepartmentList used to get the list of the departments
func DepartmentList(ctx *gin.Context) {
	var msg errors.ErrModel
	msg = department.DepartmentList()
	msg.RetErr(ctx)
}

// EditDepartment used to edit the department info
func EditDepartment(ctx *gin.Context) {
	var dep department.Pdepartment
	var msg errors.ErrModel
	if err := ctx.ShouldBind(&dep); err != nil {
		log.Error("Bind Param Error", err)
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	} else {
		msg = department.EditDepartment(&dep)
	}
	msg.RetErr(ctx)

}

// DelDepartment used to delete a department
func DelDepartment(ctx *gin.Context) {
	var dep department.Pdepartment
	var msg errors.ErrModel
	if err := ctx.ShouldBind(&dep); err != nil {
		log.Error("Bind Param Error", err)
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	} else {
		msg = department.DelDepartment(&dep)
	}
	msg.RetErr(ctx)
}

// DepartmentAccounts used to get the member of one department
func DepartmentAccounts(ctx *gin.Context) {
	var msg errors.ErrModel
	var dep department.Pdepartment
	if err := ctx.ShouldBind(&dep); err != nil {
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	} else {
		msg = department.GetAccountsByDepartment(dep.ID, dep.AuthID)
	}
	msg.RetErr(ctx)
}

// SortDep used to sort all the departments
func SortDep(ctx *gin.Context) {
	var msg errors.ErrModel
	var dep department.Pdepartment
	if err := ctx.ShouldBind(&dep); err != nil {
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	} else {
		msg = department.SortDep(dep.DepID)
	}
	msg.RetErr(ctx)
}
