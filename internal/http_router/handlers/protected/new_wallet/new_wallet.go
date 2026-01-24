package new_wallet

import (
	"context"
	jwtauth "e-wallet/internal/http_router/middleware/JWTAuth"
	"e-wallet/internal/storage"
	"errors"
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
	Error string `json:"error,omitempty"`
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
			log.Error("Can't find userid")
			render.Status(r, 404)
			render.JSON(w,r, Response{
				Status: "Error",
				Error: "Can't find user id",
			})
			return
		}
		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil{
			log.Error("Can't decode json")
			render.Status(r, 400)
			render.JSON(w,r,Response{
				Status: "Error",
				Error: "Can't decode json",
			})
			return 
		}
		err = s.CreateEWallet(r.Context(), userid, req.Currency)
		if err != nil{
			if errors.Is(err, storage.ErrCurrencyWalletExist){
				log.Error("Currency wallet exists")
				render.Status(r, 400)
				render.JSON(w,r,Response{
					Status: "Error",
					Error: "Currency wallet exists",
				})
				return
			}
			if errors.Is(err, storage.ErrCurencyNotInCurrencies){
				log.Error("Currency is not in currencies")
				render.Status(r, 400)
				render.JSON(w,r,Response{
					Status: "Error",
					Error: "Currency is not in currencies",
				})
				return 
			}
			if errors.Is(err, storage.ErrWalletsNotFound){
				log.Error("Wallets is not found")
				render.Status(r, 500)
				render.JSON(w,r,Response{
					Status: "Error",
					Error: "Can't make transaction",
				})
				return
			}
			log.Error("Some problem")
			render.Status(r, 500)
			render.JSON(w,r,Response{
				Status: "Error",
				Error: "Couldn't create new e_wallet",
			})
			return 
		}
		log.Info("Wallet was created")
		render.JSON(w,r, Response{
			Status: "OK",
		})
	}
}