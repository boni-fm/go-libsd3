package postgres

import (
	"bytes"
	"context"
	"testing"
)

func TestExportQueryToCSV_ClosedDB(t *testing.T) {
	db := &Database{
		isClosed: true,
	}

	var buf bytes.Buffer
	err := db.ExportQueryToCSV(context.Background(), &buf, "SELECT 1")
	if err == nil {
		t.Error("expected error when database is closed")
	}

	if buf.Len() != 0 {
		t.Error("expected empty buffer when database is closed")
	}
}
