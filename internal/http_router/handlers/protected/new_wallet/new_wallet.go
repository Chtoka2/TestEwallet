package new_wallet

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

type Request struct{
	Currency string `json:"currency"`
}

type Response struct{
	Status string `json:"status"`
}

type CreaterOfWallet interface{
	CreateEWallet(ctx context.Context, UserID uuid.UUID, currency string) (error)
}

func New(log *slog.Logger, s CreaterOfWallet) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.protected.create_wallet"
		log = log.With(
			slog.String("op",op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		userid, errbool := jwtauth.GetUserIDFromContext(r.Context())
		if errbool == false{
			ErrHandler.ErrHandler(w,r,log,storage.ErrUserNotFound)
			return
		}
		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil{
			ErrHandler.ErrHandler(w,r,log,ErrHandler.ErrFailedDecodeJSON)
			return 
		}
		err = s.CreateEWallet(r.Context(), userid, req.Currency)
		if err != nil{
			ErrHandler.ErrHandler(w,r,log,err)
			return 
		}
		log.Info("Wallet was created")
		render.JSON(w,r, Response{
			Status: "OK",
		})
	}
}