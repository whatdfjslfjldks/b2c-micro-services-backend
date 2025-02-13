package productRoutes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
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
	var resp interface{}
	err := instance.GrpcClient.CallProductService(model2, "getProductList", &req, &resp)
	if err != nil {
		c.JSON(500, gin.H{
			"code":        500,
			"status_code": "GLB-002",
			"msg":         "grpc调用错误: " + err.Error(),
		})
		return
	}
	respCopy := resp.(*productServerProto.GetProductListResponse)
	c.JSON(int(respCopy.Code), gin.H{
		"code":        respCopy.Code,
		"status_code": respCopy.StatusCode,
		"msg":         respCopy.Msg,
		"data": gin.H{
			"productList": respCopy.ProductList,
			"totalItems":  respCopy.TotalItems,
			"currentPage": respCopy.CurrentPage,
			"pageSize":    respCopy.PageSize,
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
	var resp interface{}
	err = instance.GrpcClient.CallProductService(model2, "uploadProductByExcel", &req, &resp)
	if err != nil {
		c.JSON(500, gin.H{
			"code":        500,
			"status_code": "GLB-002",
			"msg":         "grpc调用错误: " + err.Error(),
		})
		return
	}
	respCopy := resp.(*productServerProto.UploadProductByExcelResponse)
	c.JSON(int(respCopy.Code), gin.H{
		"code":        respCopy.Code,
		"status_code": respCopy.StatusCode,
		"msg":         respCopy.Msg,
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
	var resp interface{}
	err := instance.GrpcClient.CallProductService(model2, "getProductById", &req, &resp)
	if err != nil {
		c.JSON(500, gin.H{
			"code":        500,
			"status_code": "GLB-002",
			"msg":         "grpc调用错误: " + err.Error(),
		})
		return
	}
	respCopy := resp.(*productServerProto.GetProductByIdResponse)
	c.JSON(int(respCopy.Code), gin.H{
		"code":        respCopy.Code,
		"status_code": respCopy.StatusCode,
		"msg":         respCopy.Msg,
		"data": gin.H{
			"product": respCopy.Product,
		},
	})
}

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
	var resp interface{}
	err = instance.GrpcClient.CallProductService(model2, "uploadSecKillProduct", &request, &resp)
	if err != nil {
		c.JSON(500, gin.H{
			"code":        500,
			"status_code": "GLB-002",
			"msg":         "grpc调用错误: " + err.Error(),
		})
		return
	}
	respCopy := (resp).(*productServerProto.UploadSecKillProductResponse)
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
	var resp interface{}
	err := instance.GrpcClient.CallProductService(model2, "getSecKillList", &req, &resp)
	if err != nil {
		c.JSON(500, gin.H{
			"code":        500,
			"status_code": "GLB-002",
			"msg":         "grpc调用错误: " + err.Error(),
		})
		return
	}
	respCopy := resp.(*productServerProto.GetSecKillListResponse)
	c.JSON(int(respCopy.Code), gin.H{
		"code":        respCopy.Code,
		"status_code": respCopy.StatusCode,
		"msg":         respCopy.Msg,
		"data": gin.H{
			"productList": respCopy.SecList,
			"totalItems":  respCopy.TotalItems,
			"currentPage": respCopy.CurrentPage,
			"pageSize":    respCopy.PageSize,
			"sessionId":   respCopy.SessionId,
		},
	})
}

// PurchaseSecKill 购买秒杀商品
func PurchaseSecKill(c *gin.Context) {
	accessToken := c.Request.Header.Get("Access-Token")
	type a struct {
		ProductId int32 `json:"productId"`
	}
	var b a
	if err := c.ShouldBindJSON(&b); err != nil {
		log.Printf("err: %v", err)
		c.JSON(400, gin.H{
			"code":        400,
			"status_code": "GLB-001",
			"msg":         "非法输入！",
		})
		return
	}
	req := productServerProto.PurchaseSecKillRequest{
		Id:          b.ProductId,
		AccessToken: accessToken,
	}
	var resp interface{}
	err := instance.GrpcClient.CallProductService(model2, "purchaseSecKill", &req, &resp)
	if err != nil {
		c.JSON(500, gin.H{
			"code":        500,
			"status_code": "GLB-002",
			"msg":         "grpc调用错误: " + err.Error(),
		})
		return
	}
	respCopy := resp.(*productServerProto.PurchaseSecKillResponse)
	c.JSON(int(respCopy.Code), gin.H{
		"code":        respCopy.Code,
		"status_code": respCopy.StatusCode,
		"msg":         respCopy.Msg,
	})
}

func GetOrderConfirmProduct(c *gin.Context) {
	var req productServerProto.GetOrderConfirmProductRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(400, gin.H{
			"code":        400,
			"status_code": "GLB-001",
			"msg":         "非法输入！",
		})
		return
	}
	var resp interface{}
	err = instance.GrpcClient.CallProductService(model2, "getOrderConfirmProduct", &req, &resp)
	if err != nil {
		c.JSON(500, gin.H{
			"code":        500,
			"status_code": "GLB-002",
			"msg":         "grpc调用错误: " + err.Error(),
		})
		return
	}
	respCopy := resp.(*productServerProto.GetOrderConfirmProductResponse)
	c.JSON(int(respCopy.Code), gin.H{
		"code":        respCopy.Code,
		"status_code": respCopy.StatusCode,
		"msg":         respCopy.Msg,
		"data": gin.H{
			"products": respCopy.Products,
		},
	})
}
