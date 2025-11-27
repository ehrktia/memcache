package server

import (
	"fmt"
	"net"
	"sync"
	"testing"
)

func TestHostName(t *testing.T) {
	got, err := hostname()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(got)
}

func TestWriteHostName(t *testing.T) {
	// start a test TCP server

	buf := make([]byte, 512)
	wg := &sync.WaitGroup{}
	wg.Add(2)
	ln, err := net.Listen("tcp", ":9999")
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		defer func() { wg.Done() }()
		srvConn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		if n, err := srvConn.Read(buf); err != nil || n < 1 {
			panic(err)
		}
	}()
	go func() {
		defer func() {
			wg.Done()
		}()
		rAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%s", "9999"))
		if err != nil {
			panic(err)
		}
		conn, err := net.DialTCP("tcp", nil, rAddr)
		if err != nil {
			panic(err)
		}
		h, err := hostname()
		if err != nil {
			panic(err)
		}
		if err := write(conn, h, 10); err != nil {
			panic(err)
		}
	}()
	wg.Wait()
	t.Logf("buf:%s\n", buf)
}
