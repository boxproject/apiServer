// Copyright 2018. bolaxy.org authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		 http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package versionlog

import (
	"path"
	"io/ioutil"
	"encoding/json"
	"github.com/prometheus/common/log"
	"github.com/boxproject/apiServer/utils"
)

var LogMap = make(map[string]interface{})

func init() {

	ProjectDir := utils.GetCurrentDirectory()
	ErrCodeDir := "/versionlog/"
	logFile := path.Join(ProjectDir, ErrCodeDir, "log.json")
	bytes, _ := ioutil.ReadFile(logFile)
	err := json.Unmarshal(bytes, &LogMap)
	if err != nil {
		log.Error("versionlog错误", err)
	}
}
