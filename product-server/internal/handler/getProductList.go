package handler

import (
	"context"
	logServerProto "micro-services/pkg/proto/log-server"
	pb "micro-services/pkg/proto/product-server"
	"micro-services/pkg/utils"
	"micro-services/product-server/internal/repository"
	"micro-services/product-server/internal/service"
	"micro-services/product-server/pkg/instance"
)

// GetProductList 筛选条件 价格升序，时间降序
// TODO 有空把repository层的service部分分出来
func (s *Server) GetProductList(ctx context.Context, req *pb.GetProductListRequest) (
	*pb.GetProductListResponse, error) {
	resp := &pb.GetProductListResponse{}

	list, totalItems, err := repository.GetProductList(req.CurrentPage, req.PageSize, req.CategoryId, req.Sort)
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
	resp.CategoryId = req.CategoryId
	return resp, nil
}

// GetSecKillList 获取秒杀商品列表
// TODO 没有细分错误处理
func (s *Server) GetSecKillList(ctx context.Context, req *pb.GetSecKillListRequest) (
	*pb.GetSecKillListResponse, error) {
	// 判断场次是存在
	//fmt.Println("Sdfdsfsdf: ", req.Time)
	if !service.IsSessionValid(req.Time) {
		resp := &pb.GetSecKillListResponse{
			Code:        400,
			StatusCode:  "GLB-001",
			Msg:         "场次不存在！",
			CurrentPage: req.CurrentPage,
			PageSize:    req.PageSize,
			Time:        req.Time,
		}
		return resp, nil
	}
	secProducts, totalItems, err := service.GetSecListAndTotalItems(req.CurrentPage, req.PageSize, req.Time)
	if err != nil {
		a := &logServerProto.PostLogRequest{
			Level:       "ERROR",
			Msg:         err.Error(),
			RequestPath: "/getSecKillList",
			Source:      "product-server",
			StatusCode:  "GLB-003",
			Time:        utils.GetTime(),
		}
		instance.GrpcClient.PostLog(a)
		resp := &pb.GetSecKillListResponse{
			Code:        500,
			StatusCode:  "GLB-003",
			Msg:         "获取数据失败！",
			CurrentPage: req.CurrentPage,
			PageSize:    req.PageSize,
			Time:        req.Time,
		}
		return resp, nil
	}

	return &pb.GetSecKillListResponse{
		Code:        200,
		StatusCode:  "GBL-000",
		Msg:         "获取数据成功！",
		CurrentPage: req.CurrentPage,
		PageSize:    req.PageSize,
		TotalItems:  totalItems,
		SecList:     secProducts,
		Time:        req.Time,
	}, nil
}
