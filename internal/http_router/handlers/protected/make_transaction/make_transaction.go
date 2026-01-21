package maketransaction

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
	UserIDTo uuid.UUID `json:"user_id"`
	Currency string `json:"currency"`
	Summ int64 `json:"summ"`
}

type Response struct{
	Status string
	Transaction storage.Transaction
	Error string
}

type MakerOfTransaction interface{
	Transactions(ctx context.Context, userIDFrom uuid.UUID,
	userIDTo uuid.UUID, currency string, summ int64) (storage.Transaction, error)
}

func New(log *slog.Logger, s MakerOfTransaction) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.protected.make_transactions"
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
		transaction, err := s.Transactions(r.Context(), userid,
		req.UserIDTo, req.Currency, req.Summ)
		if err != nil{
			if errors.Is(err, storage.ErrWalletWithCurrencyNotFound){
				log.Error("Wallet is not found")
				render.Status(r, 400)
				render.JSON(w,r,Response{
					Status: "Error",
					Error: "Wallet is not found",
				})
				return 
			}
			if errors.Is(err, storage.ErrInsufficientFunds){
				log.Error("Insufficient funds is on wallet")
				render.Status(r, 400)
				render.JSON(w,r,Response{
					Status: "Error",
					Error: "Insufficient funds is on wallet",
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
				Error: "Can't make transaction",
			})
			return 
		}
		log.Info("Transaction made")
		render.JSON(w,r, Response{
			Status: "OK",
			Transaction: transaction,
		})
	}
}