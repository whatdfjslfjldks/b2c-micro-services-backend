package handler

import (
	"context"
	logServerProto "micro-services/pkg/proto/log-server"
	pb "micro-services/pkg/proto/product-server"
	"micro-services/pkg/utils"
	"micro-services/product-server/internal/repository"
	"micro-services/product-server/pkg/instance"
)

func (s *Server) GetProductById(ctx context.Context, req *pb.GetProductByIdRequest) (
	*pb.GetProductByIdResponse, error) {
	product, err := repository.GetProductById(req.ProductId)
	if err != nil {
		a := &logServerProto.PostLogRequest{
			Level:       "ERROR",
			Msg:         err.Error(),
			RequestPath: "/getProductById",
			Source:      "product-server",
			StatusCode:  "GLB-003",
			Time:        utils.GetTime(),
		}
		instance.GrpcClient.PostLog(a)

		return nil, err
	}

	//int32 product_id = 1;
	//string product_name = 2;
	//string product_cover = 3;
	//double product_price = 4;
	//int32 product_categoryId = 5;
	//string description = 6;

	return &pb.GetProductByIdResponse{
		ProductId:         req.ProductId,
		ProductName:       product.Name,
		ProductCover:      product.Cover,
		ProductPrice:      product.Price,
		ProductCategoryId: product.Category,
		Description:       product.Description,
	}, nil
}
