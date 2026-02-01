package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/Di-nis/gophKeeper/internal/server/middleware/compress"
	logger "github.com/Di-nis/gophKeeper/internal/server/middleware/logger/http"
)

// NewRouter - создание роутера.
func NewRouter(handler *Handler) http.Handler {
	r := chi.NewRouter()

	r.Use(logger.WithLogging, compress.GzipMiddleware)

	r.Get("/ping", handler.Ping)
	r.Post("/register", handler.Register)
	r.Post("/login", handler.Login)

	return r

}
