package service

import (
	"context"
	"testing"
	"time"

	"github.com/Konstanta100/BrokerCalculator/internal/repository"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAccountService_FindByID(t *testing.T) {
	tests := []struct {
		name          string
		accountID     string
		account       *repository.Account
		expectedError string
	}{
		{
			name:      "successful find account",
			accountID: "account123",
			account: &repository.Account{
				ID: "account123",
				UserID: pgtype.UUID{
					Bytes: [16]byte{0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xA, 0xB, 0xC, 0xD, 0xE, 0xF, 0x10},
					Valid: true,
				},
				Name:        "Основной брокерский счет",
				Status:      "ACCOUNT_STATUS_OPEN",
				Type:        "ACCOUNT_TYPE_TINKOFF",
				AccessLevel: "ACCOUNT_ACCESS_LEVEL_FULL_ACCESS",
				OpenedDate: pgtype.Timestamp{
					Time:  time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC),
					Valid: true,
				},
				ClosedDate: pgtype.Timestamp{
					Valid: false,
				},
			},
		},
		{
			name:      "account not found",
			accountID: "acc2",
			account:   nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			repoMock := repository.NewMockQuerier(ctrl)
			repoMock.EXPECT().AccountById(ctx, gomock.Any()).Return(test.account, nil)

			service := NewAccountService(nil, repoMock, nil)
			account, _ := service.FindByID(context.Background(), test.accountID)

			assert.Equal(t, test.account, account)
		})
	}
}

func TestAccountService_FindAccounts(t *testing.T) {
	validUserID := pgtype.UUID{
		Bytes: [16]byte{0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xA, 0xB, 0xC, 0xD, 0xE, 0xF, 0x10},
		Valid: true,
	}

	testAccounts := []*repository.Account{
		{
			ID:     "account1",
			UserID: validUserID,
			Name:   "Основной счет",
			Status: "ACCOUNT_STATUS_OPEN",
		},
		{
			ID:     "account2",
			UserID: validUserID,
			Name:   "Дополнительный счет",
			Status: "ACCOUNT_STATUS_OPEN",
		},
	}

	tests := []struct {
		name          string
		userID        pgtype.UUID
		accounts      []*repository.Account
		expectedError error
	}{
		{
			name:     "successful find accounts",
			userID:   validUserID,
			accounts: testAccounts,
		},
		{
			name:     "no accounts found",
			userID:   validUserID,
			accounts: []*repository.Account{},
		},
		{
			name:          "invalid user ID",
			userID:        pgtype.UUID{Valid: false},
			expectedError: ErrInvalidUserID,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			repoMock := repository.NewMockQuerier(ctrl)

			if test.userID.Valid {
				repoMock.EXPECT().
					AccountsByUserId(ctx, test.userID).
					Return(test.accounts, nil)
			}

			service := NewAccountService(nil, repoMock, nil)
			accounts, err := service.FindAccounts(ctx, test.userID)

			if test.expectedError != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, test.expectedError)
				assert.Nil(t, accounts)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.accounts, accounts)
			}
		})
	}
}
