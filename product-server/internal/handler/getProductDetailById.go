package handler

// GetProductDetailById 获取详情页商品信息
//func (s *Server) GetProductDetailById(ctx context.Context, req *pb.GetProductDetailByIdRequest) (
//	*pb.GetProductDetailByIdResponse, error) {
//product, err := repository.GetProductDetailById(req.ProductId)
//if err != nil {
//	a := &logServerProto.PostLogRequest{
//		Level:       "ERROR",
//		Msg:         err.Error(),
//		RequestPath: "/getProductById",
//		Source:      "product-server",
//		StatusCode:  "GLB-003",
//		Time:        utils.GetTime(),
//	}
//	instance.GrpcClient.PostLog(a)
//	return nil, err
//}
//
////log.Printf("product: %v", product)
//// TODO sold is not defined
//return &pb.GetProductDetailByIdResponse{
//	Code:         200,
//	StatusCode:   "GLB-000",
//	Msg:          "商品信息获取成功！",
//	ProductId:    req.ProductId,
//	ProductName:  product.ProductName,
//	ProductImg:   product.ProductImg,
//	ProductPrice: product.ProductPrice,
//	ProductType:  product.ProductType,
//	Sold:         product.Sold,
//}, nil

//	return nil, nil
//}
