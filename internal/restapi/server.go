package restapi

import (
	"fmt"

	"github.com/Konstanta100/BrokerCalculator/internal/config"
	"github.com/Konstanta100/BrokerCalculator/internal/handler/rest"
	"github.com/Konstanta100/BrokerCalculator/internal/repository"
	"github.com/Konstanta100/BrokerCalculator/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	conf             *config.Config
	OperationHandler rest.OperationHandler
	AccountHandler   rest.AccountHandler
	UserHandler      rest.UserHandler
}

func New(conf *config.Config, db *pgxpool.Pool) (*Server, error) {
	server := Server{conf: conf}
	brokerService, err := service.New(conf.BrokerConfig)
	if err != nil {
		return nil, fmt.Errorf("could not create broker service: %w", err)
	}

	repo := repository.New(db)
	operationClient := brokerService.Client.NewOperationsServiceClient()
	accountClient := brokerService.Client.NewUsersServiceClient()

	operationService := service.NewOperationService(operationClient, repo, db)
	accountService := service.NewAccountService(accountClient, repo, db)
	userService := service.NewUserService(repo)

	operationHandler := rest.OperationHandler{OperationService: operationService}
	accountHandler := rest.AccountHandler{AccountService: accountService}
	userHandler := rest.UserHandler{UserService: userService}

	server.OperationHandler = operationHandler
	server.AccountHandler = accountHandler
	server.UserHandler = userHandler

	return &server, err
}
