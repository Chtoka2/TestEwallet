package auth

import (
	"e-wallet/internal/http_router/handlers/auth/login"
	"e-wallet/internal/http_router/handlers/auth/regist"
	"e-wallet/internal/lib/jwt"
	"e-wallet/internal/storage"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func New(log *slog.Logger, s *storage.Storage, jwtSvc *jwt.Service) http.Handler{
	r := chi.NewRouter()
	r.Post("/register", regist.New(log, s, jwtSvc)) //regist handler
	r.Post("/login", login.New(log, s, jwtSvc)) //login handler
	return r
}