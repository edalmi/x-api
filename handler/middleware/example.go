package middleware

import (
	"log"
	"net/http"
)

func Nop(user, pw string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println()
			next.ServeHTTP(w, r)
		})
	}
}
