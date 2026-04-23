package log

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/boni-fm/go-libsd3/pkg/config/constant"
	"github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
)

// Logger adalah wrapper di atas logrus.Logger dengan kemampuan rotasi file log.
// Mendukung konfigurasi timezone melalui environment variable TZ.
type Logger struct {
	*logrus.Logger
}

// NewLogger membuat instance Logger baru dengan nama file log default.
// Timezone diambil dari environment variable TZ; jika kosong atau tidak valid,
// menggunakan Asia/Jakarta sebagai fallback.
func NewLogger() *Logger {
	loc := resolveTimezone()
	log := logrus.New()
	log.SetLevel(logrus.InfoLevel)

	filename := generateLogFilename("", loc)
	fp := filepath.Join(getLogFilePath(), filename)

	rotateFileHook := generateRotateFileHook(fp, "", loc)
	log.AddHook(rotateFileHook)

	return &Logger{log}
}

// NewLoggerWithFilename membuat instance Logger baru dengan nama aplikasi kustom.
// Nama aplikasi akan disertakan dalam nama file log dan setiap baris log.
// Parameter:
//   - AppName: nama aplikasi yang digunakan sebagai bagian dari nama file log
func NewLoggerWithFilename(AppName string) *Logger {
	loc := resolveTimezone()
	log := logrus.New()
	log.SetLevel(logrus.InfoLevel)

	filename := generateLogFilename(AppName, loc)
	fp := filepath.Join(getLogFilePath(), filename)

	rotateFileHook := generateRotateFileHook(fp, AppName, loc)
	log.AddHook(rotateFileHook)

	return &Logger{log}
}

// resolveTimezone membaca env var TZ dan mengembalikan *time.Location yang sesuai.
// Jika TZ kosong atau tidak valid, menggunakan Asia/Jakarta sebagai fallback.
// Jika TZ di-set tetapi tidak valid, peringatan dicetak ke stderr.
func resolveTimezone() *time.Location {
	tz := os.Getenv("TZ")
	if tz != "" {
		loc, err := time.LoadLocation(tz)
		if err == nil {
			return loc
		}
		fmt.Fprintf(os.Stderr, "[log] TZ=%q tidak valid, menggunakan Asia/Jakarta: %v\n", tz, err)
	}
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return time.UTC
	}
	return loc
}

// generateLogFilename membuat nama file log berdasarkan AppName dan tanggal saat ini.
// Parameter:
//   - AppName: nama aplikasi (opsional); jika kosong, tidak disertakan dalam nama file
//   - loc: zona waktu untuk penentuan tanggal
func generateLogFilename(AppName string, loc *time.Location) string {
	if loc == nil {
		loc = time.UTC
	}
	date := time.Now().In(loc).Format(constant.DATE_FORMAT)
	if AppName != "" {
		return "logs_" + AppName + "_" + date + ".log"
	}
	return "logs_" + date + ".log"
}

func getLogFilePath() string {
	if runtime.GOOS == "windows" {
		return constant.FILEPATH_LOG_WINDOWS
	}
	return constant.FILEPATH_LOG_LINUX
}

func generateRotateFileHook(fp, appName string, loc *time.Location) *rotatefilehook.RotateFileHook {
	hook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
		Filename:   fp,
		MaxSize:    50,
		MaxBackups: 7,
		MaxAge:     28,
		Level:      logrus.InfoLevel,
		Formatter:  &CustomLogFormatter{AppName: &appName, Location: loc},
	})

	if err != nil {
		panic(err)
	}

	return hook.(*rotatefilehook.RotateFileHook)
}

// Say mencatat pesan di level Info.
func (l *Logger) Say(msg string) {
	l.Info(msg)
}

// Sayf mencatat pesan terformat di level Info.
func (l *Logger) Sayf(fmt string, args ...interface{}) {
	l.Infof(fmt, args...)
}

// SayWithField mencatat pesan dengan satu field tambahan di level Info.
func (l *Logger) SayWithField(msg string, k string, v interface{}) {
	l.WithField(k, v).Info(msg)
}

// SayWithFields mencatat pesan dengan beberapa field tambahan di level Info.
func (l *Logger) SayWithFields(msg string, fields map[string]interface{}) {
	l.WithFields(fields).Info(msg)
}

// SayFatal mencatat pesan di level Fatal dan menghentikan proses.
func (l *Logger) SayFatal(msg string) {
	l.Fatal(msg)
}

// SayFatalf mencatat pesan terformat di level Fatal dan menghentikan proses.
func (l *Logger) SayFatalf(fmt string, args ...interface{}) {
	l.Fatalf(fmt, args...)
}

// SayError mencatat pesan di level Error.
func (l *Logger) SayError(msg string) {
	l.Error(msg)
}

// SayErrorf mencatat pesan terformat di level Error.
func (l *Logger) SayErrorf(fmt string, args ...interface{}) {
	l.Errorf(fmt, args...)
}
