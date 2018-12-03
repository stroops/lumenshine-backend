// Code generated by protoc-gen-go. DO NOT EDIT.
// source: global.proto

package pb

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
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

type EmailContentType int32

const (
	EmailContentType_text EmailContentType = 0
	EmailContentType_html EmailContentType = 1
)

var EmailContentType_name = map[int32]string{
	0: "text",
	1: "html",
}

var EmailContentType_value = map[string]int32{
	"text": 0,
	"html": 1,
}

func (x EmailContentType) String() string {
	return proto.EnumName(EmailContentType_name, int32(x))
}

func (EmailContentType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_4baa8fc7dedf329e, []int{0}
}

type NotificationType int32

const (
	NotificationType_ios     NotificationType = 0
	NotificationType_android NotificationType = 1
	NotificationType_mail    NotificationType = 2
)

var NotificationType_name = map[int32]string{
	0: "ios",
	1: "android",
	2: "mail",
}

var NotificationType_value = map[string]int32{
	"ios":     0,
	"android": 1,
	"mail":    2,
}

func (x NotificationType) String() string {
	return proto.EnumName(NotificationType_name, int32(x))
}

func (NotificationType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_4baa8fc7dedf329e, []int{1}
}

type NotificationStatusCode int32

const (
	NotificationStatusCode_success NotificationStatusCode = 0
	NotificationStatusCode_error   NotificationStatusCode = 1
)

var NotificationStatusCode_name = map[int32]string{
	0: "success",
	1: "error",
}

var NotificationStatusCode_value = map[string]int32{
	"success": 0,
	"error":   1,
}

func (x NotificationStatusCode) String() string {
	return proto.EnumName(NotificationStatusCode_name, int32(x))
}

func (NotificationStatusCode) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_4baa8fc7dedf329e, []int{2}
}

type DeviceType int32

const (
	DeviceType_apple  DeviceType = 0
	DeviceType_google DeviceType = 1
)

var DeviceType_name = map[int32]string{
	0: "apple",
	1: "google",
}

var DeviceType_value = map[string]int32{
	"apple":  0,
	"google": 1,
}

func (x DeviceType) String() string {
	return proto.EnumName(DeviceType_name, int32(x))
}

func (DeviceType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_4baa8fc7dedf329e, []int{3}
}

type DocumentType int32

const (
	DocumentType_passport           DocumentType = 0
	DocumentType_drivers_license    DocumentType = 1
	DocumentType_id_card            DocumentType = 2
	DocumentType_proof_of_residence DocumentType = 3
)

var DocumentType_name = map[int32]string{
	0: "passport",
	1: "drivers_license",
	2: "id_card",
	3: "proof_of_residence",
}

var DocumentType_value = map[string]int32{
	"passport":           0,
	"drivers_license":    1,
	"id_card":            2,
	"proof_of_residence": 3,
}

func (x DocumentType) String() string {
	return proto.EnumName(DocumentType_name, int32(x))
}

func (DocumentType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_4baa8fc7dedf329e, []int{4}
}

type DocumentFormat int32

const (
	DocumentFormat_png  DocumentFormat = 0
	DocumentFormat_pdf  DocumentFormat = 1
	DocumentFormat_jpg  DocumentFormat = 2
	DocumentFormat_jpeg DocumentFormat = 3
)

var DocumentFormat_name = map[int32]string{
	0: "png",
	1: "pdf",
	2: "jpg",
	3: "jpeg",
}

var DocumentFormat_value = map[string]int32{
	"png":  0,
	"pdf":  1,
	"jpg":  2,
	"jpeg": 3,
}

func (x DocumentFormat) String() string {
	return proto.EnumName(DocumentFormat_name, int32(x))
}

func (DocumentFormat) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_4baa8fc7dedf329e, []int{5}
}

type DocumentSide int32

const (
	DocumentSide_front DocumentSide = 0
	DocumentSide_back  DocumentSide = 1
)

var DocumentSide_name = map[int32]string{
	0: "front",
	1: "back",
}

var DocumentSide_value = map[string]int32{
	"front": 0,
	"back":  1,
}

func (x DocumentSide) String() string {
	return proto.EnumName(DocumentSide_name, int32(x))
}

func (DocumentSide) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_4baa8fc7dedf329e, []int{6}
}

type WalletType int32

const (
	WalletType_internal WalletType = 0
	WalletType_external WalletType = 1
)

var WalletType_name = map[int32]string{
	0: "internal",
	1: "external",
}

var WalletType_value = map[string]int32{
	"internal": 0,
	"external": 1,
}

func (x WalletType) String() string {
	return proto.EnumName(WalletType_name, int32(x))
}

func (WalletType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_4baa8fc7dedf329e, []int{7}
}

type NotificationParameterType int32

const (
	NotificationParameterType_ios_title_localized_key NotificationParameterType = 0
	NotificationParameterType_ios_category            NotificationParameterType = 1
	NotificationParameterType_ios_wallet_key          NotificationParameterType = 2
)

var NotificationParameterType_name = map[int32]string{
	0: "ios_title_localized_key",
	1: "ios_category",
	2: "ios_wallet_key",
}

var NotificationParameterType_value = map[string]int32{
	"ios_title_localized_key": 0,
	"ios_category":            1,
	"ios_wallet_key":          2,
}

func (x NotificationParameterType) String() string {
	return proto.EnumName(NotificationParameterType_name, int32(x))
}

func (NotificationParameterType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_4baa8fc7dedf329e, []int{8}
}

type BaseRequest struct {
	RequestId            string   `protobuf:"bytes,1,opt,name=request_id,json=requestId,proto3" json:"request_id,omitempty"`
	UpdateBy             string   `protobuf:"bytes,2,opt,name=update_by,json=updateBy,proto3" json:"update_by,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BaseRequest) Reset()         { *m = BaseRequest{} }
func (m *BaseRequest) String() string { return proto.CompactTextString(m) }
func (*BaseRequest) ProtoMessage()    {}
func (*BaseRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4baa8fc7dedf329e, []int{0}
}

func (m *BaseRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BaseRequest.Unmarshal(m, b)
}
func (m *BaseRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BaseRequest.Marshal(b, m, deterministic)
}
func (m *BaseRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BaseRequest.Merge(m, src)
}
func (m *BaseRequest) XXX_Size() int {
	return xxx_messageInfo_BaseRequest.Size(m)
}
func (m *BaseRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_BaseRequest.DiscardUnknown(m)
}

var xxx_messageInfo_BaseRequest proto.InternalMessageInfo

func (m *BaseRequest) GetRequestId() string {
	if m != nil {
		return m.RequestId
	}
	return ""
}

func (m *BaseRequest) GetUpdateBy() string {
	if m != nil {
		return m.UpdateBy
	}
	return ""
}

type IDRequest struct {
	Base                 *BaseRequest `protobuf:"bytes,1,opt,name=base,proto3" json:"base,omitempty"`
	Id                   int64        `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *IDRequest) Reset()         { *m = IDRequest{} }
func (m *IDRequest) String() string { return proto.CompactTextString(m) }
func (*IDRequest) ProtoMessage()    {}
func (*IDRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4baa8fc7dedf329e, []int{1}
}

func (m *IDRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_IDRequest.Unmarshal(m, b)
}
func (m *IDRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_IDRequest.Marshal(b, m, deterministic)
}
func (m *IDRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IDRequest.Merge(m, src)
}
func (m *IDRequest) XXX_Size() int {
	return xxx_messageInfo_IDRequest.Size(m)
}
func (m *IDRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_IDRequest.DiscardUnknown(m)
}

var xxx_messageInfo_IDRequest proto.InternalMessageInfo

func (m *IDRequest) GetBase() *BaseRequest {
	if m != nil {
		return m.Base
	}
	return nil
}

func (m *IDRequest) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

type IDResponse struct {
	Id                   int64    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *IDResponse) Reset()         { *m = IDResponse{} }
func (m *IDResponse) String() string { return proto.CompactTextString(m) }
func (*IDResponse) ProtoMessage()    {}
func (*IDResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_4baa8fc7dedf329e, []int{2}
}

func (m *IDResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_IDResponse.Unmarshal(m, b)
}
func (m *IDResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_IDResponse.Marshal(b, m, deterministic)
}
func (m *IDResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IDResponse.Merge(m, src)
}
func (m *IDResponse) XXX_Size() int {
	return xxx_messageInfo_IDResponse.Size(m)
}
func (m *IDResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_IDResponse.DiscardUnknown(m)
}

var xxx_messageInfo_IDResponse proto.InternalMessageInfo

func (m *IDResponse) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

type IDString struct {
	Base                 *BaseRequest `protobuf:"bytes,1,opt,name=base,proto3" json:"base,omitempty"`
	Id                   string       `protobuf:"bytes,2,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *IDString) Reset()         { *m = IDString{} }
func (m *IDString) String() string { return proto.CompactTextString(m) }
func (*IDString) ProtoMessage()    {}
func (*IDString) Descriptor() ([]byte, []int) {
	return fileDescriptor_4baa8fc7dedf329e, []int{3}
}

func (m *IDString) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_IDString.Unmarshal(m, b)
}
func (m *IDString) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_IDString.Marshal(b, m, deterministic)
}
func (m *IDString) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IDString.Merge(m, src)
}
func (m *IDString) XXX_Size() int {
	return xxx_messageInfo_IDString.Size(m)
}
func (m *IDString) XXX_DiscardUnknown() {
	xxx_messageInfo_IDString.DiscardUnknown(m)
}

var xxx_messageInfo_IDString proto.InternalMessageInfo

func (m *IDString) GetBase() *BaseRequest {
	if m != nil {
		return m.Base
	}
	return nil
}

func (m *IDString) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type KeyRequest struct {
	Base                 *BaseRequest `protobuf:"bytes,1,opt,name=base,proto3" json:"base,omitempty"`
	Key                  string       `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *KeyRequest) Reset()         { *m = KeyRequest{} }
func (m *KeyRequest) String() string { return proto.CompactTextString(m) }
func (*KeyRequest) ProtoMessage()    {}
func (*KeyRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4baa8fc7dedf329e, []int{4}
}

func (m *KeyRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_KeyRequest.Unmarshal(m, b)
}
func (m *KeyRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_KeyRequest.Marshal(b, m, deterministic)
}
func (m *KeyRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_KeyRequest.Merge(m, src)
}
func (m *KeyRequest) XXX_Size() int {
	return xxx_messageInfo_KeyRequest.Size(m)
}
func (m *KeyRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_KeyRequest.DiscardUnknown(m)
}

var xxx_messageInfo_KeyRequest proto.InternalMessageInfo

func (m *KeyRequest) GetBase() *BaseRequest {
	if m != nil {
		return m.Base
	}
	return nil
}

func (m *KeyRequest) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

type StringResponse struct {
	Value                string   `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StringResponse) Reset()         { *m = StringResponse{} }
func (m *StringResponse) String() string { return proto.CompactTextString(m) }
func (*StringResponse) ProtoMessage()    {}
func (*StringResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_4baa8fc7dedf329e, []int{5}
}

func (m *StringResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StringResponse.Unmarshal(m, b)
}
func (m *StringResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StringResponse.Marshal(b, m, deterministic)
}
func (m *StringResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StringResponse.Merge(m, src)
}
func (m *StringResponse) XXX_Size() int {
	return xxx_messageInfo_StringResponse.Size(m)
}
func (m *StringResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_StringResponse.DiscardUnknown(m)
}

var xxx_messageInfo_StringResponse proto.InternalMessageInfo

func (m *StringResponse) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

type BoolResponse struct {
	Value                bool     `protobuf:"varint,1,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BoolResponse) Reset()         { *m = BoolResponse{} }
func (m *BoolResponse) String() string { return proto.CompactTextString(m) }
func (*BoolResponse) ProtoMessage()    {}
func (*BoolResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_4baa8fc7dedf329e, []int{6}
}

func (m *BoolResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BoolResponse.Unmarshal(m, b)
}
func (m *BoolResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BoolResponse.Marshal(b, m, deterministic)
}
func (m *BoolResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BoolResponse.Merge(m, src)
}
func (m *BoolResponse) XXX_Size() int {
	return xxx_messageInfo_BoolResponse.Size(m)
}
func (m *BoolResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_BoolResponse.DiscardUnknown(m)
}

var xxx_messageInfo_BoolResponse proto.InternalMessageInfo

func (m *BoolResponse) GetValue() bool {
	if m != nil {
		return m.Value
	}
	return false
}

type IntResponse struct {
	Value                int64    `protobuf:"varint,1,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *IntResponse) Reset()         { *m = IntResponse{} }
func (m *IntResponse) String() string { return proto.CompactTextString(m) }
func (*IntResponse) ProtoMessage()    {}
func (*IntResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_4baa8fc7dedf329e, []int{7}
}

func (m *IntResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_IntResponse.Unmarshal(m, b)
}
func (m *IntResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_IntResponse.Marshal(b, m, deterministic)
}
func (m *IntResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IntResponse.Merge(m, src)
}
func (m *IntResponse) XXX_Size() int {
	return xxx_messageInfo_IntResponse.Size(m)
}
func (m *IntResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_IntResponse.DiscardUnknown(m)
}

var xxx_messageInfo_IntResponse proto.InternalMessageInfo

func (m *IntResponse) GetValue() int64 {
	if m != nil {
		return m.Value
	}
	return 0
}

type Empty struct {
	Base                 *BaseRequest `protobuf:"bytes,1,opt,name=base,proto3" json:"base,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *Empty) Reset()         { *m = Empty{} }
func (m *Empty) String() string { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()    {}
func (*Empty) Descriptor() ([]byte, []int) {
	return fileDescriptor_4baa8fc7dedf329e, []int{8}
}

func (m *Empty) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Empty.Unmarshal(m, b)
}
func (m *Empty) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Empty.Marshal(b, m, deterministic)
}
func (m *Empty) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Empty.Merge(m, src)
}
func (m *Empty) XXX_Size() int {
	return xxx_messageInfo_Empty.Size(m)
}
func (m *Empty) XXX_DiscardUnknown() {
	xxx_messageInfo_Empty.DiscardUnknown(m)
}

var xxx_messageInfo_Empty proto.InternalMessageInfo

func (m *Empty) GetBase() *BaseRequest {
	if m != nil {
		return m.Base
	}
	return nil
}

type KeyValueRequest struct {
	Base                 *BaseRequest `protobuf:"bytes,1,opt,name=base,proto3" json:"base,omitempty"`
	Key                  string       `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	Value                string       `protobuf:"bytes,3,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *KeyValueRequest) Reset()         { *m = KeyValueRequest{} }
func (m *KeyValueRequest) String() string { return proto.CompactTextString(m) }
func (*KeyValueRequest) ProtoMessage()    {}
func (*KeyValueRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_4baa8fc7dedf329e, []int{9}
}

func (m *KeyValueRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_KeyValueRequest.Unmarshal(m, b)
}
func (m *KeyValueRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_KeyValueRequest.Marshal(b, m, deterministic)
}
func (m *KeyValueRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_KeyValueRequest.Merge(m, src)
}
func (m *KeyValueRequest) XXX_Size() int {
	return xxx_messageInfo_KeyValueRequest.Size(m)
}
func (m *KeyValueRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_KeyValueRequest.DiscardUnknown(m)
}

var xxx_messageInfo_KeyValueRequest proto.InternalMessageInfo

func (m *KeyValueRequest) GetBase() *BaseRequest {
	if m != nil {
		return m.Base
	}
	return nil
}

func (m *KeyValueRequest) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *KeyValueRequest) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

func init() {
	proto.RegisterEnum("pb.EmailContentType", EmailContentType_name, EmailContentType_value)
	proto.RegisterEnum("pb.NotificationType", NotificationType_name, NotificationType_value)
	proto.RegisterEnum("pb.NotificationStatusCode", NotificationStatusCode_name, NotificationStatusCode_value)
	proto.RegisterEnum("pb.DeviceType", DeviceType_name, DeviceType_value)
	proto.RegisterEnum("pb.DocumentType", DocumentType_name, DocumentType_value)
	proto.RegisterEnum("pb.DocumentFormat", DocumentFormat_name, DocumentFormat_value)
	proto.RegisterEnum("pb.DocumentSide", DocumentSide_name, DocumentSide_value)
	proto.RegisterEnum("pb.WalletType", WalletType_name, WalletType_value)
	proto.RegisterEnum("pb.NotificationParameterType", NotificationParameterType_name, NotificationParameterType_value)
	proto.RegisterType((*BaseRequest)(nil), "pb.BaseRequest")
	proto.RegisterType((*IDRequest)(nil), "pb.IDRequest")
	proto.RegisterType((*IDResponse)(nil), "pb.IDResponse")
	proto.RegisterType((*IDString)(nil), "pb.IDString")
	proto.RegisterType((*KeyRequest)(nil), "pb.KeyRequest")
	proto.RegisterType((*StringResponse)(nil), "pb.StringResponse")
	proto.RegisterType((*BoolResponse)(nil), "pb.BoolResponse")
	proto.RegisterType((*IntResponse)(nil), "pb.IntResponse")
	proto.RegisterType((*Empty)(nil), "pb.Empty")
	proto.RegisterType((*KeyValueRequest)(nil), "pb.KeyValueRequest")
}

func init() { proto.RegisterFile("global.proto", fileDescriptor_4baa8fc7dedf329e) }

var fileDescriptor_4baa8fc7dedf329e = []byte{
	// 570 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x54, 0x4d, 0x6b, 0x1b, 0x31,
	0x10, 0xf5, 0xee, 0xc6, 0x89, 0x77, 0x6c, 0x1c, 0xb1, 0x2d, 0x69, 0x4a, 0x5a, 0x28, 0x76, 0x09,
	0x61, 0x29, 0xa1, 0xa4, 0xf4, 0xdc, 0x92, 0x38, 0x05, 0x13, 0x28, 0xc5, 0x29, 0xe9, 0x2d, 0x8b,
	0x76, 0x35, 0xde, 0x2a, 0x91, 0x25, 0x55, 0x92, 0xd3, 0x6c, 0x7f, 0x7d, 0x91, 0xfc, 0x11, 0x1f,
	0x72, 0x68, 0xe8, 0x6d, 0x34, 0xf3, 0xde, 0x9b, 0xf7, 0xe6, 0x20, 0xe8, 0xd5, 0x42, 0x95, 0x54,
	0x1c, 0x6b, 0xa3, 0x9c, 0xca, 0x62, 0x5d, 0x0e, 0xc6, 0xd0, 0x3d, 0xa5, 0x16, 0x27, 0xf8, 0x6b,
	0x8e, 0xd6, 0x65, 0xaf, 0x01, 0xcc, 0xa2, 0x2c, 0x38, 0xdb, 0x8f, 0xde, 0x44, 0x47, 0xe9, 0x24,
	0x5d, 0x76, 0xc6, 0x2c, 0x3b, 0x80, 0x74, 0xae, 0x19, 0x75, 0x58, 0x94, 0xcd, 0x7e, 0x1c, 0xa6,
	0x9d, 0x45, 0xe3, 0xb4, 0x19, 0x7c, 0x86, 0x74, 0x3c, 0x5a, 0x09, 0x0d, 0x61, 0xab, 0xa4, 0x16,
	0x83, 0x44, 0xf7, 0x64, 0xf7, 0x58, 0x97, 0xc7, 0x1b, 0x7b, 0x26, 0x61, 0x98, 0xf5, 0x21, 0xe6,
	0x2c, 0xe8, 0x24, 0x93, 0x98, 0xb3, 0xc1, 0x2b, 0x00, 0xaf, 0x60, 0xb5, 0x92, 0xeb, 0x69, 0xb4,
	0x9e, 0x7e, 0x82, 0xce, 0x78, 0x74, 0xe9, 0x0c, 0x97, 0xf5, 0x53, 0xe5, 0xd3, 0x20, 0x70, 0x06,
	0x70, 0x81, 0xcd, 0x93, 0x1c, 0x12, 0x48, 0x6e, 0x71, 0x15, 0xd5, 0x97, 0x83, 0x43, 0xe8, 0x2f,
	0x3c, 0xac, 0x7d, 0x3e, 0x87, 0xf6, 0x1d, 0x15, 0x73, 0x5c, 0x9e, 0x6b, 0xf1, 0x18, 0xbc, 0x85,
	0xde, 0xa9, 0x52, 0xe2, 0x71, 0x54, 0x67, 0x85, 0x1a, 0x42, 0x77, 0x2c, 0xdd, 0xe3, 0xa0, 0x64,
	0x05, 0x7a, 0x07, 0xed, 0xf3, 0x99, 0x76, 0xcd, 0x3f, 0x59, 0x1e, 0x5c, 0xc3, 0xee, 0x05, 0x36,
	0x57, 0x9e, 0xf9, 0x7f, 0x51, 0x1f, 0xdc, 0x24, 0x1b, 0xc1, 0xf2, 0x43, 0x20, 0xe7, 0x33, 0xca,
	0xc5, 0x99, 0x92, 0x0e, 0xa5, 0xfb, 0xde, 0x68, 0xcc, 0x3a, 0xb0, 0xe5, 0xf0, 0xde, 0x91, 0x96,
	0xaf, 0x7e, 0xba, 0x99, 0x20, 0x51, 0x7e, 0x02, 0xe4, 0xab, 0x72, 0x7c, 0xca, 0x2b, 0xea, 0xb8,
	0x92, 0x01, 0xb7, 0x03, 0x09, 0x57, 0x96, 0xb4, 0xb2, 0x2e, 0xec, 0x50, 0xc9, 0x8c, 0xe2, 0x8c,
	0x44, 0x9e, 0xe3, 0x05, 0x49, 0x9c, 0xbf, 0x87, 0xbd, 0x4d, 0xce, 0xa5, 0xa3, 0x6e, 0x6e, 0xcf,
	0x14, 0x43, 0x4f, 0xb0, 0xf3, 0xaa, 0x42, 0xeb, 0xd9, 0x29, 0xb4, 0xd1, 0x18, 0x65, 0x48, 0x94,
	0x0f, 0x01, 0x46, 0x78, 0xc7, 0x2b, 0x0c, 0xfa, 0x29, 0xb4, 0xa9, 0xd6, 0x02, 0x49, 0x2b, 0x03,
	0xd8, 0xae, 0x95, 0xaa, 0x05, 0x92, 0x28, 0xbf, 0x82, 0xde, 0x48, 0x55, 0xf3, 0xd9, 0xca, 0x6e,
	0x0f, 0x3a, 0x9a, 0x5a, 0xab, 0x95, 0xf1, 0x96, 0x9f, 0xc1, 0x2e, 0x33, 0xfc, 0x0e, 0x8d, 0x2d,
	0x04, 0xaf, 0x50, 0x5a, 0x24, 0x91, 0xdf, 0xc7, 0x59, 0x51, 0x51, 0xc3, 0x48, 0x9c, 0xed, 0x41,
	0xa6, 0x8d, 0x52, 0xd3, 0x42, 0x4d, 0x0b, 0x83, 0x96, 0x33, 0x94, 0x15, 0x92, 0x24, 0xff, 0x08,
	0xfd, 0x95, 0xee, 0x17, 0x65, 0x66, 0xd4, 0xf9, 0x80, 0x5a, 0xd6, 0xa4, 0x15, 0x0a, 0x36, 0x25,
	0x91, 0x2f, 0x6e, 0x74, 0x4d, 0x62, 0x9f, 0xf2, 0x46, 0x63, 0x4d, 0x92, 0x7c, 0xf8, 0x60, 0xe7,
	0x92, 0xb3, 0xe0, 0x7a, 0x6a, 0x94, 0x5c, 0x9e, 0xaf, 0xa4, 0xd5, 0x2d, 0x89, 0xf2, 0x23, 0x80,
	0x1f, 0x54, 0x08, 0x5c, 0x3b, 0xe6, 0xd2, 0xa1, 0x91, 0x54, 0x90, 0x96, 0x7f, 0xe1, 0xfd, 0xf2,
	0x15, 0xe5, 0xd7, 0xf0, 0x72, 0xf3, 0x68, 0xdf, 0xa8, 0xa1, 0x33, 0x74, 0x68, 0x02, 0xf1, 0x00,
	0x5e, 0x70, 0x65, 0x0b, 0xc7, 0x9d, 0xc0, 0x42, 0xa8, 0x8a, 0x0a, 0xfe, 0x07, 0x59, 0x71, 0x8b,
	0x0d, 0x69, 0x65, 0x04, 0x7a, 0x7e, 0x58, 0x51, 0x87, 0xb5, 0x32, 0x0d, 0x89, 0xb2, 0x0c, 0xfa,
	0xbe, 0xf3, 0x3b, 0x6c, 0x0e, 0xa8, 0xb8, 0xdc, 0x0e, 0xbf, 0xc5, 0x87, 0xbf, 0x01, 0x00, 0x00,
	0xff, 0xff, 0xf5, 0xe6, 0x58, 0xa0, 0x3d, 0x04, 0x00, 0x00,
}
