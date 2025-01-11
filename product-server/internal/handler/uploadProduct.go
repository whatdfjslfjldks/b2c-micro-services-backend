package handler

import (
	"context"
	logServerProto "micro-services/pkg/proto/log-server"
	pb "micro-services/pkg/proto/product-server"
	"micro-services/pkg/utils"
	"micro-services/product-server/internal/service/file"
	"micro-services/product-server/pkg/instance"
)

// UploadProductByExcel TODO 注意，不依靠返回的error判断
// 只读取第一个 sheet 文件
func (s *Server) UploadProductByExcel(ctx context.Context, req *pb.UploadProductByExcelRequest) (
	*pb.UploadProductByExcelResponse, error) {
	//log.Println("接受到的文件比特流： ", req.File)
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
