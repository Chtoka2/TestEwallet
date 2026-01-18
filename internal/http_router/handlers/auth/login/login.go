package login

import (
	"context"
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

type EnterInterface interface{
	EnterAuth(ctx context.Context, email string, password string) (uuid.UUID, error)
}

func New(log *slog.Logger, s EnterInterface, jwtSvc *jwt.Service) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.auth.login"
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		var req Request
		if err := render.DecodeJSON(r.Body, &req); err != nil{
			log.Error("Failed to decode json")
			render.Status(r, 500)
			render.JSON(w,r,
				Response{
					Status: "Error",
					Error: "Failed to decode json",
				})
			return
		}
		userid, err := s.EnterAuth(r.Context(), req.Email, req.Password)
		if err != nil{
			log.Error("Email or password incorrect")
			render.Status(r, 400)
			render.JSON(w,r,Response{
				Status: "Error",
				Error: "Data incorrect",
			})
			return
		}
		jwttoken, err := jwtSvc.Generate(userid, 15*time.Minute)
		if err != nil{
			render.Status(r, 500)
			render.JSON(w,r,Response{Status: "Error", Error: "Internal server"})
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name: "Authorization",
			Value: jwttoken,
			Path: "/",
			Domain: "",
			MaxAge: 15*60,
			HttpOnly: true,
			Secure: false,
		})
		render.JSON(w,r,Response{Status: "OK", UserID: userid})
	}
}