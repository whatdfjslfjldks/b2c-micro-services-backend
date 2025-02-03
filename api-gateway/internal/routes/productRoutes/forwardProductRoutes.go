package productRoutes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
	"micro-services/api-gateway/internal/instance"
	"micro-services/pkg/utils"
	"path/filepath"
	"strings"

	productServerProto "micro-services/pkg/proto/product-server"
	"strconv"
)

var model2 = "product-server"

// GetProductList 获取商品列表
func GetProductList(c *gin.Context) {
	currentPage := c.DefaultQuery("currentPage", "1")
	pageSize := c.DefaultQuery("pageSize", "10")
	categoryId := c.DefaultQuery("categoryId", "0") // 0 is all
	sort := c.DefaultQuery("sort", "0")             // 0 all, 1 price, 2 time
	//time := c.DefaultQuery("time", "0")             // 0 all, 1 asc, 2 des
	//fmt.Println("currentPage:", currentPage, "pageSize:", pageSize, "categoryId:", categoryId)
	page, e := strconv.Atoi(currentPage)
	if e != nil {
		c.JSON(400, gin.H{
			"code":        400,
			"status_code": "GLB-001",
			"msg":         "非法输入！",
		})
		return
	}
	size, e := strconv.Atoi(pageSize)
	if e != nil {
		c.JSON(400, gin.H{
			"code":        400,
			"status_code": "GLB-001",
			"msg":         "非法输入！",
		})
		return
	}
	id, e := strconv.Atoi(categoryId)
	if e != nil {
		c.JSON(400, gin.H{
			"code":        400,
			"status_code": "GLB-001",
			"msg":         "非法输入！",
		})
		return
	}
	s, e := strconv.Atoi(sort)
	if e != nil {
		c.JSON(400, gin.H{
			"code":        400,
			"status_code": "GLB-001",
			"msg":         "非法输入！",
		})
		return
	}
	req := productServerProto.GetProductListRequest{
		CurrentPage: int32(page),
		PageSize:    int32(size),
		CategoryId:  int32(id),
		Sort:        int32(s),
	}
	var resp productServerProto.GetProductListResponse
	err := instance.GrpcClient.CallProductService(model2, "getProductList", &req, &resp)
	if err != nil {
		c.JSON(500, gin.H{
			"code":        500,
			"status_code": "GLB-002",
			"msg":         "grpc调用错误: " + err.Error(),
		})
		return
	}
	c.JSON(int(resp.Code), gin.H{
		"code":        resp.Code,
		"status_code": resp.StatusCode,
		"msg":         resp.Msg,
		"data": gin.H{
			"productList": resp.ProductList,
			"totalItems":  resp.TotalItems,
			"currentPage": resp.CurrentPage,
			"pageSize":    resp.PageSize,
		},
	})
}

// UploadProductByExcel 批量上传商品excel
func UploadProductByExcel(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		fmt.Println("err:", err)
		c.JSON(400, gin.H{
			"code":        400,
			"status_code": "GLB-001",
			"msg":         "非法输入！",
		})
		return
	}
	if file.Size > 1024*1024*5 {
		c.JSON(400, gin.H{
			"code":        400,
			"status_code": "GLB-001",
			"msg":         "文件大小不能超过5MB！",
		})
		return
	}
	// 检查文件扩展名，确保是 Excel 文件
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".xlsx" && ext != ".xls" {
		c.JSON(400, gin.H{
			"code":        400,
			"status_code": "GLB-002",
			"msg":         "文件格式不正确！必须是 Excel 格式！",
		})
		return
	}

	// 检查 MIME 类型（确保文件内容是 Excel 格式）
	contentType := file.Header.Get("Content-Type")
	if contentType != "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" &&
		contentType != "application/vnd.ms-excel" {
		c.JSON(400, gin.H{
			"code":        400,
			"status_code": "GLB-003",
			"msg":         "文件类型不正确！必须是 Excel 文件！",
		})
		return
	}

	// grpc只能传输比特流，首先读取文件内容
	fileContent, e := utils.ReadFileContent(file)
	if e != nil {
		c.JSON(500, gin.H{
			"code":        500,
			"status_code": "PRT-001",
			"msg":         "读取 excel 文件内容错误！",
		})
		return
	}

	// TODO excel格式检查
	req := productServerProto.UploadProductByExcelRequest{
		File: fileContent,
	}
	var resp productServerProto.UploadProductByExcelResponse
	err = instance.GrpcClient.CallProductService(model2, "uploadProductByExcel", &req, &resp)
	if err != nil {
		c.JSON(500, gin.H{
			"code":        500,
			"status_code": "GLB-002",
			"msg":         "grpc调用错误: " + err.Error(),
		})
		return
	}
	c.JSON(int(resp.Code), gin.H{
		"code":        resp.Code,
		"status_code": resp.StatusCode,
		"msg":         resp.Msg,
	})
}

func GetProductById(c *gin.Context) {
	productId := c.DefaultQuery("productId", "0")
	id, e := strconv.Atoi(productId)
	if e != nil {
		c.JSON(500, gin.H{
			"code":        500,
			"status_code": "GLB-002",
			"msg":         "参数错误",
		})
		return
	}
	req := productServerProto.GetProductByIdRequest{
		ProductId: int32(id),
	}
	var resp productServerProto.GetProductByIdResponse
	err := instance.GrpcClient.CallProductService(model2, "getProductById", &req, &resp)
	if err != nil {
		c.JSON(500, gin.H{
			"code":        500,
			"status_code": "GLB-002",
			"msg":         "grpc调用错误: " + err.Error(),
		})
		return
	}
	c.JSON(int(resp.Code), gin.H{
		"code":        resp.Code,
		"status_code": resp.StatusCode,
		"msg":         resp.Msg,
		"data": gin.H{
			"product": resp.Product,
		},
	})
}

//// GetProductDetailById 获取详情界面商品信息
//func GetProductDetailById(c *gin.Context) {
//	productId := c.DefaultQuery("productId", "0")
//	id, e := strconv.Atoi(productId)
//	if e != nil {
//		c.JSON(500, gin.H{
//			"code":        500,
//			"status_code": "GLB-002",
//			"msg":         "参数错误",
//		})
//		return
//	}
//	req := productServerProto.GetProductDetailByIdRequest{
//		ProductId: int32(id),
//	}
//	var resp productServerProto.GetProductDetailByIdResponse
//	err := instance.GrpcClient.CallProductService(model2, "getProductDetailById", &req, &resp)
//	if err != nil {
//		c.JSON(500, gin.H{
//			"code":        500,
//			"status_code": "GLB-002",
//			"msg":         "grpc调用错误: " + err.Error(),
//		})
//		return
//	}
//	c.JSON(int(resp.Code), gin.H{
//		"code":        resp.Code,
//		"status_code": resp.StatusCode,
//		"msg":         resp.Msg,
//		"data": gin.H{
//			"product_id":    resp.ProductId,
//			"product_name":  resp.ProductName,
//			"product_price": resp.ProductPrice,
//			"product_img":   resp.ProductImg,
//			"product_type":  resp.ProductType,
//			"product_sold":  resp.Sold,
//		},
//	})
//}

// UploadSecKillProduct 上传秒杀商品（非批量）
func UploadSecKillProduct(c *gin.Context) {
	// TODO 请求头带token，测试不加
	var request productServerProto.UploadSecKillProductRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(400, gin.H{
			"code":        400,
			"status_code": "GLB-001",
			"msg":         "非法输入！",
		})
		return
	}
	var resp productServerProto.UploadSecKillProductResponse
	err = instance.GrpcClient.CallProductService(model2, "uploadSecKillProduct", &request, &resp)
	if err != nil {
		c.JSON(500, gin.H{
			"code":        500,
			"status_code": "GLB-002",
			"msg":         "grpc调用错误: " + err.Error(),
		})
		return
	}
	respCopy := proto.Clone(&resp).(*productServerProto.UploadSecKillProductResponse)
	c.JSON(int(respCopy.Code), gin.H{
		"code":        respCopy.Code,
		"status_code": respCopy.StatusCode,
		"msg":         respCopy.Msg,
	})
}

// GetSecKillList 获取秒杀商品列表
func GetSecKillList(c *gin.Context) {
	currentPage := c.DefaultQuery("currentPage", "1")
	pageSize := c.DefaultQuery("pageSize", "10")
	time := c.DefaultQuery("sessionId", "1") // 场次
	page, e := strconv.Atoi(currentPage)
	if e != nil {
		c.JSON(400, gin.H{
			"code":        400,
			"status_code": "GLB-001",
			"msg":         "非法输入！",
		})
		return
	}
	size, e := strconv.Atoi(pageSize)
	if e != nil {
		c.JSON(400, gin.H{
			"code":        400,
			"status_code": "GLB-001",
			"msg":         "非法输入！",
		})
		return
	}
	t, e := strconv.Atoi(time)
	if e != nil {
		c.JSON(400, gin.H{
			"code":        400,
			"status_code": "GLB-001",
			"msg":         "非法输入！",
		})
		return
	}
	req := productServerProto.GetSecKillListRequest{
		CurrentPage: int32(page),
		PageSize:    int32(size),
		SessionId:   int32(t),
	}
	var resp productServerProto.GetSecKillListResponse
	err := instance.GrpcClient.CallProductService(model2, "getSecKillList", &req, &resp)
	if err != nil {
		c.JSON(500, gin.H{
			"code":        500,
			"status_code": "GLB-002",
			"msg":         "grpc调用错误: " + err.Error(),
		})
		return
	}
	c.JSON(int(resp.Code), gin.H{
		"code":        resp.Code,
		"status_code": resp.StatusCode,
		"msg":         resp.Msg,
		"data": gin.H{
			"productList": resp.SecList,
			"totalItems":  resp.TotalItems,
			"currentPage": resp.CurrentPage,
			"pageSize":    resp.PageSize,
			"sessionId":   resp.SessionId,
		},
	})
}
