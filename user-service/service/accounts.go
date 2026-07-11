package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/ranefattesingh/ecommerce-platform/user-service/repository"
	"github.com/ranefattesingh/ecommerce-platform/user-service/service/models"
)

var ErrAccountAlreadyExists = errors.New("account with the email already exists")
var ErrPublicIDGenerationFail = errors.New("fail to generate publicID")

type AccountsService interface {
	CreateAccount(ctx context.Context, account models.Account) (uuid.UUID, error)
}

type accountsService struct {
	accountsRepository repository.AccountsRepository
}

func NewAccountsService(ar repository.AccountsRepository) *accountsService {
	return &accountsService{
		accountsRepository: ar,
	}
}

func (ac *accountsService) CreateAccount(ctx context.Context, account models.Account) (uuid.UUID, error) {
	// Verify whether account already exists.
	_, err := ac.accountsRepository.GetAccountWithEmail(ctx, account.Email)
	if err == nil {
		return uuid.Nil, ErrAccountAlreadyExists
	}

	if err != pgx.ErrNoRows {
		return uuid.Nil, fmt.Errorf("service:{%w}", err)
	}

	// Generate and assign public id for the account.
	publicID, err := uuid.NewV7()
	if err != nil {
		return uuid.Nil, ErrPublicIDGenerationFail
	}

	account.PublicID = publicID

	// Create account record in database.
	_, err = ac.CreateAccount(ctx, account)
	if err != nil {
		return uuid.Nil, fmt.Errorf("service:{%w}", err)
	}

	return account.PublicID, nil
}
