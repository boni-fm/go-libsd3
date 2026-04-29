package example

// ======================================================
// Contoh pemakaian package yaml (pkg/yaml)
// ======================================================
//
// Package yaml nyediain helper untuk baca konfigurasi dari file YAML
// menggunakan Viper di balik layar, dengan cache thread-safe.
//
// Fitur yang di-demo:
//   - ReadConfigDynamic — baca seluruh config file ke map
//   - ReadConfigDynamicWithKey — baca nilai spesifik berdasarkan key
//   - GetKunciConfigFilepath — cari path config.yaml secara otomatis
//   - ClearCache — bersihkan cache (berguna di testing)
//   - Penanganan error: file tidak ada, key tidak ada, value null

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/boni-fm/go-libsd3/pkg/yaml"
)

// ContohReadConfigDynamic mendemonstrasikan ReadConfigDynamic.
// Membaca SELURUH isi config file ke dalam map[string]interface{}.
func ContohReadConfigDynamic() {
	// Buat file config YAML sementara
	tmpDir, err := os.MkdirTemp("", "yaml_demo_*")
	if err != nil {
		fmt.Printf("gagal buat tmpdir: %v\n", err)
		return
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "config.yaml")
	configContent := `
kunci: kunciG009SIM
environment: production
database:
  host: db.internal
  port: 5432
  name: mydb
kafka:
  brokers:
    - kafka-1:9092
    - kafka-2:9092
  topic: my-topic
app:
  port: 8080
  debug: false
  name: my-api-service
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		fmt.Printf("gagal buat config: %v\n", err)
		return
	}

	// ─── ReadConfigDynamic — baca semua key sekaligus ───
	allConfig, err := yaml.ReadConfigDynamic(configPath)
	if err != nil {
		fmt.Printf("ReadConfigDynamic gagal: %v\n", err)
		return
	}

	fmt.Println("seluruh konfigurasi:")
	for k, v := range allConfig {
		fmt.Printf("  %s = %v\n", k, v)
	}

	// Akses nested config
	if dbSection, ok := allConfig["database"].(map[string]interface{}); ok {
		fmt.Printf("database host: %v\n", dbSection["host"])
		fmt.Printf("database port: %v\n", dbSection["port"])
	}

	// ─── Cache: panggil kedua kalinya gratis (dari cache) ───
	allConfig2, err := yaml.ReadConfigDynamic(configPath)
	if err == nil {
		fmt.Printf("dari cache: %d key terbaca\n", len(allConfig2))
	}
}

// ContohReadConfigDynamicWithKey mendemonstrasikan pembacaan nilai per key.
// Lebih efisien kalau cuma butuh satu nilai.
func ContohReadConfigDynamicWithKey() {
	// Buat file config YAML sementara
	tmpDir, _ := os.MkdirTemp("", "yaml_key_*")
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "config.yaml")
	os.WriteFile(configPath, []byte(`
kunci: kunciG010SIM
app_name: super-service
port: 9090
debug: true
max_connections: 50
`), 0644)

	// ─── Baca nilai string ───
	kunci, err := yaml.ReadConfigDynamicWithKey(configPath, "kunci")
	if err != nil {
		fmt.Printf("gagal baca kunci: %v\n", err)
	} else {
		fmt.Printf("kunci: %v (tipe: %T)\n", kunci, kunci)
	}

	// ─── Baca nilai boolean ───
	debug, err := yaml.ReadConfigDynamicWithKey(configPath, "debug")
	if err != nil {
		fmt.Printf("gagal baca debug: %v\n", err)
	} else {
		fmt.Printf("debug: %v (tipe: %T)\n", debug, debug)
	}

	// ─── Baca nilai integer ───
	port, err := yaml.ReadConfigDynamicWithKey(configPath, "port")
	if err != nil {
		fmt.Printf("gagal baca port: %v\n", err)
	} else {
		fmt.Printf("port: %v (tipe: %T)\n", port, port)
	}

	// ─── Key tidak ada → ErrNullKeyValue ───
	_, err = yaml.ReadConfigDynamicWithKey(configPath, "key_yang_ga_ada")
	if err != nil {
		fmt.Printf("expected error key ga ada: %v\n", err)
	}

	// ─── Bersihkan cache setelah selesai (berguna di testing) ───
	yaml.ClearCache()
	fmt.Println("cache dibersihkan")
}

// ContohGetKunciConfigFilepath mendemonstrasikan pencarian otomatis config.yaml.
//
// GetKunciConfigFilepath mencari config.yaml di:
//  1. Direktori executable
//  2. Current working directory
//  3. Path dari env var CONFIG_PATH
func ContohGetKunciConfigFilepath() {
	// Kalau config.yaml ada di salah satu lokasi yang dicari,
	// fungsi ini langsung return path absolutnya

	// ─── Set CONFIG_PATH env var buat demo ───
	tmpDir, _ := os.MkdirTemp("", "yaml_path_*")
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "config.yaml")
	os.WriteFile(configPath, []byte("kunci: kunciG009SIM\n"), 0644)

	// Set env var CONFIG_PATH
	os.Setenv("CONFIG_PATH", tmpDir)
	defer os.Unsetenv("CONFIG_PATH")

	foundPath, err := yaml.GetKunciConfigFilepath()
	if err != nil {
		fmt.Printf("config.yaml tidak ditemukan: %v\n", err)
	} else {
		fmt.Printf("config.yaml ditemukan di: %s\n", foundPath)
	}

	// ─── Tanpa env var, mencari di CWD ───
	os.Unsetenv("CONFIG_PATH")

	// Buat config.yaml di CWD sementara
	cwdConfig := "config.yaml"
	os.WriteFile(cwdConfig, []byte("kunci: test\n"), 0644)
	defer os.Remove(cwdConfig)

	foundPath2, err := yaml.GetKunciConfigFilepath()
	if err != nil {
		fmt.Printf("tidak ketemu di CWD: %v\n", err)
	} else {
		fmt.Printf("ketemu di CWD: %s\n", foundPath2)
	}
}

// ContohYAMLErrorHandling mendemonstrasikan error handling yang dihasilkan package yaml.
func ContohYAMLErrorHandling() {
	// ─── File tidak ada ───
	_, err := yaml.ReadConfigDynamic("/path/ga/ada/config.yaml")
	if err != nil {
		fmt.Printf("expected error file ga ada: %v\n", err)
	}

	_, err = yaml.ReadConfigDynamicWithKey("/path/ga/ada/config.yaml", "key")
	if err != nil {
		fmt.Printf("expected error dengan key: %v\n", err)
	}

	// ─── File kosong → ErrNullConfigValue ───
	tmpFile, _ := os.CreateTemp("", "empty_*.yaml")
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	_, err = yaml.ReadConfigDynamic(tmpFile.Name())
	if err != nil {
		fmt.Printf("expected error file kosong: %v\n", err)
	}

	// ─── Key tidak ada → ErrNullKeyValue ───
	tmpDir, _ := os.MkdirTemp("", "yaml_err_*")
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "config.yaml")
	os.WriteFile(configPath, []byte("ada_key: ada_value\n"), 0644)

	_, err = yaml.ReadConfigDynamicWithKey(configPath, "ga_ada_key")
	if err != nil {
		fmt.Printf("expected error key ga ada: %v\n", err)
	}
}

// ContohYAMLClearCache mendemonstrasikan penggunaan ClearCache.
// Berguna di testing buat reset state antar test case.
func ContohYAMLClearCache() {
	tmpDir, _ := os.MkdirTemp("", "yaml_cache_*")
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "config.yaml")
	os.WriteFile(configPath, []byte("version: 1\n"), 0644)

	// Load pertama → masuk ke cache
	v1, _ := yaml.ReadConfigDynamicWithKey(configPath, "version")
	fmt.Printf("versi pertama: %v\n", v1)

	// Update file
	os.WriteFile(configPath, []byte("version: 2\n"), 0644)

	// Load kedua → masih dari cache (version: 1)
	v2, _ := yaml.ReadConfigDynamicWithKey(configPath, "version")
	fmt.Printf("versi dari cache (harusnya 1): %v\n", v2)

	// ClearCache → baca ulang dari disk
	yaml.ClearCache()

	// Load ketiga → fresh dari disk (version: 2)
	v3, _ := yaml.ReadConfigDynamicWithKey(configPath, "version")
	fmt.Printf("versi setelah clear cache (harusnya 2): %v\n", v3)
}
