package postgres

import (
	"context"
	"fmt"
)

// ConnectionManager adalah pembungkus praktis untuk mengelola koneksi database via ConnectionPool.
// Cocok digunakan sebagai satu-satunya titik akses koneksi dalam sebuah aplikasi.
type ConnectionManager struct {
	pool *ConnectionPool
}

// NewConnectionManager membuat ConnectionManager baru yang menggunakan singleton ConnectionPool.
func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		pool: GetConnectionPool(),
	}
}

// InitializeConnections menginisialisasi beberapa koneksi database sekaligus dari map konfigurasi.
// Setiap koneksi diverifikasi dengan melakukan ping setelah inisialisasi.
// Parameter:
//   - ctx: context untuk pembatalan operasi
//   - configs: map dari kode DC ke konfigurasi database
func (cm *ConnectionManager) InitializeConnections(ctx context.Context, configs map[string]Config) error {
	for kodeDc, cfg := range configs {
		if err := cm.pool.RegisterConfig(kodeDc, cfg); err != nil {
			return fmt.Errorf("error registering config for kode dc '%s': %w", kodeDc, err)
		}

		// Pre-connect to verify configurations
		if _, err := cm.pool.Connect(ctx, kodeDc); err != nil {
			return fmt.Errorf("error initializing connection for kode dc '%s': %w", kodeDc, err)
		}
	}
	return nil
}

// GetDB mengembalikan koneksi database berdasarkan kode DC.
// Jika koneksi belum ada dan konfigurasi sudah terdaftar, koneksi baru akan dibuat.
// Parameter:
//   - ctx: context untuk pembatalan operasi
//   - kodeDc: kode DC yang menentukan koneksi mana yang akan dikembalikan
func (cm *ConnectionManager) GetDB(ctx context.Context, kodeDc string) (*Database, error) {
	return cm.pool.Connect(ctx, kodeDc)
}

// Close menutup satu koneksi berdasarkan kode DC.
// Parameter:
//   - kodeDc: kode DC koneksi yang akan ditutup
func (cm *ConnectionManager) Close(kodeDc string) error {
	return cm.pool.CloseConnection(kodeDc)
}

// CloseAllConnections menutup semua koneksi yang dikelola oleh ConnectionManager.
// Cocok digunakan dengan defer:
//
//	cm := postgres.NewConnectionManager()
//	defer cm.CloseAllConnections()
func (cm *ConnectionManager) CloseAllConnections() error {
	return cm.pool.CloseAll()
}

// PrintPoolStatus mencetak status semua koneksi dalam pool ke stdout.
// Parameter:
//   - ctx: context untuk operasi pengambilan statistik
func (cm *ConnectionManager) PrintPoolStatus(ctx context.Context) {
	pool := cm.pool
	pool.mu.RLock()
	keys := make([]string, 0, len(pool.connections))
	for kodeDc := range pool.connections {
		keys = append(keys, kodeDc)
	}
	pool.mu.RUnlock()

	fmt.Println("=== Connection Pool Status ===")
	fmt.Printf("Total connections: %d\n", len(keys))

	for _, kodeDc := range keys {
		stats, err := cm.pool.GetConnectionStats(ctx, kodeDc)
		if err == nil {
			fmt.Printf("  - %s: %v\n", kodeDc, stats)
		} else {
			fmt.Printf("  - %s: Error - %v\n", kodeDc, err)
		}
	}
}

// HealthCheck mengecek kesehatan semua koneksi dalam pool.
// Mengembalikan map dari kode DC ke error (nil jika sehat).
// Parameter:
//   - ctx: context untuk operasi
func (cm *ConnectionManager) HealthCheck(ctx context.Context) map[string]error {
	return cm.pool.HealthCheckAll(ctx)
}

// PrintHealthStatus mencetak hasil health check semua koneksi ke stdout.
// Parameter:
//   - ctx: context untuk operasi
func (cm *ConnectionManager) PrintHealthStatus(ctx context.Context) {
	results := cm.HealthCheck(ctx)
	fmt.Println("=== Health Check Status ===")
	for kodeDc, err := range results {
		if err == nil {
			fmt.Printf("  ✓ %s: Healthy\n", kodeDc)
		} else {
			fmt.Printf("  ✗ %s: Unhealthy - %v\n", kodeDc, err)
		}
	}
}
