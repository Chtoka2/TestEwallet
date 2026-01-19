package protected

import (
	"e-wallet/internal/http_router/handlers/protected/wallets"
	jwtauth "e-wallet/internal/http_router/middleware/JWTAuth"
	"e-wallet/internal/lib/jwt"
	"e-wallet/internal/storage"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func New(log *slog.Logger, s *storage.Storage, jwtSvc *jwt.Service) http.Handler{
	r := chi.NewRouter()
	r.Use(jwtauth.New(log, jwtSvc))
	r.Get("/wallets", wallets.New(log, s))
	return r
}