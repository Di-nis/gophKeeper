package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/Di-nis/gophKeeper/internal/server/middleware/compress"
	logger "github.com/Di-nis/gophKeeper/internal/server/middleware/logger/http"

	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/Di-nis/gophKeeper/docs"
)

// @title Chi API
// @version 1.0
// @description Swagger для chi router
// @host localhost:8080
// @BasePath /
// NewRouter - создание роутера.
func NewRouter(handler *Handler) http.Handler {
	r := chi.NewRouter()

	r.Use(logger.WithLogging, compress.GzipMiddleware)

	r.Get("/ping", handler.Ping)
	r.Post("/register", handler.Register)
	r.Post("/login", handler.Login)

	// swagger
	r.Handle("/swagger/*", httpSwagger.WrapHandler)

	return r

}
