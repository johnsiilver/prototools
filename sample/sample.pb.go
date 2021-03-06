// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.15.2
// source: sample.proto

package sample

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type EnumValues int32

const (
	EnumValues_EV_Unknown EnumValues = 0
	EnumValues_EV_Ok      EnumValues = 1
	EnumValues_EV_Not_Ok  EnumValues = 2
	EnumValues_EV_Eh      EnumValues = 3
)

// Enum value maps for EnumValues.
var (
	EnumValues_name = map[int32]string{
		0: "EV_Unknown",
		1: "EV_Ok",
		2: "EV_Not_Ok",
		3: "EV_Eh",
	}
	EnumValues_value = map[string]int32{
		"EV_Unknown": 0,
		"EV_Ok":      1,
		"EV_Not_Ok":  2,
		"EV_Eh":      3,
	}
)

func (x EnumValues) Enum() *EnumValues {
	p := new(EnumValues)
	*p = x
	return p
}

func (x EnumValues) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (EnumValues) Descriptor() protoreflect.EnumDescriptor {
	return file_sample_proto_enumTypes[0].Descriptor()
}

func (EnumValues) Type() protoreflect.EnumType {
	return &file_sample_proto_enumTypes[0]
}

func (x EnumValues) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use EnumValues.Descriptor instead.
func (EnumValues) EnumDescriptor() ([]byte, []int) {
	return file_sample_proto_rawDescGZIP(), []int{0}
}

type Layer0_EnumEmbedded int32

const (
	Layer0_EE_UNKNOWN  Layer0_EnumEmbedded = 0
	Layer0_EE_WHATEVER Layer0_EnumEmbedded = 1
)

// Enum value maps for Layer0_EnumEmbedded.
var (
	Layer0_EnumEmbedded_name = map[int32]string{
		0: "EE_UNKNOWN",
		1: "EE_WHATEVER",
	}
	Layer0_EnumEmbedded_value = map[string]int32{
		"EE_UNKNOWN":  0,
		"EE_WHATEVER": 1,
	}
)

func (x Layer0_EnumEmbedded) Enum() *Layer0_EnumEmbedded {
	p := new(Layer0_EnumEmbedded)
	*p = x
	return p
}

func (x Layer0_EnumEmbedded) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Layer0_EnumEmbedded) Descriptor() protoreflect.EnumDescriptor {
	return file_sample_proto_enumTypes[1].Descriptor()
}

func (Layer0_EnumEmbedded) Type() protoreflect.EnumType {
	return &file_sample_proto_enumTypes[1]
}

func (x Layer0_EnumEmbedded) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Layer0_EnumEmbedded.Descriptor instead.
func (Layer0_EnumEmbedded) EnumDescriptor() ([]byte, []int) {
	return file_sample_proto_rawDescGZIP(), []int{1, 0}
}

type Supported struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ev      EnumValues `protobuf:"varint,1,opt,name=ev,proto3,enum=r3.EnumValues" json:"ev,omitempty"`
	Vstring string     `protobuf:"bytes,2,opt,name=vstring,proto3" json:"vstring,omitempty"`
	Vint32  int32      `protobuf:"varint,3,opt,name=vint32,proto3" json:"vint32,omitempty"`
	Vint64  int64      `protobuf:"varint,4,opt,name=vint64,proto3" json:"vint64,omitempty"`
	Vbool   bool       `protobuf:"varint,5,opt,name=vbool,proto3" json:"vbool,omitempty"`
	VTime   int64      `protobuf:"varint,6,opt,name=v_time,json=vTime,proto3" json:"v_time,omitempty"`
	Vfloat  float32    `protobuf:"fixed32,7,opt,name=vfloat,proto3" json:"vfloat,omitempty"`
	Vdouble float64    `protobuf:"fixed64,8,opt,name=vdouble,proto3" json:"vdouble,omitempty"`
}

func (x *Supported) Reset() {
	*x = Supported{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sample_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Supported) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Supported) ProtoMessage() {}

func (x *Supported) ProtoReflect() protoreflect.Message {
	mi := &file_sample_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Supported.ProtoReflect.Descriptor instead.
func (*Supported) Descriptor() ([]byte, []int) {
	return file_sample_proto_rawDescGZIP(), []int{0}
}

func (x *Supported) GetEv() EnumValues {
	if x != nil {
		return x.Ev
	}
	return EnumValues_EV_Unknown
}

func (x *Supported) GetVstring() string {
	if x != nil {
		return x.Vstring
	}
	return ""
}

func (x *Supported) GetVint32() int32 {
	if x != nil {
		return x.Vint32
	}
	return 0
}

func (x *Supported) GetVint64() int64 {
	if x != nil {
		return x.Vint64
	}
	return 0
}

func (x *Supported) GetVbool() bool {
	if x != nil {
		return x.Vbool
	}
	return false
}

func (x *Supported) GetVTime() int64 {
	if x != nil {
		return x.VTime
	}
	return 0
}

func (x *Supported) GetVfloat() float32 {
	if x != nil {
		return x.Vfloat
	}
	return 0
}

func (x *Supported) GetVdouble() float64 {
	if x != nil {
		return x.Vdouble
	}
	return 0
}

type Layer0 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Layer1 *Layer1             `protobuf:"bytes,1,opt,name=layer1,proto3" json:"layer1,omitempty"`
	Vint32 int32               `protobuf:"varint,2,opt,name=vint32,proto3" json:"vint32,omitempty"`
	Ee     Layer0_EnumEmbedded `protobuf:"varint,3,opt,name=ee,proto3,enum=r3.Layer0_EnumEmbedded" json:"ee,omitempty"`
}

func (x *Layer0) Reset() {
	*x = Layer0{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sample_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Layer0) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Layer0) ProtoMessage() {}

func (x *Layer0) ProtoReflect() protoreflect.Message {
	mi := &file_sample_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Layer0.ProtoReflect.Descriptor instead.
func (*Layer0) Descriptor() ([]byte, []int) {
	return file_sample_proto_rawDescGZIP(), []int{1}
}

func (x *Layer0) GetLayer1() *Layer1 {
	if x != nil {
		return x.Layer1
	}
	return nil
}

func (x *Layer0) GetVint32() int32 {
	if x != nil {
		return x.Vint32
	}
	return 0
}

func (x *Layer0) GetEe() Layer0_EnumEmbedded {
	if x != nil {
		return x.Ee
	}
	return Layer0_EE_UNKNOWN
}

type Layer1 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Supported *Supported `protobuf:"bytes,1,opt,name=supported,proto3" json:"supported,omitempty"`
	Vstring   string     `protobuf:"bytes,2,opt,name=vstring,proto3" json:"vstring,omitempty"`
}

func (x *Layer1) Reset() {
	*x = Layer1{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sample_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Layer1) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Layer1) ProtoMessage() {}

func (x *Layer1) ProtoReflect() protoreflect.Message {
	mi := &file_sample_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Layer1.ProtoReflect.Descriptor instead.
func (*Layer1) Descriptor() ([]byte, []int) {
	return file_sample_proto_rawDescGZIP(), []int{2}
}

func (x *Layer1) GetSupported() *Supported {
	if x != nil {
		return x.Supported
	}
	return nil
}

func (x *Layer1) GetVstring() string {
	if x != nil {
		return x.Vstring
	}
	return ""
}

type BunchOTypes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ev       EnumValues   `protobuf:"varint,1,opt,name=ev,proto3,enum=r3.EnumValues" json:"ev,omitempty"`
	Vstring  string       `protobuf:"bytes,2,opt,name=vstring,proto3" json:"vstring,omitempty"`
	Vint32   int32        `protobuf:"varint,3,opt,name=vint32,proto3" json:"vint32,omitempty"`
	Vint64   int64        `protobuf:"varint,4,opt,name=vint64,proto3" json:"vint64,omitempty"`
	Vbool    bool         `protobuf:"varint,5,opt,name=vbool,proto3" json:"vbool,omitempty"`
	VTime    int64        `protobuf:"varint,6,opt,name=v_time,json=vTime,proto3" json:"v_time,omitempty"`
	Vfloat   float32      `protobuf:"fixed32,7,opt,name=vfloat,proto3" json:"vfloat,omitempty"`
	Vdouble  float64      `protobuf:"fixed64,8,opt,name=vdouble,proto3" json:"vdouble,omitempty"`
	LEv      []EnumValues `protobuf:"varint,9,rep,packed,name=l_ev,json=lEv,proto3,enum=r3.EnumValues" json:"l_ev,omitempty"`
	LString  []string     `protobuf:"bytes,10,rep,name=l_string,json=lString,proto3" json:"l_string,omitempty"`
	LInt32   []int32      `protobuf:"varint,11,rep,packed,name=l_int32,json=lInt32,proto3" json:"l_int32,omitempty"`
	LInt64   []int64      `protobuf:"varint,12,rep,packed,name=l_int64,json=lInt64,proto3" json:"l_int64,omitempty"`
	LBool    []bool       `protobuf:"varint,13,rep,packed,name=l_bool,json=lBool,proto3" json:"l_bool,omitempty"`
	LFloat   []float32    `protobuf:"fixed32,14,rep,packed,name=l_float,json=lFloat,proto3" json:"l_float,omitempty"`
	LDouble  []float64    `protobuf:"fixed64,15,rep,packed,name=l_double,json=lDouble,proto3" json:"l_double,omitempty"`
	LMessage []*Supported `protobuf:"bytes,16,rep,name=l_message,json=lMessage,proto3" json:"l_message,omitempty"`
}

func (x *BunchOTypes) Reset() {
	*x = BunchOTypes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_sample_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BunchOTypes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BunchOTypes) ProtoMessage() {}

func (x *BunchOTypes) ProtoReflect() protoreflect.Message {
	mi := &file_sample_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BunchOTypes.ProtoReflect.Descriptor instead.
func (*BunchOTypes) Descriptor() ([]byte, []int) {
	return file_sample_proto_rawDescGZIP(), []int{3}
}

func (x *BunchOTypes) GetEv() EnumValues {
	if x != nil {
		return x.Ev
	}
	return EnumValues_EV_Unknown
}

func (x *BunchOTypes) GetVstring() string {
	if x != nil {
		return x.Vstring
	}
	return ""
}

func (x *BunchOTypes) GetVint32() int32 {
	if x != nil {
		return x.Vint32
	}
	return 0
}

func (x *BunchOTypes) GetVint64() int64 {
	if x != nil {
		return x.Vint64
	}
	return 0
}

func (x *BunchOTypes) GetVbool() bool {
	if x != nil {
		return x.Vbool
	}
	return false
}

func (x *BunchOTypes) GetVTime() int64 {
	if x != nil {
		return x.VTime
	}
	return 0
}

func (x *BunchOTypes) GetVfloat() float32 {
	if x != nil {
		return x.Vfloat
	}
	return 0
}

func (x *BunchOTypes) GetVdouble() float64 {
	if x != nil {
		return x.Vdouble
	}
	return 0
}

func (x *BunchOTypes) GetLEv() []EnumValues {
	if x != nil {
		return x.LEv
	}
	return nil
}

func (x *BunchOTypes) GetLString() []string {
	if x != nil {
		return x.LString
	}
	return nil
}

func (x *BunchOTypes) GetLInt32() []int32 {
	if x != nil {
		return x.LInt32
	}
	return nil
}

func (x *BunchOTypes) GetLInt64() []int64 {
	if x != nil {
		return x.LInt64
	}
	return nil
}

func (x *BunchOTypes) GetLBool() []bool {
	if x != nil {
		return x.LBool
	}
	return nil
}

func (x *BunchOTypes) GetLFloat() []float32 {
	if x != nil {
		return x.LFloat
	}
	return nil
}

func (x *BunchOTypes) GetLDouble() []float64 {
	if x != nil {
		return x.LDouble
	}
	return nil
}

func (x *BunchOTypes) GetLMessage() []*Supported {
	if x != nil {
		return x.LMessage
	}
	return nil
}

var File_sample_proto protoreflect.FileDescriptor

var file_sample_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x73, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02,
	0x72, 0x33, 0x22, 0xd4, 0x01, 0x0a, 0x09, 0x53, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x65, 0x64,
	0x12, 0x1e, 0x0a, 0x02, 0x65, 0x76, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0e, 0x2e, 0x72,
	0x33, 0x2e, 0x45, 0x6e, 0x75, 0x6d, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x52, 0x02, 0x65, 0x76,
	0x12, 0x18, 0x0a, 0x07, 0x76, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x76, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x12, 0x16, 0x0a, 0x06, 0x76, 0x69,
	0x6e, 0x74, 0x33, 0x32, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x76, 0x69, 0x6e, 0x74,
	0x33, 0x32, 0x12, 0x16, 0x0a, 0x06, 0x76, 0x69, 0x6e, 0x74, 0x36, 0x34, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x06, 0x76, 0x69, 0x6e, 0x74, 0x36, 0x34, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x62,
	0x6f, 0x6f, 0x6c, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x76, 0x62, 0x6f, 0x6f, 0x6c,
	0x12, 0x15, 0x0a, 0x06, 0x76, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x05, 0x76, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x76, 0x66, 0x6c, 0x6f, 0x61,
	0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x02, 0x52, 0x06, 0x76, 0x66, 0x6c, 0x6f, 0x61, 0x74, 0x12,
	0x18, 0x0a, 0x07, 0x76, 0x64, 0x6f, 0x75, 0x62, 0x6c, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x01,
	0x52, 0x07, 0x76, 0x64, 0x6f, 0x75, 0x62, 0x6c, 0x65, 0x22, 0x9e, 0x01, 0x0a, 0x06, 0x4c, 0x61,
	0x79, 0x65, 0x72, 0x30, 0x12, 0x22, 0x0a, 0x06, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x31, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x72, 0x33, 0x2e, 0x4c, 0x61, 0x79, 0x65, 0x72, 0x31,
	0x52, 0x06, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x31, 0x12, 0x16, 0x0a, 0x06, 0x76, 0x69, 0x6e, 0x74,
	0x33, 0x32, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x76, 0x69, 0x6e, 0x74, 0x33, 0x32,
	0x12, 0x27, 0x0a, 0x02, 0x65, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x17, 0x2e, 0x72,
	0x33, 0x2e, 0x4c, 0x61, 0x79, 0x65, 0x72, 0x30, 0x2e, 0x45, 0x6e, 0x75, 0x6d, 0x45, 0x6d, 0x62,
	0x65, 0x64, 0x64, 0x65, 0x64, 0x52, 0x02, 0x65, 0x65, 0x22, 0x2f, 0x0a, 0x0c, 0x45, 0x6e, 0x75,
	0x6d, 0x45, 0x6d, 0x62, 0x65, 0x64, 0x64, 0x65, 0x64, 0x12, 0x0e, 0x0a, 0x0a, 0x45, 0x45, 0x5f,
	0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x0f, 0x0a, 0x0b, 0x45, 0x45, 0x5f,
	0x57, 0x48, 0x41, 0x54, 0x45, 0x56, 0x45, 0x52, 0x10, 0x01, 0x22, 0x4f, 0x0a, 0x06, 0x4c, 0x61,
	0x79, 0x65, 0x72, 0x31, 0x12, 0x2b, 0x0a, 0x09, 0x73, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x65,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x72, 0x33, 0x2e, 0x53, 0x75, 0x70,
	0x70, 0x6f, 0x72, 0x74, 0x65, 0x64, 0x52, 0x09, 0x73, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x65,
	0x64, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x76, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x22, 0xbd, 0x03, 0x0a, 0x0b,
	0x42, 0x75, 0x6e, 0x63, 0x68, 0x4f, 0x54, 0x79, 0x70, 0x65, 0x73, 0x12, 0x1e, 0x0a, 0x02, 0x65,
	0x76, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0e, 0x2e, 0x72, 0x33, 0x2e, 0x45, 0x6e, 0x75,
	0x6d, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x52, 0x02, 0x65, 0x76, 0x12, 0x18, 0x0a, 0x07, 0x76,
	0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x73,
	0x74, 0x72, 0x69, 0x6e, 0x67, 0x12, 0x16, 0x0a, 0x06, 0x76, 0x69, 0x6e, 0x74, 0x33, 0x32, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x76, 0x69, 0x6e, 0x74, 0x33, 0x32, 0x12, 0x16, 0x0a,
	0x06, 0x76, 0x69, 0x6e, 0x74, 0x36, 0x34, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x76,
	0x69, 0x6e, 0x74, 0x36, 0x34, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x62, 0x6f, 0x6f, 0x6c, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x05, 0x76, 0x62, 0x6f, 0x6f, 0x6c, 0x12, 0x15, 0x0a, 0x06, 0x76,
	0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x76, 0x54, 0x69,
	0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x76, 0x66, 0x6c, 0x6f, 0x61, 0x74, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x02, 0x52, 0x06, 0x76, 0x66, 0x6c, 0x6f, 0x61, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x76, 0x64,
	0x6f, 0x75, 0x62, 0x6c, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x01, 0x52, 0x07, 0x76, 0x64, 0x6f,
	0x75, 0x62, 0x6c, 0x65, 0x12, 0x21, 0x0a, 0x04, 0x6c, 0x5f, 0x65, 0x76, 0x18, 0x09, 0x20, 0x03,
	0x28, 0x0e, 0x32, 0x0e, 0x2e, 0x72, 0x33, 0x2e, 0x45, 0x6e, 0x75, 0x6d, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x73, 0x52, 0x03, 0x6c, 0x45, 0x76, 0x12, 0x19, 0x0a, 0x08, 0x6c, 0x5f, 0x73, 0x74, 0x72,
	0x69, 0x6e, 0x67, 0x18, 0x0a, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x6c, 0x53, 0x74, 0x72, 0x69,
	0x6e, 0x67, 0x12, 0x17, 0x0a, 0x07, 0x6c, 0x5f, 0x69, 0x6e, 0x74, 0x33, 0x32, 0x18, 0x0b, 0x20,
	0x03, 0x28, 0x05, 0x52, 0x06, 0x6c, 0x49, 0x6e, 0x74, 0x33, 0x32, 0x12, 0x17, 0x0a, 0x07, 0x6c,
	0x5f, 0x69, 0x6e, 0x74, 0x36, 0x34, 0x18, 0x0c, 0x20, 0x03, 0x28, 0x03, 0x52, 0x06, 0x6c, 0x49,
	0x6e, 0x74, 0x36, 0x34, 0x12, 0x15, 0x0a, 0x06, 0x6c, 0x5f, 0x62, 0x6f, 0x6f, 0x6c, 0x18, 0x0d,
	0x20, 0x03, 0x28, 0x08, 0x52, 0x05, 0x6c, 0x42, 0x6f, 0x6f, 0x6c, 0x12, 0x17, 0x0a, 0x07, 0x6c,
	0x5f, 0x66, 0x6c, 0x6f, 0x61, 0x74, 0x18, 0x0e, 0x20, 0x03, 0x28, 0x02, 0x52, 0x06, 0x6c, 0x46,
	0x6c, 0x6f, 0x61, 0x74, 0x12, 0x19, 0x0a, 0x08, 0x6c, 0x5f, 0x64, 0x6f, 0x75, 0x62, 0x6c, 0x65,
	0x18, 0x0f, 0x20, 0x03, 0x28, 0x01, 0x52, 0x07, 0x6c, 0x44, 0x6f, 0x75, 0x62, 0x6c, 0x65, 0x12,
	0x2a, 0x0a, 0x09, 0x6c, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x10, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x72, 0x33, 0x2e, 0x53, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x65,
	0x64, 0x52, 0x08, 0x6c, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2a, 0x41, 0x0a, 0x0a, 0x45,
	0x6e, 0x75, 0x6d, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x12, 0x0e, 0x0a, 0x0a, 0x45, 0x56, 0x5f,
	0x55, 0x6e, 0x6b, 0x6e, 0x6f, 0x77, 0x6e, 0x10, 0x00, 0x12, 0x09, 0x0a, 0x05, 0x45, 0x56, 0x5f,
	0x4f, 0x6b, 0x10, 0x01, 0x12, 0x0d, 0x0a, 0x09, 0x45, 0x56, 0x5f, 0x4e, 0x6f, 0x74, 0x5f, 0x4f,
	0x6b, 0x10, 0x02, 0x12, 0x09, 0x0a, 0x05, 0x45, 0x56, 0x5f, 0x45, 0x68, 0x10, 0x03, 0x42, 0x2a,
	0x5a, 0x28, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6a, 0x6f, 0x68,
	0x6e, 0x73, 0x69, 0x69, 0x6c, 0x76, 0x65, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x74, 0x6f,
	0x6f, 0x6c, 0x73, 0x2f, 0x73, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_sample_proto_rawDescOnce sync.Once
	file_sample_proto_rawDescData = file_sample_proto_rawDesc
)

func file_sample_proto_rawDescGZIP() []byte {
	file_sample_proto_rawDescOnce.Do(func() {
		file_sample_proto_rawDescData = protoimpl.X.CompressGZIP(file_sample_proto_rawDescData)
	})
	return file_sample_proto_rawDescData
}

var file_sample_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_sample_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_sample_proto_goTypes = []interface{}{
	(EnumValues)(0),          // 0: r3.EnumValues
	(Layer0_EnumEmbedded)(0), // 1: r3.Layer0.EnumEmbedded
	(*Supported)(nil),        // 2: r3.Supported
	(*Layer0)(nil),           // 3: r3.Layer0
	(*Layer1)(nil),           // 4: r3.Layer1
	(*BunchOTypes)(nil),      // 5: r3.BunchOTypes
}
var file_sample_proto_depIdxs = []int32{
	0, // 0: r3.Supported.ev:type_name -> r3.EnumValues
	4, // 1: r3.Layer0.layer1:type_name -> r3.Layer1
	1, // 2: r3.Layer0.ee:type_name -> r3.Layer0.EnumEmbedded
	2, // 3: r3.Layer1.supported:type_name -> r3.Supported
	0, // 4: r3.BunchOTypes.ev:type_name -> r3.EnumValues
	0, // 5: r3.BunchOTypes.l_ev:type_name -> r3.EnumValues
	2, // 6: r3.BunchOTypes.l_message:type_name -> r3.Supported
	7, // [7:7] is the sub-list for method output_type
	7, // [7:7] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_sample_proto_init() }
func file_sample_proto_init() {
	if File_sample_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_sample_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Supported); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_sample_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Layer0); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_sample_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Layer1); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_sample_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BunchOTypes); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_sample_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_sample_proto_goTypes,
		DependencyIndexes: file_sample_proto_depIdxs,
		EnumInfos:         file_sample_proto_enumTypes,
		MessageInfos:      file_sample_proto_msgTypes,
	}.Build()
	File_sample_proto = out.File
	file_sample_proto_rawDesc = nil
	file_sample_proto_goTypes = nil
	file_sample_proto_depIdxs = nil
}
