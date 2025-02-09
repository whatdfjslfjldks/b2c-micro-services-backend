package kafka

import (
	"bytes"
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"gopkg.in/yaml.v3"
	"log"
	"micro-services/pkg/utils"
	"os"
	"path/filepath"
	"time"
)

// 订单状态常量
const (
	OrderPending       = 0
	OrderPaid          = 1
	OrderCompleted     = 2
	OrderShipped       = 3
	OrderDelivered     = 4
	OrderCancelled     = 5
	OrderPaymentFailed = 6
	OrderRefunded      = 7
)

// OrderKafkaConfig Kafka 配置结构体
type OrderKafkaConfig struct {
	Brokers []string `yaml:"brokers"`
	GroupID string   `yaml:"group_id"`
	Topic   string   `yaml:"topic"`
}

// OrderConfig 读取整个配置文件的结构体
type OrderConfig struct {
	KafkaUserServer OrderKafkaConfig `yaml:"kafka-order-server"`
}

// KafkaConfigOrder Kafka 配置文件
var KafkaConfigOrder *OrderConfig

// InitKafkaConfig 初始化 Kafka 配置
func InitKafkaConfig() error {
	rootPath := utils.GetCurrentPath(1)
	configPath := filepath.Join(rootPath, "../../../pkg/config", "config.yml")

	KafkaConfigOrder = &OrderConfig{}
	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}
	err = yaml.Unmarshal(data, KafkaConfigOrder)
	if err != nil {
		return fmt.Errorf("failed to unmarshal config: %v", err)
	}
	return nil
}

// NewKafkaProducer 创建一个 Kafka 生产者
func NewKafkaProducer() (sarama.SyncProducer, error) {
	producerConfig := sarama.NewConfig()
	producerConfig.Producer.RequiredAcks = sarama.WaitForAll
	producerConfig.Producer.Retry.Max = 5
	producerConfig.Producer.Return.Successes = true

	return sarama.NewSyncProducer(KafkaConfigOrder.KafkaUserServer.Brokers, producerConfig)
}

// SendMessageToPartition 发送消息到指定分区
func SendMessageToPartition(partition int32, orderId string) error {
	producer, err := NewKafkaProducer()
	if err != nil {
		return fmt.Errorf("failed to create Kafka producer: %v", err)
	}
	defer func(producer sarama.SyncProducer) {
		err := producer.Close()
		if err != nil {
			log.Printf("Failed to close Kafka producer: %v", err)
		}
	}(producer)

	// 将 orderId 转换为 []byte 类型的 key
	key := []byte(orderId)

	// 将订单状态转为 []byte 类型的 value
	value := []byte(fmt.Sprintf("%d", partition))

	message := &sarama.ProducerMessage{
		Topic:     KafkaConfigOrder.KafkaUserServer.Topic,
		Key:       sarama.ByteEncoder(key),
		Value:     sarama.ByteEncoder(value),
		Partition: partition,
	}

	// 发送消息
	_, _, err = producer.SendMessage(message)
	return err
}

// NewKafkaConsumerGroup 创建一个 Kafka 消费者组
func NewKafkaConsumerGroup() (sarama.ConsumerGroup, error) {
	consumerConfig := sarama.NewConfig()
	consumerConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	consumerConfig.Consumer.Return.Errors = true

	return sarama.NewConsumerGroup(KafkaConfigOrder.KafkaUserServer.Brokers, KafkaConfigOrder.KafkaUserServer.GroupID, consumerConfig)
}

// ConsumePartition 消费指定分区并过滤指定key的消息
func ConsumePartition(partition int32, orderId string) error {
	consumerGroup, err := NewKafkaConsumerGroup()
	if err != nil {
		return fmt.Errorf("failed to create Kafka consumer group: %v", err)
	}
	defer consumerGroup.Close()

	keyToConsume := []byte(orderId)
	handler := &PartitionConsumerHandler{
		Partition:    partition,
		KeyToConsume: keyToConsume, // 传入目标 key
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = consumerGroup.Consume(ctx, []string{KafkaConfigOrder.KafkaUserServer.Topic}, handler)
	if err != nil {
		return fmt.Errorf("failed to consume messages: %v", err)
	}

	// 检查是否处理了消息
	if !handler.Handled {
		return fmt.Errorf("no message with key %s found", orderId)
	}

	return nil
}

// PartitionConsumerHandler 实现 sarama.ConsumerGroupHandler 接口
type PartitionConsumerHandler struct {
	Partition    int32
	KeyToConsume []byte
	Handled      bool // 标志，表示是否已经处理了一个消息
}

func (h *PartitionConsumerHandler) Setup(session sarama.ConsumerGroupSession) error {
	// 你可以在这里做一些初始化工作
	return nil
}

func (h *PartitionConsumerHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	// 清理工作
	return nil
}

func (h *PartitionConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// 消费消息并过滤指定 key
	for message := range claim.Messages() {
		if bytes.Equal(message.Key, h.KeyToConsume) {
			// 仅当 key 匹配时才处理消息
			fmt.Printf("Received matching message from partition %d: key=%s, value=%s\n", message.Partition, string(message.Key), string(message.Value))
			session.MarkMessage(message, "")
			h.Handled = true
			return nil // 处理完一个消息后立即返回
		}
	}
	return nil
}

func (h *PartitionConsumerHandler) Err() error {
	return nil
}
