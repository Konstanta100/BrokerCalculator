package service

import (
	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	pb "github.com/russianinvestments/invest-api-go-sdk/proto"
)

type AccountClient interface {
	GetAccounts(status *pb.AccountStatus) (*investgo.GetAccountsResponse, error)
	GetMarginAttributes(accountID string) (*investgo.GetMarginAttributesResponse, error)
	GetUserTariff() (*investgo.GetUserTariffResponse, error)
	GetInfo() (*investgo.GetInfoResponse, error)
}

var _ AccountClient = (*investgo.UsersServiceClient)(nil)
