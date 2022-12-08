package middleware

import (
	"net/http"
	"time"
)

var (
	cacheSince = time.Now().Format(http.TimeFormat)
	cacheUntil = time.Now().AddDate(0, 0, 1).Format(http.TimeFormat)
)

func CacheControlWrapper(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "max-age=86400, must-revalidate, public") // 30 days
		w.Header().Set("Last-Modified", cacheSince)
		w.Header().Set("Expires", cacheUntil)
		next.ServeHTTP(w, r)
	})
}
