// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v3.20.1
// source: goods.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// 瀹氫箟璇锋眰娑堟伅 GetGoodsByRoomReq锛岀敤浜庤幏鍙栫洿鎾棿鍟嗗搧鍒楄〃
type GetGoodsByRoomReq struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        int64                  `protobuf:"varint,1,opt,name=UserId,proto3" json:"UserId,omitempty"` // 鐢ㄦ埛 ID
	RoomId        int64                  `protobuf:"varint,2,opt,name=RoomId,proto3" json:"RoomId,omitempty"` // 鐩存挱闂� ID
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetGoodsByRoomReq) Reset() {
	*x = GetGoodsByRoomReq{}
	mi := &file_goods_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetGoodsByRoomReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetGoodsByRoomReq) ProtoMessage() {}

func (x *GetGoodsByRoomReq) ProtoReflect() protoreflect.Message {
	mi := &file_goods_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetGoodsByRoomReq.ProtoReflect.Descriptor instead.
func (*GetGoodsByRoomReq) Descriptor() ([]byte, []int) {
	return file_goods_proto_rawDescGZIP(), []int{0}
}

func (x *GetGoodsByRoomReq) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *GetGoodsByRoomReq) GetRoomId() int64 {
	if x != nil {
		return x.RoomId
	}
	return 0
}

// 瀹氫箟鍝嶅簲娑堟伅 GoodsListResp锛岀敤浜庤繑鍥炲晢鍝佸垪琛�
type GoodsListResp struct {
	state          protoimpl.MessageState `protogen:"open.v1"`
	CurrentGoodsId int64                  `protobuf:"varint,1,opt,name=CurrentGoodsId,proto3" json:"CurrentGoodsId,omitempty"` // 褰撳墠鍟嗗搧鐨� ID
	Data           []*GoodsInfo           `protobuf:"bytes,2,rep,name=Data,proto3" json:"Data,omitempty"`                      // 鍟嗗搧鍒楄〃锛屽寘鍚涓� GoodsInfo 娑堟伅
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *GoodsListResp) Reset() {
	*x = GoodsListResp{}
	mi := &file_goods_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GoodsListResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GoodsListResp) ProtoMessage() {}

func (x *GoodsListResp) ProtoReflect() protoreflect.Message {
	mi := &file_goods_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GoodsListResp.ProtoReflect.Descriptor instead.
func (*GoodsListResp) Descriptor() ([]byte, []int) {
	return file_goods_proto_rawDescGZIP(), []int{1}
}

func (x *GoodsListResp) GetCurrentGoodsId() int64 {
	if x != nil {
		return x.CurrentGoodsId
	}
	return 0
}

func (x *GoodsListResp) GetData() []*GoodsInfo {
	if x != nil {
		return x.Data
	}
	return nil
}

// 瀹氫箟鍟嗗搧鍒楄〃椤电殑鏁版嵁缁撴瀯 GoodsInfo
type GoodsInfo struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	GoodsId       int64                  `protobuf:"varint,1,opt,name=GoodsId,proto3" json:"GoodsId,omitempty"`        // 鍟嗗搧 ID
	CategoryId    int64                  `protobuf:"varint,2,opt,name=CategoryId,proto3" json:"CategoryId,omitempty"`  // 鍒嗙被 ID
	Status        int32                  `protobuf:"varint,3,opt,name=Status,proto3" json:"Status,omitempty"`          // 鍟嗗搧鐘舵€�
	Title         string                 `protobuf:"bytes,4,opt,name=Title,proto3" json:"Title,omitempty"`             // 鍟嗗搧鏍囬
	MarketPrice   string                 `protobuf:"bytes,5,opt,name=MarketPrice,proto3" json:"MarketPrice,omitempty"` // 甯傚満浠锋牸
	Price         string                 `protobuf:"bytes,6,opt,name=Price,proto3" json:"Price,omitempty"`             // 閿€鍞环鏍�
	Brief         string                 `protobuf:"bytes,7,opt,name=Brief,proto3" json:"Brief,omitempty"`             // 鍟嗗搧绠€浠�
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GoodsInfo) Reset() {
	*x = GoodsInfo{}
	mi := &file_goods_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GoodsInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GoodsInfo) ProtoMessage() {}

func (x *GoodsInfo) ProtoReflect() protoreflect.Message {
	mi := &file_goods_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GoodsInfo.ProtoReflect.Descriptor instead.
func (*GoodsInfo) Descriptor() ([]byte, []int) {
	return file_goods_proto_rawDescGZIP(), []int{2}
}

func (x *GoodsInfo) GetGoodsId() int64 {
	if x != nil {
		return x.GoodsId
	}
	return 0
}

func (x *GoodsInfo) GetCategoryId() int64 {
	if x != nil {
		return x.CategoryId
	}
	return 0
}

func (x *GoodsInfo) GetStatus() int32 {
	if x != nil {
		return x.Status
	}
	return 0
}

func (x *GoodsInfo) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *GoodsInfo) GetMarketPrice() string {
	if x != nil {
		return x.MarketPrice
	}
	return ""
}

func (x *GoodsInfo) GetPrice() string {
	if x != nil {
		return x.Price
	}
	return ""
}

func (x *GoodsInfo) GetBrief() string {
	if x != nil {
		return x.Brief
	}
	return ""
}

// 瀹氫箟璇锋眰娑堟伅 GetGoodsDetailReq锛岀敤浜庤幏鍙栧晢鍝佽鎯�
type GetGoodsDetailReq struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	GoodsId       int64                  `protobuf:"varint,1,opt,name=GoodsId,proto3" json:"GoodsId,omitempty"` // 鍟嗗搧 ID
	UserId        int64                  `protobuf:"varint,2,opt,name=UserId,proto3" json:"UserId,omitempty"`   // 鐢ㄦ埛 ID
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetGoodsDetailReq) Reset() {
	*x = GetGoodsDetailReq{}
	mi := &file_goods_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetGoodsDetailReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetGoodsDetailReq) ProtoMessage() {}

func (x *GetGoodsDetailReq) ProtoReflect() protoreflect.Message {
	mi := &file_goods_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetGoodsDetailReq.ProtoReflect.Descriptor instead.
func (*GetGoodsDetailReq) Descriptor() ([]byte, []int) {
	return file_goods_proto_rawDescGZIP(), []int{3}
}

func (x *GetGoodsDetailReq) GetGoodsId() int64 {
	if x != nil {
		return x.GoodsId
	}
	return 0
}

func (x *GetGoodsDetailReq) GetUserId() int64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

// 瀹氫箟鍝嶅簲娑堟伅 GoodsDetail锛岀敤浜庤繑鍥炲晢鍝佽鎯�
type GoodsDetail struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	GoodsId       int64                  `protobuf:"varint,1,opt,name=GoodsId,proto3" json:"GoodsId,omitempty"`        // 鍟嗗搧 ID
	CategoryId    int64                  `protobuf:"varint,2,opt,name=CategoryId,proto3" json:"CategoryId,omitempty"`  // 鍒嗙被 ID
	Status        int32                  `protobuf:"varint,3,opt,name=Status,proto3" json:"Status,omitempty"`          // 鍟嗗搧鐘舵€�
	Title         string                 `protobuf:"bytes,4,opt,name=Title,proto3" json:"Title,omitempty"`             // 鍟嗗搧鏍囬
	Code          string                 `protobuf:"bytes,5,opt,name=Code,proto3" json:"Code,omitempty"`               // 鍟嗗搧缂栫爜
	BrandName     string                 `protobuf:"bytes,6,opt,name=BrandName,proto3" json:"BrandName,omitempty"`     // 鍝佺墝鍚嶇О
	MarketPrice   string                 `protobuf:"bytes,7,opt,name=MarketPrice,proto3" json:"MarketPrice,omitempty"` // 甯傚満浠锋牸
	Price         string                 `protobuf:"bytes,8,opt,name=Price,proto3" json:"Price,omitempty"`             // 閿€鍞环鏍�
	Brief         string                 `protobuf:"bytes,9,opt,name=Brief,proto3" json:"Brief,omitempty"`             // 鍟嗗搧绠€浠�
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GoodsDetail) Reset() {
	*x = GoodsDetail{}
	mi := &file_goods_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GoodsDetail) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GoodsDetail) ProtoMessage() {}

func (x *GoodsDetail) ProtoReflect() protoreflect.Message {
	mi := &file_goods_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GoodsDetail.ProtoReflect.Descriptor instead.
func (*GoodsDetail) Descriptor() ([]byte, []int) {
	return file_goods_proto_rawDescGZIP(), []int{4}
}

func (x *GoodsDetail) GetGoodsId() int64 {
	if x != nil {
		return x.GoodsId
	}
	return 0
}

func (x *GoodsDetail) GetCategoryId() int64 {
	if x != nil {
		return x.CategoryId
	}
	return 0
}

func (x *GoodsDetail) GetStatus() int32 {
	if x != nil {
		return x.Status
	}
	return 0
}

func (x *GoodsDetail) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *GoodsDetail) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

func (x *GoodsDetail) GetBrandName() string {
	if x != nil {
		return x.BrandName
	}
	return ""
}

func (x *GoodsDetail) GetMarketPrice() string {
	if x != nil {
		return x.MarketPrice
	}
	return ""
}

func (x *GoodsDetail) GetPrice() string {
	if x != nil {
		return x.Price
	}
	return ""
}

func (x *GoodsDetail) GetBrief() string {
	if x != nil {
		return x.Brief
	}
	return ""
}

var File_goods_proto protoreflect.FileDescriptor

var file_goods_proto_rawDesc = string([]byte{
	0x0a, 0x0b, 0x67, 0x6f, 0x6f, 0x64, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x43, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x47, 0x6f, 0x6f, 0x64, 0x73,
	0x42, 0x79, 0x52, 0x6f, 0x6f, 0x6d, 0x52, 0x65, 0x71, 0x12, 0x16, 0x0a, 0x06, 0x55, 0x73, 0x65,
	0x72, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x55, 0x73, 0x65, 0x72, 0x49,
	0x64, 0x12, 0x16, 0x0a, 0x06, 0x52, 0x6f, 0x6f, 0x6d, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x06, 0x52, 0x6f, 0x6f, 0x6d, 0x49, 0x64, 0x22, 0x5d, 0x0a, 0x0d, 0x47, 0x6f, 0x6f,
	0x64, 0x73, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x12, 0x26, 0x0a, 0x0e, 0x43, 0x75,
	0x72, 0x72, 0x65, 0x6e, 0x74, 0x47, 0x6f, 0x6f, 0x64, 0x73, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x0e, 0x43, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x47, 0x6f, 0x6f, 0x64, 0x73,
	0x49, 0x64, 0x12, 0x24, 0x0a, 0x04, 0x44, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x6f, 0x6f, 0x64, 0x73, 0x49, 0x6e,
	0x66, 0x6f, 0x52, 0x04, 0x44, 0x61, 0x74, 0x61, 0x22, 0xc1, 0x01, 0x0a, 0x09, 0x47, 0x6f, 0x6f,
	0x64, 0x73, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x18, 0x0a, 0x07, 0x47, 0x6f, 0x6f, 0x64, 0x73, 0x49,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x47, 0x6f, 0x6f, 0x64, 0x73, 0x49, 0x64,
	0x12, 0x1e, 0x0a, 0x0a, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x49, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x49, 0x64,
	0x12, 0x16, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x54, 0x69, 0x74, 0x6c,
	0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x54, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x20,
	0x0a, 0x0b, 0x4d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x50, 0x72, 0x69, 0x63, 0x65, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0b, 0x4d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x50, 0x72, 0x69, 0x63, 0x65,
	0x12, 0x14, 0x0a, 0x05, 0x50, 0x72, 0x69, 0x63, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x50, 0x72, 0x69, 0x63, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x42, 0x72, 0x69, 0x65, 0x66, 0x18,
	0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x42, 0x72, 0x69, 0x65, 0x66, 0x22, 0x45, 0x0a, 0x11,
	0x47, 0x65, 0x74, 0x47, 0x6f, 0x6f, 0x64, 0x73, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x52, 0x65,
	0x71, 0x12, 0x18, 0x0a, 0x07, 0x47, 0x6f, 0x6f, 0x64, 0x73, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x07, 0x47, 0x6f, 0x6f, 0x64, 0x73, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x55,
	0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x55, 0x73, 0x65,
	0x72, 0x49, 0x64, 0x22, 0xf5, 0x01, 0x0a, 0x0b, 0x47, 0x6f, 0x6f, 0x64, 0x73, 0x44, 0x65, 0x74,
	0x61, 0x69, 0x6c, 0x12, 0x18, 0x0a, 0x07, 0x47, 0x6f, 0x6f, 0x64, 0x73, 0x49, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x47, 0x6f, 0x6f, 0x64, 0x73, 0x49, 0x64, 0x12, 0x1e, 0x0a,
	0x0a, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x0a, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x49, 0x64, 0x12, 0x16, 0x0a,
	0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x53,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x54, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x54, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x43,
	0x6f, 0x64, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x43, 0x6f, 0x64, 0x65, 0x12,
	0x1c, 0x0a, 0x09, 0x42, 0x72, 0x61, 0x6e, 0x64, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x42, 0x72, 0x61, 0x6e, 0x64, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a,
	0x0b, 0x4d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x50, 0x72, 0x69, 0x63, 0x65, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0b, 0x4d, 0x61, 0x72, 0x6b, 0x65, 0x74, 0x50, 0x72, 0x69, 0x63, 0x65, 0x12,
	0x14, 0x0a, 0x05, 0x50, 0x72, 0x69, 0x63, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x50, 0x72, 0x69, 0x63, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x42, 0x72, 0x69, 0x65, 0x66, 0x18, 0x09,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x42, 0x72, 0x69, 0x65, 0x66, 0x32, 0x89, 0x01, 0x0a, 0x05,
	0x47, 0x6f, 0x6f, 0x64, 0x73, 0x12, 0x40, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x47, 0x6f, 0x6f, 0x64,
	0x73, 0x42, 0x79, 0x52, 0x6f, 0x6f, 0x6d, 0x12, 0x18, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x47, 0x65, 0x74, 0x47, 0x6f, 0x6f, 0x64, 0x73, 0x42, 0x79, 0x52, 0x6f, 0x6f, 0x6d, 0x52, 0x65,
	0x71, 0x1a, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x6f, 0x6f, 0x64, 0x73, 0x4c,
	0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x12, 0x3e, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x47, 0x6f,
	0x6f, 0x64, 0x73, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x12, 0x18, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x47, 0x65, 0x74, 0x47, 0x6f, 0x6f, 0x64, 0x73, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c,
	0x52, 0x65, 0x71, 0x1a, 0x12, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x6f, 0x6f, 0x64,
	0x73, 0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x42, 0x09, 0x5a, 0x07, 0x2e, 0x3b, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_goods_proto_rawDescOnce sync.Once
	file_goods_proto_rawDescData []byte
)

func file_goods_proto_rawDescGZIP() []byte {
	file_goods_proto_rawDescOnce.Do(func() {
		file_goods_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_goods_proto_rawDesc), len(file_goods_proto_rawDesc)))
	})
	return file_goods_proto_rawDescData
}

var file_goods_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_goods_proto_goTypes = []any{
	(*GetGoodsByRoomReq)(nil), // 0: proto.GetGoodsByRoomReq
	(*GoodsListResp)(nil),     // 1: proto.GoodsListResp
	(*GoodsInfo)(nil),         // 2: proto.GoodsInfo
	(*GetGoodsDetailReq)(nil), // 3: proto.GetGoodsDetailReq
	(*GoodsDetail)(nil),       // 4: proto.GoodsDetail
}
var file_goods_proto_depIdxs = []int32{
	2, // 0: proto.GoodsListResp.Data:type_name -> proto.GoodsInfo
	0, // 1: proto.Goods.GetGoodsByRoom:input_type -> proto.GetGoodsByRoomReq
	3, // 2: proto.Goods.GetGoodsDetail:input_type -> proto.GetGoodsDetailReq
	1, // 3: proto.Goods.GetGoodsByRoom:output_type -> proto.GoodsListResp
	4, // 4: proto.Goods.GetGoodsDetail:output_type -> proto.GoodsDetail
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_goods_proto_init() }
func file_goods_proto_init() {
	if File_goods_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_goods_proto_rawDesc), len(file_goods_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_goods_proto_goTypes,
		DependencyIndexes: file_goods_proto_depIdxs,
		MessageInfos:      file_goods_proto_msgTypes,
	}.Build()
	File_goods_proto = out.File
	file_goods_proto_goTypes = nil
	file_goods_proto_depIdxs = nil
}
