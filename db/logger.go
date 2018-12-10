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

// LogService 操作日志 db
type LogService struct {
}

// AddLog 添加日志
func (*LogService) AddLog(logger *Logger) error {
	err := rdb.Create(logger).Error
	return err
}

// GetLogs 获取日志
func (*LogService) GetLogs(logType string, limit int, offset int, pos string) ([]Logger, error) {
	var logger []Logger
	if limit == 0 {
		limit = -1
	}
	err := rdb.Limit(limit).Offset(offset).Where("logType=? and pos=?", logType, pos).Order("createdAt desc").Order("id desc").Find(&logger).Error
	return logger, err
}
