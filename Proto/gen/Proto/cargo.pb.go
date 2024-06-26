// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        v4.25.3
// source: cargo.proto

package golang

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

type Cargo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name   string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Amount int64  `protobuf:"varint,2,opt,name=amount,proto3" json:"amount,omitempty"`
}

func (x *Cargo) Reset() {
	*x = Cargo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cargo_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Cargo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Cargo) ProtoMessage() {}

func (x *Cargo) ProtoReflect() protoreflect.Message {
	mi := &file_cargo_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Cargo.ProtoReflect.Descriptor instead.
func (*Cargo) Descriptor() ([]byte, []int) {
	return file_cargo_proto_rawDescGZIP(), []int{0}
}

func (x *Cargo) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Cargo) GetAmount() int64 {
	if x != nil {
		return x.Amount
	}
	return 0
}

type RequestSendCargo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cargo *Cargo `protobuf:"bytes,1,opt,name=cargo,proto3" json:"cargo,omitempty"`
}

func (x *RequestSendCargo) Reset() {
	*x = RequestSendCargo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cargo_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestSendCargo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestSendCargo) ProtoMessage() {}

func (x *RequestSendCargo) ProtoReflect() protoreflect.Message {
	mi := &file_cargo_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestSendCargo.ProtoReflect.Descriptor instead.
func (*RequestSendCargo) Descriptor() ([]byte, []int) {
	return file_cargo_proto_rawDescGZIP(), []int{1}
}

func (x *RequestSendCargo) GetCargo() *Cargo {
	if x != nil {
		return x.Cargo
	}
	return nil
}

type ResponseSendCargo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id int32 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *ResponseSendCargo) Reset() {
	*x = ResponseSendCargo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cargo_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResponseSendCargo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResponseSendCargo) ProtoMessage() {}

func (x *ResponseSendCargo) ProtoReflect() protoreflect.Message {
	mi := &file_cargo_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResponseSendCargo.ProtoReflect.Descriptor instead.
func (*ResponseSendCargo) Descriptor() ([]byte, []int) {
	return file_cargo_proto_rawDescGZIP(), []int{2}
}

func (x *ResponseSendCargo) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

type RequestReceivedCargo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id int32 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *RequestReceivedCargo) Reset() {
	*x = RequestReceivedCargo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cargo_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestReceivedCargo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestReceivedCargo) ProtoMessage() {}

func (x *RequestReceivedCargo) ProtoReflect() protoreflect.Message {
	mi := &file_cargo_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestReceivedCargo.ProtoReflect.Descriptor instead.
func (*RequestReceivedCargo) Descriptor() ([]byte, []int) {
	return file_cargo_proto_rawDescGZIP(), []int{3}
}

func (x *RequestReceivedCargo) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

type ResponseReceivedCargo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *ResponseReceivedCargo) Reset() {
	*x = ResponseReceivedCargo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_cargo_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResponseReceivedCargo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResponseReceivedCargo) ProtoMessage() {}

func (x *ResponseReceivedCargo) ProtoReflect() protoreflect.Message {
	mi := &file_cargo_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResponseReceivedCargo.ProtoReflect.Descriptor instead.
func (*ResponseReceivedCargo) Descriptor() ([]byte, []int) {
	return file_cargo_proto_rawDescGZIP(), []int{4}
}

func (x *ResponseReceivedCargo) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

var File_cargo_proto protoreflect.FileDescriptor

var file_cargo_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x63, 0x61, 0x72, 0x67, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x10, 0x62,
	0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2e, 0x63, 0x61, 0x72, 0x67, 0x6f, 0x2e, 0x76, 0x31, 0x22,
	0x33, 0x0a, 0x05, 0x43, 0x61, 0x72, 0x67, 0x6f, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a, 0x06,
	0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x61, 0x6d,
	0x6f, 0x75, 0x6e, 0x74, 0x22, 0x41, 0x0a, 0x10, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x53,
	0x65, 0x6e, 0x64, 0x43, 0x61, 0x72, 0x67, 0x6f, 0x12, 0x2d, 0x0a, 0x05, 0x63, 0x61, 0x72, 0x67,
	0x6f, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e,
	0x64, 0x2e, 0x63, 0x61, 0x72, 0x67, 0x6f, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x61, 0x72, 0x67, 0x6f,
	0x52, 0x05, 0x63, 0x61, 0x72, 0x67, 0x6f, 0x22, 0x23, 0x0a, 0x11, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x53, 0x65, 0x6e, 0x64, 0x43, 0x61, 0x72, 0x67, 0x6f, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x69, 0x64, 0x22, 0x26, 0x0a, 0x14,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x64, 0x43,
	0x61, 0x72, 0x67, 0x6f, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x02, 0x69, 0x64, 0x22, 0x2b, 0x0a, 0x15, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x64, 0x43, 0x61, 0x72, 0x67, 0x6f, 0x12, 0x12, 0x0a,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x32, 0xca, 0x01, 0x0a, 0x0c, 0x43, 0x61, 0x72, 0x67, 0x6f, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x12, 0x56, 0x0a, 0x09, 0x53, 0x65, 0x6e, 0x64, 0x43, 0x61, 0x72, 0x67, 0x6f, 0x12,
	0x22, 0x2e, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2e, 0x63, 0x61, 0x72, 0x67, 0x6f, 0x2e,
	0x76, 0x31, 0x2e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x53, 0x65, 0x6e, 0x64, 0x43, 0x61,
	0x72, 0x67, 0x6f, 0x1a, 0x23, 0x2e, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2e, 0x63, 0x61,
	0x72, 0x67, 0x6f, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x53,
	0x65, 0x6e, 0x64, 0x43, 0x61, 0x72, 0x67, 0x6f, 0x22, 0x00, 0x12, 0x62, 0x0a, 0x0d, 0x52, 0x65,
	0x63, 0x65, 0x69, 0x76, 0x65, 0x64, 0x43, 0x61, 0x72, 0x67, 0x6f, 0x12, 0x26, 0x2e, 0x62, 0x61,
	0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2e, 0x63, 0x61, 0x72, 0x67, 0x6f, 0x2e, 0x76, 0x31, 0x2e, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x52, 0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x64, 0x43, 0x61,
	0x72, 0x67, 0x6f, 0x1a, 0x27, 0x2e, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2e, 0x63, 0x61,
	0x72, 0x67, 0x6f, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x52,
	0x65, 0x63, 0x65, 0x69, 0x76, 0x65, 0x64, 0x43, 0x61, 0x72, 0x67, 0x6f, 0x22, 0x00, 0x42, 0x39,
	0x5a, 0x37, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x66, 0x61, 0x72,
	0x7a, 0x61, 0x64, 0x38, 0x30, 0x72, 0x61, 0x64, 0x2f, 0x68, 0x65, 0x69, 0x6d, 0x64, 0x61, 0x6c,
	0x6c, 0x2f, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x67, 0x6f, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x2f, 0x67, 0x6f, 0x6c, 0x61, 0x6e, 0x67, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_cargo_proto_rawDescOnce sync.Once
	file_cargo_proto_rawDescData = file_cargo_proto_rawDesc
)

func file_cargo_proto_rawDescGZIP() []byte {
	file_cargo_proto_rawDescOnce.Do(func() {
		file_cargo_proto_rawDescData = protoimpl.X.CompressGZIP(file_cargo_proto_rawDescData)
	})
	return file_cargo_proto_rawDescData
}

var file_cargo_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_cargo_proto_goTypes = []interface{}{
	(*Cargo)(nil),                 // 0: backend.cargo.v1.Cargo
	(*RequestSendCargo)(nil),      // 1: backend.cargo.v1.RequestSendCargo
	(*ResponseSendCargo)(nil),     // 2: backend.cargo.v1.ResponseSendCargo
	(*RequestReceivedCargo)(nil),  // 3: backend.cargo.v1.RequestReceivedCargo
	(*ResponseReceivedCargo)(nil), // 4: backend.cargo.v1.ResponseReceivedCargo
}
var file_cargo_proto_depIdxs = []int32{
	0, // 0: backend.cargo.v1.RequestSendCargo.cargo:type_name -> backend.cargo.v1.Cargo
	1, // 1: backend.cargo.v1.CargoService.SendCargo:input_type -> backend.cargo.v1.RequestSendCargo
	3, // 2: backend.cargo.v1.CargoService.ReceivedCargo:input_type -> backend.cargo.v1.RequestReceivedCargo
	2, // 3: backend.cargo.v1.CargoService.SendCargo:output_type -> backend.cargo.v1.ResponseSendCargo
	4, // 4: backend.cargo.v1.CargoService.ReceivedCargo:output_type -> backend.cargo.v1.ResponseReceivedCargo
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_cargo_proto_init() }
func file_cargo_proto_init() {
	if File_cargo_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_cargo_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Cargo); i {
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
		file_cargo_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequestSendCargo); i {
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
		file_cargo_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResponseSendCargo); i {
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
		file_cargo_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequestReceivedCargo); i {
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
		file_cargo_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResponseReceivedCargo); i {
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
			RawDescriptor: file_cargo_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_cargo_proto_goTypes,
		DependencyIndexes: file_cargo_proto_depIdxs,
		MessageInfos:      file_cargo_proto_msgTypes,
	}.Build()
	File_cargo_proto = out.File
	file_cargo_proto_rawDesc = nil
	file_cargo_proto_goTypes = nil
	file_cargo_proto_depIdxs = nil
}
