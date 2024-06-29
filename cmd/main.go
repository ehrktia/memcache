package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"github.com/ehrktia/memcache/datastructure"
	"github.com/ehrktia/memcache/server"
	"github.com/ehrktia/memcache/wal"
	"golang.org/x/sync/errgroup"
)

func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGABRT, syscall.SIGBUS, syscall.SIGINT, syscall.SIGQUIT)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// initialize cache and queue
	once := &sync.Once{}
	queueSize := getQueueSize()
	createDataStruct(once, queueSize)
	// start http
	mux := http.NewServeMux()
	w := wal.NewWal()
	walFile := w.WalFileName()
	webServer := server.NewWebServer(w, server.NewHTTPServer())
	webServer.Server.Handler = mux
	registerHandler(mux, webServer)
	// wait for interrupt
	shutdown(sig)
	eg, _ := errgroup.WithContext(ctx)
	// create wal file
	eg.Go(func() error {
		return wal.CreateFile(walFile)
	})
	eg.Go(func() error {
		return wal.Compact(w)
	})
	eg.Go(func() error {
		return webServer.Server.ListenAndServe()
	})
	if err := eg.Wait(); err != nil {
		fmt.Fprintf(os.Stderr, "error:%v\n", err)
		os.Exit(1)
	}
}

func createDataStruct(once *sync.Once, queueSize int) {
	datastructure.NewQueue(once, queueSize)

}

func registerHandler(mux *http.ServeMux, w *server.WebServer) {
	mux.HandleFunc("/save", w.Store)
	mux.HandleFunc("/get", w.Get)
	mux.HandleFunc("/getall", w.GetAll)
}

func shutdown(s chan os.Signal) {
	go func() {
		interrup := <-s
		fmt.Fprintf(os.Stderr, "stopping service -%v\n", fmt.Errorf("%v", interrup))
		os.Exit(0)
	}()

}

func getQueueSize() int {
	size := os.Getenv("QUEUE_SIZE")
	if size == "" {
		size = "50"
	}
	s, err := strconv.Atoi(size)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}

	return s
}
