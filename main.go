package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/sjsafranek/gosimpleserver/httpfileserver"
	"github.com/sjsafranek/gosimpleserver/middleware"
	"github.com/sjsafranek/logger"
)

const (
	DEFAULT_PORT      int    = 8000
	DEFAULT_HOST      string = "0.0.0.0"
	DEFAULT_DIRECTORY string = "."
)

var (
	PORT      int    = DEFAULT_PORT
	HOST      string = DEFAULT_HOST
	DIRECTORY string = DEFAULT_DIRECTORY
)

func main() {
	flag.StringVar(&DIRECTORY, "d", DEFAULT_DIRECTORY, "directory")
	flag.StringVar(&HOST, "h", DEFAULT_HOST, "server host")
	flag.IntVar(&PORT, "p", DEFAULT_PORT, "server port")
	flag.Parse()

	logger.Infof("Serving HTTP on %v port %v", HOST, PORT)

	server, err := httpfileserver.New("/", DIRECTORY)
	if nil != err {
		logger.Error(err)
		os.Exit(1)
	}

	http.Handle("/", middleware.Adapt(server, middleware.RequestIdMiddleWare, middleware.LoggingMiddleWare, middleware.SetHeadersMiddleWare, middleware.CORSMiddleWare))

	err = http.ListenAndServe(fmt.Sprintf("%v:%v", HOST, PORT), nil)
	if err != nil {
		logger.Error("ListenAndServe: ", err)
		os.Exit(1)
	}
}
