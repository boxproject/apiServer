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
	"sync"
	"time"

	log "github.com/alecthomas/log4go"
	"github.com/boxproject/boxwallet/cli"
	"github.com/boxproject/apiServer/db"
)

var Wallet = cli.NewAppServer()

var USDT_PPID = db.USDT_PPID

var timerLock = new(sync.Mutex)
//是否循环获取主地址
var masterRet = true

//获取所有主地址
var Time2 = time.NewTicker(3 * time.Second)
func TransferTimer() {
	
	//充值
	timer3 := time.NewTicker(30 * time.Second)
	//获取余额
	timer5 := time.NewTicker(20 * time.Second)
	//获取签名机状态
	timer6 := time.NewTicker(3 * time.Second)
	// 更新余额
	timer7 := time.NewTicker(60 * time.Second)
	//更新各个地址余额
	timer8 := time.NewTicker(5 * time.Second)




	go func(){
		for{
			select{
				case <-timer3.C:
					{
						log.Debug(333)
						depositBusiness()
					}
			}
		}
	}()

	go func(){
		for{
			select{
			case <-timer6.C:
				{
					log.Debug(666)
					updateConnectStatus()
					updateAsyncStatus()
				}
			}
		}
	}()

	go func(){
		for{
			select{
			case <-timer5.C:
				{
					log.Debug(555)
					updateDepositBalance()
				}
			}
		}
	}()

	go func(){
		for{
			select{
			case <-timer7.C:
				{
					log.Debug(777)
					sumBalance()
				}
			}
		}
	}()

	go func(){
		for{
			select{
			case <-timer8.C:
				{
					log.Debug(888)
					updateBalance()
				}
			}
		}
	}()

	go func() {
		<-time.After(time.Second*30)
		transferChan()
	}()


	for masterRet{
		select {
	
		case <-Time2.C:
			{
				log.Debug(222)
				go getAllmasterAddress()
			}
		}

	}
}
