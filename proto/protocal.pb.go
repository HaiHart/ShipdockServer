// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.21.1
// source: proto/protocal.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Header struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Time *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=time,proto3" json:"time,omitempty"`
}

func (x *Header) Reset() {
	*x = Header{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_protocal_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Header) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Header) ProtoMessage() {}

func (x *Header) ProtoReflect() protoreflect.Message {
	mi := &file_proto_protocal_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Header.ProtoReflect.Descriptor instead.
func (*Header) Descriptor() ([]byte, []int) {
	return file_proto_protocal_proto_rawDescGZIP(), []int{0}
}

func (x *Header) GetTime() *timestamppb.Timestamp {
	if x != nil {
		return x.Time
	}
	return nil
}

type Container struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id       string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Time     *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=time,proto3" json:"time,omitempty"`
	Place    int32                  `protobuf:"varint,3,opt,name=place,proto3" json:"place,omitempty"`
	Name     string                 `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`
	NewPlace int32                  `protobuf:"varint,5,opt,name=new_place,json=newPlace,proto3" json:"new_place,omitempty"`
}

func (x *Container) Reset() {
	*x = Container{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_protocal_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Container) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Container) ProtoMessage() {}

func (x *Container) ProtoReflect() protoreflect.Message {
	mi := &file_proto_protocal_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Container.ProtoReflect.Descriptor instead.
func (*Container) Descriptor() ([]byte, []int) {
	return file_proto_protocal_proto_rawDescGZIP(), []int{1}
}

func (x *Container) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Container) GetTime() *timestamppb.Timestamp {
	if x != nil {
		return x.Time
	}
	return nil
}

func (x *Container) GetPlace() int32 {
	if x != nil {
		return x.Place
	}
	return 0
}

func (x *Container) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Container) GetNewPlace() int32 {
	if x != nil {
		return x.NewPlace
	}
	return 0
}

type Pack struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	List []*Container `protobuf:"bytes,1,rep,name=list,proto3" json:"list,omitempty"`
	Err  string       `protobuf:"bytes,2,opt,name=err,proto3" json:"err,omitempty"`
	Swap bool         `protobuf:"varint,3,opt,name=swap,proto3" json:"swap,omitempty"`
}

func (x *Pack) Reset() {
	*x = Pack{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_protocal_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Pack) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Pack) ProtoMessage() {}

func (x *Pack) ProtoReflect() protoreflect.Message {
	mi := &file_proto_protocal_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Pack.ProtoReflect.Descriptor instead.
func (*Pack) Descriptor() ([]byte, []int) {
	return file_proto_protocal_proto_rawDescGZIP(), []int{2}
}

func (x *Pack) GetList() []*Container {
	if x != nil {
		return x.List
	}
	return nil
}

func (x *Pack) GetErr() string {
	if x != nil {
		return x.Err
	}
	return ""
}

func (x *Pack) GetSwap() bool {
	if x != nil {
		return x.Swap
	}
	return false
}

type ShipAccess struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ShipName string  `protobuf:"bytes,1,opt,name=ship_name,json=shipName,proto3" json:"ship_name,omitempty"`
	ManId    string  `protobuf:"bytes,2,opt,name=man_id,json=manId,proto3" json:"man_id,omitempty"`
	Floor    int32   `protobuf:"varint,3,opt,name=floor,proto3" json:"floor,omitempty"`
	Head     *Header `protobuf:"bytes,5,opt,name=head,proto3" json:"head,omitempty"`
}

func (x *ShipAccess) Reset() {
	*x = ShipAccess{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_protocal_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ShipAccess) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ShipAccess) ProtoMessage() {}

func (x *ShipAccess) ProtoReflect() protoreflect.Message {
	mi := &file_proto_protocal_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ShipAccess.ProtoReflect.Descriptor instead.
func (*ShipAccess) Descriptor() ([]byte, []int) {
	return file_proto_protocal_proto_rawDescGZIP(), []int{3}
}

func (x *ShipAccess) GetShipName() string {
	if x != nil {
		return x.ShipName
	}
	return ""
}

func (x *ShipAccess) GetManId() string {
	if x != nil {
		return x.ManId
	}
	return ""
}

func (x *ShipAccess) GetFloor() int32 {
	if x != nil {
		return x.Floor
	}
	return 0
}

func (x *ShipAccess) GetHead() *Header {
	if x != nil {
		return x.Head
	}
	return nil
}

type ShipResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Size          int32        `protobuf:"varint,1,opt,name=size,proto3" json:"size,omitempty"`
	Err           string       `protobuf:"bytes,2,opt,name=err,proto3" json:"err,omitempty"`
	Head          *Header      `protobuf:"bytes,3,opt,name=head,proto3" json:"head,omitempty"`
	ContainerList []*Container `protobuf:"bytes,4,rep,name=container_list,json=containerList,proto3" json:"container_list,omitempty"`
}

func (x *ShipResponse) Reset() {
	*x = ShipResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_protocal_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ShipResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ShipResponse) ProtoMessage() {}

func (x *ShipResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_protocal_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ShipResponse.ProtoReflect.Descriptor instead.
func (*ShipResponse) Descriptor() ([]byte, []int) {
	return file_proto_protocal_proto_rawDescGZIP(), []int{4}
}

func (x *ShipResponse) GetSize() int32 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *ShipResponse) GetErr() string {
	if x != nil {
		return x.Err
	}
	return ""
}

func (x *ShipResponse) GetHead() *Header {
	if x != nil {
		return x.Head
	}
	return nil
}

func (x *ShipResponse) GetContainerList() []*Container {
	if x != nil {
		return x.ContainerList
	}
	return nil
}

var File_proto_protocal_proto protoreflect.FileDescriptor

var file_proto_protocal_proto_rawDesc = []byte{
	0x0a, 0x14, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x61, 0x6c,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x38,
	0x0a, 0x06, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x2e, 0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x52, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x22, 0x92, 0x01, 0x0a, 0x09, 0x43, 0x6f, 0x6e,
	0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x2e, 0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x52, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x70, 0x6c, 0x61, 0x63, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x70, 0x6c, 0x61, 0x63, 0x65, 0x12, 0x12, 0x0a, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x1b, 0x0a, 0x09, 0x6e, 0x65, 0x77, 0x5f, 0x70, 0x6c, 0x61, 0x63, 0x65, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x08, 0x6e, 0x65, 0x77, 0x50, 0x6c, 0x61, 0x63, 0x65, 0x22, 0x52, 0x0a,
	0x04, 0x50, 0x61, 0x63, 0x6b, 0x12, 0x24, 0x0a, 0x04, 0x6c, 0x69, 0x73, 0x74, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x6f, 0x6e, 0x74,
	0x61, 0x69, 0x6e, 0x65, 0x72, 0x52, 0x04, 0x6c, 0x69, 0x73, 0x74, 0x12, 0x10, 0x0a, 0x03, 0x65,
	0x72, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x65, 0x72, 0x72, 0x12, 0x12, 0x0a,
	0x04, 0x73, 0x77, 0x61, 0x70, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x04, 0x73, 0x77, 0x61,
	0x70, 0x22, 0x79, 0x0a, 0x0a, 0x53, 0x68, 0x69, 0x70, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12,
	0x1b, 0x0a, 0x09, 0x73, 0x68, 0x69, 0x70, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x73, 0x68, 0x69, 0x70, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x15, 0x0a, 0x06,
	0x6d, 0x61, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6d, 0x61,
	0x6e, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x66, 0x6c, 0x6f, 0x6f, 0x72, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x05, 0x66, 0x6c, 0x6f, 0x6f, 0x72, 0x12, 0x21, 0x0a, 0x04, 0x68, 0x65, 0x61,
	0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x52, 0x04, 0x68, 0x65, 0x61, 0x64, 0x22, 0x90, 0x01, 0x0a,
	0x0c, 0x53, 0x68, 0x69, 0x70, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a,
	0x04, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x73, 0x69, 0x7a,
	0x65, 0x12, 0x10, 0x0a, 0x03, 0x65, 0x72, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x65, 0x72, 0x72, 0x12, 0x21, 0x0a, 0x04, 0x68, 0x65, 0x61, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x0d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72,
	0x52, 0x04, 0x68, 0x65, 0x61, 0x64, 0x12, 0x37, 0x0a, 0x0e, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x69,
	0x6e, 0x65, 0x72, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x10,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72,
	0x52, 0x0d, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x4c, 0x69, 0x73, 0x74, 0x32,
	0x69, 0x0a, 0x03, 0x43, 0x6f, 0x6d, 0x12, 0x2d, 0x0a, 0x0d, 0x4d, 0x6f, 0x76, 0x65, 0x43, 0x6f,
	0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72, 0x12, 0x0b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x50, 0x61, 0x63, 0x6b, 0x1a, 0x0b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x61, 0x63,
	0x6b, 0x28, 0x01, 0x30, 0x01, 0x12, 0x33, 0x0a, 0x09, 0x46, 0x65, 0x74, 0x63, 0x68, 0x53, 0x68,
	0x69, 0x70, 0x12, 0x11, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x68, 0x69, 0x70, 0x41,
	0x63, 0x63, 0x65, 0x73, 0x73, 0x1a, 0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x68,
	0x69, 0x70, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x23, 0x5a, 0x21, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x48, 0x61, 0x69, 0x48, 0x61, 0x72, 0x74,
	0x2f, 0x73, 0x68, 0x69, 0x70, 0x64, 0x6f, 0x63, 0x6b, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_protocal_proto_rawDescOnce sync.Once
	file_proto_protocal_proto_rawDescData = file_proto_protocal_proto_rawDesc
)

func file_proto_protocal_proto_rawDescGZIP() []byte {
	file_proto_protocal_proto_rawDescOnce.Do(func() {
		file_proto_protocal_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_protocal_proto_rawDescData)
	})
	return file_proto_protocal_proto_rawDescData
}

var file_proto_protocal_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_proto_protocal_proto_goTypes = []interface{}{
	(*Header)(nil),                // 0: proto.Header
	(*Container)(nil),             // 1: proto.Container
	(*Pack)(nil),                  // 2: proto.Pack
	(*ShipAccess)(nil),            // 3: proto.ShipAccess
	(*ShipResponse)(nil),          // 4: proto.ShipResponse
	(*timestamppb.Timestamp)(nil), // 5: google.protobuf.Timestamp
}
var file_proto_protocal_proto_depIdxs = []int32{
	5, // 0: proto.Header.time:type_name -> google.protobuf.Timestamp
	5, // 1: proto.Container.time:type_name -> google.protobuf.Timestamp
	1, // 2: proto.Pack.list:type_name -> proto.Container
	0, // 3: proto.ShipAccess.head:type_name -> proto.Header
	0, // 4: proto.ShipResponse.head:type_name -> proto.Header
	1, // 5: proto.ShipResponse.container_list:type_name -> proto.Container
	2, // 6: proto.Com.MoveContainer:input_type -> proto.Pack
	3, // 7: proto.Com.FetchShip:input_type -> proto.ShipAccess
	2, // 8: proto.Com.MoveContainer:output_type -> proto.Pack
	4, // 9: proto.Com.FetchShip:output_type -> proto.ShipResponse
	8, // [8:10] is the sub-list for method output_type
	6, // [6:8] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_proto_protocal_proto_init() }
func file_proto_protocal_proto_init() {
	if File_proto_protocal_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_protocal_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Header); i {
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
		file_proto_protocal_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Container); i {
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
		file_proto_protocal_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Pack); i {
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
		file_proto_protocal_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ShipAccess); i {
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
		file_proto_protocal_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ShipResponse); i {
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
			RawDescriptor: file_proto_protocal_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_protocal_proto_goTypes,
		DependencyIndexes: file_proto_protocal_proto_depIdxs,
		MessageInfos:      file_proto_protocal_proto_msgTypes,
	}.Build()
	File_proto_protocal_proto = out.File
	file_proto_protocal_proto_rawDesc = nil
	file_proto_protocal_proto_goTypes = nil
	file_proto_protocal_proto_depIdxs = nil
}
