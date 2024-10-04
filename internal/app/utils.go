package app

import (
	"context"

	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
)

func (s *app) connectKafkaBrokers(ctx context.Context) error {
	kafkaConn, err := kafka.DialContext(ctx, "tcp", s.cfg.Kafka.Brokers[0])
	if err != nil{
		return errors.Wrap(err, "kafka.DialContext")
	}

	s.kafkaConn = kafkaConn

	brokers, err := kafkaConn.Brokers()
	if err != nil{
		return errors.Wrap(err, "kafkaConn.Brokers")
	}

	s.log.Infof("kafka connected to brokers: %+v", brokers)

	return nil
}