package example

// ======================================================
// Contoh pemakaian package log (pkg/log)
// ======================================================
//
// Logger ini adalah wrapper di atas logrus dengan rotasi file via lumberjack.
// Timezone diambil dari env var TZ; kalau kosong/invalid pakai Asia/Jakarta.
//
// Fitur yang di-demo:
//   - NewLogger — logger default tanpa nama app
//   - NewLoggerWithFilename — logger dengan nama app (nama file log ikut nama app)
//   - NewLoggerWithPath — logger dengan direktori log kustom (berguna di testing)
//   - Say / Sayf — log level Info
//   - SayWithField / SayWithFields — log Info dengan field tambahan
//   - SayError / SayErrorf — log level Error
//   - SayFatal / SayFatalf — log level Fatal (proses keluar!)
//   - Infof / Warnf / Errorf / Panicf — akses langsung ke logrus method
//   - SetLevel — ubah log level on the fly
//   - Penggunaan TZ env untuk mengatur timezone log

import (
	"fmt"
	"os"

	"github.com/boni-fm/go-libsd3/pkg/log"
	"github.com/sirupsen/logrus"
)

// ContohNewLogger mendemonstrasikan penggunaan NewLogger() tanpa nama app.
// File log akan dinamai "logs_YYYY-MM-DD.log" di direktori default OS.
func ContohNewLogger() {
	// Set timezone dulu via env var (opsional, default Asia/Jakarta)
	os.Setenv("TZ", "Asia/Jakarta")

	l := log.NewLogger()

	// Say — info sederhana, ga ada format
	l.Say("service dimulai")

	// Sayf — info dengan format printf
	l.Sayf("menerima request dari %s", "192.168.1.100")

	// SayError — log error, proses ga mati
	l.SayError("koneksi database timeout, coba lagi nanti")

	// SayErrorf — error dengan format
	l.SayErrorf("user id=%d gagal login, percobaan ke-%d", 42, 3)

	// Infof / Warnf / Errorf — akses method logrus langsung
	l.Infof("server jalan di port %d", 8080)
	l.Warnf("memori hampir penuh: %d%%", 87)
	l.Errorf("file tidak ditemukan: %s", "/etc/config.yaml")

	// WithField / WithFields — logrus fields, dipake buat structured logging
	l.WithField("request_id", "abc-123").Info("request masuk")
	l.WithFields(logrus.Fields{
		"user_id":   42,
		"endpoint":  "/api/users",
		"method":    "GET",
		"status":    200,
	}).Info("request selesai")

	// SetLevel — ubah log level (berguna buat debug di dev, matiin di prod)
	l.SetLevel(logrus.DebugLevel)
	l.Debugf("ini cuma keliatan kalo level Debug: %s", "detail internal")
	l.SetLevel(logrus.InfoLevel) // kembaliin ke Info
}

// ContohNewLoggerWithFilename mendemonstrasikan NewLoggerWithFilename(appName).
// File log dinamai "logs_<appName>_YYYY-MM-DD.log".
func ContohNewLoggerWithFilename() {
	// Biasanya dipanggil di main(), satu kali aja
	l := log.NewLoggerWithFilename("my-api-service")

	l.Say("service my-api-service jalan")
	l.SayWithField("versi mulai", "version", "1.2.3")
	l.SayWithFields("config loaded", map[string]interface{}{
		"env":  "production",
		"port": 8080,
	})

	// SayFatal — log FATAL lalu os.Exit(1), proses BENERAN mati!
	// Uncomment baris di bawah hanya kalau emang mau matiin proses:
	// l.SayFatal("config file tidak ditemukan, abort!")
	// l.SayFatalf("port %d sudah dipakai proses lain", 8080)

	fmt.Println("logger berhasil dibuat, cek file log di direktori default")
}

// ContohNewLoggerWithPath mendemonstrasikan NewLoggerWithPath(appName, dirPath).
// Direktori log bisa dikonfigurasi bebas — cocok buat testing atau deployment custom.
func ContohNewLoggerWithPath(logDir string) {
	// logDir bisa "/var/log/myapp", "D:/logs", atau t.TempDir() di test
	l := log.NewLoggerWithPath("worker-service", logDir)

	l.Say("worker service dimulai")
	l.Sayf("membaca tugas dari antrian, dir log: %s", logDir)

	// Contoh penggunaan dalam goroutine (thread-safe karena logrus pake mutex)
	go func() {
		l.Say("goroutine worker berjalan")
		l.SayError("goroutine ketemu error tapi ga panic")
	}()
}

// ContohTZEnv mendemonstrasikan pengaruh TZ environment variable.
// Logger akan pakai timezone yang diset di TZ buat timestamp di log.
func ContohTZEnv() {
	// Coba set timezone Jakarta
	os.Setenv("TZ", "Asia/Jakarta")
	l1 := log.NewLogger()
	l1.Say("log ini pake Asia/Jakarta timezone")

	// Coba set timezone UTC
	os.Setenv("TZ", "UTC")
	l2 := log.NewLogger()
	l2.Say("log ini pake UTC timezone")

	// Timezone ga valid — otomatis fallback ke Asia/Jakarta
	os.Setenv("TZ", "Invalid/Timezone")
	l3 := log.NewLogger() // bakal cetak warning ke stderr, lalu fallback
	l3.Say("log ini pake fallback Asia/Jakarta karena TZ ga valid")

	// Bersiin env
	os.Unsetenv("TZ")
}
