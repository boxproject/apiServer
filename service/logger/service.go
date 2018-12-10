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
package logger

import (
	"github.com/boxproject/apiServer/db"
	reterror "github.com/boxproject/apiServer/errors"
	logger "github.com/alecthomas/log4go"
	"github.com/spf13/viper"
	"github.com/boxproject/apiServer/common"
	"reflect"
	"text/template"
	"bytes"
	"github.com/gin-gonic/gin"
	"strings"
	"github.com/boxproject/apiServer/utils"
)

// Logger log db struct
type Logger struct {
	LogType string `form:"log_type"`
	Limit   int    `form:"limit"`
	Start   int    `form:"start"`
	Pos     string `form:"pos"`
}

// AddLog add log
func AddLog(logType string, note string, account string, loggerType string, pos string) error {
	logInfo := &db.Logger{
		Operator: account,
		LogType:  logType,
		Detail:   loggerType,
		Note:     note,
		Pos:      pos,
	}
	ds := &db.LogService{}
	err := ds.AddLog(logInfo)
	if err != nil {
		return reterror.MSG_1007
	}
	return nil
}

// GetLogs get logs
func GetLogs(ctx *gin.Context, log *Logger) reterror.ErrModel {
	if log.LogType == "" {
		return reterror.ErrModel{Code: reterror.ParamsNil, Err: reterror.PARAMS_NULL}
	}
	ds := &db.LogService{}
	list, err := ds.GetLogs(log.LogType, log.Limit, log.Start, log.Pos)
	if err != nil {
		logger.Error("get log info error", err)
		return reterror.ErrModel{Code: reterror.Failed, Err: reterror.SYSTEM_ERROR}
	}
	result := formatLoggerMsg(ctx, list)
	return reterror.ErrModel{Code: reterror.Success, Err: err, Data: result}
}

func formatLoggerMsg(ctx *gin.Context, logList []db.Logger) []db.Logger {
	language := utils.ConvertLangTag(ctx.GetHeader(common.HeaderLangKey))
	// 更新detail
	for i := 0; i < len(logList); i++ {
		value := reflect.ValueOf(logList[i].Detail)
		key := strings.ToLower(value.String())
		tmpl, _ := template.New("content").Parse(OperLoggerTempMap[language][key].(string))
		buf := new(bytes.Buffer)
		tmpl.Execute(buf, logList[i])
		logList[i].Detail = buf.String()
	}
	return logList
}

func init() {
	// zh
	vzh := viper.New()
	vzh.SetConfigType(common.TOML_FILE)
	vzh.SetConfigFile(ZhLoggerFilePath)
	if err := vzh.ReadInConfig(); err != nil {
		logger.Error("read logger_zh.toml error", err)
	}
	OperLoggerTempMap[common.LangZhType] = vzh.AllSettings()

	// en
	ven := viper.New()
	ven.SetConfigType(common.TOML_FILE)
	ven.SetConfigFile(EnLoggerFilePath)
	if err := ven.ReadInConfig(); err != nil {
		logger.Error("read logger_en.toml error", err)
	}
	OperLoggerTempMap[common.LangEnType] = ven.AllSettings()
}
