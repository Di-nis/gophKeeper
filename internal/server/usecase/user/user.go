// Package user - реализация usecase для работы с пользователями.
package user

import (
	"context"
	"errors"

	"github.com/Di-nis/gophKeeper/internal/model"
	"github.com/Di-nis/gophKeeper/internal/server/auth"
	crpt "github.com/Di-nis/gophKeeper/internal/server/crypto"
	"github.com/samborkent/uuidv7"

	cfg "github.com/Di-nis/gophKeeper/internal/server/config"
	"github.com/Di-nis/gophKeeper/internal/server/repository"
)

const (
	// passwordHashLength - длина хэша пароля.
	passwordHashLength = 16
)

var (
	// ErrLoginAlreadyExist - user with this login already exists.
	ErrLoginAlreadyExist = errors.New("user with this login already exists")
	// ErrUserAlreadyExist - user already exists.
	ErrUserAlreadyExist = errors.New("user already exists")
	// ErrUserNotFound - user not found.
	ErrUserNotFound = errors.New("user not found")
	// ErrBuildingToken - ошибка создания токена.
	ErrBuildingToken = errors.New("error building token")
)

// Pinger - интерфейс для проверки соединения с базой данных.
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

// Getter - интерфейс для получения пользователя.
type Getter interface {
	SelectUserID(context.Context, *model.Auth) error
}

// Repository - интерфейс для базы данных.
type Repository interface {
	Pinger
	Inserter
	Exister
	Getter
	Close() error
}

// Auth - интерфейс для аутентификации пользователя.
type Auth interface {
	BuildJWT(string, model.UserID) (string, error)
}

// Usecase - какое-то описание.
type Usecase struct {
	repo   Repository
	crypto Crypter
	auth   Auth
	config *cfg.Config
}

// New - создание структуры Usecase.
func New(repo Repository, config *cfg.Config) *Usecase {
	return &Usecase{
		repo:   repo,
		crypto: crpt.New(),
		auth:   auth.New(),
		config: config,
	}
}

// Ping - проверка соединения с базой данных.
func (u *Usecase) Ping(ctx context.Context) error {
	return u.repo.Ping(ctx)
}

// Register - регистрация/аутентификация пользователя.
func (u *Usecase) Register(ctx context.Context, auth model.Auth) error {
	var exists bool
	exists, err := u.repo.ExistLogin(ctx, auth.Login)
	if err != nil {
		return err
	}
	if exists {
		return ErrLoginAlreadyExist
	}

	uuid := uuidv7.New()
	hash := u.crypto.Create(auth.Password, passwordHashLength)

	auth.SetID(uuid).SetPasswordHash(hash).SetRole(model.RoleUser)

	if err := u.repo.InsertUser(ctx, auth); err != nil {
		if errors.Is(err, repository.ErrUserAlreadyExists) {
			return ErrLoginAlreadyExist
		}
		return err
	}

	return nil
}

// Login - аутентификация пользователя.
func (u *Usecase) Login(ctx context.Context, auth *model.Auth) (string, error) {
	var exists bool
	exists, err := u.repo.ExistLogin(ctx, auth.Login)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", ErrUserNotFound
	}

	auth.PasswordHash = u.crypto.Create(auth.Password, passwordHashLength)

	exists, err = u.repo.ExistHashPassword(ctx, auth.PasswordHash)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", ErrUserAlreadyExist
	}

	err = u.repo.SelectUserID(ctx, auth)
	if err != nil {
		return "", ErrBuildingToken
	}

	token, err := u.auth.BuildJWT(u.config.JWTSecret, auth.GetID())
	if err != nil {
		return "", ErrBuildingToken
	}

	return token, nil
}
