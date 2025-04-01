package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Konstanta100/BrokerCalculator/internal/dto"
	"github.com/Konstanta100/BrokerCalculator/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	pb "github.com/russianinvestments/invest-api-go-sdk/proto"
)

type OperationService struct {
	operationClient *investgo.OperationsServiceClient
	repository      *repository.Queries
	db              *pgxpool.Pool
}

func NewOperationService(
	operationClient *investgo.OperationsServiceClient,
	repo *repository.Queries,
	db *pgxpool.Pool,
) *OperationService {
	return &OperationService{
		operationClient: operationClient,
		repository:      repo,
		db:              db,
	}
}

func (s *OperationService) CalculateCommission(
	ctx context.Context,
	accountID string,
	dateTimeFrom, dateTimeTo time.Time,
) (*dto.CalculateCommission, error) {
	queryParam := repository.OperationsByInstrumentAndDateRangeParams{
		AccountID:      accountID,
		InstrumentType: pb.OperationType_OPERATION_TYPE_BROKER_FEE.String(),
		Date: pgtype.Timestamp{
			Time:  dateTimeFrom,
			Valid: !dateTimeFrom.IsZero(),
		},
		Date_2: pgtype.Timestamp{
			Time:  dateTimeTo,
			Valid: !dateTimeTo.IsZero(),
		},
	}

	operations, err := s.repository.OperationsByInstrumentAndDateRange(ctx, queryParam)
	if err != nil {
		return nil, fmt.Errorf("failed to get operations: %w", err)
	}

	result := &dto.CalculateCommission{
		DateFrom:        dateTimeFrom.Format(time.DateOnly),
		DateTo:          dateTimeTo.Format(time.DateOnly),
		DateCommissions: make(map[string]map[string]float64),
		TotalPayments:   make(map[string]float64),
	}

	if len(operations) == 0 {
		return result, nil
	}

	for _, op := range operations {
		payment, convertErr := repository.NumericToFloat64(op.Payment)
		if convertErr != nil {
			log.Printf("Invalid payment in operation %s: %v", op.ID, convertErr)
			continue
		}

		dateCommission := op.Date.Time.Format(time.DateOnly)

		if _, exists := result.DateCommissions[dateCommission]; !exists {
			result.DateCommissions[dateCommission] = make(map[string]float64)
		}

		result.DateCommissions[dateCommission][op.Currency] += payment
		result.TotalPayments[op.Currency] += payment
	}

	return result, nil
}

func (s *OperationService) LoadOperationsFromBroker(
	ctx context.Context,
	accountID string,
) ([]*repository.Operation, error) {
	cursorRequest := investgo.GetOperationsByCursorRequest{
		AccountId:          accountID,
		State:              pb.OperationState_OPERATION_STATE_EXECUTED,
		OperationTypes:     []pb.OperationType{pb.OperationType_OPERATION_TYPE_BROKER_FEE},
		WithoutTrades:      true,
		WithoutCommissions: false,
		WithoutOvernights:  true,
	}

	operationListDto, err := s.findAllOperationsByCursor(&cursorRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to load operations: %w", err)
	}

	operationParams := make([]repository.BulkInsertOperationsParams, 0, len(operationListDto))
	for _, operationDto := range operationListDto {
		var payment pgtype.Numeric
		strPayment := fmt.Sprintf("%.2f", operationDto.Payment)
		if err = payment.Scan(strPayment); err != nil {
			return nil, fmt.Errorf("failed to convert payment to numeric: %w", err)
		}

		operationParams = append(operationParams, repository.BulkInsertOperationsParams{
			ID:             operationDto.ID,
			Figi:           operationDto.Figi,
			InstrumentType: operationDto.InstrumentType,
			Quantity:       operationDto.Quantity,
			Payment:        payment,
			Currency:       operationDto.Currency,
			Date:           pgtype.Timestamp{Time: operationDto.Date, Valid: true},
			AccountID:      accountID,
		})
	}

	err = s.insertWithTransaction(ctx, operationParams)
	if err != nil {
		return nil, fmt.Errorf("failed to insert operations: %w", err)
	}

	operations, err := s.repository.OperationsByAccountId(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to load operations: %w", err)
	}

	return operations, nil
}

func (s *OperationService) findAllOperationsByCursor(
	cursorRequest *investgo.GetOperationsByCursorRequest,
) ([]*dto.Operation, error) {
	var (
		allOperations []*dto.Operation
		cursor        string
	)

	for {
		cursorRequest.Cursor = cursor
		operationsResp, err := s.operationClient.GetOperationsByCursor(cursorRequest)
		if err != nil {
			return nil, fmt.Errorf("failed to get operations: %w", err)
		}

		ops := operationsResp.GetItems()
		if len(ops) > 0 {
			for _, op := range ops {
				allOperations = append(allOperations,
					&dto.Operation{
						ID:             op.GetId(),
						Figi:           op.GetFigi(),
						InstrumentType: op.GetType().String(),
						Date:           op.GetDate().AsTime(),
						Quantity:       op.GetQuantity(),
						Currency:       op.GetPrice().GetCurrency(),
						Payment:        op.GetPayment().ToFloat(),
					})
			}
		}

		if !operationsResp.GetHasNext() {
			break
		}
		cursor = operationsResp.GetNextCursor()
	}

	return allOperations, nil
}

func (s *OperationService) insertWithTransaction(
	ctx context.Context,
	operations []repository.BulkInsertOperationsParams,
) error {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead})
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	dbTX := repository.Queries{}
	repo := dbTX.WithTx(tx)

	_, err = repo.BulkInsertOperations(ctx, operations)
	if err != nil {
		return fmt.Errorf("operation not created: %w", err)
	}

	return tx.Commit(ctx)
}
