package message

import (
	"strconv"
	"strings"
	voucher "github.com/boxproject/apiServer/rpc"
	"github.com/boxproject/apiServer/db"
	"github.com/boxproject/apiServer/common"
	log "github.com/alecthomas/log4go"
	"github.com/spf13/viper"
)

func init(){
	// zh
	vzh := viper.New()
	vzh.SetConfigType(common.TOML_FILE)
	vzh.SetConfigFile(ZhMsgFilePath)
	if err := vzh.ReadInConfig(); err != nil {
		log.Error("read message_zh.toml error", err)
	}

	MsgTitleTempMap[common.LangZhType] = vzh.GetStringMap(MsgTitleTag)
	MsgAccountTypeMap[common.LangZhType] = vzh.GetStringMap(MsgAccountTypeTag)
	MsgContentTempMap[common.LangZhType] = vzh.GetStringMap(MsgContentTag)

	// en
	ven := viper.New()
	ven.SetConfigType(common.TOML_FILE)
	ven.SetConfigFile(EnMsgFilePath)
	if err := ven.ReadInConfig(); err != nil {
		log.Error("read message_en.toml error", err)
	}
	MsgTitleTempMap[common.LangEnType] = ven.GetStringMap(MsgTitleTag)
	MsgAccountTypeMap[common.LangEnType] = ven.GetStringMap(MsgAccountTypeTag)
	MsgContentTempMap[common.LangEnType] = ven.GetStringMap(MsgContentTag)

}
// AddAdmin send message when add admin
func AddAdmin(add string, accountId int) error {
	ms := &db.MessageService{}

	parsedAdd := strings.Split(add, ",")
	ids := make([]int, 0, len(parsedAdd))
	for _, v := range parsedAdd {
		id, _ := strconv.Atoi(v)
		ids = append(ids, id)
	}
	err := ms.AddAdmin(ids, accountId)
	return err
}

// DelAdmin send message when delete admin
func DelAdmin(accountID, id int) error {
	ms := &db.MessageService{}
	err := ms.DelAdmin(accountID, id)
	return err
}

// ChangeNode send message when node changed
func ChangeNode() error {
	ms := &db.MessageService{}
	err := ms.ChangeNode()
	return err
}

// SuccessVerifyTemp send message when verify template
func SuccessVerifyTemp(name, id string, creatorId int) error {
	ms := &db.MessageService{}
	err := ms.SuccessVerifyTemp(name, id, creatorId)
	return err
}

// FailVerifyTemp send message when verify template failed
func FailVerifyTemp(opername, name, id string, creatorId int) error {
	ms := &db.MessageService{}
	err := ms.FailVerifyTemp(opername, name, id, creatorId)
	return err
}

// SuccessVerifyTrans send message when verify transfer
func SuccessVerifyTrans(id string) error {
	ms := &db.MessageService{}
	ts := &db.TransferDBService{}
	//通过orderid查找转账信息
	transferOrder, err := ts.FindOrderById(id)
	err = ms.SuccessVerifyTrans(transferOrder.CoinName, transferOrder.Amount, transferOrder.ApplyerId, id)
	return err
}

// FailVerifyTrans send message when verify transfer failed
func FailVerifyTrans(id, reason string) error {
	ms := &db.MessageService{}
	ts := &db.TransferDBService{}
	//通过orderid查找转账信息
	transferOrder, err := ts.FindOrderById(id)
	err = ms.FailVerifyTrans(transferOrder.CoinName, transferOrder.Amount, reason, transferOrder.ApplyerId, id)
	return err
}

// TransferSuccess send message when transfer success
func TransferSuccess(CoinName, Amount, TransferId string, ApplyerId int) error {
	ms := &db.MessageService{}
	err := ms.TransferSuccess(CoinName, Amount, TransferId, ApplyerId)
	return err
}

// TransferFail send message when transfer failed
func TransferFail(CoinName, Amount, TransferId string, ApplyerId int) error {
	ms := &db.MessageService{}
	err := ms.TransferFail(CoinName, Amount, TransferId, ApplyerId,common.NormalMessage)
	return err
}

func TransferPartialSuccessful(coinName, amount, orderId string, applyerId int) error {
	ms := &db.MessageService{}
	err := ms.TransferPartialSuccessful(coinName, amount, orderId, applyerId)
	return err
}

// TransferFail send message when transfer failed cause out of date
func TransferFailOutOfDate(CoinName, Amount, TransferId string, ApplyerId int) error {
	ms := &db.MessageService{}
	err := ms.TransferFailOutOfDate(CoinName, Amount, TransferId, ApplyerId,common.NormalMessage)
	return err
}

func TransferApprovalTimeOut(CoinName, Amount, TransferId string, ApplyerId int) error {
	ms := &db.MessageService{}
	err := ms.TransferApprovalTimeOut(CoinName, Amount, TransferId, ApplyerId,common.NormalMessage)
	return err
}

// CancelTemplate send message when disable template
func CancelTemplate(opername, userAccount, name, id string) error {
	ms := &db.MessageService{}
	err := ms.CancelTemplate(opername, userAccount, name, id)
	return err
}

// CancelTransfer send message when disable transfer
func CancelTransfer(ApplyerId int, CoinName, Amount, id string, isCancelUser bool) error {
	ms := &db.MessageService{}
	err := ms.CancelTransfer(ApplyerId, CoinName, Amount, id, isCancelUser)
	return err
}

// AddAuth send message when add auth to account
func AddAuth(add, authId []int, userType int, account string) error {
	ms := &db.MessageService{}
	err := ms.AddAuth(add, authId, userType, account)
	return err
}

// DelAuth send message when delete auth from account
func DelAuth(id, accountId, userType int, account string) error {
	ms := &db.MessageService{}
	err := ms.DelAuth(id, accountId, userType, account)
	return err
}

// 签名机报警
func VoucherWarn(msgCount, receiverAccId, pageType int, param string) (*db.Message, error) {
	ms := &db.MessageService{}
	return ms.VoucherWarn(receiverAccId, msgCount, pageType, param)
}

// 发送签名机异常站内信
func VoucherFail(msgType, vRes int, tempInfo *db.Template, transApply *db.TransferSend ) error {
	log.Debug("VoucherFail message...")
	ms := &db.MessageService{}
	var err error
	if vRes == voucher.STATUS_APP_LOSEBASEINFO ||  vRes == voucher.STATUS_APP_RSA_SIGN_ERROR || vRes == voucher.STATUS_APP_VERIFY_ERROR ||
		vRes == voucher.STATUS_APP_RSA_DECRYPT_ERROR || vRes == voucher.STATUS_APP_AES_DECRYPT_ERROR || vRes == voucher.STATUS_FLOW_VERIFY_ERROR ||
			vRes == voucher.STATUS_HASH_LOSEBASEINFO || vRes == voucher.STATUS_TRANS_LOSEBASEINFO || vRes == voucher.STATUS_TRANS_VERIFY_ERROR ||
				vRes == voucher.STATUS_TRANS_SIGN_ERROR || vRes == voucher.STATUS_TRANS_UNKNOW_CURRENCY || vRes == voucher.STATUS_TRANS_CURRENCY_NOTEQUAL ||
					vRes == voucher.STATUS_TRANS_MARSHAL_ERROR || vRes == voucher.STATUS_TRANS_GET_INFO_ERROR || vRes == voucher.STATUS_TRANS_TOADDR_NOT_INTERIOR ||
						vRes == voucher.STATUS_TRANS_SIGN_REPEAT {
							if msgType == VoucherRetTempFailMsgType {
								// 审批流失败
								err = ms.FailVerifyTemp("", tempInfo.Name, tempInfo.ID, tempInfo.CreatorID)
								} else {
									// 转账失败
									err = ms.TransferFail(transApply.CoinName, transApply.Amount, transApply.TransferId, transApply.ApplyerId, common.WarnMessage)
									}
								return err
								}
							return nil
}
//同步节点状态  -- 1.恢复正常 2。异常
func AysncBlockChainStatus(types int){
	ms := &db.MessageService{}
	err := ms.AysncBlockChainStatus(types)
	if err !=nil {
		log.Error("同步节点站内信错误")
	}
}


