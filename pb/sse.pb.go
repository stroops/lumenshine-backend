// Code generated by protoc-gen-go. DO NOT EDIT.
// source: sse.proto

package pb

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

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

type SSEListenForRequest struct {
	Base                 *BaseRequest `protobuf:"bytes,1,opt,name=base,proto3" json:"base,omitempty"`
	OpTypes              int64        `protobuf:"varint,2,opt,name=op_types,json=opTypes,proto3" json:"op_types,omitempty"`
	SourceReciver        string       `protobuf:"bytes,3,opt,name=source_reciver,json=sourceReciver,proto3" json:"source_reciver,omitempty"`
	StellarAccount       string       `protobuf:"bytes,4,opt,name=stellar_account,json=stellarAccount,proto3" json:"stellar_account,omitempty"`
	WithResume           bool         `protobuf:"varint,5,opt,name=with_resume,json=withResume,proto3" json:"with_resume,omitempty"`
	ReturnData           bool         `protobuf:"varint,6,opt,name=return_data,json=returnData,proto3" json:"return_data,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *SSEListenForRequest) Reset()         { *m = SSEListenForRequest{} }
func (m *SSEListenForRequest) String() string { return proto.CompactTextString(m) }
func (*SSEListenForRequest) ProtoMessage()    {}
func (*SSEListenForRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_302dad79fff882c0, []int{0}
}

func (m *SSEListenForRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SSEListenForRequest.Unmarshal(m, b)
}
func (m *SSEListenForRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SSEListenForRequest.Marshal(b, m, deterministic)
}
func (m *SSEListenForRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SSEListenForRequest.Merge(m, src)
}
func (m *SSEListenForRequest) XXX_Size() int {
	return xxx_messageInfo_SSEListenForRequest.Size(m)
}
func (m *SSEListenForRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SSEListenForRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SSEListenForRequest proto.InternalMessageInfo

func (m *SSEListenForRequest) GetBase() *BaseRequest {
	if m != nil {
		return m.Base
	}
	return nil
}

func (m *SSEListenForRequest) GetOpTypes() int64 {
	if m != nil {
		return m.OpTypes
	}
	return 0
}

func (m *SSEListenForRequest) GetSourceReciver() string {
	if m != nil {
		return m.SourceReciver
	}
	return ""
}

func (m *SSEListenForRequest) GetStellarAccount() string {
	if m != nil {
		return m.StellarAccount
	}
	return ""
}

func (m *SSEListenForRequest) GetWithResume() bool {
	if m != nil {
		return m.WithResume
	}
	return false
}

func (m *SSEListenForRequest) GetReturnData() bool {
	if m != nil {
		return m.ReturnData
	}
	return false
}

type SSEGetDataRequest struct {
	Base                 *BaseRequest `protobuf:"bytes,1,opt,name=base,proto3" json:"base,omitempty"`
	SourceReciver        string       `protobuf:"bytes,2,opt,name=source_reciver,json=sourceReciver,proto3" json:"source_reciver,omitempty"`
	Count                int64        `protobuf:"varint,3,opt,name=count,proto3" json:"count,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *SSEGetDataRequest) Reset()         { *m = SSEGetDataRequest{} }
func (m *SSEGetDataRequest) String() string { return proto.CompactTextString(m) }
func (*SSEGetDataRequest) ProtoMessage()    {}
func (*SSEGetDataRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_302dad79fff882c0, []int{1}
}

func (m *SSEGetDataRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SSEGetDataRequest.Unmarshal(m, b)
}
func (m *SSEGetDataRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SSEGetDataRequest.Marshal(b, m, deterministic)
}
func (m *SSEGetDataRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SSEGetDataRequest.Merge(m, src)
}
func (m *SSEGetDataRequest) XXX_Size() int {
	return xxx_messageInfo_SSEGetDataRequest.Size(m)
}
func (m *SSEGetDataRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SSEGetDataRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SSEGetDataRequest proto.InternalMessageInfo

func (m *SSEGetDataRequest) GetBase() *BaseRequest {
	if m != nil {
		return m.Base
	}
	return nil
}

func (m *SSEGetDataRequest) GetSourceReciver() string {
	if m != nil {
		return m.SourceReciver
	}
	return ""
}

func (m *SSEGetDataRequest) GetCount() int64 {
	if m != nil {
		return m.Count
	}
	return 0
}

type SSEGetData struct {
	SseConfigId          int64    `protobuf:"varint,1,opt,name=sse_config_id,json=sseConfigId,proto3" json:"sse_config_id,omitempty"`
	SourceReceiver       string   `protobuf:"bytes,2,opt,name=source_receiver,json=sourceReceiver,proto3" json:"source_receiver,omitempty"`
	StellarAccount       string   `protobuf:"bytes,3,opt,name=stellar_account,json=stellarAccount,proto3" json:"stellar_account,omitempty"`
	OperationType        int64    `protobuf:"varint,4,opt,name=operation_type,json=operationType,proto3" json:"operation_type,omitempty"`
	OperationData        string   `protobuf:"bytes,5,opt,name=operation_data,json=operationData,proto3" json:"operation_data,omitempty"`
	TransactionId        int64    `protobuf:"varint,6,opt,name=transaction_id,json=transactionId,proto3" json:"transaction_id,omitempty"`
	OperationId          int64    `protobuf:"varint,7,opt,name=operation_id,json=operationId,proto3" json:"operation_id,omitempty"`
	LedgerId             int64    `protobuf:"varint,8,opt,name=ledger_id,json=ledgerId,proto3" json:"ledger_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SSEGetData) Reset()         { *m = SSEGetData{} }
func (m *SSEGetData) String() string { return proto.CompactTextString(m) }
func (*SSEGetData) ProtoMessage()    {}
func (*SSEGetData) Descriptor() ([]byte, []int) {
	return fileDescriptor_302dad79fff882c0, []int{2}
}

func (m *SSEGetData) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SSEGetData.Unmarshal(m, b)
}
func (m *SSEGetData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SSEGetData.Marshal(b, m, deterministic)
}
func (m *SSEGetData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SSEGetData.Merge(m, src)
}
func (m *SSEGetData) XXX_Size() int {
	return xxx_messageInfo_SSEGetData.Size(m)
}
func (m *SSEGetData) XXX_DiscardUnknown() {
	xxx_messageInfo_SSEGetData.DiscardUnknown(m)
}

var xxx_messageInfo_SSEGetData proto.InternalMessageInfo

func (m *SSEGetData) GetSseConfigId() int64 {
	if m != nil {
		return m.SseConfigId
	}
	return 0
}

func (m *SSEGetData) GetSourceReceiver() string {
	if m != nil {
		return m.SourceReceiver
	}
	return ""
}

func (m *SSEGetData) GetStellarAccount() string {
	if m != nil {
		return m.StellarAccount
	}
	return ""
}

func (m *SSEGetData) GetOperationType() int64 {
	if m != nil {
		return m.OperationType
	}
	return 0
}

func (m *SSEGetData) GetOperationData() string {
	if m != nil {
		return m.OperationData
	}
	return ""
}

func (m *SSEGetData) GetTransactionId() int64 {
	if m != nil {
		return m.TransactionId
	}
	return 0
}

func (m *SSEGetData) GetOperationId() int64 {
	if m != nil {
		return m.OperationId
	}
	return 0
}

func (m *SSEGetData) GetLedgerId() int64 {
	if m != nil {
		return m.LedgerId
	}
	return 0
}

type SSEGetDataResponse struct {
	Data                 []*SSEGetData `protobuf:"bytes,1,rep,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *SSEGetDataResponse) Reset()         { *m = SSEGetDataResponse{} }
func (m *SSEGetDataResponse) String() string { return proto.CompactTextString(m) }
func (*SSEGetDataResponse) ProtoMessage()    {}
func (*SSEGetDataResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_302dad79fff882c0, []int{3}
}

func (m *SSEGetDataResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SSEGetDataResponse.Unmarshal(m, b)
}
func (m *SSEGetDataResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SSEGetDataResponse.Marshal(b, m, deterministic)
}
func (m *SSEGetDataResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SSEGetDataResponse.Merge(m, src)
}
func (m *SSEGetDataResponse) XXX_Size() int {
	return xxx_messageInfo_SSEGetDataResponse.Size(m)
}
func (m *SSEGetDataResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_SSEGetDataResponse.DiscardUnknown(m)
}

var xxx_messageInfo_SSEGetDataResponse proto.InternalMessageInfo

func (m *SSEGetDataResponse) GetData() []*SSEGetData {
	if m != nil {
		return m.Data
	}
	return nil
}

type SSERemoveListeningRequest struct {
	Base                 *BaseRequest `protobuf:"bytes,1,opt,name=base,proto3" json:"base,omitempty"`
	SourceReciver        string       `protobuf:"bytes,2,opt,name=source_reciver,json=sourceReciver,proto3" json:"source_reciver,omitempty"`
	StellarAccount       string       `protobuf:"bytes,3,opt,name=stellar_account,json=stellarAccount,proto3" json:"stellar_account,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *SSERemoveListeningRequest) Reset()         { *m = SSERemoveListeningRequest{} }
func (m *SSERemoveListeningRequest) String() string { return proto.CompactTextString(m) }
func (*SSERemoveListeningRequest) ProtoMessage()    {}
func (*SSERemoveListeningRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_302dad79fff882c0, []int{4}
}

func (m *SSERemoveListeningRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SSERemoveListeningRequest.Unmarshal(m, b)
}
func (m *SSERemoveListeningRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SSERemoveListeningRequest.Marshal(b, m, deterministic)
}
func (m *SSERemoveListeningRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SSERemoveListeningRequest.Merge(m, src)
}
func (m *SSERemoveListeningRequest) XXX_Size() int {
	return xxx_messageInfo_SSERemoveListeningRequest.Size(m)
}
func (m *SSERemoveListeningRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SSERemoveListeningRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SSERemoveListeningRequest proto.InternalMessageInfo

func (m *SSERemoveListeningRequest) GetBase() *BaseRequest {
	if m != nil {
		return m.Base
	}
	return nil
}

func (m *SSERemoveListeningRequest) GetSourceReciver() string {
	if m != nil {
		return m.SourceReciver
	}
	return ""
}

func (m *SSERemoveListeningRequest) GetStellarAccount() string {
	if m != nil {
		return m.StellarAccount
	}
	return ""
}

type SSEClearSourceReciversRequest struct {
	Base                 *BaseRequest `protobuf:"bytes,1,opt,name=base,proto3" json:"base,omitempty"`
	SourceReciver        string       `protobuf:"bytes,2,opt,name=source_reciver,json=sourceReciver,proto3" json:"source_reciver,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *SSEClearSourceReciversRequest) Reset()         { *m = SSEClearSourceReciversRequest{} }
func (m *SSEClearSourceReciversRequest) String() string { return proto.CompactTextString(m) }
func (*SSEClearSourceReciversRequest) ProtoMessage()    {}
func (*SSEClearSourceReciversRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_302dad79fff882c0, []int{5}
}

func (m *SSEClearSourceReciversRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SSEClearSourceReciversRequest.Unmarshal(m, b)
}
func (m *SSEClearSourceReciversRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SSEClearSourceReciversRequest.Marshal(b, m, deterministic)
}
func (m *SSEClearSourceReciversRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SSEClearSourceReciversRequest.Merge(m, src)
}
func (m *SSEClearSourceReciversRequest) XXX_Size() int {
	return xxx_messageInfo_SSEClearSourceReciversRequest.Size(m)
}
func (m *SSEClearSourceReciversRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SSEClearSourceReciversRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SSEClearSourceReciversRequest proto.InternalMessageInfo

func (m *SSEClearSourceReciversRequest) GetBase() *BaseRequest {
	if m != nil {
		return m.Base
	}
	return nil
}

func (m *SSEClearSourceReciversRequest) GetSourceReciver() string {
	if m != nil {
		return m.SourceReciver
	}
	return ""
}

func init() {
	proto.RegisterType((*SSEListenForRequest)(nil), "pb.SSEListenForRequest")
	proto.RegisterType((*SSEGetDataRequest)(nil), "pb.SSEGetDataRequest")
	proto.RegisterType((*SSEGetData)(nil), "pb.SSEGetData")
	proto.RegisterType((*SSEGetDataResponse)(nil), "pb.SSEGetDataResponse")
	proto.RegisterType((*SSERemoveListeningRequest)(nil), "pb.SSERemoveListeningRequest")
	proto.RegisterType((*SSEClearSourceReciversRequest)(nil), "pb.SSEClearSourceReciversRequest")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// SSEServiceClient is the client API for SSEService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type SSEServiceClient interface {
	ListenFor(ctx context.Context, in *SSEListenForRequest, opts ...grpc.CallOption) (*Empty, error)
	RemoveListening(ctx context.Context, in *SSERemoveListeningRequest, opts ...grpc.CallOption) (*Empty, error)
	ClearSourceRecivers(ctx context.Context, in *SSEClearSourceReciversRequest, opts ...grpc.CallOption) (*Empty, error)
	GetData(ctx context.Context, in *SSEGetDataRequest, opts ...grpc.CallOption) (*SSEGetDataResponse, error)
}

type sSEServiceClient struct {
	cc *grpc.ClientConn
}

func NewSSEServiceClient(cc *grpc.ClientConn) SSEServiceClient {
	return &sSEServiceClient{cc}
}

func (c *sSEServiceClient) ListenFor(ctx context.Context, in *SSEListenForRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/pb.SSEService/ListenFor", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sSEServiceClient) RemoveListening(ctx context.Context, in *SSERemoveListeningRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/pb.SSEService/RemoveListening", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sSEServiceClient) ClearSourceRecivers(ctx context.Context, in *SSEClearSourceReciversRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/pb.SSEService/ClearSourceRecivers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *sSEServiceClient) GetData(ctx context.Context, in *SSEGetDataRequest, opts ...grpc.CallOption) (*SSEGetDataResponse, error) {
	out := new(SSEGetDataResponse)
	err := c.cc.Invoke(ctx, "/pb.SSEService/GetData", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SSEServiceServer is the server API for SSEService service.
type SSEServiceServer interface {
	ListenFor(context.Context, *SSEListenForRequest) (*Empty, error)
	RemoveListening(context.Context, *SSERemoveListeningRequest) (*Empty, error)
	ClearSourceRecivers(context.Context, *SSEClearSourceReciversRequest) (*Empty, error)
	GetData(context.Context, *SSEGetDataRequest) (*SSEGetDataResponse, error)
}

func RegisterSSEServiceServer(s *grpc.Server, srv SSEServiceServer) {
	s.RegisterService(&_SSEService_serviceDesc, srv)
}

func _SSEService_ListenFor_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SSEListenForRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SSEServiceServer).ListenFor(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.SSEService/ListenFor",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SSEServiceServer).ListenFor(ctx, req.(*SSEListenForRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SSEService_RemoveListening_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SSERemoveListeningRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SSEServiceServer).RemoveListening(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.SSEService/RemoveListening",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SSEServiceServer).RemoveListening(ctx, req.(*SSERemoveListeningRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SSEService_ClearSourceRecivers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SSEClearSourceReciversRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SSEServiceServer).ClearSourceRecivers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.SSEService/ClearSourceRecivers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SSEServiceServer).ClearSourceRecivers(ctx, req.(*SSEClearSourceReciversRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SSEService_GetData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SSEGetDataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SSEServiceServer).GetData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.SSEService/GetData",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SSEServiceServer).GetData(ctx, req.(*SSEGetDataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _SSEService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.SSEService",
	HandlerType: (*SSEServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListenFor",
			Handler:    _SSEService_ListenFor_Handler,
		},
		{
			MethodName: "RemoveListening",
			Handler:    _SSEService_RemoveListening_Handler,
		},
		{
			MethodName: "ClearSourceRecivers",
			Handler:    _SSEService_ClearSourceRecivers_Handler,
		},
		{
			MethodName: "GetData",
			Handler:    _SSEService_GetData_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "sse.proto",
}

func init() { proto.RegisterFile("sse.proto", fileDescriptor_302dad79fff882c0) }

var fileDescriptor_302dad79fff882c0 = []byte{
	// 534 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x94, 0xc1, 0x8e, 0xd3, 0x3c,
	0x14, 0x85, 0x27, 0x4d, 0x3b, 0x6d, 0x6f, 0xa7, 0xa9, 0x7e, 0xcf, 0x0f, 0x64, 0x06, 0x8d, 0xe8,
	0x04, 0x55, 0x74, 0x55, 0x89, 0xb2, 0x41, 0x48, 0x2c, 0x60, 0x08, 0xa8, 0x12, 0x2b, 0x87, 0x7d,
	0xe4, 0x26, 0x97, 0x62, 0x91, 0xc6, 0xc1, 0x76, 0x8a, 0x66, 0xc9, 0x1b, 0xf0, 0x36, 0x3c, 0x12,
	0xaf, 0x81, 0x62, 0x87, 0x4c, 0x5b, 0x3a, 0x12, 0x2c, 0x66, 0xe9, 0xef, 0xde, 0xf8, 0xfa, 0x9c,
	0x63, 0x07, 0xfa, 0x4a, 0xe1, 0xac, 0x90, 0x42, 0x0b, 0xd2, 0x2a, 0x96, 0xe7, 0x27, 0xab, 0x4c,
	0x2c, 0x59, 0x66, 0x49, 0xf0, 0xd3, 0x81, 0xd3, 0x28, 0x0a, 0xdf, 0x73, 0xa5, 0x31, 0x7f, 0x2b,
	0x24, 0xc5, 0x2f, 0x25, 0x2a, 0x4d, 0x1e, 0x43, 0x7b, 0xc9, 0x14, 0xfa, 0xce, 0xd8, 0x99, 0x0e,
	0xe6, 0xa3, 0x59, 0xb1, 0x9c, 0xbd, 0x66, 0x0a, 0xeb, 0x32, 0x35, 0x45, 0x72, 0x06, 0x3d, 0x51,
	0xc4, 0xfa, 0xba, 0x40, 0xe5, 0xb7, 0xc6, 0xce, 0xd4, 0xa5, 0x5d, 0x51, 0x7c, 0xa8, 0x96, 0x64,
	0x02, 0x9e, 0x12, 0xa5, 0x4c, 0x30, 0x96, 0x98, 0xf0, 0x0d, 0x4a, 0xdf, 0x1d, 0x3b, 0xd3, 0x3e,
	0x1d, 0x5a, 0x4a, 0x2d, 0x24, 0x4f, 0x60, 0xa4, 0x34, 0x66, 0x19, 0x93, 0x31, 0x4b, 0x12, 0x51,
	0xe6, 0xda, 0x6f, 0x9b, 0x3e, 0xaf, 0xc6, 0xaf, 0x2c, 0x25, 0x8f, 0x60, 0xf0, 0x95, 0xeb, 0x4f,
	0xb1, 0x44, 0x55, 0xae, 0xd1, 0xef, 0x8c, 0x9d, 0x69, 0x8f, 0x42, 0x85, 0xa8, 0x21, 0x55, 0x83,
	0x44, 0x5d, 0xca, 0x3c, 0x4e, 0x99, 0x66, 0xfe, 0xb1, 0x6d, 0xb0, 0xe8, 0x0d, 0xd3, 0x2c, 0x28,
	0xe1, 0xbf, 0x28, 0x0a, 0xdf, 0xa1, 0xae, 0x56, 0xff, 0x24, 0xf3, 0x4f, 0x2d, 0xad, 0x43, 0x5a,
	0xfe, 0x87, 0x8e, 0x55, 0xe0, 0x1a, 0x2b, 0xec, 0x22, 0xf8, 0xd1, 0x02, 0xb8, 0x99, 0x4b, 0x02,
	0x18, 0x2a, 0x85, 0x71, 0x22, 0xf2, 0x8f, 0x7c, 0x15, 0xf3, 0xd4, 0x4c, 0x76, 0xe9, 0x40, 0x29,
	0xbc, 0x32, 0x6c, 0x91, 0x1a, 0x53, 0x9a, 0x79, 0xb8, 0x35, 0xd0, 0x6b, 0x06, 0xe2, 0x6d, 0xee,
	0xb9, 0x07, 0xdd, 0x9b, 0x80, 0x27, 0x0a, 0x94, 0x4c, 0x73, 0x91, 0x9b, 0xbc, 0x8c, 0xcb, 0x2e,
	0x1d, 0x36, 0xb4, 0x4a, 0x6d, 0xb7, 0xcd, 0xd8, 0xd8, 0xb1, 0x42, 0x1b, 0x6a, 0x34, 0x4c, 0xc0,
	0xd3, 0x92, 0xe5, 0x8a, 0x25, 0xa6, 0x91, 0xa7, 0xc6, 0x6d, 0x97, 0x0e, 0xb7, 0xe8, 0x22, 0x25,
	0x97, 0x70, 0x72, 0xb3, 0x1b, 0x4f, 0xfd, 0xae, 0x55, 0xda, 0xb0, 0x45, 0x4a, 0x1e, 0x42, 0x3f,
	0xc3, 0x74, 0x85, 0xb2, 0xaa, 0xf7, 0x4c, 0xbd, 0x67, 0xc1, 0x22, 0x0d, 0x9e, 0x03, 0xd9, 0x0e,
	0x4c, 0x15, 0x22, 0x57, 0x48, 0x02, 0x68, 0x9b, 0x93, 0x39, 0x63, 0x77, 0x3a, 0x98, 0x7b, 0x55,
	0x62, 0x5b, 0x5d, 0xa6, 0x16, 0x7c, 0x77, 0xe0, 0x2c, 0x8a, 0x42, 0x8a, 0x6b, 0xb1, 0x41, 0x7b,
	0xb5, 0x79, 0xbe, 0xba, 0x8b, 0xcc, 0xff, 0x36, 0x81, 0xe0, 0x33, 0x5c, 0x44, 0x51, 0x78, 0x95,
	0x21, 0x93, 0xd1, 0xf6, 0x0e, 0xea, 0x0e, 0x4e, 0x35, 0xff, 0x66, 0xef, 0x5c, 0x84, 0x72, 0xc3,
	0x13, 0x24, 0x4f, 0xa1, 0xdf, 0xbc, 0x6f, 0xf2, 0xa0, 0x76, 0x6c, 0xff, 0xc5, 0x9f, 0xf7, 0xab,
	0x42, 0xb8, 0x2e, 0xf4, 0x75, 0x70, 0x44, 0x5e, 0xc2, 0x68, 0xcf, 0x3d, 0x72, 0x51, 0x7f, 0x78,
	0xd8, 0xd5, 0xdd, 0xcf, 0x43, 0x38, 0x3d, 0x20, 0x95, 0x5c, 0xd6, 0x5b, 0xdc, 0x6e, 0xc3, 0xee,
	0x36, 0x2f, 0xa0, 0xfb, 0xfb, 0xdd, 0xdc, 0xdb, 0x0b, 0xba, 0x6e, 0xbf, 0xbf, 0x8f, 0xed, 0x2d,
	0x09, 0x8e, 0x96, 0xc7, 0xe6, 0xff, 0xf6, 0xec, 0x57, 0x00, 0x00, 0x00, 0xff, 0xff, 0x1a, 0x9e,
	0xc9, 0x04, 0xfe, 0x04, 0x00, 0x00,
}
