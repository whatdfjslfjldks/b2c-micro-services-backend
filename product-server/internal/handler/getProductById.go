package handler

import (
	"context"
	logServerProto "micro-services/pkg/proto/log-server"
	pb "micro-services/pkg/proto/product-server"
	"micro-services/pkg/utils"
	"micro-services/product-server/internal/service"
	"micro-services/product-server/pkg/instance"
)

// GetProductById 通过id获取详情页商品信息，普通，秒杀，预售都用这个函数
func (s *Server) GetProductById(ctx context.Context, req *pb.GetProductByIdRequest) (
	*pb.GetProductByIdResponse, error) {

	product, err := service.GetProductById(req.ProductId)
	if err != nil {
		//log.Printf("failed to get product by id: %v", err)
		if err.Error() == "GLB-001" {
			return &pb.GetProductByIdResponse{
				Code:       404,
				StatusCode: "GLB-001",
				Msg:        "商品不存在",
			}, nil
		}
		a := &logServerProto.PostLogRequest{
			Level:       "ERROR",
			Msg:         err.Error(),
			RequestPath: "/getProductById",
			Source:      "product-server",
			StatusCode:  "GLB-003",
			Time:        utils.GetTime(),
		}
		instance.GrpcClient.PostLog(a)
		return &pb.GetProductByIdResponse{
			Code:       500,
			StatusCode: "GLB-003",
			Msg:        "查询数据库失败！",
		}, nil
	}

	//fmt.Println("dsfsdf:", product)

	return &pb.GetProductByIdResponse{
		Code:       200,
		StatusCode: "GLB-000",
		Msg:        "获取数据成功！",
		Product:    product,
	}, nil

}
