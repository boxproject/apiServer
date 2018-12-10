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
package utils

import (
	"path"
	"os"
	"os/exec"
	"path/filepath"
	"os/user"
	"runtime"
	logger "github.com/alecthomas/log4go"
)

var (
	rootPath string
	filePath string
)

func init() {
	main, _ := exec.LookPath(os.Args[0])
	file, _ := filepath.Abs(main)
	rootPath = path.Dir(file)
}

func DefaultConfigDir() string {
	home := homeDir()
	if home != "" {
		if runtime.GOOS == "darwin" {
			return filepath.Join(home, ".bcmonitor")
		} else if runtime.GOOS == "windows" {
			return filepath.Join(home, "AppData", "Roaming", "bcmonitor")
		} else {
			return filepath.Join(home, ".bcmonitor")
		}
	}

	return ""
}

func homeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}

	if usr, err := user.Current(); err == nil {
		return usr.HomeDir
	}

	return ""
}

func InitLogger() {
	logFile := path.Join(rootPath, "log.xml")
	for i := 0; i < 3; i++ {
		if _, err := os.Stat(logFile); !os.IsNotExist(err) {
			break
		}
		if i == 0 {
			logFile = path.Join(filePath, "log.xml")
		} else if i == 1 {
			logFile = path.Join(DefaultConfigDir(), "log.xml")
		}
	}
	logger.LoadConfiguration(logFile)
}
