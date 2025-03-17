package restapi

import (
	"fmt"
	"github.com/Konstanta100/BrokerCalculator/internal/broker"
	"github.com/Konstanta100/BrokerCalculator/internal/config"
	"github.com/Konstanta100/BrokerCalculator/internal/handler/rest"
)

type Server struct {
	conf             config.Config
	brokerService    broker.Service
	OperationHandler rest.OperationHandler
	AccountHandler   rest.AccountHandler
}

func New(conf config.Config) (*Server, error) {
	server := Server{conf: conf}
	brokerService, err := broker.New(conf.BrokerConfig)

	if err != nil {
		return nil, fmt.Errorf("could not create broker service: %w", err)
	}

	operationHandler := rest.OperationHandler{OperationService: brokerService.NewOperationService()}
	accountHandler := rest.AccountHandler{AccountService: brokerService.NewAccountService()}

	server.OperationHandler = operationHandler
	server.AccountHandler = accountHandler

	return &server, err
}
