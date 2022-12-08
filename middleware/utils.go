package middleware

import (
	"net/http"
)

// func Attach(next http.Handler) http.Handler {
// 	return LoggingMiddleWare(SetHeadersMiddleWare(CORSMiddleWare(next)))
// }

// Adapter wraps an http.Handler with additional
// functionality.
type Adapter func(http.Handler) http.Handler

// Adapt h with all specified adapters.
func Adapt(handler http.Handler, adapters ...Adapter) http.Handler {
	for _, adapter := range adapters {
		handler = adapter(handler)
	}
	return handler
}
