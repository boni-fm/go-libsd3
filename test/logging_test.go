package test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/boni-fm/go-libsd3/helper/logging"
)

func TestLoggerDailyRotation(t *testing.T) {
	logname := "Testing"
	log := logging.NewLoggerWithFilename(logname)
	log.Say("Test log message")
	log.Sayf("Test logf %d", 123)
	log.SayWithField("Test with field", "foo", "bar")
	log.SayWithFields("Test with fields", map[string]interface{}{"a": 1, "b": "c"})
	log.Warn("Test warn level")
	log.Errorf("Test errorf %s", "err")

	// Check log file exists for today
	homedir, _ := os.UserHomeDir()
	filename := "logs" + logname + time.Now().Format("2006-01-02") + ".log"
	logDir := filepath.Join(homedir, "_docker", "_app", "logs")
	logPath := filepath.Join(logDir, filename)

	// Ensure log directory exists
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		if err := os.MkdirAll(logDir, 0755); err != nil {
			t.Fatalf("Failed to create log directory: %v", err)
		}
	}

	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		t.Errorf("Expected log file %s to exist", logPath)
	}
}

func TestLoggerDefault(t *testing.T) {
	log := logging.NewLogger()
	log.Say("Default logger test message")
	log.Warn("Default logger warning")
	log.Errorf("Default logger error")
	log.SayWithField("Default logger with field", "key", "value")
	log.SayWithFields("Default logger with fields", map[string]interface{}{"foo": 42, "bar": "baz"})
	// No file check here, just ensure no panic and output is produced
}
