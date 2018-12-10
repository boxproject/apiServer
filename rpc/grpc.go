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
	 "context"
	 "crypto/tls"
	 "crypto/x509"
	 "encoding/json"
	 "fmt"
	 "github.com/boxproject/apiServer/static/pb"
	 "github.com/boxproject/apiServer/config"
	 "github.com/go-errors/errors"
	 "google.golang.org/grpc"
	 "google.golang.org/grpc/credentials"
	 "google.golang.org/grpc/keepalive"
	 "google.golang.org/grpc/reflection"
	 "io"
	 "io/ioutil"
	 "net"
	 "strings"
	 "sync"
	 "time"
	 log "github.com/alecthomas/log4go"
 )


type streamCtrl struct {
	stream   pb.Synchronizer_ListenServer
	quitChan chan struct{}
}

type serverStruct struct {
	mu                  *sync.Mutex // protects routeNotes
	client              map[string]*streamCtrl
	handleShakeQuitChan chan struct{}
}

var defaultstream = &serverStruct{
	client:              make(map[string]*streamCtrl),
	mu:                  new(sync.Mutex),
	handleShakeQuitChan: make(chan struct{}),
}

type routerChan struct {
	routerMu  *sync.Mutex
	routerMap map[string]chan interface{}
}

var gRouterChan = &routerChan{
	routerMu:  new(sync.Mutex),
	routerMap: make(map[string]chan interface{}),
}

func IsConnect() bool {
	lockLastTime.Lock()
	defer lockLastTime.Unlock()
	now := time.Now().UnixNano() / time.Millisecond.Nanoseconds()
	if (now - lastTime) > timeout {
		//超时
		return false
	} else {
		return true
	}

}

var HeartStatus []int

func NewClientChan(key string) chan interface{} {
	defer gRouterChan.routerMu.Unlock()
	gRouterChan.routerMu.Lock()
	if routerChan, ok := gRouterChan.routerMap[key]; !ok {
		gRouterChan.routerMap[key] = make(chan interface{}, 1)
		return gRouterChan.routerMap[key]
	} else {
		return routerChan
	}
	return nil
}
func GetClientChan(key string) chan interface{} {
	defer gRouterChan.routerMu.Unlock()
	gRouterChan.routerMu.Lock()
	if routerChan, ok := gRouterChan.routerMap[key]; ok {
		return routerChan
	}
	return nil
}
func UpdateClientChan(key string, info interface{}) {
	defer gRouterChan.routerMu.Unlock()
	gRouterChan.routerMu.Lock()
	if sendChan, ok := gRouterChan.routerMap[key]; ok {
		sendChan <- info
	}
}

func RemoveChan(key string) {
	defer gRouterChan.routerMu.Unlock()
	gRouterChan.routerMu.Lock()
	delete(gRouterChan.routerMap, key)
}

func WaitClientRep(oper chan interface{}, timeout time.Duration) (interface{}, error) {
	timeOutTicker := time.NewTicker(timeout)
	if oper == nil {
		return nil, errors.New("chan is nil")
	}
	for {
		select {
		case <-timeOutTicker.C:
			timeOutTicker.Stop()
			close(oper)
			return nil, errors.New("timeout")
		case data, ok := <-oper:
			if ok {
				close(oper)
				return data, nil
			}
			return nil, errors.New("chan error")
		}
	}
}

// send to client
func SendToClient(sendType string, clientName string, msg []byte) bool {
	defer defaultstream.mu.Unlock()
	defaultstream.mu.Lock()
	status := false
	for key, value := range defaultstream.client {
		if strings.HasPrefix(key, clientName) {
			if err := value.stream.Send(&pb.StreamRsp{Type: sendType, Msg: msg}); err != nil {
				close(value.quitChan)
				return false
			} else {
				status = true
			}
		}
	}
	return status
}

func SendToClientRetErr(sendType string, clientName string, msg []byte) error {
	defer defaultstream.mu.Unlock()
	defaultstream.mu.Lock()
	var err error = fmt.Errorf("no client ...")
	for key, value := range defaultstream.client {
		if strings.HasPrefix(key, clientName) {
			if err = value.stream.Send(&pb.StreamRsp{Type: sendType, Msg: msg}); err != nil {
				close(value.quitChan)
				return err
			} else {
				err = nil
			}
		}
	}
	return err
}

func (s *serverStruct) Router(ctx context.Context, req *pb.RouterRequest) (*pb.RouterResponse, error) {
	//fmt.Println("Router")
	oper := &GrpcClient{}
	if err := json.Unmarshal(req.Msg, oper); err == nil {
		//fmt.Println(oper)
		go UpdateClientChan(req.RouterType, oper)
	} else {
		log.Error("json.Unmarshal error:", err)
	}
	//if routerChan ,ok := routerMap[oper.Timestamp];ok {
	//	routerChan <- oper
	//}
	return new(pb.RouterResponse), nil
}

func (s *serverStruct) Listen(stream pb.Synchronizer_ListenServer) (err error) {
	defer fmt.Println("Listen end ...")
	fmt.Println("Listen begin ...")
	listReq, err := stream.Recv()
	if err != nil {
		return stream.Context().Err()
	}
	key := fmt.Sprintf("%v_%v_%v", listReq.ServerName, listReq.Name, listReq.Ip)
	fmt.Printf("listReq: %s\n", key)

	// limit connect number
	if len(s.client) >= 1 {
		//fmt.Println(" stream done...")
		return stream.Context().Err()
	}

	s.mu.Lock()
	quit := make(chan struct{}, 1)
	if _, ok := s.client[key]; !ok {
		s.client[key] = &streamCtrl{
			stream:   stream,
			quitChan: quit,
		}
	} else {
		fmt.Printf("clent:%v is already have\n", key)
		return nil
	}
	s.mu.Unlock()

	//监控连接情况
	rets := true
	for rets {
		select {
		case <-quit:
			quit = nil
			rets = false
			break
		default:
			heartMsg, err := stream.Recv()
			if err == io.EOF {
				log.Error("err EOF...", err)
				log.Error("clien:%v closed", key)
				rets = false
				break
			}
			if err != nil {
				log.Error("[LISTEN ERR] %v\n", err)
				log.Error("clien:%v closed", key)
				rets = false
				break
			}
			UpdateVoucherStatus(heartMsg.Msg)

			fmt.Printf("heart:%+v\n", GetVoucherStatus())
		}
	}
	s.mu.Lock()
	log.Error("delete key-----------------------")
	delete(s.client, key)
	s.mu.Unlock()

	return nil
}

func InitGrpc() error {
	fmt.Println("load grpc start------")
	//cfg,_ := config.LoadConfig()
	cfg := config.Conf
	// load cred
	cred, err := loadCredential(cfg.Voucher.ServerPem, cfg.Voucher.ServerKey, cfg.Voucher.ClientPem)
	if err != nil {
		fmt.Println("load tls cert failed. cause: %v\n", err)
		return err
	}
	options := []grpc.ServerOption{
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{MinTime: time.Minute, PermitWithoutStream: true}),
		grpc.Creds(cred),
	}
	server := grpc.NewServer(options...)
	pb.RegisterSynchronizerServer(server, defaultstream)
	reflection.Register(server)
	defer server.Stop()
	defer defaultstream.Stop()

	fmt.Println("load grpc start1------", "localhost:"+cfg.Voucher.Port)
	lis, err := net.Listen("tcp", ":"+cfg.Voucher.Port)
	if err != nil {
		fmt.Printf("Can not listen to the port %v, cause: %v\n", "localhost:"+cfg.Voucher.Port, err)
		return err
	}
	go defaultstream.HandleShake()
	// start server , begin listen
	if err = server.Serve(lis); err != nil {
		fmt.Printf("gRPC service error, cause: %v\n", err)
		return err
	}
	return nil
}

//加载服务端证书
func loadCredential(ServerCert, ServerKey, ClientCert string) (credentials.TransportCredentials, error) {
	cert, err := tls.LoadX509KeyPair(ServerCert, ServerKey)
	if err != nil {
		return nil, err
	}

	certBytes, err := ioutil.ReadFile(ClientCert)
	if err != nil {
		return nil, err
	}

	clientCertPool := x509.NewCertPool()
	ok := clientCertPool.AppendCertsFromPEM(certBytes)
	if !ok {
		return nil, err
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    clientCertPool,
	}

	return credentials.NewTLS(config), nil
}

func (s *serverStruct) HandleShake() error {
	//监控连接情况
	timeTicker := time.NewTicker(time.Nanosecond)
	loop := true
	for loop {
		select {
		case <-timeTicker.C:
			timeTicker.Stop()
			timeTicker = time.NewTicker(time.Second * 5)
			SendToClient(HANDLE_SHAKE, "voucher", nil)
			break
		case <-s.handleShakeQuitChan:
			s.handleShakeQuitChan = nil
			return nil
		}
	}
	return nil
}
func (s *serverStruct) Stop() error {
	fmt.Println("stop grpc")
	if s.handleShakeQuitChan != nil {
		close(s.handleShakeQuitChan)
		s.handleShakeQuitChan = nil
	}
	return nil
}
