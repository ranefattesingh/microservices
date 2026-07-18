package models

import "time"

type AccessType string

type User struct {
	ID         int64      `db:"id"`
	FirstName  string     `db:"first_name"`
	LastName   string     `db:"last_name"`
	Email      string     `db:"email"`
	Phone      string     `db:"phone"`
	AccessType AccessType `db:"access_type"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  time.Time  `db:"updated_at"`
}
