// Package main - точка входа в реализацию клиента для сервиса GophKeeper.
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Di-nis/gophKeeper/internal/client/config"
	"github.com/Di-nis/gophKeeper/internal/client/usecase"
	"github.com/Di-nis/gophKeeper/pkg/logger"

	gc "github.com/Di-nis/gophKeeper/internal/client/api/grpc"
	hc "github.com/Di-nis/gophKeeper/internal/client/api/http"
	repository "github.com/Di-nis/gophKeeper/internal/client/repository/local"

	"github.com/joho/godotenv"
)

var (
	BuildVersion = "N/A"
	BuildTime    = "N/A"
	BuildCommit  = "N/A"
)

func main() {
	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n", BuildVersion, BuildTime, BuildCommit)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := godotenv.Load(); err != nil {
		log.Printf("failed to load env: %v", err)
	}

	cfg := config.New()
	if err := cfg.Load(); err != nil {
		log.Fatal("config loading error:", err)
	}

	var err error
	if err = logger.New(cfg.LogLevel); err != nil {
		log.Fatal("logger initialization error:", err)
	}

	httpClient := hc.New(cfg)
	gRPCClient, err := gc.New(cfg)
	if err != nil {
		log.Fatal("gRPC client initialization error:", err)
	}

	repo, err := repository.New(cfg.DatabasePath)
	if err != nil {
		log.Fatal("repo initialization error:", err)
	}

	uc, err := usecase.New(cfg, httpClient, gRPCClient, repo)
	if err != nil {
		log.Fatal("service initialization error` ", err)
	}

	err = uc.Execute(ctx)
	if err != nil {
		log.Fatal("request error: ", err)
	}

	defer cancel()

	httpClient.Client.CloseIdleConnections()
	if err = gRPCClient.Close(); err != nil {
		log.Fatal("gRPC client close error: ", err)
	}

	if err = repo.Close(); err != nil {
		log.Fatal("repo close error: ", err)
	}
}
