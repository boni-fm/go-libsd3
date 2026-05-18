// Package fileutil menyediakan helper sistem file yang dapat digunakan kembali untuk CSV, ZIP, FTP, dan
// operasi file umum. Semua fungsi menerima afero.Fs untuk kemudahan pengujian —
// gunakan afero.NewMemMapFs() pada pengujian dan afero.NewOsFs() pada produksi.
package fileutil

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/afero"
)

// EnsureDir membuat dirPath dan semua direktori induk jika belum ada.
// Aman dipanggil ketika direktori sudah ada.
func EnsureDir(fs afero.Fs, dirPath string) error {
	if err := fs.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("fileutil/file: ensure dir %q: %w", dirPath, err)
	}
	return nil
}

// FileExists melaporkan apakah path ada dan merupakan file biasa (bukan direktori).
// Mengembalikan false untuk semua error stat.
func FileExists(fs afero.Fs, path string) bool {
	info, err := fs.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// DirExists melaporkan apakah path ada dan merupakan direktori.
// Mengembalikan false untuk semua error stat.
func DirExists(fs afero.Fs, path string) bool {
	info, err := fs.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// CopyFile menyalin src ke dst, membuat direktori induk dst jika diperlukan.
// dst akan ditimpa jika sudah ada. Penyalinan dilakukan secara streaming via io.Copy.
func CopyFile(fs afero.Fs, src, dst string) error {
	srcFile, err := fs.Open(src)
	if err != nil {
		return fmt.Errorf("fileutil/file: open src %q: %w", src, err)
	}
	defer srcFile.Close()

	if err := EnsureDir(fs, aferoDir(dst)); err != nil {
		return err
	}

	dstFile, err := fs.Create(dst)
	if err != nil {
		return fmt.Errorf("fileutil/file: create dst %q: %w", dst, err)
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("fileutil/file: copy %q → %q: %w", src, dst, err)
	}
	return nil
}

// ReadLines membaca file teks UTF-8 dan mengembalikan baris-barisnya sebagai slice string.
// Newline di akhir setiap baris dihapus.
func ReadLines(fs afero.Fs, path string) ([]string, error) {
	f, err := fs.Open(path)
	if err != nil {
		return nil, fmt.Errorf("fileutil/file: open %q: %w", path, err)
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("fileutil/file: read lines from %q: %w", path, err)
	}
	return lines, nil
}

// WriteLines menulis baris ke path yang dipisahkan newline, menimpa jika sudah ada.
// Direktori induk dibuat sesuai kebutuhan.
func WriteLines(fs afero.Fs, path string, lines []string) error {
	if err := EnsureDir(fs, aferoDir(path)); err != nil {
		return err
	}
	f, err := fs.Create(path)
	if err != nil {
		return fmt.Errorf("fileutil/file: create %q: %w", path, err)
	}
	defer f.Close()

	content := strings.Join(lines, "\n")
	if len(lines) > 0 {
		content += "\n"
	}
	if _, err := f.WriteString(content); err != nil {
		return fmt.Errorf("fileutil/file: write lines to %q: %w", path, err)
	}
	return nil
}

// aferoDir mengembalikan komponen direktori dari path menggunakan logika yang memperhatikan pemisah path.
func aferoDir(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' || path[i] == '\\' {
			return path[:i]
		}
	}
	return "."
}
