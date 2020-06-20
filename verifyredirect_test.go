package main

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const testServerAddress = "localhost:10099"

//nolint:funlen
func TestVerifyRedirect(t *testing.T) {
	server := newTestServer()

	defer func() {
		err := server.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	testDataTable := map[string]struct {
		source string
		target string
		err    string
	}{
		"source invalid, expect error": {
			source: "",
			target: "http://localhost:10099",
			err:    "unsupported protocol scheme",
		},
		"target invalid, expect error": {
			source: "http://localhost:10099",
			target: "",
			err:    "unsupported protocol scheme",
		},
		"both invalid, expect error": {
			source: "",
			target: "",
			err:    "unsupported protocol scheme",
		},
		"source does not resolve target": {
			source: "http://localhost:10099/test1",
			target: "http://localhost:10099/test4",
			err:    `source uri http://localhost:10099/test1 does resolve to http://localhost:10099/redirect1,not to targetUri http://localhost:10099/test4 which resolves to http://localhost:10099/deadend`, //nolint:lll
		},
		"valid urls, same urls": {
			source: "http://localhost:10099",
			target: "http://localhost:10099",
			err:    "",
		},
		"valid urls, source is resolved in target": {
			source: "http://localhost:10099/test1",
			target: "http://localhost:10099/redirect1",
			err:    "",
		},
		"valid urls, source and target resolve in the same url": {
			source: "http://localhost:10099/test2",
			target: "http://localhost:10099/test3",
			err:    "",
		},
	}

	for testName, testData := range testDataTable {
		t.Run(testName, func(t *testing.T) {
			redirect := Redirect{
				source: testData.source,
				target: testData.target,
			}
			err := VerifyRedirect(redirect)

			if testData.err != "" {
				if err == nil {
					t.Fatalf("Expected error, but got nil")
				}
				assert.Contains(t, err.Error(), testData.err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func newTestServer() *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/test1", http301("redirect1"))
	mux.Handle("/redirect1", httpOkHandler())
	mux.Handle("/target2", httpOkHandler())
	mux.Handle("/test2", http301("target2"))
	mux.Handle("/test3", http301("target2"))
	mux.Handle("/test4", http301("deadend"))
	mux.Handle("/deadend", httpOkHandler())

	srv := &http.Server{Addr: testServerAddress, Handler: mux}

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			fmt.Println(err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	return srv
}

func http301(redirectTarget string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, redirectTarget, http.StatusMovedPermanently)
	}
}

func httpOkHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}
