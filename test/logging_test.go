package test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/boni-fm/go-libsd3/helper/logging"
)

func TestLoggerDailyRotation(t *testing.T) {
	log := logging.NewLogger()
	log.Say("Test log message")
	log.Sayf("Test logf %d", 123)
	log.SayWithField("Test with field", "foo", "bar")
	log.SayWithFields("Test with fields", map[string]interface{}{"a": 1, "b": "c"})
	log.Warn("Test warn level")
	log.Errorf("Test errorf %s", "err")

	// Check log file exists for today
	homedir, _ := os.UserHomeDir()
	filename := "logs" + time.Now().Format("2006-01-02") + ".log"
	logPath := filepath.Join(homedir, "_docker", "_app", "logs", filename)

	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		t.Errorf("Expected log file %s to exist", logPath)
	}
}
