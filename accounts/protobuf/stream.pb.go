// Code generated by protoc-gen-go. DO NOT EDIT.
// source: accounts/protobuf/stream.proto

package protobuf

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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type Account struct {
	Address              string            `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	PublicKey            string            `protobuf:"bytes,2,opt,name=public_key,json=publicKey,proto3" json:"public_key,omitempty"`
	Nonce                uint64            `protobuf:"varint,3,opt,name=nonce,proto3" json:"nonce,omitempty"`
	Balance              uint64            `protobuf:"varint,4,opt,name=balance,proto3" json:"balance,omitempty"`
	StorageRoot          string            `protobuf:"bytes,5,opt,name=storage_root,json=storageRoot,proto3" json:"storage_root,omitempty"`
	Balances             map[string]uint64 `protobuf:"bytes,6,rep,name=balances,proto3" json:"balances,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Account) Reset()         { *m = Account{} }
func (m *Account) String() string { return proto.CompactTextString(m) }
func (*Account) ProtoMessage()    {}
func (*Account) Descriptor() ([]byte, []int) {
	return fileDescriptor_a171dd46f4a5f190, []int{0}
}

func (m *Account) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Account.Unmarshal(m, b)
}
func (m *Account) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Account.Marshal(b, m, deterministic)
}
func (m *Account) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Account.Merge(m, src)
}
func (m *Account) XXX_Size() int {
	return xxx_messageInfo_Account.Size(m)
}
func (m *Account) XXX_DiscardUnknown() {
	xxx_messageInfo_Account.DiscardUnknown(m)
}

var xxx_messageInfo_Account proto.InternalMessageInfo

func (m *Account) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *Account) GetPublicKey() string {
	if m != nil {
		return m.PublicKey
	}
	return ""
}

func (m *Account) GetNonce() uint64 {
	if m != nil {
		return m.Nonce
	}
	return 0
}

func (m *Account) GetBalance() uint64 {
	if m != nil {
		return m.Balance
	}
	return 0
}

func (m *Account) GetStorageRoot() string {
	if m != nil {
		return m.StorageRoot
	}
	return ""
}

func (m *Account) GetBalances() map[string]uint64 {
	if m != nil {
		return m.Balances
	}
	return nil
}

type AccountRegisterRequest struct {
	PublicKey            string   `protobuf:"bytes,1,opt,name=public_key,json=publicKey,proto3" json:"public_key,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AccountRegisterRequest) Reset()         { *m = AccountRegisterRequest{} }
func (m *AccountRegisterRequest) String() string { return proto.CompactTextString(m) }
func (*AccountRegisterRequest) ProtoMessage()    {}
func (*AccountRegisterRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_a171dd46f4a5f190, []int{1}
}

func (m *AccountRegisterRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AccountRegisterRequest.Unmarshal(m, b)
}
func (m *AccountRegisterRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AccountRegisterRequest.Marshal(b, m, deterministic)
}
func (m *AccountRegisterRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AccountRegisterRequest.Merge(m, src)
}
func (m *AccountRegisterRequest) XXX_Size() int {
	return xxx_messageInfo_AccountRegisterRequest.Size(m)
}
func (m *AccountRegisterRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AccountRegisterRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AccountRegisterRequest proto.InternalMessageInfo

func (m *AccountRegisterRequest) GetPublicKey() string {
	if m != nil {
		return m.PublicKey
	}
	return ""
}

func init() {
	proto.RegisterType((*Account)(nil), "protobuf.Account")
	proto.RegisterMapType((map[string]uint64)(nil), "protobuf.Account.BalancesEntry")
	proto.RegisterType((*AccountRegisterRequest)(nil), "protobuf.AccountRegisterRequest")
}

func init() { proto.RegisterFile("accounts/protobuf/stream.proto", fileDescriptor_a171dd46f4a5f190) }

var fileDescriptor_a171dd46f4a5f190 = []byte{
	// 281 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x90, 0x4d, 0x4b, 0xf3, 0x40,
	0x10, 0xc7, 0xd9, 0x36, 0x7d, 0x9b, 0xf6, 0x81, 0xc7, 0x45, 0x64, 0x11, 0xd4, 0x58, 0x2f, 0x39,
	0xa5, 0xa0, 0x07, 0xc5, 0x9e, 0x0c, 0x08, 0x82, 0x97, 0xb2, 0x5f, 0xa0, 0x6c, 0x92, 0xb1, 0x06,
	0xd3, 0x6c, 0xdd, 0x17, 0x21, 0x9f, 0xc2, 0xaf, 0x2c, 0xbb, 0x9b, 0x08, 0xf6, 0xb6, 0xbf, 0xdf,
	0xcc, 0x0e, 0xff, 0x19, 0xb8, 0x14, 0x45, 0x21, 0x6d, 0x63, 0xf4, 0xea, 0xa0, 0xa4, 0x91, 0xb9,
	0x7d, 0x5b, 0x69, 0xa3, 0x50, 0xec, 0x53, 0xcf, 0x74, 0xda, 0xeb, 0xe5, 0xf7, 0x00, 0x26, 0x4f,
	0xa1, 0x99, 0x32, 0x98, 0x88, 0xb2, 0x54, 0xa8, 0x35, 0x23, 0x31, 0x49, 0x66, 0xbc, 0x47, 0x7a,
	0x01, 0x70, 0xb0, 0x79, 0x5d, 0x15, 0xdb, 0x0f, 0x6c, 0xd9, 0xc0, 0x17, 0x67, 0xc1, 0xbc, 0x62,
	0x4b, 0x4f, 0x61, 0xd4, 0xc8, 0xa6, 0x40, 0x36, 0x8c, 0x49, 0x12, 0xf1, 0x00, 0x6e, 0x5c, 0x2e,
	0x6a, 0xe1, 0x7c, 0xe4, 0x7d, 0x8f, 0xf4, 0x1a, 0x16, 0xda, 0x48, 0x25, 0x76, 0xb8, 0x55, 0x52,
	0x1a, 0x36, 0xf2, 0x03, 0xe7, 0x9d, 0xe3, 0x52, 0x1a, 0xba, 0x86, 0x69, 0xd7, 0xad, 0xd9, 0x38,
	0x1e, 0x26, 0xf3, 0xdb, 0xab, 0xb4, 0x0f, 0x9d, 0x76, 0x81, 0xd3, 0xac, 0xeb, 0x78, 0x6e, 0x8c,
	0x6a, 0xf9, 0xef, 0x87, 0xf3, 0x35, 0xfc, 0xfb, 0x53, 0xa2, 0xff, 0x61, 0xe8, 0x82, 0x87, 0xad,
	0xdc, 0xd3, 0x45, 0xfe, 0x12, 0xb5, 0x45, 0xbf, 0x4c, 0xc4, 0x03, 0x3c, 0x0e, 0x1e, 0xc8, 0xf2,
	0x1e, 0xce, 0xba, 0xf9, 0x1c, 0x77, 0x95, 0x36, 0xa8, 0x38, 0x7e, 0x5a, 0xd4, 0xe6, 0xe8, 0x0a,
	0xe4, 0xe8, 0x0a, 0xd9, 0x0d, 0x9c, 0x14, 0x72, 0x9f, 0xbe, 0xa3, 0x2a, 0x2b, 0xab, 0x43, 0xda,
	0x6c, 0xf1, 0x12, 0x70, 0xe3, 0x68, 0x43, 0xf2, 0xb1, 0xd7, 0x77, 0x3f, 0x01, 0x00, 0x00, 0xff,
	0xff, 0xbd, 0x03, 0x6c, 0x98, 0xa2, 0x01, 0x00, 0x00,
}
