package user

import (
	"context"
	"errors"

	"github.com/Di-nis/gophKeeper/internal/model"
	crpt "github.com/Di-nis/gophKeeper/internal/server/crypto"
)

const (
	// passwordHashLength - длина хэша пароля.
	passwordHashLength = 16
)

var (
	// ErrLoginAlreadyExist - user with this login already exists.
	ErrLoginAlreadyExist = errors.New("user with this login already exists")
	// ErrCredentialsAlreadyExist - user with this credentials already exists.
	ErrCredentialsAlreadyExist = errors.New("user with this credentials already exists")
)

type Pinger interface {
	Ping(context.Context) error
}

// Crypter - интерфейс для шифрования данных.
type Crypter interface {
	Create(string, int) string
}

// Inserter - интерфейс для добавления пользователя.
type Inserter interface {
	InsertUser(context.Context, model.Auth) error
}

// Exister - интерфейс для проверки пользователя.
type Exister interface {
	ExistLogin(context.Context, string) (bool, error)
	ExistHashPassword(context.Context, string) (bool, error)
}

// Repository - интерфейс для базы данных.
type Repository interface {
	Pinger
	Inserter
	Exister
	Close() error
}

// Usecase - какое-то описание.
type Usecase struct {
	Repo   Repository
	Crypto Crypter
}

// New - создание структуры Usecase.
func New(repo Repository) *Usecase {
	return &Usecase{
		Repo:   repo,
		Crypto: crpt.New(),
	}
}

// Ping - проверка соединения с базой данных.
func (u *Usecase) Ping(ctx context.Context) error {
	return u.Repo.Ping(ctx)
}

// Register - регистрация пользователя.
func (u *Usecase) Register(ctx context.Context, auth model.Auth) error {
	exists, err := u.Repo.ExistLogin(ctx, auth.Login)
	if err != nil {
		return err
	}
	if exists {
		return ErrLoginAlreadyExist
	}
	auth.PasswordHash = u.Crypto.Create(auth.Password, passwordHashLength)

	if err := u.Repo.InsertUser(ctx, auth); err != nil {
		return err
	}

	return nil
}

// Authentication - аутентификация пользователя.
func (u *Usecase) Authentication(ctx context.Context, auth model.Auth) error {
	var exists bool
	exists, err := u.Repo.ExistLogin(ctx, auth.Login)
	if err != nil {
		return err
	}
	if !exists {
		return ErrCredentialsAlreadyExist
	}

	auth.PasswordHash = u.Crypto.Create(auth.Password, passwordHashLength)

	exists, err = u.Repo.ExistHashPassword(ctx, auth.PasswordHash)
	if err != nil {
		return err
	}
	if !exists {
		return ErrCredentialsAlreadyExist
	}
	return nil
}
