package storage

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/Di-nis/gophKeeper/internal/model"
)

// Producer - структура для записи данных в файл.
type Producer struct {
	file   *os.File
	writer *bufio.Writer
}

// NewProducer - создание нового объекта Producer.
func NewProducer(path string) (*Producer, error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return nil, err
	}
	return &Producer{
		file:   file,
		writer: bufio.NewWriter(file),
	}, nil
}

// Write - запись данных в файл.
func (p *Producer) Write(common model.Common) error {
	data, err := json.Marshal(&common)
	if err != nil {
		return err
	}

	if _, err := p.writer.Write(data); err != nil {
		return err
	}

	if err := p.writer.WriteByte('\n'); err != nil {
		return err
	}

	return p.writer.Flush()
}

// Close - закрытие файла.
func (p *Producer) Close() error {
	return p.file.Close()
}
