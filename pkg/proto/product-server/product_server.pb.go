// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.1
// 	protoc        v5.29.2
// source: product_server.proto

package product_server_proto

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

type GetProductListRequest struct {
	state       protoimpl.MessageState `protogen:"open.v1"`
	CurrentPage int32                  `protobuf:"varint,1,opt,name=currentPage,proto3" json:"currentPage,omitempty"`
	PageSize    int32                  `protobuf:"varint,2,opt,name=pageSize,proto3" json:"pageSize,omitempty"`
	// 前后端维护一个种类id对应表
	CategoryId int32 `protobuf:"varint,3,opt,name=categoryId,proto3" json:"categoryId,omitempty"`
	// 前后端维护一个价格区间对应表
	PriceRange    int32 `protobuf:"varint,4,opt,name=priceRange,proto3" json:"priceRange,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetProductListRequest) Reset() {
	*x = GetProductListRequest{}
	mi := &file_product_server_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetProductListRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetProductListRequest) ProtoMessage() {}

func (x *GetProductListRequest) ProtoReflect() protoreflect.Message {
	mi := &file_product_server_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetProductListRequest.ProtoReflect.Descriptor instead.
func (*GetProductListRequest) Descriptor() ([]byte, []int) {
	return file_product_server_proto_rawDescGZIP(), []int{0}
}

func (x *GetProductListRequest) GetCurrentPage() int32 {
	if x != nil {
		return x.CurrentPage
	}
	return 0
}

func (x *GetProductListRequest) GetPageSize() int32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

func (x *GetProductListRequest) GetCategoryId() int32 {
	if x != nil {
		return x.CategoryId
	}
	return 0
}

func (x *GetProductListRequest) GetPriceRange() int32 {
	if x != nil {
		return x.PriceRange
	}
	return 0
}

type GetProductListResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Code          int32                  `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	StatusCode    string                 `protobuf:"bytes,2,opt,name=status_code,json=statusCode,proto3" json:"status_code,omitempty"`
	Msg           string                 `protobuf:"bytes,3,opt,name=msg,proto3" json:"msg,omitempty"`
	CurrentPage   int32                  `protobuf:"varint,4,opt,name=currentPage,proto3" json:"currentPage,omitempty"`
	PageSize      int32                  `protobuf:"varint,5,opt,name=pageSize,proto3" json:"pageSize,omitempty"`
	TotalItems    int32                  `protobuf:"varint,6,opt,name=totalItems,proto3" json:"totalItems,omitempty"`
	ProductList   []*ProductListItem     `protobuf:"bytes,7,rep,name=productList,proto3" json:"productList,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetProductListResponse) Reset() {
	*x = GetProductListResponse{}
	mi := &file_product_server_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetProductListResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetProductListResponse) ProtoMessage() {}

func (x *GetProductListResponse) ProtoReflect() protoreflect.Message {
	mi := &file_product_server_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetProductListResponse.ProtoReflect.Descriptor instead.
func (*GetProductListResponse) Descriptor() ([]byte, []int) {
	return file_product_server_proto_rawDescGZIP(), []int{1}
}

func (x *GetProductListResponse) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *GetProductListResponse) GetStatusCode() string {
	if x != nil {
		return x.StatusCode
	}
	return ""
}

func (x *GetProductListResponse) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

func (x *GetProductListResponse) GetCurrentPage() int32 {
	if x != nil {
		return x.CurrentPage
	}
	return 0
}

func (x *GetProductListResponse) GetPageSize() int32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

func (x *GetProductListResponse) GetTotalItems() int32 {
	if x != nil {
		return x.TotalItems
	}
	return 0
}

func (x *GetProductListResponse) GetProductList() []*ProductListItem {
	if x != nil {
		return x.ProductList
	}
	return nil
}

type ProductListItem struct {
	state             protoimpl.MessageState `protogen:"open.v1"`
	ProductId         int32                  `protobuf:"varint,1,opt,name=product_id,json=productId,proto3" json:"product_id,omitempty"`
	ProductName       string                 `protobuf:"bytes,2,opt,name=product_name,json=productName,proto3" json:"product_name,omitempty"`
	ProductCover      string                 `protobuf:"bytes,3,opt,name=product_cover,json=productCover,proto3" json:"product_cover,omitempty"`
	ProductPrice      string                 `protobuf:"bytes,4,opt,name=product_price,json=productPrice,proto3" json:"product_price,omitempty"`
	ProductCategoryId int32                  `protobuf:"varint,5,opt,name=product_categoryId,json=productCategoryId,proto3" json:"product_categoryId,omitempty"`
	unknownFields     protoimpl.UnknownFields
	sizeCache         protoimpl.SizeCache
}

func (x *ProductListItem) Reset() {
	*x = ProductListItem{}
	mi := &file_product_server_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ProductListItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProductListItem) ProtoMessage() {}

func (x *ProductListItem) ProtoReflect() protoreflect.Message {
	mi := &file_product_server_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProductListItem.ProtoReflect.Descriptor instead.
func (*ProductListItem) Descriptor() ([]byte, []int) {
	return file_product_server_proto_rawDescGZIP(), []int{2}
}

func (x *ProductListItem) GetProductId() int32 {
	if x != nil {
		return x.ProductId
	}
	return 0
}

func (x *ProductListItem) GetProductName() string {
	if x != nil {
		return x.ProductName
	}
	return ""
}

func (x *ProductListItem) GetProductCover() string {
	if x != nil {
		return x.ProductCover
	}
	return ""
}

func (x *ProductListItem) GetProductPrice() string {
	if x != nil {
		return x.ProductPrice
	}
	return ""
}

func (x *ProductListItem) GetProductCategoryId() int32 {
	if x != nil {
		return x.ProductCategoryId
	}
	return 0
}

var File_product_server_proto protoreflect.FileDescriptor

var file_product_server_proto_rawDesc = []byte{
	0x0a, 0x14, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x95, 0x01,
	0x0a, 0x15, 0x47, 0x65, 0x74, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x4c, 0x69, 0x73, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x20, 0x0a, 0x0b, 0x63, 0x75, 0x72, 0x72, 0x65,
	0x6e, 0x74, 0x50, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0b, 0x63, 0x75,
	0x72, 0x72, 0x65, 0x6e, 0x74, 0x50, 0x61, 0x67, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61, 0x67,
	0x65, 0x53, 0x69, 0x7a, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x70, 0x61, 0x67,
	0x65, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72,
	0x79, 0x49, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x63, 0x61, 0x74, 0x65, 0x67,
	0x6f, 0x72, 0x79, 0x49, 0x64, 0x12, 0x1e, 0x0a, 0x0a, 0x70, 0x72, 0x69, 0x63, 0x65, 0x52, 0x61,
	0x6e, 0x67, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x70, 0x72, 0x69, 0x63, 0x65,
	0x52, 0x61, 0x6e, 0x67, 0x65, 0x22, 0xf7, 0x01, 0x0a, 0x16, 0x47, 0x65, 0x74, 0x50, 0x72, 0x6f,
	0x64, 0x75, 0x63, 0x74, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04,
	0x63, 0x6f, 0x64, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x5f, 0x63,
	0x6f, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x03, 0x6d, 0x73, 0x67, 0x12, 0x20, 0x0a, 0x0b, 0x63, 0x75, 0x72, 0x72, 0x65,
	0x6e, 0x74, 0x50, 0x61, 0x67, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0b, 0x63, 0x75,
	0x72, 0x72, 0x65, 0x6e, 0x74, 0x50, 0x61, 0x67, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61, 0x67,
	0x65, 0x53, 0x69, 0x7a, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x70, 0x61, 0x67,
	0x65, 0x53, 0x69, 0x7a, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x49, 0x74,
	0x65, 0x6d, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x74, 0x6f, 0x74, 0x61, 0x6c,
	0x49, 0x74, 0x65, 0x6d, 0x73, 0x12, 0x38, 0x0a, 0x0b, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74,
	0x4c, 0x69, 0x73, 0x74, 0x18, 0x07, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x4c, 0x69, 0x73, 0x74, 0x49, 0x74,
	0x65, 0x6d, 0x52, 0x0b, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x4c, 0x69, 0x73, 0x74, 0x22,
	0xcc, 0x01, 0x0a, 0x0f, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x4c, 0x69, 0x73, 0x74, 0x49,
	0x74, 0x65, 0x6d, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74,
	0x49, 0x64, 0x12, 0x21, 0x0a, 0x0c, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x5f, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63,
	0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74,
	0x5f, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x70, 0x72,
	0x6f, 0x64, 0x75, 0x63, 0x74, 0x43, 0x6f, 0x76, 0x65, 0x72, 0x12, 0x23, 0x0a, 0x0d, 0x70, 0x72,
	0x6f, 0x64, 0x75, 0x63, 0x74, 0x5f, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0c, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x50, 0x72, 0x69, 0x63, 0x65, 0x12,
	0x2d, 0x0a, 0x12, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x5f, 0x63, 0x61, 0x74, 0x65, 0x67,
	0x6f, 0x72, 0x79, 0x49, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x11, 0x70, 0x72, 0x6f,
	0x64, 0x75, 0x63, 0x74, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x49, 0x64, 0x32, 0x5f,
	0x0a, 0x0e, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x4d, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x50, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x4c, 0x69,
	0x73, 0x74, 0x12, 0x1c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x50, 0x72,
	0x6f, 0x64, 0x75, 0x63, 0x74, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x1d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x47, 0x65, 0x74, 0x50, 0x72, 0x6f, 0x64,
	0x75, 0x63, 0x74, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42,
	0x27, 0x5a, 0x25, 0x2e, 0x2f, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x2d, 0x73, 0x65, 0x72,
	0x76, 0x65, 0x72, 0x3b, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x5f, 0x73, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x5f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_product_server_proto_rawDescOnce sync.Once
	file_product_server_proto_rawDescData = file_product_server_proto_rawDesc
)

func file_product_server_proto_rawDescGZIP() []byte {
	file_product_server_proto_rawDescOnce.Do(func() {
		file_product_server_proto_rawDescData = protoimpl.X.CompressGZIP(file_product_server_proto_rawDescData)
	})
	return file_product_server_proto_rawDescData
}

var file_product_server_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_product_server_proto_goTypes = []any{
	(*GetProductListRequest)(nil),  // 0: proto.GetProductListRequest
	(*GetProductListResponse)(nil), // 1: proto.GetProductListResponse
	(*ProductListItem)(nil),        // 2: proto.ProductListItem
}
var file_product_server_proto_depIdxs = []int32{
	2, // 0: proto.GetProductListResponse.productList:type_name -> proto.ProductListItem
	0, // 1: proto.ProductService.GetProductList:input_type -> proto.GetProductListRequest
	1, // 2: proto.ProductService.GetProductList:output_type -> proto.GetProductListResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_product_server_proto_init() }
func file_product_server_proto_init() {
	if File_product_server_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_product_server_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_product_server_proto_goTypes,
		DependencyIndexes: file_product_server_proto_depIdxs,
		MessageInfos:      file_product_server_proto_msgTypes,
	}.Build()
	File_product_server_proto = out.File
	file_product_server_proto_rawDesc = nil
	file_product_server_proto_goTypes = nil
	file_product_server_proto_depIdxs = nil
}
