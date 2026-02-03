// Package http - реализация HTTP-клиента.
package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/Di-nis/gophKeeper/internal/model"
	"github.com/Di-nis/gophKeeper/internal/client/config"
	"github.com/Di-nis/gophKeeper/pkg/logger"
)

var (
	errWtiteFile            = errors.New("error write to file")
	errEmptyLoginOrPassword = errors.New("empty login or password")
	errEmptyCookie          = errors.New("empty cookie")
)

// attemptsCount - максимальное количество попыток при запросах на удаленный сервер.
const attemptsCount = 5

// URLPath - пути для запросов.
type URLPath struct {
	register string
	login    string
}

// NewURLPath - конструктор UrlPath.
func NewURLPath() *URLPath {
	return &URLPath{
		register: "/register",
		login:    "/login",
	}
}

// Client - HTTP-клиент.
type Client struct {
	HTTPClient    *http.Client
	serverAddress string
	tokenStorage  string
	urlPath       *URLPath
}

// New - конструктор клиента.
func New(cfg *config.Config) *Client {
	httpClient := &http.Client{}

	return &Client{
		HTTPClient:    httpClient,
		serverAddress: cfg.ServerAddressHTTP,
		tokenStorage:  cfg.TokenStorage,
		urlPath:       NewURLPath(),
	}
}

// Register - регистрация пользователя.
func (c *Client) Register(ctx context.Context, auth model.Auth) error {
	if auth.Login == "" || auth.Password == "" {
		return errEmptyLoginOrPassword
	}

	var (
		lastResErr error
		resp       *http.Response
	)

	for range attemptsCount {
		body, err := createBody(auth)
		if err != nil {
			return err
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.serverAddress+c.urlPath.register, body)
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err = c.HTTPClient.Do(req)
		if err != nil {
			lastResErr = err
			continue
		}

		defer resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			lastResErr = nil
			break
		}
	}

	if lastResErr != nil {
		logger.Sugar.Warnw(
			"path: internal/client/api/http/http.go, func Register(), error request",
			"attempt", attemptsCount,
			"err", lastResErr,
		)
		return lastResErr
	}

	return nil
}

// Login - авторизация пользователя.
func (c *Client) Login(ctx context.Context, auth model.Auth) error {
	if auth.Login == "" || auth.Password == "" {
		return errEmptyLoginOrPassword
	}

	var (
		lastResErr error
		resp       *http.Response
	)

	for range attemptsCount {
		body, err := createBody(auth)
		if err != nil {
			return err
		}

		if ctx.Err() != nil {
			fmt.Println("ctx already canceled:", ctx.Err())
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.serverAddress+c.urlPath.login, body)
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err = c.HTTPClient.Do(req)
		if err != nil {
			lastResErr = err
			continue
		}

		defer resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			lastResErr = nil
			break
		}
	}

	if lastResErr != nil {
		logger.Sugar.Warnw(
			"path: internal/client/api/http/http.go, func Login(), error request",
			"attempt", attemptsCount,
			"err", lastResErr,
		)
		return lastResErr
	}

	var token string
	cookie := resp.Cookies()
	if len(cookie) == 0 {
		logger.Sugar.Warnw("path: internal/client/api/http/http.go, func Login(), cookie is empty")
		return errEmptyCookie
	}

	token = cookie[0].Value

	err := writeToFile(c.tokenStorage, token)
	if err != nil {
		return errWtiteFile
	}

	return nil
}
