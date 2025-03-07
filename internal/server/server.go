package server

import (
	"context"
	"fmt"
	"github.com/Konstanta100/BrokerCalculator/internal/config"
	"github.com/Konstanta100/BrokerCalculator/internal/operation"
	"github.com/gorilla/mux"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

type Server struct {
	conf config.Config
}

func RestAPIServer(conf config.Config) (*http.Server, error) {
	const timeout = time.Millisecond * 25
	router := mux.NewRouter()

	server := http.Server{
		Addr:              conf.HTTPAddr(),
		ReadHeaderTimeout: timeout,
		Handler:           router,
		ErrorLog:          log.New(os.Stderr, "", 0),
		BaseContext: func(net.Listener) context.Context {
			return context.Background()
		},
	}

	operationHandler, err := operation.NewService(conf)

	if err != nil {
		return nil, fmt.Errorf("could not create operations handler: %f", err)
	}

	apiRoute := router.PathPrefix("/api").Subrouter()
	apiRoute.HandleFunc("/operations/commission", operationHandler.CalculateCommission).Methods("POST")

	return &server, nil
}
