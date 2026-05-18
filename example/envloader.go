package example

// ======================================================
// Contoh pemakaian package envloader (pkg/envloader)
// ======================================================
//
// Package envloader buat muat environment variables dari file .env.
// Berguna banget di development lokal biar ga perlu set env manual.
//
// Fitur yang di-demo:
//   - Load — muat .env, variabel yang sudah ada TIDAK ditimpa
//   - LoadOverride — muat .env, variabel yang sudah ada DITIMPA
//   - Format yang didukung: KEY=VALUE, KEY="VALUE", KEY='VALUE'
//   - Komentar (#) dan baris kosong diabaikan
//   - Error handling kalo file ga ketemu

import (
	"fmt"
	"os"

	"github.com/boni-fm/go-libsd3/pkg/envloader"
)

// ContohEnvloaderLoad mendemonstrasikan envloader.Load().
// Variabel yang sudah ada di env TIDAK akan ditimpa (safe for existing env).
func ContohEnvloaderLoad() {
	// Buat file .env sementara buat demo
	envContent := `
# ini komentar, diabaikan aja

APP_NAME=my-service
APP_PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_NAME="my_database"
DB_USER='postgres'
DB_PASS="rahasia123"

# baris kosong juga diabaikan

KAFKA_BROKERS=localhost:9092
KAFKA_TOPIC=my-topic
`
	tmpFile, err := os.CreateTemp("", "demo_*.env")
	if err != nil {
		fmt.Printf("gagal buat file temp: %v\n", err)
		return
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(envContent); err != nil {
		fmt.Printf("gagal tulis ke temp: %v\n", err)
		return
	}
	tmpFile.Close()

	// Pastiin variabel belum ada sebelum load
	os.Unsetenv("APP_NAME")
	os.Unsetenv("APP_PORT")
	os.Unsetenv("DB_HOST")

	// Load — muat env dari file
	if err := envloader.Load(tmpFile.Name()); err != nil {
		fmt.Printf("Load gagal: %v\n", err)
		return
	}

	// Sekarang variabel udah keisi
	fmt.Println("APP_NAME:", os.Getenv("APP_NAME"))   // → my-service
	fmt.Println("APP_PORT:", os.Getenv("APP_PORT"))   // → 8080
	fmt.Println("DB_HOST:", os.Getenv("DB_HOST"))     // → localhost
	fmt.Println("DB_NAME:", os.Getenv("DB_NAME"))     // → my_database (tanpa quotes)
	fmt.Println("DB_USER:", os.Getenv("DB_USER"))     // → postgres (tanpa quotes)
	fmt.Println("DB_PASS:", os.Getenv("DB_PASS"))     // → rahasia123 (tanpa quotes)

	// Load TIDAK menimpa yang sudah ada
	os.Setenv("APP_NAME", "sudah-ada-duluan")

	if err := envloader.Load(tmpFile.Name()); err != nil {
		fmt.Printf("Load kedua gagal: %v\n", err)
		return
	}

	// APP_NAME harusnya tetap "sudah-ada-duluan", bukan "my-service"
	fmt.Println("APP_NAME setelah load kedua:", os.Getenv("APP_NAME")) // → sudah-ada-duluan
}

// ContohEnvloaderLoadOverride mendemonstrasikan envloader.LoadOverride().
// Variabel yang sudah ada AKAN ditimpa oleh nilai dari file .env.
func ContohEnvloaderLoadOverride() {
	tmpFile, err := os.CreateTemp("", "demo_override_*.env")
	if err != nil {
		fmt.Printf("gagal buat file temp: %v\n", err)
		return
	}
	defer os.Remove(tmpFile.Name())

	tmpFile.WriteString("APP_ENV=production\nLOG_LEVEL=warn\n")
	tmpFile.Close()

	// Set nilai awal
	os.Setenv("APP_ENV", "development")
	os.Setenv("LOG_LEVEL", "debug")

	fmt.Println("SEBELUM override:")
	fmt.Println("APP_ENV:", os.Getenv("APP_ENV"))   // → development
	fmt.Println("LOG_LEVEL:", os.Getenv("LOG_LEVEL")) // → debug

	// LoadOverride — TIMPA yang sudah ada
	if err := envloader.LoadOverride(tmpFile.Name()); err != nil {
		fmt.Printf("LoadOverride gagal: %v\n", err)
		return
	}

	fmt.Println("SESUDAH override:")
	fmt.Println("APP_ENV:", os.Getenv("APP_ENV"))    // → production
	fmt.Println("LOG_LEVEL:", os.Getenv("LOG_LEVEL")) // → warn
}

// ContohEnvloaderErrorHandling mendemonstrasikan error handling envloader.
func ContohEnvloaderErrorHandling() {
	// File tidak ada → error
	err := envloader.Load("/path/yang/tidak/ada/.env")
	if err != nil {
		fmt.Printf("expected error file ga ada: %v\n", err)
	}

	// LoadOverride juga sama
	err = envloader.LoadOverride("/nonexistent/.env")
	if err != nil {
		fmt.Printf("expected error dari LoadOverride: %v\n", err)
	}
}

// ContohEnvloaderFormatYangDukung mendemonstrasikan berbagai format yang bisa dibaca.
func ContohEnvloaderFormatYangDukung() {
	tmpFile, err := os.CreateTemp("", "format_*.env")
	if err != nil {
		fmt.Printf("gagal buat temp: %v\n", err)
		return
	}
	defer os.Remove(tmpFile.Name())

	// Berbagai format yang didukung
	content := `
# ── komentar pakai hash ──
TANPA_QUOTES=nilai_langsung
DENGAN_DOUBLE_QUOTES="nilai dengan spasi"
DENGAN_SINGLE_QUOTES='nilai pakai single'
DENGAN_SPASI_AWAL=   nilai dengan spasi sebelum dan sesudah trim   

# baris ini bukan key=value jadi diabaikan saja
ini-bukan-env-var

ANGKA=42
URL=http://localhost:8080/api/v1

# key kosong diabaikan
=tidak_valid
`
	tmpFile.WriteString(content)
	tmpFile.Close()

	// Bersiin dulu
	for _, k := range []string{"TANPA_QUOTES", "DENGAN_DOUBLE_QUOTES", "DENGAN_SINGLE_QUOTES",
		"DENGAN_SPASI_AWAL", "ANGKA", "URL"} {
		os.Unsetenv(k)
	}

	if err := envloader.Load(tmpFile.Name()); err != nil {
		fmt.Printf("Load gagal: %v\n", err)
		return
	}

	fmt.Println("TANPA_QUOTES:", os.Getenv("TANPA_QUOTES"))             // → nilai_langsung
	fmt.Println("DENGAN_DOUBLE_QUOTES:", os.Getenv("DENGAN_DOUBLE_QUOTES")) // → nilai dengan spasi
	fmt.Println("DENGAN_SINGLE_QUOTES:", os.Getenv("DENGAN_SINGLE_QUOTES")) // → nilai pakai single
	fmt.Println("ANGKA:", os.Getenv("ANGKA"))                            // → 42
	fmt.Println("URL:", os.Getenv("URL"))                                // → http://localhost:8080/api/v1
}
