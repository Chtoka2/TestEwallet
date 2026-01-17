package jwtauth

import (
	"log/slog"
	"net/http"
)

func New(log *slog.Logger) (func (next http.Handler) http.Handler){
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request){
			
		}
		return http.HandlerFunc(fn)
	}	
}