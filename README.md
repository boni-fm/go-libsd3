## github.com/boni-fm/go-libsd3

Library utilitas Go untuk aplikasi SD3.

### Struktur Folder

- `cmd/` : Entry point aplikasi utama.
- `helper/` : Helper, misal pengelolaan kunci/setting, logging, dsb.
- `pkg/` : Package utilitas seperti dbutil dan logutil.
- `test/` : File-file unit test.

### Fitur Utama

- Membaca konfigurasi koneksi database dari file XML (`helper/kunci`).
- Utilitas koneksi database PostgreSQL (`pkg/dbutil`).
- Logging dengan rotasi file harian dan format log yang mudah dibaca (`helper/logging`).

### Contoh Penggunaan Logger


```go
import "github.com/boni-fm/github.com/boni-fm/go-libsd3/helper/logging"

func main() {
	log := logging.NewLogger()
	log.Say("Contoh log info")
	log.SayError("Contoh log error")
	log.SayWithField("Log dengan field", "user", "admin")
}
```

Log akan otomatis tersimpan di folder logs harian dengan format rapi.

### Dependensi Eksternal

- github.com/sirupsen/logrus
- github.com/snowzach/rotatefilehook
- github.com/mattn/go-colorable

### Cara Menjalankan Test

Jalankan perintah berikut di root folder:

```bash
go test ./test/...
```

### Catatan

- Pastikan file konfigurasi `SettingWeb.xml` tersedia di folder home user sesuai struktur yang diharapkan.
- Untuk logging, pastikan folder tujuan log dapat ditulis oleh aplikasi.
