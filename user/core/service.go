package core

import (
	"context"
	"net/http"
	"time"

	"github.com/ranefattesingh/microservices/user/models"
	"github.com/ranefattesingh/pkg/json"
)

var (
	ErrUserDoesNotExist = &json.Error{
		HTTPStatusCode: http.StatusNotFound,
		Code:           1,
		Message:        "user does not exist",
	}
)

var userX = models.User{
	Id:         0,
	Name:       "UserName",
	Email:      "username@email.com",
	Password:   "UserPassword",
	IsAdmin:    false,
	CreateDate: time.Now().String(),
	UpdateDate: time.Now().String(),
}

var users models.Users

type UserService interface {
	Login(ctx context.Context, login *models.Login) (string, error)
	Authorize(s string) (map[string]interface{}, error)
	AddUser(ctx context.Context, user *models.User) (int, error)
	GetUserById(ctx context.Context, id int) (*models.User, error)
	EditUser(ctx context.Context, user *models.User) error
	GetAllUsers(ctx context.Context) (models.Users, error)
}

type userService struct {
}

func NewUserService() *userService {
	return &userService{}
}

func (u userService) Authorize(s string) (map[string]interface{}, error) {
	return nil, nil
}

func (u userService) GetAllUsers(ctx context.Context) (models.Users, error) {
	return users, nil
}

func (u userService) GetUserById(ctx context.Context, id int) (*models.User, error) {
	return &userX, nil
}

func (u userService) AddUser(ctx context.Context, user *models.User) (id int, err error) {
	user.Id++
	users = append(users, user)
	return user.Id, nil
}

func (u userService) EditUser(ctx context.Context, user *models.User) error {
	return ErrUserDoesNotExist
}

func (u userService) Login(ctx context.Context, login *models.Login) (string, error) {
	return "loggedin", nil
}
