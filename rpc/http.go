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

import (
	"encoding/json"
	"fmt"
	"time"
	log "github.com/alecthomas/log4go"
	"github.com/satori/go.uuid"
	"github.com/boxproject/apiServer/errors"
	"sync"
)

var lastTime int64
var lockLastTime = new(sync.Mutex)

const timeout = int64(5000)

var vouchersatus *VoucherStatus

func UpdateVoucherStatus(msg []byte) {
	lockLastTime.Lock()
	defer lockLastTime.Unlock()
	status := &VoucherStatus{}
	if err := json.Unmarshal(msg, status); err == nil {
		vouchersatus = status
		lastTime = time.Now().UnixNano() / time.Millisecond.Nanoseconds()

	} else {
		fmt.Println("json Unmarshal error:", err)
	}
}

func GetVoucherStatus() VoucherStatus {
	lockLastTime.Lock()
	defer lockLastTime.Unlock()
	if vouchersatus == nil {
		return VoucherStatus{}
	}
	return *vouchersatus
}

func SendVoucherData(oper *GrpcServer) (*VoucherStatus, *GrpcClient, int) {
	switch oper.Type {
	case VOUCHER_OPERATE_ADDKEY:
		break
	case VOUCHER_OPERATE_CREATE:
		break
	case VOUCHER_OPERATE_DEPLOY:
		break
	case VOUCHER_OPERATE_START:
		break
	case VOUCHER_OPERATE_STOP:
		break
	case VOUCHER_OPERATE_APP_PUBKEY:
		break
	case VOUCHER_OPERATE_HDBAKUP:
		break
	case VOUCHER_OPERATE_SIGN:
		break
	case VOUCHER_OPERATE_CHECK_PASS:
		break
	case VOUCHER_OPERATE_MASTER_PUBKEY:
		break
	case VOUCHER_OPERATE_ADDHASH:
		break
	case VOUCHER_OPERATE_DISHASH:
		break
	case VOUCHER_OPERATE_HASH_AVAILABLE:
		break
	case VOUCHER_OPERATE_RECOVER:
		break
	case VOUCHER_OPERATE_TRANSSIGN_INTERIOR:
		break
	case VOUCHER_OPERATE_KEY_EXCHANGE:
		break

	case VOUCHER_STATUS:
		//encoder := json.NewEncoder(rw)
		vs := GetVoucherStatus()
		now := time.Now().UnixNano() / time.Millisecond.Nanoseconds()
		if (now - lastTime) > timeout {
			//超时
			return &vs, nil, VRET_TIMEOUT
		}
		return &vs, nil, VRET_STATUS
	default:
		{
			log.Error("无此类别")
			return nil, nil, VRET_ERR
		}

	}
	sendType := fmt.Sprintf("%v", uuid.Must(uuid.NewV4()).String())
	msg, err := json.Marshal(oper)
	if err != nil {
		log.Error("json转化失败", err)
		return nil, nil, VRET_ERR
	} else if SendToClient(sendType, "voucher", msg) == false {
		log.Error("发送签名机消息失败")
		return nil, nil, VRET_TIMEOUT
	} else {
		req, err := WaitClientRep(NewClientChan(sendType), time.Second*3)
		RemoveChan(sendType)
		if err != nil {
			log.Error("与签名机连接失败", err)
			return nil, nil, VRET_TIMEOUT
		} else {
			if req != nil {
				operReq := req.(*GrpcClient)
				return nil, operReq, VRET_CLIENT

			} else {
				log.Error("返回值错误")
				return nil, nil, VRET_ERR
			}
		}
	}
	//return nil,nil,VRET_ERR
}

// 获取签名机状态
func SendToVoucher(oper *GrpcServer) (*GrpcClient, error) {
	switch oper.Type {
	case VOUCHER_OPERATE_ADDKEY:
		break
	case VOUCHER_OPERATE_CREATE:
		break
	case VOUCHER_OPERATE_DEPLOY:
		break
	case VOUCHER_OPERATE_START:
		break
	case VOUCHER_OPERATE_STOP:
		break
	case VOUCHER_OPERATE_APP_PUBKEY:
		break
	case VOUCHER_OPERATE_HDBAKUP:
		break
	case VOUCHER_OPERATE_SIGN:
		break
	case VOUCHER_OPERATE_CHECK_PASS:
		break
	case VOUCHER_OPERATE_MASTER_PUBKEY:
		break
	case VOUCHER_OPERATE_ADDHASH:
		break
	case VOUCHER_OPERATE_DISHASH:
		break
	case VOUCHER_OPERATE_HASH_AVAILABLE:
		break
	case VOUCHER_OPERATE_RECOVER:
		break
	case VOUCHER_OPERATE_TRANSSIGN_INTERIOR:
		break
	default:
		{
			log.Error("无此类别")
			return nil, errors.New("Invalidate RPC Code.")
		}

	}
	sendType := fmt.Sprintf("%v", uuid.Must(uuid.NewV4()).String())
	msg, err := json.Marshal(oper)
	if err != nil {
		log.Error("json转化失败", err)
		return nil, err
	} else if err = SendToClientRetErr(sendType, "voucher", msg); err != nil {
		log.Error("发送签名机消息失败")
		return nil, err
	} else {
		req, err := WaitClientRep(NewClientChan(sendType), time.Second*3)
		RemoveChan(sendType)
		if err != nil {
			log.Error("与签名机连接失败", err)
			return nil, err
		} else {
			if req != nil {
				operReq := req.(*GrpcClient)
				return operReq, nil
			} else {
				log.Error("返回值错误")
				return nil, errors.New("Voucher Return Error")
			}
		}
	}
}
