package local

import (
	"github.com/Di-nis/gophKeeper/internal/client/storage"
	"github.com/Di-nis/gophKeeper/internal/model"
)

// Repo - структура базы данных.
type Repo struct {
	Storage *storage.Storage
}

// New - создание структуры Repo.
func New(path string) (*Repo, error) {
	storage, err := storage.New(path)
	if err != nil {
		return nil, err
	}
	return &Repo{Storage: storage}, nil
}

// Close - закрытие файла.
func (repo *Repo) Close() error {
	if repo.Storage != nil {
		if err := repo.Storage.Producer.Close(); err != nil {
			return err
		}
	}
	return nil
}

// Write - запись в файл.
func (repo *Repo) Write(data model.Common) error {
	return repo.Storage.Producer.Write(data)
}
