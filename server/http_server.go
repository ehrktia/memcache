package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
	"strings"

	"github.com/ehrktia/memcache/datastructure"
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


func Store(res http.ResponseWriter, req *http.Request) {
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
	// extract req data
	d, err := extractReqData(buf)
	if err != nil {
		if err := writeResponse(res, []byte(err.Error())); err != nil {
			return
		}
	}

	result, _ := datastructure.Add(d.Key, d.Value)
	if result != nil {
		_, err := json.Marshal(result)
		if err != nil {
			if err := writeResponse(res, []byte(err.Error())); err != nil {
				return
			}
		}
		if err := writeResponse(res, []byte("successfully added to cache")); err != nil {
			return
		}

	}
}

func Get(res http.ResponseWriter, req *http.Request) {
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
	v:= datastructure.Get(d.Key)
	// data not found
	if strings.EqualFold(v.(string),datastructure.NotFound){
		err := fmt.Errorf("%v matching value not found", d.Key)
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
