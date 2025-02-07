package orderRoutes

import (
	"github.com/gin-gonic/gin"
	"log"
	"micro-services/api-gateway/internal/instance"
	orderServerProto "micro-services/pkg/proto/order-server"
	"strconv"
)

var model6 = "order-server"

type Order struct {
	Address    string `json:"address"`
	Detail     string `json:"detail"`
	Name       string `json:"name"`
	Phone      string `json:"phone"`
	Note       string `json:"note"`
	OrderItems Cart   `json:"orderItems"`
}

// Cart 表示购物车
type Cart struct {
	ProductList []Product `json:"productList"`
}

// Product 表示购物车中的商品
type Product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Cover    string  `json:"cover"`     // 封面，选列表第一个
	TypeName string  `json:"type_name"` // 类型名称
	Price    float64 `json:"price"`     // 单价
	Amount   int     `json:"amount"`    // 购买数量
}

// CreateOrder 创建订单
func CreateOrder(c *gin.Context) {
	accessToken := c.Request.Header.Get("Access-Token")
	var a Order
	if err := c.ShouldBindJSON(&a); err != nil {
		log.Printf("err: %v", err)
		c.JSON(400, gin.H{
			"code":        400,
			"status_code": "GLB-001",
			"msg":         "非法输入！",
		})
		return
	}

	var productIDs, productAmounts []int32
	var typeNames []string
	for _, product := range a.OrderItems.ProductList {
		productIDs = append(productIDs, int32(product.ID))
		typeNames = append(typeNames, product.TypeName)
		productAmounts = append(productAmounts, int32(product.Amount))
	}

	ph, err := strconv.Atoi(a.Phone)
	if err != nil {
		c.JSON(400, gin.H{
			"code":        400,
			"status_code": "GLB-001",
			"msg":         "非法输入！",
		})
		return
	}

	req := orderServerProto.CreateOrderRequest{
		AccessToken:   accessToken,
		Address:       a.Address,
		Detail:        a.Detail,
		Name:          a.Name,
		Phone:         int32(ph),
		Note:          a.Note,
		ProductId:     productIDs,
		TypeName:      typeNames,
		ProductAmount: productAmounts,
	}
	var resp interface{}
	err = instance.GrpcClient.CallOrderService(model6, "createOrder", &req, &resp)
	if err != nil {
		c.JSON(500, gin.H{
			"code":        500,
			"status_code": "GLB-002",
			"msg":         "grpc调用错误: " + err.Error(),
		})
		return
	}
	respCopy := (resp).(*orderServerProto.CreateOrderResponse)
	c.JSON(int(respCopy.Code), gin.H{
		"code":        respCopy.Code,
		"status_code": respCopy.StatusCode,
		"msg":         respCopy.Msg,
		"data": gin.H{
			"orderId": respCopy.OrderId,
		},
	})
}
