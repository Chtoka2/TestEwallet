package protected

import (
	"e-wallet/internal/storage"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func New(log *slog.Logger, s *storage.Storage) http.Handler{
	r := chi.NewRouter()

}

