package productRoutes

import (
	"github.com/gin-gonic/gin"
	"micro-services/api-gateway/internal/instance"

	productServerProto "micro-services/pkg/proto/product-server"
	"strconv"
)

var model2 = "product-server"

// GetProductList 获取商品列表
func GetProductList(c *gin.Context) {
	currentPage := c.DefaultQuery("currentPage", "1")
	pageSize := c.DefaultQuery("pageSize", "10")
	categoryId := c.DefaultQuery("categoryId", "1")
	priceRange := c.DefaultQuery("range", "1")
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
	pr, e := strconv.Atoi(priceRange)
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
		PriceRange:  int32(pr),
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
		"data":        resp.ProductList,
		"totalItems":  resp.TotalItems,
		"currentPage": resp.CurrentPage,
		"pageSize":    resp.PageSize,
	})
}
