// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.14.0
// source: packets/proto/raw.proto

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

type Raw struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type int32 `protobuf:"varint,1,opt,name=type,proto3" json:"type,omitempty"`
}

func (x *Raw) Reset() {
	*x = Raw{}
	if protoimpl.UnsafeEnabled {
		mi := &file_packets_proto_raw_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Raw) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Raw) ProtoMessage() {}

func (x *Raw) ProtoReflect() protoreflect.Message {
	mi := &file_packets_proto_raw_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Raw.ProtoReflect.Descriptor instead.
func (*Raw) Descriptor() ([]byte, []int) {
	return file_packets_proto_raw_proto_rawDescGZIP(), []int{0}
}

func (x *Raw) GetType() int32 {
	if x != nil {
		return x.Type
	}
	return 0
}

var File_packets_proto_raw_proto protoreflect.FileDescriptor

var file_packets_proto_raw_proto_rawDesc = []byte{
	0x0a, 0x17, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f,
	0x72, 0x61, 0x77, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x70, 0x61, 0x63, 0x6b, 0x65,
	0x74, 0x73, 0x22, 0x19, 0x0a, 0x03, 0x52, 0x61, 0x77, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x42, 0x0e, 0x5a,
	0x0c, 0x64, 0x69, 0x73, 0x74, 0x3b, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x73, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_packets_proto_raw_proto_rawDescOnce sync.Once
	file_packets_proto_raw_proto_rawDescData = file_packets_proto_raw_proto_rawDesc
)

func file_packets_proto_raw_proto_rawDescGZIP() []byte {
	file_packets_proto_raw_proto_rawDescOnce.Do(func() {
		file_packets_proto_raw_proto_rawDescData = protoimpl.X.CompressGZIP(file_packets_proto_raw_proto_rawDescData)
	})
	return file_packets_proto_raw_proto_rawDescData
}

var file_packets_proto_raw_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_packets_proto_raw_proto_goTypes = []interface{}{
	(*Raw)(nil), // 0: packets.Raw
}
var file_packets_proto_raw_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_packets_proto_raw_proto_init() }
func file_packets_proto_raw_proto_init() {
	if File_packets_proto_raw_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_packets_proto_raw_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Raw); i {
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
			RawDescriptor: file_packets_proto_raw_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_packets_proto_raw_proto_goTypes,
		DependencyIndexes: file_packets_proto_raw_proto_depIdxs,
		MessageInfos:      file_packets_proto_raw_proto_msgTypes,
	}.Build()
	File_packets_proto_raw_proto = out.File
	file_packets_proto_raw_proto_rawDesc = nil
	file_packets_proto_raw_proto_goTypes = nil
	file_packets_proto_raw_proto_depIdxs = nil
}
