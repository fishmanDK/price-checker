package kafka

import (
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/compress"
)

func NewWriter(brokers []string, errLogger kafka.Logger) *kafka.Writer{
	return &kafka.Writer{
		Addr: kafka.TCP(brokers...),
		Balancer: &kafka.LeastBytes{},
		RequiredAcks: writerRequiredAcks, //FIXME
		MaxAttempts:  writerMaxAttempts, //FIXME
		ErrorLogger:  errLogger, //FIXME
		Compression:  compress.Snappy, //FIXME
		ReadTimeout:  writerReadTimeout, //FIXME
		WriteTimeout: writerWriteTimeout, //FIXME
		Async:        false, //FIXME
	}
}