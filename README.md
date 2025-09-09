## go-libsd3

Library utilitas Go untuk aplikasi SD3.

### Struktur Folder

- `cmd/` : Entry point aplikasi utama.
- `helper/` : Helper, misal pengelolaan kunci/setting, logging, dsb.
- `pkg/` : Package utilitas seperti multidbutil (multiton DB pool), dbutil, dan logutil.
- `test/` : File-file unit test.

### Fitur Utama

- Membaca konfigurasi koneksi database (PostgreSQL & SQL Server) dari file XML (`helper/kunci`).
- Utilitas multiton database pool untuk koneksi per data center (`pkg/multidbutil`).
- Logging dengan rotasi file harian dan format log yang mudah dibaca (`helper/logging`).
- Helper untuk pengelolaan setting/kunci aplikasi (`helper/settinglibgooo`).
- Producer/consumer helper untuk kebutuhan event atau message queue (jika ada di project Anda).

### Contoh Penggunaan Multiton DB

```go
import "github.com/boni-fm/go-libsd3/pkg/multidbutil"

func main() {
    multiDB := &multidbutil.MultiDB{}
    multiDB.SetupMultiDB("DC01")
    db, err := multiDB.GetDB("DC01")
    if err != nil {
        panic(err)
    }
    // Gunakan db untuk query
}
```

### Contoh Penggunaan Logger

```go
import "github.com/boni-fm/go-libsd3/helper/logging"

func main() {
    log := logging.NewLogger()
    log.Say("Contoh log info")
    log.SayError("Contoh log error")
    log.SayWithField("Log dengan field", "user", "admin")
}
```

### Contoh Penggunaan SettingLibGooo

```go
import "github.com/boni-fm/go-libsd3/helper/settinglibgooo"

func main() {
    // Inisialisasi setting/kunci
    setting := settinglibgooo.NewSettingLib("DC01")
    
    // Membaca konfigurasi dari SettingWeb.xml
    connInfo := setting.SettingWebClient.GetConnectionInfoPostgre()
    fmt.Println("IP Postgres:", connInfo.IPPostgres)
}
```

### Contoh Penggunaan Producer (jika ada)

```go
import "github.com/boni-fm/go-libsd3/helper/producer"

func main() {
    prod := producer.NewProducer("my-topic")
    err := prod.SendMessage("Hello, SD3!")
    if err != nil {
        fmt.Println("Gagal kirim pesan:", err)
    }
}
```

Log akan otomatis tersimpan di folder log harian dengan format rapi.

### Dependensi Eksternal

- github.com/sirupsen/logrus
- github.com/snowzach/rotatefilehook
- github.com/mattn/go-colorable
- github.com/lib/pq (driver PostgreSQL)

### Cara Menjalankan Test

Jalankan perintah berikut di root folder:

```bash
go test ./test/...
```

### Catatan

- Pastikan file konfigurasi `SettingWeb.xml` tersedia di folder yang sesuai dengan OS dan struktur yang diharapkan.
- Untuk logging, pastikan folder tujuan log (`/var/log/nginx/api` di Linux atau path yang sesuai di Windows) dapat ditulis oleh aplikasi.
- Untuk penggunaan multiton DB, pastikan environment dan koneksi database sudah dikonfigurasi dengan benar.
- Untuk producer/consumer, pastikan service/broker yang digunakan sudah berjalan dan dapat diakses.

---
**Ini README.md di generate oleh AI**
