package login

import (
	"e-wallet/internal/lib/jwt"
	"e-wallet/internal/storage"
	"log/slog"
	"net/http"
)

type Request struct{
	Email, Password string
}

func New(log *slog.Logger, s *storage.Storage, jwtSvc *jwt.Service) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		
	}
}