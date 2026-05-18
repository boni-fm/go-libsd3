# Audit Keamanan & Kualitas Kode — go-libsd3

## KRITIS

### AUDIT-001
**File:** `pkg/settinglibgo/settinglib.go:GetConnStringPostgre()`
**Masalah:** Type assertion `key.(string)` tidak aman — akan panic jika value adalah nil atau bukan string.
**Dampak:** Crash proses saat startup jika nilai key tidak sesuai tipe.
**Perbaikan:** Gunakan type assertion dengan ok-check: `if v, ok := key.(string); ok { ... }`.

### AUDIT-002
**File:** `pkg/settinglibgo/xml.go:DynamicSettingWebXMLReader()`
**Masalah:** (1) Memanggil `log.Fatalf()` yang akan menghentikan seluruh proses. (2) Cache diperbarui SETELAH loop pencarian sehingga tidak pernah terisi saat cache hit.
**Dampak:** Process exit yang tidak terkontrol; cache tidak pernah berfungsi dengan benar.
**Perbaikan:** Ganti `log.Fatalf` dengan `fmt.Fprintf(os.Stderr, ...)` + `return ""`; pindahkan update cache ke SEBELUM loop pencarian.

### AUDIT-003
**File:** `pkg/db/postgres/pool.go:GetConnection()`
**Masalah:** TOCTOU (Time-of-Check-Time-of-Use) race — map dicek dengan RLock, koneksi baru dibuat tanpa lock, lalu disimpan dengan Lock baru. Dua goroutine dapat membuat koneksi duplikat.
**Dampak:** Kebocoran koneksi database; koneksi duplikat untuk kode DC yang sama.
**Perbaikan:** Gunakan pola double-checked locking.

### AUDIT-004
**File:** `pkg/db/postgres/db.go:CopyFrom()`
**Masalah:** Pengecekan `isClosed` terjadi SETELAH koneksi diperoleh dari pool (baris 521-529), sehingga guard closed tidak efektif.
**Dampak:** Operasi copy mungkin dijalankan pada koneksi yang sudah ditutup.
**Perbaikan:** Pindahkan pengecekan `isClosed` ke SEBELUM `d.Pool.Acquire(ctx)`.

## TINGGI

### AUDIT-005
**File:** `pkg/log/formatter.go:CustomLogFormatter.Format()`
**Masalah:** `entry.Time` tidak menggunakan timezone yang dikonfigurasi; logger selalu menggunakan timezone lokal server tanpa memperhatikan env var `TZ`.
**Dampak:** Timestamp log salah di lingkungan produksi dengan zona waktu berbeda (mis. Asia/Jakarta).
**Perbaikan:** Tambahkan field `Location *time.Location` ke `CustomLogFormatter`; gunakan lokasi tersebut untuk memformat `entry.Time`.

### AUDIT-006
**File:** `pkg/db/postgres/db.go:Query()`
**Masalah:** Mengembalikan `nil` saat database ditutup; pemanggil yang memanggil `.Scan()` pada nilai nil akan panic.
**Dampak:** Panic di goroutine pemanggil saat database ditutup.
**Catatan:** Tanda tangan fungsi tidak dapat diubah tanpa breaking change; tambahkan dokumentasi agar pemanggil menggunakan `IsClosed()` terlebih dahulu.

### AUDIT-007
**File:** `pkg/kafkautil/producer.go:NewProducer()`
**Masalah:** `kgo.RequiredAcks` diset dua kali (baris 63: `LeaderAck`, baris 66: `AllISRAcks`); pemanggilan kedua secara diam-diam menimpa yang pertama, menghasilkan pengaturan durabilitas yang lebih ketat dari yang diharapkan.
**Dampak:** Performa Kafka tidak sesuai konfigurasi yang dimaksud.
**Perbaikan:** Hapus baris 63 `kgo.RequiredAcks(kgo.LeaderAck())`.

### AUDIT-008
**File:** `pkg/settinglibgo/client.go:GetVariable()`
**Masalah:** Variabel level paket `BASEURL` dimutasi tanpa lock di dalam `GetVariable()`; ini adalah data race jika beberapa goroutine memanggil `GetVariable()` secara bersamaan.
**Dampak:** Data race yang dapat menyebabkan nilai BASEURL tidak konsisten.
**Perbaikan:** Baca env var secara lokal tanpa menulis ke variabel paket.

## MENENGAH

### AUDIT-009
**File:** `pkg/settinglibgo/xml.go`
**Masalah:** Error dari `io.ReadAll` diabaikan secara diam-diam (baris 62: `byteValue, _ := io.ReadAll(xmlFile)`).
**Dampak:** Kegagalan baca file tidak dilaporkan; data XML mungkin tidak lengkap.
**Perbaikan:** Tangani error dan kembalikan `""` dengan pesan error ke stderr.

### AUDIT-010
**File:** `pkg/settinglibgo/client.go`
**Masalah:** Error dari `json.Marshal` diabaikan secara diam-diam (baris 84: `bodyByte, _ := json.Marshal(params)`).
**Dampak:** Request dengan body rusak dikirim tanpa peringatan.
**Perbaikan:** Tangani error marshal dan kembalikan error.

### AUDIT-011
**File:** `pkg/versi/versi.go`, `pkg/settinglibgo/kunci.go`
**Masalah:** Variabel `var log` level paket diinisialisasi saat waktu import, menyebabkan pembuatan file log segera saat package diimport.
**Dampak:** File log dibuat bahkan jika fitur logging tidak digunakan; efek samping yang tidak diinginkan saat import.
**Perbaikan:** Gunakan inisialisasi lazy atau dependency injection untuk logger.

### AUDIT-012
**File:** `pkg/db/postgres/pool.go:ConnectionPool`
**Masalah:** Hanya menggunakan `kodeDc` sebagai composite key; dimensi `kunci` tidak ada, membuat beberapa tenant dengan kunci/credential berbeda tetapi kodeDc sama menjadi tidak mungkin.
**Dampak:** Tidak mendukung multi-tenant dengan kodeDC yang sama.
**Perbaikan:** Implementasikan registry dengan composite key `(kunci, kodedc)` — lihat registry.go.

## RENDAH

### AUDIT-013
**File:** `pkg/db/postgres/db.go`
**Masalah:** Banyak pemanggilan `fmt.Printf` untuk logging alih-alih menggunakan structured logger.
**Dampak:** Log tidak konsisten; tidak dapat difilter atau diformat secara terpusat.
**Perbaikan:** Gunakan structured logger (logrus) di seluruh file.

### AUDIT-014
**File:** `pkg/log/log.go`
**Masalah:** `generateLogFilename` menggunakan `time.Now()` yang mengambil waktu lokal server; tidak ada kesadaran timezone.
**Dampak:** Nama file log mungkin tidak sesuai dengan timezone yang diharapkan di lingkungan produksi.
**Perbaikan:** Gunakan lokasi timezone yang sudah di-resolve saat membuat nama file log.

### AUDIT-015
**File:** Semua file exported symbols
**Masalah:** Semua simbol yang diekspor tidak memiliki komentar dokumentasi dalam Bahasa Indonesia.
**Dampak:** Developer tidak mendapat panduan penggunaan API yang jelas.
**Perbaikan:** Tambahkan komentar dokumentasi Bahasa Indonesia ke semua simbol yang diekspor.
