package middleware

import "net/http"

func Nop(user, pw string) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return h
	}
}
