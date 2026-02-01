// Package main - точка входа в реализацию клиента для сервиса GophKeeper.
package main

import (
	"context"
	"fmt"
	"log"

	"time"

	"github.com/Di-nis/gophKeeper/internal/client/config"
	"github.com/Di-nis/gophKeeper/internal/client/usecase"

	gc "github.com/Di-nis/gophKeeper/internal/client/api/grpc"
	hc "github.com/Di-nis/gophKeeper/internal/client/api/http"
	"github.com/Di-nis/gophKeeper/pkg/logger"

	in "github.com/Di-nis/gophKeeper/internal/client/cli/input"
	out "github.com/Di-nis/gophKeeper/internal/client/cli/output"

	"github.com/joho/godotenv"
)

var (
	BuildVersion = "N/A"
	BuildTime    = "N/A"
	BuildCommit  = "N/A"
)

func main() {
	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n", BuildVersion, BuildTime, BuildCommit)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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

	// repo, err := app.InitRepoPostgres(cfg)
	// if err != nil {
	// 	log.Fatal("repo initialization error:", err)
	// }

	input := in.New()
	output := out.New()
	input.Parser()

	httpClient := hc.New(cfg.BaseURLHTTP, cfg.TokenStorage)
	gRPCClient, err := gc.New(cfg)
	if err != nil {
		log.Fatal("gRPC client initialization error:", err)
	}

	uc := usecase.New(gRPCClient, httpClient, input, output)

	err = uc.Execute(ctx)
	if err != nil {
		log.Fatal("request error: ", err)
	}

	uc.Output.Print(err)

}
