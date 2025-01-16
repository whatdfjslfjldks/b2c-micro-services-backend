// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.2
// source: product_server.proto

package product_server_proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	ProductService_GetProductList_FullMethodName       = "/proto.ProductService/GetProductList"
	ProductService_UploadProductByExcel_FullMethodName = "/proto.ProductService/UploadProductByExcel"
	ProductService_GetProductById_FullMethodName       = "/proto.ProductService/GetProductById"
)

// ProductServiceClient is the client API for ProductService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ProductServiceClient interface {
	// TODO: test 获取商品列表
	GetProductList(ctx context.Context, in *GetProductListRequest, opts ...grpc.CallOption) (*GetProductListResponse, error)
	// 注意限制文件大小，不能超过5MB，以免影响通信速度，估算，excel不能超过500行（实际可以1500行左右吧）
	// TODO 身份校验
	UploadProductByExcel(ctx context.Context, in *UploadProductByExcelRequest, opts ...grpc.CallOption) (*UploadProductByExcelResponse, error)
	// 通过product_id获取商品
	GetProductById(ctx context.Context, in *GetProductByIdRequest, opts ...grpc.CallOption) (*GetProductByIdResponse, error)
}

type productServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewProductServiceClient(cc grpc.ClientConnInterface) ProductServiceClient {
	return &productServiceClient{cc}
}

func (c *productServiceClient) GetProductList(ctx context.Context, in *GetProductListRequest, opts ...grpc.CallOption) (*GetProductListResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetProductListResponse)
	err := c.cc.Invoke(ctx, ProductService_GetProductList_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *productServiceClient) UploadProductByExcel(ctx context.Context, in *UploadProductByExcelRequest, opts ...grpc.CallOption) (*UploadProductByExcelResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UploadProductByExcelResponse)
	err := c.cc.Invoke(ctx, ProductService_UploadProductByExcel_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *productServiceClient) GetProductById(ctx context.Context, in *GetProductByIdRequest, opts ...grpc.CallOption) (*GetProductByIdResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetProductByIdResponse)
	err := c.cc.Invoke(ctx, ProductService_GetProductById_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProductServiceServer is the server API for ProductService service.
// All implementations must embed UnimplementedProductServiceServer
// for forward compatibility.
type ProductServiceServer interface {
	// TODO: test 获取商品列表
	GetProductList(context.Context, *GetProductListRequest) (*GetProductListResponse, error)
	// 注意限制文件大小，不能超过5MB，以免影响通信速度，估算，excel不能超过500行（实际可以1500行左右吧）
	// TODO 身份校验
	UploadProductByExcel(context.Context, *UploadProductByExcelRequest) (*UploadProductByExcelResponse, error)
	// 通过product_id获取商品
	GetProductById(context.Context, *GetProductByIdRequest) (*GetProductByIdResponse, error)
	mustEmbedUnimplementedProductServiceServer()
}

// UnimplementedProductServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedProductServiceServer struct{}

func (UnimplementedProductServiceServer) GetProductList(context.Context, *GetProductListRequest) (*GetProductListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProductList not implemented")
}
func (UnimplementedProductServiceServer) UploadProductByExcel(context.Context, *UploadProductByExcelRequest) (*UploadProductByExcelResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UploadProductByExcel not implemented")
}
func (UnimplementedProductServiceServer) GetProductById(context.Context, *GetProductByIdRequest) (*GetProductByIdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetProductById not implemented")
}
func (UnimplementedProductServiceServer) mustEmbedUnimplementedProductServiceServer() {}
func (UnimplementedProductServiceServer) testEmbeddedByValue()                        {}

// UnsafeProductServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ProductServiceServer will
// result in compilation errors.
type UnsafeProductServiceServer interface {
	mustEmbedUnimplementedProductServiceServer()
}

func RegisterProductServiceServer(s grpc.ServiceRegistrar, srv ProductServiceServer) {
	// If the following call pancis, it indicates UnimplementedProductServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&ProductService_ServiceDesc, srv)
}

func _ProductService_GetProductList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetProductListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProductServiceServer).GetProductList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProductService_GetProductList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProductServiceServer).GetProductList(ctx, req.(*GetProductListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProductService_UploadProductByExcel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UploadProductByExcelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProductServiceServer).UploadProductByExcel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProductService_UploadProductByExcel_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProductServiceServer).UploadProductByExcel(ctx, req.(*UploadProductByExcelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProductService_GetProductById_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetProductByIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProductServiceServer).GetProductById(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProductService_GetProductById_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProductServiceServer).GetProductById(ctx, req.(*GetProductByIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ProductService_ServiceDesc is the grpc.ServiceDesc for ProductService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ProductService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.ProductService",
	HandlerType: (*ProductServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetProductList",
			Handler:    _ProductService_GetProductList_Handler,
		},
		{
			MethodName: "UploadProductByExcel",
			Handler:    _ProductService_UploadProductByExcel_Handler,
		},
		{
			MethodName: "GetProductById",
			Handler:    _ProductService_GetProductById_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "product_server.proto",
}
