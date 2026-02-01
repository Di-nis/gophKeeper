package http

import (
	"bytes"
	"encoding/json"
	"io"
	"os"

	"github.com/Di-nis/gophKeeper/internal/model"
)

// WriteToFile - функция записи в файл.
func writeToFile(path string, data string) error {
	return os.WriteFile(path, []byte(data), 0644)
}

// // ReadFromFile - функция чтения из файла.
// func ReadFromFile(filename string) (string, error) {
// 	data, err := os.ReadFile(filename)
// 	if err != nil {
// 		return "", err
// 	}
// 	return string(data), nil
// }

// createBody - функция создания тела запроса.
func createBody(auth model.Auth) (io.Reader, error) {
	b, err := json.Marshal(auth)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}
