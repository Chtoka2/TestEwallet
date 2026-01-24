package between

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
	FromCurrency string `json:"from_currency"`
	ToCurrency string `json:"to_currency"`
	AmountFrom int64 `json:"amount_from"`
}

type Response struct{
	Status string `json:"status"`
	Error string `json:"error,omitempty"`
}

type BetweenTransaction interface{
	ConvertCurrency(
	ctx context.Context,
	userID uuid.UUID,
	fromCurrency, toCurrency string,
	amountFrom int64,
	rate float64,
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
			log.Error("Can't find user id")
			render.Status(r, 500)
			render.JSON(w,r, Response{
				Status: "Error",
				Error: "Can't find user id",
			})
		}
		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil{
			log.Error("Can't decode request json")
			render.Status(r, 400)
			render.JSON(w,r, Response{
				Status: "Error",
				Error: "Can't decode request json",
			})
		}
		err = s.ConvertCurrency(r.Context(), userID,
		req.FromCurrency, req.ToCurrency, req.AmountFrom, )
		if err != nil{
			if errors.Is(err, storage.ErrWalletsNotFound){
				log.Error("Can't find wallets")
			}else if errors.Is(err, storage.ErrInsufficientFunds){
				log.Error("In wallet insufficient funds")
			}else if errors.Is(err, storage.ErrSameCurrency){
				log.Error("Can't convert same currency")
			}else if errors.Is(err, storage.ErrInvalidExchangeRate){
				log.Error("Can't find wallets")
			}else if errors.Is(err, storage.ErrInvalidAmount){
				log.Error("Can't find wallets")
			}

		}
	}
}