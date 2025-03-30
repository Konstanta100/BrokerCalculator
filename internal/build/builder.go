package build

import (
	"github.com/Konstanta100/BrokerCalculator/internal/config"
	"github.com/gorilla/mux"
	"net/http"
)

type Builder struct {
	config *config.Config
	router *mux.Router
	server *http.Server
}

func New(conf *config.Config) *Builder {
	b := Builder{config: conf}
	return &b
}
