package auth

import (
	"os"
)

// ReadFileAsString - чтение файла.
func ReadFileAsString(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
