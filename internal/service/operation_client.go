package service

import (
	"time"

	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	pb "github.com/russianinvestments/invest-api-go-sdk/proto"
)

type OperationsClient interface {
	GetOperations(req *investgo.GetOperationsRequest) (*investgo.OperationsResponse, error)
	GetPortfolio(accountID string, currency pb.PortfolioRequest_CurrencyRequest) (*investgo.PortfolioResponse, error)
	GetPositions(accountID string) (*investgo.PositionsResponse, error)
	GetWithdrawLimits(accountID string) (*investgo.WithdrawLimitsResponse, error)
	GetBrokerReport(taskID string, page int32) (*investgo.GetBrokerReportResponse, error)
	GenerateBrokerReport(accountID string, from, to time.Time) (*investgo.GenerateBrokerReportResponse, error)
	GetDividentsForeignIssuer(taskID string, page int32) (*investgo.GetDividendsForeignIssuerResponse, error)
	GenerateDividentsForeignIssuer(
		accountID string,
		from, to time.Time,
	) (*investgo.GetDividendsForeignIssuerResponse, error)
	GetOperationsByCursorShort(accountID string) (*investgo.GetOperationsByCursorResponse, error)
	GetOperationsByCursor(req *investgo.GetOperationsByCursorRequest) (*investgo.GetOperationsByCursorResponse, error)
}

var _ OperationsClient = (*investgo.OperationsServiceClient)(nil)
