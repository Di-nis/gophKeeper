// Package main - точка входа в реализацию клиента для сервиса GophKeeper.
package main

import (
	// "context"
	"fmt"
	"log"

	// "time"

	"github.com/Di-nis/gophKeeper/internal/client/config"
	"github.com/Di-nis/gophKeeper/internal/client/usecase"

	gc "github.com/Di-nis/gophKeeper/internal/client/api/grpc"
	hc "github.com/Di-nis/gophKeeper/internal/client/api/http"
	"github.com/Di-nis/gophKeeper/pkg/logger"

	"github.com/Di-nis/gophKeeper/internal/client/cli"

	"github.com/joho/godotenv"
)

var (
	BuildVersion = "N/A"
	BuildTime    = "N/A"
	BuildCommit  = "N/A"
)

func main() {
	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n", BuildVersion, BuildTime, BuildCommit)

	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	// defer cancel()

	if err := godotenv.Load(); err != nil {
		log.Printf("failed to load env: %v", err)
	}

	cfg := config.New()
	cfg.Load()

	var err error
	if err = logger.New(cfg.LogLevel); err != nil {
		log.Fatal("logger initialization error:", err)
	}

	cli := cli.New()
	cli.Parser()

	httpClient := hc.New(cfg.BaseURLHTTP, cfg.TokenStorage)
	gRPCClient, err := gc.New(cfg)
	if err != nil {
		log.Fatal("gRPC client initialization error:", err)
	}

	uc := usecase.New(gRPCClient, httpClient, cli)

	err = uc.Execute()
	if err != nil {
		log.Fatal("request error: ", err)
	}

}
