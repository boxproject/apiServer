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
	"strings"
	"time"
	"encoding/json"
	log "github.com/alecthomas/log4go"
	voucher "github.com/boxproject/apiServer/rpc"
	"github.com/go-errors/errors"
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
	"database/sql"
	"github.com/satori/go.uuid"
	"github.com/boxproject/apiServer/common"
)

// CoinAccount represents the coin statistical data
// it is applicable to account and asset information
// corresponding to the statistical currency.
type CoinAccount struct {
	ID          int    // currency id
	Name        string // currency name
	Main        int    // number of main accounts
	Child       int    // number of sub accounts
	HasToken    bool   // whether it contains tokens
	HasMerge    bool   // whether it can be merged
	HasAddChild bool   // whether it can add sub accounts
}

// AsCapital represents the top 10 assets
type AsCapital struct {
	ID       int    // address id
	Address  string // address
	Tag      string // address tag
	Balance  string // balance of the address
	Amount   int    // the number of address tags
	Type     int    // address type 0-master address 1-normal address 2-contract address
	TagIndex int    `gorm:"column:tagIndex"` // tag index
}

// TrRecord represents the struct of transfer record info
type TrRecord struct {
	ApplyReason string `gorm:"column:applyReason"` // reason to apply this transfer
	Status      int    `gorm:"column:status"`      // status of this record 0-Approval 1-Approved 2-Reject 3-Transfer pending 4-Transfer failed
	// 5.Transfer success 6-Cancel 7-Invalid 8-Approval Expired 9-Transfer Expired 10-Disable User
	// 11-Disable template 12-Insufficient balance 100-Deposit success
	Amount    string    `gorm:"column:amount"`    // transfer amount
	CreatedAt time.Time `gorm:"column:createdAt"` // Order creation time
	Precise   int       `gorm:"column:precise"`   // Coin accuracy
	Type      int       // data type 0-deposit 1-transfer
	Link      string    //  transfer id info on blockchain explorer
	TxID      string    `gorm:"column:txId"` // transfer id on blockchain
}

// TokenBalance represents the struct of the balance of tokens
type TokenBalance struct {
	CoinName string `gorm:"column:coinName"` // coin name
	CoinId   int    `gorm:"column:coinId"`   // coin id in sql
	Balance  string `gorm:"column:balance"`  // balance of the coin
}

type AccManageService struct {
}

// when transfer the money from the sub-account to the main account, the balance of sub-account
// must no less than `minBalance`
const minBalance = "0.1"

// Statistical Statistical account information
func (*AccManageService) Statistical() []CoinAccount {
	var accounts []CoinAccount
	rows, err := rdb.Raw("select c.id id, c.name ,sum(type=0) main ,sum(type =1) child from coin c join address ad on c.id = ad.coinId and tokenType=0 and available = 0  group by name,id").Rows()
	defer rows.Close()
	for rows.Next() {
		var account CoinAccount
		rdb.ScanRows(rows, &account)
		accounts = append(accounts, account)
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error("多账户统计", err)
		return nil
	}
	return accounts
}

// ChildStatistical Statistical sub-account information
func (*AccManageService) ChildStatistical(id int) (int, int) {
	var childCount int
	var otherCount int
	row := rdb.Table("address").Select("count(0)").Where("coinId = ? and isDeleted = 0 and type = 1 ", id).Row()
	row.Scan(&childCount)
	row = rdb.Table("address").Select("count(0)").Where("coinId != ? and isDeleted = 0 and type = 1 ", id).Row()
	row.Scan(&otherCount)
	return childCount, otherCount
}

// CountAllChild Statistical sub-address number and assets of them
func (*AccManageService) CountAllChild(id int) (string, int, error) {
	// amount of balance
	var balance string
	var amount int
	// coin info
	var coin Coin
	// 获取ETH信息
	var ethInfo Coin
	err := rdb.First(&ethInfo, "name = ?", "ETH").Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Error("未获取到ETH币种信息")
		}
		log.Error("获取ETH信息失败")
		return "0", 0, err
	}
	err = rdb.First(&coin, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Error("未获取到对应币种信息", id)
		}
		return "0", 0, err
	}
	if coin.ID != ethInfo.ID {
		if coin.TokenType != ethInfo.ID {
			return "0", 0, errors.New("目前账户合并仅支持ETH.")
		}
	}
	row := rdb.Raw("SELECT count(0) amount FROM `capital` join address on capital.addressId=address.id where capital.balance>0 and capital.coinId= ? and address.type =1", coin.ID).Row()
	row.Scan(&amount)
	row = rdb.Raw("SELECT sum(capital.balance) balance FROM `capital` join address on capital.addressId=address.id where capital.balance>0 and capital.coinId= ? and address.type =1", coin.ID).Row()
	row.Scan(&balance)
	return balance, amount, nil
}

// AddressTopTen statistical the 10 assets address
func (*AccManageService) AddressTopTen(id int) ([]AsCapital, string, error) {
	//查出资金总额
	var total string
	row := rdb.Table("address").Select("sum(c.balance) as total").Joins("join capital c on address.id=c.addressId and address.coinId=? and address.coinName = c.coinName and address.type != 2", id).Row()
	row.Scan(&total)
	var asCapitals []AsCapital
	//查列表
	rows, err := rdb.Table("address").Select("address.id, address.type ,address.address,address.tag,c.balance").Joins("left join capital c on address.id = c.addressId and address.coinName = c.coinName").Where("address.isDeleted=0  and address.coinId=? and address.type != 2", id).Order("balance+0 desc,id desc").Limit(9).Rows()
	defer rows.Close()
	for rows.Next() {
		var asCapital AsCapital
		var address Address
		rdb.ScanRows(rows, &asCapital)
		if asCapital.Type == common.MainAddressType {
			asCapital.Tag = "主账户"
		}
		rdb.Last(&address, "tag = ?", asCapital.Tag).Order("tagIndex asc")
		asCapital.Amount = address.TagIndex
		asCapitals = append(asCapitals, asCapital)
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error("账户资金top10", err)
		return nil, "", err
	}
	return asCapitals, total, nil
}

// AddressList list of address info
func (*AccManageService) AddressList(id int, condition string, page, tp int) ([]AsCapital, string, error) {

	// main currency
	var coin Coin
	coinName := "ETH"
	err := rdb.First(&coin, "id = ?", id).Error
	if coin.TokenType != 0 {
		coinName = coin.Name
		var main Coin
		err = rdb.First(&main, "id = ?", coin.TokenType).Error
		coin = main
	}
	if err != nil {
		return nil, "", err
	}
	// balance of amount
	var total string
	row := rdb.Table("capital").Select("sum(capital.balance) as total").Where("coinId = ?", id).Row()

	row.Scan(&total)
	var asCapitals []AsCapital

	Db := rdb.Table("address")
	// 查询列表
	if condition != "" {
		Db = Db.Where("address.tag like ? or address.address like ?", "%"+condition+"%", "%"+condition+"%")
	}
	Db = Db.Limit(10).Offset(page * 10)
	var where string
	if tp == 0 {
		where = "address.isDeleted=0 and address.coinId=? and address.type !=2 "
	} else { //子账户
		Db = rdb.Table("capital")
		//查询列表
		if condition != "" {
			Db = Db.Where("address.tag like ? or address.address like ?", "%"+condition+"%", "%"+condition+"%")
		}
		Db = Db.Limit(10).Offset(page * 10)
		var rows *sql.Rows
		if coinName == "ETH" {
			rows, err = Db.Select("address.id,address.type,address.address,address.tag,capital.balance,address.tagIndex").Joins("join address on capital.addressId=address.id").Where("capital.balance>0.1 and capital.coinId=? and address.type =1", id).Rows()
		} else {
			rows, err = Db.Select("address.id,address.type,address.address,address.tag,capital.balance,address.tagIndex").Joins("join address on capital.addressId=address.id").Where("capital.balance>0 and capital.coinId=? and address.type =1", id).Rows()
		}
		for rows.Next() {
			var asCapital AsCapital
			var address Address
			rdb.ScanRows(rows, &asCapital)
			rdb.Last(&address, "tag = ?", asCapital.Tag).Order("tagIndex asc")
			asCapital.Amount = address.TagIndex
			asCapitals = append(asCapitals, asCapital)
		}
		if err != nil && err != gorm.ErrRecordNotFound {
			log.Error("账户列表", err)
			return nil, "", err
		}
		return asCapitals, total, nil

		//where = "address.isDeleted=0 and address.coinId=? and address.type =1 and c.balance >0.1"
	}
	rows, err := Db.Select("address.id,address.type,address.address,address.tag,c.balance,address.tagIndex").Joins("left join capital c on address.id = c.addressId and address.coinName = c.coinName").Where(where, coin.ID).Order("type asc,balance+0 desc,id desc").Rows()
	defer rows.Close()
	for rows.Next() {
		var asCapital AsCapital
		var address Address
		rdb.ScanRows(rows, &asCapital)
		rdb.Last(&address, "tag = ?", asCapital.Tag).Order("tagIndex asc")
		asCapital.Amount = address.TagIndex
		if asCapital.Type == common.MainAddressType {
			asCapital.Tag = "主账户"
			asCapital.Amount = 1
		}
		asCapitals = append(asCapitals, asCapital)
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error("账户列表", err)
		return nil, "", err
	}
	return asCapitals, total, nil
}

// FindAddressById get the address info by address id
func (*AccManageService) FindAddressById(id int) (*Address, error) {
	//查出资金总额
	var address Address
	err := rdb.First(&address, "id=?", id).Error
	return &address, err
}

// SetTag set tag for address
func (*AccManageService) SetTag(id int, tag string) int64 {
	//查出资金总额
	var address Address
	num := rdb.Model(&address).Where("id = ?", id).Update("tag", tag).RowsAffected
	return num
}

// TransferRecord transfer order info, include deposit and withdraw
func (*AccManageService) TransferRecord(id int) ([]TrRecord, error) {
	var trRecordRows []TrRecord
	var address Address
	err := rdb.First(&address, "id=?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	transferStatus := []int{common.TransferSomeSuccess, common.TransferFailed, common.TransferSuccess, common.TransferApproveExpire, common.TransferInsufficientBalance, common.TransferMinerToFloatError, common.TransferCoinNotFound,
		common.TransferGenTxError, common.TransferVoucherReturnError, common.TransferDicemalError, common.TransferAddressIndexError, common.TransferDecodeMsgError, common.TransferDecodeDataError}
	rows, err := rdb.Raw("(select dep.amount, 0 as type,100 as status, dep.createdAt ,"+
		" (case tr.types when 2 then '从子账户转入' else '转入成功' end) as applyReason,"+
		" dep.precise,dep.txId as txId from deposit dep left join transfer tr on dep.txId = tr.txId"+
		" where toAddr = ? and isUpdate=2 and coinId = ?) "+
		" union all (select amount,1 as type,status,createdAt,applyReason,0 as precise, txId from transfer where fromAddress=? and coinName = ? and status in (?)) order by createdAt desc", address.Address, address.CoinId, address.Address, address.CoinName, transferStatus).Rows()
	defer rows.Close()
	for rows.Next() {
		var trRecord TrRecord
		err = rdb.ScanRows(rows, &trRecord)
		if err != nil {
			log.Error("scan transfer error", err)
			return nil, err
		}
		switch strings.ToUpper(address.CoinName) {
		case "ETH":
			trRecord.Link = ETH_LINK + trRecord.TxID
			break
		case "BTC":
			trRecord.Link = BTC_LINK + trRecord.TxID
			break
		case "LTC":
			trRecord.Link = LTC_LINK + trRecord.TxID + ".htm"
			break
		default:
			trRecord.Link = ETH_LINK + trRecord.TxID
		}
		if trRecord.Type == 0 { //转入时计算精度
			numDecimal, _ := decimal.NewFromString(trRecord.Amount)
			amountDecimal := numDecimal.Shift(int32(-trRecord.Precise))
			trRecord.Amount = amountDecimal.String()
		}
		trRecordRows = append(trRecordRows, trRecord)
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error("转账记录", err)
		return nil, err
	}
	return trRecordRows, nil
}

// TokenRecord transfer record of token
type tokenRecord struct {
	TrRecord
	CoinName string `gorm:"column:coinName"`
}

func (*AccManageService) TokenRecord(id, coinId int) ([]TrRecord, error) {
	var trRecordRows []TrRecord
	var coin Coin
	var err error
	err = rdb.First(&coin, "id = ?", coinId).Error
	if err != nil {
		log.Debug("转账记录获取币种失败", coinId)
		return nil, err
	}
	if id != 0 { //指定地址代币记录
		var address Address
		err = rdb.First(&address, "id=?", id).Error
		rows, err := rdb.Raw("(select amount, 0 as type,100 as status, createdAt ,'转入成功' as applyReason,precise from deposit where toAddr=? and coinId=? and isUpdate=2) union all (select amount,1 as type,status,createdAt,applyReason,0 as precise from transfer where fromAddress=? and coinName = ? and status > 0) order by createdAt desc", address.Address, coinId, address.Address, coin.Name).Rows()
		defer rows.Close()
		if err != nil {
			log.Debug("获取指定地址代币记录失败", err)
			return nil, err
		}
		for rows.Next() {
			var trRecord TrRecord
			rdb.ScanRows(rows, &trRecord)
			if trRecord.Type == 0 { //转入时计算精度
				numDecimal, _ := decimal.NewFromString(trRecord.Amount)
				amountDecimal := numDecimal.Shift(int32(-trRecord.Precise))
				trRecord.Amount = amountDecimal.String()
			}
			switch strings.ToUpper(address.CoinName) {
			case "ETH":
				trRecord.Link = ETH_LINK + trRecord.TxID
				break
			case "BTC":
				trRecord.Link = BTC_LINK + trRecord.TxID
				break
			case "LTC":
				trRecord.Link = LTC_LINK + trRecord.TxID + ".htm"
				break
			default:
				trRecord.Link = ETH_LINK + trRecord.TxID
			}
			trRecordRows = append(trRecordRows, trRecord)
		}
	} else {
		var coinInfo Coin
		err = rdb.First(&coinInfo, "id = ?", coinId).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				log.Error("对应币种未找到", coinId)
				return nil, nil
			}
			log.Error("查找币种信息失败", coinId, err)
		}
		rows, err := rdb.Raw(`(select amount, 0 as type,100 as status, createdAt ,'转入成功' as applyReason,precise, coinName, txId from deposit where coinId in (?) and isUpdate=2) union all (select amount, 1 as type, status, createdAt, applyReason,0 as precise, coinName, txId from transfer where coinName in (?) and status > 0) order by createdAt desc`, coinId, coinInfo.Name).Rows()
		if err != nil {
			log.Error("获取主链币及其代币交易记录错误", err)
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			var trRecord TrRecord
			var ttrRecord tokenRecord
			rdb.ScanRows(rows, &ttrRecord)
			switch strings.ToUpper(ttrRecord.CoinName) {
			case "ETH":
				trRecord.Link = ETH_LINK + ttrRecord.TxID
				break
			case "BTC":
				trRecord.Link = BTC_LINK + ttrRecord.TxID
				break
			case "LTC":
				trRecord.Link = LTC_LINK + ttrRecord.TxID + ".htm"
				break
			default:
				trRecord.Link = ETH_LINK + ttrRecord.TxID
			}
			trRecord.Type = ttrRecord.Type
			trRecord.Status = ttrRecord.Status
			trRecord.Amount = ttrRecord.Amount
			trRecord.CreatedAt = ttrRecord.CreatedAt
			trRecord.TxID = ttrRecord.TxID
			trRecord.Precise = ttrRecord.Precise
			trRecord.ApplyReason = ttrRecord.ApplyReason
			if ttrRecord.Type == 0 { //转入时计算精度
				numDecimal, _ := decimal.NewFromString(trRecord.Amount)
				amountDecimal := numDecimal.Shift(int32(-trRecord.Precise))
				trRecord.Amount = amountDecimal.String()
			}

			trRecordRows = append(trRecordRows, trRecord)
		}
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error("代币转账记录", err)
		return nil, err
	}

	return trRecordRows, nil
}

// TokenList token info list
func (*AccManageService) TokenList(id int) ([]TokenBalance, error) {
	var tokenBalances []TokenBalance
	rows, err := rdb.Table("capital").Select("capital.balance,capital.coinName,coin.id coinId").Joins("join coin on capital.coinName=coin.name").Where("capital.addressId =? and coin.tokenType != 0 ", id).Rows()
	defer rows.Close()
	for rows.Next() {
		var tokenBalance TokenBalance
		rdb.ScanRows(rows, &tokenBalance)
		tokenBalances = append(tokenBalances, tokenBalance)
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return tokenBalances, nil
}

// AddChildAccount add sub-address by coin id, auto set tag for created addresses
func (*AccManageService) AddChildAccount(coinId, amount int, tag string, adds, deeps []string) int {
	var address Address
	var coin Coin
	var coinName string
	//var data string
	valueStrings := make([]string, 0, amount)
	valueArgs := make([]interface{}, 0, amount*6)
	//查询是否已存在别名并获取tagindex
	num := rdb.Last(&address, "coinId = ? and tag = ?", coinId, tag).RowsAffected
	number := rdb.First(&coin, "id = ? ", coinId).RowsAffected
	if number == 1 {
		coinName = coin.Name
	} else {
		return 0
	}
	if num == 1 {
		//已存在别名 延用tagIndex
		tagIndex := address.TagIndex
		for i := 0; i < amount; i++ {
			valueStrings = append(valueStrings, "(?,?,?,?,?,?,?)")
			valueArgs = append(valueArgs, coinName)
			valueArgs = append(valueArgs, coinId)
			valueArgs = append(valueArgs, tag)
			valueArgs = append(valueArgs, adds[i])
			valueArgs = append(valueArgs, 1)
			valueArgs = append(valueArgs, tagIndex+i+1)
			valueArgs = append(valueArgs, deeps[i])
		}
	} else if num == 0 {
		for i := 0; i < amount; i++ {
			valueStrings = append(valueStrings, "(?,?,?,?,?,?,?)")
			valueArgs = append(valueArgs, coinName)
			valueArgs = append(valueArgs, coinId)
			valueArgs = append(valueArgs, tag)
			valueArgs = append(valueArgs, adds[i])
			valueArgs = append(valueArgs, 1)
			valueArgs = append(valueArgs, i+1)
			valueArgs = append(valueArgs, deeps[i])
		}
	} else {
		return 0
	}
	tx := rdb.Begin()
	stmt := fmt.Sprintf("insert into address (coinName,coinId,tag,address,type,tagIndex,deep) VALUES %s", strings.Join(valueStrings, ","))
	err := tx.Exec(stmt, valueArgs...).Error
	if err != nil {
		tx.Rollback()
		log.Error("添加失败", err)
		return 0
	}
	tx.Commit()
	return 1
}

type EthBalance struct {
	Balance string
	Address string
}

// MergeAccount transfer asset from sub-address to main address, param mold defines the payment account
// mold:
//	0: all the sub-address
// 	1: the sub-address specified
type transferMsgJson struct {
	FromAddress AddressJson
	ToAddress   AddressJson
	Token       string
}

func (*AccManageService) MergeAccount(coinId, mold int, ids []int) int {
	var main []Address
	var accounts []MergeAccount
	var err error
	//查询主账户
	//rdb.Table("coin").Select("address.coinId,address.address,address.id,address.deep").Joins("join address on coin.id = address.coinId").Where("coin.id = ? and address.type = 0 and address.isDeleted = 0", coinId).Scan(&main)
	err = rdb.Raw("select * from address where type = 0 and coinName = ?", "ETH").Scan(&main).Error
	if err != nil || (len(main) == 0) { //未查询到主账户
		return common.MergeAccountAddressNotFound
	}
	tx := rdb.Begin()
	//取出每个地址对应的eth余额
	var ethBalance []EthBalance
	var ethAddressMap = make(map[string]string)
	err = rdb.Raw("select cap.balance as balance,addr.address as address from capital cap left join address addr on addr.id = cap.addressId left join coin c on c.id = cap.coinId where c.name = 'ETH'  and addr.type = 1 and cap.balance > ?", "0").Scan(&ethBalance).Error
	if len(ethBalance) == 0 {
		log.Error("没有可用的子账户")
		return common.MergeAccountNoAddrToMegge
	}
	for i := 0; i < len(ethBalance); i++ {
		ethAddressMap[ethBalance[i].Address] = ethBalance[i].Balance
	}

	if mold == 0 {
		//合并所有子账号
		//查询出所有子账号
		err = rdb.Raw("select cap.coinName as coinName,cap.balance as balance,addr.address as fromAddr,addr.deep,c.tokenAddress from capital cap "+
			"left join address addr on addr.id = cap.addressId "+
			"left join coin c on c.id = cap.coinId where cap.coinId = ? and cap.balance>0 and addr.type = 1", coinId).Scan(&accounts).Error
	} else {
		err = rdb.Raw("select cap.coinName as coinName,cap.balance as balance,addr.address as fromAddr,addr.deep,c.tokenAddress from capital cap "+
			"left join address addr on addr.id = cap.addressId "+
			"left join coin c on c.id = cap.coinId where cap.coinId = ? and cap.balance>0 and addr.type = 1 and cap.addressId in (?)", coinId, ids).Scan(&accounts).Error
	}

	if err != nil {
		log.Error("没有查询到子账户", err)
		return common.MergeAccountNoAddrToMegge //没有可合并的子账户
	}
	if len(accounts) == 0 {
		return common.MergeAccountNoAddrToMegge //没有可合并的子账户
	}
	valueStrings := make([]string, 0, len(accounts))
	valueArgs := make([]interface{}, 0, len(accounts)*7)
	//可用数量
	use := 0
	//合并的订单
	orderId := uuid.Must(uuid.NewV4()).String()
	for _, account := range accounts {
		if ethAddressMap[account.FromAddr] == "" {
			continue
		}
		if account.CoinName == "ETH" {

			decimalBalance, error := decimal.NewFromString(account.Balance)
			if error != nil {
				log.Error("转数字失败", error)
				continue
			}
			aa, _ := decimal.NewFromString(minBalance)
			if decimalBalance.LessThanOrEqual(aa) {
				continue
			}
			bb := decimalBalance.Sub(aa)
			account.Balance = bb.String()
		}
		transferId := uuid.Must(uuid.NewV4()).String()
		valueStrings = append(valueStrings, "(?,?,?,?,?,?,?,?,?,?)")
		valueArgs = append(valueArgs, transferId)
		valueArgs = append(valueArgs, orderId)
		valueArgs = append(valueArgs, account.CoinName)
		valueArgs = append(valueArgs, account.Balance)
		valueArgs = append(valueArgs, main[0].Address)
		valueArgs = append(valueArgs, 2)

		var msgJson transferMsgJson
		var uintDeep []uint32
		err = json.Unmarshal([]byte(account.Deep), &uintDeep)
		if err != nil {
			log.Error("deep解析错误", account.FromAddr)
			return common.MergeAccountSqlErr
		}

		msgJson.FromAddress = AddressJson{
			Address: account.FromAddr,
			Deep:    uintDeep,
		}
		msgJson.ToAddress = AddressJson{
			Address: main[0].Address,
			Deep:    nil,
		}
		msgJson.Token = account.TokenAddress
		json, _ := json.Marshal(msgJson)
		valueArgs = append(valueArgs, string(json))
		valueArgs = append(valueArgs, 1)
		valueArgs = append(valueArgs, account.FromAddr)
		valueArgs = append(valueArgs, "向总账户转账")
		use++
		transferSend := TransferSend{
			TransferId: transferId,
			CoinName:   account.CoinName,
			Amount:     account.Balance,
			ToAddress:  main[0].Address,
			Types:      2,
			AddressMsg: string(json),
			Status:     1, //1审批成功，转账中
			OrderId:    orderId,
		}
		MergeETHTransferCh <- transferSend
	}
	if use == 0 {
		log.Error("没有可用账户")
		return common.MergeAccountNoAddrToMegge
	}
	stmt := fmt.Sprintf("insert into transfer (`id`,`orderId`,`coinName`,`amount`,`toAddress`,`types`,`msg`,`status`,`fromAddress`,`applyReason`) VALUES %s", strings.Join(valueStrings, ","))
	err = tx.Exec(stmt, valueArgs...).Error
	if err != nil {
		tx.Rollback()
		log.Error("插入记录失败", err)
		return common.MergeAccountSqlErr
	}

	tx.Commit()

	return common.MergeAccountSuccess
}

type addrBaseInfo struct {
	Address string `json:"address"`
	Balance string `json:"balance"`
	Type    int    `json:"type"`
}
type output struct {
	Type          int          `json:"type"`
	Result        addrBaseInfo `json:"result"`
	Sign          string       `json:"sign"`
	Account       string       `json:"account"`
	MasterAddress string       `json:"masterAddress"`
}

// ContractAddr get contract address from voucher
const (
	needToCreateContract int = iota
	hasNoContract
	contractSaved
)

func (*AccManageService) ContractAddr() (*output, error) {
	// 获取币种信息
	coinInfo, err := getCoinInfoByName("ETH")
	if err != nil {
		log.Error("ETH NOT FOUND")
		return nil, errors.New("Eth Not Found.")
	}
	// 获取主账户地址和合约地址
	var data addrBaseInfo
	var result output
	// 未申请创建合约地址
	result.Type = needToCreateContract
	var count_contract int
	// 从签名机获取合约地址并保存
	err = updateContractAddr(coinInfo)
	if err != nil {
		log.Error("存储合约地址", err)
		return nil, err
	}

	rows, err := rdb.Table("address").Select("address.address, address.type, ifnull(capital.balance, 0)").Joins("left join capital on capital.addressId = address.id and capital.coinName = address.coinName").Where(&Address{CoinName: "ETH", IsDeleted: 0}).Where("type in (?)", []int{0, 2}).Rows()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Debug("contract address not found")
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		r := addrBaseInfo{}
		rows.Scan(&r.Address, &r.Type, &r.Balance)
		if r.Type == common.ContractAddressType {
			count_contract++
			if r.Address == "" {
				result.Type = hasNoContract // 合约未生成

			} else {
				result.Type = contractSaved
			}
		}
		if r.Type == common.MainAddressType {
			data.Address = r.Address
			data.Balance = r.Balance
			data.Type = r.Type
		}
	}
	// 合约地址数目不合法
	if count_contract > 1 {
		log.Debug("Too Many Contract Address Count(contractAddr) = %v, Count = %v", count_contract, data)
		return nil, errors.New("Too Many Contract Address")
	}
	result.Result = data
	return &result, nil
}

// MainAccAssets get the assets info of main address
func (*AccManageService) MainAccAssets() (string, error) {
	var asset string
	row := rdb.Table("address").Select("ifnull(c.balance, 0)").Joins("left join capital as c on c.addressId = address.id and c.coinName = address.coinName").Where("c.CoinName = ?", "ETH").Where(" address.type = ?", common.MainAddressType).Row()
	err := row.Scan(&asset)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Debug("No Main Address Yet")
			return "0", nil
		}
		return "", err
	}
	return asset, nil
}

func getCoinInfoByName(name string) (*Coin, error) {
	var coinInfo Coin
	err := rdb.Where(&Coin{Name: name}).First(&coinInfo).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Error("ETH not found")
			return nil, errors.New("ETH NOT FOUND")
		}
		return nil, err
	}
	return &coinInfo, nil
}

func updateContractAddr(coinInfo *Coin) error {
	// 从签名机获取合约地址
	vres, _, voucherRet := voucher.SendVoucherData(&voucher.GrpcServer{Type: voucher.VOUCHER_STATUS})
	log.Debug("voucher return ", voucherRet)
	if voucherRet == voucher.VRET_ERR || voucherRet == voucher.VRET_TIMEOUT {
		log.Error("Voucher Return Fail", voucherRet)
		return errors.New("Voucher Return Fail")
	}
	if voucherRet == voucher.VRET_STATUS {
		if vres == nil {
			log.Error("Get Contract Address Error")
			return errors.New("Voucher Return Fail")
		}
	}
	var addr Address
	err := rdb.Where(Address{Type: common.ContractAddressType}).First(&addr).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}

	if vres.ContractAddress != "" && addr.Address == "" {
		// 存储合约地址
		err := rdb.Model(&Address{}).Where(&Address{Type: common.ContractAddressType, CoinId: coinInfo.ID}).Update("address", vres.ContractAddress).Error
		if err != nil {
			log.Error("存储合约地址", err)
			return errors.New("Voucher Return Fail")
		}
		return nil
	}
	return nil
}


// ApplyForContract When creating thetemplate flow for the first time, if there is no contract, you should apply to create
// contract from voucher.
func (*AccManageService) ApplyForContract() error {
	// 获取币种信息
	coinInfo, err := getCoinInfoByName("ETH")
	if err != nil {
		log.Error("ETH NOT FOUND")
		return err
	}
	oper := &voucher.GrpcServer{
		Type: voucher.VOUCHER_OPERATE_DEPLOY,
	}
	_, cr, voucherRet := voucher.SendVoucherData(oper)
	if voucherRet == voucher.VRET_ERR || voucherRet == voucher.VRET_TIMEOUT {
		log.Error("Voucher Return Fail", voucherRet)
		return errors.New("Voucher Return Error")
	}
	if voucherRet == voucher.VRET_CLIENT {
		//if cr.Status != 0 && cr.Status != voucher.STATUS_NO_CONTRACT {
		if cr.Status != 0 && cr.Status != voucher.STATUS_CONTRACT_ADDR_PENDING {
			log.Error("创建合约失败")
			return errors.New("Failed To Create Contract")
		}
	}
	// 是否已创建过合约地址
	var addr Address
	err = rdb.Where(Address{Type: common.ContractAddressType}).First(&addr).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 初始化空合约地址
			err = initEmptyContractAddr(coinInfo)
			if err != nil {
				log.Error("初始化空合约地址")
				return err
			}
		} else {
			log.Error("获取合约地址", err)
			return err
		}
	}
	if addr.Type == common.ContractAddressType {
		// 已创建过合约地址
		log.Debug("已创建过合约地址")
		return errors.New("Duplicate Init Contract")
	}

	return nil
}

func initEmptyContractAddr(coinInfo *Coin) error {
	err := rdb.Create(&Address{Address: "", Type: common.ContractAddressType, CoinId: coinInfo.ID, CoinName: coinInfo.Name, Tag: "", TagIndex: 0, IsDeleted: 0}).Error
	if err != nil {
		log.Error("c初始化合约地址合约地址", err)
		return err
	}
	return nil
}

// GetMasterAddress get the main address infos of all coin
func (*AccManageService) GetMasterAddress() []Address {
	var address []Address
	err := rdb.Raw("select * from address where type = 0 union (select a.* from address a left join capital cap on a.id = cap.addressId where cap.balance >'0')").Scan(&address).Error
	//err := rdb.Where("type = 0").Find(&address).Error
	if err != nil {
		log.Error("获取所有address失败")
		return nil
	}
	return address
}

// SaveMasterAddress save main address to db
func (*AccManageService) SaveMasterAddress(address *Address) bool {

	err := rdb.Create(address).Error
	if err != nil {
		log.Error("保存address失败", err)
		return false
	}
	return true
}

// FindAddressByAddress get address detail info by address
func (*AccManageService) FindAddressByAddress(addr string, coinName string) bool {
	var address Address
	err := rdb.Where("address = ? and coinName = ?", addr, coinName).First(&address).Error
	if err != nil {
		log.Error("查询地址失败", err)
		return false
	}
	return true
}

// FindAddressByBTCAddress get the BTC and LTC address info by address
func (*AccManageService) FindAddressByBTCAddress(addr []string, coinId int) ([]Address, bool) {
	var address []Address
	err := rdb.Raw("SELECT * FROM `address`  WHERE address in (?) and coinId = ?", addr, coinId).Scan(&address).Error
	if err != nil {
		log.Error("查询地址失败", err)
		return address, false
	}
	return address, true
}

// UpdateBalance update the balance by address
func (*AccManageService) UpdateBalance(capital *Capital) bool {
	tx := rdb.Begin()
	err := tx.Exec("delete from capital where addressId = ? and coinId = ?", capital.AddressID, capital.CoinId).Error
	if err != nil {
		log.Error("更新余额", err)
		tx.Rollback()
		return false
	}
	err = tx.Create(capital).Error

	if err != nil {
		log.Error("更新余额", err)
		tx.Rollback()
		return false
	}
	tx.Commit()
	return true

}
