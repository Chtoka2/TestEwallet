package regist

import (
	"context"
	"e-wallet/internal/lib/jwt"
	"e-wallet/internal/storage"
	"errors"
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
			log.Error("Failed to decode JSON")
			render.Status(r, 400)
			render.JSON(w,r, Response{Status: "Error", Error: "Incorrect json"})
			return
		}
		userID, err := s.RegistAUTH(r.Context(), req.Email, req.Password)
		if err != nil{
			if errors.Is(err, storage.ErrEmailAlredyExists){
				render.Status(r, http.StatusConflict)
				render.JSON(w,r,Response{Status: "Error", Error: "Email alredy exists"})
				return 
			}
			if errors.Is(err, storage.ErrInvalidInput){
				render.Status(r, 400)
				render.JSON(w,r,Response{Status: "Error", Error: "Invalid input"})
				return 
			}
			render.Status(r, 500)
			render.JSON(w,r, Response{Status: "Error", Error: "Server problem"})
			return
		}
		token, err := jwtSvc.Generate(userID, 15*time.Minute)
		if err != nil{
			render.Status(r,500)
			render.JSON(w,r,Response{Status: "Error", Error: "Internal server"})
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