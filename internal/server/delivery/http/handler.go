// Package http предоставляет HTTP-транспорт для приложения, реализуя обработчики
// маршрутов и преобразуя сетевые запросы в вызовы сервисного слоя.
package http

import (
	"errors"
	"net/http"

	// "net/http/pprof"
	// "reflect"

	// "github.com/Di-nis/gophKeeper/internal/server/middleware/auth"
	// "github.com/Di-nis/gophKeeper/internal/server/middleware/cidr"
	"github.com/Di-nis/gophKeeper/internal/model"
	"github.com/Di-nis/gophKeeper/internal/server/config"
	"github.com/Di-nis/gophKeeper/internal/server/usecase/user"
	"github.com/Di-nis/gophKeeper/pkg/logger"

	"context"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// Pinger - интерфейс для проверки соединения с БД.
type Pinger interface {
	Ping(context.Context) error
}

// Auth - интерфейс для регистрации/авторизации.
type Auth interface {
	Register(context.Context, model.Auth) error
	// Login(context.Context, model.Auth) error
}

// Handler - структура HTTP-хендлера.
type Handler struct {
	Pinger Pinger
	Auth   Auth
	Config *config.Config
}

// New - создание структуры Controller.
func New(pinger Pinger, auth Auth, config *config.Config) *Handler {
	return &Handler{
		Pinger: pinger,
		Auth:   auth,
		Config: config,
	}
}

// ping - пинг БД.
func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	err := h.Pinger.Ping(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Register - регистрация пользователя.
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var err error
	a := model.Auth{}

	if err := readReq(r, &a); err != nil {
		logger.Sugar.Errorf("cannot read request: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = h.Auth.Register(r.Context(), a); err != nil {
		if errors.Is(err, user.ErrLoginAlreadyExist) {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		logger.Sugar.Errorf("cannot register user: %v", err)
	}
	w.WriteHeader(http.StatusOK)
}

// // Login - аутентификация пользователя.
// func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
// 	c := model.Creds{}

// 	if err := handler.ReadReq(r, &c); err != nil {
// 		logger.Sugar.Errorf("cannot read request: %v", err)
// 		w.WriteHeader(http.StatusBadRequest)

// 		return
// 	}

// 	jwtToken, err := h.service.Login(r.Context(), c)
// 	if err != nil {
// 		if errors.Is(err, repository.ErrUserNotFound) {
// 			logger.Sugar.Errorf("user not found registered: %v", err)
// 			w.WriteHeader(http.StatusUnauthorized)

// 			return
// 		}

// 		logger.Sugar.Errorf("can't register user: %v", err)
// 		w.WriteHeader(http.StatusInternalServerError)

// 		return
// 	}

// 	cookieAuth := &http.Cookie{
// 		Name:  "auth_user",
// 		Value: jwtToken,
// 		Path:  "/",
// 	}

// 	http.SetCookie(w, cookieAuth)

// 	w.WriteHeader(http.StatusOK)
// }
