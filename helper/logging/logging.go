package logging

import (
	"path/filepath"
	"runtime"
	"time"

	"github.com/boni-fm/go-libsd3/config"
	"github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
)

type Logger struct{ *logrus.Logger }

func (l *Logger) Say(msg string) {
	l.Info(msg)
}
func (l *Logger) Sayf(fmt string, args ...interface{}) {
	l.Infof(fmt, args...)
}
func (l *Logger) SayWithField(msg string, k string, v interface{}) {
	l.WithField(k, v).Info(msg)
}
func (l *Logger) SayWithFields(msg string, fields map[string]interface{}) {
	l.WithFields(fields).Info(msg)
}

func (l *Logger) SayFatal(msg string) {
	l.Fatal(msg)
}
func (l *Logger) SayFatalf(fmt string, args ...interface{}) {
	l.Fatalf(fmt, args...)
}
func (l *Logger) SayError(msg string) {
	l.Error(msg)
}
func (l *Logger) SayErrorf(fmt string, args ...interface{}) {
	l.Errorf(fmt, args...)
}

func NewLogger() *Logger {
	log := logrus.New()
	log.SetLevel(logrus.InfoLevel)

	filename := generateLogFilename("")
	filepath := filepath.Join(getLogFilePath(), filename)

	rotateFileHook := generateRotateFileHook(filepath, "")
	log.AddHook(rotateFileHook)

	return &Logger{log}
}

func NewLoggerWithFilename(AppName string) *Logger {
	log := logrus.New()
	log.SetLevel(logrus.InfoLevel)

	filename := generateLogFilename(AppName)
	filepath := filepath.Join(getLogFilePath(), filename)

	rotateFileHook := generateRotateFileHook(filepath, AppName)
	log.AddHook(rotateFileHook)

	return &Logger{log}
}

// fungsi setup Loggernya
func getLogFilePath() string {
	if runtime.GOOS == "windows" {
		return config.FILEPATH_LOG_WINDOWS
	}
	return config.FILEPATH_LOG_LINUX
}

func generateLogFilename(appName string) string {
	dateStr := time.Now().Format(config.DATE_FORMAT)
	if appName != "" {
		return "logs_" + appName + "_" + dateStr + ".log"
	}
	return "logs_" + dateStr + ".log"
}

func generateRotateFileHook(filepath, appName string) *rotatefilehook.RotateFileHook {
	hook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
		Filename:   filepath,
		MaxSize:    50, // megabytes
		MaxBackups: 3,  // amouts
		MaxAge:     28, //days
		Level:      logrus.InfoLevel,
		Formatter:  &CustomLogFormatter{AppName: &appName},
	})

	if err != nil {
		panic(err)
	}

	return hook.(*rotatefilehook.RotateFileHook)
}
