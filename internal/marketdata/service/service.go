package service

import (
	"log"

	"github.com/fishmanDK/price_checker/internal/kafka"
	"github.com/fishmanDK/price_checker/internal/marketdata/storage"
)

type Service struct {
	storage    storage.Storage
	kafkaProducer *kafka.KafkaProducer
}

func NewService(storage storage.Storage) *Service {
	kafkaProducer, err := kafka.NewKafkaProducer("localhost:29092", "ratios-info")
	if err != nil{
		log.Fatal(err.Error())
	}

	return &Service{
		storage:    storage,
		kafkaProducer: kafkaProducer,
	}
}
