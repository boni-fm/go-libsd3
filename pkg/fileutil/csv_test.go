package fileutil_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/boni-fm/go-libsd3/pkg/fileutil"
	"github.com/spf13/afero"
)

type testRow struct {
	Name  string `csv:"name"`
	Score int    `csv:"score"`
}

func TestWriteAndReadCSVStructs(t *testing.T) {
	fs := afero.NewMemMapFs()
	data := []testRow{
		{Name: "Alice", Score: 90},
		{Name: "Bob", Score: 85},
	}

	if err := fileutil.WriteStructsToCSVFs(fs, data, "/out/data.csv"); err != nil {
		t.Fatalf("WriteStructsToCSVFs: %v", err)
	}

	got, err := fileutil.ReadCSVToStructsFs[testRow](fs, "/out/data.csv")
	if err != nil {
		t.Fatalf("ReadCSVToStructsFs: %v", err)
	}
	if !reflect.DeepEqual(got, data) {
		t.Errorf("got %v, want %v", got, data)
	}
}

func TestWriteSliceToCSV(t *testing.T) {
	fs := afero.NewMemMapFs()
	rows := [][]string{
		{"id", "value"},
		{"1", "hello"},
		{"2", "world"},
	}

	if err := fileutil.WriteSliceToCSVFs(fs, rows, "/out/slice.csv"); err != nil {
		t.Fatalf("WriteSliceToCSVFs: %v", err)
	}

	content, err := afero.ReadFile(fs, "/out/slice.csv")
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if !strings.Contains(string(content), "id,value") {
		t.Errorf("expected header row in CSV, got: %s", string(content))
	}
}

func TestWriteStructsToCSV_CreatesParentDir(t *testing.T) {
	fs := afero.NewMemMapFs()
	data := []testRow{{Name: "Test", Score: 1}}

	err := fileutil.WriteStructsToCSVFs(fs, data, "/deep/nested/dir/out.csv")
	if err != nil {
		t.Fatalf("WriteStructsToCSVFs should auto-create dirs: %v", err)
	}
	if !fileutil.FileExists(fs, "/deep/nested/dir/out.csv") {
		t.Error("expected file to exist after write")
	}
}
