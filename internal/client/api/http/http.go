package http

import (
	"net/http"
)

// attemptsCount - максимальное количество попыток при запросах на удаленный сервер.
const attemptsCount = 5

// Client - HTTP-клиент.
type Client struct {
	HTTPClient    *http.Client
	ServerAddress string
}

// New - конструктор клиента.
func New(serverAddress string) *Client {
	httpClient := &http.Client{}

	return &Client{
		HTTPClient:    httpClient,
		ServerAddress: serverAddress,
	}
}
