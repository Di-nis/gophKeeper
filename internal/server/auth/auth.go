package auth

import (
	"errors"
	"fmt"

	"github.com/Di-nis/gophKeeper/internal/model"

	"github.com/golang-jwt/jwt/v5"

	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	// ErrUserIDEmpty - user id is empty.
	ErrUserIDEmpty = errors.New("user id is empty")
	// ErrUnexpectedSighMethod - unexpected signing method.
	ErrUnexpectedSighMethod = errors.New("unexpected signing method")
)

const (
	// HeaderAuthorization - заголовок авторизации.
	HeaderAuthorization = "Cookie"
	// TokenExp - время жизни токена.
	TokenExp = time.Hour * 24
	// Key - ключ для создания токена.
	Key contextKey = "userID"
)

type contextKey string

// Claims — структура утверждений, которая включает стандартные утверждения и одно пользовательское UserID.
type Claims struct {
	jwt.RegisteredClaims
	UserID model.UserID
}

// New - создание экземпляра Claims.
func New() *Claims {
	return &Claims{}
}

// BuildJWT - создание JWT токена.
func (c *Claims) BuildJWT(secretKey string, userID model.UserID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		UserID: userID,
	})

	tokenOut, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenOut, nil
}

// GetClaims - получение утверждений.
func GetClaims(tokenIn, secretKey string) (*Claims, bool) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenIn, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("`%v:` `%s`", ErrUnexpectedSighMethod, t.Header["alg"])
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
