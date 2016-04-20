package dev

import (
	//"fmt"
	"io"
	"net"
	//"strings"
	"net/http"
	"testing"
	"time"

	"github.com/udhos/jazigo/temp"
)

func TestHTTP1(t *testing.T) {

	t.Logf("TestHTTP1: starting")

	// launch bogus test server
	addr := ":2001"
	s, listenErr := spawnServerHTTP(t, addr)
	if listenErr != nil {
		t.Errorf("could not spawn bogus HTTP server: %v", listenErr)
	}
	t.Logf("TestHTTP1: server running on %s", addr)

	// run client test
	logger := &testLogger{t}
	/*
		app := &bogusApp{
			models:  map[string]*Model{},
			devices: map[string]*Device{},
		}
	*/
	app := NewDeviceTable()
	RegisterModels(logger, app)
	CreateDevice(app, logger, "http", "lab1", "localhost"+addr, "", "", "", "")

	repo := temp.TempRepo()
	defer temp.CleanupTempRepo()

	good, bad, skip := ScanDevices(app, logger, 3, 100*time.Millisecond, 200*time.Millisecond, repo, 10, 0)
	if good != 1 || bad != 0 || skip != 0 {
		t.Errorf("good=%d bad=%d skip=%d", good, bad, skip)
	}

	s.close()

	<-s.done // wait termination of accept loop goroutine
}

func spawnServerHTTP(t *testing.T, addr string) (*testServer, error) {

	t.Logf("spawnServerHTTP: will listen on %s", addr)

	http.HandleFunc("/", rootHandler) // default handler

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	t.Logf("spawnServerHTTP: listening on %s", addr)

	s := &testServer{listener: ln, done: make(chan int)}

	go func() {
		err := http.Serve(ln, nil)
		t.Logf("spawnServerHTTP: http.Serve: exited: %v", err)
		close(s.done)
	}()

	return s, nil
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello web client\n")
}