package handler

import (
	"context"
	logServerProto "micro-services/pkg/proto/log-server"
	pb "micro-services/pkg/proto/product-server"
	"micro-services/pkg/utils"
	"micro-services/product-server/internal/repository"
	"micro-services/product-server/internal/service"
	"micro-services/product-server/internal/service/file"
	"micro-services/product-server/pkg/instance"
	"micro-services/product-server/pkg/model/dto"
)

// UploadProductByExcel TODO 注意，不依靠返回的error判断
// 只读取第一个 sheet 文件
func (s *Server) UploadProductByExcel(ctx context.Context, req *pb.UploadProductByExcelRequest) (
	*pb.UploadProductByExcelResponse, error) {
	// 首先检查文件格式是否符合要求
	b, e := file.IsFileValid(req.File)
	if e != nil || !b {
		// 也可能是服务端问题。。
		return &pb.UploadProductByExcelResponse{
			Code:       400,
			Msg:        e.Error(),
			StatusCode: "PRT-002",
		}, nil
	}
	// 符合要求，开始读取并存储数据
	err := file.UploadProduct(req.File)
	if err != nil {
		if err.Error() == "GLB-003" {
			a := &logServerProto.PostLogRequest{
				Level:       "ERROR",
				Msg:         "数据库错误！",
				RequestPath: "/uploadProductByExcel",
				Source:      "product-server",
				StatusCode:  "GLB-003",
				Time:        utils.GetTime(),
			}
			instance.GrpcClient.PostLog(a)
			return &pb.UploadProductByExcelResponse{
				Code:       500,
				Msg:        "数据库错误！",
				StatusCode: "GLB-003",
			}, nil
		}
		return &pb.UploadProductByExcelResponse{
			Code:       500,
			Msg:        err.Error(),
			StatusCode: "PRT-001",
		}, nil
	}

	return nil, nil
}

// UploadSecKillProduct 上传秒杀商品
func (s *Server) UploadSecKillProduct(c context.Context, req *pb.UploadSecKillProductRequest) (
	*pb.UploadSecKillProductResponse, error) {
	// 判断场次是否存在
	if !service.IsSessionValid(req.Time) {
		resp := &pb.UploadSecKillProductResponse{
			Code:       400,
			StatusCode: "GLB-001",
			Msg:        "场次不存在！",
		}
		return resp, nil
	}
	// 判断图片,类别列表是否为空
	if len(req.SecImg) == 0 || len(req.SecType) == 0 {
		resp := &pb.UploadSecKillProductResponse{
			Code:       400,
			StatusCode: "GLB-001",
			Msg:        "图片或类别列表为空！",
		}
		return resp, nil
	}
	secProduct := dto.SecKillProduct{
		SecName:          req.SecName,
		SecDescription:   req.SecDescription,
		SecPrice:         req.SecPrice,
		SecOriginalPrice: req.SecOriginalPrice,
		Stock:            req.SecStock,
		StartTime:        req.SecStartTime,
		EndTime:          req.SecEndTime,
		SecType:          req.SecType,
		SecImg:           req.SecImg,
		Time:             req.Time,
	}
	err := repository.UploadSecProduct(secProduct)
	if err != nil {
		if err.Error() == "GLB-003" {
			a := &logServerProto.PostLogRequest{
				Level:       "ERROR",
				Msg:         "数据库错误！",
				RequestPath: "/uploadSecKillProduct",
				Source:      "product-server",
				StatusCode:  "GLB-003",
				Time:        utils.GetTime(),
			}
			instance.GrpcClient.PostLog(a)
			return &pb.UploadSecKillProductResponse{
				Code:       500,
				Msg:        "数据库错误！",
				StatusCode: "GLB-003",
			}, nil
		}
		return &pb.UploadSecKillProductResponse{
			Code:       500,
			Msg:        "数据插入失败！",
			StatusCode: err.Error(),
		}, nil
	}
	return &pb.UploadSecKillProductResponse{
		Code:       200,
		Msg:        "数据插入成功！",
		StatusCode: "GLB-000",
	}, nil
}
