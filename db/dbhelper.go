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
	"fmt"
	log "github.com/alecthomas/log4go"
	"github.com/boxproject/apiServer/config"
	"github.com/jinzhu/gorm"
	"github.com/go-sql-driver/mysql"
	_"github.com/jinzhu/gorm/dialects/mysql"
	"github.com/boxproject/apiServer/common"
	"github.com/spf13/viper"
	"github.com/theplant/batchputs"
	"os"
)

var rdb *gorm.DB

func init() {
	var cfg *config.Config
	var err error
	if config.Conf == nil {
		cfg, err = config.LoadConfig()
		if err != nil {
			log.Error("Init DB while loading config error.", err)
			panic(err)
		}
		config.Conf = cfg
	}

	dbOpts := mysql.Config{
		User:      cfg.Database.User,
		Passwd:    cfg.Database.Password,
		DBName:    cfg.Database.DbName,
		Net:       "tcp",
		Addr:      cfg.Database.Host,
		ParseTime: true,
		Params:    map[string]string{"charset": "utf8", "loc": "Local"},
	}
	fmt.Println("db", dbOpts.FormatDSN())
	db, err := gorm.Open("mysql", dbOpts.FormatDSN())
	if err != nil {
		log.Error("Init DB Error", err)
		db.Close()
	}

	db.DB().SetMaxOpenConns(cfg.Database.MaxOpen)
	db.DB().SetMaxIdleConns(cfg.Database.MaxIdle)
	rdb = db
	if config.Conf.Database.DebugMode == true {
		rdb = db.Debug()
	}
	// 更新权限表auth
	vAuthZh := viper.New()
	vAuthZh.SetConfigType(common.TOML_FILE)
	vAuthZh.SetConfigFile(AuthTypeZhPath)
	if err := vAuthZh.ReadInConfig(); err != nil {
		log.Error("read auth_zh.toml error", err)
	}
	AuthMaps[common.LangZhType] = vAuthZh.AllSettings()
	vAuthEn := viper.New()
	vAuthEn.SetConfigType(common.TOML_FILE)
	vAuthEn.SetConfigFile(AuthTypeEnPath)
	if err := vAuthEn.ReadInConfig(); err != nil {
		log.Error("read auth_en.toml error", err)
	}
	AuthMaps[common.LangEnType] = vAuthEn.AllSettings()
	columns := []string{"id", "name", "authType"}
	rows := [][]interface{}{}
	i := 1
	for k, v := range AuthMaps[common.LangZhType] {
		rows = append(rows, []interface{}{
			i, v, k,
		})
		i++
	}
	err = batchputs.Put(rdb.DB(), os.Getenv("DB_DIALECT"), "auth", "id", columns, rows)
	if err != nil {
		log.Error("初始权限信息错误", err)
	}
}
