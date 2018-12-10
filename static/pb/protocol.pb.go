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
// Code generated by protoc-gen-go. DO NOT EDIT.
// source: pb/protocol.proto

package pb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type RouterRequest struct {
	RouterType           string   `protobuf:"bytes,1,opt,name=routerType,proto3" json:"routerType,omitempty"`
	RouterName           string   `protobuf:"bytes,2,opt,name=routerName,proto3" json:"routerName,omitempty"`
	Msg                  []byte   `protobuf:"bytes,3,opt,name=msg,proto3" json:"msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RouterRequest) Reset()         { *m = RouterRequest{} }
func (m *RouterRequest) String() string { return proto.CompactTextString(m) }
func (*RouterRequest) ProtoMessage()    {}
func (*RouterRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_protocol_9525ba12f913ea0d, []int{0}
}
func (m *RouterRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RouterRequest.Unmarshal(m, b)
}
func (m *RouterRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RouterRequest.Marshal(b, m, deterministic)
}
func (dst *RouterRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RouterRequest.Merge(dst, src)
}
func (m *RouterRequest) XXX_Size() int {
	return xxx_messageInfo_RouterRequest.Size(m)
}
func (m *RouterRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_RouterRequest.DiscardUnknown(m)
}

var xxx_messageInfo_RouterRequest proto.InternalMessageInfo

func (m *RouterRequest) GetRouterType() string {
	if m != nil {
		return m.RouterType
	}
	return ""
}

func (m *RouterRequest) GetRouterName() string {
	if m != nil {
		return m.RouterName
	}
	return ""
}

func (m *RouterRequest) GetMsg() []byte {
	if m != nil {
		return m.Msg
	}
	return nil
}

type RouterResponse struct {
	Code                 string   `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg                  []byte   `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RouterResponse) Reset()         { *m = RouterResponse{} }
func (m *RouterResponse) String() string { return proto.CompactTextString(m) }
func (*RouterResponse) ProtoMessage()    {}
func (*RouterResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_protocol_9525ba12f913ea0d, []int{1}
}
func (m *RouterResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RouterResponse.Unmarshal(m, b)
}
func (m *RouterResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RouterResponse.Marshal(b, m, deterministic)
}
func (dst *RouterResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RouterResponse.Merge(dst, src)
}
func (m *RouterResponse) XXX_Size() int {
	return xxx_messageInfo_RouterResponse.Size(m)
}
func (m *RouterResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RouterResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RouterResponse proto.InternalMessageInfo

func (m *RouterResponse) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *RouterResponse) GetMsg() []byte {
	if m != nil {
		return m.Msg
	}
	return nil
}

type ListenReq struct {
	ServerName           string   `protobuf:"bytes,1,opt,name=serverName,proto3" json:"serverName,omitempty"`
	Name                 string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Ip                   string   `protobuf:"bytes,3,opt,name=ip,proto3" json:"ip,omitempty"`
	Type                 string   `protobuf:"bytes,4,opt,name=type,proto3" json:"type,omitempty"`
	Msg                  []byte   `protobuf:"bytes,5,opt,name=msg,proto3" json:"msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListenReq) Reset()         { *m = ListenReq{} }
func (m *ListenReq) String() string { return proto.CompactTextString(m) }
func (*ListenReq) ProtoMessage()    {}
func (*ListenReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_protocol_9525ba12f913ea0d, []int{2}
}
func (m *ListenReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListenReq.Unmarshal(m, b)
}
func (m *ListenReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListenReq.Marshal(b, m, deterministic)
}
func (dst *ListenReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListenReq.Merge(dst, src)
}
func (m *ListenReq) XXX_Size() int {
	return xxx_messageInfo_ListenReq.Size(m)
}
func (m *ListenReq) XXX_DiscardUnknown() {
	xxx_messageInfo_ListenReq.DiscardUnknown(m)
}

var xxx_messageInfo_ListenReq proto.InternalMessageInfo

func (m *ListenReq) GetServerName() string {
	if m != nil {
		return m.ServerName
	}
	return ""
}

func (m *ListenReq) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *ListenReq) GetIp() string {
	if m != nil {
		return m.Ip
	}
	return ""
}

func (m *ListenReq) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *ListenReq) GetMsg() []byte {
	if m != nil {
		return m.Msg
	}
	return nil
}

type StreamRsp struct {
	Type                 string   `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	Msg                  []byte   `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StreamRsp) Reset()         { *m = StreamRsp{} }
func (m *StreamRsp) String() string { return proto.CompactTextString(m) }
func (*StreamRsp) ProtoMessage()    {}
func (*StreamRsp) Descriptor() ([]byte, []int) {
	return fileDescriptor_protocol_9525ba12f913ea0d, []int{3}
}
func (m *StreamRsp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StreamRsp.Unmarshal(m, b)
}
func (m *StreamRsp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StreamRsp.Marshal(b, m, deterministic)
}
func (dst *StreamRsp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StreamRsp.Merge(dst, src)
}
func (m *StreamRsp) XXX_Size() int {
	return xxx_messageInfo_StreamRsp.Size(m)
}
func (m *StreamRsp) XXX_DiscardUnknown() {
	xxx_messageInfo_StreamRsp.DiscardUnknown(m)
}

var xxx_messageInfo_StreamRsp proto.InternalMessageInfo

func (m *StreamRsp) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *StreamRsp) GetMsg() []byte {
	if m != nil {
		return m.Msg
	}
	return nil
}

func init() {
	proto.RegisterType((*RouterRequest)(nil), "pb.RouterRequest")
	proto.RegisterType((*RouterResponse)(nil), "pb.RouterResponse")
	proto.RegisterType((*ListenReq)(nil), "pb.ListenReq")
	proto.RegisterType((*StreamRsp)(nil), "pb.StreamRsp")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// SynchronizerClient is the client API for Synchronizer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type SynchronizerClient interface {
	Router(ctx context.Context, in *RouterRequest, opts ...grpc.CallOption) (*RouterResponse, error)
	Listen(ctx context.Context, opts ...grpc.CallOption) (Synchronizer_ListenClient, error)
}

type synchronizerClient struct {
	cc *grpc.ClientConn
}

func NewSynchronizerClient(cc *grpc.ClientConn) SynchronizerClient {
	return &synchronizerClient{cc}
}

func (c *synchronizerClient) Router(ctx context.Context, in *RouterRequest, opts ...grpc.CallOption) (*RouterResponse, error) {
	out := new(RouterResponse)
	err := c.cc.Invoke(ctx, "/pb.Synchronizer/router", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *synchronizerClient) Listen(ctx context.Context, opts ...grpc.CallOption) (Synchronizer_ListenClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Synchronizer_serviceDesc.Streams[0], "/pb.Synchronizer/listen", opts...)
	if err != nil {
		return nil, err
	}
	x := &synchronizerListenClient{stream}
	return x, nil
}

type Synchronizer_ListenClient interface {
	Send(*ListenReq) error
	Recv() (*StreamRsp, error)
	grpc.ClientStream
}

type synchronizerListenClient struct {
	grpc.ClientStream
}

func (x *synchronizerListenClient) Send(m *ListenReq) error {
	return x.ClientStream.SendMsg(m)
}

func (x *synchronizerListenClient) Recv() (*StreamRsp, error) {
	m := new(StreamRsp)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// SynchronizerServer is the server API for Synchronizer service.
type SynchronizerServer interface {
	Router(context.Context, *RouterRequest) (*RouterResponse, error)
	Listen(Synchronizer_ListenServer) error
}

func RegisterSynchronizerServer(s *grpc.Server, srv SynchronizerServer) {
	s.RegisterService(&_Synchronizer_serviceDesc, srv)
}

func _Synchronizer_Router_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RouterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SynchronizerServer).Router(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Synchronizer/Router",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SynchronizerServer).Router(ctx, req.(*RouterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Synchronizer_Listen_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(SynchronizerServer).Listen(&synchronizerListenServer{stream})
}

type Synchronizer_ListenServer interface {
	Send(*StreamRsp) error
	Recv() (*ListenReq, error)
	grpc.ServerStream
}

type synchronizerListenServer struct {
	grpc.ServerStream
}

func (x *synchronizerListenServer) Send(m *StreamRsp) error {
	return x.ServerStream.SendMsg(m)
}

func (x *synchronizerListenServer) Recv() (*ListenReq, error) {
	m := new(ListenReq)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _Synchronizer_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.Synchronizer",
	HandlerType: (*SynchronizerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "router",
			Handler:    _Synchronizer_Router_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "listen",
			Handler:       _Synchronizer_Listen_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "pb/protocol.proto",
}

func init() { proto.RegisterFile("pb/protocol.proto", fileDescriptor_protocol_9525ba12f913ea0d) }

var fileDescriptor_protocol_9525ba12f913ea0d = []byte{
	// 268 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x51, 0xc1, 0x4e, 0xc3, 0x30,
	0x0c, 0x5d, 0xb2, 0x51, 0xa9, 0xd6, 0x36, 0x31, 0x9f, 0xaa, 0x1d, 0xd0, 0xd4, 0x53, 0x0f, 0xa8,
	0x30, 0x90, 0xf8, 0x0a, 0xc4, 0x21, 0xe3, 0x07, 0xda, 0x62, 0x41, 0xa5, 0x35, 0xc9, 0x92, 0x14,
	0x69, 0x7c, 0x3d, 0x4a, 0xba, 0x2e, 0x45, 0xdc, 0x5e, 0x9e, 0xfd, 0xfc, 0x9e, 0x1d, 0xd8, 0xe8,
	0xfa, 0x41, 0x1b, 0xe5, 0x54, 0xa3, 0x8e, 0x65, 0x00, 0xc8, 0x75, 0x9d, 0x57, 0xb0, 0x12, 0xaa,
	0x77, 0x64, 0x04, 0x9d, 0x7a, 0xb2, 0x0e, 0xef, 0x00, 0x4c, 0x20, 0xde, 0xcf, 0x9a, 0x32, 0xb6,
	0x63, 0x45, 0x2a, 0x26, 0x4c, 0xac, 0xbf, 0x55, 0x1d, 0x65, 0x7c, 0x5a, 0xf7, 0x0c, 0xde, 0xc2,
	0xbc, 0xb3, 0x9f, 0xd9, 0x7c, 0xc7, 0x8a, 0xa5, 0xf0, 0x30, 0x7f, 0x81, 0xf5, 0x68, 0x61, 0xb5,
	0x92, 0x96, 0x10, 0x61, 0xd1, 0xa8, 0x8f, 0x71, 0x7a, 0xc0, 0xa3, 0x8e, 0x47, 0x5d, 0x0f, 0xe9,
	0x6b, 0x6b, 0x1d, 0x49, 0x41, 0x27, 0x6f, 0x6b, 0xc9, 0x7c, 0x5f, 0x6c, 0x2f, 0xb1, 0x22, 0xe3,
	0x47, 0xca, 0x18, 0x28, 0x60, 0x5c, 0x03, 0x6f, 0x75, 0x48, 0x92, 0x0a, 0xde, 0x6a, 0xdf, 0xe3,
	0xfc, 0x52, 0x8b, 0xa1, 0xc7, 0xe3, 0xd1, 0xf6, 0x26, 0xda, 0xee, 0x21, 0x3d, 0x38, 0x43, 0x55,
	0x27, 0x6c, 0x94, 0xb0, 0xff, 0x92, 0x98, 0xf4, 0x49, 0xc1, 0xf2, 0x70, 0x96, 0xcd, 0x97, 0x51,
	0xb2, 0xfd, 0x21, 0x83, 0x7b, 0x48, 0x86, 0x8b, 0xe0, 0xa6, 0xd4, 0x75, 0xf9, 0xe7, 0xc0, 0x5b,
	0x9c, 0x52, 0xc3, 0x41, 0xf2, 0x19, 0xde, 0x43, 0x72, 0x0c, 0xcb, 0xe2, 0xca, 0xd7, 0xaf, 0x8b,
	0x6f, 0xc3, 0xf3, 0x1a, 0x28, 0x9f, 0x15, 0xec, 0x91, 0xd5, 0x49, 0xf8, 0xc0, 0xe7, 0xdf, 0x00,
	0x00, 0x00, 0xff, 0xff, 0x1f, 0xa5, 0x2f, 0x53, 0xd5, 0x01, 0x00, 0x00,
}
