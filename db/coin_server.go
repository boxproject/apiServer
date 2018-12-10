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
	"math/rand"
	"time"
	rerr "github.com/boxproject/apiServer/errors"
	"github.com/jinzhu/gorm"
	"github.com/prometheus/common/log"
	"github.com/boxproject/apiServer/common"
)

type CoinService struct {
}

// CoinIsExist check if the coin is existed by given token address and its name
func (*CoinService) CoinIsExist(address, name string) (bool, error) {
	var coin Coin
	err := rdb.First(&coin, "tokenAddress = ? or name = ?", address, name).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, rerr.MSG_6001
}

// AddCoin add noe token to this system
func (*CoinService) AddCoin(symbol, name string, precise int, tokenAddress string) error {
	var eth Coin
	tx := rdb.Begin()
	err := tx.First(&eth, "name = 'ETH'").Error
	if err != nil {
		tx.Rollback()
		return err
	}
	coin := Coin{
		Name:         symbol,
		FullName:     name,
		Precise:      precise,
		TokenAddress: tokenAddress,
		TokenType:    eth.ID,
		Balance:      "0",
	}
	err = tx.Create(&coin).Error
	if err != nil {
		tx.Rollback()
		log.Error(err)
		return err
	}
	var coinDB Coin
	err = tx.Raw("select * from coin where tokenAddress = ?", coin.TokenAddress).Scan(&coinDB).Error
	if err != nil {
		tx.Rollback()
		log.Error(err)
		return err
	}
	err = tx.Exec("update deposit set coinId = ?,coinName = ?,isUpdate = 3,precise = ? where tokenAddress = ?", coinDB.ID, coinDB.Name, precise, coinDB.TokenAddress).Error
	if err != nil {
		tx.Rollback()
		log.Error(err)
		return err
	}
	tx.Commit()
	return err
}

// CoinList get the coin list from db
func (*CoinService) CoinList(typecoin int) ([]Coin, error) {
	var coins []Coin
	var err error
	if typecoin == common.CoinListMasterCoinType { //主链币
		err = rdb.Find(&coins, "tokenType= 0 ").Error
	} else if typecoin == common.CoinListTokenType { //代币
		err = rdb.Find(&coins, "tokenType != 0").Error
	} else if typecoin == common.CoinListETHCoinType { //ETH
		var eth Coin
		err = rdb.First(&eth, "name = 'ETH'").Error
		err = rdb.Find(&coins, "name = 'ETH' or tokenType = ?", eth.ID).Error
	} else if typecoin == common.CoinListAllCoinType { //全部
		err = rdb.Find(&coins).Error
	} else if typecoin == common.CoinListAvailableCoinType { //启用
		err = rdb.Find(&coins, "available = 0").Error
	}
	return coins, err
}

// CoinBalance the struct of balance info by coin
type CoinBalance struct {
	Name     string `gorm:"column:name"`
	Balance  string `gorm:"column:balance"`
	FullName string `gorm:"column:fullName"`
	Id       int
}

// CoinBalance get all master coin balance info
func (*CoinService) CoinBalance() ([]CoinBalance, error) {
	var coin []Coin
	var err error
	var cbs []CoinBalance
	rows, err := rdb.Find(&coin, "tokenType= 0").Rows()
	defer rows.Close()
	for rows.Next() {
		var cb CoinBalance
		rdb.ScanRows(rows, &cb)
		cbs = append(cbs, cb)
	}
	return cbs, err
}

// QRcode get the deposit address
func (*CoinService) QRcode(id int) (map[string]string, error) {
	var addr Address
	var childs []Address
	mapss := make(map[string]string)
	err := rdb.First(&addr, "coinid = ? and type = 0", id).Error
	err = rdb.Find(&childs, "coinid = ? and type = 1", id).Error
	mapss["MainAddress"] = addr.Address
	rand.Seed(time.Now().Unix())
	if len(childs) != 0 {
		randNum := rand.Intn(len(childs))
		mapss["ChildAddress"] = childs[randNum].Address
		mapss["Index"] = childs[randNum].Deep
		mapss["CoinName"] = childs[randNum].CoinName
	}
	return mapss, err
}

// CoinStauts enable or disable the given coin
func (*CoinService) CoinStauts(id, status int) (*Coin, error) {
	coin := Coin{
		ID: id,
	}
	err := rdb.First(&coin).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			//币种不存在
			return nil, nil
		}
		//查询错误
		return nil, err
	}
	err = rdb.Model(&coin).Where("id = ?", id).Update("available", status).Error
	if err != nil {
		log.Error("update coin status error", err)
		return nil, err
	}
	return &coin, nil
}

// CoinVo coins info of available coin
type CoinVo struct {
	ID           int
	Name         string `gorm:"column:name"`
	FullName     string `gorm:"column:fullName"`
	Precise      int    `gorm:"column:precise"`
	Balance      string `gorm:"column:balance"`
	TokenType    int    `gorm:"column:tokenType"`
	TokenAddress string `gorm:"column:tokenAddress"`
	Available    int    `gorm:"column:available"`
	Currency     string `gorm:"column:currency"`
}

// CoinListByTransfer get coin info available
func (*CoinService) CoinListByTransfer() ([]CoinVo, error) {
	var coins []CoinVo
	err := rdb.Raw("select c.*,(case cc.name when 'ETH' then 'ERC20' else c.name end) as currency from coin c " +
		"left join(select id,name,tokenType from coin) cc on c.tokenType = cc.id where c.available = 0").Scan(&coins).Error
	return coins, err
}

// GetCoinById get coin info by id
func (*CoinService) GetCoinById(id int) (*Coin, error) {
	var coin Coin
	err := rdb.First(&coin, id).Error
	if err != nil {
		log.Error("getCoin err ", err)
	}
	fmt.Println("name: ", coin.Name)
	return &coin, err
}

// GetCoins get all coins
func (*CoinService) GetCoins() ([]Coin, error) {
	var coins []Coin
	err := rdb.Find(&coins).Error
	return coins, err
}

// GetETH get eth info from db
func (*CoinService) GetETH() (*Coin, error) {
	var coin Coin
	err := rdb.First(&coin, "name= 'ETH'").Error
	return &coin, err
}
