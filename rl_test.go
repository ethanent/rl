package rl

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

type simHTTPResponseWriter struct {
	statusCode int
}

func (s *simHTTPResponseWriter) Header() http.Header {
	return http.Header{}
}

func (s *simHTTPResponseWriter) WriteHeader(c int) {
	s.statusCode = c
}

func (s *simHTTPResponseWriter) Write(d []byte) (int, error) {
	return len(d), nil
}

func simReq(h http.HandlerFunc, req *http.Request) int {
	simRW := &simHTTPResponseWriter{
		statusCode: 200,
	}

	h(simRW, req)

	return simRW.statusCode
}

func TestRl(t *testing.T) {
	count := 0

	testHandler := func(w http.ResponseWriter, r *http.Request) {
		count++
	}

	wrappedHandler := LimitWrap(time.Millisecond * 400, 4, testHandler)

	for i := 0; i < 4; i++ {
		testRw := simReq(wrappedHandler, &http.Request{
			RemoteAddr: "0.0.0.0:57575",
		})

		if testRw != 200 {
			fmt.Println("S1", testRw)
			t.FailNow()
		}
	}

	for i := 0; i < 4; i++ {
		testRw := simReq(wrappedHandler, &http.Request{
			RemoteAddr: "6.6.6.6:74747",
		})

		if testRw != 200 {
			fmt.Println("S2", testRw)
			t.FailNow()
		}
	}

	// Use different port to check for 429 in order to ensure that
	// rl determines IP accurately.

	testRw := simReq(wrappedHandler, &http.Request{
		RemoteAddr: "0.0.0.0:57775",
	})

	if testRw != 429 {
		fmt.Println("S3", testRw)
		t.FailNow()
	}

	time.Sleep(time.Millisecond * 600)

	testRw = simReq(wrappedHandler, &http.Request{
		RemoteAddr: "0.0.0.0:57575",
	})

	if testRw != 200 {
		fmt.Println("S4", testRw)
		t.FailNow()
	}
}
