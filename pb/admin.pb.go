// Code generated by protoc-gen-go. DO NOT EDIT.
// source: admin.proto

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

type GetKnownCurrencyRequest struct {
	Base                 *BaseRequest `protobuf:"bytes,1,opt,name=base,proto3" json:"base,omitempty"`
	Id                   int64        `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *GetKnownCurrencyRequest) Reset()         { *m = GetKnownCurrencyRequest{} }
func (m *GetKnownCurrencyRequest) String() string { return proto.CompactTextString(m) }
func (*GetKnownCurrencyRequest) ProtoMessage()    {}
func (*GetKnownCurrencyRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_admin_fa42f3cd8e7c79bf, []int{0}
}
func (m *GetKnownCurrencyRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetKnownCurrencyRequest.Unmarshal(m, b)
}
func (m *GetKnownCurrencyRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetKnownCurrencyRequest.Marshal(b, m, deterministic)
}
func (dst *GetKnownCurrencyRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetKnownCurrencyRequest.Merge(dst, src)
}
func (m *GetKnownCurrencyRequest) XXX_Size() int {
	return xxx_messageInfo_GetKnownCurrencyRequest.Size(m)
}
func (m *GetKnownCurrencyRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetKnownCurrencyRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetKnownCurrencyRequest proto.InternalMessageInfo

func (m *GetKnownCurrencyRequest) GetBase() *BaseRequest {
	if m != nil {
		return m.Base
	}
	return nil
}

func (m *GetKnownCurrencyRequest) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

type GetKnownCurrencyResponse struct {
	Id                   int64    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name                 string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	IssuerPublicKey      string   `protobuf:"bytes,3,opt,name=issuer_public_key,json=issuerPublicKey,proto3" json:"issuer_public_key,omitempty"`
	AssetCode            string   `protobuf:"bytes,4,opt,name=asset_code,json=assetCode,proto3" json:"asset_code,omitempty"`
	ShortDescription     string   `protobuf:"bytes,5,opt,name=short_description,json=shortDescription,proto3" json:"short_description,omitempty"`
	LongDescription      string   `protobuf:"bytes,6,opt,name=long_description,json=longDescription,proto3" json:"long_description,omitempty"`
	OrderIndex           int64    `protobuf:"varint,7,opt,name=order_index,json=orderIndex,proto3" json:"order_index,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetKnownCurrencyResponse) Reset()         { *m = GetKnownCurrencyResponse{} }
func (m *GetKnownCurrencyResponse) String() string { return proto.CompactTextString(m) }
func (*GetKnownCurrencyResponse) ProtoMessage()    {}
func (*GetKnownCurrencyResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_admin_fa42f3cd8e7c79bf, []int{1}
}
func (m *GetKnownCurrencyResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetKnownCurrencyResponse.Unmarshal(m, b)
}
func (m *GetKnownCurrencyResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetKnownCurrencyResponse.Marshal(b, m, deterministic)
}
func (dst *GetKnownCurrencyResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetKnownCurrencyResponse.Merge(dst, src)
}
func (m *GetKnownCurrencyResponse) XXX_Size() int {
	return xxx_messageInfo_GetKnownCurrencyResponse.Size(m)
}
func (m *GetKnownCurrencyResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetKnownCurrencyResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetKnownCurrencyResponse proto.InternalMessageInfo

func (m *GetKnownCurrencyResponse) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *GetKnownCurrencyResponse) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *GetKnownCurrencyResponse) GetIssuerPublicKey() string {
	if m != nil {
		return m.IssuerPublicKey
	}
	return ""
}

func (m *GetKnownCurrencyResponse) GetAssetCode() string {
	if m != nil {
		return m.AssetCode
	}
	return ""
}

func (m *GetKnownCurrencyResponse) GetShortDescription() string {
	if m != nil {
		return m.ShortDescription
	}
	return ""
}

func (m *GetKnownCurrencyResponse) GetLongDescription() string {
	if m != nil {
		return m.LongDescription
	}
	return ""
}

func (m *GetKnownCurrencyResponse) GetOrderIndex() int64 {
	if m != nil {
		return m.OrderIndex
	}
	return 0
}

type GetKnownCurrenciesResponse struct {
	Currencies           []*GetKnownCurrencyResponse `protobuf:"bytes,1,rep,name=currencies,proto3" json:"currencies,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                    `json:"-"`
	XXX_unrecognized     []byte                      `json:"-"`
	XXX_sizecache        int32                       `json:"-"`
}

func (m *GetKnownCurrenciesResponse) Reset()         { *m = GetKnownCurrenciesResponse{} }
func (m *GetKnownCurrenciesResponse) String() string { return proto.CompactTextString(m) }
func (*GetKnownCurrenciesResponse) ProtoMessage()    {}
func (*GetKnownCurrenciesResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_admin_fa42f3cd8e7c79bf, []int{2}
}
func (m *GetKnownCurrenciesResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetKnownCurrenciesResponse.Unmarshal(m, b)
}
func (m *GetKnownCurrenciesResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetKnownCurrenciesResponse.Marshal(b, m, deterministic)
}
func (dst *GetKnownCurrenciesResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetKnownCurrenciesResponse.Merge(dst, src)
}
func (m *GetKnownCurrenciesResponse) XXX_Size() int {
	return xxx_messageInfo_GetKnownCurrenciesResponse.Size(m)
}
func (m *GetKnownCurrenciesResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetKnownCurrenciesResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetKnownCurrenciesResponse proto.InternalMessageInfo

func (m *GetKnownCurrenciesResponse) GetCurrencies() []*GetKnownCurrencyResponse {
	if m != nil {
		return m.Currencies
	}
	return nil
}

func init() {
	proto.RegisterType((*GetKnownCurrencyRequest)(nil), "pb.GetKnownCurrencyRequest")
	proto.RegisterType((*GetKnownCurrencyResponse)(nil), "pb.GetKnownCurrencyResponse")
	proto.RegisterType((*GetKnownCurrenciesResponse)(nil), "pb.GetKnownCurrenciesResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// AdminApiServiceClient is the client API for AdminApiService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AdminApiServiceClient interface {
	GetKnownCurrency(ctx context.Context, in *GetKnownCurrencyRequest, opts ...grpc.CallOption) (*GetKnownCurrencyResponse, error)
	GetKnownCurrencies(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*GetKnownCurrenciesResponse, error)
}

type adminApiServiceClient struct {
	cc *grpc.ClientConn
}

func NewAdminApiServiceClient(cc *grpc.ClientConn) AdminApiServiceClient {
	return &adminApiServiceClient{cc}
}

func (c *adminApiServiceClient) GetKnownCurrency(ctx context.Context, in *GetKnownCurrencyRequest, opts ...grpc.CallOption) (*GetKnownCurrencyResponse, error) {
	out := new(GetKnownCurrencyResponse)
	err := c.cc.Invoke(ctx, "/pb.AdminApiService/GetKnownCurrency", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminApiServiceClient) GetKnownCurrencies(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*GetKnownCurrenciesResponse, error) {
	out := new(GetKnownCurrenciesResponse)
	err := c.cc.Invoke(ctx, "/pb.AdminApiService/GetKnownCurrencies", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AdminApiServiceServer is the server API for AdminApiService service.
type AdminApiServiceServer interface {
	GetKnownCurrency(context.Context, *GetKnownCurrencyRequest) (*GetKnownCurrencyResponse, error)
	GetKnownCurrencies(context.Context, *Empty) (*GetKnownCurrenciesResponse, error)
}

func RegisterAdminApiServiceServer(s *grpc.Server, srv AdminApiServiceServer) {
	s.RegisterService(&_AdminApiService_serviceDesc, srv)
}

func _AdminApiService_GetKnownCurrency_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetKnownCurrencyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminApiServiceServer).GetKnownCurrency(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.AdminApiService/GetKnownCurrency",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminApiServiceServer).GetKnownCurrency(ctx, req.(*GetKnownCurrencyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdminApiService_GetKnownCurrencies_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminApiServiceServer).GetKnownCurrencies(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.AdminApiService/GetKnownCurrencies",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminApiServiceServer).GetKnownCurrencies(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _AdminApiService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.AdminApiService",
	HandlerType: (*AdminApiServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetKnownCurrency",
			Handler:    _AdminApiService_GetKnownCurrency_Handler,
		},
		{
			MethodName: "GetKnownCurrencies",
			Handler:    _AdminApiService_GetKnownCurrencies_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "admin.proto",
}

func init() { proto.RegisterFile("admin.proto", fileDescriptor_admin_fa42f3cd8e7c79bf) }

var fileDescriptor_admin_fa42f3cd8e7c79bf = []byte{
	// 359 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x92, 0xcf, 0x4f, 0xea, 0x40,
	0x10, 0xc7, 0x69, 0xe1, 0xf1, 0xc2, 0xf4, 0xe5, 0x01, 0x7b, 0x79, 0x9b, 0x3e, 0x7f, 0x90, 0x7a,
	0x41, 0x4d, 0x38, 0xe0, 0xd5, 0x0b, 0xa2, 0x31, 0x86, 0x44, 0x4d, 0xbd, 0x79, 0x69, 0xda, 0xee,
	0x04, 0x37, 0x96, 0xdd, 0x75, 0xb7, 0x55, 0x7b, 0xf5, 0x7f, 0xf1, 0xff, 0x34, 0x5d, 0x14, 0x11,
	0xd0, 0x5b, 0xf3, 0x99, 0xcf, 0xcc, 0x66, 0xbe, 0x53, 0xf0, 0x62, 0x36, 0xe3, 0x62, 0xa0, 0xb4,
	0xcc, 0x25, 0x71, 0x55, 0xe2, 0xff, 0x99, 0x66, 0x32, 0x89, 0xb3, 0x39, 0x09, 0x2e, 0xe1, 0xdf,
	0x39, 0xe6, 0x13, 0x21, 0x9f, 0xc4, 0xb8, 0xd0, 0x1a, 0x45, 0x5a, 0x86, 0xf8, 0x50, 0xa0, 0xc9,
	0xc9, 0x1e, 0x34, 0x92, 0xd8, 0x20, 0x75, 0x7a, 0x4e, 0xdf, 0x1b, 0xb6, 0x07, 0x2a, 0x19, 0x9c,
	0xc4, 0x06, 0xdf, 0xcb, 0xa1, 0x2d, 0x92, 0xbf, 0xe0, 0x72, 0x46, 0xdd, 0x9e, 0xd3, 0xaf, 0x87,
	0x2e, 0x67, 0xc1, 0x8b, 0x0b, 0x74, 0x7d, 0xa0, 0x51, 0x52, 0x2c, 0x64, 0xe7, 0x43, 0x26, 0x04,
	0x1a, 0x22, 0x9e, 0xa1, 0x6d, 0x6f, 0x85, 0xf6, 0x9b, 0x1c, 0x40, 0x97, 0x1b, 0x53, 0xa0, 0x8e,
	0x54, 0x91, 0x64, 0x3c, 0x8d, 0xee, 0xb1, 0xa4, 0x75, 0x2b, 0xb4, 0xe7, 0x85, 0x6b, 0xcb, 0x27,
	0x58, 0x92, 0x6d, 0x80, 0xd8, 0x18, 0xcc, 0xa3, 0x54, 0x32, 0xa4, 0x0d, 0x2b, 0xb5, 0x2c, 0x19,
	0x4b, 0x86, 0xe4, 0x10, 0xba, 0xe6, 0x4e, 0xea, 0x3c, 0x62, 0x68, 0x52, 0xcd, 0x55, 0xce, 0xa5,
	0xa0, 0xbf, 0xac, 0xd5, 0xb1, 0x85, 0xd3, 0x4f, 0x4e, 0xf6, 0xa1, 0x93, 0x49, 0x31, 0xfd, 0xe2,
	0x36, 0xe7, 0xcf, 0x56, 0x7c, 0x59, 0xdd, 0x05, 0x4f, 0x6a, 0x86, 0x3a, 0xe2, 0x82, 0xe1, 0x33,
	0xfd, 0x6d, 0xf7, 0x01, 0x8b, 0x2e, 0x2a, 0x12, 0xdc, 0x82, 0xbf, 0x92, 0x01, 0x47, 0xb3, 0x48,
	0xe1, 0x18, 0x20, 0x5d, 0x50, 0xea, 0xf4, 0xea, 0x7d, 0x6f, 0xb8, 0x55, 0xa5, 0xfb, 0x5d, 0x6e,
	0xe1, 0x92, 0x3f, 0x7c, 0x75, 0xa0, 0x3d, 0xaa, 0x4e, 0x3a, 0x52, 0xfc, 0x06, 0xf5, 0x23, 0x4f,
	0x91, 0x5c, 0x41, 0x67, 0xb5, 0x97, 0xfc, 0xdf, 0x3c, 0xd1, 0xde, 0xce, 0xff, 0xf1, 0xb9, 0xa0,
	0x46, 0x46, 0x40, 0xd6, 0x17, 0x20, 0xad, 0xaa, 0xeb, 0x6c, 0xa6, 0xf2, 0xd2, 0xdf, 0xd9, 0x30,
	0x60, 0x69, 0xc7, 0xa0, 0x96, 0x34, 0xed, 0xff, 0x75, 0xf4, 0x16, 0x00, 0x00, 0xff, 0xff, 0x00,
	0xf3, 0xe5, 0x00, 0x80, 0x02, 0x00, 0x00,
}