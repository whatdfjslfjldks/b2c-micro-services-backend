package handler

import (
	"context"
	logServerProto "micro-services/pkg/proto/log-server"
	pb "micro-services/pkg/proto/product-server"
	"micro-services/pkg/utils"
	"micro-services/product-server/internal/service"
	"micro-services/product-server/pkg/instance"
)

func (s *Server) GetProductList(ctx context.Context, req *pb.GetProductListRequest) (
	*pb.GetProductListResponse, error) {
	resp := &pb.GetProductListResponse{}
	//fmt.Println("sdfds", req.CurrentPage, req.PageSize, req.CategoryId, req.PriceRange)
	list, totalItems, err := service.GetProductList(req.CurrentPage, req.PageSize, req.CategoryId, req.PriceRange)
	if err != nil {
		a := &logServerProto.PostLogRequest{
			Level:       "ERROR",
			Msg:         err.Error(),
			RequestPath: "/getProductList",
			Source:      "product-server",
			StatusCode:  "GLB-003",
			Time:        utils.GetTime(),
		}
		instance.GrpcClient.PostLog(a)
		resp.Code = 500
		resp.StatusCode = "GLB-003"
		resp.Msg = err.Error()
		return resp, nil
	}
	resp.ProductList = list
	resp.TotalItems = totalItems
	resp.Code = 200
	resp.StatusCode = "GLB-000"
	resp.Msg = "获取商品列表成功！"
	resp.CurrentPage = req.CurrentPage
	resp.PageSize = req.PageSize
	return resp, nil
}
