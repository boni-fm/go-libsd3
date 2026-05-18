package fileutil

import (
	"archive/zip"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
)

// ZipFiles membuat arsip ZIP di destPath yang berisi semua file dan direktori
// yang terdaftar di sourcePaths (secara rekursif). Menggunakan OS filesystem.
func ZipFiles(destPath string, sourcePaths []string) error {
	return ZipFilesFs(afero.NewOsFs(), destPath, sourcePaths)
}

// ZipFilesFs adalah versi ZipFiles yang mendukung afero.
func ZipFilesFs(fs afero.Fs, destPath string, sourcePaths []string) error {
	if err := EnsureDir(fs, aferoDir(destPath)); err != nil {
		return fmt.Errorf("fileutil/zip: %w", err)
	}
	destFile, err := fs.Create(destPath)
	if err != nil {
		return fmt.Errorf("fileutil/zip: create dest %q: %w", destPath, err)
	}
	defer destFile.Close()

	zw := zip.NewWriter(destFile)
	defer zw.Close()

	for _, src := range sourcePaths {
		if err := addToZipFs(fs, zw, src, ""); err != nil {
			return fmt.Errorf("fileutil/zip: add %q: %w", src, err)
		}
	}
	return nil
}

func addToZipFs(fs afero.Fs, zw *zip.Writer, path, baseInZip string) error {
	info, err := fs.Stat(path)
	if err != nil {
		return err
	}

	if info.IsDir() {
		entries, err := afero.ReadDir(fs, path)
		if err != nil {
			return err
		}
		for _, e := range entries {
			child := path + "/" + e.Name()
			var zipBase string
			if baseInZip == "" {
				zipBase = info.Name()
			} else {
				zipBase = baseInZip + "/" + info.Name()
			}
			if err := addToZipFs(fs, zw, child, zipBase); err != nil {
				return err
			}
		}
		return nil
	}

	var entryName string
	if baseInZip == "" {
		entryName = info.Name()
	} else {
		entryName = baseInZip + "/" + info.Name()
	}

	w, err := zw.Create(entryName)
	if err != nil {
		return err
	}
	f, err := fs.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(w, f)
	return err
}

// UnzipTo mengekstrak arsip ZIP di srcPath ke dalam destDir.
// Dilindungi terhadap serangan zip-slip (path traversal).
func UnzipTo(srcPath, destDir string) error {
	return UnzipToFs(afero.NewOsFs(), srcPath, destDir)
}

// UnzipToFs adalah versi UnzipTo yang mendukung afero.
func UnzipToFs(fs afero.Fs, srcPath, destDir string) error {
	srcFile, err := fs.Open(srcPath)
	if err != nil {
		return fmt.Errorf("fileutil/zip: open src %q: %w", srcPath, err)
	}
	defer srcFile.Close()

	info, err := fs.Stat(srcPath)
	if err != nil {
		return fmt.Errorf("fileutil/zip: stat src %q: %w", srcPath, err)
	}

	// zip.NewReader needs an io.ReaderAt and the size.
	type readerAt interface {
		io.Reader
		io.ReaderAt
	}
	ra, ok := srcFile.(readerAt)
	if !ok {
		return fmt.Errorf("fileutil/zip: source file does not implement io.ReaderAt")
	}

	zr, err := zip.NewReader(ra, info.Size())
	if err != nil {
		return fmt.Errorf("fileutil/zip: open zip reader: %w", err)
	}

	// Normalise destDir so our prefix check works reliably on all platforms.
	cleanDest := filepath.Clean(destDir)
	// Use slash-normalised paths for the prefix check so mixed separators
	// (e.g. on Windows) cannot bypass the protection.
	cleanDestSlash := filepath.ToSlash(cleanDest) + "/"

	for _, f := range zr.File {
		// Zip-slip protection: normalise the target and check it is inside destDir.
		target := filepath.Join(cleanDest, f.Name)
		targetSlash := filepath.ToSlash(filepath.Clean(target)) + "/"
		if !strings.HasPrefix(targetSlash, cleanDestSlash) {
			return fmt.Errorf("fileutil/zip: zip-slip detected: %s", f.Name)
		}

		if f.FileInfo().IsDir() {
			if err := fs.MkdirAll(target, 0755); err != nil {
				return fmt.Errorf("fileutil/zip: mkdir %q: %w", target, err)
			}
			continue
		}

		if err := EnsureDir(fs, filepath.Dir(target)); err != nil {
			return fmt.Errorf("fileutil/zip: %w", err)
		}

		outFile, err := fs.Create(target)
		if err != nil {
			return fmt.Errorf("fileutil/zip: create %q: %w", target, err)
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return fmt.Errorf("fileutil/zip: open entry %q: %w", f.Name, err)
		}

		_, copyErr := io.Copy(outFile, rc)
		rc.Close()
		outFile.Close()
		if copyErr != nil {
			return fmt.Errorf("fileutil/zip: extract %q: %w", f.Name, copyErr)
		}
	}
	return nil
}
