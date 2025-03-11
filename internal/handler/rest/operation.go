package rest

import (
	"fmt"
	"github.com/Konstanta100/BrokerCalculator/internal/service/operation"
	"net/http"
)

type Handler struct {
	OperationService *operation.Service
}

func (s *Handler) CalculateCommission(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("[INFO request]", request.Host, request.URL.Path, request.Method)
}

func (s *Handler) GetOperation(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("[INFO request]", request.Host, request.URL.Path, request.Method)
}
