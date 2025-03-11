package build

import (
	"fmt"
	"github.com/Konstanta100/BrokerCalculator/internal/restapi"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"golang.org/x/net/context"
	"log"
	"net"
	"net/http"
	"time"
)

func (b *Builder) RestAPIServer() (*http.Server, error) {
	server, err := b.httpServer()
	if err != nil {
		return nil, fmt.Errorf("could not create http server: %w", err)
	}

	err = b.registerHandlers()
	if err != nil {
		return nil, fmt.Errorf("could not register handlers: %w", err)
	}

	return server, nil
}

func (b *Builder) httpRouter() *mux.Router {
	if b.router != nil {
		return b.router
	}

	b.router = mux.NewRouter()

	return b.router
}

func (b *Builder) httpServer() (*http.Server, error) {
	const timeout = time.Millisecond * 25

	server := http.Server{
		Addr:              b.config.HTTPAddr(),
		ReadHeaderTimeout: timeout,
		Handler:           b.httpRouter(),
		ErrorLog:          log.New(zerolog.Nop(), "", 0),
		BaseContext: func(net.Listener) context.Context {
			return context.Background()
		},
	}

	return &server, nil
}

func (b *Builder) registerHandlers() error {
	router := b.httpRouter()
	apiRoute := router.PathPrefix("/api").Subrouter()

	server, err := restapi.New(b.config)

	if err != nil {
		return fmt.Errorf("could not create rest api server: %w", err)
	}

	apiRoute.HandleFunc("/operations/commission", server.OperationHandler.CalculateCommission).Methods("GET")
	apiRoute.HandleFunc("/operations", server.OperationHandler.GetOperation).Methods("GET")

	return nil
}
