package login

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
	Error string `json:"error"`
}

type EnterInterface interface{
	EnterAuth(context.Context, string, string) (uuid.UUID, error)
}

func New(log *slog.Logger, s EnterInterface, jwtSvc jwt.JWTGeneratorInterface) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.auth.login"
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var req Request
		if err := render.DecodeJSON(r.Body, &req); err != nil{
			ErrHandler.ErrHandler(w,r,log, ErrHandler.ErrFailedDecodeJSON)
			return
		}
		userid, err := s.EnterAuth(r.Context(), req.Email, req.Password)
		if err != nil{
			ErrHandler.ErrHandler(w,r,log,err)
			return
		}
		jwttoken, err := jwtSvc.Generate(userid, 15*time.Minute)
		if err != nil{
			ErrHandler.ErrHandler(w,r,log, ErrHandler.ErrFailedCodeJSON)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name: "Authorization",
			Value: jwttoken,
			Path: "/",
			Domain: "",
			MaxAge: 15*60,
			HttpOnly: true,
			Secure: true,
		})
		render.JSON(w,r,Response{Status: "OK", UserID: userid})
	}
}