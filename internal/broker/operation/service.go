package operation

import (
	"fmt"
	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	"net/http"
)

type Service struct {
	Client *investgo.OperationsServiceClient
}

func (s *Service) CalculateCommission(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("[INFO request]", request.Host, request.URL.Path, request.Method)
}
