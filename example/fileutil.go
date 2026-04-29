package example

// ======================================================
// Contoh pemakaian package fileutil (pkg/fileutil)
// ======================================================
//
// Package fileutil nyediain helper untuk:
//   - Operasi file umum: EnsureDir, FileExists, DirExists, CopyFile, ReadLines, WriteLines
//   - CSV: WriteStructsToCSV, ReadCSVToStructs, WriteSliceToCSV (dengan versi afero)
//   - ZIP: ZipFiles, UnzipTo (dengan perlindungan zip-slip)
//   - FTP: NewFTPClient, Upload, Download, Delete, List, MakeDir
//
// Semua fungsi file/csv/zip punya versi afero (Fs) buat kemudahan testing.
// Gunakan afero.NewMemMapFs() di test dan afero.NewOsFs() di produksi.

import (
	"fmt"
	"os"
	"time"

	"github.com/boni-fm/go-libsd3/pkg/fileutil"
	"github.com/spf13/afero"
)

// ──────────────────────────────────────────────────────────
// Contoh 1: Operasi File Umum
// ──────────────────────────────────────────────────────────

// ContohFileUtil mendemonstrasikan EnsureDir, FileExists, DirExists, CopyFile,
// ReadLines, dan WriteLines menggunakan in-memory filesystem (afero).
func ContohFileUtil() {
	// Di produksi pakai afero.NewOsFs()
	// Di sini pakai MemMapFs biar ga perlu nulis ke disk beneran
	fs := afero.NewMemMapFs()

	// EnsureDir — buat direktori (termasuk induknya), aman kalau udah ada
	err := fileutil.EnsureDir(fs, "/data/logs/app")
	if err != nil {
		fmt.Printf("EnsureDir gagal: %v\n", err)
		return
	}
	fmt.Println("direktori /data/logs/app berhasil dibuat")

	// Buat juga direktori bersarang
	fileutil.EnsureDir(fs, "/data/backup/2024/01")
	fileutil.EnsureDir(fs, "/data/backup/2024/02")

	// DirExists — cek apakah direktori ada
	fmt.Println("DirExists /data/logs/app:", fileutil.DirExists(fs, "/data/logs/app"))   // true
	fmt.Println("DirExists /data/tidakada:", fileutil.DirExists(fs, "/data/tidakada"))   // false

	// FileExists — cek apakah file (bukan dir) ada
	fmt.Println("FileExists /data/logs/app:", fileutil.FileExists(fs, "/data/logs/app")) // false, itu dir bukan file

	// Buat file dulu via WriteLines
	lines := []string{
		"baris pertama",
		"baris kedua",
		"baris ketiga dengan unicode: ☕",
	}
	err = fileutil.WriteLines(fs, "/data/logs/app/app.log", lines)
	if err != nil {
		fmt.Printf("WriteLines gagal: %v\n", err)
		return
	}
	fmt.Println("berhasil nulis", len(lines), "baris ke /data/logs/app/app.log")

	// FileExists — sekarang file ada
	fmt.Println("FileExists /data/logs/app/app.log:", fileutil.FileExists(fs, "/data/logs/app/app.log")) // true

	// ReadLines — baca file teks, return slice string
	readBack, err := fileutil.ReadLines(fs, "/data/logs/app/app.log")
	if err != nil {
		fmt.Printf("ReadLines gagal: %v\n", err)
		return
	}
	for i, l := range readBack {
		fmt.Printf("baris[%d]: %s\n", i, l)
	}

	// CopyFile — salin file, direktori tujuan dibuat otomatis
	err = fileutil.CopyFile(fs, "/data/logs/app/app.log", "/data/backup/2024/01/app.log.bak")
	if err != nil {
		fmt.Printf("CopyFile gagal: %v\n", err)
		return
	}
	fmt.Println("file berhasil dicopy ke backup")

	// Verifikasi hasil copy
	fmt.Println("FileExists backup:", fileutil.FileExists(fs, "/data/backup/2024/01/app.log.bak"))

	// Contoh pakai OS filesystem (beneran ke disk)
	// fs := afero.NewOsFs()
	// fileutil.EnsureDir(fs, "/tmp/myapp/data")
}

// ──────────────────────────────────────────────────────────
// Contoh 2: CSV — WriteStructsToCSV, ReadCSVToStructs, WriteSliceToCSV
// ──────────────────────────────────────────────────────────

// Karyawan adalah struct yang akan disimpan ke CSV.
// Tag `csv` dipake oleh library gocsv buat header kolom.
type Karyawan struct {
	ID        int    `csv:"id"`
	Nama      string `csv:"nama"`
	Jabatan   string `csv:"jabatan"`
	Gaji      int    `csv:"gaji"`
}

// ContohFileUtilCSV mendemonstrasikan semua fungsi CSV.
func ContohFileUtilCSV() {
	fs := afero.NewMemMapFs()

	// ─── WriteStructsToCSVFs ───
	karyawans := []Karyawan{
		{ID: 1, Nama: "Budi Santoso", Jabatan: "Engineer", Gaji: 10000000},
		{ID: 2, Nama: "Siti Rahayu", Jabatan: "Designer", Gaji: 9000000},
		{ID: 3, Nama: "Ahmad Fauzi", Jabatan: "Manager", Gaji: 15000000},
	}

	csvPath := "/data/csv/karyawan.csv"
	err := fileutil.WriteStructsToCSVFs[Karyawan](fs, karyawans, csvPath)
	if err != nil {
		fmt.Printf("WriteStructsToCSVFs gagal: %v\n", err)
		return
	}
	fmt.Println("berhasil tulis", len(karyawans), "karyawan ke CSV")

	// ─── ReadCSVToStructsFs ───
	bacaKaryawan, err := fileutil.ReadCSVToStructsFs[Karyawan](fs, csvPath)
	if err != nil {
		fmt.Printf("ReadCSVToStructsFs gagal: %v\n", err)
		return
	}
	for _, k := range bacaKaryawan {
		fmt.Printf("karyawan: ID=%d Nama=%s Jabatan=%s Gaji=%d\n",
			k.ID, k.Nama, k.Jabatan, k.Gaji)
	}

	// ─── WriteSliceToCSVFs ─── (raw 2D slice, baris pertama jadi header)
	rawData := [][]string{
		{"kode", "nama_produk", "harga", "stok"},   // header
		{"P001", "Laptop Gaming", "15000000", "10"},
		{"P002", "Mouse Wireless", "250000", "50"},
		{"P003", "Keyboard Mechanical", "500000", "30"},
	}

	rawCsvPath := "/data/csv/produk.csv"
	err = fileutil.WriteSliceToCSVFs(fs, rawData, rawCsvPath)
	if err != nil {
		fmt.Printf("WriteSliceToCSVFs gagal: %v\n", err)
		return
	}
	fmt.Println("berhasil tulis CSV dari raw slice")

	// ─── Versi OS (tanpa afero) — langsung ke disk ───
	// fileutil.WriteStructsToCSV(karyawans, "/tmp/karyawan.csv")
	// hasil, _ := fileutil.ReadCSVToStructs[Karyawan]("/tmp/karyawan.csv")
	// fileutil.WriteSliceToCSV(rawData, "/tmp/produk.csv")
}

// ──────────────────────────────────────────────────────────
// Contoh 3: ZIP — ZipFiles, UnzipTo
// ──────────────────────────────────────────────────────────

// ContohFileUtilZip mendemonstrasikan ZipFilesFs dan UnzipToFs.
func ContohFileUtilZip() {
	fs := afero.NewMemMapFs()

	// Buat beberapa file yang akan di-zip
	fileutil.WriteLines(fs, "/source/readme.txt", []string{"ini readme", "baris dua"})
	fileutil.WriteLines(fs, "/source/data/users.csv", []string{"id,name", "1,Budi"})
	fileutil.WriteLines(fs, "/source/data/orders.csv", []string{"id,total", "1,50000"})
	fileutil.WriteLines(fs, "/source/config.json", []string{`{"env":"prod"}`})

	// ─── ZipFilesFs ─── buat arsip ZIP dari beberapa source
	zipPath := "/output/archive.zip"
	err := fileutil.ZipFilesFs(fs, zipPath, []string{
		"/source/readme.txt",
		"/source/data",     // direktori dizip rekursif
		"/source/config.json",
	})
	if err != nil {
		fmt.Printf("ZipFilesFs gagal: %v\n", err)
		return
	}
	fmt.Println("berhasil membuat arsip zip di", zipPath)
	fmt.Println("file zip ada?", fileutil.FileExists(fs, zipPath))

	// ─── UnzipToFs ─── ekstrak arsip ke direktori tujuan
	extractDir := "/extracted"
	err = fileutil.UnzipToFs(fs, zipPath, extractDir)
	if err != nil {
		fmt.Printf("UnzipToFs gagal: %v\n", err)
		return
	}
	fmt.Println("berhasil mengekstrak arsip ke", extractDir)

	// Verifikasi hasil ekstrak
	fmt.Println("DirExists /extracted:", fileutil.DirExists(fs, extractDir))

	// ─── Versi OS ───
	// fileutil.ZipFiles("/tmp/backup.zip", []string{"/var/log/app", "/etc/myapp"})
	// fileutil.UnzipTo("/tmp/backup.zip", "/tmp/restore")
}

// ──────────────────────────────────────────────────────────
// Contoh 4: FTP — NewFTPClient, Upload, Download, Delete, List, MakeDir
// ──────────────────────────────────────────────────────────
//
// Note: Bagian ini butuh FTP server yang beneran buat jalan.
// Dalam contoh ini kita nunjukin cara pakainya aja.

// ContohFTPClient mendemonstrasikan semua method FTPClient.
// Ganti cfg dengan koneksi FTP yang beneran di environment lo.
func ContohFTPClient() {
	// ─── Buat koneksi FTP ───
	cfg := fileutil.FTPConfig{
		Host:     "ftp.server.contoh.id",
		Port:     21,              // default 21, bisa dikosongkan
		Username: "ftpuser",
		Password: "ftppassword",
		Timeout:  30 * time.Second, // default 30 detik, bisa dikosongkan
	}

	client, err := fileutil.NewFTPClient(cfg)
	if err != nil {
		fmt.Printf("gagal konek ke FTP: %v\n", err)
		return
	}
	defer client.Close() // selalu tutup sesi FTP

	fmt.Println("berhasil login ke FTP server")

	// ─── MakeDir ─── buat direktori di server
	err = client.MakeDir("/upload/2024/01")
	if err != nil {
		fmt.Printf("MakeDir gagal (mungkin udah ada): %v\n", err)
	} else {
		fmt.Println("direktori /upload/2024/01 berhasil dibuat")
	}

	// ─── Upload ─── kirim file lokal ke server FTP (streaming, hemat memori)
	// Buat file lokal dulu buat demo
	localFile := "/tmp/report_januari.csv"
	os.WriteFile(localFile, []byte("id,tanggal,nominal\n1,2024-01-01,50000\n"), 0644)
	defer os.Remove(localFile)

	err = client.Upload(localFile, "/upload/2024/01/report_januari.csv")
	if err != nil {
		fmt.Printf("Upload gagal: %v\n", err)
	} else {
		fmt.Println("file berhasil diupload")
	}

	// ─── List ─── lihat isi direktori remote
	entries, err := client.List("/upload/2024/01")
	if err != nil {
		fmt.Printf("List gagal: %v\n", err)
	} else {
		fmt.Printf("isi direktori /upload/2024/01 (%d entri):\n", len(entries))
		for _, name := range entries {
			fmt.Printf("  - %s\n", name)
		}
	}

	// ─── Download ─── unduh file dari server ke lokal (streaming, hemat memori)
	localDownload := "/tmp/downloaded_report.csv"
	err = client.Download("/upload/2024/01/report_januari.csv", localDownload)
	if err != nil {
		fmt.Printf("Download gagal: %v\n", err)
	} else {
		fmt.Println("file berhasil didownload ke", localDownload)
		defer os.Remove(localDownload)
	}

	// ─── Delete ─── hapus file di server
	err = client.Delete("/upload/2024/01/report_januari.csv")
	if err != nil {
		fmt.Printf("Delete gagal: %v\n", err)
	} else {
		fmt.Println("file berhasil dihapus dari server FTP")
	}
}
