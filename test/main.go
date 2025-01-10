package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	logServerProto "micro-services/pkg/proto/log-server"
)

//
//import (
//	"fmt"
//	"micro-services/log-server/kafka/logServerKafka"
//	"micro-services/log-server/kafka/logServerKafka/model"
//	"net/http"
//	"strings"
//)
//
//func getIPFromRequest(req *http.Request) string {
//	// 1. 尝试从 "X-Forwarded-For" 获取用户的 IP 地址
//	xForwardedFor := req.Header.Get("X-Forwarded-For")
//	if xForwardedFor != "" {
//		// 如果有多个代理，"X-Forwarded-For" 会包含多个 IP 地址，逗号分隔，取第一个即为客户端的真实 IP
//		ips := strings.Split(xForwardedFor, ",")
//		return strings.TrimSpace(ips[0])
//	}
//
//	// 2. 如果没有 "X-Forwarded-For"，则直接获取 "RemoteAddr"
//	return req.RemoteAddr
//}
//
//func handler(w http.ResponseWriter, req *http.Request) {
//	// 获取客户端的 IP 地址
//	clientIP := getIPFromRequest(req)
//
//	// 返回 IP 地址
//	fmt.Fprintf(w, "Your IP Address is: %s", clientIP)
//}
//
//func main() {
//	// 设置路由
//	//http.HandleFunc("/get-ip", handler)
//	//
//	//// 启动 HTTP 服务器
//	//fmt.Println("Server started at :8080")
//	//http.ListenAndServe(":8080", nil)
//
//	consumer := &logServerKafka.LogKafkaConsumer{}
//	err := consumer.ConsumeMessages(func(message model.Log) error {
//		// 在这里处理每条消息
//		fmt.Println("处理消息: ", message)
//		return nil
//	})
//	if err != nil {
//		fmt.Println("消费消息失败: ", err)
//	}
//
//}

// 测试
func main() {

	//// -----------------创建etcd客户端--------------------------------
	//etcdClient, err := etcd.NewEtcdService(5 * time.Second)
	//if err != nil {
	//	// 丸辣！🌶
	//	panic(err)
	//}
	//// -----------------创建GRPCClient实例--------------------------------
	//GrpcClient = gClient.NewGRPCClient(etcdClient)
	//// 获取服务地址
	//serviceAddr, err := c.etcdClient.GetService("log-server")
	//if err != nil {
	//	log.Fatal("Asdf: ",err)
	//}
	//
	//fmt.Println("服务地址-----------------------： ", serviceAddr)

	// 与 gRPC 服务建立连接
	conn, err := grpc.Dial("127.0.0.1:50052", grpc.WithInsecure()) // 可改成加密连接
	if err != nil {
		log.Fatalf("failed to connect to gRPC %v", err)
	}
	fmt.Println("连接成功")
	defer conn.Close()

	// 创建 gRPC 客户端
	client := logServerProto.NewLogServiceClient(conn)

	req := &logServerProto.PostLogRequest{
		Level:       "ERROR",
		Msg:         "hello world",
		RequestPath: "/test",
		Source:      "test",
		StatusCode:  "200",
		Time:        "2023-07-01 12:00:00",
	}
	for i := 1; i <= 5; i++ {
		fmt.Println("第", i, "次发送")
		resp, err := client.PostLog(context.Background(), req)
		if err != nil {
			log.Printf("failed to send log: %v", err)
		} else {
			fmt.Println("log sent successfully:", resp)
		}
	}

	req = &logServerProto.PostLogRequest{
		Level:       "INFO",
		Msg:         "hello world",
		RequestPath: "/test",
		Source:      "test",
		StatusCode:  "200",
		Time:        "2023-07-01 12:00:00",
	}
	for i := 1; i <= 5; i++ {
		fmt.Println("第", i, "次发送")
		resp, err := client.PostLog(context.Background(), req)
		if err != nil {
			log.Printf("failed to send log: %v", err)
		} else {
			fmt.Println("log sent successfully:", resp)
		}
	}

	req = &logServerProto.PostLogRequest{
		Level:       "WARN",
		Msg:         "hello world",
		RequestPath: "/test",
		Source:      "test",
		StatusCode:  "200",
		Time:        "2023-07-01 12:00:00",
	}
	for i := 1; i <= 5; i++ {
		fmt.Println("第", i, "次发送")
		resp, err := client.PostLog(context.Background(), req)
		if err != nil {
			log.Printf("failed to send log: %v", err)
		} else {
			fmt.Println("log sent successfully:", resp)
		}
	}

}
