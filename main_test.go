package main

import (
	"testing"
	"os"
	"net/http"
	"time"
)

var originals = struct {
	osArgs []string
}{
	osArgs: os.Args,
}

func restoreOriginals() {
	os.Args = originals.osArgs
}

// This test does not make assertions, it's just to setup a test scenario
func TestMainFunc(t *testing.T) {
	defer restoreOriginals()

	server := startServer()

	os.Args = []string{"","-c=4"}
	main()

	err := server.Close()
	if err != nil {
		t.Fatalf("Did not expect server.Close to return an error, but got: %v ", err)
	}
}

func startServer() *http.Server {

	mux := http.NewServeMux()
	mux.Handle("/test1", http301("redirect1"))
	mux.Handle("/redirect1", httpOkHandler())
	mux.Handle("/test2", httpOkHandler())

	srv := &http.Server{Addr: "localhost:10099", Handler: mux}

	go srv.ListenAndServe()

	time.Sleep(100 * time.Millisecond)

	return srv
}

func http301(redirectTarget string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Second)
		http.Redirect(w, r, redirectTarget, http.StatusMovedPermanently)
	})
}

func httpOkHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	})
}
