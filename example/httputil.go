package example

// ======================================================
// Contoh pemakaian package httputil (pkg/httputil)
// ======================================================
//
// Package httputil nyediain middleware HTTP dan helper untuk tracing request.
//
// Fitur yang di-demo:
//   - RequestIDMiddleware — tambah X-Request-ID ke setiap request
//   - GetRequestID — ambil request-id dari context
//   - Chaining middleware dengan handler lain
//   - Penggunaan bersama net/http standard library

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/boni-fm/go-libsd3/pkg/httputil"
)

// ContohRequestIDMiddleware mendemonstrasikan RequestIDMiddleware.
//
// Middleware ini:
//   - Ambil X-Request-ID dari header request kalau ada
//   - Generate random ID baru kalau headernya kosong
//   - Simpen request-id ke context
//   - Set X-Request-ID ke header response
func ContohRequestIDMiddleware() {
	// Handler aplikasi utama
	appHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ambil request-id dari context yang udah diset oleh middleware
		reqID := httputil.GetRequestID(r.Context())
		fmt.Printf("request masuk, id: %s\n", reqID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"request_id": reqID,
			"pesan":      "berhasil",
		})
	})

	// Bungkus dengan RequestIDMiddleware
	handlerDenganMiddleware := httputil.RequestIDMiddleware(appHandler)

	// Pasang ke HTTP server standard Go
	mux := http.NewServeMux()
	mux.Handle("/api/data", handlerDenganMiddleware)
	mux.Handle("/api/users", handlerDenganMiddleware)

	// Jalankan server di port 8080
	// Uncomment baris di bawah untuk benar-benar menjalankan server:
	// fmt.Println("server jalan di :8080")
	// http.ListenAndServe(":8080", mux)

	fmt.Println("middleware RequestID terdaftar di /api/data dan /api/users")
	_ = handlerDenganMiddleware // hindari compile error saat demo
}

// ContohGetRequestID mendemonstrasikan GetRequestID dari context.
//
// GetRequestID aman dipanggil kapan aja — kalau request-id ga ada di context,
// bakal return string kosong tanpa panik.
func ContohGetRequestID() {
	// Context kosong — GetRequestID return ""
	emptyCtx := context.Background()
	id1 := httputil.GetRequestID(emptyCtx)
	fmt.Printf("request-id dari context kosong: %q\n", id1) // → ""

	// Context dengan request-id (biasanya diset oleh RequestIDMiddleware)
	ctxDenganID := context.WithValue(context.Background(), httputil.RequestIDKey, "trace-abc-123-xyz")
	id2 := httputil.GetRequestID(ctxDenganID)
	fmt.Printf("request-id dari context berisi: %q\n", id2) // → "trace-abc-123-xyz"
}

// ContohMiddlewareChaining mendemonstrasikan cara chain beberapa middleware.
//
// Urutan: RequestIDMiddleware → LoggingMiddleware (custom) → Handler
func ContohMiddlewareChaining() {
	// Middleware logging custom (contoh, bukan dari library)
	loggingMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqID := httputil.GetRequestID(r.Context())
			fmt.Printf("[LOG] %s %s | request_id=%s\n", r.Method, r.URL.Path, reqID)
			next.ServeHTTP(w, r)
		})
	}

	// Middleware auth custom (contoh)
	authMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Authorization")
			if token == "" {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}

	// Handler akhir
	apiHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := httputil.GetRequestID(r.Context())
		fmt.Fprintf(w, `{"request_id":"%s","status":"ok"}`, reqID)
	})

	// Chain: request-id → logging → auth → handler
	// Urutan di-apply dari dalam ke luar (terakhir dipasang = pertama dieksekusi)
	finalHandler := httputil.RequestIDMiddleware(
		loggingMiddleware(
			authMiddleware(apiHandler),
		),
	)

	mux := http.NewServeMux()
	mux.Handle("/api/protected", finalHandler)

	fmt.Println("middleware chain terdaftar: RequestID → Logging → Auth → Handler")
	_ = mux
}

// ContohHeaderXRequestID mendemonstrasikan bagaimana X-Request-ID dari client
// diteruskan dan dikembalikan di response.
//
// Kalau client kirim X-Request-ID di header request:
//   → middleware pakai ID dari client tersebut (untuk distributed tracing)
//
// Kalau client TIDAK kirim X-Request-ID:
//   → middleware generate ID baru secara otomatis
func ContohHeaderXRequestID() {
	dieksekusi := false

	// Simulasi handler yang verifikasi request-id
	handler := httputil.RequestIDMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqID := httputil.GetRequestID(r.Context())
			fmt.Printf("handler menerima request_id: %s\n", reqID)

			// Response juga punya X-Request-ID header
			responseID := w.Header().Get("X-Request-ID")
			fmt.Printf("X-Request-ID di response header: %s\n", responseID)

			dieksekusi = true
		}),
	)

	// Simulasi request dengan X-Request-ID dari client
	req1, _ := http.NewRequest("GET", "/api/test", nil)
	req1.Header.Set("X-Request-ID", "client-generated-id-123")

	// Simulasi request tanpa X-Request-ID (middleware akan generate)
	req2, _ := http.NewRequest("GET", "/api/test", nil)

	// Pakai ResponseRecorder buat simulasi (tanpa test package)
	fmt.Println("── Request 1 (dengan X-Request-ID dari client) ──")
	handler.ServeHTTP(&dummyResponseWriter{header: http.Header{}}, req1)

	fmt.Println("── Request 2 (tanpa X-Request-ID, auto-generate) ──")
	handler.ServeHTTP(&dummyResponseWriter{header: http.Header{}}, req2)

	_ = dieksekusi
}

// dummyResponseWriter adalah implementasi http.ResponseWriter minimal buat demo.
type dummyResponseWriter struct {
	header http.Header
	status int
}

func (d *dummyResponseWriter) Header() http.Header         { return d.header }
func (d *dummyResponseWriter) Write(b []byte) (int, error) { return len(b), nil }
func (d *dummyResponseWriter) WriteHeader(statusCode int)  { d.status = statusCode }
