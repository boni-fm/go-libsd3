package fileutil_test

import (
	"testing"

	"github.com/boni-fm/go-libsd3/pkg/fileutil"
	"github.com/spf13/afero"
)

func TestEnsureDir_CreatesNested(t *testing.T) {
	fs := afero.NewMemMapFs()
	if err := fileutil.EnsureDir(fs, "/a/b/c"); err != nil {
		t.Fatalf("EnsureDir: %v", err)
	}
	if !fileutil.DirExists(fs, "/a/b/c") {
		t.Error("directory /a/b/c should exist")
	}
}

func TestFileExists_TrueForFile_FalseForDir_FalseForMissing(t *testing.T) {
	fs := afero.NewMemMapFs()
	_ = afero.WriteFile(fs, "/f.txt", []byte("x"), 0644)
	_ = fs.MkdirAll("/d", 0755)

	if !fileutil.FileExists(fs, "/f.txt") {
		t.Error("expected FileExists true for /f.txt")
	}
	if fileutil.FileExists(fs, "/d") {
		t.Error("expected FileExists false for directory /d")
	}
	if fileutil.FileExists(fs, "/missing") {
		t.Error("expected FileExists false for missing path")
	}
}

func TestDirExists_TrueForDir_FalseForFile_FalseForMissing(t *testing.T) {
	fs := afero.NewMemMapFs()
	_ = fs.MkdirAll("/d", 0755)
	_ = afero.WriteFile(fs, "/f.txt", []byte("x"), 0644)

	if !fileutil.DirExists(fs, "/d") {
		t.Error("expected DirExists true for /d")
	}
	if fileutil.DirExists(fs, "/f.txt") {
		t.Error("expected DirExists false for file /f.txt")
	}
	if fileutil.DirExists(fs, "/missing") {
		t.Error("expected DirExists false for missing path")
	}
}

func TestCopyFile_ContentMatches(t *testing.T) {
	fs := afero.NewMemMapFs()
	content := []byte("hello, world")
	_ = afero.WriteFile(fs, "/src.txt", content, 0644)

	if err := fileutil.CopyFile(fs, "/src.txt", "/sub/dst.txt"); err != nil {
		t.Fatalf("CopyFile: %v", err)
	}
	got, _ := afero.ReadFile(fs, "/sub/dst.txt")
	if string(got) != string(content) {
		t.Errorf("content mismatch: got %q, want %q", got, content)
	}
}

func TestReadWriteLines_RoundTrip(t *testing.T) {
	fs := afero.NewMemMapFs()
	lines := []string{"line one", "line two", "line three"}

	if err := fileutil.WriteLines(fs, "/lines.txt", lines); err != nil {
		t.Fatalf("WriteLines: %v", err)
	}
	got, err := fileutil.ReadLines(fs, "/lines.txt")
	if err != nil {
		t.Fatalf("ReadLines: %v", err)
	}
	if len(got) != len(lines) {
		t.Fatalf("got %d lines, want %d", len(got), len(lines))
	}
	for i, l := range lines {
		if got[i] != l {
			t.Errorf("line %d: got %q, want %q", i, got[i], l)
		}
	}
}
