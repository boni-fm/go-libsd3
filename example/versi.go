package example

// ======================================================
// Contoh pemakaian package versi (pkg/versi)
// ======================================================
//
// Package versi nyediain fungsi untuk ngecek dan memvalidasi versi program
// terhadap master program yang tersimpan di database Postgres.
//
// Fitur yang di-demo:
//   - GetVersiProgramPostgre — cek versi program, insert monitoring, return status
//   - PostgreConstrBuilder — konversi berbagai format connection string ke format pgx/Postgres
//
// Cara kerja GetVersiProgramPostgre:
//  1. Konek ke Postgres menggunakan connection string yang diberikan
//  2. Query versi program dari tabel dc_program_vbdtl_t
//  3. Bandingkan versi yang diinput dengan versi di database
//  4. Kalau sama → insert ke tabel monitoring, return "OKE..."
//  5. Kalau beda → return pesan update/mismatch
//
// Note:
//   - Butuh koneksi Postgres beneran buat jalan.
//   - Tabel yang dipakai: dc_program_vbdtl_t dan dc_monitoring_program_t

import (
	"fmt"

	"github.com/boni-fm/go-libsd3/pkg/versi"
)

// ContohGetVersiProgramPostgre mendemonstrasikan pengecekan versi program.
//
// Parameter:
//   - Constr: connection string Postgres (bisa berbagai format, akan dikonversi otomatis)
//   - Kodedc: kode DC (misalnya "G009SIM")
//   - NamaProgram: nama program yang dicek (case-insensitive, diubah ke UPPER)
//   - Versi: versi program yang sedang berjalan (misal "1.0.0" → dikonversi ke int)
//   - IPKomputer: IP address client untuk monitoring
//
// Return value:
//   - "OKE..." → versi cocok, sudah dicatat di monitoring
//   - String pesan update → versi tidak cocok, perlu update
//   - String pesan error → koneksi gagal atau program tidak terdaftar
func ContohGetVersiProgramPostgre() {
	// Format connection string yang langsung didukung oleh pgx
	connStr := "host=db.internal port=5432 user=appuser password=secret dbname=mydb sslmode=disable"

	// Cek versi program "KASIR" versi "1.2.3" di DC "G009SIM"
	result := versi.GetVersiProgramPostgre(
		connStr,    // connection string
		"G009SIM",  // kode DC
		"kasir",    // nama program (akan di-upper otomatis)
		"1.2.3",    // versi program saat ini (titik akan dihapus → 123)
		"192.168.1.100", // IP komputer client
	)

	switch result {
	case "OKE...":
		fmt.Println("versi program cocok, aktivitas berhasil dicatat di monitoring")
	default:
		fmt.Printf("hasil cek versi: %s\n", result)
		// Kemungkinan output:
		// "    Program .:KASIR:. belum update,\r\n      Versi update \r\n--==::>> 1.2.4 <<::==--"
		// "    Program .:KASIR:. Versi program tidak sama dengan master ..."
		// "    Program .:KASIR:. belum terdaftar di Master Program DC ..."
		// "Koneksi DB Gagal..."
	}
}

// ContohGetVersiDenganConstrBuilder mendemonstrasikan penggunaan
// PostgreConstrBuilder SEBELUM memanggil GetVersiProgramPostgre.
//
// Kalau connection string kamu dalam format C#/ADO.NET (Server=...;Database=...),
// konversi dulu dengan PostgreConstrBuilder sebelum dipakai.
func ContohGetVersiDenganConstrBuilder() {
	// Format connection string gaya C# / ADO.NET (dari settinglib lama)
	constrGayaLama := "Server=db.internal;Port=5432;Database=mydb;User Id=appuser;Password=secret123"

	// PostgreConstrBuilder konversi ke format standar Postgres/pgx
	constrTerkonversi := versi.PostgreConstrBuilder(constrGayaLama)
	fmt.Printf("sebelum: %s\n", constrGayaLama)
	fmt.Printf("sesudah: %s\n", constrTerkonversi)
	// Output: host=db.internal port=5432 user=appuser password=secret123 dbname=mydb sslmode=disable

	// Sekarang bisa dipakai di GetVersiProgramPostgre
	result := versi.GetVersiProgramPostgre(
		constrTerkonversi,
		"G010SIM",
		"ABSENSI",
		"2.0.1",
		"10.0.0.55",
	)
	fmt.Printf("hasil versi ABSENSI: %s\n", result)
}

// ContohPostgreConstrBuilder mendemonstrasikan konversi berbagai format
// connection string ke format standar yang bisa dipakai pgx/Postgres.
func ContohPostgreConstrBuilder() {
	testCases := []struct {
		nama  string
		input string
	}{
		{
			nama:  "format Server=",
			input: "Server=192.168.1.10;Port=5432;Database=proddb;Username=pguser;Password=pgsecret",
		},
		{
			nama:  "format Host=",
			input: "Host=pg.server.internal;Port=5432;Database=appdb;User Id=admin;Password=admin123",
		},
		{
			nama:  "format campuran spasi tidak konsisten",
			input: "Server = 10.0.0.1 ; Port = 5432 ; Database = mydb ; Username = usr ; Password = pwd",
		},
		{
			nama:  "hanya beberapa field (sisanya kosong)",
			input: "Server=localhost;Database=testdb",
		},
	}

	for _, tc := range testCases {
		hasil := versi.PostgreConstrBuilder(tc.input)
		fmt.Printf("[%s]\n  input : %s\n  output: %s\n\n", tc.nama, tc.input, hasil)
	}
}

// ContohVersiTidakCocok mendemonstrasikan berbagai skenario hasil pengecekan versi.
func ContohVersiTidakCocok() {
	connStr := "host=db.internal port=5432 user=pguser password=pgpass dbname=proddb sslmode=disable"

	skenario := []struct {
		nama      string
		program   string
		versi     string
		deskripsi string
	}{
		{
			nama:      "versi sama — sukses",
			program:   "KASIR",
			versi:     "1.2.3",
			deskripsi: "versi di DB juga 1.2.3 → OKE",
		},
		{
			nama:      "versi lebih lama dari master",
			program:   "KASIR",
			versi:     "1.1.0",
			deskripsi: "versi di DB adalah 1.2.3, client masih 1.1.0 → perlu update",
		},
		{
			nama:      "versi lebih baru dari master",
			program:   "KASIR",
			versi:     "1.9.9",
			deskripsi: "versi client lebih baru dari master → mismatch",
		},
		{
			nama:      "program tidak terdaftar",
			program:   "PROGRAM_BARU_BELUM_DAFTAR",
			versi:     "1.0.0",
			deskripsi: "program belum ada di dc_program_vbdtl_t → pesan tidak terdaftar",
		},
	}

	for _, s := range skenario {
		fmt.Printf("── skenario: %s ──\n", s.nama)
		fmt.Printf("   deskripsi: %s\n", s.deskripsi)

		result := versi.GetVersiProgramPostgre(
			connStr,
			"G009SIM",
			s.program,
			s.versi,
			"192.168.1.100",
		)
		fmt.Printf("   hasil: %q\n\n", result)
	}
}
