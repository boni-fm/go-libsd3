// Package retry menyediakan helper untuk mengulang operasi yang gagal dengan strategi exponential back-off.
// Berguna untuk operasi jaringan atau database yang mungkin gagal sementara.
package retry

import (
	"context"
	"fmt"
	"math"
	"time"
)

// Config menyimpan konfigurasi untuk mekanisme retry.
type Config struct {
	// MaxAttempts adalah jumlah maksimal percobaan (termasuk percobaan pertama).
	MaxAttempts int
	// BaseDelay adalah waktu tunggu awal sebelum percobaan ulang.
	BaseDelay time.Duration
	// MaxDelay adalah batas maksimum waktu tunggu antar percobaan.
	MaxDelay time.Duration
	// Multiplier adalah pengali untuk exponential back-off.
	Multiplier float64
}

// DefaultConfig mengembalikan konfigurasi retry bawaan yang umum digunakan.
// MaxAttempts=3, BaseDelay=100ms, MaxDelay=10s, Multiplier=2.0
func DefaultConfig() Config {
	return Config{
		MaxAttempts: 3,
		BaseDelay:   100 * time.Millisecond,
		MaxDelay:    10 * time.Second,
		Multiplier:  2.0,
	}
}

// Do menjalankan fn dengan mekanisme retry sesuai konfigurasi c.
// Mengembalikan error terakhir jika semua percobaan gagal.
// Menghormati pembatalan context — akan berhenti jika context dibatalkan.
// Parameter:
//   - ctx: context untuk pembatalan operasi
//   - c: konfigurasi retry
//   - fn: fungsi yang akan dieksekusi ulang jika gagal
func Do(ctx context.Context, c Config, fn func() error) error {
	if c.MaxAttempts <= 0 {
		c.MaxAttempts = 1
	}

	var lastErr error
	for attempt := 0; attempt < c.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("retry: context dibatalkan sebelum percobaan ke-%d: %w", attempt+1, err)
		}

		lastErr = fn()
		if lastErr == nil {
			return nil
		}

		if attempt < c.MaxAttempts-1 {
			delay := calculateDelay(c, attempt)
			select {
			case <-ctx.Done():
				return fmt.Errorf("retry: context dibatalkan saat menunggu percobaan ke-%d: %w", attempt+2, ctx.Err())
			case <-time.After(delay):
			}
		}
	}

	return fmt.Errorf("retry: semua %d percobaan gagal, error terakhir: %w", c.MaxAttempts, lastErr)
}

// DoWithResult menjalankan fn dan mengembalikan hasil generik dengan mekanisme retry.
// Mengembalikan nilai zero dan error jika semua percobaan gagal.
// Parameter:
//   - ctx: context untuk pembatalan operasi
//   - c: konfigurasi retry
//   - fn: fungsi yang akan dieksekusi ulang jika gagal, mengembalikan nilai dan error
func DoWithResult[T any](ctx context.Context, c Config, fn func() (T, error)) (T, error) {
	var zero T
	var lastErr error

	if c.MaxAttempts <= 0 {
		c.MaxAttempts = 1
	}

	for attempt := 0; attempt < c.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return zero, fmt.Errorf("retry: context dibatalkan sebelum percobaan ke-%d: %w", attempt+1, err)
		}

		result, err := fn()
		if err == nil {
			return result, nil
		}
		lastErr = err

		if attempt < c.MaxAttempts-1 {
			delay := calculateDelay(c, attempt)
			select {
			case <-ctx.Done():
				return zero, fmt.Errorf("retry: context dibatalkan saat menunggu percobaan ke-%d: %w", attempt+2, ctx.Err())
			case <-time.After(delay):
			}
		}
	}

	return zero, fmt.Errorf("retry: semua %d percobaan gagal, error terakhir: %w", c.MaxAttempts, lastErr)
}

// calculateDelay menghitung waktu tunggu untuk percobaan berikutnya menggunakan exponential back-off.
func calculateDelay(c Config, attempt int) time.Duration {
	delay := float64(c.BaseDelay) * math.Pow(c.Multiplier, float64(attempt))
	if delay > float64(c.MaxDelay) {
		delay = float64(c.MaxDelay)
	}
	return time.Duration(delay)
}
