// Code generated by protoc-gen-go. DO NOT EDIT.
// source: clustermessage.proto

package clustermessage

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

type CommandType int32

const (
	CommandType_Reserved        CommandType = 0
	CommandType_ClusterRegist   CommandType = 1
	CommandType_ClusterUnregist CommandType = 2
	CommandType_NeighborRoute   CommandType = 3
	CommandType_SubTreeRoute    CommandType = 4
	CommandType_DeployReq       CommandType = 5
	CommandType_DeployResp      CommandType = 6
	CommandType_ControlReq      CommandType = 7
	CommandType_ControlResp     CommandType = 8
	CommandType_EdgeReport      CommandType = 9
)

var CommandType_name = map[int32]string{
	0: "Reserved",
	1: "ClusterRegist",
	2: "ClusterUnregist",
	3: "NeighborRoute",
	4: "SubTreeRoute",
	5: "DeployReq",
	6: "DeployResp",
	7: "ControlReq",
	8: "ControlResp",
	9: "EdgeReport",
}

var CommandType_value = map[string]int32{
	"Reserved":        0,
	"ClusterRegist":   1,
	"ClusterUnregist": 2,
	"NeighborRoute":   3,
	"SubTreeRoute":    4,
	"DeployReq":       5,
	"DeployResp":      6,
	"ControlReq":      7,
	"ControlResp":     8,
	"EdgeReport":      9,
}

func (x CommandType) String() string {
	return proto.EnumName(CommandType_name, int32(x))
}

func (CommandType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_cb5c8b0b58767cdb, []int{0}
}

// ClusterMessage is the message between cluster controllers and maybe cc and cluster shim.
type ClusterMessage struct {
	Head                 *MessageHead `protobuf:"bytes,1,opt,name=Head,proto3" json:"Head,omitempty"`
	Body                 []byte       `protobuf:"bytes,2,opt,name=Body,proto3" json:"Body,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *ClusterMessage) Reset()         { *m = ClusterMessage{} }
func (m *ClusterMessage) String() string { return proto.CompactTextString(m) }
func (*ClusterMessage) ProtoMessage()    {}
func (*ClusterMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_cb5c8b0b58767cdb, []int{0}
}

func (m *ClusterMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ClusterMessage.Unmarshal(m, b)
}
func (m *ClusterMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ClusterMessage.Marshal(b, m, deterministic)
}
func (m *ClusterMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ClusterMessage.Merge(m, src)
}
func (m *ClusterMessage) XXX_Size() int {
	return xxx_messageInfo_ClusterMessage.Size(m)
}
func (m *ClusterMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_ClusterMessage.DiscardUnknown(m)
}

var xxx_messageInfo_ClusterMessage proto.InternalMessageInfo

func (m *ClusterMessage) GetHead() *MessageHead {
	if m != nil {
		return m.Head
	}
	return nil
}

func (m *ClusterMessage) GetBody() []byte {
	if m != nil {
		return m.Body
	}
	return nil
}

type MessageHead struct {
	// MessageID is the uuid of a cluster message.
	// if the message comes from a crd, the messageid is the name of the crd.
	MessageID            string      `protobuf:"bytes,1,opt,name=MessageID,proto3" json:"MessageID,omitempty"`
	Command              CommandType `protobuf:"varint,2,opt,name=Command,proto3,enum=clustermessage.CommandType" json:"Command,omitempty"`
	ClusterSelector      string      `protobuf:"bytes,3,opt,name=ClusterSelector,proto3" json:"ClusterSelector,omitempty"`
	ClusterName          string      `protobuf:"bytes,4,opt,name=ClusterName,proto3" json:"ClusterName,omitempty"`
	ParentClusterName    string      `protobuf:"bytes,5,opt,name=ParentClusterName,proto3" json:"ParentClusterName,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *MessageHead) Reset()         { *m = MessageHead{} }
func (m *MessageHead) String() string { return proto.CompactTextString(m) }
func (*MessageHead) ProtoMessage()    {}
func (*MessageHead) Descriptor() ([]byte, []int) {
	return fileDescriptor_cb5c8b0b58767cdb, []int{1}
}

func (m *MessageHead) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MessageHead.Unmarshal(m, b)
}
func (m *MessageHead) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MessageHead.Marshal(b, m, deterministic)
}
func (m *MessageHead) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MessageHead.Merge(m, src)
}
func (m *MessageHead) XXX_Size() int {
	return xxx_messageInfo_MessageHead.Size(m)
}
func (m *MessageHead) XXX_DiscardUnknown() {
	xxx_messageInfo_MessageHead.DiscardUnknown(m)
}

var xxx_messageInfo_MessageHead proto.InternalMessageInfo

func (m *MessageHead) GetMessageID() string {
	if m != nil {
		return m.MessageID
	}
	return ""
}

func (m *MessageHead) GetCommand() CommandType {
	if m != nil {
		return m.Command
	}
	return CommandType_Reserved
}

func (m *MessageHead) GetClusterSelector() string {
	if m != nil {
		return m.ClusterSelector
	}
	return ""
}

func (m *MessageHead) GetClusterName() string {
	if m != nil {
		return m.ClusterName
	}
	return ""
}

func (m *MessageHead) GetParentClusterName() string {
	if m != nil {
		return m.ParentClusterName
	}
	return ""
}

type ControllerTask struct {
	Destination          string   `protobuf:"bytes,1,opt,name=Destination,proto3" json:"Destination,omitempty"`
	Method               string   `protobuf:"bytes,2,opt,name=Method,proto3" json:"Method,omitempty"`
	URI                  string   `protobuf:"bytes,3,opt,name=URI,proto3" json:"URI,omitempty"`
	Body                 []byte   `protobuf:"bytes,4,opt,name=Body,proto3" json:"Body,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ControllerTask) Reset()         { *m = ControllerTask{} }
func (m *ControllerTask) String() string { return proto.CompactTextString(m) }
func (*ControllerTask) ProtoMessage()    {}
func (*ControllerTask) Descriptor() ([]byte, []int) {
	return fileDescriptor_cb5c8b0b58767cdb, []int{2}
}

func (m *ControllerTask) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ControllerTask.Unmarshal(m, b)
}
func (m *ControllerTask) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ControllerTask.Marshal(b, m, deterministic)
}
func (m *ControllerTask) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ControllerTask.Merge(m, src)
}
func (m *ControllerTask) XXX_Size() int {
	return xxx_messageInfo_ControllerTask.Size(m)
}
func (m *ControllerTask) XXX_DiscardUnknown() {
	xxx_messageInfo_ControllerTask.DiscardUnknown(m)
}

var xxx_messageInfo_ControllerTask proto.InternalMessageInfo

func (m *ControllerTask) GetDestination() string {
	if m != nil {
		return m.Destination
	}
	return ""
}

func (m *ControllerTask) GetMethod() string {
	if m != nil {
		return m.Method
	}
	return ""
}

func (m *ControllerTask) GetURI() string {
	if m != nil {
		return m.URI
	}
	return ""
}

func (m *ControllerTask) GetBody() []byte {
	if m != nil {
		return m.Body
	}
	return nil
}

type ControllerTaskResponse struct {
	Timestamp            int64    `protobuf:"varint,1,opt,name=Timestamp,proto3" json:"Timestamp,omitempty"`
	StatusCode           int32    `protobuf:"varint,2,opt,name=StatusCode,proto3" json:"StatusCode,omitempty"`
	Body                 []byte   `protobuf:"bytes,3,opt,name=Body,proto3" json:"Body,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ControllerTaskResponse) Reset()         { *m = ControllerTaskResponse{} }
func (m *ControllerTaskResponse) String() string { return proto.CompactTextString(m) }
func (*ControllerTaskResponse) ProtoMessage()    {}
func (*ControllerTaskResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_cb5c8b0b58767cdb, []int{3}
}

func (m *ControllerTaskResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ControllerTaskResponse.Unmarshal(m, b)
}
func (m *ControllerTaskResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ControllerTaskResponse.Marshal(b, m, deterministic)
}
func (m *ControllerTaskResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ControllerTaskResponse.Merge(m, src)
}
func (m *ControllerTaskResponse) XXX_Size() int {
	return xxx_messageInfo_ControllerTaskResponse.Size(m)
}
func (m *ControllerTaskResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ControllerTaskResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ControllerTaskResponse proto.InternalMessageInfo

func (m *ControllerTaskResponse) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *ControllerTaskResponse) GetStatusCode() int32 {
	if m != nil {
		return m.StatusCode
	}
	return 0
}

func (m *ControllerTaskResponse) GetBody() []byte {
	if m != nil {
		return m.Body
	}
	return nil
}

type DeployTask struct {
	Replicas             int32             `protobuf:"varint,1,opt,name=Replicas,proto3" json:"Replicas,omitempty"`
	PodParams            map[string]string `protobuf:"bytes,2,rep,name=PodParams,proto3" json:"PodParams,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Status               string            `protobuf:"bytes,3,opt,name=Status,proto3" json:"Status,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *DeployTask) Reset()         { *m = DeployTask{} }
func (m *DeployTask) String() string { return proto.CompactTextString(m) }
func (*DeployTask) ProtoMessage()    {}
func (*DeployTask) Descriptor() ([]byte, []int) {
	return fileDescriptor_cb5c8b0b58767cdb, []int{4}
}

func (m *DeployTask) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeployTask.Unmarshal(m, b)
}
func (m *DeployTask) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeployTask.Marshal(b, m, deterministic)
}
func (m *DeployTask) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeployTask.Merge(m, src)
}
func (m *DeployTask) XXX_Size() int {
	return xxx_messageInfo_DeployTask.Size(m)
}
func (m *DeployTask) XXX_DiscardUnknown() {
	xxx_messageInfo_DeployTask.DiscardUnknown(m)
}

var xxx_messageInfo_DeployTask proto.InternalMessageInfo

func (m *DeployTask) GetReplicas() int32 {
	if m != nil {
		return m.Replicas
	}
	return 0
}

func (m *DeployTask) GetPodParams() map[string]string {
	if m != nil {
		return m.PodParams
	}
	return nil
}

func (m *DeployTask) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func init() {
	proto.RegisterEnum("clustermessage.CommandType", CommandType_name, CommandType_value)
	proto.RegisterType((*ClusterMessage)(nil), "clustermessage.ClusterMessage")
	proto.RegisterType((*MessageHead)(nil), "clustermessage.MessageHead")
	proto.RegisterType((*ControllerTask)(nil), "clustermessage.ControllerTask")
	proto.RegisterType((*ControllerTaskResponse)(nil), "clustermessage.ControllerTaskResponse")
	proto.RegisterType((*DeployTask)(nil), "clustermessage.DeployTask")
	proto.RegisterMapType((map[string]string)(nil), "clustermessage.DeployTask.PodParamsEntry")
}

func init() { proto.RegisterFile("clustermessage.proto", fileDescriptor_cb5c8b0b58767cdb) }

var fileDescriptor_cb5c8b0b58767cdb = []byte{
	// 507 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x53, 0xd1, 0x6e, 0xda, 0x30,
	0x14, 0x5d, 0x08, 0xd0, 0xe6, 0x42, 0xd3, 0xf4, 0xae, 0xaa, 0x50, 0x37, 0x4d, 0x88, 0x27, 0x36,
	0x4d, 0x4c, 0xea, 0x34, 0x69, 0x9a, 0xf6, 0x34, 0xa8, 0xb6, 0x3e, 0xb4, 0x42, 0x06, 0x3e, 0xc0,
	0x90, 0x2b, 0x9a, 0x35, 0x89, 0x33, 0xdb, 0x54, 0xe2, 0xcf, 0xf6, 0x07, 0xfb, 0x89, 0x7d, 0x4c,
	0x65, 0xc7, 0x0d, 0x81, 0xbe, 0xf9, 0x1c, 0x9f, 0xdc, 0x73, 0xef, 0xb9, 0x0e, 0x9c, 0xaf, 0xd2,
	0x8d, 0xd2, 0x24, 0x33, 0x52, 0x8a, 0xaf, 0x69, 0x54, 0x48, 0xa1, 0x05, 0x86, 0xfb, 0xec, 0x60,
	0x01, 0xe1, 0xb8, 0x64, 0x6e, 0x4b, 0x06, 0x3f, 0x41, 0xf3, 0x17, 0xf1, 0xb8, 0xe7, 0xf5, 0xbd,
	0x61, 0xe7, 0xea, 0xcd, 0xe8, 0xa0, 0x8c, 0x93, 0x19, 0x09, 0xb3, 0x42, 0x44, 0x68, 0xfe, 0x10,
	0xf1, 0xb6, 0xd7, 0xe8, 0x7b, 0xc3, 0x2e, 0xb3, 0xe7, 0xc1, 0x7f, 0x0f, 0x3a, 0x35, 0x25, 0xbe,
	0x85, 0xc0, 0xc1, 0x9b, 0x89, 0xad, 0x1c, 0xb0, 0x1d, 0x81, 0x5f, 0xe0, 0x68, 0x2c, 0xb2, 0x8c,
	0xe7, 0xb1, 0x2d, 0x12, 0xbe, 0x74, 0x75, 0xd7, 0xf3, 0x6d, 0x41, 0xec, 0x59, 0x8b, 0x43, 0x38,
	0x75, 0xbd, 0xcf, 0x28, 0xa5, 0x95, 0x16, 0xb2, 0xe7, 0xdb, 0xd2, 0x87, 0x34, 0xf6, 0xa1, 0xe3,
	0xa8, 0x3b, 0x9e, 0x51, 0xaf, 0x69, 0x55, 0x75, 0x0a, 0x3f, 0xc2, 0xd9, 0x94, 0x4b, 0xca, 0x75,
	0x5d, 0xd7, 0xb2, 0xba, 0x97, 0x17, 0x83, 0x02, 0xc2, 0xb1, 0xc8, 0xb5, 0x14, 0x69, 0x4a, 0x72,
	0xce, 0xd5, 0x83, 0x71, 0x98, 0x90, 0xd2, 0x49, 0xce, 0x75, 0x22, 0x72, 0x37, 0x62, 0x9d, 0xc2,
	0x0b, 0x68, 0xdf, 0x92, 0xbe, 0x17, 0xe5, 0x8c, 0x01, 0x73, 0x08, 0x23, 0xf0, 0x17, 0xec, 0xc6,
	0x75, 0x6e, 0x8e, 0x55, 0xa0, 0xcd, 0x5a, 0xa0, 0xbf, 0xe1, 0x62, 0xdf, 0x91, 0x91, 0x2a, 0x44,
	0xae, 0xc8, 0x44, 0x3b, 0x4f, 0x32, 0x52, 0x9a, 0x67, 0x85, 0xf5, 0xf5, 0xd9, 0x8e, 0xc0, 0x77,
	0x00, 0x33, 0xcd, 0xf5, 0x46, 0x8d, 0x45, 0x4c, 0xd6, 0xb9, 0xc5, 0x6a, 0x4c, 0xe5, 0xe5, 0xd7,
	0xbc, 0xfe, 0x79, 0x00, 0x13, 0x2a, 0x52, 0xb1, 0xb5, 0xa3, 0x5d, 0xc2, 0x31, 0xa3, 0x22, 0x4d,
	0x56, 0x5c, 0xd9, 0xfa, 0x2d, 0x56, 0x61, 0xfc, 0x09, 0xc1, 0x54, 0xc4, 0x53, 0x2e, 0x79, 0xa6,
	0x7a, 0x8d, 0xbe, 0x3f, 0xec, 0x5c, 0xbd, 0x3f, 0xdc, 0xdd, 0xae, 0xd4, 0xa8, 0xd2, 0x5e, 0xe7,
	0x5a, 0x6e, 0xd9, 0xee, 0x5b, 0x93, 0x4e, 0xd9, 0x95, 0x0b, 0xc2, 0xa1, 0xcb, 0xef, 0x10, 0xee,
	0x7f, 0x64, 0xf2, 0x7a, 0xa0, 0xad, 0x4b, 0xd8, 0x1c, 0xf1, 0x1c, 0x5a, 0x8f, 0x3c, 0xdd, 0x90,
	0x0b, 0xb6, 0x04, 0xdf, 0x1a, 0x5f, 0xbd, 0x0f, 0x7f, 0x3d, 0xe8, 0xd4, 0x9e, 0x0e, 0x76, 0xcd,
	0x28, 0x8a, 0xe4, 0x23, 0xc5, 0xd1, 0x2b, 0x3c, 0x83, 0x13, 0xb7, 0x54, 0x46, 0xeb, 0x44, 0xe9,
	0xc8, 0xc3, 0xd7, 0xd5, 0x93, 0x5a, 0xe4, 0xb2, 0x24, 0x1b, 0x46, 0x77, 0x47, 0xc9, 0xfa, 0x7e,
	0x29, 0x24, 0x13, 0x1b, 0x4d, 0x91, 0x8f, 0x11, 0x74, 0x67, 0x9b, 0xe5, 0x5c, 0x12, 0x95, 0x4c,
	0x13, 0x4f, 0x20, 0x28, 0x07, 0x65, 0xf4, 0x27, 0x6a, 0x61, 0xf8, 0x1c, 0xa1, 0xd9, 0x53, 0xd4,
	0x36, 0xd8, 0xed, 0xcf, 0xdc, 0x1f, 0xe1, 0xa9, 0x69, 0xcc, 0x61, 0x55, 0x44, 0xc7, 0x46, 0x70,
	0x1d, 0xaf, 0x89, 0x51, 0x21, 0xa4, 0x8e, 0x82, 0x65, 0xdb, 0xfe, 0xaf, 0x9f, 0x9f, 0x02, 0x00,
	0x00, 0xff, 0xff, 0x40, 0xc0, 0x9c, 0x7f, 0xc7, 0x03, 0x00, 0x00,
}
