package logger

import (
	"log/slog"
	"os"
	"time"
)

type Config struct {
	LogLevel string `yaml:"level"`
	DevMode  bool   `yaml:"devMode"`
	Encoder  string `yaml:"encoder"`
}


type Logger interface {
	InitLogger()
	Sync() error
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Warn(args ...interface{})
	Warnf(template string, args ...interface{})
	WarnMsg(msg string, err error)
	Error(args ...interface{})
	Errorf(template string, args ...interface{})
	Err(msg string, err error)
	DPanic(args ...interface{})
	DPanicf(template string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})
	Printf(template string, args ...interface{})
	WithName(name string)
	HttpMiddlewareAccessLogger(method string, uri string, status int, size int64, time time.Duration)
	GrpcMiddlewareAccessLogger(method string, time time.Duration, metaData map[string][]string, err error)
	GrpcClientInterceptorLogger(method string, req interface{}, reply interface{}, time time.Duration, metaData map[string][]string, err error)
	KafkaProcessMessage(topic string, partition int, message string, workerID int, offset int64, time time.Time)
	KafkaLogCommittedMessage(topic string, partition int, offset int64)
}

type appLogger struct {
	level       string
	devMode     bool
	encoding    string
	// sugarLogger *zap.SugaredLogger //FIXME
	logger      *slog.Logger
}

func NewAppLogger(cfg *Config) *appLogger{
	return &appLogger{level: cfg.LogLevel, devMode: cfg.DevMode, encoding: cfg.Encoder}
}

func (l *appLogger) InitLogger() {
	var opts slog.HandlerOptions
	if l.devMode{
		opts = slog.HandlerOptions{Level: slog.LevelDebug}
	} else {
		opts = slog.HandlerOptions{Level: slog.LevelInfo}
	}

	var slogHandler slog.Handler
	if l.encoding == "console"{
		slogHandler = slog.NewTextHandler(os.Stdout, &opts)
	} else {
		slogHandler = slog.NewJSONHandler(os.Stdout, &opts)
	}

	l.logger = slog.New(slogHandler)
}

func (l *appLogger) Sync() error{return nil}

func (l *appLogger) Debug(args ...interface{}){}

func (l *appLogger) Debugf(template string, args ...interface{}){}

func (l *appLogger) Info(args ...interface{}){}

func (l *appLogger) Infof(template string, args ...interface{}){}

func (l *appLogger) Warn(args ...interface{}){}

func (l *appLogger) Warnf(template string, args ...interface{}){}

func (l *appLogger) WarnMsg(msg string, err error){}

func (l *appLogger) Error(args ...interface{}){}

func (l *appLogger) Errorf(template string, args ...interface{}){}

func (l *appLogger) Err(msg string, err error){}

func (l *appLogger) DPanic(args ...interface{}){}

func (l *appLogger) DPanicf(template string, args ...interface{}){}

func (l *appLogger) Fatal(args ...interface{}){}

func (l *appLogger) Fatalf(template string, args ...interface{}){}

func (l *appLogger) Printf(template string, args ...interface{}){}

func (l *appLogger) WithName(name string){}

func (l *appLogger) HttpMiddlewareAccessLogger(method string, uri string, status int, size int64, time time.Duration){}

func (l *appLogger) GrpcMiddlewareAccessLogger(method string, time time.Duration, metaData map[string][]string, err error){}

func (l *appLogger) GrpcClientInterceptorLogger(method string, req interface{}, reply interface{}, time time.Duration, metaData map[string][]string, err error){}

func (l *appLogger) KafkaProcessMessage(topic string, partition int, message string, workerID int, offset int64, time time.Time){}

func (l *appLogger) KafkaLogCommittedMessage(topic string, partition int, offset int64){}
