package auth

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"

	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type contextKey string

const (
	HeaderAuthorization = "authorization"
	TokenExp             = time.Hour * 3
	UserIDKey contextKey = "userID"
)

// Claims — структура утверждений, которая включает стандартные утверждения
// и одно пользовательское — UserID
type Claims struct {
	jwt.RegisteredClaims
	SID    string
	UserID string
}

// BuildJWTString создаёт токен и возвращает его в виде строки.
func BuildJWTString(secretKey, userID, sessionID string) (string, error) {
	if userID == "" {
		return "", errors.New("user id is empty")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		SID:    sessionID,
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GenerateSessionID - генерация уникального идентификатора сессии.
func GenerateSessionID() string {
	return uuid.NewString()
}

// GenerateUserID - генерация уникального идентификатора пользователя.
func GenerateUserID() string {
	t := time.Now()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	id := ulid.MustNew(ulid.Timestamp(t), entropy).String()
	return id
}

// GetClaims - получение утверждений.
func GetClaims(tokenString, secretKey string) (*Claims, bool) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(secretKey), nil
		})
	if err != nil {
		return nil, false
	}

	if !token.Valid {
		return nil, false
	}
	return claims, true
}
