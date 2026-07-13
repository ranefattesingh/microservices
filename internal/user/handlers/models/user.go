package models

type AccessType string

type User struct {
	ID         int32      `json:"id"`
	FirstName  string     `json:"firstName"`
	LastName   string     `json:"lastName"`
	Email      string     `json:"email"`
	Phone      string     `json:"phone"`
	AccessType AccessType `json:"accessType"`
	CreatedAt  string     `json:"createdAt"`
	UpdatedAt  string     `json:"updatedAt"`
}
