package between

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
	FromCurrency string `json:"from_currency"`
	ToCurrency string `json:"to_currency"`
	AmountFrom int64 `json:"amount_from"`
}

type Response struct{
	Status string `json:"status"`
}

type BetweenTransaction interface{
	ConvertCurrency(
	ctx context.Context,
	userID uuid.UUID,
	log *slog.Logger,
	fromCurrency, toCurrency string,
	amountFrom int64,
	) error 
}

func New(log *slog.Logger, s BetweenTransaction) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.protected.BetweenTransaction"
		log = log.With(
			slog.String("op",op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		userID, errbool := jwtauth.GetUserIDFromContext(r.Context())
		if errbool == false{
			ErrHandler.ErrHandler(w,r,log, storage.ErrUserNotFound)
			return
		}
		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil{
			ErrHandler.ErrHandler(w,r,log,ErrHandler.ErrFailedDecodeJSON)
			return 
		}
		err = s.ConvertCurrency(
			r.Context(),
			userID,
			log,
			req.FromCurrency,
			req.ToCurrency,
			req.AmountFrom,
		)
		if err != nil{
			ErrHandler.ErrHandler(w,r,log,err)
			return
		}
		render.JSON(w,r,Response{
			Status: "OK",
		})
	}
}