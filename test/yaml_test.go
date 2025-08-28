package test

import (
	"os"
	"testing"

	yaml "github.com/boni-fm/go-libsd3/helper/yamlreader"
)

func TestReadConfigDynamic(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "testconfig-*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())

	yamlContent := []byte("kunci: testvalue\nfoo: bar\nnum: 42\n")
	if _, err := tmpfile.Write(yamlContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpfile.Close()

	cfg, err := yaml.ReadConfigDynamic(tmpfile.Name())
	if err != nil {
		t.Fatalf("ReadConfigDynamic failed: %v", err)
	}

	if cfg["kunci"] != "testvalue" {
		t.Errorf("Expected kunci to be 'testvalue', got %v", cfg["kunci"])
	}
	if cfg["foo"] != "bar" {
		t.Errorf("Expected foo to be 'bar', got %v", cfg["foo"])
	}
	if cfg["num"].(int) != 42 {
		t.Errorf("Expected num to be 42, got %v", cfg["num"])
	}
}
