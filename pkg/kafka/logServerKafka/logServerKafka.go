package logServerKafka

import (
	"encoding/json"
	"fmt"
	"log"
	"micro-services/pkg/kafka/logServerKafka/model"
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

// kafka 生产者客户端
type LogKafkaProducer struct {
	Producer sarama.SyncProducer
	once     sync.Once
}

// 读取配置文件,初始化 kafka，可以在程序初启动执行
func InitKafkaConfig() error {
	rootPath := utils.GetCurrentPath(1)
	configPath := filepath.Join(rootPath, "../config", "config.yml")

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

// 生产者客户端
func (k *LogKafkaProducer) InitProducer() error {
	var err error
	// 定义生产者配置
	config := sarama.NewConfig()
	// 等待所有副本确认,可以提高消息的容错性，但是会导致速度减慢
	// 日志采集不太需要消息高可靠，允许丢失
	//config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Timeout = 10 * time.Second // 发送消息的超时时间

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
	// 将自定义log格式的message转为json字符串
	logBytes, er := json.Marshal(message)
	if er != nil {
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
		return fmt.Errorf("failed to send message: %v", err)
	}

	log.Printf("Message sent successfully: %s", message)
	return nil
}

//
//// Kafka 消费消息（消费者）
//type ConsumerHandler struct{}
//
//func (h *ConsumerHandler) Setup(sarama.ConsumerGroupSession) error {
//	fmt.Println("set up")
//	return nil
//}
//func (h *ConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
//	fmt.Println("clean up")
//	return nil
//}
//func (h *ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
//	fmt.Println("来消息了！！")
//	for message := range claim.Messages() {
//		// TODO 处理消息
//		fmt.Println("test: ", string(message.Value))
//		session.MarkMessage(message, "")
//	}
//	return nil
//}
//
//func ConsumeMessages() error {
//	// 检查配置是否完整
//	if kafkaConfig == nil {
//		return fmt.Errorf("Kafka config is nil")
//	}
//
//	if len(kafkaConfig.KafkaUserServer.Brokers) == 0 || kafkaConfig.KafkaUserServer.GroupID == "" || kafkaConfig.KafkaUserServer.Topic == "" {
//		return fmt.Errorf("kafka configuration is missing or incomplete. Brokers: %v, GroupID: %s, Topic: %s",
//			kafkaConfig.KafkaUserServer.Brokers, kafkaConfig.KafkaUserServer.GroupID, kafkaConfig.KafkaUserServer.Topic)
//	}
//
//	// 创建消费者客户端
//	consumer, err := sarama.NewConsumerGroup(kafkaConfig.KafkaUserServer.Brokers, kafkaConfig.KafkaUserServer.GroupID, nil)
//	if err != nil {
//		return fmt.Errorf("failed to create consumer group: %v", err)
//	}
//	defer consumer.Close()
//
//	// 创建消费者处理器
//	handler := &ConsumerHandler{}
//	for {
//		// 消费消息
//		err := consumer.Consume(context.Background(), []string{kafkaConfig.KafkaUserServer.Topic}, handler)
//		if err != nil {
//			return fmt.Errorf("failed to consume messages: %v", err)
//		}
//	}
//}
