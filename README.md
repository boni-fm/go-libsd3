# go-libsd3

Library Go bersama untuk layanan backend berbasis arsitektur DC (Data Center). Menyediakan utilitas siap pakai untuk koneksi database, logging, konfigurasi, Kafka, retry, HTTP middleware, dan lainnya.

## Instalasi

```bash
go get github.com/boni-fm/go-libsd3
```

## Daftar Paket

| Paket | Deskripsi |
|-------|-----------|
| `pkg/db/postgres` | Koneksi PostgreSQL dengan pool, registry multi-tenant, dan bulk export CSV |
| `pkg/log` | Logger berbasis logrus dengan rotasi file dan dukungan timezone |
| `pkg/retry` | Mekanisme retry dengan exponential back-off |
| `pkg/envloader` | Pemuat environment variable dari file `.env` |
| `pkg/httputil` | Middleware HTTP (request-id tracing) |
| `pkg/kafkautil` | Producer Kafka berbasis franz-go |
| `pkg/settinglibgo` | Klien layanan kunci (SettingLib) |
| `pkg/auth` | Helper autentikasi AWS/JWT |
| `pkg/config/constant` | Konstanta konfigurasi global |
| `pkg/yaml` | Helper pembaca konfigurasi YAML |
| `pkg/versi` | Informasi versi program |

## Contoh Penggunaan

### Database PostgreSQL

```go
import "github.com/boni-fm/go-libsd3/pkg/db/postgres"

// Daftarkan database dengan composite key (kunci, kodedc)
db, err := postgres.RegisterDB(ctx, "kunci-tenant-a", "DC001", postgres.Config{
    KodeDC:  "DC001",
    AppName: "myapp",
})

// Ambil kembali instance yang sudah terdaftar
db, err := postgres.GetDB("kunci-tenant-a", "DC001")

// Jalankan query
var result MyStruct
err = db.SelectOne(ctx, &result, "SELECT * FROM tabel WHERE id = $1", 1)

// Export ke CSV
var buf bytes.Buffer
err = db.ExportQueryToCSV(ctx, &buf, "SELECT * FROM tabel")
```

### Logger

```go
import logger "github.com/boni-fm/go-libsd3/pkg/log"

// Buat logger baru (timezone dari env TZ, fallback Asia/Jakarta)
log := logger.NewLoggerWithFilename("myapp")
log.Say("Aplikasi dimulai")
log.SayWithField("koneksi berhasil", "host", "localhost")
log.SayError("terjadi kesalahan")
```

### Retry

```go
import "github.com/boni-fm/go-libsd3/pkg/retry"

err := retry.Do(ctx, retry.DefaultConfig(), func() error {
    return callExternalService()
})

result, err := retry.DoWithResult(ctx, retry.DefaultConfig(), func() (string, error) {
    return fetchData()
})
```

### Environment Loader

```go
import "github.com/boni-fm/go-libsd3/pkg/envloader"

// Muat .env tanpa menimpa variabel yang sudah ada
err := envloader.Load(".env")

// Muat .env dan timpa variabel yang sudah ada
err := envloader.LoadOverride(".env.local")
```

### HTTP Middleware

```go
import "github.com/boni-fm/go-libsd3/pkg/httputil"

// Tambahkan request-id ke setiap request
mux := http.NewServeMux()
handler := httputil.RequestIDMiddleware(mux)
http.ListenAndServe(":8080", handler)

// Ambil request-id dari context
id := httputil.GetRequestID(r.Context())
```

### Kafka Producer

```go
import "github.com/boni-fm/go-libsd3/pkg/kafkautil"

config := &kafkautil.ProducerConfig{
    Brokers: []string{"localhost:9092"},
    Topic:   "my-topic",
}
producer, err := kafkautil.NewProducer(config, logger)
err = producer.SendJSON(ctx, "key", myData)
```

## Konfigurasi Environment Variable

| Variabel | Deskripsi | Default |
|----------|-----------|---------|
| `TZ` | Timezone untuk logger (format IANA, mis. `Asia/Jakarta`) | `Asia/Jakarta` |
| `KunciWeb` | Kunci identifikasi layanan kunci (Docker) | - |
| `KUNCI_IP_DOMAIN` | Alamat IP/domain layanan kunci | `localhost` |

## Panduan Kontribusi

1. Fork repository ini
2. Buat branch fitur: `git checkout -b feat/nama-fitur`
3. Pastikan semua test lulus: `go test ./...`
4. Pastikan tidak ada peringatan: `go vet ./...`
5. Commit dengan pesan deskriptif dalam Bahasa Indonesia
6. Buat Pull Request ke branch `main`

### Konvensi

- Semua komentar dokumentasi pada simbol yang diekspor **harus dalam Bahasa Indonesia**
- Gunakan `fmt.Fprintf(os.Stderr, ...)` untuk error logging — hindari `log.Fatalf`
- Semua akses map konkuren harus dilindungi dengan `sync.RWMutex`
- Gunakan pola double-checked locking untuk inisialisasi malas yang thread-safe

## Lisensi

Hak cipta © 2025 boni-fm. Seluruh hak dilindungi.
