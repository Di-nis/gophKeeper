// Package http предоставляет HTTP-транспорт для приложения, реализуя обработчики
// маршрутов и преобразуя сетевые запросы в вызовы сервисного слоя.
package http

import (
	"errors"
	"net/http"

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
	Login(context.Context, *model.Auth) (string, error)
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

// @Summary Ping
// @Description Проверка сервиса
// @Tags test
// @Accept json
// @Success 200
// @Failure 500 "внутренная ошибка сервера"
// @Router /ping [get]
// Ping - пинг БД.
func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	err := h.Pinger.Ping(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// @Summary Register
// @Description Регистрация пользователя
// @Tags users
// @Accept json
// @Param Auth body model.Auth true "данные пользователя"
// @Success 200
// @Failure 400 "неверный запрос"
// @Failure 409 "пользователя с таким логином уже существует"
// @Failure 500 "внутренная ошибка сервера"
// @Router /register [post]
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

// @Summary Login
// @Description Получение токена для авторизации
// @Tags users
// @Accept json
// @Param Auth body model.Auth true "данные пользователя"
// @Success 200
// @Header 200 {string} Set-Cookie "name=auth_user; value=your_token; expires=Mon, 02 Feb 2026 13:08:01 GMT; Path=/; Domain=localhost; session_id=abc123; HttpOnly secure=true"
// @Failure 400 "неверный запрос"
// @Failure 404 "пользователя с таким логином не найден"
// @Failure 500 "внутренная ошибка сервера"
// @Router /login [post]
// Login - Получение токена для авторизации.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	a := &model.Auth{}

	if err := readReq(r, a); err != nil {
		logger.Sugar.Errorf("cannot read request: %v", err)
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	jwtToken, err := h.Auth.Login(r.Context(), a)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			logger.Sugar.Errorf("user not found registered: %v", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		logger.Sugar.Errorf("internal error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	cookie := &http.Cookie{
		Name:     "auth_user",
		Value:    jwtToken,
		Expires:  time.Now().Add(24 * time.Hour),
		Path:     "/",
		Domain:   "localhost",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   true,
	}

	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusOK)
}
