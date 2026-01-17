package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/caarlos0/env/v6"
)

// Config - структура конфигурации приложения.
type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS"`
	LogLevel      string `env:"LOG_LEVEL"`
	Config        string `env:"CONFIG"`
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

	// второй приоритет - из аргументов командной строки
	c.loanFromFlags()

	// третий приоритет - из файла
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

// loanFromFlags - загрузка конфигурации из аргументов командной строки.
func (c *Config) loanFromFlags() {
	var serverAddress, logLevel, config string

	flag.StringVar(&serverAddress, "a", "", "URL")
	flag.StringVar(&logLevel, "l", "info", "log level")
	flag.StringVar(&config, "config", "", "path to the configuration file")

	flag.Parse()

	if c.ServerAddress == "" {
		c.ServerAddress = serverAddress
	}
	if c.LogLevel == "" {
		c.LogLevel = logLevel
	}
	if c.Config == "" {
		c.Config = config
	}
}

// loanFromFile - загрузка конфигурации из файла.
func (c *Config) loanFromFile() error {
	if c.Config == "" {
		return nil
	}

	type ConfigAlias struct {
		ServerAddress string `json:"server_address"`
		LogLevel      string `json:"log_level"`
	}

	var configAlias ConfigAlias

	jsonFile, err := os.Open(c.Config)
	if err != nil {
		return fmt.Errorf("path: internal/server/config/config.go, func loanFromJSON(), failed to open json file: %w", err)
	}

	jsonFileData, err := io.ReadAll(jsonFile)
	if err != nil {
		return fmt.Errorf("path: internal/server/config/config.go, func loanFromJSON(), failed to read json file: %w", err)
	}
	defer jsonFile.Close()

	if err := json.Unmarshal(jsonFileData, &configAlias); err != nil {
		return fmt.Errorf("path: internal/server/config/config.go, func loanFromJSON(), failed to unmarshal data: %w", err)
	}

	if c.ServerAddress == "" {
		c.ServerAddress = configAlias.ServerAddress
	}
	if c.LogLevel == "" {
		c.LogLevel = configAlias.LogLevel
	}

	return nil
}
