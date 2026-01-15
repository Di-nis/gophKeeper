package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Di-nis/gophKeeper/internal/server/model"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
)

// RepoPostgres - репозиторий для работы с БД Postgres.
type RepoPostgres struct {
	db *sql.DB
}

// NewRepoPostgres - конструктор репозитория.
func NewRepoPostgres(dataSourceName string) (*RepoPostgres, error) {
	db, err := sql.Open("pgx", dataSourceName)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(50)
	db.SetConnMaxLifetime(5 * time.Minute)

	return &RepoPostgres{
		db: db,
	}, nil
}

// Ping - проверка соединения с БД.
func (repo *RepoPostgres) Ping(ctx context.Context) error {
	return repo.db.PingContext(ctx)
}

// Close - закрытие соединения с БД.
func (repo *RepoPostgres) Close() error {
	return repo.db.Close()
}

// Migrations - миграции БД.
func (repo *RepoPostgres) Migrations() error {
	var err error
	driver, err := postgres.WithInstance(repo.db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("path: internal/repository/postgres_repository.go, func Migrations(), failed create driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file:migrations",
		"postgres",
		driver)
	if err != nil {
		return fmt.Errorf("path: internal/repository/postgres_repository.go, func Migrations(), failed new Migrate Instance: %w", err)
	}

	err = m.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	} else if err != nil {
		return fmt.Errorf("path: internal/repository/postgres_repository.go, func Migrations(), failed make migrations: %w", err)
	}
	return nil
}

// InsertUser - добавление пользователя в БД.
func (repo *RepoPostgres) InsertUser(ctx context.Context, user model.Auth) error {
	return nil
}

// GetUser - получение пользователя из БД.
func (repo *RepoPostgres) SelectUser(ctx context.Context, user model.Auth) error {
	return nil
}

// InsertCredentials - добавление учетных данных в БД.
func (repo *RepoPostgres) InsertCredentials(ctx context.Context, сredentials model.Credentials) error {
	return nil
}

// func (repo *RepoPostgres) SelectCredentials(ctx context.Context, user model.Credentials) error {
// 	return nil
// }

// InsertPaymentCard - добавление банковской карты в БД.
func (repo *RepoPostgres) InsertPaymentCard(ctx context.Context, paymentCard model.PaymentCard) error {
	return nil
}

// InsertBinary - добавление бинарных данных в БД.
func (repo *RepoPostgres) InsertBinary(ctx context.Context, paymentCard model.Binary) error {
	return nil
}

// InsertText - добавление текстовых данных в БД.
func (repo *RepoPostgres) InsertText(ctx context.Context, paymentCard model.Text) error {
	return nil
}
