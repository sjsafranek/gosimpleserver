package middleware

import (
	"net/http"
)

func CORSMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.WriteHeader(200)
			return
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
