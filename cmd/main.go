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
	httpServer := server.NewHTTPServer()
	registerHandler(mux)
	httpServer.Handler = mux
	// wait for interrupt
	shutdown(sig)
	// start server
	eg, _ := errgroup.WithContext(ctx)
	eg.Go(func() error { return httpServer.ListenAndServe() })
	// stop app in case of error
	if err := eg.Wait(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(0)
	}
}

func createDataStruct(once *sync.Once, queueSize int) {
	datastructure.NewQueue(once, queueSize)

}

func registerHandler(mux *http.ServeMux) {
	mux.HandleFunc("/save", server.Store)
	mux.HandleFunc("/get", server.Get)
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
		size = "10"
	}
	s, err := strconv.Atoi(size)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}

	return s
}
