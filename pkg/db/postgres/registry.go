package postgres

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// registryKey adalah kunci komposit untuk registry multiton.
type registryKey struct {
	kunci  string
	kodedc string
}

// registry adalah registry multiton level paket untuk menyimpan instance *Database.
var (
	dbRegistry    = make(map[registryKey]*Database)
	dbRegistryCfg = make(map[registryKey]Config)
	registryMu    sync.RWMutex
)

// GetDB mengembalikan *Database yang ada untuk pasangan (kunci, kodedc).
// Mengembalikan error jika pasangan kunci tidak ditemukan.
// Parameter:
//   - kunci: kunci identifikasi tenant
//   - kodedc: kode DC database
func GetDB(kunci, kodedc string) (*Database, error) {
	key := registryKey{kunci: kunci, kodedc: kodedc}
	registryMu.RLock()
	db, ok := dbRegistry[key]
	registryMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("GetDB: kunci=%s kodedc=%s tidak ditemukan di registry", kunci, kodedc)
	}
	return db, nil
}

// RegisterDB mendaftarkan *Database baru untuk pasangan (kunci, kodedc).
// Jika kunci yang sama sudah ada dengan konfigurasi yang sama, instance yang sudah ada dikembalikan.
// Jika kunci yang sama sudah ada dengan konfigurasi berbeda, error dikembalikan.
// Menggunakan pola double-checked locking untuk keamanan konkurensi.
// Parameter:
//   - ctx: context untuk pembuatan koneksi database baru
//   - kunci: kunci identifikasi tenant
//   - kodedc: kode DC database
//   - cfg: konfigurasi koneksi database
func RegisterDB(ctx context.Context, kunci, kodedc string, cfg Config) (*Database, error) {
	key := registryKey{kunci: kunci, kodedc: kodedc}

	// Fast path: read lock
	registryMu.RLock()
	if db, ok := dbRegistry[key]; ok {
		existingCfg := dbRegistryCfg[key]
		registryMu.RUnlock()
		if existingCfg != cfg {
			return nil, fmt.Errorf("RegisterDB: kunci=%s kodedc=%s sudah terdaftar dengan konfigurasi berbeda: %w", kunci, kodedc, ErrConfigExists)
		}
		return db, nil
	}
	registryMu.RUnlock()

	// Slow path: write lock with double-check
	registryMu.Lock()
	defer registryMu.Unlock()

	if db, ok := dbRegistry[key]; ok {
		existingCfg := dbRegistryCfg[key]
		if existingCfg != cfg {
			return nil, fmt.Errorf("RegisterDB: kunci=%s kodedc=%s sudah terdaftar dengan konfigurasi berbeda: %w", kunci, kodedc, ErrConfigExists)
		}
		return db, nil
	}

	db, err := NewDatabase(ctx, &cfg)
	if err != nil {
		return nil, fmt.Errorf("RegisterDB: gagal membuat koneksi untuk kunci=%s kodedc=%s: %w", kunci, kodedc, err)
	}

	dbRegistry[key] = db
	dbRegistryCfg[key] = cfg
	return db, nil
}

// MustGetDB mengembalikan *Database untuk pasangan (kunci, kodedc) atau panic jika tidak ditemukan.
// Gunakan fungsi ini hanya jika Anda yakin database sudah terdaftar sebelumnya.
// Parameter:
//   - kunci: kunci identifikasi tenant
//   - kodedc: kode DC database
func MustGetDB(kunci, kodedc string) *Database {
	db, err := GetDB(kunci, kodedc)
	if err != nil {
		panic(err)
	}
	return db
}

// CloseAll menutup setiap *Database yang terdaftar dan mengosongkan registry.
// Mengembalikan error gabungan jika ada koneksi yang gagal ditutup.
func CloseAll() error {
	registryMu.Lock()
	defer registryMu.Unlock()

	var errs []error
	for key, db := range dbRegistry {
		if err := db.Close(); err != nil && !errors.Is(err, ErrConnClose) {
			errs = append(errs, fmt.Errorf("CloseAll: gagal menutup kunci=%s kodedc=%s: %w", key.kunci, key.kodedc, err))
		}
	}

	dbRegistry = make(map[registryKey]*Database)
	dbRegistryCfg = make(map[registryKey]Config)

	if len(errs) > 0 {
		return fmt.Errorf("CloseAll: beberapa koneksi gagal ditutup: %v", errs)
	}
	return nil
}

// injectDB menyuntikkan *Database langsung ke registry untuk keperluan pengujian.
// Fungsi ini tidak diekspor dan hanya digunakan dalam test.
func injectDB(kunci, kodedc string, db *Database, cfg Config) {
	key := registryKey{kunci: kunci, kodedc: kodedc}
	registryMu.Lock()
	defer registryMu.Unlock()
	dbRegistry[key] = db
	dbRegistryCfg[key] = cfg
}
