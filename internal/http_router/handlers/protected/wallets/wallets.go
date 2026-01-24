package wallets

import (
	"context"
	jwtauth "e-wallet/internal/http_router/middleware/JWTAuth"
	"e-wallet/internal/lib/ErrHandler"
	"e-wallet/internal/storage"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type WalletGetter interface{
	UserWallets(ctx context.Context, userID uuid.UUID) ([]storage.Wallet, error)
}

type Response struct{
	Status string `json:"status"`
	Wallets []storage.Wallet `json:"wallets,omitempty"`
}

func New(log *slog.Logger, s WalletGetter) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.protected.wallets"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		userid, errbool := jwtauth.GetUserIDFromContext(r.Context())
		if errbool != true{
			ErrHandler.ErrHandler(w,r,log,storage.ErrUserNotFound)
			return
		}
		wallets, err := s.UserWallets(r.Context(), userid)
		if err != nil{
			ErrHandler.ErrHandler(w,r,log,err)
			return
		}
		render.JSON(w,r, Response{
			Status: "OK",
			Wallets: wallets,
		})
	}
}