// Package httputil menyediakan middleware HTTP untuk keperluan tracing dan logging request.
// Berguna untuk melacak setiap request yang masuk ke service.
package httputil

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

// contextKey adalah tipe untuk kunci context agar tidak konflik dengan paket lain.
type contextKey string

// RequestIDKey adalah kunci context untuk menyimpan dan mengambil request-id.
const RequestIDKey contextKey = "request-id"

// RequestIDMiddleware menambahkan request-id unik ke setiap HTTP request.
// Request-id diambil dari header X-Request-ID jika ada, atau dibuat baru jika tidak.
// Request-id disimpan dalam context request dan dikembalikan dalam header X-Request-ID response.
// Parameter:
//   - next: handler HTTP berikutnya dalam chain
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}

		ctx := context.WithValue(r.Context(), RequestIDKey, requestID)
		w.Header().Set("X-Request-ID", requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetRequestID mengambil request-id dari context.
// Mengembalikan string kosong jika tidak ditemukan.
// Parameter:
//   - ctx: context yang mungkin mengandung request-id
func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(RequestIDKey).(string); ok {
		return id
	}
	return ""
}

// generateRequestID membuat request-id baru secara acak menggunakan 16 byte hex.
func generateRequestID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "unknown"
	}
	return hex.EncodeToString(b)
}
