// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.6.1
// source: file.proto

package connection

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

type MessageType int32

const (
	MessageType_DUMMY           MessageType = 0
	MessageType_PUT             MessageType = 1
	MessageType_REGISTRATION    MessageType = 2
	MessageType_GET             MessageType = 3
	MessageType_RM              MessageType = 4
	MessageType_ACK             MessageType = 5
	MessageType_HEARTBEAT       MessageType = 6
	MessageType_ACK_LS          MessageType = 7
	MessageType_ERROR           MessageType = 8
	MessageType_ACK_PUT         MessageType = 9
	MessageType_LS              MessageType = 10
	MessageType_ACK_GET         MessageType = 11
	MessageType_CHECKSUM        MessageType = 12
	MessageType_HEARTBEAT_CHUNK MessageType = 13
)

// Enum value maps for MessageType.
var (
	MessageType_name = map[int32]string{
		0:  "DUMMY",
		1:  "PUT",
		2:  "REGISTRATION",
		3:  "GET",
		4:  "RM",
		5:  "ACK",
		6:  "HEARTBEAT",
		7:  "ACK_LS",
		8:  "ERROR",
		9:  "ACK_PUT",
		10: "LS",
		11: "ACK_GET",
		12: "CHECKSUM",
		13: "HEARTBEAT_CHUNK",
	}
	MessageType_value = map[string]int32{
		"DUMMY":           0,
		"PUT":             1,
		"REGISTRATION":    2,
		"GET":             3,
		"RM":              4,
		"ACK":             5,
		"HEARTBEAT":       6,
		"ACK_LS":          7,
		"ERROR":           8,
		"ACK_PUT":         9,
		"LS":              10,
		"ACK_GET":         11,
		"CHECKSUM":        12,
		"HEARTBEAT_CHUNK": 13,
	}
)

func (x MessageType) Enum() *MessageType {
	p := new(MessageType)
	*p = x
	return p
}

func (x MessageType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MessageType) Descriptor() protoreflect.EnumDescriptor {
	return file_file_proto_enumTypes[0].Descriptor()
}

func (MessageType) Type() protoreflect.EnumType {
	return &file_file_proto_enumTypes[0]
}

func (x MessageType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MessageType.Descriptor instead.
func (MessageType) EnumDescriptor() ([]byte, []int) {
	return file_file_proto_rawDescGZIP(), []int{0}
}

type FileData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	MessageType MessageType `protobuf:"varint,1,opt,name=messageType,proto3,enum=MessageType" json:"messageType,omitempty"`
	Path        string      `protobuf:"bytes,2,opt,name=path,proto3" json:"path,omitempty"`
	Size        []byte      `protobuf:"bytes,3,opt,name=size,proto3" json:"size,omitempty"`
	Data        string      `protobuf:"bytes,4,opt,name=data,proto3" json:"data,omitempty"`
	Checksum    []byte      `protobuf:"bytes,5,opt,name=checksum,proto3" json:"checksum,omitempty"`
	SenderId    string      `protobuf:"bytes,6,opt,name=senderId,proto3" json:"senderId,omitempty"`
	Chunk       []*Chunk    `protobuf:"bytes,7,rep,name=chunk,proto3" json:"chunk,omitempty"`
	Node        *Node       `protobuf:"bytes,8,opt,name=node,proto3" json:"node,omitempty"`
	DataSize    int64       `protobuf:"varint,9,opt,name=dataSize,proto3" json:"dataSize,omitempty"`
	Nodes       []*Node     `protobuf:"bytes,10,rep,name=nodes,proto3" json:"nodes,omitempty"`
}

func (x *FileData) Reset() {
	*x = FileData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_file_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FileData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FileData) ProtoMessage() {}

func (x *FileData) ProtoReflect() protoreflect.Message {
	mi := &file_file_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FileData.ProtoReflect.Descriptor instead.
func (*FileData) Descriptor() ([]byte, []int) {
	return file_file_proto_rawDescGZIP(), []int{0}
}

func (x *FileData) GetMessageType() MessageType {
	if x != nil {
		return x.MessageType
	}
	return MessageType_DUMMY
}

func (x *FileData) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

func (x *FileData) GetSize() []byte {
	if x != nil {
		return x.Size
	}
	return nil
}

func (x *FileData) GetData() string {
	if x != nil {
		return x.Data
	}
	return ""
}

func (x *FileData) GetChecksum() []byte {
	if x != nil {
		return x.Checksum
	}
	return nil
}

func (x *FileData) GetSenderId() string {
	if x != nil {
		return x.SenderId
	}
	return ""
}

func (x *FileData) GetChunk() []*Chunk {
	if x != nil {
		return x.Chunk
	}
	return nil
}

func (x *FileData) GetNode() *Node {
	if x != nil {
		return x.Node
	}
	return nil
}

func (x *FileData) GetDataSize() int64 {
	if x != nil {
		return x.DataSize
	}
	return 0
}

func (x *FileData) GetNodes() []*Node {
	if x != nil {
		return x.Nodes
	}
	return nil
}

type Node struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id       string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Hostname string `protobuf:"bytes,2,opt,name=hostname,proto3" json:"hostname,omitempty"`
	Port     int32  `protobuf:"varint,3,opt,name=port,proto3" json:"port,omitempty"`
}

func (x *Node) Reset() {
	*x = Node{}
	if protoimpl.UnsafeEnabled {
		mi := &file_file_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Node) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Node) ProtoMessage() {}

func (x *Node) ProtoReflect() protoreflect.Message {
	mi := &file_file_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Node.ProtoReflect.Descriptor instead.
func (*Node) Descriptor() ([]byte, []int) {
	return file_file_proto_rawDescGZIP(), []int{1}
}

func (x *Node) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Node) GetHostname() string {
	if x != nil {
		return x.Hostname
	}
	return ""
}

func (x *Node) GetPort() int32 {
	if x != nil {
		return x.Port
	}
	return 0
}

type Chunk struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name     string  `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Size     int64   `protobuf:"varint,2,opt,name=size,proto3" json:"size,omitempty"`
	Checksum int32   `protobuf:"varint,3,opt,name=checksum,proto3" json:"checksum,omitempty"`
	Nodes    []*Node `protobuf:"bytes,4,rep,name=nodes,proto3" json:"nodes,omitempty"`
	Num      int32   `protobuf:"varint,5,opt,name=num,proto3" json:"num,omitempty"`
}

func (x *Chunk) Reset() {
	*x = Chunk{}
	if protoimpl.UnsafeEnabled {
		mi := &file_file_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Chunk) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Chunk) ProtoMessage() {}

func (x *Chunk) ProtoReflect() protoreflect.Message {
	mi := &file_file_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Chunk.ProtoReflect.Descriptor instead.
func (*Chunk) Descriptor() ([]byte, []int) {
	return file_file_proto_rawDescGZIP(), []int{2}
}

func (x *Chunk) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Chunk) GetSize() int64 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *Chunk) GetChecksum() int32 {
	if x != nil {
		return x.Checksum
	}
	return 0
}

func (x *Chunk) GetNodes() []*Node {
	if x != nil {
		return x.Nodes
	}
	return nil
}

func (x *Chunk) GetNum() int32 {
	if x != nil {
		return x.Num
	}
	return 0
}

var File_file_proto protoreflect.FileDescriptor

var file_file_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x66, 0x69, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xa0, 0x02, 0x0a,
	0x08, 0x46, 0x69, 0x6c, 0x65, 0x44, 0x61, 0x74, 0x61, 0x12, 0x2e, 0x0a, 0x0b, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0c,
	0x2e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65, 0x52, 0x0b, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x74,
	0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x70, 0x61, 0x74, 0x68, 0x12, 0x12, 0x0a,
	0x04, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x73, 0x69, 0x7a,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x73, 0x75,
	0x6d, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x63, 0x68, 0x65, 0x63, 0x6b, 0x73, 0x75,
	0x6d, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x49, 0x64, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1c, 0x0a,
	0x05, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x18, 0x07, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x06, 0x2e, 0x43,
	0x68, 0x75, 0x6e, 0x6b, 0x52, 0x05, 0x63, 0x68, 0x75, 0x6e, 0x6b, 0x12, 0x19, 0x0a, 0x04, 0x6e,
	0x6f, 0x64, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x05, 0x2e, 0x4e, 0x6f, 0x64, 0x65,
	0x52, 0x04, 0x6e, 0x6f, 0x64, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x64, 0x61, 0x74, 0x61, 0x53, 0x69,
	0x7a, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x64, 0x61, 0x74, 0x61, 0x53, 0x69,
	0x7a, 0x65, 0x12, 0x1b, 0x0a, 0x05, 0x6e, 0x6f, 0x64, 0x65, 0x73, 0x18, 0x0a, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x05, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x52, 0x05, 0x6e, 0x6f, 0x64, 0x65, 0x73, 0x22,
	0x46, 0x0a, 0x04, 0x4e, 0x6f, 0x64, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x68, 0x6f, 0x73, 0x74, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x68, 0x6f, 0x73, 0x74, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x22, 0x7a, 0x0a, 0x05, 0x43, 0x68, 0x75, 0x6e, 0x6b,
	0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x68, 0x65, 0x63,
	0x6b, 0x73, 0x75, 0x6d, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x63, 0x68, 0x65, 0x63,
	0x6b, 0x73, 0x75, 0x6d, 0x12, 0x1b, 0x0a, 0x05, 0x6e, 0x6f, 0x64, 0x65, 0x73, 0x18, 0x04, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x05, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x52, 0x05, 0x6e, 0x6f, 0x64, 0x65,
	0x73, 0x12, 0x10, 0x0a, 0x03, 0x6e, 0x75, 0x6d, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03,
	0x6e, 0x75, 0x6d, 0x2a, 0xb8, 0x01, 0x0a, 0x0b, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x54,
	0x79, 0x70, 0x65, 0x12, 0x09, 0x0a, 0x05, 0x44, 0x55, 0x4d, 0x4d, 0x59, 0x10, 0x00, 0x12, 0x07,
	0x0a, 0x03, 0x50, 0x55, 0x54, 0x10, 0x01, 0x12, 0x10, 0x0a, 0x0c, 0x52, 0x45, 0x47, 0x49, 0x53,
	0x54, 0x52, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x10, 0x02, 0x12, 0x07, 0x0a, 0x03, 0x47, 0x45, 0x54,
	0x10, 0x03, 0x12, 0x06, 0x0a, 0x02, 0x52, 0x4d, 0x10, 0x04, 0x12, 0x07, 0x0a, 0x03, 0x41, 0x43,
	0x4b, 0x10, 0x05, 0x12, 0x0d, 0x0a, 0x09, 0x48, 0x45, 0x41, 0x52, 0x54, 0x42, 0x45, 0x41, 0x54,
	0x10, 0x06, 0x12, 0x0a, 0x0a, 0x06, 0x41, 0x43, 0x4b, 0x5f, 0x4c, 0x53, 0x10, 0x07, 0x12, 0x09,
	0x0a, 0x05, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x10, 0x08, 0x12, 0x0b, 0x0a, 0x07, 0x41, 0x43, 0x4b,
	0x5f, 0x50, 0x55, 0x54, 0x10, 0x09, 0x12, 0x06, 0x0a, 0x02, 0x4c, 0x53, 0x10, 0x0a, 0x12, 0x0b,
	0x0a, 0x07, 0x41, 0x43, 0x4b, 0x5f, 0x47, 0x45, 0x54, 0x10, 0x0b, 0x12, 0x0c, 0x0a, 0x08, 0x43,
	0x48, 0x45, 0x43, 0x4b, 0x53, 0x55, 0x4d, 0x10, 0x0c, 0x12, 0x13, 0x0a, 0x0f, 0x48, 0x45, 0x41,
	0x52, 0x54, 0x42, 0x45, 0x41, 0x54, 0x5f, 0x43, 0x48, 0x55, 0x4e, 0x4b, 0x10, 0x0d, 0x42, 0x0e,
	0x5a, 0x0c, 0x2e, 0x2f, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_file_proto_rawDescOnce sync.Once
	file_file_proto_rawDescData = file_file_proto_rawDesc
)

func file_file_proto_rawDescGZIP() []byte {
	file_file_proto_rawDescOnce.Do(func() {
		file_file_proto_rawDescData = protoimpl.X.CompressGZIP(file_file_proto_rawDescData)
	})
	return file_file_proto_rawDescData
}

var file_file_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_file_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_file_proto_goTypes = []interface{}{
	(MessageType)(0), // 0: MessageType
	(*FileData)(nil), // 1: FileData
	(*Node)(nil),     // 2: Node
	(*Chunk)(nil),    // 3: Chunk
}
var file_file_proto_depIdxs = []int32{
	0, // 0: FileData.messageType:type_name -> MessageType
	3, // 1: FileData.chunk:type_name -> Chunk
	2, // 2: FileData.node:type_name -> Node
	2, // 3: FileData.nodes:type_name -> Node
	2, // 4: Chunk.nodes:type_name -> Node
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_file_proto_init() }
func file_file_proto_init() {
	if File_file_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_file_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FileData); i {
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
		file_file_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Node); i {
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
		file_file_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Chunk); i {
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
			RawDescriptor: file_file_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_file_proto_goTypes,
		DependencyIndexes: file_file_proto_depIdxs,
		EnumInfos:         file_file_proto_enumTypes,
		MessageInfos:      file_file_proto_msgTypes,
	}.Build()
	File_file_proto = out.File
	file_file_proto_rawDesc = nil
	file_file_proto_goTypes = nil
	file_file_proto_depIdxs = nil
}
