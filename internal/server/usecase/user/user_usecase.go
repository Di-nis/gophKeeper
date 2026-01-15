package user

import (
	"context"

	"github.com/Di-nis/gophKeeper/internal/server/model"
)

// URLRepository - интерфейс для базы данных.
type UserRepository interface {
	Ping(context.Context) error
	Close() error
	InsertUser(context.Context, model.Auth) error
	GetUser(context.Context, model.Auth) error
}

// UserUsecase - какое-то описание.
type UserUsecase struct {
	repo UserRepository
}

// NewUserUsecase - создание структуры NewUserUsecase.
func NewUserUsecase(repo UserRepository) *UserUsecase {
	return &UserUsecase{
		repo: repo,
	}
}

// Ping - проверка соединения с базой данных.
func (u *UserUsecase) Ping(ctx context.Context) error {
	return u.repo.Ping(ctx)
}

// CreateUser - создание пользователя.
func (u *UserUsecase) CreateUser(ctx context.Context, user model.Auth) error {
	return nil
}

// GetUser - получение пользователя.
func (u *UserUsecase) GetUser(ctx context.Context, user model.Auth) error {
	return nil
}