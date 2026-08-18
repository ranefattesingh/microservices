package models

import (
	"time"

	"github.com/ranefattesingh/microservices/auth/repository/db/models"
)

type Account struct {
	ID        int64
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (a *Account) ToRepositoryModel() *models.Account {
	return &models.Account{
		ID:        a.ID,
		Email:     a.Email,
		Password:  a.Password,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}
