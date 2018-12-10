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
	"github.com/gin-gonic/gin"
	middle "github.com/boxproject/apiServer/middleware"
	"github.com/boxproject/apiServer/errors"
	log "github.com/alecthomas/log4go"
	"github.com/boxproject/apiServer/service/template"
	"github.com/boxproject/apiServer/service/logger"
	"github.com/boxproject/apiServer/common"
)

// CreateTemplate create a new template flow
func CreateTemplate(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*middle.CustomClaims)
	//claims := &middle.CustomClaims{ID:1}
	var tep template.AddTmp
	var msg errors.ErrModel
	if err := ctx.ShouldBind(&tep); err != nil {
		log.Error("Bind Param Error", err)
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	} else {
		msg = template.CreateNewTemplate(&tep, claims)
		// 操作日志
		if msg.Code == errors.Success {
			if data, ok := msg.Data.(string); ok {
				logErr := logger.AddLog("template", "", claims.Account, common.LoggerTemplateNew, data)
				if logErr != nil {
					log.Error("模板日志添加失败", claims.Account)
				}
			}
		}
	}
	msg.RetErr(ctx)
}

// TemplateList get the template flow list
func TemplateList(ctx *gin.Context) {
	var tep template.AddTmp
	var msg errors.ErrModel
	claims := ctx.MustGet("claims").(*middle.CustomClaims)
	if err := ctx.ShouldBind(&tep); err != nil {
		log.Error("Bind Param Error", err)
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	} else {
		msg = template.TemplateList(tep.Types, claims.UserType, claims.ID)
	}
	msg.RetErr(ctx)
}

// FindTemplateById get the template flow info by id
func FindTemplateById(ctx *gin.Context) {
	var tep template.TemplateVo
	var msg errors.ErrModel
	if err := ctx.ShouldBind(&tep); err != nil {
		log.Error("Bind Param Error", err)
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	} else {
		claims := ctx.MustGet("claims").(*middle.CustomClaims)
		msg = template.FindTemplateById(&tep, claims.ID)
	}
	msg.RetErr(ctx)
}

// VerifyTemplate owner approval the template flow
func VerifyTemplate(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*middle.CustomClaims)
	var tep template.TemplateVo
	var msg errors.ErrModel
	if err := ctx.ShouldBind(&tep); err != nil {
		log.Error("Bind Param Error", err)
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	} else {
		msg = template.VerifyTemplate(&tep, claims)
	}
	msg.RetErr(ctx)
}

// CancelTemplate disable one template flow
func CancelTemplate(ctx *gin.Context) {
	claims := ctx.MustGet("claims").(*middle.CustomClaims)
	var c_temp template.ParamCancelTemplate
	var msg errors.ErrModel
	if err := ctx.ShouldBind(&c_temp); err != nil {
		log.Error("Bind Param Error", err)
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	} else {
		msg = template.CancelTemp(&c_temp, claims)
	}
	msg.RetErr(ctx)
}

// TempStatistics get template info by related accounts
func TempStatistics(ctx *gin.Context) {
	var temp template.TemplateVo
	var msg errors.ErrModel
	if err := ctx.ShouldBind(&temp); err != nil {
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	}
	msg = template.TempStatistics(temp.AppName)
	msg.RetErr(ctx)
}

// TempStatistics counts the related transfer number by template
func TxNumByTemplate(ctx *gin.Context) {
	var temp template.TemplateVo
	var msg errors.ErrModel
	if err := ctx.ShouldBind(&temp); err != nil {
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	}
	msg = template.TxNumByTempID(temp.Id)
	msg.RetErr(ctx)
}





