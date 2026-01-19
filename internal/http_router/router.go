package http_router

import (
	"log/slog"
	"net/http"
	"os"

	"e-wallet/internal/http_router/handlers"
	"e-wallet/internal/http_router/handlers/auth"
	"e-wallet/internal/http_router/handlers/protected"
	"e-wallet/internal/http_router/middleware/logger"
	"e-wallet/internal/lib/jwt"
	"e-wallet/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Init_router(log *slog.Logger, s *storage.Storage) http.Handler{
	r := chi.NewRouter()
	//Added jwt service

	secret_key := os.Getenv("JWT_SECRET")
	if secret_key == "" {
		log.Error("JWT_SECRET is not added")
		os.Exit(1)
	}
	jwtSvc := jwt.NewJWTService(secret_key)
	r.Use(middleware.RequestID)
	r.Use(logger.New(log))
	r.Use(middleware.URLFormat)
	r.Mount("/auth", auth.New(log, s, jwtSvc))
	r.Mount("/users", protected.New(log, s, jwtSvc))
	r.Get("/get_id", handlers.NewGetIdByEmail(log, s))
	return r
}