package postgres_test

import (
	"fmt"

	"github.com/boni-fm/go-libsd3/pkg/db/postgres"
)

func ExampleGetDB() {
	// GetDB digunakan untuk mengambil instance database yang sudah terdaftar.
	// Contoh:
	//   db, err := postgres.GetDB("kunci-tenant-a", "DC001")
	//   if err != nil {
	//       log.Fatal(err)
	//   }
	_, err := postgres.GetDB("kunci-contoh", "DC-CONTOH")
	fmt.Println("error ketika DB tidak ditemukan:", err != nil)
	// Output:
	// error ketika DB tidak ditemukan: true
}

func ExampleCloseAll() {
	// CloseAll menutup semua koneksi yang terdaftar di registry.
	// Aman dipanggil meski tidak ada koneksi yang terdaftar.
	err := postgres.CloseAll()
	fmt.Println("CloseAll error:", err)
	// Output:
	// CloseAll error: <nil>
}
