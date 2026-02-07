package storage

import (
	"github.com/Di-nis/gophKeeper/internal/model"
)

// WriteCloser - интерфейс для записи в файл.
type WriteCloser interface {
	Write(model.Common) error
	Close() error
}

// ReadCloser - интерфейс для чтения из файла.
type ReadCloser interface {
	Load() (model.Common, error)
	Close() error
}

// Storage - структура для хранения файлов.
type Storage struct {
	Producer WriteCloser
	Consumer ReadCloser
}

// NewStorage - создание нового хранилища.
func New(path string) (*Storage, error) {
	producer, err := NewProducer(path)
	if err != nil {
		return nil, err
	}
	consumer, err := NewConsumer(path)
	if err != nil {
		return nil, err
	}

	return &Storage{
		Producer: producer,
		Consumer: consumer,
	}, nil
}
