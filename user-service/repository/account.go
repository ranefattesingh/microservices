package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ranefattesingh/ecommerce-platform/user-service/repository/postgres/accounts"
	"github.com/ranefattesingh/ecommerce-platform/user-service/repository/postgres/connection"
	"github.com/ranefattesingh/ecommerce-platform/user-service/service/models"
)

type AccountsRepository interface {
	GetAccountWithEmail(ctx context.Context, email string) (models.Account, error)
	CreateAccount(ctx context.Context, account models.Account) error
}

type accountsRepository struct {
	queries accounts.Queries
}

func NewAccountsRepository(c connection.Connection) *accountsRepository {
	return &accountsRepository{
		queries: *accounts.New(c.Pool()),
	}
}

func (ac *accountsRepository) CreateAccount(ctx context.Context, account models.Account) error {
	createAccountParams := accounts.CreateAccountParams{
		PublicID: pgtype.UUID{
			Bytes: account.PublicID,
			Valid: true,
		},
		Email:        account.Email,
		PasswordHash: account.Password,
	}

	_, err := ac.queries.CreateAccount(ctx, createAccountParams)
	if err != nil {
		return fmt.Errorf("repository: {%w}", err)
	}

	return nil
}

func (ac *accountsRepository) GetAccountWithEmail(ctx context.Context, email string) (models.Account, error) {
	account, err := ac.queries.GetAccountWithEmail(ctx, email)
	if err != nil {
		return models.Account{}, fmt.Errorf("repository: {%w}", err)
	}

	accountModel := models.Account{
		ID:        account.ID,
		PublicID:  account.PublicID.Bytes,
		Email:     account.Email,
		CreatedAt: account.CreatedAt.Time,
		UpdatedAt: account.UpdatedAt.Time,
	}

	return accountModel, nil
}
