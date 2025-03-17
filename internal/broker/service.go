package broker

import (
	"context"
	"fmt"
	"github.com/Konstanta100/BrokerCalculator/internal/operation"
	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"time"
)

type Service struct {
	Client *investgo.Client
}

func New(conf investgo.Config) (Service, error) {
	service := Service{}

	client, err := investgo.NewClient(context.Background(), conf, initLogger())
	if err != nil {
		return service, fmt.Errorf("error creating investgo client: %w", err)
	}

	service.Client = client

	return service, nil
}

func (s *Service) NewOperationService() *operation.Service {
	return &operation.Service{
		AccountId:  s.Client.Config.AccountId,
		GRPCClient: s.Client.NewOperationsServiceClient(),
	}
}

func initLogger() *zap.SugaredLogger {
	zapConfig := zap.NewDevelopmentConfig()
	zapConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.DateTime)
	zapConfig.EncoderConfig.TimeKey = "time"
	l, err := zapConfig.Build()
	logger := l.Sugar()

	defer func() {
		err = logger.Sync()
		if err != nil {
			log.Printf(err.Error())
		}
	}()

	if err != nil {
		log.Fatalf("logger creating error %v", err)
	}
	return logger
}
