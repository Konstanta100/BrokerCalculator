package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Konstanta100/BrokerCalculator/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	pb "github.com/russianinvestments/invest-api-go-sdk/proto"
)

var ErrInvalidUserID = errors.New("invalid user ID")

type AccountService struct {
	accountClient AccountClient
	repository    repository.Querier
	db            *pgxpool.Pool
}

func NewAccountService(
	accountClient AccountClient,
	repo repository.Querier,
	db *pgxpool.Pool,
) *AccountService {
	return &AccountService{
		accountClient: accountClient,
		repository:    repo,
		db:            db,
	}
}

func (s *AccountService) LoadAccountsFromBroker(
	ctx context.Context,
	userID pgtype.UUID,
) ([]*repository.Account, error) {
	user, err := s.repository.UserById(ctx, userID)
	if err != nil {
		return nil, err
	}

	accountStatusEnum := pb.AccountStatus_ACCOUNT_STATUS_ALL
	accountResp, err := s.accountClient.GetAccounts(&accountStatusEnum)
	if err != nil {
		return nil, err
	}

	accs := accountResp.GetAccounts()
	accountParamList := make([]repository.AccountCreateParams, 0, len(accs))
	for _, acc := range accs {
		accountParamList = append(accountParamList,
			repository.AccountCreateParams{
				ID:          acc.GetId(),
				UserID:      user.ID,
				Type:        acc.GetType().String(),
				Name:        acc.GetName(),
				Status:      acc.GetStatus().String(),
				OpenedDate:  repository.ConvertProtoTimestampToPgType(acc.OpenedDate),
				ClosedDate:  repository.ConvertProtoTimestampToPgType(acc.ClosedDate),
				AccessLevel: acc.GetAccessLevel().String(),
			})
	}

	err = s.insertWithTransaction(ctx, accountParamList)
	if err != nil {
		return nil, fmt.Errorf("failed to insert accounts: %w", err)
	}

	return s.FindAccounts(ctx, userID)
}

func (s *AccountService) FindByID(ctx context.Context, accountID string) (*repository.Account, error) {
	account, err := s.repository.AccountById(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	return account, nil
}

func (s *AccountService) FindAccounts(ctx context.Context, userID pgtype.UUID) ([]*repository.Account, error) {
	if !userID.Valid {
		return nil, ErrInvalidUserID
	}

	accounts, err := s.repository.AccountsByUserId(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get accounts: %w", err)
	}

	return accounts, nil
}

func (s *AccountService) insertWithTransaction(ctx context.Context, paramList []repository.AccountCreateParams) error {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	dbTX := repository.Queries{}
	repo := dbTX.WithTx(tx)
	for _, accountParam := range paramList {
		_, err = repo.AccountCreate(ctx, accountParam)
		if err != nil {
			return errors.New("accounts not created")
		}
	}

	return tx.Commit(ctx)
}
