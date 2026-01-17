package http_router

import (
	"log/slog"
	"net/http"

	"e-wallet/internal/http_router/handlers"
	"e-wallet/internal/http_router/middleware/logger"
	"e-wallet/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Init_router(log *slog.Logger, s *storage.Storage) http.Handler{
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(logger.New(log))
	r.Use(middleware.URLFormat)
	r.Get("/get_id", handlers.NewGetIdByEmail(log, s))
	return r
}