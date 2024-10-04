package config

import (
	"time"

	"github.com/fishmanDK/price_checker/internal/logger"
)

type (
	Config struct {
		Logger *logger.Config `yaml:"logger"`
		HTTP   HTTP           `yaml:"http"`
		Kafka  Kafka          `yaml:"kafka"`
	}

	HTTP struct {
		Port         string        `yaml:"port"`
		ReadTimeout  time.Duration `yaml:"readTimeout"`
		WriteTimeout time.Duration `yaml:"writeTimeout"`
		// MaxHeaderBytes //TODO
	}

	Kafka struct {
		Brokers    []string `yaml:"brokers"`
		GroupID    string   `yaml:"groupID"`
		InitTopics bool     `yaml:"initTopics"`
	}
)

func InitConfig() (*Config, error) {
	return nil, nil
}
