package test

import (
	"context"
	"testing"

	"github.com/boni-fm/go-libsd3/pkg/db/postgres"
)

func TestInitConstrByKodeDc_EnvOverride(t *testing.T) {
	connStr, err := postgres.InitConstrByKodeDc(context.Background(), "hohoho")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	t.Logf("constr %+v", connStr)
	if connStr == "" {
		t.Fatalf("expected non-empty connection string when env set")
	}
}
