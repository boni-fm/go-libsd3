package log

import (
	"os"
	"strings"
	"testing"
)

func TestResolveTimezone_NewYork(t *testing.T) {
	os.Setenv("TZ", "America/New_York")
	defer os.Unsetenv("TZ")

	loc := resolveTimezone()
	if !strings.Contains(loc.String(), "New_York") {
		t.Errorf("expected location containing New_York, got %s", loc.String())
	}
}

func TestResolveTimezone_Fallback(t *testing.T) {
	os.Unsetenv("TZ")

	loc := resolveTimezone()
	if loc.String() != "Asia/Jakarta" {
		t.Errorf("expected Asia/Jakarta fallback, got %s", loc.String())
	}
}

func TestResolveTimezone_InvalidFallback(t *testing.T) {
	os.Setenv("TZ", "Invalid/Timezone")
	defer os.Unsetenv("TZ")

	loc := resolveTimezone()
	if loc.String() != "Asia/Jakarta" {
		t.Errorf("expected Asia/Jakarta fallback for invalid TZ, got %s", loc.String())
	}
}
