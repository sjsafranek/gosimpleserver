package middleware

import (
	"net/http"
	"time"

	"github.com/sjsafranek/logger"
)

// https://ndersson.me/post/capturing_status_code_in_net_http/
// https://upgear.io/blog/golang-tip-wrapping-http-response-writer-for-middleware/
// https://www.reddit.com/r/golang/comments/7p35s4/how_do_i_get_the_response_status_for_my_middleware/
type statusWriter struct {
	http.ResponseWriter
	status int
	length int
}

func (self *statusWriter) WriteHeader(code int) {
	self.status = code
	self.ResponseWriter.WriteHeader(code)
}

func (self *statusWriter) Write(b []byte) (int, error) {
	if self.status == 0 {
		self.status = 200
	}
	n, err := self.ResponseWriter.Write(b)
	self.length += n
	return n, err
}

func (self *statusWriter) Flush() {
	self.ResponseWriter.(http.Flusher).Flush()
}

func LoggingMiddleWare(next http.Handler) http.Handler {
	log := logger.New()
	log.SetName("HTTPAPI")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log.Debugf("IN  %v %v %v - %v bytes", r.RemoteAddr, r.Method, r.URL.Path, r.ContentLength)

		start := time.Now()

		// Event Source streams
		if 0 != len(r.Header["Accept"]) {
			if "text/event-stream" == r.Header["Accept"][0] {
				next.ServeHTTP(w, r)
				log.Debugf("OUT %v %v %v [%v] %v", r.RemoteAddr, r.Method, r.URL, 200, time.Since(start))
				return
			}
		}

		// handle websockets...
		for _, header := range r.Header["Upgrade"] {
			if header == "websocket" {
				next.ServeHTTP(w, r)
				log.Debugf("OUT %v %v %v [%v] %v", r.RemoteAddr, r.Method, r.URL, 200, time.Since(start))
				return
			}
		}
		//.end

		// Initialize the status to 200 in case WriteHeader is not called
		sw := statusWriter{w, 200, 0}

		// during
		next.ServeHTTP(&sw, r)

		// end
		log.Debugf("OUT %v %v %v [%v] %v - %v bytes", r.RemoteAddr, r.Method, r.URL, sw.status, time.Since(start), sw.length)
	})
}
