package operation

import (
	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	pb "github.com/russianinvestments/invest-api-go-sdk/proto"
	"strings"
	"time"
)

type Service struct {
	AccountId  string
	GRPCClient *investgo.OperationsServiceClient
}

func (s *Service) CalculateCommission(dateTimeFrom, dateTimeTo time.Time) (*CalculateCommission, error) {
	cursorRequest := investgo.GetOperationsByCursorRequest{
		AccountId:          s.AccountId,
		From:               dateTimeFrom,
		To:                 dateTimeTo,
		State:              pb.OperationState_OPERATION_STATE_EXECUTED,
		OperationTypes:     []pb.OperationType{pb.OperationType_OPERATION_TYPE_BROKER_FEE},
		WithoutTrades:      true,
		WithoutCommissions: false,
		WithoutOvernights:  true,
		Limit:              int32(1000),
	}
	operations, err := s.findOperationsByCursor(&cursorRequest)
	if err != nil {
		return nil, err
	}

	mapSumPayment := make(map[string]SumPayment)
	mapCommission := make(map[string]Commission)
	if len(*operations) > 0 {
		for _, op := range *operations {
			dateByCommission := op.Time.Format(time.DateOnly)
			keyCom := strings.Join([]string{dateByCommission, op.Currency}, "_")
			commission, ok := mapCommission[keyCom]
			if !ok {
				commission = Commission{op.Currency, 0, dateByCommission}
			}

			commission.Payment += op.Payment
			mapCommission[keyCom] = commission

			sum, ok := mapSumPayment[op.Currency]
			if !ok {
				sum = SumPayment{op.Currency, 0}
			}

			sum.Payment += op.Payment
			mapSumPayment[op.Currency] = sum
		}
	}

	var calculateCommission CalculateCommission
	var commissions Commissions
	var sumPayments []SumPayment

	for _, com := range mapCommission {
		commissions = append(commissions, com)
	}

	for _, sum := range mapSumPayment {
		sumPayments = append(sumPayments, sum)
	}

	calculateCommission.DateFrom = dateTimeFrom.Format(time.DateOnly)
	calculateCommission.DateTo = dateTimeTo.Format(time.DateOnly)
	calculateCommission.Commissions = commissions
	calculateCommission.SumPayment = sumPayments

	return &calculateCommission, nil
}

func (s *Service) GetOperation(figi string, dateTimeFrom, dateTimeTo time.Time) (*Operations, error) {
	return s.findOperations(figi, dateTimeFrom, dateTimeTo)
}

func (s *Service) findOperations(figi string, dateTimeFrom, dateTimeTo time.Time) (*Operations, error) {
	var operations Operations
	operationsResp, err := s.GRPCClient.GetOperations(&investgo.GetOperationsRequest{
		Figi:      figi,
		AccountId: s.AccountId,
		State:     pb.OperationState_OPERATION_STATE_EXECUTED,
		From:      dateTimeFrom,
		To:        dateTimeTo,
	})

	if err != nil {
		return nil, err
	} else {
		ops := operationsResp.GetOperations()
		if len(ops) > 0 {
			for _, op := range ops {
				operations = append(operations,
					Operation{
						Id:             op.GetId(),
						Figi:           op.GetFigi(),
						Description:    op.GetType(),
						InstrumentType: op.GetOperationType().String(),
						Time:           op.GetDate().AsTime(),
						Quantity:       op.GetQuantity(),
						Currency:       op.GetCurrency(),
						Payment:        op.GetPayment().ToFloat(),
					})
			}
		}
	}

	return &operations, nil
}

func (s *Service) findOperationsByCursor(cursorRequest *investgo.GetOperationsByCursorRequest) (*Operations, error) {
	var operations Operations

	operationsResp, err := s.GRPCClient.GetOperationsByCursor(cursorRequest)

	if err != nil {
		return nil, err
	} else {
		ops := operationsResp.GetItems()
		if len(ops) > 0 {
			for _, op := range ops {
				operations = append(operations,
					Operation{
						Id:             op.GetId(),
						Figi:           op.GetFigi(),
						Description:    op.GetDescription(),
						InstrumentType: op.GetType().String(),
						Time:           op.GetDate().AsTime(),
						Quantity:       op.GetQuantity(),
						Currency:       op.GetPrice().GetCurrency(),
						Payment:        op.GetPayment().ToFloat(),
					})
			}
		}
	}

	return &operations, nil
}
