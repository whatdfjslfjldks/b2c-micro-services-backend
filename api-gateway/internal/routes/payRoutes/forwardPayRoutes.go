package payRoutes

import (
	"github.com/gin-gonic/gin"
	"micro-services/api-gateway/internal/instance"
	payServerProto "micro-services/pkg/proto/pay-server"
)

var model5 = "pay-server"

func TradePreCreate(c *gin.Context) {
	var req payServerProto.TradePreCreateRequest
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
	err = instance.GrpcClient.CallPayService(model5, "tradePreCreate", &req, &resp)
	if err != nil {
		c.JSON(500, gin.H{
			"code":        500,
			"status_code": "GLB-002",
			"msg":         "grpc调用错误: " + err.Error(),
		})
		return
	}
	// 假言推断
	respCopy := (resp).(*payServerProto.TradePreCreateResponse)
	c.JSON(int(respCopy.Code), gin.H{
		"code":        respCopy.Code,
		"status_code": respCopy.StatusCode,
		"msg":         respCopy.Msg,
		"data": gin.H{
			"code_url": respCopy.CodeUrl,
		},
	})

}
