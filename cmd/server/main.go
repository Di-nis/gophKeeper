// Package main - точка входа в приложение GophKeeper.
// @title GophKeeper
// @version 1.0
// @description API для хранения секретов

// @contact.name Denis Smirnov
// @contact.email di-nis@ya.ru

// @host localhost:8080
// @BasePath /

package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"github.com/Di-nis/gophKeeper/internal/server/app"
	"github.com/Di-nis/gophKeeper/internal/server/config"
	grpcHandler "github.com/Di-nis/gophKeeper/internal/server/delivery/grpc"
	httpHandler "github.com/Di-nis/gophKeeper/internal/server/delivery/http"

	dataUsecase "github.com/Di-nis/gophKeeper/internal/server/usecase/data"
	userUsecase "github.com/Di-nis/gophKeeper/internal/server/usecase/user"
	"github.com/Di-nis/gophKeeper/pkg/logger"

	grpcServer "github.com/Di-nis/gophKeeper/internal/server/server/grpc"
	httpServer "github.com/Di-nis/gophKeeper/internal/server/server/http"

	"github.com/joho/godotenv"
)

var (
	buildVersion = "N/A"
	BuildTime    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n", buildVersion, BuildTime, buildCommit)

	var err error

	if err := godotenv.Load(); err != nil {
		log.Fatal("failed to load env:", err)
	}

	cfg := config.New()
	if err = cfg.Load(); err != nil {
		log.Fatal("config error loading:", err)
	}

	if err := logger.New(cfg.LogLevel); err != nil {
		log.Fatal("logger initialization error:", err)
	}

	repo, err := app.InitRepoPostgres(cfg)
	if err != nil {
		log.Fatal("repo initialization error:", err)
	}

	dataUc := dataUsecase.New(repo)

	dataHandler := grpcHandler.New(dataUc, cfg)
	grpcSrv, err := grpcServer.New(cfg, dataHandler)
	if err != nil {
		logger.Sugar.Fatalf("internal/server/app/app.go, func Start(), failed to start grpc-server: %v", err)
	}

	userUc := userUsecase.New(repo, cfg)
	pingUc, authUc := userUc, userUc

	userHandler := httpHandler.New(pingUc, authUc, cfg)
	router := httpHandler.NewRouter(userHandler)

	httpSrv := httpServer.New(cfg, router)

	application := app.New(grpcSrv, httpSrv)

	application.Start()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	application.Stop(shutdownCtx)

	<-shutdownCtx.Done()
}
