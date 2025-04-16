// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v3.20.1
// source: stock.proto

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

// 定义通用响应消息
type Response struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Success       bool                   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"` // 操作是否成功
	Message       string                 `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`  // 操作结果描述
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Response) Reset() {
	*x = Response{}
	mi := &file_stock_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response) ProtoMessage() {}

func (x *Response) ProtoReflect() protoreflect.Message {
	mi := &file_stock_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response.ProtoReflect.Descriptor instead.
func (*Response) Descriptor() ([]byte, []int) {
	return file_stock_proto_rawDescGZIP(), []int{0}
}

func (x *Response) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *Response) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type GetStockReq struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	GoodsId       int64                  `protobuf:"varint,1,opt,name=goods_id,json=goodsId,proto3" json:"goods_id,omitempty"` // 商品ID
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetStockReq) Reset() {
	*x = GetStockReq{}
	mi := &file_stock_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetStockReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetStockReq) ProtoMessage() {}

func (x *GetStockReq) ProtoReflect() protoreflect.Message {
	mi := &file_stock_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetStockReq.ProtoReflect.Descriptor instead.
func (*GetStockReq) Descriptor() ([]byte, []int) {
	return file_stock_proto_rawDescGZIP(), []int{1}
}

func (x *GetStockReq) GetGoodsId() int64 {
	if x != nil {
		return x.GoodsId
	}
	return 0
}

type GoodsStockInfo struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	GoodsId       int64                  `protobuf:"varint,1,opt,name=goods_id,json=goodsId,proto3" json:"goods_id,omitempty"` // 商品ID
	Stock         int64                  `protobuf:"varint,2,opt,name=stock,proto3" json:"stock,omitempty"`                    // 库存数量
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GoodsStockInfo) Reset() {
	*x = GoodsStockInfo{}
	mi := &file_stock_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GoodsStockInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GoodsStockInfo) ProtoMessage() {}

func (x *GoodsStockInfo) ProtoReflect() protoreflect.Message {
	mi := &file_stock_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GoodsStockInfo.ProtoReflect.Descriptor instead.
func (*GoodsStockInfo) Descriptor() ([]byte, []int) {
	return file_stock_proto_rawDescGZIP(), []int{2}
}

func (x *GoodsStockInfo) GetGoodsId() int64 {
	if x != nil {
		return x.GoodsId
	}
	return 0
}

func (x *GoodsStockInfo) GetStock() int64 {
	if x != nil {
		return x.Stock
	}
	return 0
}

type ReduceStockInfo struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	GoodsId       int64                  `protobuf:"varint,1,opt,name=goods_id,json=goodsId,proto3" json:"goods_id,omitempty"` // 商品ID
	Num           int64                  `protobuf:"varint,2,opt,name=num,proto3" json:"num,omitempty"`                        // 下单商品数量（扣减库存数量）
	OrderId       int64                  `protobuf:"varint,3,opt,name=order_id,json=orderId,proto3" json:"order_id,omitempty"` // 订单ID（可选字段，仅在扣减库存时使用）
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ReduceStockInfo) Reset() {
	*x = ReduceStockInfo{}
	mi := &file_stock_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ReduceStockInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReduceStockInfo) ProtoMessage() {}

func (x *ReduceStockInfo) ProtoReflect() protoreflect.Message {
	mi := &file_stock_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReduceStockInfo.ProtoReflect.Descriptor instead.
func (*ReduceStockInfo) Descriptor() ([]byte, []int) {
	return file_stock_proto_rawDescGZIP(), []int{3}
}

func (x *ReduceStockInfo) GetGoodsId() int64 {
	if x != nil {
		return x.GoodsId
	}
	return 0
}

func (x *ReduceStockInfo) GetNum() int64 {
	if x != nil {
		return x.Num
	}
	return 0
}

func (x *ReduceStockInfo) GetOrderId() int64 {
	if x != nil {
		return x.OrderId
	}
	return 0
}

type StockInfoList struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Data          []*GoodsStockInfo      `protobuf:"bytes,1,rep,name=data,proto3" json:"data,omitempty"` // 库存信息列表
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StockInfoList) Reset() {
	*x = StockInfoList{}
	mi := &file_stock_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StockInfoList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StockInfoList) ProtoMessage() {}

func (x *StockInfoList) ProtoReflect() protoreflect.Message {
	mi := &file_stock_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StockInfoList.ProtoReflect.Descriptor instead.
func (*StockInfoList) Descriptor() ([]byte, []int) {
	return file_stock_proto_rawDescGZIP(), []int{4}
}

func (x *StockInfoList) GetData() []*GoodsStockInfo {
	if x != nil {
		return x.Data
	}
	return nil
}

var File_stock_proto protoreflect.FileDescriptor

var file_stock_proto_rawDesc = string([]byte{
	0x0a, 0x0b, 0x73, 0x74, 0x6f, 0x63, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x3e, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x22, 0x28, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x53, 0x74, 0x6f, 0x63, 0x6b,
	0x52, 0x65, 0x71, 0x12, 0x19, 0x0a, 0x08, 0x67, 0x6f, 0x6f, 0x64, 0x73, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x67, 0x6f, 0x6f, 0x64, 0x73, 0x49, 0x64, 0x22, 0x41,
	0x0a, 0x0e, 0x47, 0x6f, 0x6f, 0x64, 0x73, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x49, 0x6e, 0x66, 0x6f,
	0x12, 0x19, 0x0a, 0x08, 0x67, 0x6f, 0x6f, 0x64, 0x73, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x07, 0x67, 0x6f, 0x6f, 0x64, 0x73, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x73,
	0x74, 0x6f, 0x63, 0x6b, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x73, 0x74, 0x6f, 0x63,
	0x6b, 0x22, 0x59, 0x0a, 0x0f, 0x52, 0x65, 0x64, 0x75, 0x63, 0x65, 0x53, 0x74, 0x6f, 0x63, 0x6b,
	0x49, 0x6e, 0x66, 0x6f, 0x12, 0x19, 0x0a, 0x08, 0x67, 0x6f, 0x6f, 0x64, 0x73, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x67, 0x6f, 0x6f, 0x64, 0x73, 0x49, 0x64, 0x12,
	0x10, 0x0a, 0x03, 0x6e, 0x75, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x03, 0x6e, 0x75,
	0x6d, 0x12, 0x19, 0x0a, 0x08, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x07, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x49, 0x64, 0x22, 0x3a, 0x0a, 0x0d,
	0x53, 0x74, 0x6f, 0x63, 0x6b, 0x49, 0x6e, 0x66, 0x6f, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x29, 0x0a,
	0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x6f, 0x6f, 0x64, 0x73, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x49, 0x6e,
	0x66, 0x6f, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x32, 0xdc, 0x02, 0x0a, 0x05, 0x53, 0x74, 0x6f,
	0x63, 0x6b, 0x12, 0x32, 0x0a, 0x08, 0x53, 0x65, 0x74, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x12, 0x15,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x6f, 0x6f, 0x64, 0x73, 0x53, 0x74, 0x6f, 0x63,
	0x6b, 0x49, 0x6e, 0x66, 0x6f, 0x1a, 0x0f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x35, 0x0a, 0x08, 0x47, 0x65, 0x74, 0x53, 0x74, 0x6f,
	0x63, 0x6b, 0x12, 0x12, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x53, 0x74,
	0x6f, 0x63, 0x6b, 0x52, 0x65, 0x71, 0x1a, 0x15, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47,
	0x6f, 0x6f, 0x64, 0x73, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x36, 0x0a,
	0x0b, 0x52, 0x65, 0x64, 0x75, 0x63, 0x65, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x12, 0x16, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x52, 0x65, 0x64, 0x75, 0x63, 0x65, 0x53, 0x74, 0x6f, 0x63, 0x6b,
	0x49, 0x6e, 0x66, 0x6f, 0x1a, 0x0f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x38, 0x0a, 0x0d, 0x52, 0x6f, 0x6c, 0x6c, 0x62, 0x61, 0x63,
	0x6b, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x12, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x52,
	0x65, 0x64, 0x75, 0x63, 0x65, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x49, 0x6e, 0x66, 0x6f, 0x1a, 0x0f,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x3b, 0x0a, 0x0d, 0x42, 0x61, 0x74, 0x63, 0x68, 0x47, 0x65, 0x74, 0x53, 0x74, 0x6f, 0x63, 0x6b,
	0x12, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x49, 0x6e,
	0x66, 0x6f, 0x4c, 0x69, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53,
	0x74, 0x6f, 0x63, 0x6b, 0x49, 0x6e, 0x66, 0x6f, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x39, 0x0a, 0x10,
	0x42, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x64, 0x75, 0x63, 0x65, 0x53, 0x74, 0x6f, 0x63, 0x6b,
	0x12, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x49, 0x6e,
	0x66, 0x6f, 0x4c, 0x69, 0x73, 0x74, 0x1a, 0x0f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x09, 0x5a, 0x07, 0x2e, 0x3b, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_stock_proto_rawDescOnce sync.Once
	file_stock_proto_rawDescData []byte
)

func file_stock_proto_rawDescGZIP() []byte {
	file_stock_proto_rawDescOnce.Do(func() {
		file_stock_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_stock_proto_rawDesc), len(file_stock_proto_rawDesc)))
	})
	return file_stock_proto_rawDescData
}

var file_stock_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_stock_proto_goTypes = []any{
	(*Response)(nil),        // 0: proto.Response
	(*GetStockReq)(nil),     // 1: proto.GetStockReq
	(*GoodsStockInfo)(nil),  // 2: proto.GoodsStockInfo
	(*ReduceStockInfo)(nil), // 3: proto.ReduceStockInfo
	(*StockInfoList)(nil),   // 4: proto.StockInfoList
}
var file_stock_proto_depIdxs = []int32{
	2, // 0: proto.StockInfoList.data:type_name -> proto.GoodsStockInfo
	2, // 1: proto.Stock.SetStock:input_type -> proto.GoodsStockInfo
	1, // 2: proto.Stock.GetStock:input_type -> proto.GetStockReq
	3, // 3: proto.Stock.ReduceStock:input_type -> proto.ReduceStockInfo
	3, // 4: proto.Stock.RollbackStock:input_type -> proto.ReduceStockInfo
	4, // 5: proto.Stock.BatchGetStock:input_type -> proto.StockInfoList
	4, // 6: proto.Stock.BatchReduceStock:input_type -> proto.StockInfoList
	0, // 7: proto.Stock.SetStock:output_type -> proto.Response
	2, // 8: proto.Stock.GetStock:output_type -> proto.GoodsStockInfo
	0, // 9: proto.Stock.ReduceStock:output_type -> proto.Response
	0, // 10: proto.Stock.RollbackStock:output_type -> proto.Response
	4, // 11: proto.Stock.BatchGetStock:output_type -> proto.StockInfoList
	0, // 12: proto.Stock.BatchReduceStock:output_type -> proto.Response
	7, // [7:13] is the sub-list for method output_type
	1, // [1:7] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_stock_proto_init() }
func file_stock_proto_init() {
	if File_stock_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_stock_proto_rawDesc), len(file_stock_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_stock_proto_goTypes,
		DependencyIndexes: file_stock_proto_depIdxs,
		MessageInfos:      file_stock_proto_msgTypes,
	}.Build()
	File_stock_proto = out.File
	file_stock_proto_goTypes = nil
	file_stock_proto_depIdxs = nil
}
