// Package postgres - реализация репозитория для работы с БД Postgres.
package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/Di-nis/gophKeeper/internal/model"
	"github.com/Di-nis/gophKeeper/internal/server/repository"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	// "github.com/samborkent/uuidv7"

	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/lib/pq"
)

// InitDB - инициализация БД.
func InitDB(DSN string) (*sql.DB, error) {
	if DSN == "" {
		return nil, repository.ErrEmptyDatabaseDSN
	}

	db, err := sql.Open("pgx", DSN)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(50)
	db.SetConnMaxLifetime(5 * time.Minute)
	return db, nil
}

// Repo - репозиторий для работы с БД Postgres.
type Repo struct {
	db *sql.DB
	m  sync.RWMutex
}

// New - конструктор репозитория.
func New(db *sql.DB) (*Repo, error) {
	if err := db.Ping(); err != nil {
		db.Close()
		// TODO: какая ошибка errDBClosed
		return nil, fmt.Errorf("path: internal/server/repository/postgres/postgres.go, func New(), failed ping db: %w", err)
	}

	return &Repo{db: db}, nil
}

// Ping - проверка соединения с БД.
func (repo *Repo) Ping(ctx context.Context) error {
	return repo.db.PingContext(ctx)
}

// Close - закрытие соединения с БД.
func (repo *Repo) Close() error {
	return repo.db.Close()
}

// Migrations - миграции БД.
func (repo *Repo) Migrations() error {
	var err error
	driver, err := postgres.WithInstance(repo.db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("path: internal/server/repository/postgres/postgres.go, func Migrations(), failed create driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file:migrations",
		"postgres",
		driver)
	if err != nil {
		return fmt.Errorf("path: internal/server/repository/postgres/postgres.go, func Migrations(), failed new Migrate Instance: %w", err)
	}

	err = m.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	} else if err != nil {
		return fmt.Errorf("path: internal/server/repository/postgres/postgres.go, func Migrations(), failed make migrations: %w", err)
	}
	return nil
}

// InsertUser - добавление пользователя.
func (repo *Repo) InsertUser(ctx context.Context, auth model.Auth) error {
	repo.m.Lock()
	defer repo.m.Unlock()

	query := "INSERT INTO auth (id, login, password_hash, role) VALUES ($1, $2, $3, $4)"
	_, err := repo.db.ExecContext(ctx, query, auth.GetID(), auth.GetLogin(), auth.GetPasswordHash(), auth.GetRole())
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("path: internal/server/repository/postgres/postgres.go, func InsertUser(): %w", repository.ErrUserAlreadyExists)
		}
		return fmt.Errorf("path: internal/server/repository/postgres/postgres.go, func InsertUser(): %w", err)
	}
	return nil
}

// ExistLogin - проверка существования логина.
func (repo *Repo) ExistLogin(ctx context.Context, login string) (bool, error) {
	repo.m.RLock()
	defer repo.m.RUnlock()

	var exists bool

	query := "SELECT EXISTS (SELECT login FROM auth WHERE login = $1)"
	row := repo.db.QueryRowContext(ctx, query, login)

	err := row.Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("path: internal/server/repository/postgres/postgres.go, func ExistLogin(): %w", err)
	}

	if exists {
		return true, nil
	}
	return false, nil
}

// Insert - метод/оркестратор для вставки данных.
func (repo *Repo) Insert(ctx context.Context, data any) error {
	switch d := data.(type) {
	case *model.Credentials:
		return repo.InsertCredentials(ctx, d)
	case *model.PaymentCard:
		return repo.InsertPaymentCard(ctx, d)
	case *model.Binary:
		return repo.InsertBinary(ctx, d)
	case *model.Text:
		return repo.InsertText(ctx, d)
	default:
		return fmt.Errorf("path: internal/server/repository/postgres/postgres.go, func Insert(): %w", repository.ErrUnknownType)
	}
}

// execInsert - вставка новых данных (общий метод).
func (repo *Repo) execInsert(
	ctx context.Context,
	query string,
	args ...any,
) error {
	repo.m.Lock()
	defer repo.m.Unlock()

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

// InsertPaymentCard - метод для вставки PaymentCard.
func (repo *Repo) InsertPaymentCard(ctx context.Context, card *model.PaymentCard) error {
	query := `
		INSERT INTO payment_card (user_id, number, exp_month, exp_year, cvv, alias, info)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	return repo.execInsert(
		ctx,
		query,
		card.UUID,
		card.Number,
		card.ExpMonth,
		card.ExpYear,
		card.CVV,
		card.Alias,
		card.Info,
	)
}

// InsertBinary - метод для вставки Binary.
func (repo *Repo) InsertBinary(ctx context.Context, bin *model.Binary) error {
	query := `
		INSERT INTO binary_data (user_id, data, alias, info)
		VALUES ($1, $2, $3, $4)
	`

	return repo.execInsert(
		ctx,
		query,
		bin.UUID,
		bin.Data,
		bin.Alias,
		bin.Info,
	)
}

// InsertText - метод для вставки Text.
func (repo *Repo) InsertText(ctx context.Context, text *model.Text) error {
	query := `
		INSERT INTO binary_data (user_id, data, alias, info)
		VALUES ($1, $2, $3, $4)
	`

	return repo.execInsert(
		ctx,
		query,
		text.UUID,
		text.Data,
		text.Alias,
		text.Info,
	)
}

// selectOne - получение данных (общий метод).
func (repo *Repo) selectOne(
	ctx context.Context,
	query string,
	args []any,
	dest []any,
) error {
	repo.m.RLock()
	defer repo.m.RUnlock()

	row := repo.db.QueryRowContext(ctx, query, args...)

	err := row.Scan(dest...)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("path: internal/server/repository/postgres/postgres.go, func selectOne(): %w", repository.ErrDataNotFound)
	}

	return nil
}

// Select - получение данных.
func (repo *Repo) Select(ctx context.Context, data any) error {
	switch d := data.(type) {
	case *model.Credentials:
		return repo.SelectCredentials(ctx, d)
	case *model.PaymentCard:
		return repo.SelectPaymentCard(ctx, d)
	case *model.Binary:
		return repo.SelectBinary(ctx, d)
	case *model.Text:
		return repo.InsertText(ctx, d)
	default:
		return fmt.Errorf("path: internal/server/repository/postgres/postgres.go, func Select(): %w", repository.ErrUnknownType)
	}
}

// SelectCredentials -
func (repo *Repo) SelectCredentials(ctx context.Context, cred *model.Credentials) error {
	query := `
		SELECT login, password, info
		FROM credentials
		WHERE alias = $1 AND user_id = $2
	`

	return repo.selectOne(
		ctx,
		query,
		[]any{cred.Alias, cred.UUID},
		[]any{&cred.Login, &cred.Password, &cred.Info},
	)
}

// SelectPaymentCard - получение данных банковской карты.
func (repo *Repo) SelectPaymentCard(ctx context.Context, card *model.PaymentCard) error {
	query := `
		SELECT number, exp_month, exp_year, cvv, info
		FROM payment_card
		WHERE alias = $1 AND user_id = $2
	`

	return repo.selectOne(
		ctx,
		query,
		[]any{card.Alias, card.UUID},
		[]any{&card.Number, &card.ExpMonth, &card.ExpYear, &card.CVV, &card.Info},
	)
}

// SelectBinary - получение бинарных данных.
func (repo *Repo) SelectBinary(ctx context.Context, bin *model.Binary) error {
	query := `
		SELECT data, info
		FROM binary_data
		WHERE alias = $1 AND user_id = $2
	`

	return repo.selectOne(
		ctx,
		query,
		[]any{bin.Alias, bin.UUID},
		[]any{&bin.Data, &bin.Info},
	)
}

// SelectText - получение текстовых данных.
func (repo *Repo) SelectText(ctx context.Context, text *model.Text) error {
	query := `
		SELECT data, info
		FROM text_data
		WHERE alias = $1 AND user_id = $2
	`

	return repo.selectOne(
		ctx,
		query,
		[]any{text.Alias, text.UUID},
		[]any{&text.Data, &text.Info},
	)
}

// execDelete - удаление данных (общий метод).
func (repo *Repo) execDelete(
	ctx context.Context,
	tableName string,
	pair [2]any,
) error {
	repo.m.Lock()
	defer repo.m.Unlock()

	query := fmt.Sprintf("DELETE FROM %s WHERE alias = $1 AND user_id = $2", tableName)

	result, err := repo.db.ExecContext(ctx, query, pair[0], pair[1])
	if err != nil {
		return err
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return repository.ErrDataNotFound
	}

	return nil
}

// Delete - удаление данных.
func (repo *Repo) Delete(ctx context.Context, data any) error {
	switch d := data.(type) {
	case model.Credentials:
		return repo.DeleteCredentials(ctx, d)
	case model.PaymentCard:
		return repo.DeletePaymentCard(ctx, d)
	case model.Binary:
		return repo.DeleteBinary(ctx, d)
	case model.Text:
		return repo.DeleteText(ctx, d)
	default:
		return fmt.Errorf("path: internal/server/repository/postgres/postgres.go, func Delete(): %w", repository.ErrUnknownType)
	}
}

// DeleteCredentials - удаление данных типа "логин/пароль".
func (repo *Repo) DeleteCredentials(ctx context.Context, cred model.Credentials) error {
	pair := [2]any{cred.Alias, cred.UUID}
	return repo.execDelete(ctx, "credentials", pair)
}

// DeletePaymentCard - удаление данных банковской карты.
func (repo *Repo) DeletePaymentCard(ctx context.Context, card model.PaymentCard) error {
	pair := [2]any{card.Alias, card.UUID}
	return repo.execDelete(ctx, "payment_card", pair)
}

// DeleteBinary - удаление бинарных данных.
func (repo *Repo) DeleteBinary(ctx context.Context, bin model.Binary) error {
	pair := [2]any{bin.Alias, bin.UUID}
	return repo.execDelete(ctx, "binary_data", pair)
}

// DeleteText - удаление текстовых данных.
func (repo *Repo) DeleteText(ctx context.Context, text model.Text) error {
	pair := [2]any{text.Alias, text.UUID}
	return repo.execDelete(ctx, "text_data", pair)
}

// SelectUserData - получение данных пользователя.
func (repo *Repo) SelectUserData(ctx context.Context, userID model.UserID) ([]model.Credentials, []model.PaymentCard, []model.Binary, []model.Text, error) {
	var errs []error
	creds, err := repo.SelectUserCredentials(ctx, userID)
	if err != nil {
		errs = append(errs, err)
	}
	cards, err := repo.SelectUserPaymentCards(ctx, userID)
	if err != nil {
		errs = append(errs, err)
	}
	// bins, err := repo.SelectUserBinary(ctx, userID)
	// if err != nil {

	// }
	return creds, cards, nil, nil, errors.Join(errs...)
}

// SelectUserCredentials - получение всех учетных данных пользователя.
func (repo *Repo) SelectUserCredentials(ctx context.Context, userID model.UserID) ([]model.Credentials, error) {
	repo.m.RLock()
	defer repo.m.RUnlock()

	query := "SELECT login, password, info FROM credentials WHERE user_id = $1"
	stmt, err := repo.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("path: internal/server/repository/postgres/postgres.go, func SelectUserCredentials(): %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("path: internal/server/repository/postgres/postgres.go, func SelectUserCredentials(): %w", err)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("path: internal/server/repository/postgres/postgres.go, func SelectUserCredentials():  %w", err)
	}

	var creds []model.Credentials
	for rows.Next() {
		var cred model.Credentials
		err = rows.Scan(&cred.Login, &cred.Password, &cred.Info)
		if err != nil {
			return nil, fmt.Errorf("path: internal/server/repository/postgres/postgres.go, func SelectUserCredentials():, failed to scan url: %w", err)
		}

		creds = append(creds, cred)
	}
	return creds, nil
}

// SelectUserPaymentCards - получение всех данных типа "банковская карта" пользователя.
func (repo *Repo) SelectUserPaymentCards(ctx context.Context, userID model.UserID) ([]model.PaymentCard, error) {
	var cards []model.PaymentCard

	return cards, nil
}

// SelectUserBinary - получение всех бинарных данных пользователя.
func (repo *Repo) SelectUserBinaries(ctx context.Context, userID model.UserID) ([]model.Binary, error) {
	var bins []model.Binary

	return bins, nil
}

// SelectUserTexts - получение всех текстовых данных пользователя.
func (repo *Repo) SelectUserTexts(ctx context.Context, userID model.UserID) ([]model.Text, error) {
	var texts []model.Text

	return texts, nil
}
