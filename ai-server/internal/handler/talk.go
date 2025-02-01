package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"micro-services/pkg/proto/ai-server"
	"net/http"
)

// Talk 实现 gRPC 流式响应
func (s *Server) Talk(req *ai_server_proto.TalkRequest, stream ai_server_proto.AIService_TalkServer) error {
	// 构造 HTTP 请求体
	ollamaRequest := map[string]interface{}{
		"prompt": req.Prompt,
		"model":  "deepseek-r1:32b",
		"stream": true, // 启用流式处理
	}
	requestBody, _ := json.Marshal(ollamaRequest)

	// 向外部服务发起 HTTP 请求
	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 读取 HTTP 响应体的流
	decoder := json.NewDecoder(resp.Body)
	for {
		var data map[string]interface{}
		if err := decoder.Decode(&data); err == io.EOF {
			// 当流结束时退出
			break
		} else if err != nil {
			return err
		}

		// 提取 HTTP 响应中的 "response" 字段
		if response, exists := data["response"].(string); exists {
			// 将数据逐步发送到 gRPC 客户端
			if err := stream.Send(&ai_server_proto.TalkResponse{Response: response}); err != nil {
				return err
			}
		}
	}

	// 返回 nil 表示流结束
	return nil
}
