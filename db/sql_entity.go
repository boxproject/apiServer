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
	"database/sql"
	"time"
	"sync"
)

// account
type Account struct {
	ID           int       `gorm:"AUTO_INCREMENT"`
	AppId        string    `gorm:"column:appId"`
	Name         string    `gorm:"column:name"`
	Pwd          string    `gorm:"column:pwd"`
	Salt         string    `gorm:"column:salt"`
	PubKey       string    `gorm:"column:pubKey"`
	IsDeleted    int       `gorm:"column:isDeleted"`
	Msg          string    `gorm:"column:msg"`
	DepartmentId int       `gorm:"column:departmentId"`
	UserType     int       `gorm:"column:userType"`
	Frozen       int       `gorm:"column:frozen"`
	Attempts     int       `gorm:"column:attempts"`
	FrozenTo     time.Time `gorm:"column:frozenTo"`
	SourceAppId  string    `gorm:"column:sourceAppId"`
	Level        int       `gorm:"column:level"`
	CreatedAt    time.Time `gorm:"column:createdAt"`
	UpdatedAt    time.Time `gorm:"column:updatedAt"`
}

func (dep Account) TableName() string {
	return "account"
}

// address
type Address struct {
	ID        int
	Address   string
	Type      int
	CoinId    int `gorm:"column:coinId"`
	Tag       string
	TagIndex  int    `gorm:"column:tagIndex"`
	IsDeleted int    `gorm:"column:isDeleted"`
	CoinName  string `gorm:"column:coinName"`
	UsedBy    int    `gorm:"column:usedBy"`
	Deep      string `gorm:"column:deep"`
}

func (addr Address) TableName() string {
	return "address"
}

// auth
type Auth struct {
	ID       int    `gorm:"AUTO_INCREMENT"`
	Name     string `gorm:"column:name"`
	AuthType string `gorm:"column:authType"`
}

func (auth Auth) TableName() string {
	return "auth"
}

// authmap
type AuthMap struct {
	ID        int `gorm:"AUTO_INCREMENT"`
	AccountId int `gorm:"column:accountId"`
	AuthId    int `gorm:"column:authId"`
}

func (authMap AuthMap) TableName() string {
	return "authmap"
}

// capital
type Capital struct {
	ID        int
	AddressID int    `gorm:"column:addressId"`
	CoinName  string `gorm:"column:coinName"`
	Balance   string
	CoinId    int       `gorm:"column:coinId"`
	CreatedAt time.Time `gorm:"column:createdAt"`
	Address   string    `gorm:"column:address"`
}

func (capital Capital) TableName() string {
	return "capital"
}

// coin
type Coin struct {
	ID           int    `gorm:"AUTO_INCREMENT"`
	Name         string `gorm:"column:name"`
	FullName     string `gorm:"column:fullName"`
	Precise      int    `gorm:"column:precise"`
	Balance      string `gorm:"column:balance"`
	TokenType    int    `gorm:"column:tokenType"`
	TokenAddress string `gorm:"column:tokenAddress"`
	Available    int    `gorm:"column:available"`
}

func (dep Coin) TableName() string {
	return "coin"
}

// department
type Department struct {
	ID           int
	Name         string
	CreatorAccId int `gorm:"column:creatorAccId"`
	Available    int
	Order        int       `gorm:"column:order"`
	CreatedAt    time.Time `gorm:"column:createdAt"`
	UpdatedAt    time.Time `gorm:"column:updatedAt"`
}

func (dep Department) TableName() string {
	return "department"
}

// deposit
type Deposit struct {
	ID           string
	CoinId       int    `gorm:"column:coinId"`
	CoinName     string `gorm:"column:coinName"`
	FromAddr     string `gorm:"column:fromAddr"`
	ToAddr       string `gorm:"column:toAddr"`
	Amount       string
	TxId         string `gorm:"column:txId"`
	BlockNum     string `gorm:"column:blockNumber"`
	Confirm      int
	CreatedAt    time.Time `gorm:"column:createdAt"`
	IsUpdate     int       `gorm:"column:isUpdate"`
	TokenAddress string    `gorm:"column:tokenAddress"`
	Precise      int       `gorm:"column:precise"`
	UpdatedAt    time.Time `gorm:"column:updatedAt"`
}

func (deposit Deposit) TableName() string {
	return "deposit"
}

// registration
type Registration struct {
	ID            int       `gorm:"AUTO_INCREMENT"`
	Name          string    `gorm:"column:name"`
	Pwd           string    `gorm:"column:pwd"`
	Salt          string    `gorm:"column:salt"`
	SourceAppId   string    `gorm:"column:sourceAppId"`
	SourceAccount string    `gorm:"column:sourceAccount"`
	AppId         string    `gorm:"column:appId"`
	PubKey        string    `gorm:"column:pubKey"`
	Status        int       `gorm:"column:status"`
	Msg           string    `gorm:"column:msg"`
	CreatedAt     time.Time `gorm:"column:createdAt"`
	UpdatedAt     time.Time `gorm:"column:updatedAt"`
	RefuseReason  string    `gorm:"column:refuseReason"`
	Level         int       `gorm:"column:level"`
}

func (dep Registration) TableName() string {
	return "registration"
}

// template
type Template struct {
	ID        string `gorm:"column:id"`
	Hash      string
	Name      string
	CreatorID int `gorm:"column:creatorId"`
	Content   string
	Status    int
	Period    int
	CreatedAt time.Time `gorm:"column:createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt"`
}

func (template Template) TableName() string {
	return "template"
}

// templateOper
type TemplateOper struct {
	ID         string    `gorm:"column:id"`
	TemplateId string    `gorm:"column:templateId"`
	AccountId  int       `gorm:"column:accountId"`
	Status     int       `gorm:"column:status"`
	CreatedAt  time.Time `gorm:"column:createdAt"`
	AppId      string    `gorm:"column:appId"`
	AppName    string    `gorm:"column:appName"`
	Sign       string    `gorm:"column:sign"`
	Reason     string    `gorm:"column:reason"`
}

func (templateOper TemplateOper) TableName() string {
	return "templateOper"
}

// templateView
type TemplateView struct {
	TemplateId  string `gorm:"column:templateId"`
	CoinId      int    `gorm:"column:coinId"`
	ReferAmount string `gorm:"column:referAmount"`
	AmountLimit string `gorm:"column:amountLimit"`
	Referspan   int
	FrozenTo    time.Time `gorm:"column:frozenTo"`
}

func (tv TemplateView) TableName() string {
	return "templateView"
}

// transfer
type Transfer struct {
	ID          string
	CoinName    string `gorm:"column:coinName"`
	Amount      string
	Status      int
	TxId        string `gorm:"column:txId"`
	ToAddr      string `gorm:"column:toAddress"`
	FromAddr    string `gorm:"column:fromAddress"`
	Accepted    int
	OrderId     string `gorm:"column:orderId"`
	Tag         string
	CreatedAt   time.Time `gorm:"column:createdAt"`
	UpdatedAt   time.Time `gorm:"column:updatedAt"`
	ApplyReason string    `gorm:"column:applyReason"`
	AmountIndex int       `gorm:"column:amountIndex"`
	Types       int       `gorm:"column:types"`
	Msg         string    `gorm:"column:msg"`
	FromAddress string    `gorm:"column:fromAddress"`
}

func (tx Transfer) TableName() string {
	return "transfer"
}

// transferReview
type TransferReview struct {
	ID          int
	OrderNum    string `gorm:"column:orderNum"`
	Status      int
	Encode      string
	Reason      string
	UpdatedAt   time.Time `gorm:"column:updatedAt"`
	AccountName string    `gorm:"column:accountName"`
	Level       int       `gorm:"column:level"`
}

func (tf TransferReview) TableName() string {
	return "transferReview"
}

// transferOrder
type TransferOrder struct {
	ID            string         `gorm:"column:id"`
	CoinName      string         `gorm:"column:coinName"`
	Hash          string         `gorm:"column:hash"`
	ApplyReason   string         `gorm:"column:applyReason"`
	Deadline      sql.NullString `gorm:"column:deadline"`
	Miner         string         `gorm:"column:miner"`
	Sign          string         `gorm:"column:sign"`
	CreatedAt     time.Time      `gorm:"column:createdAt"`
	UpdatedAt     time.Time      `gorm:"column:updatedAt"`
	Status        int            `gorm:"column:status"`
	Content       string         `gorm:"column:content"`
	ApplyerId     int            `gorm:"column:applyerId"`
	Amount        string         `gorm:"column:amount"`
	NowLevel      int            `gorm:"column:nowLevel"`
	ApproversSign string         `gorm:"column:approversSign"`
}

func (to TransferOrder) TableName() string {
	return "transferOrder"
}

type OrderDetail struct {
	FullName        string    `gorm:"column:fullName"`
	CoinName        string    `gorm:"column:coinName"`
	ApplyerName     string    `gorm:"column:applyerName"`
	AppId           string    `gorm:"column:appId"`
	ApplyReason     string    `gorm:"column:applyReason"`
	Deadline        time.Time `gorm:"column:deadline"`
	CreatedAt       time.Time `gorm:"column:createdAt"`
	OrderAddress    []interface{}
	StrTemplate     string `gorm:"column:strtemplate"`
	Amount          string `gorm:"column:amount"`
	Status          int    `gorm:"column:status"`
	NowLevel        int    `gorm:"column:nowLevel"`
	Msg             string `gorm:"column:content"`
	ApprovalContent map[string]interface{}
	Template        TemplateContent
	OrderId         string `gorm:"column:orderId"`
	ApplySign       string `gorm:"column:applySign"`
}

type OwnerOper struct {
	ID        int       `gorm:"column:id"`
	Account   string    `gorm:"column:account"`
	OpAccount string    `gorm:"column:operatedAccount"`
	Status    int       `gorm:"column:status"`
	OpAppId   string    `gorm:"column:operatedAppid"`
	RegId     int       `gorm:"column:regId"`
	UpdatedAt time.Time `gorm:"column:updatedAt"`
}

func (owner OwnerOper) TableName() string {
	return "ownerOper"
}

type OwnerReg struct {
	ID           int       `gorm:"AUTO_INCREMENT"`
	Name         string    `gorm:"column:name"`
	Pwd          string    `gorm:"column:pwd"`
	Salt         string    `gorm:"column:salt"`
	AppId        string    `gorm:"column:appId"`
	PubKey       string    `gorm:"column:pubKey"`
	Status       int       `gorm:"column:status"`
	CreatedAt    time.Time `gorm:"column:createdAt"`
	UpdatedAt    time.Time `gorm:"column:updatedAt"`
	RefuseReason string    `gorm:"column:refuseReason"`
	Msg          string    `gorm:"column:msg"`
}

func (ownerReg OwnerReg) TableName() string {
	return "ownerReg"
}

// configs
type Configs struct {
	ID    int
	Key   string `gorm:"column:con_key"`
	Value string `gorm:"column:con_value"`
}

func (con Configs) TableName() string {
	return "configs"
}

// 合并账户--------------------
type MergeAccount struct {
	ToAddr       string
	FromAddr     string `gorm:"column:fromAddr"`
	CoinName     string `gorm:"column:coinName"`
	Balance      string `gorm:"column:balance"`
	Deep         string `gorm:"column:deep"`
	TokenAddress string `gorm:"column:tokenAddress"`
}

//----------------------

type Logger struct {
	ID        int       `gorm:"AUTO_INCREMENT"`
	Detail    string    `gorm:"column:detail"`
	Note      string    `gorm:"column:note"`
	Operator  string    `gorm:"column:operator"`
	LogType   string    `gorm:"column:logType"`
	CreatedAt time.Time `gorm:"column:createdAt"`
	Pos       string    `gorm:"column:pos"`
}

func (logger Logger) TableName() string {
	return "log"
}

type Message struct {
	ID        int       `gorm:"AUTO_INCREMENT"`
	Title     string    `gorm:"column:title"`
	CreatedAt time.Time `gorm:"column:createdAt"`
	Content   string    `gorm:"column:content"`
	Padding   string    `gorm:"column:padding"`
	Receiver  string    `gorm:"column:receiver"`
	Reader    string    `gorm:"column:reader"`
	Type      int       `gorm:"column:type"`
	WarnType  int       `gorm:"column:warnType"`
	Param     string    `gorm:"column:param"`
}

func (Message Message) TableName() string {
	return "message"
}

//wallet record

type TxInfos struct {
	Txid     string `gorm:"column:txId"`
	Confirm  int    `gorm:"column:confirm"`
	Token    string `gorm:"column:token"`
	Type     int    `gorm:"column:type"`
	Height   int    `gorm:"column:height"`
	Fee      string `gorm:"column:fee"`
	TxObj    string `gorm:"column:tx_obj"`
	ExtValid int    `gorm:"column:ext_valid"`

	CreateAt   time.Time `gorm:"column:created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at"`
	DeleteddAt time.Time `gorm:"column:deletedd_at"`
}

func (txInfos TxInfos) TableName() string {
	return "tx_infos"
}

type AddressSigns struct {
	Account     string        `json:"account"`
	AddressInfo []AddressSign `json:"addressSign"`
}
type AddressSign struct {
	CoinName      string `json:"coinName"`
	MasterAddress string `json:"masterAddress"`
	Sign          string `json:"sign"`
}

//需要上链的审批流

type TemplateChain struct {
	Lock           *sync.Mutex
	TemplateChains []Template
}

var TemplateChains = TemplateChain{
	Lock:           new(sync.Mutex),
	TemplateChains: []Template{},
}

func (this *TemplateChain) SetTemplateChain(tem Template) {
	this.Lock.Lock()
	defer this.Lock.Unlock()
	this.TemplateChains = append(this.TemplateChains, tem)
}

func (this *TemplateChain) GetTemplateChain() []Template {
	this.Lock.Lock()
	defer this.Lock.Unlock()
	return this.TemplateChains
}

//删除第一个元素
func (this *TemplateChain) DelTemplateChain() {
	this.Lock.Lock()
	defer this.Lock.Unlock()
	this.TemplateChains = append(this.TemplateChains[:0], this.TemplateChains[1:]...)
}

type WebTransfers struct {
	ID          int
	TransferId  string `gorm:"column:transferId"`
	Msg         string `gorm:"column:msg"`
	CreatedTime int64  `gorm:"column:createdTime"`
}

func (webTransfers WebTransfers) TableName() string {
	return "webTransfers"
}
