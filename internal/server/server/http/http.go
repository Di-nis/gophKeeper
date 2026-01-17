package http

import (
	"context"
	"net/http"

	"github.com/Di-nis/gophKeeper/internal/server/config"
	"github.com/Di-nis/gophKeeper/pkg/logger"
)

type Server struct {
	httpSrv *http.Server
	config  *config.Config
}

// New - запуск HTTP-сервера.
func New(config *config.Config, routerHandler http.Handler) *Server {
	httpSrv := &http.Server{
		Addr:    config.ServerAddressHTTP,
		Handler: routerHandler,
	}

	return &Server{
		httpSrv: httpSrv,
		config:  config,
	}
}

// Start - запуск HTTP-сервера.
func (s *Server) Start() error {
	var err error

	if s.config.EnableHTTPS {
		if err = s.httpSrv.ListenAndServeTLS(s.config.CertFilePath, s.config.KeyFilePath); err != nil && err != http.ErrServerClosed {
			logger.Sugar.Fatalf("failed start TLS-server: %v", err)
		}
		return nil
	}
	if err = s.httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Sugar.Fatalf("failed start server: %v", err)
	}
	return nil
}

// Stop - остановка HTTP-сервера.
func (s *Server) Stop(ctx context.Context) error {
	return s.httpSrv.Shutdown(ctx)
}
