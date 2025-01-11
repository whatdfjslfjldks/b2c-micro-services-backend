package instance

import (
	"micro-services/pkg/etcd"
	gClient "micro-services/product-server/pkg/grpcClient"
	"time"
)

var GrpcClient *gClient.GRPCClient

func NewInstance() {
	// -----------------创建etcd客户端--------------------------------
	etcdClient, err := etcd.NewEtcdService(5 * time.Second)
	if err != nil {
		// 丸辣！🌶
		panic(err)
	}
	// -----------------创建GRPCClient实例--------------------------------
	GrpcClient = gClient.NewGRPCClient(etcdClient)
}
