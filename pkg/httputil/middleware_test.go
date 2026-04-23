package httputil

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestIDMiddleware_GeneratesID(t *testing.T) {
	handler := RequestIDMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := GetRequestID(r.Context())
		if id == "" {
			t.Error("expected non-empty request ID in context")
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Header().Get("X-Request-ID") == "" {
		t.Error("expected X-Request-ID header in response")
	}
}

func TestRequestIDMiddleware_UsesExistingID(t *testing.T) {
	const existingID = "my-request-123"

	handler := RequestIDMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := GetRequestID(r.Context())
		if id != existingID {
			t.Errorf("expected %s, got %s", existingID, id)
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("X-Request-ID", existingID)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if got := rr.Header().Get("X-Request-ID"); got != existingID {
		t.Errorf("expected %s in response header, got %s", existingID, got)
	}
}

func TestGetRequestID_EmptyContext(t *testing.T) {
	id := GetRequestID(context.Background())
	if id != "" {
		t.Errorf("expected empty string for empty context, got %s", id)
	}
}

func TestGenerateRequestID_UniqueValues(t *testing.T) {
	id1 := generateRequestID()
	id2 := generateRequestID()
	if id1 == id2 {
		t.Error("expected unique request IDs")
	}
	if len(id1) != 32 {
		t.Errorf("expected 32 char hex string, got length %d", len(id1))
	}
}
