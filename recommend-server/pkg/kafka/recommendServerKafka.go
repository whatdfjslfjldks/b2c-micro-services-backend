package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"micro-services/recommend-server/pkg/kafka/model"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/IBM/sarama" // Kafka 客户端库
	"gopkg.in/yaml.v3"
	"micro-services/pkg/utils"
)

// RecommendKafkaConfig Kafka 配置结构体
type RecommendKafkaConfig struct {
	Brokers []string `yaml:"brokers"`
	GroupID string   `yaml:"group_id"`
	Topic   string   `yaml:"topic"`
}

// RecommendConfig 读取整个配置文件的结构体
type RecommendConfig struct {
	KafkaUserServer RecommendKafkaConfig `yaml:"kafka-recommend-server"`
}

// KafkaConfigRecommend kafka 配置文件
var KafkaConfigRecommend *RecommendConfig

// InitKafkaConfig 初始化 Kafka 配置
func InitKafkaConfig() error {
	rootPath := utils.GetCurrentPath(1)
	configPath := filepath.Join(rootPath, "../../../pkg/config", "config.yml")

	KafkaConfigRecommend = &RecommendConfig{}
	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}
	err = yaml.Unmarshal(data, KafkaConfigRecommend)
	if err != nil {
		return fmt.Errorf("failed to unmarshal config: %v", err)
	}
	return nil
}

// ------------------ 生产者部分 ------------------

// RecommendKafkaProducer Kafka 生产者客户端
type RecommendKafkaProducer struct {
	Producer sarama.SyncProducer
	once     sync.Once
}

// InitProducer 生产者初始化
func (k *RecommendKafkaProducer) InitProducer() error {
	var err error
	// 创建生产者配置
	config := sarama.NewConfig()
	// 等待所有副本确认
	config.Producer.Timeout = 10 * time.Second // 发送消息的超时时间
	config.Producer.Return.Successes = true    // 必须设置为 true，表示生产者需要等待 Kafka 返回成功的消息确认
	// TODO 非常重要！！！ sarama操纵kafka发送指定partition分区，需要配置这一项！！！！！
	config.Producer.Partitioner = sarama.NewManualPartitioner
	k.once.Do(func() {
		// 创建生产者客户端
		k.Producer, err = sarama.NewSyncProducer(KafkaConfigRecommend.KafkaUserServer.Brokers, config)
	})
	if err != nil {
		log.Printf("创建 logkafka 的生产者客户端失败： %v", err)
		return err
	}
	return nil
}

// PublishMessage Kafka 发布消息（根据日志级别发送到不同的分区）
// 0-click, 1-purchase, 2-search, 3-browse
func (k *RecommendKafkaProducer) PublishMessage(message model.Recommend, partition int32) error {
	// 将自定义 log 格式的 message 转为 JSON 字符串
	recommendBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("将消息转换为 JSON 失败： %v", err)
		return err
	}

	// 创建消息并指定分区
	msg := &sarama.ProducerMessage{
		Partition: partition,
		Topic:     KafkaConfigRecommend.KafkaUserServer.Topic,
		Value:     sarama.ByteEncoder(recommendBytes),
	}
	//fmt.Println("分区", partition)

	// 发送消息
	_, _, err = k.Producer.SendMessage(msg)
	if err != nil {
		log.Printf("发送消息失败： %v", err)
		return fmt.Errorf("failed to send message: %v", err)
	}

	log.Printf("Message sent successfully: %s", message)
	return nil
}

// --------------------- 消费者部分 -------------------------------

// RecommendKafkaConsumer Kafka 消费者客户端
type RecommendKafkaConsumer struct {
	ConsumerGroup sarama.ConsumerGroup
	once          sync.Once
}

// InitConsumer 消费者初始化
func (k *RecommendKafkaConsumer) InitConsumer() error {
	var err error
	// 创建消费者组配置
	config := sarama.NewConfig()
	config.Consumer.Fetch.Max = 100               // 每次拉取的最大消息数，限制消费者消费速度
	config.Consumer.MaxWaitTime = 1 * time.Second // 每次最多等待 10 秒

	config.Consumer.Offsets.Initial = sarama.OffsetNewest         // TODO 从最新的消息开始消费,根据需要可以做调整
	config.Consumer.Offsets.AutoCommit.Enable = true              // 默认设置是 true，表示自动提交
	config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second // 自动提交间隔

	k.once.Do(func() {
		// 创建消费者客户端
		k.ConsumerGroup, err = sarama.NewConsumerGroup(KafkaConfigRecommend.KafkaUserServer.Brokers, KafkaConfigRecommend.KafkaUserServer.GroupID, config)
	})
	if err != nil {
		log.Printf("创建 logkafka 的消费者客户端失败： %v", err)
		return err
	}
	return nil
}

// ConsumerHandler 消费消息处理
type ConsumerHandler struct {
	MessageHandler func(message model.Recommend) error
}

// Setup 设置消费者组
func (h *ConsumerHandler) Setup(sarama.ConsumerGroupSession) error {
	fmt.Println("set up")
	return nil
}

// Cleanup 清理消费者组
func (h *ConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	fmt.Println("clean up")
	return nil
}

// ConsumeClaim 消费消息
func (h *ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		// 处理消息（通过回调函数）
		recommendMessage := model.Recommend{}
		if err := json.Unmarshal(message.Value, &recommendMessage); err != nil {
			log.Printf("消息解析失败: %v", err)
			continue
		}
		//fmt.Println("分区222：:", logMessage)

		// 调用回调函数处理消息
		if err := h.MessageHandler(recommendMessage); err != nil {
			log.Printf("处理消息失败: %v", err)
		}
		// 标记消息已消费
		session.MarkMessage(message, "")
	}
	return nil
}

// ConsumeMessages 开始消费消息
func (k *RecommendKafkaConsumer) ConsumeMessages(ctx context.Context, handler func(message model.Recommend) error) error {
	// 初始化 Kafka 配置
	err := InitKafkaConfig()
	if err != nil {
		fmt.Println("init kafka config failed:", err)
		return err
	}

	// 初始化消费者
	err = k.InitConsumer()
	if err != nil {
		return fmt.Errorf("初始化消费者失败: %v", err)
	}
	defer k.ConsumerGroup.Close()

	// 创建消费者处理器
	consumerHandler := &ConsumerHandler{
		MessageHandler: handler, // 使用传入的回调函数处理消息
	}

	// 开始消费消息
	for {
		// 消费消息
		err := k.ConsumerGroup.Consume(ctx, []string{KafkaConfigRecommend.KafkaUserServer.Topic}, consumerHandler)
		if err != nil {
			return fmt.Errorf("消费消息失败: %v", err)
		}
	}
}
