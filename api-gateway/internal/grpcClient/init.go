package server_grpcClient

import (
	"micro-services/pkg/etcd"
)

// GRPCClient 封装了与 gRPC 服务的连接
type GRPCClient struct {
	etcdClient *etcd.EtcdService
}

// NewGRPCClient 构造 GRPCClient 实例
func NewGRPCClient(etcdClient *etcd.EtcdService) *GRPCClient {
	return &GRPCClient{
		etcdClient: etcdClient,
	}
}
