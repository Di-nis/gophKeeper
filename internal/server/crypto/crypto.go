package crypto

import (
	"crypto/sha256"
)

var base62Alphabet = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// Crypto - структура сервиса по шифторванию данных.
type Crypto struct{}

// New - создание структуры Service.
func New() *Crypto {
	return &Crypto{}
}

// Encode - преобразование числа в Base62.
func (c *Crypto) Encode(num uint64) string {
	if num == 0 {
		return string(base62Alphabet[0])
	}
	var encoded []byte
	for num > 0 {
		rem := num % 62
		num /= 62

		encoded = append([]byte{base62Alphabet[rem]}, encoded...)
	}
	return string(encoded)
}

// Create - создание хэш на основе входных данных.
func (c *Crypto) Create(data string, length int) string {
	var num uint64
	hash := sha256.Sum256([]byte(data))
	for i := range 8 {
		num = (num << 8) | uint64(hash[i])
	}

	b62 := c.Encode(num)
	if len(b62) > length {
		return b62[:length]
	}
	return b62
}
