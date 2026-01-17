// Package config содержит структуры и функции для загрузки и проверки
// конфигурации приложения из переменных окружения, файлов или других источников.
package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/caarlos0/env/v6"
)

// Config - структура конфигурации приложения.
type Config struct {
	ServerAddressHTTP string        `env:"SERVER_ADDRESS_HTTP"`
	ServerAddressGRPC string        `env:"SERVER_ADDRESS_GRPC"`
	BaseURLHTTP       string        `env:"BASE_URL_HTTP"`
	BaseURLGRPC       string        `env:"BASE_URL_GRPC"`
	LogLevel          string        `env:"LOG_LEVEL"`
	DatabaseDSN       string        `env:"DATABASE_DSN"`
	JWTSecret         string        `env:"JWT_SECRET"`
	EnableHTTPS       bool          `env:"ENABLE_HTTPS"`
	CertFilePath      string        `env:"CERT_FILE_PATH"`
	KeyFilePath       string        `env:"KEY_FILE_PATH"`
	Config            string        `env:"CONFIG"`
	TrustedSubnet     string        `env:"TRUSTED_SUBNET"`
	ShutdownTimeout   time.Duration `env:"SHUTDOWN_TIMEOUT"`
	UseMockAuth       bool
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
	var (
		serverAddressHTTP, baseURLHTTP, serverAddressGRPC, baseURLGRPC    string
		logLevel, dataBaseDSN, auditFile, auditURL, config, trustedSubnet string
		enableHTTPS, enableGRPC, useHeader                                bool
		shutdownTimeout                                                   time.Duration
	)
	flag.StringVar(&serverAddressHTTP, "address-http", "", "URL")
	flag.StringVar(&baseURLHTTP, "url-http", "", "HTTP-server, base URL")
	flag.StringVar(&serverAddressGRPC, "address-grpc", "", "URL")
	flag.StringVar(&baseURLGRPC, "url-grpc", "", "gRPC-server, base URL")
	flag.StringVar(&logLevel, "l", "info", "log level")
	flag.StringVar(&dataBaseDSN, "d", "", "database dource name")
	flag.StringVar(&auditFile, "audit-file", "", "path to audit file")
	flag.StringVar(&auditURL, "audit-url", "", "path to audit URL")
	flag.StringVar(&config, "config", "", "path to the configuration file")
	flag.StringVar(&trustedSubnet, "t", "", "CIDR")
	flag.BoolVar(&enableHTTPS, "s", false, "use HTTPS web-server")
	flag.BoolVar(&useHeader, "use-header", false, "using a header when parsing an IP address")
	flag.BoolVar(&enableGRPC, "grpc", false, "use gRPC server")
	flag.DurationVar(&shutdownTimeout, "timeout", 0, "shutdown timeout")

	flag.Parse()

	if c.ServerAddressHTTP == "" {
		c.ServerAddressHTTP = serverAddressHTTP
	}
	if c.BaseURLHTTP == "" {
		c.BaseURLHTTP = baseURLHTTP
	}
	if c.ServerAddressGRPC == "" {
		c.ServerAddressGRPC = serverAddressGRPC
	}
	if c.BaseURLGRPC == "" {
		c.BaseURLGRPC = baseURLGRPC
	}
	if c.LogLevel == "" {
		c.LogLevel = logLevel
	}
	if c.DatabaseDSN == "" {
		c.DatabaseDSN = dataBaseDSN
	}
	if c.Config == "" {
		c.Config = config
	}
	if c.TrustedSubnet == "" {
		c.TrustedSubnet = trustedSubnet
	}
	if !c.EnableHTTPS {
		c.EnableHTTPS = enableHTTPS
	}
	if c.ShutdownTimeout == 0 {
		c.ShutdownTimeout = shutdownTimeout
	}
}

type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	dur, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	d.Duration = dur
	return nil
}

type ConfigAlias struct {
	ServerAddressHTTP string   `json:"server_address_HTTP"`
	BaseURLHTTP       string   `json:"base_url_HTTP"`
	ServerAddressGRPC string   `json:"server_address_gRPC"`
	BaseURLGRPC       string   `json:"base_url_gRPC"`
	LogLevel          string   `json:"log_level"`
	DataBaseDSN       string   `json:"database_dsn"`
	EnableHTTPS       bool     `json:"enable_https"`
	CertFilePath      string   `json:"cert_file_path"`
	KeyFilePath       string   `json:"key_file_path"`
	TrustedSubnet     string   `json:"trusted_subnet"`
	ShutdownTimeout   Duration `json:"shutdown_timeout"`
}

// loanFromFile - загрузка конфигурации из файла.
func (c *Config) loanFromFile() error {
	if c.Config == "" {
		return nil
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

	if c.ServerAddressHTTP == "" {
		c.ServerAddressHTTP = configAlias.ServerAddressHTTP
	}
	if c.BaseURLHTTP == "" {
		c.BaseURLHTTP = configAlias.BaseURLHTTP
	}
	if c.ServerAddressGRPC == "" {
		c.ServerAddressGRPC = configAlias.ServerAddressGRPC
	}
	if c.BaseURLGRPC == "" {
		c.BaseURLGRPC = configAlias.BaseURLGRPC
	}
	if c.LogLevel == "" {
		c.LogLevel = configAlias.LogLevel
	}
	if c.DatabaseDSN == "" {
		c.DatabaseDSN = configAlias.DataBaseDSN
	}
	if !c.EnableHTTPS {
		c.EnableHTTPS = configAlias.EnableHTTPS
	}
	if c.CertFilePath == "" {
		c.CertFilePath = configAlias.CertFilePath
	}
	if c.KeyFilePath == "" {
		c.KeyFilePath = configAlias.KeyFilePath
	}
	if c.ShutdownTimeout == 0 {
		c.ShutdownTimeout = configAlias.ShutdownTimeout.Duration
	}

	return nil
}
