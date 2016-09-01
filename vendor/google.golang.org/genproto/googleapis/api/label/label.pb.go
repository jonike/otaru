// Code generated by protoc-gen-go.
// source: google.golang.org/genproto/googleapis/api/label/label.proto
// DO NOT EDIT!

/*
Package google_api is a generated protocol buffer package.

It is generated from these files:
	google.golang.org/genproto/googleapis/api/label/label.proto

It has these top-level messages:
	LabelDescriptor
*/
package google_api

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Value types that can be used as label values.
type LabelDescriptor_ValueType int32

const (
	// A variable-length string. This is the default.
	LabelDescriptor_STRING LabelDescriptor_ValueType = 0
	// Boolean; true or false.
	LabelDescriptor_BOOL LabelDescriptor_ValueType = 1
	// A 64-bit signed integer.
	LabelDescriptor_INT64 LabelDescriptor_ValueType = 2
)

var LabelDescriptor_ValueType_name = map[int32]string{
	0: "STRING",
	1: "BOOL",
	2: "INT64",
}
var LabelDescriptor_ValueType_value = map[string]int32{
	"STRING": 0,
	"BOOL":   1,
	"INT64":  2,
}

func (x LabelDescriptor_ValueType) String() string {
	return proto.EnumName(LabelDescriptor_ValueType_name, int32(x))
}
func (LabelDescriptor_ValueType) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0} }

// A description of a label.
type LabelDescriptor struct {
	// The label key.
	Key string `protobuf:"bytes,1,opt,name=key" json:"key,omitempty"`
	// The type of data that can be assigned to the label.
	ValueType LabelDescriptor_ValueType `protobuf:"varint,2,opt,name=value_type,json=valueType,enum=google.api.LabelDescriptor_ValueType" json:"value_type,omitempty"`
	// A human-readable description for the label.
	Description string `protobuf:"bytes,3,opt,name=description" json:"description,omitempty"`
}

func (m *LabelDescriptor) Reset()                    { *m = LabelDescriptor{} }
func (m *LabelDescriptor) String() string            { return proto.CompactTextString(m) }
func (*LabelDescriptor) ProtoMessage()               {}
func (*LabelDescriptor) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func init() {
	proto.RegisterType((*LabelDescriptor)(nil), "google.api.LabelDescriptor")
	proto.RegisterEnum("google.api.LabelDescriptor_ValueType", LabelDescriptor_ValueType_name, LabelDescriptor_ValueType_value)
}

func init() {
	proto.RegisterFile("google.golang.org/genproto/googleapis/api/label/label.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 234 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x5c, 0x4f, 0x4d, 0x4b, 0x03, 0x31,
	0x10, 0x35, 0xad, 0x16, 0x77, 0x84, 0x1a, 0x72, 0xda, 0x93, 0x2c, 0x05, 0xa1, 0x07, 0x49, 0x40,
	0xc5, 0x8b, 0xb7, 0xa5, 0x20, 0x85, 0xd2, 0x96, 0x75, 0xf1, 0x2a, 0x69, 0x1d, 0x42, 0x30, 0xee,
	0x84, 0x74, 0x2d, 0xec, 0x4f, 0xf3, 0xdf, 0x49, 0xd2, 0xaa, 0x8b, 0x97, 0xe1, 0x31, 0xef, 0x83,
	0xf7, 0xe0, 0xd1, 0x10, 0x19, 0x87, 0xd2, 0x90, 0xd3, 0x8d, 0x91, 0x14, 0x8c, 0x32, 0xd8, 0xf8,
	0x40, 0x2d, 0xa9, 0x03, 0xa5, 0xbd, 0xdd, 0x29, 0xed, 0xad, 0x72, 0x7a, 0x83, 0xee, 0x70, 0x65,
	0x12, 0x08, 0x38, 0x9a, 0xb5, 0xb7, 0x93, 0x2f, 0x06, 0x97, 0x8b, 0xc8, 0xcd, 0x70, 0xb7, 0x0d,
	0xd6, 0xb7, 0x14, 0x04, 0x87, 0xe1, 0x3b, 0x76, 0x39, 0x2b, 0xd8, 0x34, 0xab, 0x22, 0x14, 0x33,
	0x80, 0xbd, 0x76, 0x9f, 0xf8, 0xda, 0x76, 0x1e, 0xf3, 0x41, 0xc1, 0xa6, 0xe3, 0xdb, 0x6b, 0xf9,
	0x17, 0x23, 0xff, 0x45, 0xc8, 0x97, 0xa8, 0xae, 0x3b, 0x8f, 0x55, 0xb6, 0xff, 0x81, 0xa2, 0x80,
	0x8b, 0xb7, 0xa3, 0xc4, 0x52, 0x93, 0x0f, 0x53, 0x7e, 0xff, 0x35, 0xb9, 0x81, 0xec, 0xd7, 0x29,
	0x00, 0x46, 0xcf, 0x75, 0x35, 0x5f, 0x3e, 0xf1, 0x13, 0x71, 0x0e, 0xa7, 0xe5, 0x6a, 0xb5, 0xe0,
	0x4c, 0x64, 0x70, 0x36, 0x5f, 0xd6, 0x0f, 0xf7, 0x7c, 0x50, 0x5e, 0xc1, 0x78, 0x4b, 0x1f, 0xbd,
	0x1a, 0x25, 0xa4, 0x1e, 0xeb, 0xb8, 0x72, 0xcd, 0x36, 0xa3, 0x34, 0xf7, 0xee, 0x3b, 0x00, 0x00,
	0xff, 0xff, 0x3b, 0xdf, 0x1f, 0x56, 0x2d, 0x01, 0x00, 0x00,
}
