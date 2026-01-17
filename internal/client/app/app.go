// Package app. Инициализация зависимостей
// конфигурации, инфраструктурных компонентов и запуск сервера.
package app

import (
	// "context"
	// "os/signal"
	// "syscall"

	"fmt"

	gc "github.com/Di-nis/gophKeeper/internal/client/api/grpc"
	hc "github.com/Di-nis/gophKeeper/internal/client/api/http"
	"github.com/Di-nis/gophKeeper/internal/client/config"
	"github.com/Di-nis/gophKeeper/internal/client/cli"
	"github.com/Di-nis/gophKeeper/pkg/logger"

	"github.com/joho/godotenv"
)

// Start - запуск приложения.
func Start() error {
	// ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	// defer stop()

	if err := godotenv.Load(); err != nil {
		logger.Sugar.Infof("internal/client/app/app.go, func Start(), failed to load env: %v", err)
	}

	config, err := initConfigAndLogger()
	if err != nil {
		return err
	}

	httpClient := hc.New(config.ServerAddress)
	gRPCClient, err := gc.New(config.ServerAddress)
	if err != nil {
		return err
	}

	fmt.Println(httpClient, gRPCClient)

	// TODO: тут должна быть реализация CLI
	cli, err := cli.New(gRPCClient)
	if err != nil {
		return err
	}

	cli.Router()
	return nil

}

// initConfigAndLogger - инициализация конфигурации и логгера.
func initConfigAndLogger() (*config.Config, error) {
	cfg := config.New()
	cfg.Load()

	var err error
	if err = logger.New(cfg.LogLevel); err != nil {
		return nil, err
	}
	logger.Sugar.Infof("Run server, server address: %v", err)
	return cfg, nil
}
