package data

import (
	"context"

	"github.com/Di-nis/gophKeeper/internal/server/model"
)

type CredentialsRepository interface {
	InsertCredentials(context.Context, model.Credentials) error
	GetCredentials(context.Context, model.Credentials) error
}

type Pinger interface {
	Ping(context.Context) error
}

type DataRepository interface {
	Pinger
	CredentialsRepository
	Close() error
}

type DataUsecase struct {
	Ping     Pinger
	CredRepo CredentialsRepository
}

func NewDataUsecase(repo DataRepository) *DataUsecase {
	return &DataUsecase{
		Ping:     repo,
		CredRepo: repo,
	}
}
