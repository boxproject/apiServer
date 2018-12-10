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
	"github.com/boxproject/apiServer/service/logger"
	"github.com/boxproject/apiServer/errors"
	"github.com/gin-gonic/gin"
)

// GetLogs get log info
func GetLogs(ctx *gin.Context) {
	var logs logger.Logger
	var msg errors.ErrModel
	if err := ctx.ShouldBind(&logs); err != nil {
		log.Error("Bind Param Error", err)
		msg = errors.ErrModel{Code: errors.Failed, Err: errors.SYSTEM_ERROR}
	} else {
		msg = logger.GetLogs(ctx, &logs)
	}
	msg.RetErr(ctx)
}
