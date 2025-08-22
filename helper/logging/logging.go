package logging

import (
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
)

/*
	TODO :
	- perbaiki struktur log nya supaya lebih bagus wkwkwk
*/

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

	logLevel := logrus.InfoLevel
	log := logrus.New()
	log.SetLevel(logLevel)

	homedir, _ := os.UserHomeDir()
	filename := "logs" + time.Now().Format("2006-01-02") + ".log"
	filepath := filepath.Join(homedir, "_docker", "_app", "logs", filename)

	rotateFileHook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
		Filename:   filepath,
		MaxSize:    50, // megabytes
		MaxBackups: 3,  // amouts
		MaxAge:     28, //days
		Level:      logLevel,
		Formatter: &logrus.TextFormatter{
			PadLevelText:           true,
			ForceColors:            false,
			FullTimestamp:          true,
			TimestampFormat:        "2006-01-02 15:04:05",
			DisableLevelTruncation: true,
			QuoteEmptyFields:       true,
			DisableQuote:           true,
			DisableTimestamp:       false,
		},
	})
	if err != nil {
		log.Fatalf("Failed to initialize file rotate hook: %v", err)
	}

	log.AddHook(rotateFileHook)

	return &logger{log}
}
