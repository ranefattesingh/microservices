package models

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID        int32
	PublicID  uuid.UUID
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
