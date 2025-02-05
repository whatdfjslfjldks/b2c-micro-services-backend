package grpcClient

import (
	"micro-services/pkg/etcd"
	pb "micro-services/pkg/proto/risk-server"
)

type GRPCClient struct {
	etcdClient *etcd.EtcdService
}

// NewGRPCClient 构造 GRPCClient 实例
func NewGRPCClient(etcdClient *etcd.EtcdService) *GRPCClient {
	return &GRPCClient{
		etcdClient: etcdClient,
	}
}

func (c *GRPCClient) PostLog(request interface{}) {
	// 获取服务地址
	//serviceAddr, err := c.etcdClient.GetService("log-server")
	//if err != nil {
	//	log.Printf("failed to get service address: %v", err)
	//}
	//
	//// 与 gRPC 服务建立连接
	//conn, err := grpc.Dial(serviceAddr, grpc.WithInsecure()) // 可改成加密连接
	//if err != nil {
	//	log.Printf("failed to connect to gRPC service: %v\n", err)
	//	return
	//}
	//defer conn.Close()
	//// 创建 gRPC 客户端
	//client := logServerProto.NewLogServiceClient(conn)
	//req := request.(*logServerProto.PostLogRequest)
	//_, _ = client.PostLog(context.Background(), req)
}

func (c *GRPCClient) RiskIpAndAgentCheck(id int64, ip string, agent string, loginStatus string) *pb.RiskLoginIpAndAgentResponse {
	// 获取服务地址
	//serviceAddr, err := c.etcdClient.GetService("risk-server")
	//if err != nil {
	//	log.Printf("failed to get service address: %v", err)
	//	return nil
	//}
	//// 与 gRPC 服务建立连接
	//conn, err := grpc.Dial(serviceAddr, grpc.WithInsecure()) // 可改成加密连接
	//if err != nil {
	//	log.Printf("failed to connect to gRPC service: %v", err)
	//	return nil
	//}
	//defer conn.Close()
	//// 创建 gRPC 客户端
	//client := riskServerProto.NewRiskServiceClient(conn)
	//
	//r, _ := client.RiskLoginIpAndAgent(context.Background(), &pb.RiskLoginIpAndAgentRequest{
	//	UserId:      id,
	//	Ip:          ip,
	//	Agent:       agent,
	//	LoginStatus: loginStatus,
	//})
	//return r
	return nil
}
