package envloader

import (
	"os"
	"testing"
)

func writeTestEnvFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "envloader_test_*.env")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	f.WriteString(content)
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestLoad_BasicKeyValue(t *testing.T) {
	filename := writeTestEnvFile(t, "TEST_KEY_ENVLOADER=hello\n")
	os.Unsetenv("TEST_KEY_ENVLOADER")
	defer os.Unsetenv("TEST_KEY_ENVLOADER")

	if err := Load(filename); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if v := os.Getenv("TEST_KEY_ENVLOADER"); v != "hello" {
		t.Errorf("expected 'hello', got '%s'", v)
	}
}

func TestLoad_NoOverwrite(t *testing.T) {
	filename := writeTestEnvFile(t, "TEST_NO_OVERWRITE=newvalue\n")
	os.Setenv("TEST_NO_OVERWRITE", "original")
	defer os.Unsetenv("TEST_NO_OVERWRITE")

	if err := Load(filename); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if v := os.Getenv("TEST_NO_OVERWRITE"); v != "original" {
		t.Errorf("expected 'original' (no overwrite), got '%s'", v)
	}
}

func TestLoadOverride_Overwrites(t *testing.T) {
	filename := writeTestEnvFile(t, "TEST_OVERRIDE_VAR=newvalue\n")
	os.Setenv("TEST_OVERRIDE_VAR", "original")
	defer os.Unsetenv("TEST_OVERRIDE_VAR")

	if err := LoadOverride(filename); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if v := os.Getenv("TEST_OVERRIDE_VAR"); v != "newvalue" {
		t.Errorf("expected 'newvalue', got '%s'", v)
	}
}

func TestLoad_CommentsAndEmpty(t *testing.T) {
	content := "# this is a comment\n\nTEST_AFTER_COMMENT=value\n"
	filename := writeTestEnvFile(t, content)
	os.Unsetenv("TEST_AFTER_COMMENT")
	defer os.Unsetenv("TEST_AFTER_COMMENT")

	if err := Load(filename); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if v := os.Getenv("TEST_AFTER_COMMENT"); v != "value" {
		t.Errorf("expected 'value', got '%s'", v)
	}
}

func TestLoad_QuotedValues(t *testing.T) {
	content := "TEST_QUOTED_DOUBLE=\"double quoted\"\nTEST_QUOTED_SINGLE='single quoted'\n"
	filename := writeTestEnvFile(t, content)
	os.Unsetenv("TEST_QUOTED_DOUBLE")
	os.Unsetenv("TEST_QUOTED_SINGLE")
	defer os.Unsetenv("TEST_QUOTED_DOUBLE")
	defer os.Unsetenv("TEST_QUOTED_SINGLE")

	if err := Load(filename); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if v := os.Getenv("TEST_QUOTED_DOUBLE"); v != "double quoted" {
		t.Errorf("expected 'double quoted', got '%s'", v)
	}
	if v := os.Getenv("TEST_QUOTED_SINGLE"); v != "single quoted" {
		t.Errorf("expected 'single quoted', got '%s'", v)
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	err := Load("/nonexistent/path/.env")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestStripQuotes(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"hello"`, "hello"},
		{`'world'`, "world"},
		{"noQuotes", "noQuotes"},
		{`"`, `"`},
		{"", ""},
	}
	for _, tc := range tests {
		result := stripQuotes(tc.input)
		if result != tc.expected {
			t.Errorf("stripQuotes(%q) = %q, want %q", tc.input, result, tc.expected)
		}
	}
}
