package restapi

import (
	"fmt"
	"github.com/Konstanta100/BrokerCalculator/internal/broker"
	"github.com/Konstanta100/BrokerCalculator/internal/broker/operation"
	"github.com/Konstanta100/BrokerCalculator/internal/config"
)

type Server struct {
	conf             config.Config
	BrokerService    broker.Service
	OperationService operation.Service
}

func New(conf config.Config) (*Server, error) {
	server := Server{conf: conf}
	brokerService, err := broker.New(conf.BrokerConfig)

	if err != nil {
		return nil, fmt.Errorf("could not create broker service: %w", err)
	}

	operationService := operation.Service{Client: brokerService.Client.NewOperationsServiceClient()}

	server.OperationService = operationService

	return &server, err
}
