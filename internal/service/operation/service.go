package operation

import "github.com/russianinvestments/invest-api-go-sdk/investgo"

type Service struct {
	GRPCClient *investgo.OperationsServiceClient
}

func NewOperationService(grpcClient *investgo.OperationsServiceClient) *Service {
	return &Service{GRPCClient: grpcClient}
}
