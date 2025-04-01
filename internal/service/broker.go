package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type BrokerService struct {
	Client *investgo.Client
}

func New(conf investgo.Config) (BrokerService, error) {
	service := BrokerService{}
	client, err := investgo.NewClient(context.Background(), conf, initLogger())
	if err != nil {
		return service, fmt.Errorf("error creating investgo client: %w", err)
	}

	service.Client = client

	return service, nil
}

func initLogger() *zap.SugaredLogger {
	zapConfig := zap.NewDevelopmentConfig()
	zapConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.DateTime)
	zapConfig.EncoderConfig.TimeKey = "time"
	l, err := zapConfig.Build()
	logger := l.Sugar()

	if err != nil {
		log.Fatalf("logger creating error %v", err)
	}
	return logger
}
