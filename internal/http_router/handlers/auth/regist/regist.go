package regist

import (
	"context"
	"e-wallet/internal/lib/ErrHandler"
	"e-wallet/internal/lib/jwt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type Request struct{
	Email string `json:"email"`
	Password string `json:"password"`
}

type Response struct{
	Status string `json:"status"`
	UserID uuid.UUID `json:"user_id,omitempty"`
	Error string `json:"error,omitempty"`
}

type RegistAuth interface{
	RegistAUTH(ctx context.Context, email string, password string) (uuid.UUID, error)
}

func New(log *slog.Logger, s RegistAuth, jwtSvc jwt.JWTGeneratorInterface) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		op := "handlers.Auth.RegistHandler"
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var req Request
		if err := render.DecodeJSON(r.Body, &req); err != nil{
			ErrHandler.ErrHandler(w,r,log,ErrHandler.ErrFailedDecodeJSON)
			return
		}
		userID, err := s.RegistAUTH(r.Context(), req.Email, req.Password)
		if err != nil{
			ErrHandler.ErrHandler(w,r,log,err)
			return
		}
		token, err := jwtSvc.Generate(userID, 15*time.Minute)
		if err != nil{
			ErrHandler.ErrHandler(w,r,log,ErrHandler.ErrFailedCodeJWT)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name: "Authorization",
			Value: token,
			Path: "/",
			Domain: "",
			MaxAge: 15*60,
			HttpOnly: true,
			Secure: true,
		})
		render.JSON(w,r,Response{Status: "OK", UserID: userID})
	}
}