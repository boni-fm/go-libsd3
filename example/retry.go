package example

// ======================================================
// Contoh pemakaian package retry (pkg/retry)
// ======================================================
//
// Package retry nyediain helper buat ngulang operasi yang gagal
// pake strategi exponential back-off.
//
// Fitur yang di-demo:
//   - DefaultConfig — konfigurasi retry bawaan (3x, 100ms, 10s, 2x)
//   - Config kustom — atur MaxAttempts, BaseDelay, MaxDelay, Multiplier
//   - Do — jalanin fn yang return error
//   - DoWithResult — jalanin fn yang return (value, error)
//   - Penanganan context cancellation
//   - Contoh penggunaan nyata (simulasi HTTP call & DB query retry)

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/boni-fm/go-libsd3/pkg/retry"
)

// ContohRetryDo mendemonstrasikan retry.Do dengan berbagai konfigurasi.
func ContohRetryDo() {
	ctx := context.Background()

	// ─── Konfigurasi default ───
	// MaxAttempts=3, BaseDelay=100ms, MaxDelay=10s, Multiplier=2.0
	defaultCfg := retry.DefaultConfig()
	fmt.Printf("default config: %+v\n", defaultCfg)

	// Simulasi operasi yang gagal dua kali, sukses di percobaan ketiga
	attempt := 0
	err := retry.Do(ctx, defaultCfg, func() error {
		attempt++
		if attempt < 3 {
			return fmt.Errorf("server timeout, percobaan ke-%d", attempt)
		}
		fmt.Printf("sukses di percobaan ke-%d\n", attempt)
		return nil
	})
	if err != nil {
		fmt.Printf("semua retry gagal: %v\n", err)
	}

	// ─── Konfigurasi kustom ───
	customCfg := retry.Config{
		MaxAttempts: 5,
		BaseDelay:   200 * time.Millisecond,
		MaxDelay:    3 * time.Second,
		Multiplier:  1.5, // lebih lambat naik dari default
	}

	attempt2 := 0
	err = retry.Do(ctx, customCfg, func() error {
		attempt2++
		if attempt2 < 4 {
			return fmt.Errorf("koneksi ditolak, percobaan ke-%d", attempt2)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("retry custom gagal: %v\n", err)
	} else {
		fmt.Printf("retry custom sukses di percobaan ke-%d\n", attempt2)
	}

	// ─── Semua retry gagal ───
	errGagalMulu := errors.New("service mati total")
	err = retry.Do(ctx, retry.Config{MaxAttempts: 3, BaseDelay: 10 * time.Millisecond}, func() error {
		return errGagalMulu
	})
	if err != nil {
		fmt.Printf("expected error: %v\n", err)
		// Cek apakah error asli masih bisa di-unwrap
		if errors.Is(err, errGagalMulu) {
			fmt.Println("error asli masih bisa di-unwrap pake errors.Is")
		}
	}
}

// ContohRetryDoWithResult mendemonstrasikan retry.DoWithResult[T].
// Cocok kalo fungsinya return nilai, bukan cuma error.
func ContohRetryDoWithResult() {
	ctx := context.Background()

	cfg := retry.Config{
		MaxAttempts: 4,
		BaseDelay:   50 * time.Millisecond,
		MaxDelay:    1 * time.Second,
		Multiplier:  2.0,
	}

	// Simulasi fetch data dari DB yang sekali-sekali timeout
	callCount := 0
	result, err := retry.DoWithResult[string](ctx, cfg, func() (string, error) {
		callCount++
		if callCount < 3 {
			return "", fmt.Errorf("db query timeout attempt %d", callCount)
		}
		return "data berhasil diambil dari DB", nil
	})

	if err != nil {
		fmt.Printf("gagal semua: %v\n", err)
	} else {
		fmt.Printf("hasil: %q (sukses setelah %d percobaan)\n", result, callCount)
	}

	// Contoh dengan return tipe struct
	type APIResponse struct {
		StatusCode int
		Body       string
	}

	resp, err := retry.DoWithResult[APIResponse](ctx, retry.DefaultConfig(), func() (APIResponse, error) {
		// simulasi panggil API eksternal
		return APIResponse{StatusCode: 200, Body: `{"status":"ok"}`}, nil
	})
	if err != nil {
		fmt.Printf("API call gagal: %v\n", err)
	} else {
		fmt.Printf("API response: status=%d body=%s\n", resp.StatusCode, resp.Body)
	}
}

// ContohRetryContextCancel mendemonstrasikan perilaku retry saat context dibatalkan.
// Retry akan berhenti lebih awal kalo context di-cancel, ga bakal nunggu semua attempt selesai.
func ContohRetryContextCancel() {
	// Context dengan timeout singkat
	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	cfg := retry.Config{
		MaxAttempts: 10,               // banyak attempt
		BaseDelay:   100 * time.Millisecond, // tapi delay per-attempt juga gede
		MaxDelay:    2 * time.Second,
		Multiplier:  2.0,
	}

	err := retry.Do(ctx, cfg, func() error {
		return errors.New("selalu gagal")
	})

	if err != nil {
		fmt.Printf("retry dibatalkan oleh context: %v\n", err)
		// Error message bakal mengandung "context dibatalkan"
	}

	// Context cancel manual
	ctx2, cancel2 := context.WithCancel(context.Background())

	go func() {
		time.Sleep(80 * time.Millisecond)
		fmt.Println("context di-cancel dari goroutine lain")
		cancel2()
	}()

	attempt := 0
	err = retry.Do(ctx2, cfg, func() error {
		attempt++
		return fmt.Errorf("gagal attempt %d", attempt)
	})
	if err != nil {
		fmt.Printf("retry dihentiin paksa: %v (setelah %d attempt)\n", err, attempt)
	}
}

// ContohRetryHTTPCall contoh pemakaian nyata buat retry HTTP call ke external service.
func ContohRetryHTTPCall() {
	ctx := context.Background()

	// Config yang cocok buat HTTP retry ke external API
	httpCfg := retry.Config{
		MaxAttempts: 3,
		BaseDelay:   500 * time.Millisecond,
		MaxDelay:    5 * time.Second,
		Multiplier:  2.0,
	}

	var responseData string
	err := retry.Do(ctx, httpCfg, func() error {
		// Di sini normalnya lo manggil http.Get() atau http.Client.Do()
		// Kalau dapet 5xx atau timeout, return error biar diretry
		// Contoh:
		// resp, err := http.Get("https://api.contoh.id/data")
		// if err != nil { return fmt.Errorf("http error: %w", err) }
		// defer resp.Body.Close()
		// if resp.StatusCode >= 500 { return fmt.Errorf("server error: %d", resp.StatusCode) }
		// ...baca body...

		responseData = `{"id":1,"nama":"Budi"}` // simulasi response sukses
		return nil
	})

	if err != nil {
		fmt.Printf("HTTP call gagal setelah semua retry: %v\n", err)
	} else {
		fmt.Printf("HTTP call sukses: %s\n", responseData)
	}
}
