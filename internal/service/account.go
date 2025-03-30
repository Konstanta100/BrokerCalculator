package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/Konstanta100/BrokerCalculator/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/russianinvestments/invest-api-go-sdk/investgo"
	pb "github.com/russianinvestments/invest-api-go-sdk/proto"
)

type AccountService struct {
	accountClient *investgo.UsersServiceClient
	repository    *repository.Queries
	db            *pgxpool.Pool
}

func NewAccountService(accountClient *investgo.UsersServiceClient, repo *repository.Queries, db *pgxpool.Pool) *AccountService {
	return &AccountService{
		accountClient: accountClient,
		repository:    repo,
		db:            db,
	}
}

func (s *AccountService) LoadAccountsFromBroker(ctx context.Context, userID pgtype.UUID) ([]*repository.Account, error) {
	user, err := s.repository.UserById(ctx, userID)
	if err != nil {
		fmt.Println(err)
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

	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.RepeatableRead})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	dbTX := repository.Queries{}
	repo := dbTX.WithTx(tx)
	for _, accountParam := range accountParamList {
		_, err = repo.AccountCreate(ctx, accountParam)
		if err != nil {
			return nil, errors.New("accounts not created")
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, errors.New("accounts not created")
	}

	return s.FindAccounts(ctx, userID)
}

func (s *AccountService) FindById(ctx context.Context, id string) (*repository.Account, error) {
	account, err := s.repository.AccountById(ctx, id)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (s *AccountService) FindAccounts(ctx context.Context, userID pgtype.UUID) ([]*repository.Account, error) {
	accounts, err := s.repository.AccountsByUserId(ctx, userID)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}
