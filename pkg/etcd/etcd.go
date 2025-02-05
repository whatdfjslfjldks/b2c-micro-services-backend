package etcd

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"os"
	"time"
)

// EtcdService represents the etcd service
type EtcdService struct {
	client *clientv3.Client
	//lease  clientv3.LeaseID // 租约ID，用于续约
}

// NewEtcdService setup EtcdService instance
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

// RegisterService 60s租约，每30s检测一次状态，心跳机制
func (es *EtcdService) RegisterService(serviceName, serviceAddr string, ttl int64) error {
	// 创建一个租约（TTL 设定为 60 秒）
	leaseResp, err := es.client.Grant(context.Background(), ttl)
	if err != nil {
		return fmt.Errorf("etcd租约创建失败:%v", err)
	}
	fmt.Printf("租约ID: %d\n", leaseResp.ID)
	// 保存租约ID
	//es.lease = leaseResp.ID

	// 注册服务时，设置带有租约的键值对
	key := fmt.Sprintf("%s%s", "services:", serviceName)
	value := serviceAddr
	_, err = es.client.Put(context.Background(), key, value, clientv3.WithLease(leaseResp.ID))
	if err != nil {
		return fmt.Errorf("etcd注册服务失败:%v", err)
	}
	fmt.Printf("服务 %s 注册成功，地址是 %s\n", serviceName, serviceAddr)

	// 开始续约
	go es.keepAlive(leaseResp.ID)

	return nil
}

// keepAlive starts a goroutine that keeps renewing the lease
func (es *EtcdService) keepAlive(leaseId clientv3.LeaseID) {
	// 通过 KeepAlive 方法续约租约
	ch, err := es.client.KeepAlive(context.Background(), leaseId)
	if err != nil {
		fmt.Printf("租约续约失败: %v\n", err)
		return
	}

	// 监听租约续约的响应
	for {
		select {
		case <-ch:
			fmt.Println(leaseId, "服务续约成功")
		case <-time.After(30 * time.Second):
			// 每 30 秒检查一次租约续约的状态
			fmt.Println("租约续约正在进行中...")
		}
	}
}

// GetService get etcdService group
func (es *EtcdService) GetService(serviceName string) ([]string, error) {
	// 构造 etcd 中的键
	key := fmt.Sprintf("%s%s", "services:", serviceName)

	// 从 etcd 中获取指定键的信息
	resp, err := es.client.Get(context.Background(), key, clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("etcd获取服务失败:%v", err)
	}

	// 检查是否找到了相关的键值对
	if len(resp.Kvs) == 0 {
		return nil, fmt.Errorf("没有找到服务:%s", serviceName)
	}

	// 用于存储提取出来的 value
	values := make([]string, 0, len(resp.Kvs))
	// 遍历获取到的键值对
	for _, kv := range resp.Kvs {
		// 将每个键值对的 value 转换为字符串并添加到 values 切片中
		values = append(values, string(kv.Value))
	}

	return values, nil
}

// GetAllServices get all etcd services
func (es *EtcdService) GetAllServices() (map[string]string, error) {
	// 设置超时时间 5 秒
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // 确保在函数返回时取消上下文
	// 获取以 "services:" 为前缀的所有服务
	resp, err := es.client.Get(ctx, "services:", clientv3.WithPrefix())
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

// Close unopened etcd
func (es *EtcdService) Close() error {
	return es.client.Close()
}
