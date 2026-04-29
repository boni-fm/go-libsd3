# Changelog

Semua perubahan penting pada proyek ini akan didokumentasikan di file ini.

Format mengikuti [Keep a Changelog](https://keepachangelog.com/id/1.0.0/).

## [Unreleased]

### Ditambahkan

- **AUDIT.md** — Dokumen audit keamanan & kualitas kode yang mendetail (AUDIT-001 s/d AUDIT-015)
- **pkg/retry** — Paket baru untuk mekanisme retry dengan exponential back-off; mendukung generic `DoWithResult[T]`
- **pkg/envloader** — Paket baru untuk memuat environment variables dari file `.env`; mendukung `Load` dan `LoadOverride`
- **pkg/httputil** — Paket baru untuk middleware HTTP; termasuk `RequestIDMiddleware` dan `GetRequestID`
- **pkg/db/postgres/registry.go** — Registry multiton dengan composite key `(kunci, kodedc)` untuk mendukung multi-tenant
- **pkg/db/postgres/ExportQueryToCSV** — Metode baru untuk mengekspor hasil query ke format CSV
- **pkg/log/log_test.go** — Unit test untuk `resolveTimezone` dan fallback timezone
- **pkg/db/postgres/pool_test.go** — Unit test untuk concurrency dan double-checked locking
- **pkg/db/postgres/registry_test.go** — Unit test untuk registry multi-tenant
- **pkg/db/postgres/csv_test.go** — Unit test untuk `ExportQueryToCSV` pada DB yang tertutup
- **pkg/retry/retry_test.go** — Unit test lengkap untuk mekanisme retry
- **pkg/envloader/envloader_test.go** — Unit test untuk pemuat `.env`
- **pkg/httputil/middleware_test.go** — Unit test untuk middleware request-id
- File contoh (`example_test.go`) untuk paket `retry`, `envloader`, `httputil`, `log`, dan `postgres`

### Diperbaiki

- **AUDIT-001** (`pkg/settinglibgo/settinglib.go`) — Ganti type assertion tidak aman `key.(string)` dengan ok-check
- **AUDIT-002** (`pkg/settinglibgo/xml.go`) — Ganti `log.Fatalf` dengan `fmt.Fprintf(os.Stderr, ...)` + `return ""`; pindahkan update cache ke sebelum loop pencarian
- **AUDIT-003** (`pkg/db/postgres/pool.go`) — Implementasikan double-checked locking di `GetConnection` untuk mencegah TOCTOU race
- **AUDIT-004** (`pkg/db/postgres/db.go`) — Pindahkan pengecekan `isClosed` ke sebelum `Pool.Acquire` di `CopyFrom`
- **AUDIT-005** (`pkg/log/formatter.go`) — Tambahkan field `Location *time.Location` ke `CustomLogFormatter`; gunakan untuk format timestamp
- **AUDIT-007** (`pkg/kafkautil/producer.go`) — Hapus duplikasi `kgo.RequiredAcks(kgo.LeaderAck())` yang ditimpa oleh `AllISRAcks`
- **AUDIT-008** (`pkg/settinglibgo/client.go`) — Perbaiki data race pada variabel paket `BASEURL`; gunakan variabel lokal
- **AUDIT-009** (`pkg/settinglibgo/xml.go`) — Tangani error dari `io.ReadAll` yang sebelumnya diabaikan
- **AUDIT-010** (`pkg/settinglibgo/client.go`) — Tangani error dari `json.Marshal` yang sebelumnya diabaikan

### Diperbarui

- **pkg/log/log.go** — Tambahkan `resolveTimezone()` untuk membaca env var `TZ`; fallback ke `Asia/Jakarta`
- **pkg/log/formatter.go** — Perbarui `CustomLogFormatter` untuk mendukung `Location *time.Location`
- **pkg/db/postgres/errors.go** — Tambahkan komentar dokumentasi Bahasa Indonesia
- **pkg/settinglibgo/error.go** — Tambahkan komentar dokumentasi Bahasa Indonesia
- **pkg/config/constant/constant.go** — Tambahkan komentar dokumentasi Bahasa Indonesia
- **README.md** — Ditulis ulang sepenuhnya dalam Bahasa Indonesia dengan contoh penggunaan lengkap
- **CHANGELOG.md** — Dibuat dalam Bahasa Indonesia
