package postgres

import (
	"context"
	"fmt"
	"sync"
)

// ConnectionPool adalah struktur yang mengelola kumpulan koneksi database.
// Setiap koneksi diidentifikasi dengan kode DC unik.
type ConnectionPool struct {
	connections map[string]*Database
	configs     map[string]Config
	mu          sync.RWMutex
}

// singleton lifetime untuk kolam renang db nya
var (
	instance *ConnectionPool
	once     sync.Once
)

// GetConnectionPool mengembalikan instance singleton dari ConnectionPool.
// Aman digunakan dari beberapa goroutine secara bersamaan.
func GetConnectionPool() *ConnectionPool {
	once.Do(func() {
		instance = &ConnectionPool{
			connections: make(map[string]*Database),
			configs:     make(map[string]Config),
		}
	})
	return instance
}

// RegisterConfig mendaftarkan konfigurasi database dengan kunci unik.
// Mengembalikan error jika kunci sudah terdaftar.
// Parameter:
//   - kodeDc: kode DC unik sebagai kunci konfigurasi
//   - cfg: konfigurasi database yang akan didaftarkan
func (cp *ConnectionPool) RegisterConfig(kodeDc string, cfg Config) error {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	if _, exists := cp.configs[kodeDc]; exists {
		return fmt.Errorf("configuration with key '%s' already exists", kodeDc)
	}

	cp.configs[kodeDc] = cfg
	return nil
}

// Connect mengembalikan koneksi aktif untuk kode DC yang diberikan.
// Jika koneksi belum ada, koneksi baru akan dibuat secara otomatis menggunakan
// konfigurasi yang sudah didaftarkan lewat RegisterConfig.
// Pemanggil dapat menggunakan defer pool.Close() untuk menutup seluruh pool setelah selesai.
// Menggunakan pola double-checked locking untuk mencegah race condition.
// Parameter:
//   - ctx: context untuk pembatalan operasi
//   - kodeDc: kode DC yang menentukan koneksi mana yang akan dikembalikan
func (cp *ConnectionPool) Connect(ctx context.Context, kodeDc string) (*Database, error) {
	// Fast path: check under read lock
	cp.mu.RLock()
	if conn, exists := cp.connections[kodeDc]; exists {
		cp.mu.RUnlock()
		return conn, nil
	}
	cp.mu.RUnlock()

	// Slow path: acquire write lock and double-check
	cp.mu.Lock()
	defer cp.mu.Unlock()

	// Double-check after acquiring write lock
	if conn, exists := cp.connections[kodeDc]; exists {
		return conn, nil
	}

	// Get the configuration for this key
	cfg, configExists := cp.configs[kodeDc]
	if !configExists {
		return nil, fmt.Errorf("no configuration found for key '%s'", kodeDc)
	}

	// Create new connection while holding the write lock to prevent duplicates
	db, err := NewDatabase(ctx, &cfg)
	if err != nil {
		return nil, fmt.Errorf("error creating database connection for key '%s': %w", kodeDc, err)
	}

	cp.connections[kodeDc] = db
	return db, nil
}

// HasConnection mengecek apakah koneksi dengan kode DC tertentu ada dalam pool.
// Parameter:
//   - kodeDc: kode DC yang akan dicek
func (cp *ConnectionPool) HasConnection(kodeDc string) bool {
	cp.mu.RLock()
	defer cp.mu.RUnlock()

	_, exists := cp.connections[kodeDc]
	return exists
}

// CloseConnection menutup koneksi tertentu dari pool berdasarkan kode DC.
// Parameter:
//   - kodeDc: kode DC koneksi yang akan ditutup
func (cp *ConnectionPool) CloseConnection(kodeDc string) error {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	conn, exists := cp.connections[kodeDc]
	if !exists {
		return fmt.Errorf("no connection found for key '%s'", kodeDc)
	}

	err := conn.Close()
	delete(cp.connections, kodeDc)
	return err
}

// Close menutup semua koneksi dalam pool dan mengosongkan daftar koneksi.
// Metode ini merupakan alias dari CloseAll dan cocok digunakan dengan defer:
//
//	pool := postgres.GetConnectionPool()
//	defer pool.Close()
func (cp *ConnectionPool) Close() error {
	return cp.CloseAll()
}

// CloseAll menutup semua koneksi dalam pool dan mengosongkan daftar koneksi.
func (cp *ConnectionPool) CloseAll() error {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	var errors []error
	for key, conn := range cp.connections {
		if err := conn.Close(); err != nil {
			errors = append(errors, fmt.Errorf("error closing connection '%s': %w", key, err))
		}
	}

	cp.connections = make(map[string]*Database)

	if len(errors) > 0 {
		return fmt.Errorf("multiple errors occurred while closing connections: %v", errors)
	}

	return nil
}

// GetAllKodeDcPoolKey mengembalikan semua kunci kode DC yang terdaftar dalam pool.
func (cp *ConnectionPool) GetAllKodeDcPoolKey() []string {
	cp.mu.RLock()
	defer cp.mu.RUnlock()

	keys := make([]string, 0, len(cp.connections))
	for key := range cp.connections {
		keys = append(keys, key)
	}
	return keys
}

// GetConnectionStats mengembalikan statistik koneksi berdasarkan kode DC.
// Parameter:
//   - ctx: context untuk operasi pengambilan koneksi
//   - kodeDc: kode DC yang statistiknya akan diambil
func (cp *ConnectionPool) GetConnectionStats(ctx context.Context, kodeDc string) (map[string]interface{}, error) {
	conn, err := cp.Connect(ctx, kodeDc)
	if err != nil {
		return nil, err
	}

	stats := conn.GetStats()

	return map[string]interface{}{
		"total_connections":    stats.TotalConns(),
		"idle_connections":     stats.IdleConns(),
		"acquired_connections": stats.AcquiredConns(),
		"construction_slots":   stats.ConstructingConns(),
		"max_size":             stats.MaxConns(),
	}, nil
}

// ReinitializeConnection menutup dan membuat ulang koneksi untuk kode DC tertentu.
// Parameter:
//   - ctx: context untuk operasi
//   - kodeDc: kode DC koneksi yang akan diinisialisasi ulang
func (cp *ConnectionPool) ReinitializeConnection(ctx context.Context, kodeDc string) error {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	if conn, exists := cp.connections[kodeDc]; exists {
		if err := conn.Close(); err != nil {
			return fmt.Errorf("error closing existing connection '%s': %w", kodeDc, err)
		}
		delete(cp.connections, kodeDc)
	}

	cfg, configExists := cp.configs[kodeDc]
	if !configExists {
		return fmt.Errorf("no configuration found for key '%s'", kodeDc)
	}

	db, err := NewDatabase(ctx, &cfg)
	if err != nil {
		return fmt.Errorf("error creating new database connection for key '%s': %w", kodeDc, err)
	}

	cp.connections[kodeDc] = db
	return nil
}

// UpdateConfig memperbarui konfigurasi untuk kode DC tertentu.
// Parameter:
//   - kodeDc: kode DC yang konfigurasinya akan diperbarui
//   - cfg: konfigurasi baru
func (cp *ConnectionPool) UpdateConfig(kodeDc string, cfg Config) {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	cp.configs[kodeDc] = cfg
}

// HealthCheck mengecek kesehatan koneksi untuk kode DC tertentu dengan mengirim ping.
// Parameter:
//   - ctx: context untuk operasi
//   - kodeDc: kode DC koneksi yang akan dicek
func (cp *ConnectionPool) HealthCheck(ctx context.Context, kodeDc string) error {
	conn, err := cp.Connect(ctx, kodeDc)
	if err != nil {
		return err
	}

	return conn.Ping(ctx)
}

// HealthCheckAll mengecek kesehatan semua koneksi dalam pool.
// Mengembalikan map dari kode DC ke error (nil jika sehat).
// Parameter:
//   - ctx: context untuk operasi
func (cp *ConnectionPool) HealthCheckAll(ctx context.Context) map[string]error {
	cp.mu.RLock()
	keys := make([]string, 0, len(cp.connections))
	for kodeDc := range cp.connections {
		keys = append(keys, kodeDc)
	}
	cp.mu.RUnlock()

	results := make(map[string]error)
	for _, kodeDc := range keys {
		results[kodeDc] = cp.HealthCheck(ctx, kodeDc)
	}

	return results
}
