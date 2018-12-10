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
package rpc

//  server to voucher
type GrpcServer struct {
	Type      string //req type
	AppId     string //app id
	BakAction string //backup action
	AppName   string //app name
	AesKey    []byte //aes key
	Msg       []byte
	Sign      []byte //app msg sign
	Timestamp int64  //Timestamp,时间戳
	Hash      string //hash
	FlowInfo  []byte //flow info
	FlowSigns []byte
	//转账专用
	Currency string //Currency 币种
	AddrMap  []byte // addr map
	TxInfo   []byte //trans info no sign
	TransMsg []byte //trans info
	Other    []byte //other otherMapInfo["appPubKey"] otherMapInfo["voucherSign"]
}

type TransferInfo struct {
	TransferMsgs   string //transfer message  【交易信息,json序列化(TransferMsg{})】
	AmountIndex    int    //转账金额索引 需要添加
	ToAddress      string
	ApplyAccount   string //apply account   【申请人账户】
	ApplyPublickey string //apply public key  【申请人公钥】
	ApplySign      string //apply sign    【申请人签名】
	ApproversSign  string //Approvers sign   【审批人签名，json序列化(map[string(appName)]string(sign))】
}

type AppMsg struct {
	Password  []byte //password
	BakPasswd []byte //backup password
	PublicKey []byte //app public key
}

type FlowSign struct {
	AppId string
	Sign  string
}

//  client to server
type GrpcClient struct {
	Type         string //req type
	AppId        string //app id
	AppName      string //app name
	AesKey       string //aes key
	PublicKeys   string //app public key;map[appId]pubKeyBytes
	Sign         []byte //app msg sign
	ContractAddr string //contract address
	D            string //voucher random D value
	Status       int    //server req status
	SignTx       []byte //signed trans
	Other        []byte //other
	Timestamp    int64  //Timestamp,时间戳
	Logs         []byte // operate logs ,OperateLog
}

type OperateLog struct {
	Timestamp   int64  //time stamp   【时间戳】
	AppId       string //app id    【app ID】
	AppName     string //app name    【app ID】
	OperateType string //app operate type 【操作类型】
	Result      string //app operate result 【操作结果】
	Other       string //other message  【其他】
}

//签名机状态
type VoucherStatus struct {
	ServerStatus    int                         //系统状态
	Total           int                         //密钥数量
	BakDevStatus    int                         //备份设备状态
	BakStatus       map[string]int              //备份状态(hardware backup status) [appId]bakTimes
	EthNetState     int                         //eth网络状态
	BaseAddress     string                      //主地址
	ContractAddress string                      //合约地址
	NodesAuthorized map[string][]NodeAuthorized //授权情况[req type]{[]user info} ,[请求类型(create,start,reback)]{用户信息}
	KeyStoreStatus  map[string]KeyStoreStatu    //公约添加状态
	Timestamp       int64                       //Timestamp,时间戳
}

type NodeAuthorized struct {
	ApplyerId  string
	Authorized bool
}

type KeyStoreStatu struct {
	ApplyerId   string
	ApplyerName string
}

const (
	VOUCHER_STATUS                     = "status"
	VOUCHER_OPERATE_ADDKEY             = "0"  // add app pub key
	VOUCHER_OPERATE_CREATE             = "1"  // create
	VOUCHER_OPERATE_DEPLOY             = "2"  // deploy contract
	VOUCHER_OPERATE_START              = "3"  // start
	VOUCHER_OPERATE_STOP               = "4"  // stop
	VOUCHER_OPERATE_APP_PUBKEY         = "5"  //get app pub key
	VOUCHER_OPERATE_HDBAKUP            = "6"  //private key backup
	VOUCHER_OPERATE_SIGN               = "7"  //trans sign//构造交易签名
	VOUCHER_OPERATE_CHECK_PASS         = "8"  // app pass check
	VOUCHER_OPERATE_MASTER_PUBKEY      = "9"  // hd wallet master pub key
	VOUCHER_OPERATE_RECOVER            = "10" //app pub key recovery
	VOUCHER_OPERATE_ADDHASH            = "11" // add flow hash
	VOUCHER_OPERATE_DISHASH            = "12" // disable flow hash
	VOUCHER_OPERATE_HASH_AVAILABLE     = "13"
	VOUCHER_OPERATE_TRANSSIGN_INTERIOR = "14"
	VOUCHER_OPERATE_KEY_EXCHANGE       = "15" // voucher key exchange  【密钥交换】
	HANDLE_SHAKE                       = "HandleShake"
)

// 返回状态status
const (
	STATUS_ERROR int = iota - 1
	STATUS_OK
	STATUS_OTHER
)

const (
	STATUS_APP_LOSEBASEINFO          int = iota + 100 //lose base info 				【缺失必要信息】
	STATUS_APP_RSA_SIGN_ERROR                         //rsa sign error 				【rsa签名失败】
	STATUS_APP_VERIFY_ERROR                           //Verify error 				【验证签名错误】
	STATUS_APP_RSA_DECRYPT_ERROR                      //rsa Decrypt error			【rsa解密错误】
	STATUS_APP_AES_DECRYPT_ERROR                      //aes Decrypt error 			【aes解密错误】
	STATUS_APP_PASSWORD_ERROR                         //password error				【密码错误】
	STATUS_APP_NAME_NOTMATCH                          //app name not math			【app名称不匹配】
	STATUS_APP_SAMENAME                               //same app name				【相同的app名称】
	STATUS_APP_RECOVERY_SUCCESSS                      //app pubKey recovery success	【app公钥恢复成功】
	STATUS_APP_HDBOX_BAKUPFAILED                      //hardware box backup failed  【box硬件备份失败】
	STATUS_APP_HDBOX_DEVLOCKED                        //hardware box device locked  【box硬件锁定】
	STATUS_APP_HDBOX_NOBAKDEV                         //hardware box not find device【box硬件未发现】
	STATUS_APP_HDBOX_BLOCK_NO_ENOUGH                  //hardware box block not enough【box硬件存储空间不足】
)

// contract and flow hash status
const (
	STATUS_NO_CONTRACT           int = iota + 200 //contract address not have	【没有合约地址】
	STATUS_CONTRACT_ADDR_PENDING                  //contract address pending	【合约地址上链中】
	STATUS_INSUFFICIENT_FUNDS                     //insufficient funds 			【手续费不足】
	STATUS_FLOW_VERIFY_ERROR                      //flow Verify error  			【审批流验证签名错误】
	STATUS_HASH_LOSEBASEINFO                      //lose base info 				【审批流hash缺失必要信息】
	STATUS_HASH_AVAILABLE                         //hash Available  			【审批流hash可用】
	STATUS_HASH_UNAVAILABLE                       //hash unavailable  			【审批流hash不可用】
	STATUS_HASH_PENDING                           //hash pending  				【审批流hash pending】
)

// trans status
const (
	STATUS_TRANS_LOSEBASEINFO        int = iota + 300 //lose base info 				【缺失必要信息】
	STATUS_TRANS_VERIFY_ERROR                         //trans Verify error  		【交易验证签名错误】
	STATUS_TRANS_SIGN_ERROR                           //trans sign error 			【交易签名失败】
	STATUS_TRANS_UNKNOW_CURRENCY                      //unknow Currency 			【未知币种】
	STATUS_TRANS_CURRENCY_NOTEQUAL                    //Currency not equal 			【币种不匹配】
	STATUS_TRANS_MARSHAL_ERROR                        //trans marshal error 		【交易marshal错误】
	STATUS_TRANS_GET_INFO_ERROR                       //get trans into error		【获取交易信息错误】
	STATUS_TRANS_TOADDR_NOT_INTERIOR                  //to address not Interior		【转出地址为非内部地址】
	STATUS_TRANS_SIGN_REPEAT                          //trans sign repeat 			【交易签名重复】
)

//签名机系统状态
const (
	VOUCHER_STATUS_UNCONNETED = 0 //未连接
	VOUCHER_STATUS_UNCREATED  = 1 //未创建
	VOUCHER_STATUS_CREATED    = 2 //已创建
	VOUCHER_STATUS_BAKUP      = 3 //已备份
	VOUCHER_STATUS_STARTED    = 4 //已启动
	VOUCHER_STATUS_PAUSED     = 5 //已停止
)

// 签名机连接状态
const (
	VRET_ERR     = 0
	VRET_STATUS  = 1
	VRET_CLIENT  = 2
	VRET_TIMEOUT = 4
)
