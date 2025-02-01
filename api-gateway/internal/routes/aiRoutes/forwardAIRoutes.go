package aiRoutes

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"log"
	aiServerProto "micro-services/pkg/proto/ai-server"
	pb "micro-services/pkg/proto/ai-server"
)

var model4 = "ai-server"

func Talk(c *gin.Context) {
	// 获取请求参数
	var req pb.TalkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"code":        400,
			"status_code": "GLB-001",
			"msg":         "参数错误！",
		})
		return
	}
	conn, err := grpc.Dial("localhost:50056", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := aiServerProto.NewAIServiceClient(conn)

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Flush()

	stream, err := client.Talk(context.Background(), &aiServerProto.TalkRequest{Prompt: req.Prompt})
	if err != nil {
		log.Fatalf("could not talk: %v", err)
	}

	for {
		resp, err := stream.Recv()
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("Stream ended")
			} else {
				log.Fatalf("Failed to receive response: %v", err)
			}
			break
		}
		//fmt.Println("Received:", resp.GetResponse())
		// 直接发送数据
		c.Writer.Write([]byte(resp.GetResponse() + "\n\n"))
		c.Writer.Flush()
	}
	//c.Writer.Flush()

}
