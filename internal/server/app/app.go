// Package app. Инициализация зависимостей
// конфигурации, инфраструктурных компонентов и запуск сервера.
package app

import (
	"context"

	"github.com/Di-nis/gophKeeper/internal/server/config"
	"github.com/Di-nis/gophKeeper/pkg/logger"

	"github.com/Di-nis/gophKeeper/internal/server/repository/postgres"
)

// Runner - интерфейс для запуска и остановки сервера.
type Runner interface {
	Start() error
	Stop(ctx context.Context) error
}

// App - приложение.
type App struct {
	servers []Runner
}

// New - конструктор по созданию приложения.
func New(servers ...Runner) *App {
	return &App{servers: servers}
}

// Start - запуск приложения.
func (a *App) Start() {
	for _, s := range a.servers {
		go func(s Runner) {
			if err := s.Start(); err != nil {
				logger.Sugar.Infof("internal/server/app/app.go,  App.Start(), start server error: %v", err)
			}
		}(s)
	}
}

// Stop - остановка приложения.
func (a *App) Stop(ctx context.Context) {
	for _, s := range a.servers {
		if err := s.Stop(ctx); err != nil {
			logger.Sugar.Infof("internal/server/app/app.go,  App.Start(), stop server error: %v", err)
		}
	}
}

// InitRepoPostgres - инициализация репозитория для работы с PostgreSQL.
func InitRepoPostgres(config *config.Config) (*postgres.Repo, error) {
	db, err := postgres.InitDB(config.DatabaseDSN)
	if err != nil {
		return nil, err
	}

	repo, err := postgres.New(db)
	if err != nil {
		return nil, err
	}

	err = repo.Migrations()
	if err != nil {
		return nil, err
	}

	return repo, nil
}
