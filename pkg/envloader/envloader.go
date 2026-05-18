// Package envloader menyediakan helper untuk memuat environment variables dari file .env.
// Berguna untuk development lokal tanpa harus mengatur variabel secara manual.
package envloader

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Load memuat environment variables dari file .env ke os.Environ.
// Variabel yang sudah ada tidak akan ditimpa (overwrite=false).
// Parameter:
//   - filename: path ke file .env yang akan dimuat
func Load(filename string) error {
	return loadFile(filename, false)
}

// LoadOverride memuat environment variables dari file .env dan menimpa yang sudah ada.
// Parameter:
//   - filename: path ke file .env yang akan dimuat
func LoadOverride(filename string) error {
	return loadFile(filename, true)
}

// loadFile adalah implementasi internal untuk memuat file .env.
// Parameter:
//   - filename: path ke file .env
//   - overwrite: jika true, variabel yang sudah ada akan ditimpa
func loadFile(filename string, overwrite bool) error {
	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("envloader: gagal membuka file %s: %w", filename, err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse KEY=VALUE or KEY="VALUE"
		idx := strings.IndexByte(line, '=')
		if idx < 0 {
			continue // skip lines without '='
		}

		key := strings.TrimSpace(line[:idx])
		value := strings.TrimSpace(line[idx+1:])

		if key == "" {
			continue
		}

		// Strip surrounding quotes from value
		value = stripQuotes(value)

		if !overwrite {
			if _, exists := os.LookupEnv(key); exists {
				continue
			}
		}

		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("envloader: gagal set env %s pada baris %d: %w", key, lineNum, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("envloader: error membaca file %s: %w", filename, err)
	}

	return nil
}

// stripQuotes menghapus tanda kutip di awal dan akhir string jika ada.
func stripQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
