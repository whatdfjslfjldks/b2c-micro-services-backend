package grpcClient

import (
	"micro-services/pkg/etcd"
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
