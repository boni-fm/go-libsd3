package example

// ======================================================
// Contoh pemakaian package postgres (pkg/db/postgres)
// ======================================================
//
// Di sini kita contohnin semua fitur yang ada di paket postgres:
//   - NewDatabase — buat koneksi fresh ke Postgres
//   - ConnectionPool (singleton) — kelola banyak koneksi berdasarkan kode DC
//   - ConnectionManager — wrapper praktis buat ConnectionPool
//   - Registry (GetDB/RegisterDB/MustGetDB/CloseAll) — multiton per (kunci, kodedc)
//   - Query/QueryRows/Exec — raw pgx query
//   - SelectAll/SelectOne — scanning pake scany
//   - Insert/InsertBatch/CopyFrom — insert & bulk
//   - BeginTx/ExecuteInTransaction — transaksi
//   - ExportQueryToCSV — export hasil query ke CSV
//   - Ping/GetStats/IsClosed/GetStartTime dsb — connection management
//
// Note:
//   - Semua fungsi ini butuh koneksi Postgres beneran buat jalan.
//   - Di sini cuma nunjukin CARA pakenya, bukan actual run.
//   - Ganti placeholder DSN/kunci/kodeDC sesuai environment lo.

import (
	"context"
	"fmt"
	"time"

	"github.com/boni-fm/go-libsd3/pkg/db/postgres"
	"github.com/jackc/pgx/v5"
)

// ──────────────────────────────────────────────────────────
// Contoh 1: NewDatabase — buat koneksi tunggal ke Postgres
// ──────────────────────────────────────────────────────────
//
// NewDatabase bakal ambil connection string dari SettingLib (lewat env KunciWeb
// atau config.yaml). Kalau KodeDC dikasih, dia manggil InitConstrByKodeDc;
// kalau ga, manggil InitConstr yang baca dari env/yaml langsung.
func ContohNewDatabase() {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// config minimal — KodeDC wajib kalo mau bedain per-DC
	cfg := postgres.Config{
		KodeDC:  "G009SIM",
		AppName: "contoh-app",
		// opsional — kalau ga diisi bakal pake default value
		MaxConns:        10,
		MinConns:        2,
		ConnMaxLifetime: 30 * time.Minute,
		MaxConnIdleTime: 10 * time.Minute,
	}

	db, err := postgres.NewDatabase(ctx, &cfg)
	if err != nil {
		fmt.Printf("gagal konek ke db: %v\n", err)
		return
	}
	defer db.Close() // jangan lupa tutup!

	fmt.Println("database connected, start time:", db.GetStartTime())
	fmt.Println("kode DC:", db.GetKodeDc())
	fmt.Println("conn string:", db.GetConnString())
	fmt.Println("sudah ditutup?", db.IsClosed())

	// Ping buat ngecek koneksi masih sehat
	if err := db.Ping(ctx); err != nil {
		fmt.Printf("ping gagal: %v\n", err)
		return
	}
	fmt.Println("ping sukses!")

	// Liat statistik pool
	stats := db.GetStats()
	fmt.Printf("total conns: %d | idle: %d | acquired: %d | max: %d\n",
		stats.TotalConns(), stats.IdleConns(), stats.AcquiredConns(), stats.MaxConns())

	// Ambil waktu server Postgres
	serverTime, err := db.GetServerTime(ctx)
	if err != nil {
		fmt.Printf("gagal dapetin waktu server: %v\n", err)
	} else {
		fmt.Println("waktu server:", serverTime)
	}

	// Ambil kode DC dari tabel dc_tabel_dc_t (kalo tabelnya ada)
	kodeDcDB, err := db.GetKodeDcFromDB(ctx)
	if err != nil {
		fmt.Printf("gagal dapetin kode DC dari DB: %v\n", err)
	} else {
		fmt.Println("kode DC dari DB:", kodeDcDB)
	}
}

// ──────────────────────────────────────────────────────────
// Contoh 2: Query, QueryRows, Exec — raw pgx query
// ──────────────────────────────────────────────────────────
func ContohQueryRaw(db *postgres.Database) {
	ctx := context.Background()

	// Query — single row result (pake pgx.Row)
	row := db.Query(ctx, "SELECT NOW()")
	if row == nil {
		fmt.Println("koneksi udah ditutup bos")
		return
	}
	var now time.Time
	if err := row.Scan(&now); err != nil {
		fmt.Printf("scan gagal: %v\n", err)
		return
	}
	fmt.Println("waktu sekarang:", now)

	// QueryRows — multiple rows
	rows, err := db.QueryRows(ctx, "SELECT id, name FROM users WHERE active = $1", true)
	if err != nil {
		fmt.Printf("query rows gagal: %v\n", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			fmt.Printf("scan row gagal: %v\n", err)
			continue
		}
		fmt.Printf("user: id=%d name=%s\n", id, name)
	}

	// Exec — INSERT/UPDATE/DELETE, return jumlah baris kena
	rowsAffected, err := db.Exec(ctx,
		"UPDATE users SET last_seen = $1 WHERE id = $2",
		time.Now(), 42,
	)
	if err != nil {
		fmt.Printf("exec gagal: %v\n", err)
		return
	}
	fmt.Printf("%d baris kena update\n", rowsAffected)
}

// ──────────────────────────────────────────────────────────
// Contoh 3: SelectAll / SelectOne — scanning pake scany
// ──────────────────────────────────────────────────────────

// User adalah contoh struct yang dipakai buat scanning.
type User struct {
	ID    int    `db:"id"`
	Name  string `db:"name"`
	Email string `db:"email"`
}

func ContohSelectScany(db *postgres.Database) {
	ctx := context.Background()

	// SelectAll — hasil banyak baris masuk ke slice
	var users []User
	err := db.SelectAll(ctx, &users,
		"SELECT id, name, email FROM users WHERE active = $1 LIMIT 10", true)
	if err != nil {
		fmt.Printf("SelectAll gagal: %v\n", err)
		return

	}
	for _, u := range users {
		fmt.Printf("user: %+v\n", u)
	}

	// SelectOne — single row masuk ke struct
	var u User
	err = db.SelectOne(ctx, &u,
		"SELECT id, name, email FROM users WHERE id = $1", 1)
	if err != nil {
		fmt.Printf("SelectOne gagal: %v\n", err)
		return
	}
	fmt.Printf("user tunggal: %+v\n", u)

	// SelectRowCallback — proses baris satu-satu tanpa load semua ke memori
	err = db.SelectRowCallback(ctx,
		"SELECT id, name, email FROM users",
		func(ctx context.Context, rows pgx.Rows) error {
			for rows.Next() {
				var u User
				if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
					return err
				}
				fmt.Printf("streaming user: %+v\n", u)
			}
			return rows.Err()
		},
	)
	if err != nil {
		fmt.Printf("SelectRowCallback gagal: %v\n", err)
	}
}

// ──────────────────────────────────────────────────────────
// Contoh 4: ScanAllRows / ScanOneRow / ScanRow
// ──────────────────────────────────────────────────────────
func ContohScanHelpers(db *postgres.Database) {
	ctx := context.Background()

	// ScanAllRows — scan dari pgx.Rows yang udah ada
	rows, err := db.QueryRows(ctx, "SELECT id, name, email FROM users")
	if err != nil {
		fmt.Printf("QueryRows gagal: %v\n", err)
		return
	}
	defer rows.Close()

	var users []User
	if err := db.ScanAllRows(rows, &users); err != nil {
		fmt.Printf("ScanAllRows gagal: %v\n", err)
		return
	}
	fmt.Printf("total users dari ScanAllRows: %d\n", len(users))

	// ScanOneRow — mirip ScanAllRows tapi cuma satu baris
	rows2, err := db.QueryRows(ctx, "SELECT id, name, email FROM users WHERE id = $1", 1)
	if err != nil {
		fmt.Printf("QueryRows gagal: %v\n", err)
		return
	}
	defer rows2.Close()

	var u User
	if err := db.ScanOneRow(rows2, &u); err != nil {
		fmt.Printf("ScanOneRow gagal: %v\n", err)
		return
	}
	fmt.Printf("user dari ScanOneRow: %+v\n", u)
}

// ──────────────────────────────────────────────────────────
// Contoh 5: Insert / InsertBatch / CopyFrom
// ──────────────────────────────────────────────────────────
func ContohInsert(db *postgres.Database) {
	ctx := context.Background()

	// Insert — masukin satu baris, dapet last inserted ID
	id, err := db.Insert(ctx,
		"INSERT INTO users(name, email) VALUES($1, $2) RETURNING id",
		"Budi Santoso", "budi@contoh.id",
	)
	if err != nil {
		fmt.Printf("Insert gagal: %v\n", err)
		return
	}
	fmt.Printf("user baru ID: %d\n", id)

	// InsertBatch — masukin banyak baris dalam satu round-trip ke DB
	records := [][]interface{}{
		{"Andi", "andi@contoh.id"},
		{"Siti", "siti@contoh.id"},
		{"Rudi", "rudi@contoh.id"},
	}
	lastID, err := db.InsertBatch(ctx,
		"INSERT INTO users(name, email) VALUES($1, $2) RETURNING id",
		records,
	)
	if err != nil {
		fmt.Printf("InsertBatch gagal: %v\n", err)
		return
	}
	fmt.Printf("batch insert done, last rows affected: %d\n", lastID)

	// CopyFrom — cara tercepat buat masukin data gede (COPY FROM)
	bulkData := [][]interface{}{
		{"Dewi", "dewi@contoh.id"},
		{"Joko", "joko@contoh.id"},
	}
	rowsCopied, err := db.CopyFrom(ctx,
		"users",
		[]string{"name", "email"},
		bulkData,
	)
	if err != nil {
		fmt.Printf("CopyFrom gagal: %v\n", err)
		return
	}
	fmt.Printf("CopyFrom berhasil masukin %d baris\n", rowsCopied)
}

// ──────────────────────────────────────────────────────────
// Contoh 6: Transaksi — BeginTx, BeginTxWithOptions, ExecuteInTransaction
// ──────────────────────────────────────────────────────────
func ContohTransaksi(db *postgres.Database) {
	ctx := context.Background()

	// BeginTx — mulai transaksi manual
	tx, err := db.BeginTx(ctx)
	if err != nil {
		fmt.Printf("BeginTx gagal: %v\n", err)
		return
	}

	// Jangan lupa rollback kalo ada error (pake defer)
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback(ctx)
			panic(p)
		}
	}()

	_, err = tx.Exec(ctx, "UPDATE saldo SET jumlah = jumlah - $1 WHERE user_id = $2", 100000, 1)
	if err != nil {
		tx.Rollback(ctx)
		fmt.Printf("exec dalam tx gagal: %v\n", err)
		return
	}

	_, err = tx.Exec(ctx, "UPDATE saldo SET jumlah = jumlah + $1 WHERE user_id = $2", 100000, 2)
	if err != nil {
		tx.Rollback(ctx)
		fmt.Printf("exec dalam tx gagal: %v\n", err)
		return
	}

	if err := tx.Commit(ctx); err != nil {
		fmt.Printf("commit gagal: %v\n", err)
		return
	}
	fmt.Println("transfer sukses!")

	// BeginTxWithOptions — transaksi dengan isolation level khusus
	txOpts := pgx.TxOptions{IsoLevel: pgx.Serializable}
	tx2, err := db.BeginTxWithOptions(ctx, txOpts)
	if err != nil {
		fmt.Printf("BeginTxWithOptions gagal: %v\n", err)
		return
	}
	tx2.Rollback(ctx) // contoh rollback langsung

	// ExecuteInTransaction — cara paling praktis, commit/rollback otomatis
	err = db.ExecuteInTransaction(ctx, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, "INSERT INTO log_aktivitas(pesan) VALUES($1)", "user login")
		if err != nil {
			return fmt.Errorf("insert log gagal: %w", err)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("ExecuteInTransaction gagal: %v\n", err)
	} else {
		fmt.Println("transaksi otomatis sukses!")
	}
}

// ──────────────────────────────────────────────────────────
// Contoh 7: ExportQueryToCSV — export hasil query ke CSV
// ──────────────────────────────────────────────────────────
func ContohExportCSV(db *postgres.Database) {
	ctx := context.Background()

	// default separator koma
	err := db.ExportQueryToCSV(ctx, "/tmp/users.csv",
		"SELECT id, name, email FROM users ORDER BY id LIMIT 100", "")
	if err != nil {
		fmt.Printf("ExportQueryToCSV gagal: %v\n", err)
		return
	}
	fmt.Println("CSV berhasil diekspor ke /tmp/users.csv")

	// pake separator titik koma
	err = db.ExportQueryToCSV(ctx, "/tmp/users_semicolon.csv",
		"SELECT id, name, email FROM users ORDER BY id LIMIT 100", ";")
	if err != nil {
		fmt.Printf("ExportQueryToCSV (semicolon) gagal: %v\n", err)
		return
	}
	fmt.Println("CSV (semicolon) berhasil diekspor ke /tmp/users_semicolon.csv")
}

// ──────────────────────────────────────────────────────────
// Contoh 8: ConnectionPool — singleton pool buat multi-DC
// ──────────────────────────────────────────────────────────
func ContohConnectionPool() {
	ctx := context.Background()

	// Dapetin singleton pool
	pool := postgres.GetConnectionPool()

	// Daftarin konfigurasi per kode DC
	cfgDC1 := postgres.Config{KodeDC: "G009SIM", AppName: "my-app"}
	cfgDC2 := postgres.Config{KodeDC: "G010SIM", AppName: "my-app"}

	if err := pool.RegisterConfig("G009SIM", cfgDC1); err != nil {
		fmt.Printf("RegisterConfig DC1 gagal: %v\n", err)
	}
	if err := pool.RegisterConfig("G010SIM", cfgDC2); err != nil {
		fmt.Printf("RegisterConfig DC2 gagal: %v\n", err)
	}

	// Connect — lazy connect, koneksi dibuat kalau belum ada
	db1, err := pool.Connect(ctx, "G009SIM")
	if err != nil {
		fmt.Printf("Connect DC1 gagal: %v\n", err)
		return
	}
	fmt.Println("terkonek ke DC1:", db1.GetKodeDc())

	// HasConnection — cek apakah koneksi sudah ada
	fmt.Println("DC1 ada di pool?", pool.HasConnection("G009SIM"))
	fmt.Println("DC99 ada di pool?", pool.HasConnection("G099SIM"))

	// GetAllKodeDcPoolKey — dapetin semua kunci yang terdaftar
	keys := pool.GetAllKodeDcPoolKey()
	fmt.Println("kunci di pool:", keys)

	// GetConnectionStats — statistik koneksi per DC
	stats, err := pool.GetConnectionStats(ctx, "G009SIM")
	if err != nil {
		fmt.Printf("GetConnectionStats gagal: %v\n", err)
	} else {
		fmt.Printf("stats DC1: %v\n", stats)
	}

	// HealthCheck — ping semua koneksi
	healthResults := pool.HealthCheckAll(ctx)
	for dc, err := range healthResults {
		if err == nil {
			fmt.Printf("  ✓ %s: sehat\n", dc)
		} else {
			fmt.Printf("  ✗ %s: sakit — %v\n", dc, err)
		}
	}

	// UpdateConfig — update konfigurasi tanpa restart
	pool.UpdateConfig("G009SIM", postgres.Config{KodeDC: "G009SIM", AppName: "updated-app"})

	// ReinitializeConnection — tutup & buat ulang koneksi (klo ada masalah)
	if err := pool.ReinitializeConnection(ctx, "G009SIM"); err != nil {
		fmt.Printf("reinit gagal: %v\n", err)
	}

	// CloseConnection — tutup satu koneksi
	if err := pool.CloseConnection("G010SIM"); err != nil {
		fmt.Printf("CloseConnection DC2 gagal: %v\n", err)
	}

	// CloseAll / Close — tutup semua sekaligus
	defer pool.Close()
}

// ──────────────────────────────────────────────────────────
// Contoh 9: ConnectionManager — wrapper praktis pool
// ──────────────────────────────────────────────────────────
func ContohConnectionManager() {
	ctx := context.Background()

	cm := postgres.NewConnectionManager()

	// InitializeConnections — daftarin + pre-connect beberapa DC sekaligus
	configs := map[string]postgres.Config{
		"G009SIM": {KodeDC: "G009SIM", AppName: "my-app"},
		"G010SIM": {KodeDC: "G010SIM", AppName: "my-app"},
	}
	if err := cm.InitializeConnections(ctx, configs); err != nil {
		fmt.Printf("InitializeConnections gagal: %v\n", err)
		return
	}

	// GetDB — ambil koneksi per kode DC
	db, err := cm.GetDB(ctx, "G009SIM")
	if err != nil {
		fmt.Printf("GetDB gagal: %v\n", err)
		return
	}
	fmt.Println("dapet koneksi:", db.GetKodeDc())

	// PrintPoolStatus — cetak status pool ke stdout
	cm.PrintPoolStatus(ctx)

	// HealthCheck — cek kesehatan semua
	results := cm.HealthCheck(ctx)
	for dc, err := range results {
		if err != nil {
			fmt.Printf("DC %s sakit: %v\n", dc, err)
		}
	}
	cm.PrintHealthStatus(ctx)

	// Close satu koneksi
	if err := cm.Close("G010SIM"); err != nil {
		fmt.Printf("Close DC2 gagal: %v\n", err)
	}

	// CloseAllConnections — tutup semua pas shutdown
	defer cm.CloseAllConnections()
}

// ──────────────────────────────────────────────────────────
// Contoh 10: Registry (GetDB/RegisterDB/MustGetDB/CloseAll)
// ──────────────────────────────────────────────────────────
//
// Registry adalah multiton per pasangan (kunci, kodedc).
// Cocok buat deployment multi-tenant di mana satu service
// melayani banyak tenant (kunci) dengan banyak DC.
func ContohRegistry() {
	ctx := context.Background()

	cfg := postgres.Config{
		KodeDC:  "G009SIM",
		AppName: "tenant-app",
	}

	// RegisterDB — daftarin + buat koneksi baru ke registry
	db, err := postgres.RegisterDB(ctx, "kunciG009", "G009SIM", cfg)
	if err != nil {
		fmt.Printf("RegisterDB gagal: %v\n", err)
		return
	}
	fmt.Println("berhasil registrasi tenant:", db.GetKodeDc())

	// GetDB — ambil koneksi yang udah terdaftar
	db2, err := postgres.GetDB("kunciG009", "G009SIM")
	if err != nil {
		fmt.Printf("GetDB gagal: %v\n", err)
		return
	}
	fmt.Println("GetDB sukses:", db2.GetKodeDc())

	// MustGetDB — panic kalau ga ketemu (pake klo yakin udah terdaftar)
	db3 := postgres.MustGetDB("kunciG009", "G009SIM")
	fmt.Println("MustGetDB sukses:", db3.GetKodeDc())

	// CloseAll — tutup semua koneksi di registry
	if err := postgres.CloseAll(); err != nil {
		fmt.Printf("CloseAll gagal: %v\n", err)
	} else {
		fmt.Println("semua koneksi registry ditutup")
	}
}
