package logging

import (
	"path/filepath"
	"runtime"
	"time"

	"github.com/boni-fm/go-libsd3/config"
	"github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
)

type logger struct{ *logrus.Logger }

func (l *logger) Say(msg string) {
	l.Info(msg)
}
func (l *logger) Sayf(fmt string, args ...interface{}) {
	l.Infof(fmt, args...)
}
func (l *logger) SayWithField(msg string, k string, v interface{}) {
	l.WithField(k, v).Info(msg)
}
func (l *logger) SayWithFields(msg string, fields map[string]interface{}) {
	l.WithFields(fields).Info(msg)
}

func (l *logger) SayFatal(msg string) {
	l.Fatal(msg)
}
func (l *logger) SayFatalf(fmt string, args ...interface{}) {
	l.Fatalf(fmt, args...)
}
func (l *logger) SayError(msg string) {
	l.Error(msg)
}
func (l *logger) SayErrorf(fmt string, args ...interface{}) {
	l.Errorf(fmt, args...)
}

func NewLogger() *logger {
	log := logrus.New()
	log.SetLevel(logrus.InfoLevel)

	filename := generateLogFilename("")
	filepath := filepath.Join(getLogFilePath(), filename)

	rotateFileHook := generateRotateFileHook(filepath, "")
	log.AddHook(rotateFileHook)

	return &logger{log}
}

func NewLoggerWithFilename(AppName string) *logger {
	log := logrus.New()
	log.SetLevel(logrus.InfoLevel)

	filename := generateLogFilename(AppName)
	filepath := filepath.Join(getLogFilePath(), filename)

	rotateFileHook := generateRotateFileHook(filepath, AppName)
	log.AddHook(rotateFileHook)

	return &logger{log}
}

// fungsi setup loggernya
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
