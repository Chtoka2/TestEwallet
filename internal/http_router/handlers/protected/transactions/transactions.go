package transactions

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

type Response struct{
	Status string `json:"status"`
	TransactionsMy []storage.Transaction `json:"transactions_my,omitempty"`
	TransactionsOther []storage.Transaction `json:"transactions_other,omitempty"`
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
			ErrHandler.ErrHandler(w,r,log,storage.ErrUserNotFound)
			return 
		}
		transactionsMy, transactionsOther, err := s.GetTransactions(r.Context(), userId)
		if err != nil{
			ErrHandler.ErrHandler(w,r,log,err)
			return
		}
		render.JSON(w,r, Response{
			Status: "OK",
			TransactionsMy: transactionsMy,
			TransactionsOther: transactionsOther,
		})
	}
}