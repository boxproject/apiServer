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
package config

import (
	log "github.com/alecthomas/log4go"
	"github.com/spf13/viper"
	"os"
	"path"
	"os/exec"
	"path/filepath"
	"os/user"
	"runtime"
	"github.com/boxproject/apiServer/utils"
	"github.com/boxproject/apiServer/common"
)

var (
	rootPath string
	filePath string
	Conf     *Config
)

func init() {
	viper.SetConfigType(common.TOML_FILE)
	ProjectDir := utils.GetCurrentDirectory()
	ConfigDir := common.CONF_PATH
	config_file := path.Join(ProjectDir, ConfigDir, common.CONF_FILE)
	viper.SetConfigFile(config_file)
	main, _ := exec.LookPath(os.Args[0])
	file, _ := filepath.Abs(main)
	rootPath = path.Dir(file)
}

// LoadConfig load the config file which the system needed.
func LoadConfig() (*Config, error) {
	if err := viper.ReadInConfig(); err != nil {
		log.Error("Error reading config file", err)
		return nil, err
	}
	cfg := &Config{}
	err := viper.Unmarshal(&cfg)
	if err != nil {
		log.Error("unable to decode into struct", err)
		return nil, err
	}
	Conf = cfg
	return cfg, nil
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

// DefaultConfigDir returns the default path when the config path is not existed.
func DefaultConfigDir() string {
	home := homeDir()
	if home != "" {
		if runtime.GOOS == "darwin" {
			return filepath.Join(home, ".apiServer")
		} else if runtime.GOOS == "windows" {
			return filepath.Join(home, "AppData", "Roaming", "apiServer")
		} else {
			return filepath.Join(home, ".apiServer")
		}
	}
	return ""
}

// InitLogger initialize the configuration file for log4js
// this configuration file specifies the path where the log
// data stored.
func InitLogger() {
	logFile := path.Join(rootPath, common.LOG_CONF)
	for i := 0; i < 3; i++ {
		if _, err := os.Stat(logFile); !os.IsNotExist(err) {
			break
		}
		if i == 0 {
			logFile = path.Join(filePath, common.LOG_CONF)
		} else if i == 1 {
			logFile = path.Join(DefaultConfigDir(), common.LOG_CONF)
		}
	}
	log.LoadConfiguration(logFile)
}
