package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Di-nis/gophKeeper/internal/client/repository"
	"github.com/Di-nis/gophKeeper/internal/model"
	"github.com/jackc/pgx/v5/pgconn"

	_ "github.com/mattn/go-sqlite3"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// InitDB - инициализация БД.
func InitDB(path string) (*sql.DB, error) {
	if path == "" {
		return nil, repository.ErrEmptyDatabasePath
	}

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Repo - репозиторий для работы с БД sqlite3.
type Repo struct {
	db *sql.DB
}

// New - конструктор репозитория.
func New(db *sql.DB) (*Repo, error) {
	return &Repo{db: db}, nil
}

// Close - закрытие соединения с БД.
func (repo *Repo) Close() error {
	return repo.db.Close()
}

// Migrations - миграции БД.
func (repo *Repo) Migrations() error {
	var err error
	driver, err := sqlite3.WithInstance(repo.db, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("path: internal/client/repository/sqlite/sqlite.go, method Migrations(), failed create driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file:migrations/common",
		"sqlite3",
		driver)
	if err != nil {
		return fmt.Errorf("path: internal/client/repository/sqlite/sqlite.go, method Migrations(), failed new Migrate Instance: %w", err)
	}

	err = m.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	} else if err != nil {
		return fmt.Errorf("path: internal/client/repository/sqlite/sqlite.go, method Migrations(), failed make migrations: %w", err)
	}
	return nil
}

// execInsert - вставка новых данных (общий метод).
func (repo *Repo) execInsert(
	ctx context.Context,
	query string,
	args ...any,
) error {
	// repo.m.Lock()
	// defer repo.m.Unlock()

	_, err := repo.db.ExecContext(ctx, query, args...)
	if err != nil {
		var pqErr *pgconn.PgError
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return repository.ErrDataAlreadyExists
		}
		return err
	}
	return nil
}

// InsertCredentials - метод для вставки Credentials.
func (repo *Repo) InsertCredentials(ctx context.Context, cred *model.Credentials) error {
	query := `
		INSERT INTO credentials (user_id, login, password, alias, info)
		VALUES ($1, $2, $3, $4, $5)
	`

	return repo.execInsert(
		ctx,
		query,
		cred.UUID,
		cred.Login,
		cred.Password,
		cred.Alias,
		cred.Info,
	)
}

// // InsertPaymentCard - метод для вставки PaymentCard.
// func (repo *Repo) InsertPaymentCard(ctx context.Context, card *model.PaymentCard) error {
// 	query := `
// 		INSERT INTO payment_card (user_id, number, exp_month, exp_year, cvv, alias, info)
// 		VALUES ($1, $2, $3, $4, $5, $6, $7)
// 	`

// 	return repo.execInsert(
// 		ctx,
// 		query,
// 		card.UUID,
// 		card.Number,
// 		card.ExpMonth,
// 		card.ExpYear,
// 		card.CVV,
// 		card.Alias,
// 		card.Info,
// 	)
// }

// // InsertBinary - метод для вставки Binary.
// func (repo *Repo) InsertBinary(ctx context.Context, bin *model.Binary) error {
// 	query := `
// 		INSERT INTO binary_data (user_id, data, alias, info)
// 		VALUES ($1, $2, $3, $4)
// 	`

// 	return repo.execInsert(
// 		ctx,
// 		query,
// 		bin.UUID,
// 		bin.Data,
// 		bin.Alias,
// 		bin.Info,
// 	)
// }

// // InsertText - метод для вставки Text.
// func (repo *Repo) InsertText(ctx context.Context, text *model.Text) error {
// 	query := `
// 		INSERT INTO text_data (user_id, data, alias, info)
// 		VALUES ($1, $2, $3, $4)
// 	`

// 	return repo.execInsert(
// 		ctx,
// 		query,
// 		text.UUID,
// 		text.Data,
// 		text.Alias,
// 		text.Info,
// 	)
// }
