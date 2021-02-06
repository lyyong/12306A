// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.14.0
// source: indent.proto

package indentRPC

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

// 创建订单，存入 redis 设置过期时间，等待用户支付
type CreateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId         int32  `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	TrainId        int32  `protobuf:"varint,2,opt,name=train_id,json=trainId,proto3" json:"train_id,omitempty"`
	StartStationId int32  `protobuf:"varint,3,opt,name=start_station_id,json=startStationId,proto3" json:"start_station_id,omitempty"`
	StartTime      string `protobuf:"bytes,4,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"`
	DestStationId  int32  `protobuf:"varint,5,opt,name=dest_station_id,json=destStationId,proto3" json:"dest_station_id,omitempty"`
	ArriveTime     string `protobuf:"bytes,6,opt,name=arrive_time,json=arriveTime,proto3" json:"arrive_time,omitempty"`
	Date           string `protobuf:"bytes,7,opt,name=date,proto3" json:"date,omitempty"`
	ExpiredTime    int32  `protobuf:"varint,8,opt,name=expired_time,json=expiredTime,proto3" json:"expired_time,omitempty"`
	TicketNumber   int32  `protobuf:"varint,9,opt,name=ticket_number,json=ticketNumber,proto3" json:"ticket_number,omitempty"`
	Amount         int32  `protobuf:"varint,10,opt,name=amount,proto3" json:"amount,omitempty"`
}

func (x *CreateRequest) Reset() {
	*x = CreateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_indent_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateRequest) ProtoMessage() {}

func (x *CreateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_indent_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateRequest.ProtoReflect.Descriptor instead.
func (*CreateRequest) Descriptor() ([]byte, []int) {
	return file_indent_proto_rawDescGZIP(), []int{0}
}

func (x *CreateRequest) GetUserId() int32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *CreateRequest) GetTrainId() int32 {
	if x != nil {
		return x.TrainId
	}
	return 0
}

func (x *CreateRequest) GetStartStationId() int32 {
	if x != nil {
		return x.StartStationId
	}
	return 0
}

func (x *CreateRequest) GetStartTime() string {
	if x != nil {
		return x.StartTime
	}
	return ""
}

func (x *CreateRequest) GetDestStationId() int32 {
	if x != nil {
		return x.DestStationId
	}
	return 0
}

func (x *CreateRequest) GetArriveTime() string {
	if x != nil {
		return x.ArriveTime
	}
	return ""
}

func (x *CreateRequest) GetDate() string {
	if x != nil {
		return x.Date
	}
	return ""
}

func (x *CreateRequest) GetExpiredTime() int32 {
	if x != nil {
		return x.ExpiredTime
	}
	return 0
}

func (x *CreateRequest) GetTicketNumber() int32 {
	if x != nil {
		return x.TicketNumber
	}
	return 0
}

func (x *CreateRequest) GetAmount() int32 {
	if x != nil {
		return x.Amount
	}
	return 0
}

type CreateResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IndentOuterId string `protobuf:"bytes,1,opt,name=indent_outer_id,json=indentOuterId,proto3" json:"indent_outer_id,omitempty"`
}

func (x *CreateResponse) Reset() {
	*x = CreateResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_indent_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateResponse) ProtoMessage() {}

func (x *CreateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_indent_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateResponse.ProtoReflect.Descriptor instead.
func (*CreateResponse) Descriptor() ([]byte, []int) {
	return file_indent_proto_rawDescGZIP(), []int{1}
}

func (x *CreateResponse) GetIndentOuterId() string {
	if x != nil {
		return x.IndentOuterId
	}
	return ""
}

// 用户支付后调用该服务写入将 redis 中订单信息数据库
type PayRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId        int32  `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	IndentOuterId string `protobuf:"bytes,2,opt,name=indent_outer_id,json=indentOuterId,proto3" json:"indent_outer_id,omitempty"`
	PayAmount     int32  `protobuf:"varint,3,opt,name=pay_amount,json=payAmount,proto3" json:"pay_amount,omitempty"`
}

func (x *PayRequest) Reset() {
	*x = PayRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_indent_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PayRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PayRequest) ProtoMessage() {}

func (x *PayRequest) ProtoReflect() protoreflect.Message {
	mi := &file_indent_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PayRequest.ProtoReflect.Descriptor instead.
func (*PayRequest) Descriptor() ([]byte, []int) {
	return file_indent_proto_rawDescGZIP(), []int{2}
}

func (x *PayRequest) GetUserId() int32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *PayRequest) GetIndentOuterId() string {
	if x != nil {
		return x.IndentOuterId
	}
	return ""
}

func (x *PayRequest) GetPayAmount() int32 {
	if x != nil {
		return x.PayAmount
	}
	return 0
}

type PayResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IsOk bool `protobuf:"varint,1,opt,name=is_ok,json=isOk,proto3" json:"is_ok,omitempty"`
}

func (x *PayResponse) Reset() {
	*x = PayResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_indent_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PayResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PayResponse) ProtoMessage() {}

func (x *PayResponse) ProtoReflect() protoreflect.Message {
	mi := &file_indent_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PayResponse.ProtoReflect.Descriptor instead.
func (*PayResponse) Descriptor() ([]byte, []int) {
	return file_indent_proto_rawDescGZIP(), []int{3}
}

func (x *PayResponse) GetIsOk() bool {
	if x != nil {
		return x.IsOk
	}
	return false
}

// 根据用户 id 查询是否有未完成的订单
type UnfinishedRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId int32 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

func (x *UnfinishedRequest) Reset() {
	*x = UnfinishedRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_indent_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UnfinishedRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UnfinishedRequest) ProtoMessage() {}

func (x *UnfinishedRequest) ProtoReflect() protoreflect.Message {
	mi := &file_indent_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UnfinishedRequest.ProtoReflect.Descriptor instead.
func (*UnfinishedRequest) Descriptor() ([]byte, []int) {
	return file_indent_proto_rawDescGZIP(), []int{4}
}

func (x *UnfinishedRequest) GetUserId() int32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

type UnfinishedResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	HasUnfinishedIndent bool `protobuf:"varint,1,opt,name=has_unfinished_indent,json=hasUnfinishedIndent,proto3" json:"has_unfinished_indent,omitempty"`
}

func (x *UnfinishedResponse) Reset() {
	*x = UnfinishedResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_indent_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UnfinishedResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UnfinishedResponse) ProtoMessage() {}

func (x *UnfinishedResponse) ProtoReflect() protoreflect.Message {
	mi := &file_indent_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UnfinishedResponse.ProtoReflect.Descriptor instead.
func (*UnfinishedResponse) Descriptor() ([]byte, []int) {
	return file_indent_proto_rawDescGZIP(), []int{5}
}

func (x *UnfinishedResponse) GetHasUnfinishedIndent() bool {
	if x != nil {
		return x.HasUnfinishedIndent
	}
	return false
}

var File_indent_proto protoreflect.FileDescriptor

var file_indent_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x69, 0x6e, 0x64, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09,
	0x69, 0x6e, 0x64, 0x65, 0x6e, 0x74, 0x52, 0x50, 0x43, 0x22, 0xc9, 0x02, 0x0a, 0x0d, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x75,
	0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x75, 0x73,
	0x65, 0x72, 0x49, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x74, 0x72, 0x61, 0x69, 0x6e, 0x5f, 0x69, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x74, 0x72, 0x61, 0x69, 0x6e, 0x49, 0x64, 0x12,
	0x28, 0x0a, 0x10, 0x73, 0x74, 0x61, 0x72, 0x74, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0e, 0x73, 0x74, 0x61, 0x72, 0x74,
	0x53, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x74, 0x61,
	0x72, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73,
	0x74, 0x61, 0x72, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x26, 0x0a, 0x0f, 0x64, 0x65, 0x73, 0x74,
	0x5f, 0x73, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x0d, 0x64, 0x65, 0x73, 0x74, 0x53, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64,
	0x12, 0x1f, 0x0a, 0x0b, 0x61, 0x72, 0x72, 0x69, 0x76, 0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x61, 0x72, 0x72, 0x69, 0x76, 0x65, 0x54, 0x69, 0x6d,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x64, 0x61, 0x74, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x65, 0x78, 0x70, 0x69, 0x72, 0x65, 0x64,
	0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0b, 0x65, 0x78, 0x70,
	0x69, 0x72, 0x65, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x74, 0x69, 0x63, 0x6b,
	0x65, 0x74, 0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x09, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x0c, 0x74, 0x69, 0x63, 0x6b, 0x65, 0x74, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x16, 0x0a,
	0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x61,
	0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x38, 0x0a, 0x0e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x26, 0x0a, 0x0f, 0x69, 0x6e, 0x64, 0x65, 0x6e,
	0x74, 0x5f, 0x6f, 0x75, 0x74, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0d, 0x69, 0x6e, 0x64, 0x65, 0x6e, 0x74, 0x4f, 0x75, 0x74, 0x65, 0x72, 0x49, 0x64, 0x22,
	0x6c, 0x0a, 0x0a, 0x50, 0x61, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a,
	0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06,
	0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x26, 0x0a, 0x0f, 0x69, 0x6e, 0x64, 0x65, 0x6e, 0x74,
	0x5f, 0x6f, 0x75, 0x74, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0d, 0x69, 0x6e, 0x64, 0x65, 0x6e, 0x74, 0x4f, 0x75, 0x74, 0x65, 0x72, 0x49, 0x64, 0x12, 0x1d,
	0x0a, 0x0a, 0x70, 0x61, 0x79, 0x5f, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x09, 0x70, 0x61, 0x79, 0x41, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x22, 0x0a,
	0x0b, 0x50, 0x61, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x13, 0x0a, 0x05,
	0x69, 0x73, 0x5f, 0x6f, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x04, 0x69, 0x73, 0x4f,
	0x6b, 0x22, 0x2c, 0x0a, 0x11, 0x55, 0x6e, 0x66, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x65, 0x64, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22,
	0x48, 0x0a, 0x12, 0x55, 0x6e, 0x66, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x65, 0x64, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x32, 0x0a, 0x15, 0x68, 0x61, 0x73, 0x5f, 0x75, 0x6e, 0x66,
	0x69, 0x6e, 0x69, 0x73, 0x68, 0x65, 0x64, 0x5f, 0x69, 0x6e, 0x64, 0x65, 0x6e, 0x74, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x13, 0x68, 0x61, 0x73, 0x55, 0x6e, 0x66, 0x69, 0x6e, 0x69, 0x73,
	0x68, 0x65, 0x64, 0x49, 0x6e, 0x64, 0x65, 0x6e, 0x74, 0x32, 0xea, 0x01, 0x0a, 0x0d, 0x49, 0x6e,
	0x64, 0x65, 0x6e, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x45, 0x0a, 0x0c, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x49, 0x6e, 0x64, 0x65, 0x6e, 0x74, 0x12, 0x18, 0x2e, 0x69, 0x6e,
	0x64, 0x65, 0x6e, 0x74, 0x52, 0x50, 0x43, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x69, 0x6e, 0x64, 0x65, 0x6e, 0x74, 0x52, 0x50,
	0x43, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x12, 0x3c, 0x0a, 0x09, 0x50, 0x61, 0x79, 0x49, 0x6e, 0x64, 0x65, 0x6e, 0x74, 0x12,
	0x15, 0x2e, 0x69, 0x6e, 0x64, 0x65, 0x6e, 0x74, 0x52, 0x50, 0x43, 0x2e, 0x50, 0x61, 0x79, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x69, 0x6e, 0x64, 0x65, 0x6e, 0x74, 0x52,
	0x50, 0x43, 0x2e, 0x50, 0x61, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x12, 0x54, 0x0a, 0x13, 0x48, 0x61, 0x73, 0x55, 0x6e, 0x66, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x65,
	0x64, 0x49, 0x6e, 0x64, 0x65, 0x6e, 0x74, 0x12, 0x1c, 0x2e, 0x69, 0x6e, 0x64, 0x65, 0x6e, 0x74,
	0x52, 0x50, 0x43, 0x2e, 0x55, 0x6e, 0x66, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x65, 0x64, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x69, 0x6e, 0x64, 0x65, 0x6e, 0x74, 0x52, 0x50,
	0x43, 0x2e, 0x55, 0x6e, 0x66, 0x69, 0x6e, 0x69, 0x73, 0x68, 0x65, 0x64, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_indent_proto_rawDescOnce sync.Once
	file_indent_proto_rawDescData = file_indent_proto_rawDesc
)

func file_indent_proto_rawDescGZIP() []byte {
	file_indent_proto_rawDescOnce.Do(func() {
		file_indent_proto_rawDescData = protoimpl.X.CompressGZIP(file_indent_proto_rawDescData)
	})
	return file_indent_proto_rawDescData
}

var file_indent_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_indent_proto_goTypes = []interface{}{
	(*CreateRequest)(nil),      // 0: indentRPC.CreateRequest
	(*CreateResponse)(nil),     // 1: indentRPC.CreateResponse
	(*PayRequest)(nil),         // 2: indentRPC.PayRequest
	(*PayResponse)(nil),        // 3: indentRPC.PayResponse
	(*UnfinishedRequest)(nil),  // 4: indentRPC.UnfinishedRequest
	(*UnfinishedResponse)(nil), // 5: indentRPC.UnfinishedResponse
}
var file_indent_proto_depIdxs = []int32{
	0, // 0: indentRPC.IndentService.CreateIndent:input_type -> indentRPC.CreateRequest
	2, // 1: indentRPC.IndentService.PayIndent:input_type -> indentRPC.PayRequest
	4, // 2: indentRPC.IndentService.HasUnfinishedIndent:input_type -> indentRPC.UnfinishedRequest
	1, // 3: indentRPC.IndentService.CreateIndent:output_type -> indentRPC.CreateResponse
	3, // 4: indentRPC.IndentService.PayIndent:output_type -> indentRPC.PayResponse
	5, // 5: indentRPC.IndentService.HasUnfinishedIndent:output_type -> indentRPC.UnfinishedResponse
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_indent_proto_init() }
func file_indent_proto_init() {
	if File_indent_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_indent_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateRequest); i {
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
		file_indent_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateResponse); i {
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
		file_indent_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PayRequest); i {
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
		file_indent_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PayResponse); i {
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
		file_indent_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UnfinishedRequest); i {
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
		file_indent_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UnfinishedResponse); i {
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
			RawDescriptor: file_indent_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_indent_proto_goTypes,
		DependencyIndexes: file_indent_proto_depIdxs,
		MessageInfos:      file_indent_proto_msgTypes,
	}.Build()
	File_indent_proto = out.File
	file_indent_proto_rawDesc = nil
	file_indent_proto_goTypes = nil
	file_indent_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// IndentServiceClient is the client API for IndentService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type IndentServiceClient interface {
	CreateIndent(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateResponse, error)
	PayIndent(ctx context.Context, in *PayRequest, opts ...grpc.CallOption) (*PayResponse, error)
	HasUnfinishedIndent(ctx context.Context, in *UnfinishedRequest, opts ...grpc.CallOption) (*UnfinishedResponse, error)
}

type indentServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewIndentServiceClient(cc grpc.ClientConnInterface) IndentServiceClient {
	return &indentServiceClient{cc}
}

func (c *indentServiceClient) CreateIndent(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateResponse, error) {
	out := new(CreateResponse)
	err := c.cc.Invoke(ctx, "/indentRPC.IndentService/CreateIndent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *indentServiceClient) PayIndent(ctx context.Context, in *PayRequest, opts ...grpc.CallOption) (*PayResponse, error) {
	out := new(PayResponse)
	err := c.cc.Invoke(ctx, "/indentRPC.IndentService/PayIndent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *indentServiceClient) HasUnfinishedIndent(ctx context.Context, in *UnfinishedRequest, opts ...grpc.CallOption) (*UnfinishedResponse, error) {
	out := new(UnfinishedResponse)
	err := c.cc.Invoke(ctx, "/indentRPC.IndentService/HasUnfinishedIndent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// IndentServiceServer is the server API for IndentService service.
type IndentServiceServer interface {
	CreateIndent(context.Context, *CreateRequest) (*CreateResponse, error)
	PayIndent(context.Context, *PayRequest) (*PayResponse, error)
	HasUnfinishedIndent(context.Context, *UnfinishedRequest) (*UnfinishedResponse, error)
}

// UnimplementedIndentServiceServer can be embedded to have forward compatible implementations.
type UnimplementedIndentServiceServer struct {
}

func (*UnimplementedIndentServiceServer) CreateIndent(context.Context, *CreateRequest) (*CreateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateIndent not implemented")
}
func (*UnimplementedIndentServiceServer) PayIndent(context.Context, *PayRequest) (*PayResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PayIndent not implemented")
}
func (*UnimplementedIndentServiceServer) HasUnfinishedIndent(context.Context, *UnfinishedRequest) (*UnfinishedResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HasUnfinishedIndent not implemented")
}

func RegisterIndentServiceServer(s *grpc.Server, srv IndentServiceServer) {
	s.RegisterService(&_IndentService_serviceDesc, srv)
}

func _IndentService_CreateIndent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IndentServiceServer).CreateIndent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/indentRPC.IndentService/CreateIndent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IndentServiceServer).CreateIndent(ctx, req.(*CreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IndentService_PayIndent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PayRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IndentServiceServer).PayIndent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/indentRPC.IndentService/PayIndent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IndentServiceServer).PayIndent(ctx, req.(*PayRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IndentService_HasUnfinishedIndent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UnfinishedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IndentServiceServer).HasUnfinishedIndent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/indentRPC.IndentService/HasUnfinishedIndent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IndentServiceServer).HasUnfinishedIndent(ctx, req.(*UnfinishedRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _IndentService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "indentRPC.IndentService",
	HandlerType: (*IndentServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateIndent",
			Handler:    _IndentService_CreateIndent_Handler,
		},
		{
			MethodName: "PayIndent",
			Handler:    _IndentService_PayIndent_Handler,
		},
		{
			MethodName: "HasUnfinishedIndent",
			Handler:    _IndentService_HasUnfinishedIndent_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "indent.proto",
}