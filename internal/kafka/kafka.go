package kafka

import (
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaProducer struct {
	Producer *kafka.Producer
	Topic    string
}

func NewKafkaProducer(brokers, topic string) (*KafkaProducer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to create Kafka producer: %v", err)
	}
	return &KafkaProducer{Producer: p, Topic: topic}, nil
}

func (kp *KafkaProducer) PublishMessage(symb string, ratio float64, duration string) error {
	data := map[string]interface{}{
		"symbol":   symb,
		"ratio":    fmt.Sprintf("%.2f%%", ratio),
		"duration": duration,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("Failed to marshal data to JSON: %v", err)
	}

	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &kp.Topic, Partition: kafka.PartitionAny},
		Value:          jsonData,
		Key:            []byte(symb),
	}

	err = kp.Producer.Produce(message, nil)
	if err != nil {
		return fmt.Errorf("Failed to publish message to Kafka: %v", err)
	}

	kp.Producer.Flush(15 * 1000)
	return nil
}