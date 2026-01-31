// Package config содержит параметры запуска клиента.
package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/caarlos0/env/v6"
)

// Config - структура конфигурации приложения.
type Config struct {
	ServerAddressHTTP string `env:"SERVER_ADDRESS_HTTP"`
	ServerAddressGRPC string `env:"SERVER_ADDRESS_GRPC"`
	BaseURLHTTP       string `env:"BASE_URL_HTTP"`
	BaseURLGRPC       string `env:"BASE_URL_GRPC"`
	LogLevel          string `env:"LOG_LEVEL"`
	Config            string `env:"CONFIG"`
	TokenStorage      string `env:"TOKEN_STORAGE"`
}

// New - функция для создания конфигурации.
func New() *Config {
	return &Config{}
}

// Load - метод для парсинга конфигурации из переменных окружения и аргументов командной строки.
func (c *Config) Load() error {
	var err error

	// первый приоритет - из переменных окружения
	err = c.loanFromEnv()
	if err != nil {
		return err
	}

	// второй приоритет - из файла
	err = c.loanFromFile()
	if err != nil {
		return err
	}
	return nil
}

// loanFromEnv - загрузка конфигурации из переменных окружения.
func (c *Config) loanFromEnv() error {
	if err := env.Parse(c); err != nil {
		return fmt.Errorf("path: internal/server/config/config.go, func loanFromEnv(), failed to parse env: %w", err)
	}
	return nil
}

// loanFromFile - загрузка конфигурации из файла.
func (c *Config) loanFromFile() error {
	if c.Config == "" {
		return nil
	}

	type ConfigAlias struct {
		ServerAddressHTTP string `json:"SERVER_ADDRESS_HTTP"`
		ServerAddressGRPC string `json:"SERVER_ADDRESS_GRPC"`
		BaseURLHTTP       string `json:"BASE_URL_HTTP"`
		BaseURLGRPC       string `json:"BASE_URL_GRPC"`
		LogLevel          string `json:"log_level"`
		TokenStorage      string `json:"token_storage"`
	}

	var configAlias ConfigAlias

	jsonFile, err := os.Open(c.Config)
	if err != nil {
		return fmt.Errorf("path: internal/client/config/config.go, func loanFromJSON(), failed to open json file: %w", err)
	}

	jsonFileData, err := io.ReadAll(jsonFile)
	if err != nil {
		return fmt.Errorf("path: internal/client/config/config.go, func loanFromJSON(), failed to read json file: %w", err)
	}
	defer jsonFile.Close()

	if err := json.Unmarshal(jsonFileData, &configAlias); err != nil {
		return fmt.Errorf("path: internal/client/config/config.go, func loanFromJSON(), failed to unmarshal data: %w", err)
	}

	if c.ServerAddressHTTP == "" {
		c.ServerAddressHTTP = configAlias.ServerAddressHTTP
	}
	if c.ServerAddressGRPC == "" {
		c.ServerAddressGRPC = configAlias.ServerAddressGRPC
	}
	if c.BaseURLHTTP == "" {
		c.BaseURLHTTP = configAlias.BaseURLHTTP
	}
	if c.BaseURLGRPC == "" {
		c.BaseURLGRPC = configAlias.BaseURLGRPC
	}
	if c.LogLevel == "" {
		c.LogLevel = configAlias.LogLevel
	}
	if c.TokenStorage == "" {
		c.TokenStorage = configAlias.TokenStorage
	}

	return nil
}
