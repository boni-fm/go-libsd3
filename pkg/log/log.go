package log

import (
	"path/filepath"
	"runtime"
	"time"

	"github.com/boni-fm/go-libsd3/pkg/config/constant"
	"github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
)

/*
	Library logger
	Depedency:
	- logrus
	- rotatefilehook

	TODO:
	- buat jadi 3 hari file log nya
*/

type Logger struct {
	*logrus.Logger
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

// Generate log filename
func generateLogFilename(AppName string) string {
	date := time.Now().Format(constant.DATE_FORMAT)
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

func generateRotateFileHook(filepath, appName string) *rotatefilehook.RotateFileHook {
	hook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
		Filename:   filepath,
		MaxSize:    50, // megabytes
		MaxBackups: 7,  // amouts
		MaxAge:     28, //days
		Level:      logrus.InfoLevel,
		Formatter:  &CustomLogFormatter{AppName: &appName},
	})

	if err != nil {
		panic(err)
	}

	return hook.(*rotatefilehook.RotateFileHook)
}

// #####################################################3

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
