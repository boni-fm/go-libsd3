package httputil_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/boni-fm/go-libsd3/pkg/httputil"
)

func ExampleRequestIDMiddleware() {
	handler := httputil.RequestIDMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := httputil.GetRequestID(r.Context())
		fmt.Println("request ID present:", id != "")
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	// Output:
	// request ID present: true
}
