package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"
	//"os"

	"github.com/sjsafranek/gosimpleserver/utils"
)

// var (
// 	hostname string
// )

// func init() {
// 	name, err := os.Hostname()
// 	if err != nil {
// 		panic(err)
// 	}
// 	hostname = name
// }

func NewRequestId() string {
	return fmt.Sprintf("%x:%x", utils.GetGoRoutineID(), time.Now().UnixNano())
}

func RequestIdMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := NewRequestId()

		ctx := r.Context()
		req := r.WithContext(context.WithValue(ctx, "key", "val"))
		*r = *req

		w.Header().Set("X-RequestId", id)
		next.ServeHTTP(w, r)
	})
}
