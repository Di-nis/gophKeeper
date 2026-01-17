package data

import (
	"context"
	"errors"
	// "sync"+-

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
)

// Pinger - интерфейс для проверки соединения с базой данных.
type Pinger interface {
	Ping(context.Context) error
}

// Crypter - интерфейс для шифрования данных.
type Crypter interface {
	Create(string, int) string
}

// EventPublisher - интерфейс для публикации событий.
type EventPublisher interface {
	Publish(context.Context, any) error
}

// Credentialer - интерфейс для создания данных типа "логин/пароль".
type Credentialer interface {
	InsertCredentials(context.Context, *model.Credentials) error
	SelectCredentials(context.Context, *model.Credentials) error
	SelectUserCredentials(context.Context, model.UserID) ([]model.Credentials, error)
	DeleteCredentials(context.Context, []model.Credentials) error
}

type PaymentCarder interface {
	InsertPaymentCard(context.Context, model.PaymentCard) error
	SelectPaymentCard(context.Context, string) (model.PaymentCard, error)
	SelectUserPaymentCards(context.Context, model.UserID) ([]model.PaymentCard, error)
	DeletePaymentCard(context.Context, string) error
}

type Binarier interface {
	InsertBinary(context.Context, model.Binary) error
	SelectBinary(context.Context, string) (model.Binary, error)
	SelectUserBinaries(context.Context, model.UserID) ([]model.Binary, error)
	DeleteBinary(context.Context, string) error
}

type Texter interface {
	InsertText(context.Context, model.Text) error
	SelectText(context.Context, string) (model.Text, error)
	SelectUserTexts(context.Context, model.UserID) ([]model.Text, error)
	DeleteText(context.Context, string) error
}

// Repository - интерфейс для базы данных.
type Repository interface {
	Pinger
	Credentialer
	PaymentCarder
	Binarier
	Texter
	Close() error
}

// Usecase - реализация usecase для работы с данными.
type Usecase struct {
	Repo      Repository
	crypto    Crypter
	publisher EventPublisher
}

// CreateCredentials - создание данных типа "логин/пароль".
func (u *Usecase) CreateCredentials(ctx context.Context, cred *model.Credentials) error {
	raw := cred.Login + cred.Password + cred.Info
	cred.Alias = u.crypto.Create(raw, aliasLength)

	if err := u.Repo.InsertCredentials(ctx, cred); err != nil {
		if errors.Is(err, repo.ErrDataAlreadyExists) {
			return ErrDataRegistered
		}
		return err
	}
	return nil
}

// GetCredentials - получение информации о данных типа "логин/пароль".
func (u *Usecase) GetCredentials(ctx context.Context, cred *model.Credentials) error {
	err := u.Repo.SelectCredentials(ctx, cred)
	if err != nil {
		if errors.Is(err, repo.ErrDataNotFound) {
			return ErrDataNotFound
		}
		return err
	}
	return nil
}

// DeleteCredentials - удаление данных типа "логин/пароль".
func (u *Usecase) DeleteCredentials(ctx context.Context, creds []model.Credentials) error {
	// inChan := make(chan model.Credentials, 1024)
	// resultChan := make(chan error, numWorkers)
	// var wg sync.WaitGroup

	// go func() {
	// 	defer close(inChan)
	// 	urlUseCase.generator(ctx, creds, inChan)
	// }()

	// for w := 1; w <= numWorkers; w++ {
	// 	wg.Add(1)
	// 	go func(id int) {
	// 		defer wg.Done()
	// 		urlUseCase.worker(ctx, inChan, resultChan)
	// 	}(w)
	// }

	// go func() {
	// 	wg.Wait()
	// 	close(resultChan)
	// }()

	// var firstErr error
	// for err := range resultChan {
	// 	if err != nil && firstErr == nil {
	// 		firstErr = err
	// 	}
	// }
	// return firstErr
	return nil
}

// CreatePaymentCard - создание данных типа "банковская карта".
func (u *Usecase) CreatePaymentCard(ctx context.Context, card model.PaymentCard) error {
	return nil
}

// GetPaymentCard - получение данных типа "банковская карта".
func (u *Usecase) GetPaymentCard(ctx context.Context, alias string) (model.PaymentCard, error) {
	var card model.PaymentCard
	return card, nil
}

// DeletePaymentCard - удаление данных типа "банковская карта".
func (u *Usecase) DeletePaymentCard(ctx context.Context, alias string) error {
	return nil
}

// CreateBinary - создание бинарных данных.
func (u *Usecase) CreateBinary(ctx context.Context, bin model.Binary) error {
	return nil
}

// GetBinary - получение бинарных данных.
func (u *Usecase) GetBinary(ctx context.Context, alias string) (model.Binary, error) {
	var bin model.Binary
	return bin, nil
}

// DeleteBinary - удаление бинарных данных.
func (u *Usecase) DeleteBinary(ctx context.Context, alias string) error {
	return nil
}

// CreateText - создание текстовых данных.
func (u *Usecase) CreateText(ctx context.Context, text model.Text) error {
	return nil
}

// GetText - получение текстовых данных.
func (u *Usecase) GetText(ctx context.Context, alias string) (model.Text, error) {
	var text model.Text
	return text, nil
}

// DeleteText - удаление текстовых данных.
func (u *Usecase) DeleteText(ctx context.Context, alias string) error {
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

	if creds, err = u.Repo.SelectUserCredentials(ctx, userID); err != nil {
		return nil, nil, nil, nil, err
	}
	// if cards, err = u.Repo.SelectUserPaymentCards(ctx, userID); err != nil {

	// }

	return creds, cards, bins, texts, nil
}

// New - создание структуры Usecase.
func New(repo Repository) *Usecase {
	return &Usecase{
		Repo:   repo,
		crypto: crpt.New(),
	}
}
