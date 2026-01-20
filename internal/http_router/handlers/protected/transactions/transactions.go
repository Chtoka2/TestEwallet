package transactions

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

type Response struct{
	Status string `json:"status"`
	TransactionsMy []storage.Transaction `json:"transactions_my,omitempty"`
	TransactionsOther []storage.Transaction `json:"transactions_other,omitempty"`
	Error string `json:"error,omitempty"`
}

type TransactionGetter interface{
	GetTransactions(
		ctx context.Context,
		userID uuid.UUID)(
		[]storage.Transaction,
		[]storage.Transaction,
		error)
}

func New(log *slog.Logger, s TransactionGetter)(http.HandlerFunc){
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.protected.transactions"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		userId, errbool := jwtauth.GetUserIDFromContext(r.Context())
		if errbool != true{
			log.Error("Can't read jwt")
			render.Status(r, 500)
			render.JSON(w,r,Response{
				Status: "Error",
				Error: "Can't read json",
			})
			return 
		}
		transactionsMy, transactionsOther, err := s.GetTransactions(r.Context(), userId)
		if err != nil{
			if errors.Is(err, storage.ErrWalletsNotFound){
				log.Error("Can't find wallets")
				render.Status(r, 404)
				render.JSON(w,r, Response{
					Status: "Error",
					Error: "Can't find wallets",
				})
				return
			}
			log.Error("Some problem", slog.String("Error", err.Error()))
			render.Status(r, 500)
			render.JSON(w,r, Response{
				Status: "Error",
				Error: "Internal server",
			})
			return 
		}
		render.JSON(w,r, Response{
			Status: "OK",
			TransactionsMy: transactionsMy,
			TransactionsOther: transactionsOther,
		})
	}
}