package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/fishmanDK/price_checker/internal/config"
	"github.com/fishmanDK/price_checker/internal/logger"
	"github.com/fishmanDK/price_checker/internal/marketdata/checker"
	"github.com/fishmanDK/price_checker/internal/marketdata/service"
	"github.com/segmentio/kafka-go"
)

type app struct{
	log logger.Logger
	cfg *config.Config
	kafkaConn *kafka.Conn
	// ps *service.ProjectService //TODO
	// metrics   *metrics.WriterServiceMetrics //TODO
	s *service.Service
	c *checker.Checker
}


func NewApp(log logger.Logger, cfg *config.Config) *app {
	return &app{log: log, cfg: cfg}
}

func (a *app) Run() error{
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	// TODO: init storage
	// repo := storage.NewStorage("")
	// a.s = service.NewService(repo)

	a.connectKafkaBrokers(ctx)

	go a.c.Start() //TODO: добавить что-то вроде error в ответе
	// go func() { //TODO: что-то вроде этого
	// 	if err := s.runHttpServer(); err != nil {
	// 		s.log.Errorf(" s.runHttpServer: %v", err)
	// 		cancel()
	// 	}
	// }()

	return nil
}