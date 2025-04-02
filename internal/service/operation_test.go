package service

import (
	"context"
	"errors"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/Konstanta100/BrokerCalculator/internal/dto"
	"github.com/Konstanta100/BrokerCalculator/internal/repository"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	pb "github.com/russianinvestments/invest-api-go-sdk/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type MockInvestAPIClient struct {
	OperationsClient
	ctrl *gomock.Controller
}

func NewMockInvestAPIClient(ctrl *gomock.Controller) *MockInvestAPIClient {
	return &MockInvestAPIClient{ctrl: ctrl}
}

func (m *MockInvestAPIClient) EXPECT() *MockInvestAPIClientMockRecorder {
	return &MockInvestAPIClientMockRecorder{mock: m}
}

func (m *MockInvestAPIClient) GetOperationsByCursor(
	req *investgo.GetOperationsByCursorRequest,
) (*investgo.GetOperationsByCursorResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOperationsByCursor", req)
	ret0, _ := ret[0].(*investgo.GetOperationsByCursorResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

type MockInvestAPIClientMockRecorder struct {
	mock *MockInvestAPIClient
}

func (m *MockInvestAPIClientMockRecorder) GetOperationsByCursor(req interface{}) *gomock.Call {
	m.mock.ctrl.T.Helper()
	return m.mock.ctrl.RecordCallWithMethodType(
		m.mock,
		"GetOperationsByCursor",
		reflect.TypeOf((*MockInvestAPIClient)(nil).GetOperationsByCursor),
		req,
	)
}

func TestOperationService_CalculateCommission(t *testing.T) {
	tests := []struct {
		name         string
		operations   []*repository.Operation
		expected     *dto.CalculateCommission
		expectError  bool
		errorMessage string
	}{
		{
			name: "success with single operation",
			operations: []*repository.Operation{
				{
					ID:       "op1",
					Payment:  pgtype.Numeric{Int: big.NewInt(100), Exp: 0, Valid: true},
					Currency: "USD",
					Date: pgtype.Timestamp{
						Time:  time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
				},
			},
			expected: &dto.CalculateCommission{
				DateFrom: "2025-01-01",
				DateTo:   "2025-01-31",
				DateCommissions: map[string]map[string]float64{
					"2025-01-01": {"USD": 100},
				},
				TotalPayments: map[string]float64{
					"USD": 100,
				},
			},
		},
		{
			name: "multiple operations with different currencies",
			operations: []*repository.Operation{
				{
					ID:       "op1",
					Payment:  pgtype.Numeric{Int: big.NewInt(100), Exp: 0, Valid: true},
					Currency: "USD",
					Date: pgtype.Timestamp{
						Time:  time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
				},
				{
					ID:       "op2",
					Payment:  pgtype.Numeric{Int: big.NewInt(5000), Exp: 0, Valid: true},
					Currency: "RUB",
					Date: pgtype.Timestamp{
						Time:  time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
				},
			},
			expected: &dto.CalculateCommission{
				DateFrom: "2025-01-01",
				DateTo:   "2025-01-31",
				DateCommissions: map[string]map[string]float64{
					"2025-01-01": {
						"USD": 100.00,
						"RUB": 5000.00,
					},
				},
				TotalPayments: map[string]float64{
					"USD": 100.00,
					"RUB": 5000.00,
				},
			},
		},
		{
			name: "invalid payment value",
			operations: []*repository.Operation{
				{
					ID:       "op1",
					Payment:  pgtype.Numeric{Valid: false},
					Currency: "USD",
					Date: pgtype.Timestamp{
						Time:  time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
				},
			},
			expected: &dto.CalculateCommission{
				DateFrom:        "2025-01-01",
				DateTo:          "2025-01-31",
				DateCommissions: map[string]map[string]float64{},
				TotalPayments:   map[string]float64{},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			accountID := "test-account"
			from := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
			to := time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC)

			repoMock := repository.NewMockQuerier(ctrl)
			repoMock.EXPECT().OperationsByInstrumentAndDateRange(ctx, gomock.Any()).Return(test.operations, nil)

			service := &OperationService{
				repository: repoMock,
			}

			result, err := service.CalculateCommission(ctx, accountID, from, to)

			if test.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), test.errorMessage)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestOperationService_LoadOperationsFromBroker(t *testing.T) {
	type testCase struct {
		name          string
		prepareMocks  func(*repository.MockQuerier, *MockInvestAPIClient)
		expected      []*repository.Operation
		expectError   bool
		errorContains string
	}

	testDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	testOperation := &repository.Operation{
		ID:             "op1",
		Figi:           "figi1",
		InstrumentType: pb.OperationType_OPERATION_TYPE_BROKER_FEE.String(),
		Quantity:       1,
		Payment:        pgtype.Numeric{Int: big.NewInt(150), Exp: 2, Valid: true}, // 1.50
		Currency:       "USD",
		Date:           pgtype.Timestamp{Time: testDate, Valid: true},
		AccountID:      "test-account",
	}

	tests := []testCase{
		{
			name: "successful load and save",
			prepareMocks: func(repo *repository.MockQuerier, api *MockInvestAPIClient) {
				api.EXPECT().GetOperationsByCursor(gomock.Any()).Return(
					&investgo.GetOperationsByCursorResponse{
						Header: nil,
						GetOperationsByCursorResponse: &pb.GetOperationsByCursorResponse{
							Items: []*pb.OperationItem{
								{
									Id:             "op1",
									Figi:           "figi1",
									Type:           pb.OperationType_OPERATION_TYPE_BROKER_FEE,
									Date:           timestamppb.New(testDate),
									Quantity:       1,
									Price:          &pb.MoneyValue{Currency: "USD"},
									Payment:        &pb.MoneyValue{Currency: "USD", Units: 1, Nano: 50000000},
									InstrumentType: "share",
								},
							},
							HasNext: false,
						},
					},
					nil,
				)

				repo.EXPECT().BulkInsertOperations(gomock.Any(), gomock.Any()).Return(int64(1), nil)
				repo.EXPECT().OperationsByAccountId(
					gomock.Any(),
					"test-account",
				).Return([]*repository.Operation{testOperation}, nil)
			},
			expected:    []*repository.Operation{testOperation},
			expectError: false,
		},
		{
			name: "api error",
			prepareMocks: func(_ *repository.MockQuerier, api *MockInvestAPIClient) {
				api.EXPECT().GetOperationsByCursor(gomock.Any()).Return(nil, errors.New("api error"))
			},
			expectError:   true,
			errorContains: "failed to load operations",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repoMock := repository.NewMockQuerier(ctrl)
			apiMock := NewMockInvestAPIClient(ctrl)
			if test.prepareMocks != nil {
				test.prepareMocks(repoMock, apiMock)
			}

			service := &OperationService{
				operationClient: apiMock,
				repository:      repoMock,
			}

			result, err := service.LoadOperationsFromBroker(context.Background(), "test-account")
			if test.expectError {
				require.Error(t, err)
				if test.errorContains != "" {
					assert.Contains(t, err.Error(), test.errorContains)
				}
				return
			}

			require.NoError(t, err)
			assert.Equal(t, test.expected, result)
		})
	}
}
