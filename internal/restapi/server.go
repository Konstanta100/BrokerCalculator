package restapi

import (
	"fmt"
	"github.com/Konstanta100/BrokerCalculator/internal/config"
	"github.com/Konstanta100/BrokerCalculator/internal/handler/rest"
	"github.com/Konstanta100/BrokerCalculator/internal/service/broker"
	"github.com/Konstanta100/BrokerCalculator/internal/service/operation"
)

type Server struct {
	conf             config.Config
	brokerService    broker.Service
	OperationHandler rest.Handler
}

func New(conf config.Config) (*Server, error) {
	server := Server{conf: conf}
	brokerService, err := broker.New(conf.BrokerConfig)

	if err != nil {
		return nil, fmt.Errorf("could not create broker service: %w", err)
	}

	operationHandler := rest.Handler{
		OperationService: operation.NewOperationService(brokerService.Client.NewOperationsServiceClient()),
	}

	server.OperationHandler = operationHandler

	return &server, err
}
