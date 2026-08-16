package service

import (
	"context"
	"errors"
	"sort"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	repoModels "github.com/ranefattesingh/microservices/user/repository/db/models"
	"github.com/ranefattesingh/microservices/user/service/models"
)

type stubUserRepo struct {
	users      map[int64]*repoModels.User
	byEmail    map[string]*repoModels.User
	createID   int64
	createErr  error
	getByIDErr error
	getAllErr  error
	updateErr  error
}

func newStubUserRepo(users ...*repoModels.User) *stubUserRepo {
	repo := &stubUserRepo{
		users:   make(map[int64]*repoModels.User),
		byEmail: make(map[string]*repoModels.User),
	}

	for _, user := range users {
		repo.users[user.ID] = cloneUser(user)
		repo.byEmail[user.Email] = cloneUser(user)
	}

	return repo
}

func cloneUser(user *repoModels.User) *repoModels.User {
	if user == nil {
		return nil
	}

	copyUser := *user
	return &copyUser
}

func (s *stubUserRepo) Create(ctx context.Context, user *repoModels.User) (int64, error) {
	if s.createErr != nil {
		return 0, s.createErr
	}

	s.createID++
	user.ID = s.createID
	s.users[user.ID] = cloneUser(user)
	s.byEmail[user.Email] = cloneUser(user)

	return user.ID, nil
}

func (s *stubUserRepo) GetByID(ctx context.Context, id int64) (*repoModels.User, error) {
	if s.getByIDErr != nil {
		return nil, s.getByIDErr
	}

	user, ok := s.users[id]
	if !ok {
		return nil, pgx.ErrNoRows
	}

	return cloneUser(user), nil
}

func (s *stubUserRepo) GetByEmail(ctx context.Context, email string) (*repoModels.User, error) {
	user, ok := s.byEmail[email]
	if !ok {
		return nil, pgx.ErrNoRows
	}

	return cloneUser(user), nil
}

func (s *stubUserRepo) Update(ctx context.Context, id int64, user *repoModels.User) error {
	if s.updateErr != nil {
		return s.updateErr
	}

	if _, ok := s.users[id]; !ok {
		return pgx.ErrNoRows
	}

	s.users[id] = cloneUser(user)
	s.byEmail[user.Email] = cloneUser(user)
	return nil
}

func (s *stubUserRepo) GetAll(ctx context.Context, page, limit int) ([]*repoModels.User, int, error) {
	if s.getAllErr != nil {
		return nil, 0, s.getAllErr
	}

	all := make([]*repoModels.User, 0, len(s.users))
	for _, user := range s.users {
		all = append(all, cloneUser(user))
	}

	sort.Slice(all, func(i, j int) bool {
		return all[i].ID < all[j].ID
	})

	if page < 1 {
		page = 1
	}
	if limit <= 0 {
		limit = len(all)
	}

	offset := (page - 1) * limit
	if offset > len(all) {
		return []*repoModels.User{}, len(all), nil
	}

	end := offset + limit
	if end > len(all) {
		end = len(all)
	}

	return all[offset:end], len(all), nil
}

func TestCreateUser_Success(t *testing.T) {
	repo := newStubUserRepo()
	repo.createID = 5
	service := NewUserService(repo)

	user := &models.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Password:  "pass1234",
	}

	id, err := service.CreateUser(context.Background(), user)
	if err != nil {
		t.Fatalf("CreateUser returned unexpected error: %v", err)
	}

	if id != 6 {
		t.Fatalf("expected created user ID 6, got %d", id)
	}

	stored := repo.byEmail["john@example.com"]
	if stored == nil {
		t.Fatal("expected user to be persisted in repository")
	}

	if stored.Password == "pass1234" {
		t.Fatal("expected password to be hashed before storing")
	}

	if stored.Password == "" {
		t.Fatal("expected stored password hash to be non-empty")
	}
}

func TestCreateUser_EmailAlreadyTaken(t *testing.T) {
	now := time.Now()
	repo := newStubUserRepo(&repoModels.User{ID: 1, FirstName: "Jane", LastName: "Doe", Email: "jane@example.com", CreatedAt: now, UpdatedAt: now})
	service := NewUserService(repo)

	_, err := service.CreateUser(context.Background(), &models.User{FirstName: "John", LastName: "Doe", Email: "jane@example.com", Password: "pass1234"})
	if !errors.Is(err, ErrEmailAlreadyTaken) {
		t.Fatalf("expected ErrEmailAlreadyTaken, got %v", err)
	}
}

func TestGetUser_Success(t *testing.T) {
	now := time.Now()
	repo := newStubUserRepo(&repoModels.User{ID: 7, FirstName: "John", LastName: "Doe", Email: "john@example.com", CreatedAt: now, UpdatedAt: now})
	service := NewUserService(repo)

	user, err := service.GetUser(context.Background(), 7)
	if err != nil {
		t.Fatalf("GetUser returned unexpected error: %v", err)
	}

	if user.ID != 7 || user.Email != "john@example.com" || user.FirstName != "John" {
		t.Fatalf("unexpected user returned: %+v", user)
	}
}

func TestGetUser_NotFound(t *testing.T) {
	repo := newStubUserRepo()
	repo.getByIDErr = pgx.ErrNoRows
	service := NewUserService(repo)

	_, err := service.GetUser(context.Background(), 99)
	if !errors.Is(err, ErrUserDoesNotExist) {
		t.Fatalf("expected ErrUserDoesNotExist, got %v", err)
	}
}

func TestGetAllUsers_Success(t *testing.T) {
	now := time.Now()
	repo := newStubUserRepo(
		&repoModels.User{ID: 1, FirstName: "John", LastName: "Doe", Email: "john@example.com", CreatedAt: now, UpdatedAt: now},
		&repoModels.User{ID: 2, FirstName: "Jane", LastName: "Doe", Email: "jane@example.com", CreatedAt: now, UpdatedAt: now},
	)
	service := NewUserService(repo)

	users, total, err := service.GetAllUsers(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("GetAllUsers returned unexpected error: %v", err)
	}

	if len(users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(users))
	}

	if total != 2 {
		t.Fatalf("expected total count 2, got %d", total)
	}
}

func TestGetAllUsers_Pagination(t *testing.T) {
	now := time.Now()
	repo := newStubUserRepo(
		&repoModels.User{ID: 1, FirstName: "John", LastName: "Doe", Email: "john@example.com", CreatedAt: now, UpdatedAt: now},
		&repoModels.User{ID: 2, FirstName: "Jane", LastName: "Doe", Email: "jane@example.com", CreatedAt: now, UpdatedAt: now},
		&repoModels.User{ID: 3, FirstName: "Jake", LastName: "Doe", Email: "jake@example.com", CreatedAt: now, UpdatedAt: now},
	)
	service := NewUserService(repo)

	users, total, err := service.GetAllUsers(context.Background(), 2, 1)
	if err != nil {
		t.Fatalf("GetAllUsers returned unexpected error: %v", err)
	}

	if len(users) != 1 {
		t.Fatalf("expected 1 user in page 2, got %d", len(users))
	}

	if users[0].ID != 2 {
		t.Fatalf("expected first item on page 2 to be user ID 2, got %d", users[0].ID)
	}

	if total != 3 {
		t.Fatalf("expected total count 3, got %d", total)
	}
}

func TestGetAllUsers_RepositoryError(t *testing.T) {
	repo := newStubUserRepo()
	repo.getAllErr = errors.New("db unavailable")
	service := NewUserService(repo)

	_, _, err := service.GetAllUsers(context.Background(), 1, 10)
	if err == nil {
		t.Fatal("expected repository error to be returned")
	}
}

func TestUpdateUser_Success(t *testing.T) {
	now := time.Now()
	repo := newStubUserRepo(&repoModels.User{ID: 3, FirstName: "Old", LastName: "Name", Email: "old@example.com", CreatedAt: now, UpdatedAt: now})
	service := NewUserService(repo)

	err := service.UpdateUser(context.Background(), 3, &models.User{FirstName: "New", LastName: "Name", Email: "new@example.com"})
	if err != nil {
		t.Fatalf("UpdateUser returned unexpected error: %v", err)
	}

	updated := repo.users[3]
	if updated.FirstName != "New" || updated.Email != "new@example.com" {
		t.Fatalf("user was not updated correctly: %+v", updated)
	}
}

func TestUpdateUser_NotFound(t *testing.T) {
	repo := newStubUserRepo()
	repo.getByIDErr = pgx.ErrNoRows
	service := NewUserService(repo)

	err := service.UpdateUser(context.Background(), 55, &models.User{FirstName: "Name", LastName: "Value", Email: "value@example.com"})
	if !errors.Is(err, ErrUserDoesNotExist) {
		t.Fatalf("expected ErrUserDoesNotExist, got %v", err)
	}
}

func TestUpdateUser_EmailAlreadyTaken(t *testing.T) {
	now := time.Now()
	repo := newStubUserRepo(
		&repoModels.User{ID: 10, FirstName: "Someone", LastName: "Else", Email: "taken@example.com", CreatedAt: now, UpdatedAt: now},
		&repoModels.User{ID: 11, FirstName: "Other", LastName: "User", Email: "other@example.com", CreatedAt: now, UpdatedAt: now},
	)
	service := NewUserService(repo)

	err := service.UpdateUser(context.Background(), 11, &models.User{FirstName: "Other", LastName: "User", Email: "taken@example.com"})
	if !errors.Is(err, ErrEmailAlreadyTaken) {
		t.Fatalf("expected ErrEmailAlreadyTaken, got %v", err)
	}
}
