package httpfileserver

import (
	"net/http"
	"time"
)

type file struct {
	bytes  []byte
	header http.Header
	date   time.Time
}
