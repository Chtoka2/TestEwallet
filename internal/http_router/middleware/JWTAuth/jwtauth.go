package jwtauth

import (
	"context"
	"e-wallet/internal/lib/jwt"
	"errors"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

type contextKey string

const UserCtxKey contextKey = "UserID"

func New(log *slog.Logger, jwtSvc *jwt.Service) (func (next http.Handler) http.Handler){
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader, err := r.Cookie("Authorization")
			if err != nil{
				http.Error(w, "Hasn't cookie", 401)
				return 
			}
			token := authHeader.Value
			if token == "" {
  	  			http.Error(w, "empty token", http.StatusUnauthorized)
    			return
			}
			userid, err := jwtSvc.Validate(token)
			if err != nil{
				if errors.Is(err, jwt.ErrExpiredToken) {
					http.Error(w, "token expired", http.StatusUnauthorized)
				} else {
					http.Error(w, "invalid token", http.StatusUnauthorized)
				}
				return
			}
			ctx := context.WithValue(r.Context(), UserCtxKey, userid)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(UserCtxKey).(uuid.UUID)
	return userID, ok
}