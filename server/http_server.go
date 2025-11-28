package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"codeberg.org/ehrktia/memcache/datastructure"
	"codeberg.org/ehrktia/memcache/wal"
)

func NewHTTPServer() *http.Server {
	port := os.Getenv("PORT")
	if port == "" {
		fmt.Println("using default port `8080`")
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port)
	s := &http.Server{
		Addr:         addr,
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
	}
	fmt.Println("starting server in port:", addr)
	return s
}

type WebServer struct {
	Wal    *wal.Wal
	Server *http.Server
}

func NewWebServer(w *wal.Wal, h *http.Server) *WebServer {
	return &WebServer{
		Wal:    w,
		Server: h,
	}
}

// Store receives values which are required to be stored
// writes data to wal file
func (w *WebServer) Store(res http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	if err := reqPostMethod(req); err != nil {
		if err := writeResponse(res, []byte(err.Error())); err != nil {
			return
		}
		return
	}
	// get k,v from req body
	buf := &bytes.Buffer{}
	// read req k
	if err := readReqBody(buf, req, res); err != nil {
		return
	}
	// add data to wal
	if err := wal.UpdCache(w.Wal, buf.Bytes()); err != nil {
		if err := writeResponse(res, []byte(err.Error())); err != nil {
			return
		}
	}
	// write response
	if err := writeResponse(res, []byte(fmt.Sprintf("%s\n", "successfully added to cache"))); err != nil {
		return
	}
}

// Get retrieves the value associated with key from in-memory cache store
func (w *WebServer) Get(res http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	if err := reqGetMethod(req); err != nil {
		if err := writeResponse(res, []byte(err.Error())); err != nil {
			return
		}
		return
	}
	buf := &bytes.Buffer{}
	// read req k
	if err := readReqBody(buf, req, res); err != nil {
		return
	}
	d, err := extractReqData(buf)
	if err != nil {
		if err := writeResponse(res, []byte(err.Error())); err != nil {
			return
		}
	}
	// get data from cache
	v := datastructure.Get(d.Key)
	// data not found
	if strings.EqualFold(v.(string), datastructure.NotFound) {
		err := fmt.Errorf("[%v] matching %s", d.Key, datastructure.NotFound)
		if err := writeResponse(res, []byte(err.Error())); err != nil {
			return
		}
	}
	// encode result to response
	resultBytes, err := json.Marshal(v)
	if err != nil {
		if err := writeResponse(res, []byte(err.Error())); err != nil {
			return
		}
	}
	// write response
	if err := writeResponse(res, resultBytes); err != nil {
		return
	}
}

// GetAll emits all the data stored in-memory cache
// this is expensive
func (w *WebServer) GetAll(res http.ResponseWriter, r *http.Request) {
	v := datastructure.GetAll()
	resultBytes, err := json.Marshal(v)
	if err != nil {
		if err := writeResponse(res, []byte(fmt.Sprintf("%s\n", err.Error()))); err != nil {
			return
		}
	}
	if err := writeResponse(res, resultBytes); err != nil {
		return
	}
}
