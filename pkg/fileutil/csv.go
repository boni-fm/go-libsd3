package fileutil

import (
	"encoding/csv"
	"fmt"

	"github.com/gocarina/gocsv"
	"github.com/spf13/afero"
)

// WriteStructsToCSV menulis data (slice struct dengan tag gocsv) ke filePath pada OS filesystem.
// Direktori induk dibuat secara otomatis.
func WriteStructsToCSV[T any](data []T, filePath string) error {
	return WriteStructsToCSVFs[T](afero.NewOsFs(), data, filePath)
}

// WriteStructsToCSVFs adalah versi WriteStructsToCSV yang mendukung afero.
func WriteStructsToCSVFs[T any](fs afero.Fs, data []T, filePath string) error {
	if err := EnsureDir(fs, aferoDir(filePath)); err != nil {
		return fmt.Errorf("fileutil/csv: %w", err)
	}
	f, err := fs.Create(filePath)
	if err != nil {
		return fmt.Errorf("fileutil/csv: create %q: %w", filePath, err)
	}
	defer f.Close()

	if err := gocsv.Marshal(data, f); err != nil {
		return fmt.Errorf("fileutil/csv: marshal structs: %w", err)
	}
	return nil
}

// ReadCSVToStructs membaca filePath dan unmarshal baris-barisnya menjadi slice T menggunakan tag gocsv.
func ReadCSVToStructs[T any](filePath string) ([]T, error) {
	return ReadCSVToStructsFs[T](afero.NewOsFs(), filePath)
}

// ReadCSVToStructsFs adalah versi ReadCSVToStructs yang mendukung afero.
func ReadCSVToStructsFs[T any](fs afero.Fs, filePath string) ([]T, error) {
	f, err := fs.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("fileutil/csv: open %q: %w", filePath, err)
	}
	defer f.Close()

	var out []T
	if err := gocsv.Unmarshal(f, &out); err != nil {
		return nil, fmt.Errorf("fileutil/csv: unmarshal structs: %w", err)
	}
	return out, nil
}

// WriteSliceToCSV menulis baris (baris pertama sebagai header) ke filePath pada OS filesystem.
func WriteSliceToCSV(data [][]string, filePath string) error {
	return WriteSliceToCSVFs(afero.NewOsFs(), data, filePath)
}

// WriteSliceToCSVFs adalah versi WriteSliceToCSV yang mendukung afero.
func WriteSliceToCSVFs(fs afero.Fs, data [][]string, filePath string) error {
	if err := EnsureDir(fs, aferoDir(filePath)); err != nil {
		return fmt.Errorf("fileutil/csv: %w", err)
	}
	f, err := fs.Create(filePath)
	if err != nil {
		return fmt.Errorf("fileutil/csv: create %q: %w", filePath, err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	if err := w.WriteAll(data); err != nil {
		return fmt.Errorf("fileutil/csv: write slice: %w", err)
	}
	return nil
}
