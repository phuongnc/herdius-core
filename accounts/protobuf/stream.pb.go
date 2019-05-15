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
	Address              string               `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	PublicKey            string               `protobuf:"bytes,2,opt,name=public_key,json=publicKey,proto3" json:"public_key,omitempty"`
	Nonce                uint64               `protobuf:"varint,3,opt,name=nonce,proto3" json:"nonce,omitempty"`
	Balance              uint64               `protobuf:"varint,4,opt,name=balance,proto3" json:"balance,omitempty"`
	StorageRoot          string               `protobuf:"bytes,5,opt,name=storage_root,json=storageRoot,proto3" json:"storage_root,omitempty"`
	Balances             map[string]uint64    `protobuf:"bytes,6,rep,name=balances,proto3" json:"balances,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	EBalances            map[string]*EBalance `protobuf:"bytes,7,rep,name=eBalances,proto3" json:"eBalances,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
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

func (m *Account) GetEBalances() map[string]*EBalance {
	if m != nil {
		return m.EBalances
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

type EBalance struct {
	Address              string   `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	Balance              uint64   `protobuf:"varint,2,opt,name=balance,proto3" json:"balance,omitempty"`
	LastBlockHeight      uint64   `protobuf:"varint,3,opt,name=last_block_height,json=lastBlockHeight,proto3" json:"last_block_height,omitempty"`
	Nonce                uint64   `protobuf:"varint,4,opt,name=nonce,proto3" json:"nonce,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EBalance) Reset()         { *m = EBalance{} }
func (m *EBalance) String() string { return proto.CompactTextString(m) }
func (*EBalance) ProtoMessage()    {}
func (*EBalance) Descriptor() ([]byte, []int) {
	return fileDescriptor_a171dd46f4a5f190, []int{2}
}

func (m *EBalance) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EBalance.Unmarshal(m, b)
}
func (m *EBalance) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EBalance.Marshal(b, m, deterministic)
}
func (m *EBalance) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EBalance.Merge(m, src)
}
func (m *EBalance) XXX_Size() int {
	return xxx_messageInfo_EBalance.Size(m)
}
func (m *EBalance) XXX_DiscardUnknown() {
	xxx_messageInfo_EBalance.DiscardUnknown(m)
}

var xxx_messageInfo_EBalance proto.InternalMessageInfo

func (m *EBalance) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *EBalance) GetBalance() uint64 {
	if m != nil {
		return m.Balance
	}
	return 0
}

func (m *EBalance) GetLastBlockHeight() uint64 {
	if m != nil {
		return m.LastBlockHeight
	}
	return 0
}

func (m *EBalance) GetNonce() uint64 {
	if m != nil {
		return m.Nonce
	}
	return 0
}

func init() {
	proto.RegisterType((*Account)(nil), "protobuf.Account")
	proto.RegisterMapType((map[string]uint64)(nil), "protobuf.Account.BalancesEntry")
	proto.RegisterMapType((map[string]*EBalance)(nil), "protobuf.Account.EBalancesEntry")
	proto.RegisterType((*AccountRegisterRequest)(nil), "protobuf.AccountRegisterRequest")
	proto.RegisterType((*EBalance)(nil), "protobuf.EBalance")
}

func init() { proto.RegisterFile("accounts/protobuf/stream.proto", fileDescriptor_a171dd46f4a5f190) }

var fileDescriptor_a171dd46f4a5f190 = []byte{
	// 367 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x52, 0xdd, 0x4a, 0xf3, 0x40,
	0x10, 0x25, 0x4d, 0xfa, 0x37, 0xed, 0xf7, 0xd3, 0xe5, 0xe3, 0x23, 0x14, 0xd4, 0x58, 0x6f, 0x82,
	0x17, 0x29, 0xd4, 0x0b, 0xc5, 0x82, 0x60, 0xa0, 0x50, 0xf0, 0xa6, 0xe4, 0x05, 0xc2, 0x26, 0x1d,
	0xdb, 0xd0, 0x34, 0x5b, 0x77, 0x37, 0x42, 0xef, 0x7c, 0x15, 0xdf, 0x54, 0x76, 0x93, 0xd8, 0xb4,
	0x82, 0xde, 0xed, 0x39, 0x33, 0x67, 0x38, 0xb3, 0x67, 0xe0, 0x9c, 0xc6, 0x31, 0xcb, 0x33, 0x29,
	0xc6, 0x3b, 0xce, 0x24, 0x8b, 0xf2, 0xe7, 0xb1, 0x90, 0x1c, 0xe9, 0xd6, 0xd3, 0x98, 0x74, 0x2a,
	0x7a, 0xf4, 0x6e, 0x42, 0xfb, 0xb1, 0x68, 0x26, 0x36, 0xb4, 0xe9, 0x72, 0xc9, 0x51, 0x08, 0xdb,
	0x70, 0x0c, 0xb7, 0x1b, 0x54, 0x90, 0x9c, 0x01, 0xec, 0xf2, 0x28, 0x4d, 0xe2, 0x70, 0x83, 0x7b,
	0xbb, 0xa1, 0x8b, 0xdd, 0x82, 0x79, 0xc2, 0x3d, 0xf9, 0x07, 0xcd, 0x8c, 0x65, 0x31, 0xda, 0xa6,
	0x63, 0xb8, 0x56, 0x50, 0x00, 0x35, 0x2e, 0xa2, 0x29, 0x55, 0xbc, 0xa5, 0xf9, 0x0a, 0x92, 0x4b,
	0xe8, 0x0b, 0xc9, 0x38, 0x5d, 0x61, 0xc8, 0x19, 0x93, 0x76, 0x53, 0x0f, 0xec, 0x95, 0x5c, 0xc0,
	0x98, 0x24, 0x53, 0xe8, 0x94, 0xdd, 0xc2, 0x6e, 0x39, 0xa6, 0xdb, 0x9b, 0x5c, 0x78, 0x95, 0x69,
	0xaf, 0x34, 0xec, 0xf9, 0x65, 0xc7, 0x2c, 0x93, 0x7c, 0x1f, 0x7c, 0x0a, 0xc8, 0x03, 0x74, 0xb1,
	0xaa, 0xd9, 0x6d, 0xad, 0x76, 0xbe, 0xaa, 0x67, 0xc7, 0xf2, 0x83, 0x64, 0x38, 0x85, 0x5f, 0x47,
	0x35, 0xf2, 0x17, 0x4c, 0xb5, 0x78, 0xf1, 0x2b, 0xea, 0xa9, 0x56, 0x7e, 0xa5, 0x69, 0x8e, 0xfa,
	0x33, 0xac, 0xa0, 0x00, 0xf7, 0x8d, 0x3b, 0x63, 0xb8, 0x80, 0xdf, 0xb3, 0x9f, 0xd4, 0x6e, 0x5d,
	0xdd, 0x9b, 0x90, 0x83, 0xb9, 0x4a, 0x5a, 0x9b, 0x38, 0xba, 0x85, 0xff, 0xa5, 0xe7, 0x00, 0x57,
	0x89, 0x90, 0xc8, 0x03, 0x7c, 0xc9, 0x51, 0xc8, 0x93, 0x5c, 0x8c, 0x93, 0x5c, 0x46, 0x6f, 0x06,
	0x74, 0xaa, 0x81, 0xdf, 0xa4, 0x5b, 0x0b, 0xaa, 0x71, 0x1c, 0xd4, 0x35, 0x0c, 0x52, 0x2a, 0x64,
	0x18, 0xa5, 0x2c, 0xde, 0x84, 0x6b, 0x4c, 0x56, 0x6b, 0x59, 0x86, 0xfc, 0x47, 0x15, 0x7c, 0xc5,
	0xcf, 0x35, 0x7d, 0x38, 0x02, 0xab, 0x76, 0x04, 0xfe, 0x15, 0x0c, 0x62, 0xb6, 0xf5, 0xd6, 0xc8,
	0x97, 0x49, 0x2e, 0x8a, 0x3d, 0xfd, 0xfe, 0xbc, 0x80, 0x0b, 0x85, 0x16, 0x46, 0xd4, 0xd2, 0xf4,
	0xcd, 0x47, 0x00, 0x00, 0x00, 0xff, 0xff, 0x34, 0xfd, 0xf2, 0x5c, 0xb7, 0x02, 0x00, 0x00,
}
