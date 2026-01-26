package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Di-nis/gophKeeper/internal/model"
	"github.com/Di-nis/gophKeeper/internal/server/repository"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

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
func (repo *Repo) InsertUser(ctx context.Context, user model.Auth) error {
	repo.m.Lock()
	defer repo.m.Unlock()
	// TODO: доработать, пока ошибка
	query := "INSERT INTO auth (id, login, password_hash) VALUES ($1, $2, $3)"
	_, err := repo.db.ExecContext(ctx, query, user.ID, user.Login, user.PasswordHash)
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

// ExistHashPassword - проверка существования хэша пароля.
func (repo *Repo) ExistHashPassword(ctx context.Context, passwordHash string) (bool, error) {
	repo.m.RLock()
	defer repo.m.RUnlock()

	var exists bool

	query := "SELECT EXISTS (SELECT password_hash FROM auth WHERE password_hash = $1)"
	row := repo.db.QueryRowContext(ctx, query, passwordHash)

	err := row.Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("path: internal/server/repository/postgres/postgres.go, func ExistHashPassword(): %w", err)
	}

	if exists {
		return true, nil
	}
	return false, nil
}

// Insert - добавление данных.
func (repo *Repo) Insert(ctx context.Context, data any) error {
	repo.m.Lock()
	defer repo.m.Unlock()

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

// Select - получение данных.
func (repo *Repo) Select(ctx context.Context, data any) error {
	repo.m.RLock()
	defer repo.m.RUnlock()

	switch d := data.(type) {
	case *model.Credentials:
		return repo.SelectCredentials(ctx, d)
	case *model.PaymentCard:
		return repo.SelectPaymentCard(ctx, d)
	default:
		return fmt.Errorf("path: internal/server/repository/postgres/postgres.go, func Select(): %w", repository.ErrUnknownType)
	}
}

// Delete - удаление данных.
func (repo *Repo) Delete(ctx context.Context, data any) error {
	repo.m.Lock()
	defer repo.m.Unlock()

	switch d := data.(type) {
	case []*model.Credentials:
		return repo.DeleteCredentials(ctx, d)
	case []*model.PaymentCard:
		return repo.DeletePaymentCard(ctx, d)
	case []*model.Binary:
		return repo.DeleteBinary(ctx, d)
	case []*model.Text:
		return repo.DeleteText(ctx, d)
	}

	return nil
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

// InsertCredentials - добавление учетных данных.
func (repo *Repo) InsertCredentials(ctx context.Context, cred *model.Credentials) error {
	query := "INSERT INTO credentials (user_id, login, password, alias, info) VALUES ($1, $2, $3, $4, $5)"

	_, err := repo.db.ExecContext(ctx, query, cred.UUID, cred.Login, cred.Password, cred.Alias, cred.Info)
	if err != nil {
		var pqErr *pgconn.PgError
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return fmt.Errorf("path: internal/server/repository/postgres/postgres.go, func InsertCredentials(): %w", repository.ErrDataAlreadyExists)
		}
		return err
	}

	return nil
}

// SelectCredentials - получение учетных данных.
func (repo *Repo) SelectCredentials(ctx context.Context, cred *model.Credentials) error {
	query := "SELECT login, password, info FROM credentials WHERE alias = $1 AND user_id = $2"
	row := repo.db.QueryRowContext(ctx, query, cred.Alias, cred.UUID)

	err := row.Scan(&cred.Login, &cred.Password, &cred.Info)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("path: internal/server/repository/postgres/postgres.go, func SelectCredentials(): %w", repository.ErrDataNotFound)
	}

	return nil
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

// DeleteCredentials - удаление учетных данных.
func (repo *Repo) DeleteCredentials(ctx context.Context, creds []*model.Credentials) error {
	repo.m.Lock()
	defer repo.m.Unlock()

	if len(creds) == 0 {
		return repository.ErrNoData
	}

	var values []string
	var args []any

	for i, cred := range creds {
		base := i * 2
		params := fmt.Sprintf("($%d, $%d)", base+1, base+2)
		values = append(values, params)
		args = append(args, cred.Alias, cred.UUID)
	}

	query := `
	DELETE FROM AS u FROM (VALUES ` + strings.Join(values, ",") + `) AS v(alias, user_id) WHERE u.alias = v.alias AND u.user_id = v.user_id;`

	result, err := repo.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("path: internal/server/repository/postgres/postgres.go, func DeleteCredentials(), failed to delete credentials: %w", err)
	}

	count, _ := result.RowsAffected()
	if count == 0 {
		return fmt.Errorf("path: internal/server/repository/postgres/postgres.go, func DeleteCredentials(): %w", repository.ErrDataNotFound)
	}
	return nil
}

// InsertPaymentCard - добавление банковской карты.
func (repo *Repo) InsertPaymentCard(ctx context.Context, paymentCard *model.PaymentCard) error {
	return nil
}

// SelectPaymentCard - получение данных банковской карты.
func (repo *Repo) SelectPaymentCard(ctx context.Context, card *model.PaymentCard) error {
	return nil
}

// SelectUserPaymentCards - получение всех данных типа "банковская карта" пользователя.
func (repo *Repo) SelectUserPaymentCards(ctx context.Context, userID model.UserID) ([]model.PaymentCard, error) {
	var cards []model.PaymentCard

	return cards, nil
}

// DeletePaymentCard - удаление данных банковской карты.
func (repo *Repo) DeletePaymentCard(ctx context.Context, card []*model.PaymentCard) error {
	return nil
}

// InsertBinary - добавление бинарных данных.
func (repo *Repo) InsertBinary(ctx context.Context, paymentCard *model.Binary) error {
	return nil
}

// SelectBinary - получение бинарных данных.
func (repo *Repo) SelectBinary(ctx context.Context, bin *model.Binary) error {
	return nil
}

// SelectUserBinary - получение всех бинарных данных пользователя.
func (repo *Repo) SelectUserBinaries(ctx context.Context, userID model.UserID) ([]model.Binary, error) {
	var bins []model.Binary

	return bins, nil
}

// DeleteBinary - удаление бинарных данных.
func (repo *Repo) DeleteBinary(ctx context.Context, bin []*model.Binary) error {
	return nil
}

// InsertText - добавление текстовых данных в БД.
func (repo *Repo) InsertText(ctx context.Context, paymentCard *model.Text) error {
	return nil
}

// SelectText - получение текстовых данных.
func (repo *Repo) SelectText(ctx context.Context, text *model.Text) error {
	return nil
}

// SelectUserTexts - получение всех текстовых данных пользователя.
func (repo *Repo) SelectUserTexts(ctx context.Context, userID model.UserID) ([]model.Text, error) {
	var texts []model.Text

	return texts, nil
}

// DeleteText - удаление текстовых данных.
func (repo *Repo) DeleteText(ctx context.Context, text []*model.Text) error {
	return nil
}
