package test

import (
	"os"
	"testing"

	"github.com/boni-fm/go-libsd3/helper/kunci"
)

func TestGetConnStringDockerPostgre_WithEnv(t *testing.T) {
	os.Setenv("KUNCI_ENV_KUNCI", "dummy_kunci_value")
	defer os.Unsetenv("KUNCI_ENV_KUNCI")

	connStr := kunci.GetConnStringDockerPostgre("POSTGRE")
	if connStr == "" {
		t.Error("Expected non-empty connection string when env is set")
	}
}

func TestGetConnStringDockerPostgre_WithoutEnv(t *testing.T) {
	os.Unsetenv("KUNCI_ENV_KUNCI")
	// This will fallback to yamlreader, which may need to be mocked for a real test
	_ = kunci.GetConnStringDockerPostgre("POSTGRE")
	// No assertion here, as it depends on the yamlreader's behavior
}
