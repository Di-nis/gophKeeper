package storage

import (
	"bufio"
	"encoding/json"
	"os"

	"github.com/Di-nis/gophKeeper/internal/model"
)

// Consumer - структура для чтения данных из файла.
type Consumer struct {
	file    *os.File
	scanner *bufio.Scanner
}

// NewConsumer - создание нового объекта Consumer.
func NewConsumer(filename string) (*Consumer, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		file:    file,
		scanner: bufio.NewScanner(file),
	}, nil
}

// Read - чтение данных из файла.
func (c *Consumer) Read() (*model.Common, error) {
	data := c.scanner.Bytes()

	common := model.Common{}
	err := json.Unmarshal(data, &common)
	if err != nil {
		return nil, err
	}

	return &common, nil
}

// Close - закрытие файла.
func (c *Consumer) Close() error {
	return c.file.Close()
}

// Load - загрузка данных из файла.
func (c *Consumer) Load() (model.Common, error) {
	// URLArray := make([]models.URLBase, 0)

	// for c.scanner.Scan() {
	// 	urlData, err := c.Read()
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	URLArray = append(URLArray, *urlData)
	// }
	// err := c.Close()
	// if err != nil {
	// 	return nil, err
	// }
	// return URLArray, nil
	return model.Common{}, nil
}
