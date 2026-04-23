package test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/boni-fm/go-libsd3/pkg/log"
)

func TestNewLogger_NotNil(t *testing.T) {
	l := log.NewLogger()
	if l == nil {
		t.Fatal("NewLogger() returned nil")
	}
}

func TestNewLoggerWithFilename_NotNil(t *testing.T) {
	l := log.NewLoggerWithFilename("testapp")
	if l == nil {
		t.Fatal("NewLoggerWithFilename() returned nil")
	}
}

func TestNewLoggerWithPath_WritesToCustomDir(t *testing.T) {
	dir := t.TempDir()
	l := log.NewLoggerWithPath("testapp", dir)
	if l == nil {
		t.Fatal("NewLoggerWithPath() returned nil")
	}
	l.Say("hello from test")
	// give lumberjack time to flush
	time.Sleep(50 * time.Millisecond)
	// check a log file was created
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir: %v", err)
	}
	if len(entries) == 0 {
		t.Fatal("expected log file to be created in temp dir")
	}
	// read and check content
	logFile := filepath.Join(dir, entries[0].Name())
	data, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if !strings.Contains(string(data), "[testapp]") {
		t.Errorf("expected log to contain [testapp], got: %s", string(data))
	}
}

func TestLogger_Say_DoesNotPanic(t *testing.T) {
	dir := t.TempDir()
	l := log.NewLoggerWithPath("testapp", dir)
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Say panicked: %v", r)
		}
	}()
	l.Say("test message")
}

func TestLogger_SayError_DoesNotPanic(t *testing.T) {
	dir := t.TempDir()
	l := log.NewLoggerWithPath("testapp", dir)
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("SayError panicked: %v", r)
		}
	}()
	l.SayError("test error")
}

func TestLogger_SayWithField_DoesNotPanic(t *testing.T) {
	dir := t.TempDir()
	l := log.NewLoggerWithPath("testapp", dir)
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("SayWithField panicked: %v", r)
		}
	}()
	l.SayWithField("test", "key", "value")
}

func TestLogger_SayWithFields_DoesNotPanic(t *testing.T) {
	dir := t.TempDir()
	l := log.NewLoggerWithPath("testapp", dir)
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("SayWithFields panicked: %v", r)
		}
	}()
	l.SayWithFields("test", map[string]interface{}{"k": "v"})
}
