package test

import (
	"testing"

	"github.com/boni-fm/go-libsd3/pkg/dbutil"
)

func TestConnectAndClose(t *testing.T) {
	db, err := dbutil.SetupConnectionDatabase()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	if db == nil {
		t.Fatal("Connect returned nil db")
	}
	err = db.Close()
	if err != nil {
		t.Errorf("Failed to close db: %v", err)
	}
}

func TestHealthCheck(t *testing.T) {
	db, err := dbutil.SetupConnectionDatabase()
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	status := db.HealthCheck()
	if status == "" {
		t.Error("HealthCheck returned empty string")
	}
	t.Logf("HealthCheck status: %s", status)
	db.Close()
}
