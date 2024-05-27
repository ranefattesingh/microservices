package psql

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ranefattesingh/microservices/user/core"
	"github.com/ranefattesingh/pkg/postgresql/pgx/pool"
)

type psqlRepo struct {
	pool *pool.Pool
}

func NewRepo(pool *pool.Pool) *psqlRepo {
	return &psqlRepo{
		pool: pool,
	}
}

// func (p psqlRepo) All(ctx context.Context) (core.Users, error) {
// 	rows, err := p.pool.Connection().Query(ctx, "id, name, email, is_admin, create_date FROM users")
// 	if err != nil {
// 		return nil, err
// 	}

// 	var users core.Users
// 	for rows.Next() {
// 		user := &core.User{}

// 		if err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.IsAdmin, &user.CreateDate); err != nil {
// 			return users, err
// 		}

// 		users = append(users, user)
// 	}

// 	return users, nil
// }

// func (p psqlRepo) ByID(ctx context.Context, id int) (*core.User, error) {

// 	user := &core.User{}

// 	err := p.pool.Connection().QueryRow(ctx, "SELECT id, name, email, is_admin, create_date FROM users WHERE id=$1", id).
// 		Scan(&user.ID, &user.Name, &user.Email, &user.IsAdmin, &user.CreateDate)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return user, nil
// }

func (p psqlRepo) CreateUser(ctx context.Context, user core.User) (uuid.UUID, error) {
	userID := uuid.New()

	_, err := p.pool.Connection().Exec(ctx, createUserQuery, userID, user.Name, user.Email, user.Password, user.IsAdmin)
	if err != nil {
		return uuid.Nil, fmt.Errorf("CreateUser error: %w", err)
	}

	return userID, nil
}

// func (p psqlRepo) Edit(ctx context.Context, user *core.User) error {
// 	return nil
// }

// func (p psqlRepo) Authenticate(ctx context.Context, email, pass string) (int, bool, error) {
// 	return 0, true, nil
// }
