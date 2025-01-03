package etcd

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"os"
	"time"
)

type EtcdService struct {
	client *clientv3.Client
}

// setup EtcdService instance
func NewEtcdService(dialTimeout time.Duration) (*EtcdService, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{os.Getenv("api") + ":2379"}, // etcd节点地址
		DialTimeout: dialTimeout,                          // 连接超时时间
	})
	if err != nil {
		return nil, fmt.Errorf("etcd连接失败:%v", err)
	}
	return &EtcdService{
		client: client,
	}, nil
}

// register etcdservice
func (es *EtcdService) RegisterService(serviceName, serviceAddr string) error {
	key := fmt.Sprintf("%s%s", "services:", serviceName)
	value := serviceAddr
	_, err := es.client.Put(context.Background(), key, value)
	if err != nil {
		return fmt.Errorf("etcd注册服务失败:%v", err)
	}
	fmt.Printf("服务 %s 注册成功，地址是 %s\n", serviceName, serviceAddr)
	return nil
}

// get etcdservice
func (es *EtcdService) GetService(serviceName string) (string, error) {
	key := fmt.Sprintf("%s%s", "services:", serviceName)
	resp, err := es.client.Get(context.Background(), key)
	if err != nil {
		return "", fmt.Errorf("etcd获取服务失败:%v", err)
	}
	if len(resp.Kvs) == 0 {
		return "", fmt.Errorf("没有找到服务:%s", serviceName)
	}
	return string(resp.Kvs[0].Value), nil
}

// get all etcd services
func (es *EtcdService) GetAllServices() (map[string]string, error) {
	// 获取以 "services:" 为前缀的所有服务
	resp, err := es.client.Get(context.Background(), "services:", clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("etcd获取所有服务失败:%v", err)
	}

	// 保存所有服务地址的映射
	serviceAddresses := make(map[string]string)

	// 遍历所有服务的键值对
	for _, kv := range resp.Kvs {
		// 从键中提取服务名称和地址
		serviceName, serviceAddr := extractServiceDetails(string(kv.Key), string(kv.Value))
		if serviceName != "" {
			// 将服务名称和地址保存到映射中
			serviceAddresses[serviceName] = serviceAddr
		}
	}

	return serviceAddresses, nil
}

// 从键中提取服务名称和地址
func extractServiceDetails(key, value string) (string, string) {
	// 假设服务路径格式是 "services:{service_name}:{port}"
	// 从键中提取服务名称（即"{service_name}"）
	parts := splitKey(key, ":")
	if len(parts) >= 2 {
		serviceName := parts[1]
		// 返回服务名称和地址
		return serviceName, value
	}
	return "", ""
}

// 简单的字符串分割方法
func splitKey(key, delimiter string) []string {
	var result []string
	start := 0
	for i := 0; i < len(key); i++ {
		if string(key[i]) == delimiter {
			if i > start {
				result = append(result, key[start:i])
			}
			start = i + 1
		}
	}
	if start < len(key) {
		result = append(result, key[start:])
	}
	return result
}

// close etcd
func (es *EtcdService) Close() error {
	return es.client.Close()
}
