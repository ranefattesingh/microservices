package models

type User struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	IsAdmin    bool   `json:"is_admin"`
	CreateDate string `json:"created_date"`
	UpdateDate string `json:"updated_date"`
}

type Users []*User
