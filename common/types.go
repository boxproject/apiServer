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
package common

// file types define
const (
	YML_FILE 	= "yml"
	TOML_FILE 	= "toml"
	JSON_FILE 	= "json"
	XML_FILE 	= "xml"
	HeaderLangKey = "content-language"
	HeaderZhLang = "zh-Hans"
	HeaderEnLang = "en"
)

const (
	LangZhType = "zh"
	LangEnType = "en"
)

// apiServer config path
const CONF_PATH = "/config/"

// voucher key file
const VOUCHER_KEY_FILE = "voucher.key"

// config file name
const CONF_FILE = "config.toml"

// The configuration file for log4go
const LOG_CONF = "log.xml"

// max password attempts
const MAX_PASWD_ATTEMPTS = 5

// template period Upper limit
const TemplateExpireUpperLimit = 240

// MergeAccount returns
const (
	MergeAccountAddressNotFound int = iota
	MergeAccountSuccess
	MergeAccountOther
	MergeAccountTransfering
	MergeAccountNoAddrToMegge
	MergeAccountSqlErr
)

// get coin list, list type
const (
	CoinListMasterCoinType int = iota
	CoinListTokenType
	CoinListAllCoinType
	CoinListETHCoinType
	CoinListAvailableCoinType
)

// address type in db
const (
	MainAddressType int = iota
	NormalAddressType
	ContractAddressType
)

// template status
const (
	TemplatePending    int = iota // 0 waiting to approval
	TemplateApprovaled            // 1 approved
	TemplateRejected              // 2 reject
	TemplateToDisable             // 3 to be disable
	TemplateAvaliable             // 4 available on blockchain
	TemplateDisabled              // 5 disable on blockchain
	TemplateFailed                // 6 voucher return error
)

// transfer status
const (
	TransferToApproval          int = iota     // 0 waiting to approval
	TransferApprovaled                         // 1 approved, transfer pending
	TransferReject                             // 2 reject
	TransferSomeSuccess                        // 3 some of the transfer success
	TransferFailed                             // 4 transfer failed
	TransferSuccess                            // 5 transfer success
	TransferCancel                             // 6 cancel transfer 撤回
	TransferInvalid                            // 7 invalid transfer 非法转账
	TransferDisuse                             // 8 approval out of date 审批过期
	TransferApproveExpire                      // 9 transfer out of date 转账过期
	TransferDisableEmployee                    // 10 tx canceled case user account disabled 禁用员工作废
	TransferCancelTemp                         // 11  tx canceled case template canceled 模板停用作废
	TransferInsufficientBalance                // 12 tx Insufficient balance
	TransferMinerToFloatError   int = iota + 8 // 21 矿工费转化float错误
	TransferCoinNotFound                       // 22 无此币种类别
	TransferGenTxError                         // 23 构建交易失败
	TransferVoucherReturnError                 // 24 签名机返回错误
	TransferDicemalError                       // 25 精度不够
	TransferAddressIndexError                  // 26 地址索引解析错误
	TransferDecodeMsgError                     // 27 msg解析错误
	TransferDecodeDataError                    // 28 数据解析错误
)

// user registration record status
const (
	RegWaitingToApproval int = iota // 0 注册
	RegToApproval                   // 1 扫过码，等待审核
	RegApprovaled                   // 2 同意
	RegRejected                     // 3 拒绝
	OwnerRegNotFound
	OwnerRegError
)

// owner account operation status
const (
	OwnerAccApprovaling int = iota
	OwnerAccApprovaled
	OwnerAccRejected
	RecoveryReject // owner register status
	HasApproved
	RecoveryInvalid
)

// account type
const (
	NormalAccType int = iota
	AdminAccType
	OwnerAccType
)

// message type
const (
	NormalMessage int = iota
	WarnMessage
	ErrorMessage
)

const (
	NodeStable int = iota
	NodeSyncing
)

const (
	LoggerOwnerRecovery  = "ownerRecovery"  // 申请恢复股东账户
	LoggerOwnerSubmit    = "ownerSubmit"    // 扫码提交
	LoggerOwnerAccept    = "ownerAccept"    // 同意恢复
	LoggerOwnerReject    = "ownerReject"    // 拒绝恢复
	LoggerTemplateNew    = "templateNew"    // 新建审批流模板
	LoggerTemplateAccept = "templateAccept" // 审批通过审批流模板
	LoggerTemplateReject = "templateReject" // 审批拒绝审批流模板
	LoggerTemplateCancel = "templateCancel" // 撤销审批流模板
	LoggerTemplateActive = "templateActive" // 审批流生效
	LoggerVoucherKeyWord = "voucherKeyWord" // 输入启动关键词
	LoggerVoucherStart   = "voucherStart"   // 签名机启动成功
	LoggerVoucherStop    = "voucherStop"    // 签名机关停成功
	LoggerVoucherTimeout = "voucherTimeout" // 启动超时
)

const (
	MsgOwnerAccountType = "ownerType" // 股东
	MsgAdminAccountType = "adminType" // 管理员
)

const (
	MsgTitleAppointAdmin            = "appointAdmin"            // 管理员任命通知
	MsgTitleDeleteAdmin             = "deleteAdmin"             // 取消管理员通知
	MsgTitleAuthorization           = "authorization"           // 功能授权通知
	MsgTitleNodeSwitch              = "nodeSwitch"              // 公链节点切换通知
	MsgTitleApprovalTemplateSuccess = "approvalTemplateSuccess" // 审批流模板审批成功
	MsgTitleApprovalTemplateFail    = "approvalTemplateFail"    // 审批流模板审批失败
	MsgTitleCancelTemplate          = "cancelTemplate"          // 审批流作废通知
	MsgTitleTransferApproved        = "transferApproved"        // 转账申请审批成功
	MsgTitleTransferReject          = "transferReject"          // 转账申请审批失败
	MsgTitleTransferSuccess         = "transferSuccess"         // 转账成功通知
	MsgTitleTransferFail            = "transferFail"            // 转账失败
	MsgTitleTransferCanceled        = "transferCanceled"        // 转账申请作废通知
	MsgTitleVoucherWarn             = "voucherWarn"             // 签名机验签失败警告
	MsgTitleNodeSync                = "nodeSync"                // 节点同步通知
)

const (
	MsgContentAppointAdmin                     = "appointAdmin"                     // 设置管理员
	MSgContentAppointAdminNoticeOthers         = "appointAdminNoticeOthers"         // 设置管理员通知其他管理员
	MsgContentDeleteAdmin                      = "deleteAdmin"                      // 已取消管理员权限
	MsgContentDeleteAdminNoticeOthers          = "deleteAdminNoticeOthers"          // 取消管理员权限通知其他管理员
	MsgContentAuthGrant                        = "authGrant"                        // 授予功能权限
	MsgContentAuthCancel                       = "authCancel"                       // 取消功能权限
	MsgContentNodeSwitch                       = "nodeSwitch"                       // 公链节点切换
	MsgContentApprovalTemplateSuccess          = "approvalTemplateSuccess"          // 审批流已全票通过审核
	MsgContentApprovalTemplateFailByVoucher    = "approvalTemplateFailByVoucher"    // 审批流由于签名机验签失败而被否决
	MsgContentApprovalTemplateFailByOwner      = "approvalTemplateFailByOwner"      // 审批流被股东否决
	MsgContentCancelTempByOwner                = "cancelTempByOwner"                // 审批流被作废
	MsgContentCancelTemp                       = "cancelTemp"                       // 审批流被作废
	MsgContentCancelTempByDisableUser          = "cancelTempByDisableUser"          // 停用账户时系统自动作废审批流
	MsgContentTransferApproved                 = "transferApproved"                 // 转账申请通过审批
	MsgContentTransferApprovalFaliCauseTimeOut = "transferApprovalFaliCauseTimeOut" // 转账由于超时审批失败
	MsgContentTransferReject                   = "transferReject"                   // 转账申请由于未达到最小审批通过人数而审批失败
	MsgContentTransferSuccess                  = "transferSuccess"                  // 转账成功
	MsgContentTransferFail                     = "transferFail"                     // 转账失败
	MsgContentTransferPartiallySuceeded        = "transferPartiallySuceeded"        // 部分转账成功
	MsgContentTransferFailCauseTimeout         = "transferFailCauseTimeout"         // 转账超过截至时间失败
	MsgContentTransferCanceled                 = "transferCanceled"                 // 转账申请被作废
	MsgContentTransferCanceledCauseDisableUser = "transferCanceledCauseDisableUser" // 由于审批流中的用户账户被停用而转账被作废
	MsgContentVoucherWarn                      = "voucherWarn"                      // 签名机报警
	MsgContentNodeSyncFail                     = "nodeSyncFail"                     // 节点同步异常
	MsgContentNodeSyncSuccess                  = "nodeSyncSUccess"                  // 节点同步恢复正常
)

const (
	AuthBalance = "balanceAuth" // 资产查询权限
	AuthAddress = "addressAuth" // 多账户管理权限
)

type MsgTemplatePadding struct {
	AccountType         string   `json:"accType"`
	OperAccountName     string   `json:"operaterAcc"`
	PromoterAccountName string   `json:"promoterAcc"`
	TemplateName        string   `json:"templateName"`
	DisabledAccount     string   `json:"disabledAcc"`
	TransOrderCoinName  string   `json:"transOrderCoinName"`
	TransOrderAmount    string   `json:"transOrderAmount"`
	VoucherWarnCount    int      `json:"voucherWarnCount"`
	AuthNames           []string `json:"authName"`
	AuthName            string
}
