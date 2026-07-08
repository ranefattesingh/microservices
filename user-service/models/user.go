package models

type CreateUser struct {
	ID              string
	FirstName       string
	LastName        string
	Email           string
	Password        string
	ConfirmPassword string
}
