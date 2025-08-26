package logging

import (
	"bytes"
	"fmt"

	"github.com/sirupsen/logrus"
)

type CustomLogFormatter struct {
	AppName *string
}

func (c *CustomLogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	// ini buat setting format lognya, sekarang ngikutin format log service abstraction + data
	fmt.Fprintf(b, "[%s] [%s] [%s] - %s - Data: %v\n", *c.AppName, entry.Time.Format("2006-01-02 15:04:05"), entry.Level, entry.Message, entry.Data)

	return b.Bytes(), nil
}
