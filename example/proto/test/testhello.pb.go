// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.22.3
// source: testhello.proto

package test

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

type TestRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string    `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	One  []*Person `protobuf:"bytes,2,rep,name=one,proto3" json:"one,omitempty"`
}

func (x *TestRequest) Reset() {
	*x = TestRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_testhello_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TestRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TestRequest) ProtoMessage() {}

func (x *TestRequest) ProtoReflect() protoreflect.Message {
	mi := &file_testhello_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TestRequest.ProtoReflect.Descriptor instead.
func (*TestRequest) Descriptor() ([]byte, []int) {
	return file_testhello_proto_rawDescGZIP(), []int{0}
}

func (x *TestRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *TestRequest) GetOne() []*Person {
	if x != nil {
		return x.One
	}
	return nil
}

type TestReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *TestReply) Reset() {
	*x = TestReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_testhello_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TestReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TestReply) ProtoMessage() {}

func (x *TestReply) ProtoReflect() protoreflect.Message {
	mi := &file_testhello_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TestReply.ProtoReflect.Descriptor instead.
func (*TestReply) Descriptor() ([]byte, []int) {
	return file_testhello_proto_rawDescGZIP(), []int{1}
}

func (x *TestReply) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type Person struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name    string  `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Age     int32   `protobuf:"varint,2,opt,name=age,proto3" json:"age,omitempty"`
	Income  float32 `protobuf:"fixed32,3,opt,name=income,proto3" json:"income,omitempty"`
	Mystroy *Stroy  `protobuf:"bytes,4,opt,name=mystroy,proto3" json:"mystroy,omitempty"`
}

func (x *Person) Reset() {
	*x = Person{}
	if protoimpl.UnsafeEnabled {
		mi := &file_testhello_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Person) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Person) ProtoMessage() {}

func (x *Person) ProtoReflect() protoreflect.Message {
	mi := &file_testhello_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Person.ProtoReflect.Descriptor instead.
func (*Person) Descriptor() ([]byte, []int) {
	return file_testhello_proto_rawDescGZIP(), []int{2}
}

func (x *Person) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Person) GetAge() int32 {
	if x != nil {
		return x.Age
	}
	return 0
}

func (x *Person) GetIncome() float32 {
	if x != nil {
		return x.Income
	}
	return 0
}

func (x *Person) GetMystroy() *Stroy {
	if x != nil {
		return x.Mystroy
	}
	return nil
}

type Stroy struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Two string `protobuf:"bytes,1,opt,name=two,proto3" json:"two,omitempty"`
}

func (x *Stroy) Reset() {
	*x = Stroy{}
	if protoimpl.UnsafeEnabled {
		mi := &file_testhello_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Stroy) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Stroy) ProtoMessage() {}

func (x *Stroy) ProtoReflect() protoreflect.Message {
	mi := &file_testhello_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Stroy.ProtoReflect.Descriptor instead.
func (*Stroy) Descriptor() ([]byte, []int) {
	return file_testhello_proto_rawDescGZIP(), []int{3}
}

func (x *Stroy) GetTwo() string {
	if x != nil {
		return x.Two
	}
	return ""
}

var File_testhello_proto protoreflect.FileDescriptor

var file_testhello_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x74, 0x65, 0x73, 0x74, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x42, 0x0a, 0x0b, 0x54, 0x65, 0x73, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1f, 0x0a, 0x03, 0x6f,
	0x6e, 0x65, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x70, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x52, 0x03, 0x6f, 0x6e, 0x65, 0x22, 0x25, 0x0a, 0x09,
	0x54, 0x65, 0x73, 0x74, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x22, 0x6e, 0x0a, 0x06, 0x70, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x12, 0x12, 0x0a,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x12, 0x10, 0x0a, 0x03, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03,
	0x61, 0x67, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x69, 0x6e, 0x63, 0x6f, 0x6d, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x02, 0x52, 0x06, 0x69, 0x6e, 0x63, 0x6f, 0x6d, 0x65, 0x12, 0x26, 0x0a, 0x07, 0x6d,
	0x79, 0x73, 0x74, 0x72, 0x6f, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x73, 0x74, 0x72, 0x6f, 0x79, 0x52, 0x07, 0x6d, 0x79, 0x73, 0x74,
	0x72, 0x6f, 0x79, 0x22, 0x19, 0x0a, 0x05, 0x73, 0x74, 0x72, 0x6f, 0x79, 0x12, 0x10, 0x0a, 0x03,
	0x74, 0x77, 0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x74, 0x77, 0x6f, 0x32, 0x40,
	0x0a, 0x0a, 0x4e, 0x65, 0x77, 0x47, 0x72, 0x65, 0x65, 0x74, 0x65, 0x72, 0x12, 0x32, 0x0a, 0x08,
	0x53, 0x61, 0x79, 0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x12, 0x12, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x54, 0x65, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x10, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00,
	0x42, 0x07, 0x5a, 0x05, 0x74, 0x65, 0x73, 0x74, 0x2f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_testhello_proto_rawDescOnce sync.Once
	file_testhello_proto_rawDescData = file_testhello_proto_rawDesc
)

func file_testhello_proto_rawDescGZIP() []byte {
	file_testhello_proto_rawDescOnce.Do(func() {
		file_testhello_proto_rawDescData = protoimpl.X.CompressGZIP(file_testhello_proto_rawDescData)
	})
	return file_testhello_proto_rawDescData
}

var file_testhello_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_testhello_proto_goTypes = []interface{}{
	(*TestRequest)(nil), // 0: proto.TestRequest
	(*TestReply)(nil),   // 1: proto.TestReply
	(*Person)(nil),      // 2: proto.person
	(*Stroy)(nil),       // 3: proto.stroy
}
var file_testhello_proto_depIdxs = []int32{
	2, // 0: proto.TestRequest.one:type_name -> proto.person
	3, // 1: proto.person.mystroy:type_name -> proto.stroy
	0, // 2: proto.NewGreeter.SayHello:input_type -> proto.TestRequest
	1, // 3: proto.NewGreeter.SayHello:output_type -> proto.TestReply
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_testhello_proto_init() }
func file_testhello_proto_init() {
	if File_testhello_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_testhello_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TestRequest); i {
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
		file_testhello_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TestReply); i {
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
		file_testhello_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Person); i {
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
		file_testhello_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Stroy); i {
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
			RawDescriptor: file_testhello_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_testhello_proto_goTypes,
		DependencyIndexes: file_testhello_proto_depIdxs,
		MessageInfos:      file_testhello_proto_msgTypes,
	}.Build()
	File_testhello_proto = out.File
	file_testhello_proto_rawDesc = nil
	file_testhello_proto_goTypes = nil
	file_testhello_proto_depIdxs = nil
}
