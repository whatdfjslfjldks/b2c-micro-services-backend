package recommendRoutes

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/proto"
	"micro-services/api-gateway/internal/instance"
	recommendServerProto "micro-services/pkg/proto/recommend-server"
	"strconv"
)

var model3 = "recommend-server"

// ClickProduct 点击
func ClickProduct(c *gin.Context) {
	var request recommendServerProto.ClickProductRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(400, gin.H{
			"code":        400,
			"status_code": "GLB-001",
			"msg":         "非法输入！",
		})
		return
	}
	req := recommendServerProto.ClickProductRequest{
		UserId:    request.UserId,
		ProductId: request.ProductId,
	}
	var resp recommendServerProto.ClickProductResponse
	err = instance.GrpcClient.CallRecommendService(model3, "clickProduct", &req, &resp)
	if err != nil {
		c.JSON(500, gin.H{
			"code":        500,
			"status_code": "GLB-002",
			"msg":         "grpc调用错误: " + err.Error(),
		})
		return
	}
	// 数据埋点不返回实际内容
	c.JSON(200, gin.H{})
}

// SearchProduct 搜索
func SearchProduct(c *gin.Context) {
	var request recommendServerProto.SearchProductRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(400, gin.H{
			"code":        400,
			"status_code": "GLB-001",
			"msg":         "非法输入！",
		})
		return
	}
	req := recommendServerProto.SearchProductRequest{
		UserId:  request.UserId,
		Keyword: request.Keyword,
	}
	var resp recommendServerProto.SearchProductResponse
	err = instance.GrpcClient.CallRecommendService(model3, "searchProduct", &req, &resp)
	if err != nil {
		c.JSON(500, gin.H{
			"code":        500,
			"status_code": "GLB-002",
			"msg":         "grpc调用错误: " + err.Error(),
		})
		return
	}
	// 数据埋点不返回实际内容
	c.JSON(200, gin.H{})
}

// BrowseProduct 浏览
func BrowseProduct(c *gin.Context) {
	var request recommendServerProto.BrowseProductRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(400, gin.H{
			"code":        400,
			"status_code": "GLB-001",
			"msg":         "非法输入！",
		})
		return
	}
	req := recommendServerProto.BrowseProductRequest{
		UserId:    request.UserId,
		ProductId: request.ProductId,
	}
	var resp recommendServerProto.BrowseProductResponse
	err = instance.GrpcClient.CallRecommendService(model3, "browseProduct", &req, &resp)
	if err != nil {
		c.JSON(500, gin.H{
			"code":        500,
			"status_code": "GLB-002",
			"msg":         "grpc调用错误: " + err.Error(),
		})
		return
	}
	// 数据埋点不返回实际内容
	c.JSON(200, gin.H{})
}

// PurchaseProduct 购买
func PurchaseProduct(c *gin.Context) {
	var request recommendServerProto.PurchaseProductRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(400, gin.H{
			"code":        400,
			"status_code": "GLB-001",
			"msg":         "非法输入！",
		})
		return
	}
	req := recommendServerProto.PurchaseProductRequest{
		UserId:    request.UserId,
		ProductId: request.ProductId,
		Quantity:  request.Quantity,
	}
	var resp recommendServerProto.PurchaseProductResponse
	err = instance.GrpcClient.CallRecommendService(model3, "purchaseProduct", &req, &resp)
	if err != nil {
		c.JSON(500, gin.H{
			"code":        500,
			"status_code": "GLB-002",
			"msg":         "grpc调用错误: " + err.Error(),
		})
		return
	}
	// 数据埋点不返回实际内容
	c.JSON(200, gin.H{})
}

// GetRecommendProductList 获取推荐商品
func GetRecommendProductList(c *gin.Context) {
	userId := c.DefaultQuery("userId", "-1")
	id, e := strconv.Atoi(userId)
	if e != nil {
		c.JSON(400, gin.H{
			"code":        400,
			"status_code": "GLB-001",
			"msg":         "非法输入！",
		})
		return
	}
	req := recommendServerProto.GetRecommendProductListRequest{
		UserId: int64(id),
		//CurrentPage: int32(size),
		//PageSize:    int32(page),
	}
	var resp recommendServerProto.GetRecommendProductListResponse
	err := instance.GrpcClient.CallRecommendService(model3, "getRecommendProductList", &req, &resp)
	if err != nil {
		c.JSON(500, gin.H{
			"code":        500,
			"status_code": "GLB-002",
			"msg":         "grpc调用错误: " + err.Error(),
		})
		return
	}
	respCopy := proto.Clone(&resp).(*recommendServerProto.GetRecommendProductListResponse)
	c.JSON(int(resp.Code), gin.H{
		"code":        respCopy.Code,
		"status_code": respCopy.StatusCode,
		"msg":         respCopy.Msg,
		"data":        respCopy.ProductList,
	})
}
