package auth

import (
	"errors"
	"fmt"

	"github.com/Di-nis/gophKeeper/internal/model"

	"github.com/golang-jwt/jwt/v5"

	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// ErrUserIDEmpty - ошибка, когда пользователь не найден.
var (
	ErrUserIDEmpty          = errors.New("user id is empty")
	ErrUnexpectedSighMethod = errors.New("unexpected signing method")
)

type contextKey string

const (
	HeaderAuthorization            = "authorization"
	TokenExp                       = time.Hour * 3
	UserIDKey           contextKey = "userID"
)

// Claims — структура утверждений, которая включает стандартные утверждения и одно пользовательское UserID.
type Claims struct {
	jwt.RegisteredClaims
	SID    string
	UserID string
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
		UserID: userID.String(),
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
