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
			authHeader := r.Header.Get("Authorization")
			if authHeader == ""{
				http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
				return
			}
			var token string
			if len(authHeader) > 7 && authHeader[:7] == "Bearer "{
				token = authHeader[7:] 
			}else{
				http.Error(w, "invalid Authorization header format", http.StatusUnauthorized)
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