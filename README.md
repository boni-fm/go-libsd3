## go-libsd3

Library utilitas Go untuk aplikasi SD3.

### Struktur Folder

- `cmd/` : Berisi entry point aplikasi utama.
- `helper/` : Berisi helper, misal pengelolaan kunci/setting.
- `pkg/` : Berisi package utilitas seperti dbutil dan logutil.
- `test/` : Berisi file-file unit test.

### Fitur Utama

- Membaca konfigurasi koneksi database dari file XML.
- Utilitas koneksi database PostgreSQL.

### Cara Menjalankan Test

Jalankan perintah berikut di root folder:

```bash
go test ./test/...
```

### Catatan

- Pastikan file konfigurasi `SettingWeb.xml` tersedia di folder home user sesuai struktur yang diharapkan.
