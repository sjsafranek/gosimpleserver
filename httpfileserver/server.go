package httpfileserver

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
	
	"github.com/sjsafranek/gosimpleserver/cache"
	"github.com/sjsafranek/gosimpleserver/filesystemwatcher"

	"github.com/sjsafranek/logger"
)

type FileServer struct {
	dir     string
	route   string
	cache   cache.ICache[string, file]
	watcher *filesystemwatcher.FileSystemWatcher
}

// New returns a new file server that can handle requests for
// files using an in-memory store with gzipping
// func New(route, dir string, options ...Option) *FileServer {
func New(route, dir string) (*FileServer, error) {

	var server FileServer

	watcher, err := filesystemwatcher.New()
	if nil != err {
		return nil, err
	}

	err = watcher.Add(dir)
	if nil != err {
		return nil, err
	}

	server = FileServer{
		dir:     dir,
		route:   route,
		watcher: watcher,
		cache:   cache.New[string, file](64),
	}

	watcher.OnChange = server.UpdateFile
	watcher.OnDelete = server.DeleteFile
	watcher.OnRename = server.DeleteFile

	return &server, nil
}

func (self *FileServer) Close() error {
	return self.watcher.Close()
}

// Handle gives a handlerfunc for the file server
func (self *FileServer) Handle() http.HandlerFunc {
	return self.ServeHTTP
}

// ServeHTTP is the server of the file server
func (self *FileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	filename := strings.TrimPrefix(r.URL.Path, self.route)

	ctx := r.Context()

	select {

	case <-ctx.Done():
		logger.Warn("request canceled")

	default:

		if "" == filename {
			http.FileServer(http.Dir(self.dir)).ServeHTTP(w, r)
			return
		}

		item, err := self.GetFile(filename)
		if nil != err {
			http.FileServer(http.Dir(self.dir)).ServeHTTP(w, r)
			return
		}

		accept_gzip := strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")
		if accept_gzip {
			w.Header().Set("Content-Encoding", "gzip")
			w.Write(item.bytes)
			return
		}

		buff := bytes.NewBuffer(item.bytes)
		reader, err := gzip.NewReader(buff)
		if nil != err {
			panic(err)
		}

		data, err := ioutil.ReadAll(reader)
		if nil != err {
			panic(err)
		}

		w.Write(data)
	}

}

func (self *FileServer) GetFile(filename string) (*file, error) {
	logger.Tracef("checking for '%v' in cache", filename)
	value := self.cache.Get(filename)
	if nil == value {
		return self.FetchFile(filename)
	}
	return value, nil
}

func (self *FileServer) FetchFile(filename string) (*file, error) {
	if _, err := os.Stat(filename); nil != err {
		return nil, err
	}

	body, err := ioutil.ReadFile(filename)
	if nil != err {
		return nil, err
	}

	var buf bytes.Buffer
	gzw := gzip.NewWriter(&buf)
	gzw.Write(body)
	gzw.Close()

	value := file{
		bytes: buf.Bytes(),
		date:  time.Now(),
	}

	logger.Tracef("adding '%v' to cache", filename)
	self.cache.Set(filename, value)

	return &value, nil
}

func (self *FileServer) DeleteFile(filename string) error {
	logger.Tracef("removing '%v' from cache", filename)
	self.cache.Del(filename)
	return nil
}

func (self *FileServer) UpdateFile(filename string) error {
	logger.Tracef("updating '%v' cache", filename)
	if self.cache.Has(filename) {
		_, err := self.FetchFile(filename)
		return err
	}
	return nil
}
