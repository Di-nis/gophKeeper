// Package data - реализация usecase для работы с данными.
package data

import (
	"context"
	"errors"

	"sync"

	"github.com/Di-nis/gophKeeper/internal/model"
	crpt "github.com/Di-nis/gophKeeper/internal/server/crypto"
	repo "github.com/Di-nis/gophKeeper/internal/server/repository"
)

const (
	// aliasLength - длина алиаса.
	aliasLength = 16
	// numWorkers - количество воркеров.
	numWorkers = 3
)

var (
	// ErrDataNotFound - data not found.
	ErrDataNotFound = errors.New("data not found")
	// ErrDataRegistered - data already registered.
	ErrDataRegistered = errors.New("data already registered")
	// ErrGetRawValue - error get raw value.
	ErrGetRawValue = errors.New("error get raw value")
	// ErrDataUnsupportedType - unsupported data type.
	ErrDataUnsupportedType = errors.New("data unsupported type")
)

// Data - тип данных.
type Data[T model.Credentials | model.PaymentCard | model.Binary | model.Text] []T

// Pinger - интерфейс для проверки соединения с базой данных.
type Pinger interface {
	Ping(context.Context) error
}

// Crypter - интерфейс для шифрования данных.
type Crypter interface {
	Create(string, int) string
}

// Inserter - интерфейс для вставки данных в базу данных.
type Inserter interface {
	Insert(context.Context, any) error
}

// Selecter - интерфейс для получения данных из базы данных.
type Selecter interface {
	Select(context.Context, any) error
	SelectUserData(context.Context, model.UserID) ([]model.Credentials, []model.PaymentCard, []model.Binary, []model.Text, error)
}

// Deleter - интерфейс для удаления данных из базы данных.
type Deleter interface {
	Delete(context.Context, any) error
}

// Repository - интерфейс для базы данных.
type Repository interface {
	Pinger
	Inserter
	Selecter
	Deleter
	Close() error
}

// Usecase - реализация usecase для работы с данными.
type Usecase struct {
	Repo   Repository
	crypto Crypter
}

// New - создание структуры Usecase.
func New(repo Repository) *Usecase {
	return &Usecase{
		Repo:   repo,
		crypto: crpt.New(),
	}
}

// CreateCredentials - создание данных типа "логин/пароль".
func (u *Usecase) Create(ctx context.Context, data any) error {
	var err error
	raw, err := getRawValue(data)
	if err != nil {
		return err
	}
	alias := u.crypto.Create(raw, aliasLength)

	err = updateData(data, alias)
	if err != nil {
		return err
	}

	if err := u.Repo.Insert(ctx, data); err != nil {
		if errors.Is(err, repo.ErrDataAlreadyExists) {
			return ErrDataRegistered
		}
		return err
	}
	return nil
}

// GetCredentials - получение информации о данных типа "логин/пароль".
func (u *Usecase) Get(ctx context.Context, data any) error {
	err := u.Repo.Select(ctx, data)
	if err != nil {
		if errors.Is(err, repo.ErrDataNotFound) {
			return ErrDataNotFound
		}
		return err
	}
	return nil
}

// GetAll - получение всех данных пользователя.
func (u *Usecase) GetAll(ctx context.Context, userID model.UserID) (
	[]model.Credentials, []model.PaymentCard, []model.Binary, []model.Text, error) {
	var (
		err   error
		creds []model.Credentials
		cards []model.PaymentCard
		bins  []model.Binary
		texts []model.Text
	)

	if creds, cards, bins, texts, err = u.Repo.SelectUserData(ctx, userID); err != nil {
		return nil, nil, nil, nil, err
	}

	return creds, cards, bins, texts, nil
}

// Delete - удаление данных.
func (u *Usecase) Delete(ctx context.Context, items []any) error {
	inChan := make(chan any, 1024)
	resultChan := make(chan error, numWorkers)
	var wg sync.WaitGroup

	go func() {
		defer close(inChan)
		u.generator(ctx, items, inChan)
	}()

	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			u.worker(ctx, inChan, resultChan)
		}(w)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	var firstErr error
	for err := range resultChan {
		if err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

// generator - генерирует сообщения в канал.
func (u *Usecase) generator(ctx context.Context, items []any, inChan chan any) {
	for _, d := range items {
		select {
		case <-ctx.Done():
			return
		case inChan <- d:
		}
	}
}

// worker - работник.
func (u *Usecase) worker(ctx context.Context, items <-chan any, result chan error) {
	itemsDB := make([]any, 0, 100)

	for {
		select {
		case <-ctx.Done():
			if len(itemsDB) > 0 {
				result <- u.Repo.Delete(ctx, itemsDB)
			}
			return

		case item, ok := <-items:
			if !ok {
				if len(itemsDB) > 0 {
					result <- u.Repo.Delete(ctx, itemsDB)
				}
				return
			}
			itemsDB = append(itemsDB, item)
			if len(itemsDB) >= 1 {
				result <- u.Repo.Delete(ctx, itemsDB)
				itemsDB = itemsDB[:0]
			}
		}
	}

}
