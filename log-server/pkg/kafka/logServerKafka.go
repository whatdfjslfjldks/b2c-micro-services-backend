package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"micro-services/log-server/pkg/kafka/model"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/IBM/sarama" // Kafka 客户端库
	"gopkg.in/yaml.v3"
	"micro-services/pkg/utils"
)

// Kafka 配置结构体
type LogKafkaConfig struct {
	Brokers []string `yaml:"brokers"`
	GroupID string   `yaml:"group_id"`
	Topic   string   `yaml:"topic"`
}

// 读取整个配置文件的结构体
type LogConfig struct {
	KafkaUserServer LogKafkaConfig `yaml:"kafka-log-server"`
}

// kafka 配置文件
var KafkaConfigLog *LogConfig

// 初始化 Kafka 配置
func InitKafkaConfig() error {
	rootPath := utils.GetCurrentPath(1)
	configPath := filepath.Join(rootPath, "../../../pkg/config", "config.yml")

	KafkaConfigLog = &LogConfig{}
	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}
	err = yaml.Unmarshal(data, KafkaConfigLog)
	if err != nil {
		return fmt.Errorf("failed to unmarshal config: %v", err)
	}
	return nil
}

// ------------------ 生产者部分 ------------------

// Kafka 生产者客户端
type LogKafkaProducer struct {
	Producer sarama.SyncProducer
	once     sync.Once
}

// 生产者初始化
func (k *LogKafkaProducer) InitProducer() error {
	var err error
	// 创建生产者配置
	config := sarama.NewConfig()
	// 等待所有副本确认
	config.Producer.Timeout = 10 * time.Second // 发送消息的超时时间
	config.Producer.Return.Successes = true    // 必须设置为 true，表示生产者需要等待 Kafka 返回成功的消息确认

	k.once.Do(func() {
		// 创建生产者客户端
		k.Producer, err = sarama.NewSyncProducer(KafkaConfigLog.KafkaUserServer.Brokers, config)
	})
	if err != nil {
		log.Printf("创建 logkafka 的生产者客户端失败： %v", err)
		return err
	}
	return nil
}

// Kafka 发布消息（生产者）
func (k *LogKafkaProducer) PublishMessage(message model.Log) error {
	// 将自定义 log 格式的 message 转为 JSON 字符串
	logBytes, er := json.Marshal(message)
	if er != nil {
		log.Printf("将自定义 log 格式的 message 转为 JSON 字符串失败： %v", er)
		return er
	}
	// 创建消息
	msg := &sarama.ProducerMessage{
		Topic: KafkaConfigLog.KafkaUserServer.Topic,
		Value: sarama.ByteEncoder(logBytes),
	}

	// 发送消息
	_, _, err := k.Producer.SendMessage(msg)
	if err != nil {
		log.Printf("发送消息失败： %v", err)
		return fmt.Errorf("failed to send message: %v", err)
	}

	log.Printf("Message sent successfully: %s", message)
	return nil
}

// --------------------- 消费者部分 -------------------------------

// Kafka 消费者客户端
type LogKafkaConsumer struct {
	ConsumerGroup sarama.ConsumerGroup
	once          sync.Once
}

// 消费者初始化
func (k *LogKafkaConsumer) InitConsumer() error {
	var err error
	// 创建消费者组配置
	config := sarama.NewConfig()
	config.Consumer.Fetch.Max = 100                       // 每次拉取的最大消息数，限制消费者消费速度
	config.Consumer.Offsets.Initial = sarama.OffsetNewest // 从最新的消息开始消费

	k.once.Do(func() {
		// 创建消费者客户端
		k.ConsumerGroup, err = sarama.NewConsumerGroup(KafkaConfigLog.KafkaUserServer.Brokers, KafkaConfigLog.KafkaUserServer.GroupID, config)
	})
	if err != nil {
		log.Printf("创建 logkafka 的消费者客户端失败： %v", err)
		return err
	}
	return nil
}

// 消费消息处理
type ConsumerHandler struct {
	MessageHandler func(message model.Log) error
}

// 设置消费者组
func (h *ConsumerHandler) Setup(sarama.ConsumerGroupSession) error {
	fmt.Println("set up")
	return nil
}

// 清理消费者组
func (h *ConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	fmt.Println("clean up")
	return nil
}

// 消费消息
func (h *ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	fmt.Println("来消息了！！")
	for message := range claim.Messages() {
		// 处理消息（通过回调函数）
		logMessage := model.Log{}
		if err := json.Unmarshal(message.Value, &logMessage); err != nil {
			log.Printf("消息解析失败: %v", err)
			continue
		}
		// 调用回调函数处理消息
		if err := h.MessageHandler(logMessage); err != nil {
			log.Printf("处理消息失败: %v", err)
		}
		// 标记消息已消费
		session.MarkMessage(message, "")
	}
	return nil
}

// 开始消费消息
func (k *LogKafkaConsumer) ConsumeMessages(handler func(message model.Log) error) error {
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
		err := k.ConsumerGroup.Consume(context.Background(), []string{KafkaConfigLog.KafkaUserServer.Topic}, consumerHandler)
		if err != nil {
			return fmt.Errorf("消费消息失败: %v", err)
		}
	}
}
