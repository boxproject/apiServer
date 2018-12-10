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
	"strconv"
	"github.com/jinzhu/gorm"
	log "github.com/alecthomas/log4go"
	"strings"
)

type TxinfosService struct {
}

var ConfigMap map[string]string
var AllMasterAddress map[string]string
var USDT_PPID = ""

func init() {
	configs := InitConfigs()
	if len(configs) > 0 {
		ConfigMap = make(map[string]string)
	}
	for i := 0; i < len(configs); i++ {
		ConfigMap[configs[i].Key] = configs[i].Value
	}
	ConfigMap["combine_account"] = "0"
	var coin Coin
	err := rdb.First(&coin, "name = ?", "USDT").Error
	if err != nil {
		log.Error("初始化币种加载失败")
	}
	USDT_PPID = coin.TokenAddress

}

// FindTxinfos
func (*TxinfosService) FindTxinfos() []TxInfos {
	var txinfos []TxInfos
	var err error
	if ConfigMap["getbalance_time"] == "" {
		err = rdb.Find(&txinfos).Error
	} else {
		err = rdb.Where("created_at > ?", ConfigMap["getbalance_time"]).Find(&txinfos).Error
	}
	if err != nil {
		log.Error("查询交易失败", err)
	}
	return txinfos
}

// FindCoin 查询币种
//  0.没有记录 1.查询失败 2.有结果
func (*TxinfosService) FindCoin(coinName string, tokenAddress string) (*Coin, int) {
	var coin Coin
	var err error
	if tokenAddress == "" {
		err = rdb.Where("name = ? ", coinName).First(&coin).Error
	} else {
		err = rdb.Where("tokenAddress = ?", tokenAddress).First(&coin).Error
	}
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, 1
		} else {
			log.Error("查询币种失败", err)
			return nil, 0
		}
	}
	return &coin, 2
}

// SaveDeposit 充值
func (*TxinfosService) SaveDeposit(deposit *Deposit) {
	err := rdb.Create(deposit).Error
	if err != nil {
		log.Error("充值", err)
	}
}

type DepositEntity struct {
	CoinName     string `gorm:"column:coinName"`
	TokenType    int    `gorm:"column:tokenType"`
	TokenAddress string `gorm:"column:tokenAddress"`
	ToAddr       string `gorm:"column:toAddr"`
	CoinId       int    `gorm:"column:coinId"`
	AddressId    int    `gorm:"column:addressId"`
	DepositId    string `gorm:"column:depositId"`
	Confirm      int    `gorm:"column:confirm"`
	Nowconfirm   int    `gorm:"column:nowconfirm"`
	ExtValid     int    `gorm:"column:ext_valid"`
	Deep         string `gorm:"column:deep"`
}

// FindDeposit 充值列表
func (*TxinfosService) FindDeposit() []DepositEntity {
	var deposits []DepositEntity
	err := rdb.Raw("select addr.deep,ti.ext_valid as ext_valid,ti.confirm as nowconfirm,dep.confirm as confirm,dep.id as depositId,addr.id as addressId,c.name as coinName," +
		"c.tokenType,c.tokenAddress,dep.toAddr as toAddr,c.id as coinId from deposit dep" +
		" left join coin c on dep.coinId = c.id left join address addr on addr.address = dep.toAddr and addr.coinId = if(c.tokenType=0,c.id,c.tokenType)  left join tx_infos ti on ti.txId = dep.txId  " +
		"where isUpdate = 3  order by updatedAt limit 1 ").Scan(&deposits).Error
	if err != nil {
		log.Error("充值列表", err)
	}
	return deposits
}

// UpdateDeposit 充值更新
func (*TxinfosService) UpdateDeposit(deposit *Deposit) {
	err := rdb.Model(&Deposit{}).Update(deposit).Error
	if err != nil {
		log.Error("充值", err)
	}
}

type BalanceEntity struct {
	CoinName     string `gorm:"column:coinName"`
	TokenAddress string `gorm:"column:tokenAddress"`
	Address      string `gorm:"column:address"`
	CapitalId    int    `gorm:"column:capitalId"`
	TokenType    int    `gorm:"column:tokenType"`
}

// GetBalanceList 应该获取的余额列表
func (*TxinfosService) GetBalanceList() []BalanceEntity {
	var deposits []BalanceEntity
	err := rdb.Raw("select c.name as coinName,c.tokenType,c.tokenAddress,cap.id as capitalId,addr.address from capital cap left join address addr " +
		"on addr.id = cap.addressId left join coin c on cap.coinId = c.id  order by createdAt limit 1").Scan(&deposits).Error
	if err != nil {
		log.Error("列表", err)
	}
	return deposits
}

//UpdateBalance  更新余额
func (*TxinfosService) UpdateBalance(capital *Capital) {
	err := rdb.Exec("update capital set balance = ?,createdAt = ? where id = ?", capital.Balance, capital.CreatedAt, capital.ID).Error

	if err != nil {
		log.Error("更新余额", err)
	}
}

//UpdateBalance  更新余额
func (*TxinfosService) UpdateBalances(coinName string, address string, balance string) {
	err := rdb.Exec("update capital cap,address addr set cap.balance = ? "+
		"where  addr.id = cap.addressId and cap.coinId = addr.coinId and addr.address = ? "+
		"and addr.coinName=?", balance, address, coinName).Error
	if err != nil {
		log.Error("更新余额", err)
	}
}

// InitConfigs 获取配置
func InitConfigs() []Configs {
	var configs []Configs
	err := rdb.Find(&configs).Error
	if err != nil {
		log.Error("更新余额", err)
	}
	return configs
}

// UpadateConfigs 更新
func (*TxinfosService) UpadateConfigs(key string, value string) error {
	err := rdb.Model(&Configs{}).Where("con_key = ?", key).Update("con_value", value).Error
	if err != nil {
		log.Error("更新余额", err)
	}
	return err
}

func (*TxinfosService) SumBalance() error {
	var capital []Capital
	conn := rdb.Begin()
	capitalErr := conn.Raw("select replace(rtrim(replace(cast(sum(cast(balance as decimal(38,20))) as char),'0',' ')),' ','0') as balance , coinId from capital group by coinId").Find(&capital).Error
	if capitalErr != nil {
		conn.Rollback()
		return capitalErr
	}
	if len(capital) == 0 {
		conn.Rollback()
		return nil
	}
	query := ""
	for _, v := range capital {
		query += " WHEN " + strconv.Itoa(v.CoinId) + " THEN " + v.Balance
	}
	sql := "CASE id" + query + " else balance END"
	err := conn.Table("coin").Updates(map[string]interface{}{"balance": gorm.Expr(sql)}).Error
	if err != nil {
		conn.Rollback()
		return err
	}
	conn.Commit()
	return nil
}

type Txinfos struct {
	Address   string `gorm:"column:address"`
	AddressId int    `gorm:"column:addressId"`
	CoinId    int    `gorm:"column:coinId"`
	CoinName  string `gorm:"column:coinName"`
}

// FindTokenAddress 添加币种专用, 更新代币余额
func (*TxinfosService) FindTokenAddress(token string) []Txinfos {
	var txinfos []Txinfos
	err := rdb.Raw("select JSON_UNQUOTE(json_extract(tx_obj,'$.ExtOut[0].Addr')) as address from tx_infos where "+
		"token = ?", token).Scan(&txinfos).Error
	if err != nil {
		log.Error("添加币种专用，更新代笔余额", err)
	}
	return txinfos
}

// UpdateTokenBalance 更新代笔余额
func (*TxinfosService) UpdateTokenBalance(maps map[string]string, token string) {
	var txinfos []Txinfos
	var addrs = map[string]Txinfos{}
	err := rdb.Raw("select c.id as coinId,a.id as addressId,c.name as coinName,a.address as address "+
		"from coin c left join address a on c.tokenType = a.coinId where c.tokenAddress= ?", token).Scan(&txinfos).Error
	if err != nil {
		log.Error("添加，更新代笔余额", err)
		return
	}
	for j := 0; j < len(txinfos); j++ {
		addrs[txinfos[j].Address] = txinfos[j]
	}
	sqlStr := "insert into capital(addressId,coinName,balance,coinId) values "
	vals := []interface{}{}
	rowSQL := "(?,?,?,?)"
	var inserts []string
	for k, v := range maps {
		inserts = append(inserts, rowSQL)
		vals = append(vals, addrs[k], addrs[k].CoinName, v, addrs[k].CoinId)
	}
	sqlStr = sqlStr + strings.Join(inserts, ",")
	err = rdb.Exec(sqlStr, vals...).Error
	if err != nil {
		log.Error("更新代币余额保存错误", err)
	}
}

// UpdateDeposits 更新充值 --添加代币时专用
func (*TxinfosService) UpdateDeposits(coin *Coin) {
	err := rdb.Exec("update deposit set coinId = ?,coinName = ?,isUpdate = 1 "+
		"where tokenAddress = ?", coin.ID, coin.Name, coin.TokenAddress).Error
	if err != nil {
		log.Error("更新充值错误", err)
	}

}
