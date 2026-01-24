package maketransaction

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
	UserIDTo uuid.UUID `json:"user_id"`
	Currency string `json:"currency"`
	Summ int64 `json:"summ"`
}

type Response struct{
	Status string
	Transaction storage.Transaction
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
			ErrHandler.ErrHandler(w,r,log,storage.ErrUserNotFound)
			return
		}
		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil{
			ErrHandler.ErrHandler(w,r,log,ErrHandler.ErrFailedDecodeJSON)
			return
		}
		transaction, err := s.Transactions(r.Context(), userid,
		req.UserIDTo, req.Currency, req.Summ)
		if err != nil{
			ErrHandler.ErrHandler(w,r,log,err)
			return 
		}
		log.Info("Transaction made")
		render.JSON(w,r, Response{
			Status: "OK",
			Transaction: transaction,
		})
	}
}