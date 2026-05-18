package fileutil_test

import (
	"archive/zip"
	"bytes"
	"testing"

	"github.com/boni-fm/go-libsd3/pkg/fileutil"
	"github.com/spf13/afero"
)

func TestZipAndUnzip(t *testing.T) {
	fs := afero.NewMemMapFs()

	_ = afero.WriteFile(fs, "/src/a.txt", []byte("file a"), 0644)
	_ = afero.WriteFile(fs, "/src/b.txt", []byte("file b"), 0644)

	if err := fileutil.ZipFilesFs(fs, "/out/archive.zip", []string{"/src/a.txt", "/src/b.txt"}); err != nil {
		t.Fatalf("ZipFilesFs: %v", err)
	}
	if !fileutil.FileExists(fs, "/out/archive.zip") {
		t.Fatal("expected archive.zip to exist")
	}

	if err := fileutil.UnzipToFs(fs, "/out/archive.zip", "/extracted"); err != nil {
		t.Fatalf("UnzipToFs: %v", err)
	}

	for _, name := range []string{"a.txt", "b.txt"} {
		if !fileutil.FileExists(fs, "/extracted/"+name) {
			t.Errorf("expected /extracted/%s to exist after unzip", name)
		}
	}
}

func TestZipDirectory_Recursive(t *testing.T) {
	fs := afero.NewMemMapFs()

	_ = fs.MkdirAll("/dir/sub", 0755)
	_ = afero.WriteFile(fs, "/dir/root.txt", []byte("root"), 0644)
	_ = afero.WriteFile(fs, "/dir/sub/child.txt", []byte("child"), 0644)

	if err := fileutil.ZipFilesFs(fs, "/out/dir.zip", []string{"/dir"}); err != nil {
		t.Fatalf("ZipFilesFs (dir): %v", err)
	}
	if !fileutil.FileExists(fs, "/out/dir.zip") {
		t.Fatal("expected dir.zip to exist")
	}
}

func TestUnzipTo_ZipSlipRejected(t *testing.T) {
	fs := afero.NewMemMapFs()

	// Create a zip with a path-traversal entry in memory.
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	w, err := zw.Create("../etc/passwd")
	if err != nil {
		t.Fatalf("create zip entry: %v", err)
	}
	_, _ = w.Write([]byte("root:x:0:0"))
	_ = zw.Close()

	// Write the malicious zip to the memfs.
	_ = afero.WriteFile(fs, "/malicious.zip", buf.Bytes(), 0644)

	err = fileutil.UnzipToFs(fs, "/malicious.zip", "/safe")
	if err == nil {
		t.Fatal("expected zip-slip error, got nil")
	}
}
