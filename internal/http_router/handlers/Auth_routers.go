package handlers

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type IDGetter_by_email interface{
	GetUserIDByEmail(ctx context.Context, email string) (uuid.UUID, error)
}

type Response struct {
	Status string `json:"status"`
	Result string `json:"result,omitempty"`
	Error string `json:"error,omitempty"`
}

func NewGetIdByEmail(log *slog.Logger, id_getter IDGetter_by_email) (http.HandlerFunc){
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "Get url"
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)
		email := r.URL.Query().Get("email")
		if email == ""{
			log.Info("Invalid request")
			render.JSON(w, r, Response{Status: "Error", Error: "Invalid request"})
			return
		}
		resID, err := id_getter.GetUserIDByEmail(r.Context(), email)
		if err != nil{
			log.Info("User not found")
			render.JSON(w,r, Response{Status: "Error", Error: err.Error()})
			return
		}
		log.Info("User id was founded", slog.String("userID", resID.String()))
		render.JSON(w, r, Response{Status: "OK", Result: resID.String()})
	}
}