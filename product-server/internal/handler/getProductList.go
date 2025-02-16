package handler

import (
	"context"
	"database/sql"
	"errors"
	"log"
	logServerProto "micro-services/pkg/proto/log-server"
	pb "micro-services/pkg/proto/product-server"
	"micro-services/pkg/utils"
	"micro-services/product-server/internal/repository"
	"micro-services/product-server/internal/service"
	"micro-services/product-server/pkg/instance"
)

// GetProductList 筛选条件 价格升序，时间降序
func (s *Server) GetProductList(ctx context.Context, req *pb.GetProductListRequest) (
	*pb.GetProductListResponse, error) {
	resp := &pb.GetProductListResponse{}

	productList, totalItems, err := repository.GetProductList(req.CurrentPage, req.PageSize, req.CategoryId, 1, req.Sort, req.Keyword)
	if err != nil {
		log.Printf("GetProductList error: %v", err)
		if errors.Is(err, sql.ErrNoRows) {
			resp = &pb.GetProductListResponse{
				Code:        400,
				StatusCode:  "GLB-004",
				Msg:         "没有商品！",
				CurrentPage: req.CurrentPage,
				PageSize:    req.PageSize,
			}
			return resp, nil
		} else {
			a := &logServerProto.PostLogRequest{
				Level:       "ERROR",
				Msg:         err.Error(),
				RequestPath: "/getProductList",
				Source:      "product-server",
				StatusCode:  "GLB-003",
				Time:        utils.GetTime(),
			}
			instance.GrpcClient.PostLog(a)
			resp := &pb.GetProductListResponse{
				Code:        500,
				StatusCode:  "GLB-003",
				Msg:         "获取数据失败！",
				CurrentPage: req.CurrentPage,
				PageSize:    req.PageSize,
			}
			return resp, nil
		}
	}
	resp = &pb.GetProductListResponse{
		Code:        200,
		StatusCode:  "GLB-000",
		Msg:         "获取数据成功！",
		CurrentPage: req.CurrentPage,
		PageSize:    req.PageSize,
		TotalItems:  totalItems,
		ProductList: productList,
	}
	return resp, nil
}

// GetSecKillList 获取秒杀商品列表
func (s *Server) GetSecKillList(ctx context.Context, req *pb.GetSecKillListRequest) (
	*pb.GetSecKillListResponse, error) {
	// 判断场次是存在
	if !service.IsSessionValid(req.SessionId) {
		resp := &pb.GetSecKillListResponse{
			Code:        400,
			StatusCode:  "GLB-001",
			Msg:         "场次不存在！",
			CurrentPage: req.CurrentPage,
			PageSize:    req.PageSize,
			SessionId:   req.SessionId,
		}
		return resp, nil
	}
	secProducts, totalItems, err := repository.GetSecList(req.CurrentPage, req.PageSize, req.SessionId)
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
			SessionId:   req.SessionId,
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
		SessionId:   req.SessionId,
	}, nil
}
