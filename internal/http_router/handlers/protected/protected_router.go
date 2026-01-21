package protected

import (
	maketransaction "e-wallet/internal/http_router/handlers/protected/make_transaction"
	"e-wallet/internal/http_router/handlers/protected/transactions"
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
	r.Get("/transactions", transactions.New(log, s))
	r.Post("/make_transaction", maketransaction.New(log, s))
	return r
}