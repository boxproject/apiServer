// Copyright 2018. bolaxy.org authors.
// 
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// 
// 		 http://www.apache.org/licenses/LICENSE-2.0
// 
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package timer

import (
	"fmt"
	"testing"
	"encoding/json"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"

	"github.com/boxproject/boxwallet/bccore"

)

//初始化gprc

func TestTransfer(t *testing.T) {


	//go server.InitGrpc()
	//TransferTimer()
	//
	//quitCh := make(chan os.Signal, 1)
	//signal.Notify(quitCh,
	//	syscall.SIGINT, syscall.SIGTERM,
	//	syscall.SIGHUP, syscall.SIGKILL,
	//	syscall.SIGUSR1, syscall.SIGUSR2)
	//// todo: add your code
	//<-quitCh

	www,_ := Wallet.GetBalance(bccore.STR_ETH,"0xE7eA18d7E9d6A2374028E0741FC83C4447d3c93D","")
	fmt.Println(www.String())
}





//验证转账签名
//msg 结构体GrpcServer里的 transMsg 
//该审批流下的所有员工公钥
func ValidateSign (msg string,approversPubkey map[string]string) bool{
	transferMsg := &TransferInfo{}
	err := json.Unmarshal([]byte(msg),transferMsg)
	if err != nil {
		fmt.Println("json解析错误")
		return false
	}
	if RsaVerify([]byte(transferMsg.ApplyPublickey),[]byte(transferMsg.TransferMsgs),[]byte(transferMsg.ApplySign))!=nil {
		fmt.Println("申请人签名错误")
		return false
	}
	//获取审批人签名
	approversSign := make(map[string]string)
    err = json.Unmarshal([]byte(transferMsg.ApproversSign), &approversSign)
	if err != nil {
		fmt.Println("审批人map解析错误")
		return false
	}
	for k, v := range approversSign {
		//获取该审批人的公钥
		if RsaVerify([]byte(approversPubkey[k]),[]byte(transferMsg.TransferMsgs),[]byte(v))!=nil {
			fmt.Println("审批人签名错误")
			return false
		}
	}
	return true
}

func RsaVerify(pubKeyBytes, msg, signature []byte) error {
	bSignData, err := base64.StdEncoding.DecodeString(string(signature))
	if err != nil {
		return fmt.Errorf("DecodeString sign error:%v", err)
	}
	hashed := sha256.Sum256(msg)
	bKey, err := base64.StdEncoding.DecodeString(string(pubKeyBytes))
	if err != nil {
		return fmt.Errorf("DecodeString pubKeyBytes  error:%v", err)
	}
	pub, err := x509.ParsePKCS1PublicKey(bKey)
	if err != nil {
		return fmt.Errorf("ParsePKCS1PublicKey error:%v", err)
	}
	err = rsa.VerifyPKCS1v15(pub, crypto.SHA256, hashed[:], bSignData)
	if err != nil {
		return fmt.Errorf("VerifyPKCS1v15 error:%v", err)
	}
	return nil
}

