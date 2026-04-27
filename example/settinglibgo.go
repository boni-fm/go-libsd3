package example

// ======================================================
// Contoh pemakaian package settinglibgo (pkg/settinglibgo)
// ======================================================
//
// Package settinglibgo nyediain client HTTP buat ngambil konfigurasi
// dari service kunci (internal API). Fungsi utamanya adalah:
//   - Dapetin connection string Postgres
//   - Dapetin nilai variabel konfigurasi (IP, port, password, dsb)
//
// Fitur yang di-demo:
//   - NewSettingLib — buat MainConfiguration dari kode DC
//   - NewSettingLibWithAppName — sama tapi dengan nama aplikasi
//   - NewSettingLibClient — buat client HTTP langsung
//   - GetConnectionString — dapetin connection string Postgres/DB
//   - GetVariable — dapetin satu nilai variabel dari kunci service
//   - SetPGConStringFromWebservice — ambil semua komponen Postgres dari kunci
//   - SetAppName — update nama aplikasi setelah dibuat
//   - GetConnStringPostgre — helper shortcut (baca dari env/yaml otomatis)
//
// Note:
//   - Butuh service kunci yang beneran buat jalan.
//   - Secara default, kunci service ada di localhost, bisa dioverride
//     via env var KUNCI_IP_DOMAIN.
//   - Kunci service dipanggil via HTTP POST ke /<kunci>/GetVariabel.

import (
	"fmt"
	"os"

	"github.com/boni-fm/go-libsd3/pkg/settinglibgo"
)

// ContohNewSettingLib mendemonstrasikan pembuatan MainConfiguration
// dan pengambilan connection string Postgres.
func ContohNewSettingLib() {
	// ─── NewSettingLib ───
	// kodedc biasanya berformat "kunciG009SIM" atau "kunciG010SIM"
	// kalau kodedc ga pake prefix "kunci", otomatis ditambahkan
	lib := settinglibgo.NewSettingLib("kunciG009SIM")

	// GetConnectionString — dapetin connection string Postgres
	// Ini akan memanggil GetVariable ke service kunci buat dapetin
	// IPPostgres, PortPostgres, UserPostgres, PasswordPostgres, DatabasePostgres
	conStr := lib.GetConnectionString("POSTGRE")
	if conStr == "" {
		fmt.Println("connection string kosong — pastikan service kunci jalan")
		return
	}
	fmt.Printf("connection string: %s\n", conStr)
	// Output contoh:
	// host=db.internal port=5432 user=appuser password=secret dbname=mydb sslmode=disable application_name=GOAPPS
}

// ContohNewSettingLibWithAppName mendemonstrasikan pembuatan dengan nama aplikasi.
// Nama aplikasi akan tampil di connection string sebagai "application_name".
// Berguna buat monitoring query dari sisi Postgres (pg_stat_activity).
func ContohNewSettingLibWithAppName() {
	// Dengan nama aplikasi — muncul di pg_stat_activity.application_name
	lib := settinglibgo.NewSettingLibWithAppName("kunciG009SIM", "order-service")

	conStr := lib.GetConnectionString("POSTGRE")
	if conStr != "" {
		fmt.Printf("conn string (dengan app name): %s\n", conStr)
		// ...application_name=order-service
	}

	// SetAppName — update nama setelah dibuat
	lib.SetAppName("payment-service")
	fmt.Println("app name diupdate ke payment-service")
}

// ContohNewSettingLibClient mendemonstrasikan penggunaan SettingLibClient langsung
// untuk mengambil variabel konfigurasi satu per satu.
func ContohNewSettingLibClient() {
	// NewSettingLibClient menerima kode kunci (dengan/tanpa prefix "kunci")
	// Kalau ga ada prefix "kunci", otomatis ditambahkan
	client := settinglibgo.NewSettingLibClient("G009SIM")
	// ekuivalen dengan NewSettingLibClient("kunciG009SIM")

	// ─── GetVariable ─── ambil satu variabel dari kunci service
	// Kunci yang umum: IPPostgres, PortPostgres, UserPostgres, PasswordPostgres,
	//                  DatabasePostgres, BaseUrlCloud, IPKafka, dst

	ipPostgres, err := client.GetVariable("IPPostgres")
	if err != nil {
		fmt.Printf("GetVariable IPPostgres gagal: %v\n", err)
	} else {
		fmt.Printf("IP Postgres: %s\n", ipPostgres)
	}

	portPostgres, err := client.GetVariable("PortPostgres")
	if err != nil {
		fmt.Printf("GetVariable PortPostgres gagal: %v\n", err)
	} else {
		fmt.Printf("Port Postgres: %s\n", portPostgres)
	}

	// Ambil base URL cloud buat service lain
	baseURLCloud, err := client.GetVariable("BaseUrlCloud")
	if err != nil {
		fmt.Printf("GetVariable BaseUrlCloud gagal: %v\n", err)
	} else {
		fmt.Printf("Base URL Cloud: %s\n", baseURLCloud)
	}
}

// ContohKUNCIIPDomainEnv mendemonstrasikan override alamat kunci service via env var.
//
// Default: service kunci ada di localhost.
// Untuk Docker/produksi, set KUNCI_IP_DOMAIN ke hostname container kunci.
func ContohKUNCIIPDomainEnv() {
	// Set env var buat demo
	os.Setenv("KUNCI_IP_DOMAIN", "docker-hub-nginx-1")
	defer os.Unsetenv("KUNCI_IP_DOMAIN")

	client := settinglibgo.NewSettingLibClient("kunciG009SIM")

	// GetVariable sekarang akan hit http://docker-hub-nginx-1/kunciG009SIM/GetVariabel
	val, err := client.GetVariable("IPPostgres")
	if err != nil {
		fmt.Printf("gagal hit kunci service di docker: %v\n", err)
	} else {
		fmt.Printf("IP Postgres dari docker kunci: %s\n", val)
	}

	// Reset ke default (localhost)
	os.Unsetenv("KUNCI_IP_DOMAIN")
	fmt.Println("KUNCI_IP_DOMAIN diunset, kembali ke localhost")
}

// ContohSetPGConStringFromWebservice mendemonstrasikan pengisian semua komponen
// Postgres sekaligus dari service kunci.
func ContohSetPGConStringFromWebservice() {
	lib := settinglibgo.NewSettingLib("kunciG009SIM")

	// SetPGConStringFromWebservice memanggil kunci service untuk mengisi
	// semua field Postgres: IP, Port, User, Password, Database
	updatedLib, err := lib.SetPGConStringFromWebservice()
	if err != nil {
		fmt.Printf("SetPGConStringFromWebservice gagal: %v\n", err)
		fmt.Println("pastikan service kunci jalan dan dapat diakses")
		return
	}

	// Sekarang GetConnectionString bisa langsung dipanggil tanpa hit network lagi
	conStr := updatedLib.GetConnectionString("POSTGRE")
	fmt.Printf("connection string setelah SetPG: %s\n", conStr)
}

// ContohGetConnStringPostgre mendemonstrasikan fungsi shortcut GetConnStringPostgre.
//
// Ini adalah fungsi level package yang otomatis:
//  1. Baca kunci dari env var KunciWeb (untuk deployment Docker)
//  2. Kalau ga ada, baca dari config.yaml (untuk non-Docker)
//  3. Panggil service kunci dengan kunci yang didapat
//  4. Return connection string Postgres
func ContohGetConnStringPostgre() {
	// Cara 1: Via environment variable (Docker deployment)
	os.Setenv("KunciWeb", "kunciG009SIM")
	defer os.Unsetenv("KunciWeb")

	conStr := settinglibgo.GetConnStringPostgre()
	if conStr == "" {
		fmt.Println("conn string kosong, service kunci mungkin ga jalan")
		return
	}
	fmt.Printf("conn string via env var: %s\n", conStr)

	// Cara 2: Via config.yaml (non-Docker / bare metal)
	// Pastikan ada file config.yaml di lokasi yang dicari dengan isi:
	// kunci: kunciG009SIM
	os.Unsetenv("KunciWeb")
	conStr2 := settinglibgo.GetConnStringPostgre()
	if conStr2 == "" {
		fmt.Println("conn string kosong, cek config.yaml punya key 'kunci'")
	} else {
		fmt.Printf("conn string via config.yaml: %s\n", conStr2)
	}
}
