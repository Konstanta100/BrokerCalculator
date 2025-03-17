package account

import (
	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	pb "github.com/russianinvestments/invest-api-go-sdk/proto"
)

type Service struct {
	AccountId  string
	GRPCClient *investgo.UsersServiceClient
}

func (s *Service) GetAccounts(status string) (*[]Account, error) {
	accountStatus, ok := pb.AccountStatus_value[status]
	if !ok {
		accountStatus = int32(pb.AccountStatus_ACCOUNT_STATUS_ALL)
	}

	accountStatusEnum := pb.AccountStatus(accountStatus)

	var accounts []Account
	accsResp, err := s.GRPCClient.GetAccounts(&accountStatusEnum)

	if err != nil {
		return nil, err
	} else {
		accs := accsResp.GetAccounts()
		for _, acc := range accs {
			accounts = append(accounts,
				Account{
					Id:          acc.GetId(),
					Type:        acc.GetType().String(),
					Name:        acc.GetName(),
					Status:      acc.GetStatus().String(),
					OpenedDate:  acc.OpenedDate.AsTime(),
					ClosedDate:  acc.ClosedDate.AsTime(),
					AccessLevel: acc.GetAccessLevel().String(),
				})
		}
	}

	return &accounts, nil
}
