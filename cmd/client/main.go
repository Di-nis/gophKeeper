package main

import (
	"fmt"
	"log"

	"github.com/Di-nis/gophKeeper/internal/client/config"
	"github.com/Di-nis/gophKeeper/internal/client/usecase"

	gc "github.com/Di-nis/gophKeeper/internal/client/api/grpc"
	hc "github.com/Di-nis/gophKeeper/internal/client/api/http"
	"github.com/Di-nis/gophKeeper/pkg/logger"

	"github.com/Di-nis/gophKeeper/internal/client/cli"

	"github.com/joho/godotenv"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n", buildVersion, buildDate, buildCommit)

	if err := godotenv.Load(); err != nil {
		log.Fatal("failed to load env: %w", err)
	}

	cfg := config.New()
	cfg.Load()

	var err error
	if err = logger.New(cfg.LogLevel); err != nil {
		log.Fatal("logger initialization error: %w", err)
	}

	command := cli.Parser()


	httpClient := hc.New(cfg.ServerAddressHTTP)
	gRPCClient, err := gc.New(cfg.ServerAddressGRPC)
	if err != nil {
		log.Fatal("gRPC client initialization error: %w", err)
	}

	uc := usecase.New(gRPCClient, httpClient, *command)

	uc.Execute()

}
