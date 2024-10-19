// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.14.0
// source: packets/proto/connection.proto

package packets

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

type Join struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ptype int32 `protobuf:"varint,1,opt,name=ptype,proto3" json:"ptype,omitempty"`
	Wtf   bool  `protobuf:"varint,2,opt,name=wtf,proto3" json:"wtf,omitempty"`
}

func (x *Join) Reset() {
	*x = Join{}
	if protoimpl.UnsafeEnabled {
		mi := &file_packets_proto_connection_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Join) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Join) ProtoMessage() {}

func (x *Join) ProtoReflect() protoreflect.Message {
	mi := &file_packets_proto_connection_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Join.ProtoReflect.Descriptor instead.
func (*Join) Descriptor() ([]byte, []int) {
	return file_packets_proto_connection_proto_rawDescGZIP(), []int{0}
}

func (x *Join) GetPtype() int32 {
	if x != nil {
		return x.Ptype
	}
	return 0
}

func (x *Join) GetWtf() bool {
	if x != nil {
		return x.Wtf
	}
	return false
}

type PlayerJoined struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type int32 `protobuf:"varint,1,opt,name=type,proto3" json:"type,omitempty"`
	Id   int64 `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"`
	X    int64 `protobuf:"varint,3,opt,name=x,proto3" json:"x,omitempty"`
	Y    int64 `protobuf:"varint,4,opt,name=y,proto3" json:"y,omitempty"`
}

func (x *PlayerJoined) Reset() {
	*x = PlayerJoined{}
	if protoimpl.UnsafeEnabled {
		mi := &file_packets_proto_connection_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PlayerJoined) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PlayerJoined) ProtoMessage() {}

func (x *PlayerJoined) ProtoReflect() protoreflect.Message {
	mi := &file_packets_proto_connection_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PlayerJoined.ProtoReflect.Descriptor instead.
func (*PlayerJoined) Descriptor() ([]byte, []int) {
	return file_packets_proto_connection_proto_rawDescGZIP(), []int{1}
}

func (x *PlayerJoined) GetType() int32 {
	if x != nil {
		return x.Type
	}
	return 0
}

func (x *PlayerJoined) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *PlayerJoined) GetX() int64 {
	if x != nil {
		return x.X
	}
	return 0
}

func (x *PlayerJoined) GetY() int64 {
	if x != nil {
		return x.Y
	}
	return 0
}

type PlayerDisconnected struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type int32 `protobuf:"varint,1,opt,name=type,proto3" json:"type,omitempty"`
	Id   int64 `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *PlayerDisconnected) Reset() {
	*x = PlayerDisconnected{}
	if protoimpl.UnsafeEnabled {
		mi := &file_packets_proto_connection_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PlayerDisconnected) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PlayerDisconnected) ProtoMessage() {}

func (x *PlayerDisconnected) ProtoReflect() protoreflect.Message {
	mi := &file_packets_proto_connection_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PlayerDisconnected.ProtoReflect.Descriptor instead.
func (*PlayerDisconnected) Descriptor() ([]byte, []int) {
	return file_packets_proto_connection_proto_rawDescGZIP(), []int{2}
}

func (x *PlayerDisconnected) GetType() int32 {
	if x != nil {
		return x.Type
	}
	return 0
}

func (x *PlayerDisconnected) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

var File_packets_proto_connection_proto protoreflect.FileDescriptor

var file_packets_proto_connection_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f,
	0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x07, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x73, 0x22, 0x2e, 0x0a, 0x04, 0x4a, 0x6f, 0x69,
	0x6e, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x05, 0x70, 0x74, 0x79, 0x70, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x77, 0x74, 0x66, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x03, 0x77, 0x74, 0x66, 0x22, 0x4e, 0x0a, 0x0c, 0x50, 0x6c, 0x61,
	0x79, 0x65, 0x72, 0x4a, 0x6f, 0x69, 0x6e, 0x65, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x0c, 0x0a,
	0x01, 0x78, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x01, 0x78, 0x12, 0x0c, 0x0a, 0x01, 0x79,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x01, 0x79, 0x22, 0x38, 0x0a, 0x12, 0x50, 0x6c, 0x61,
	0x79, 0x65, 0x72, 0x44, 0x69, 0x73, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x65, 0x64, 0x12,
	0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x74,
	0x79, 0x70, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x02, 0x69, 0x64, 0x42, 0x0e, 0x5a, 0x0c, 0x64, 0x69, 0x73, 0x74, 0x3b, 0x70, 0x61, 0x63, 0x6b,
	0x65, 0x74, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_packets_proto_connection_proto_rawDescOnce sync.Once
	file_packets_proto_connection_proto_rawDescData = file_packets_proto_connection_proto_rawDesc
)

func file_packets_proto_connection_proto_rawDescGZIP() []byte {
	file_packets_proto_connection_proto_rawDescOnce.Do(func() {
		file_packets_proto_connection_proto_rawDescData = protoimpl.X.CompressGZIP(file_packets_proto_connection_proto_rawDescData)
	})
	return file_packets_proto_connection_proto_rawDescData
}

var file_packets_proto_connection_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_packets_proto_connection_proto_goTypes = []interface{}{
	(*Join)(nil),               // 0: packets.Join
	(*PlayerJoined)(nil),       // 1: packets.PlayerJoined
	(*PlayerDisconnected)(nil), // 2: packets.PlayerDisconnected
}
var file_packets_proto_connection_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_packets_proto_connection_proto_init() }
func file_packets_proto_connection_proto_init() {
	if File_packets_proto_connection_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_packets_proto_connection_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Join); i {
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
		file_packets_proto_connection_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PlayerJoined); i {
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
		file_packets_proto_connection_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PlayerDisconnected); i {
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
			RawDescriptor: file_packets_proto_connection_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_packets_proto_connection_proto_goTypes,
		DependencyIndexes: file_packets_proto_connection_proto_depIdxs,
		MessageInfos:      file_packets_proto_connection_proto_msgTypes,
	}.Build()
	File_packets_proto_connection_proto = out.File
	file_packets_proto_connection_proto_rawDesc = nil
	file_packets_proto_connection_proto_goTypes = nil
	file_packets_proto_connection_proto_depIdxs = nil
}