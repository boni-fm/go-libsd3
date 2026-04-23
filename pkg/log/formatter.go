package log

import (
	"bytes"
	"fmt"
	"time"

	"github.com/boni-fm/go-libsd3/pkg/config/constant"
	"github.com/sirupsen/logrus"
)

// CustomLogFormatter adalah formatter log kustom yang mendukung konfigurasi timezone.
// Digunakan oleh rotatefilehook untuk memformat setiap baris log.
type CustomLogFormatter struct {
	// AppName adalah nama aplikasi yang akan ditampilkan di setiap baris log.
	AppName *string
	// Location adalah zona waktu yang digunakan untuk memformat timestamp log.
	// Jika nil, menggunakan time.UTC.
	Location *time.Location
}

// Format memformat entri log menjadi byte slice sesuai format standar service.
// Format: [AppName] [datetime] [level] - message - Data: fields
func (c *CustomLogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	loc := c.Location
	if loc == nil {
		loc = time.UTC
	}

	ts := entry.Time.In(loc).Format(constant.DATETIME_FORMAT)
	fmt.Fprintf(b, "[%s] [%s] [%s] - %s - Data: %v\n", *c.AppName, ts, entry.Level, entry.Message, entry.Data)

	return b.Bytes(), nil
}
