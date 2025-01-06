package pkgConfig

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/IBM/sarama" // Kafka 客户端库
	"gopkg.in/yaml.v3"
	"micro-services/pkg/utils"
)

// Kafka 配置结构体
type KafkaConfig struct {
	Brokers []string `yaml:"brokers"`
	GroupID string   `yaml:"group_id"`
	Topic   string   `yaml:"topic"`
}

// 读取整个配置文件的结构体
type Config struct {
	KafkaUserServer KafkaConfig `yaml:"kafka-user-server"`
}

var config *Config

// 读取配置文件
func InitKafkaConfig() error {
	rootPath := utils.GetCurrentPath(1)
	configPath := filepath.Join(rootPath, "../config", "config.yml")

	config = &Config{}
	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return fmt.Errorf("failed to unmarshal config: %v", err)
	}
	return nil
}

// Kafka 发布消息（生产者）
func PublishMessage(message string) error {
	// 创建生产者客户端
	producer, err := sarama.NewSyncProducer(config.KafkaUserServer.Brokers, nil)
	if err != nil {
		return fmt.Errorf("failed to create producer: %v", err)
	}
	defer producer.Close()

	// 创建消息
	msg := &sarama.ProducerMessage{
		Topic: config.KafkaUserServer.Topic,
		Value: sarama.StringEncoder(message),
	}

	// 发送消息
	_, _, err = producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	log.Printf("Message sent successfully: %s", message)
	return nil
}

// Kafka 消费消息（消费者）
type ConsumerHandler struct{}

func (h *ConsumerHandler) Setup(sarama.ConsumerGroupSession) error {
	fmt.Println("set up")
	return nil
}
func (h *ConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	fmt.Println("clean up")
	return nil
}
func (h *ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	fmt.Println("来消息了！！")
	for message := range claim.Messages() {
		// TODO 处理消息
		fmt.Println("test: ", string(message.Value))
		session.MarkMessage(message, "")
	}
	return nil
}

func ConsumeMessages() error {
	// 检查配置是否完整
	if config == nil {
		return fmt.Errorf("Kafka config is nil")
	}

	if len(config.KafkaUserServer.Brokers) == 0 || config.KafkaUserServer.GroupID == "" || config.KafkaUserServer.Topic == "" {
		return fmt.Errorf("Kafka configuration is missing or incomplete. Brokers: %v, GroupID: %s, Topic: %s",
			config.KafkaUserServer.Brokers, config.KafkaUserServer.GroupID, config.KafkaUserServer.Topic)
	}

	// 创建消费者客户端
	consumer, err := sarama.NewConsumerGroup(config.KafkaUserServer.Brokers, config.KafkaUserServer.GroupID, nil)
	if err != nil {
		return fmt.Errorf("failed to create consumer group: %v", err)
	}
	defer consumer.Close()

	// 创建消费者处理器
	handler := &ConsumerHandler{}
	for {
		// 消费消息
		err := consumer.Consume(context.Background(), []string{config.KafkaUserServer.Topic}, handler)
		if err != nil {
			return fmt.Errorf("failed to consume messages: %v", err)
		}
	}
}
