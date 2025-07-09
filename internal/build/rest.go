package build

import (
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/Konstanta100/BrokerCalculator/internal/restapi"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"golang.org/x/net/context"
)

func (b *Builder) RestAPIServer(ctx context.Context) (*http.Server, error) {
	server := b.httpServer(ctx)
	db, err := b.NewPostgresDB(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not create postgres db: %w", err)
	}

	err = b.registerHandlers(db)
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

func (b *Builder) httpServer(ctx context.Context) *http.Server {
	server := &http.Server{
		Addr:              net.JoinHostPort(b.config.HTTP.Host, b.config.HTTP.Port),
		ReadHeaderTimeout: time.Millisecond * 5,
		Handler:           b.httpRouter(),
		ErrorLog:          log.New(zerolog.Nop(), "", 0),
		BaseContext: func(net.Listener) context.Context {
			return ctx
		},
	}

	return server
}

func (b *Builder) registerHandlers(db *pgxpool.Pool) error {
	router := b.httpRouter()
	apiRoute := router.PathPrefix("/api").Subrouter()

	server, err := restapi.New(b.config, db)
	if err != nil {
		return fmt.Errorf("could not create rest api server: %w", err)
	}

	operationHandler := server.OperationHandler
	accountHandler := server.AccountHandler
	userHandler := server.UserHandler
	router.HandleFunc("/", b.homeHandler).Methods(http.MethodGet)
	router.HandleFunc("/broker/calculator", b.brokerCalculator).Methods(http.MethodGet)
	apiRoute.HandleFunc("/operations/commission", operationHandler.CommissionFromBroker).Methods(http.MethodPost)
	apiRoute.HandleFunc("/operations/load", operationHandler.LoadOperations).Methods(http.MethodPost)
	apiRoute.HandleFunc("/accounts/load", accountHandler.LoadAccounts).Methods(http.MethodGet)
	apiRoute.HandleFunc("/accounts", accountHandler.Accounts).Methods(http.MethodGet)
	apiRoute.HandleFunc("/account", accountHandler.Account).Methods(http.MethodGet)
	apiRoute.HandleFunc("/user", userHandler.User).Methods(http.MethodGet)
	apiRoute.HandleFunc("/user/create", userHandler.CreateUser).Methods(http.MethodPost)

	return nil
}

func (b *Builder) homeHandler(w http.ResponseWriter, r *http.Request) {
	b.renderTemplate(w, "main.html", nil)
}
func (b *Builder) brokerCalculator(w http.ResponseWriter, r *http.Request) {
	b.renderTemplate(w, "broker_calculator.html", nil)
}

func (b *Builder) renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	t, err := template.ParseFiles("app/templates/" + tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
