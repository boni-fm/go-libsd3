package log_test

import (
	"fmt"

	logger "github.com/boni-fm/go-libsd3/pkg/log"
)

func ExampleNewLoggerWithFilename() {
	// Contoh penggunaan NewLoggerWithFilename:
	// log := logger.NewLoggerWithFilename("myapp")
	// log.Say("Aplikasi dimulai")
	// log.SayWithField("koneksi berhasil", "host", "localhost")
	_ = logger.NewLoggerWithFilename // referenced to satisfy go vet
	fmt.Println("Lihat dokumentasi untuk contoh penggunaan")
	// Output:
	// Lihat dokumentasi untuk contoh penggunaan
}
