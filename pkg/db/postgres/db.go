package postgres

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/boni-fm/go-libsd3/pkg/config/constant"
	"github.com/boni-fm/go-libsd3/pkg/settinglibgo"
	"github.com/boni-fm/go-libsd3/pkg/yaml"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// -=-=-=-=-=-=-=-=-=-
// TODO:
// - tambahin transaction untuk eksekusi query kyk di c#
// - tambahin query builder
// - tambahin logging, kalo bisa make hooks ? ~
// - pelajarin fungsi hooks

// Struct buat config database nya 🔥
type Config struct {
	// dc ~ 🏢
	KodeDC string

	// dibawah ini optional,
	// nanti di initialize connection akan ada default value nya
	// kalo bisa diisi dari awal buat confignya, biar jelas
	MaxConns        int
	MinConns        int
	ConnMaxLifetime time.Duration
	MaxConnIdleTime time.Duration

	// buat ngecek kesehatan debe,
	// kalo meleduk jadi ketauan
	// TODO:
	// - add alert kalo selama healthcheck gagal
	HealthCheckInterval time.Duration
}

// struct database nya 🔥
// ini di-initialize pas buat connection db nya
type Database struct {
	Pool       *pgxpool.Pool
	ConfigDB   *Config
	ConnString string
	mu         sync.RWMutex
	isClosed   bool
	startTime  time.Time
}

// 😡 ~ di cache biar gk ada 🐛 bug
var connStrCache = struct {
	sync.RWMutex
	m map[string]string
}{
	m: make(map[string]string),
}

// -=-=-=-=-=-=-=-=-=-
// INITIALIZE METHODS
// -=-=-=-=-=-=-=-=-=-

// fungsi inisiasi koneksi database baru
// note:
// - bisa dijalankan tanpa connection pool nya
// - config nya gk harus semua diisi (ada default value)
// - ambil connstring dari settinglib, baca kunci nginx docker atau non-docker
func NewDatabase(ctx context.Context, cfg *Config) (*Database, error) {
	connStr, err := func() (string, error) {
		if cfg.KodeDC != "" {
			return InitConstrByKodeDc(ctx, cfg.KodeDC)
		}
		return InitConstr(ctx)
	}()
	if err != nil {
		return nil, err
	}

	// dapetin kunci dari settinglib
	db := &Database{
		ConfigDB:   initDefaultConfig(cfg),
		ConnString: connStr,
	}

	poolConfig, err := db.GetPool(ctx)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db.Pool = pool
	db.startTime = time.Now()
	db.isClosed = false

	return db, nil
}

func InitConstr(ctx context.Context) (string, error) {
	var strKunci string
	kodeDc := "Gxxx"

	if strKunciDocker := os.Getenv(constant.KEY_ENV_KUNCI); strKunciDocker != "" {
		strKunci = strKunciDocker
	} else {
		kuncipath, _ := yaml.GetKunciConfigFilepath()
		if _, err := os.Stat(kuncipath); err != nil {
			fmt.Printf("[WARN] Config file does not exist at %s: %v\n", kuncipath, err)
		} else {
			keyYaml, readErr := yaml.ReadConfigDynamicWithKey(kuncipath, "kunci")
			if readErr != nil {
				fmt.Printf("[WARN] Failed to read 'kunci' from YAML: %v\n", readErr)
			} else if keyYaml == nil {
				fmt.Printf("[WARN] 'kunci' key is nil in YAML\n")
			} else {
				if kunciValue, ok := keyYaml.(string); ok && kunciValue != "" {
					strKunci = kunciValue
					kodeDc = strings.ToUpper(strings.TrimPrefix(kunciValue, "kunci"))
				} else {
					fmt.Printf("[WARN] 'kunci' is not a non-empty string in YAML, got type: %T, value: %v\n", keyYaml, keyYaml)
				}
			}
		}
	}

	connStrCache.RLock()
	if cs, ok := connStrCache.m[kodeDc]; ok {
		connStrCache.RUnlock()
		return cs, nil
	}
	connStrCache.RUnlock()

	if strKunci == "" {
		strKunci = "kunci" + strings.ToLower(kodeDc)
		fmt.Printf("[WARN] Using fallback kunci: %s\n", strKunci)
	}

	// Allow context cancellation
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	// Get connection string from SettingLib
	kunciManager := settinglibgo.NewSettingLib(strKunci)
	constr := kunciManager.GetConnectionString(constant.DBTYPE_POSTGRE)
	if constr == "" {
		return "", fmt.Errorf("empty connection string for kunci=%s", strKunci)
	}

	// Cache result
	connStrCache.Lock()
	connStrCache.m[kodeDc] = constr
	connStrCache.Unlock()

	return constr, nil
}

// dapetin connection string dari settinglib
// baca kunci dari env variable (docker) atau file yaml (non-docker)
func InitConstrByKodeDc(ctx context.Context, kodeDc string) (string, error) {
	// Check cache first
	connStrCache.RLock()
	if cs, ok := connStrCache.m[kodeDc]; ok {
		connStrCache.RUnlock()
		return cs, nil
	}
	connStrCache.RUnlock()

	var strKunci string

	// Priority 1: Environment variable
	if strKunciDocker := os.Getenv(constant.KEY_ENV_KUNCI); strKunciDocker != "" {
		strKunci = strKunciDocker
		fmt.Printf("[INFO] Using kunci from ENV: %s\n", strKunci)
	} else {
		// Priority 2: YAML config file
		kuncipath, err := yaml.GetKunciConfigFilepath()
		if err != nil {
			fmt.Printf("[WARN] Failed to get kunci config filepath: %v\n", err)
		} else if kuncipath == "" {
			fmt.Printf("[WARN] Kunci config filepath is empty\n")
		} else {
			fmt.Printf("[DEBUG] Kunci config path: %s\n", kuncipath)

			// Verify file exists
			if _, err := os.Stat(kuncipath); err != nil {
				fmt.Printf("[WARN] Config file does not exist at %s: %v\n", kuncipath, err)
			} else {
				keyYaml, readErr := yaml.ReadConfigDynamicWithKey(kuncipath, "kunci")
				if readErr != nil {
					fmt.Printf("[WARN] Failed to read 'kunci' from YAML: %v\n", readErr)
				} else if keyYaml == nil {
					fmt.Printf("[WARN] 'kunci' key is nil in YAML\n")
				} else {
					if kunciValue, ok := keyYaml.(string); ok && kunciValue != "" {
						strKunci = kunciValue
						fmt.Printf("[INFO] Using kunci from YAML: %s\n", strKunci)
					} else {
						fmt.Printf("[WARN] 'kunci' is not a non-empty string in YAML, got type: %T, value: %v\n", keyYaml, keyYaml)
					}
				}
			}
		}
	}

	// Priority 3: Fallback to default
	if strKunci == "" {
		strKunci = "kunci" + strings.ToLower(kodeDc)
		fmt.Printf("[WARN] Using fallback kunci: %s\n", strKunci)
	}

	fmt.Printf("[DEBUG] Final kunci key: %s\n", strKunci)

	// Allow context cancellation
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	default:
	}

	// Get connection string from SettingLib
	kunciManager := settinglibgo.NewSettingLib(strKunci)
	constr := kunciManager.GetConnectionString(constant.DBTYPE_POSTGRE)
	if constr == "" {
		return "", fmt.Errorf("empty connection string for kunci=%s", strKunci)
	}

	// Cache result
	connStrCache.Lock()
	connStrCache.m[kodeDc] = constr
	connStrCache.Unlock()

	return constr, nil
}

// isi default value dari confignya
// kalo diisi sendiri dari awal yaudah gpp 👍
func initDefaultConfig(cfg *Config) *Config {
	if cfg.MaxConns == 0 {
		cfg.MaxConns = 10
	}

	if cfg.MinConns == 0 {
		cfg.MinConns = 2
	}

	if cfg.ConnMaxLifetime == 0 {
		cfg.ConnMaxLifetime = 30 * time.Minute
	}

	if cfg.MaxConnIdleTime == 0 {
		cfg.MaxConnIdleTime = 10 * time.Minute
	}

	if cfg.HealthCheckInterval == 0 {
		cfg.HealthCheckInterval = 5 * time.Minute
	}

	return cfg
}

// dapetin config pgx pool (bawaan pgx)
// disesuain config yg dipunya dengan config punya pgx
// mulai koneksi make config mereka
func (m *Database) GetPool(ctx context.Context) (*pgxpool.Config, error) {
	config, err := pgxpool.ParseConfig(m.ConnString)
	if err != nil {
		return nil, err
	}

	config.MaxConns = int32(m.ConfigDB.MaxConns)
	config.MinConns = int32(m.ConfigDB.MinConns)
	config.MaxConnLifetime = m.ConfigDB.ConnMaxLifetime
	config.MaxConnIdleTime = m.ConfigDB.MaxConnIdleTime

	return config, nil
}

// -=-=-=-=-=-=-=-=-=-
// QUERIES METHODS
// -=-=-=-=-=-=-=-=-=-

// kumpulan fungsi query bawaan dari pgx tanpa make scanny
// kalo make ini sama aja kyk make fungsi si pgx
// bedanya kalo disini dibantu mutex biar thread-safe (ceritanya)
// jadi kalo jalan di gorutine gk bakal tabrakan
// TODO:
// - logging?

// single row query result ~
func (d *Database) Query(ctx context.Context, query string, args ...interface{}) pgx.Row {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.isClosed {
		return nil
	}

	return d.Pool.QueryRow(ctx, query, args...)
}

// multiple rows query result ~
func (d *Database) QueryRows(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.isClosed {
		return nil, fmt.Errorf("database connection is closed")
	}

	rows, err := d.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}

	return rows, nil
}

// ini buat execute query (insert, update, delete)
// hasil return, jumlah rows yg terpengaruh
func (d *Database) Exec(ctx context.Context, query string, args ...interface{}) (int64, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.isClosed {
		return 0, fmt.Errorf("database connection is closed")
	}

	result, err := d.Pool.Exec(ctx, query, args...)
	if err != nil {
		return 0, fmt.Errorf("error executing query: %w", err)
	}

	return result.RowsAffected(), nil
}

// -=-=-=-=-=-=-=-=-=-
// SCANNING METHODS
// -=-=-=-=-=-=-=-=-=-

// kalo bagian ini itu fungsi select yang bisa masuk ke dalem variable
// bisa single row atau multiple rows
// scan nya make library exeternal scany
// TODO:
// - pastiin gk perlu ada tipe data spesifik (sql.NullString dst)

// SelectAll scans multiple rows into a slice using scany
// Usage: db.SelectAll(ctx, &users, "SELECT * FROM users WHERE status = $1", "active")
func (d *Database) SelectAll(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.isClosed {
		return fmt.Errorf("database connection is closed")
	}

	err := pgxscan.Select(ctx, d.Pool, dest, query, args...)
	if err != nil {
		return fmt.Errorf("error scanning all rows: %w", err)
	}

	return nil
}

// SelectOne scans a single row using scany
// Usage: db.SelectOne(ctx, &user, "SELECT * FROM users WHERE id = $1", 1)
func (d *Database) SelectOne(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.isClosed {
		return fmt.Errorf("database connection is closed")
	}

	err := pgxscan.Get(ctx, d.Pool, dest, query, args...)
	if err != nil {
		return fmt.Errorf("error scanning single row: %w", err)
	}

	return nil
}

// ScanAllRows scans all rows from pgx.Rows using scany
// Usage: rows, _ := db.QueryRows(ctx, query, args...); db.ScanAllRows(ctx, rows, &users)
func (d *Database) ScanAllRows(ctx context.Context, rows pgx.Rows, dest interface{}) error {
	err := pgxscan.ScanAll(dest, rows)
	if err != nil {
		return fmt.Errorf("error scanning all rows: %w", err)
	}
	return nil
}

// ScanOneRow scans a single row from pgx.Rows using scany
// Usage: rows, _ := db.QueryRows(ctx, query, args...); db.ScanOneRow(ctx, rows, &user)
func (d *Database) ScanOneRow(ctx context.Context, rows pgx.Rows, dest interface{}) error {
	err := pgxscan.ScanOne(dest, rows)
	if err != nil {
		return fmt.Errorf("error scanning single row: %w", err)
	}
	return nil
}

// ScanRow scans a single row directly from pgx.Row
// Usage: row := db.Query(ctx, query, args...); db.ScanRow(ctx, &user, row)
func (d *Database) ScanRow(ctx context.Context, dest interface{}, row pgx.Rows) error {
	err := pgxscan.ScanRow(dest, row)
	if err != nil {
		return fmt.Errorf("error scanning row: %w", err)
	}
	return nil
}

// SelectRowCallback executes a query and calls a callback for each row
// Useful for processing rows one at a time without loading all into memory
func (d *Database) SelectRowCallback(ctx context.Context, query string, fn func(context.Context, pgx.Rows) error, args ...interface{}) error {
	rows, err := d.QueryRows(ctx, query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	return fn(ctx, rows)
}

// -=-=-=-=-=-=-=-=-=-=-=-=-=-=-=
// Insert and Bulk Insert METHODS
// -=-=-=-=-=-=-=-=-=-=-=-=-=-=-=

// buat insert sama bulk insert ke tabel
// alternatif dari copyhelper di c# 👩‍💻

// Insert inserts a record and returns the last inserted ID
// Usage: id, err := db.Insert(ctx, "INSERT INTO users(name, email) VALUES($1, $2) RETURNING id", "John", "john@example.com")
func (d *Database) Insert(ctx context.Context, query string, args ...interface{}) (int64, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.isClosed {
		return 0, fmt.Errorf("database connection is closed")
	}

	var id int64
	err := d.Pool.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error inserting record: %w", err)
	}

	return id, nil
}

// InsertBatch inserts multiple records using batch operation (single round-trip to DB)
// Usage: db.InsertBatch(ctx, "INSERT INTO users(name, email) VALUES($1, $2)", [][]interface{}{{"John", "john@example.com"}, {"Jane", "jane@example.com"}})
func (d *Database) InsertBatch(ctx context.Context, query string, records [][]interface{}) (int64, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.isClosed {
		return 0, fmt.Errorf("database connection is closed")
	}

	batch := &pgx.Batch{}
	for _, record := range records {
		batch.Queue(query, record...)
	}

	results := d.Pool.SendBatch(ctx, batch)
	defer results.Close()

	var lastID int64
	for i := 0; i < len(records); i++ {
		tag, err := results.Exec()
		if err != nil {
			return 0, fmt.Errorf("error executing batch insert: %w", err)
		}

		// Extract ID from the tag if available
		if i == len(records)-1 {
			if rows := tag.RowsAffected(); rows > 0 {
				lastID = int64(rows)
			}
		}
	}

	return lastID, nil
}

// CopyFrom performs a bulk copy operation (COPY FROM)
// This is the fastest way to insert large amounts of data (10-100x faster than INSERT)
// Usage: rows, err := db.CopyFrom(ctx, "users", []string{"name", "email"}, [][]interface{}{{"John", "john@example.com"}})
func (d *Database) CopyFrom(ctx context.Context, tableName string, columnNames []string, records [][]interface{}) (int64, error) {
	d.mu.RLock()
	conn, err := d.Pool.Acquire(ctx)
	d.mu.RUnlock()

	if err != nil {
		return 0, fmt.Errorf("error acquiring connection: %w", err)
	}
	defer conn.Release()

	if d.isClosed {
		return 0, fmt.Errorf("database connection is closed")
	}

	rows, err := conn.Conn().CopyFrom(ctx, pgx.Identifier{tableName}, columnNames, pgx.CopyFromSlice(len(records), func(i int) ([]interface{}, error) {
		return records[i], nil
	}))

	if err != nil {
		return 0, fmt.Errorf("error copying data: %w", err)
	}

	return rows, nil
}

// -=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// TRANSACTION MANAGEMENT METHODS
// -=-=-=-=-=-=-=-=-=-=-=-=-=-=-

// BeginTx starts a new transaction
// Usage: tx, err := db.BeginTx(ctx); defer tx.Rollback(ctx)
func (d *Database) BeginTx(ctx context.Context) (pgx.Tx, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.isClosed {
		return nil, fmt.Errorf("database connection is closed")
	}

	tx, err := d.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("error beginning transaction: %w", err)
	}

	return tx, nil
}

// BeginTxWithOptions starts a transaction with specific options
// Usage: tx, err := db.BeginTxWithOptions(ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
func (d *Database) BeginTxWithOptions(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if d.isClosed {
		return nil, fmt.Errorf("database connection is closed")
	}

	tx, err := d.Pool.BeginTx(ctx, txOptions)
	if err != nil {
		return nil, fmt.Errorf("error beginning transaction: %w", err)
	}

	return tx, nil
}

// ExecuteInTransaction executes a function within a transaction
// Automatically handles commit/rollback
// Usage: db.ExecuteInTransaction(ctx, func(tx pgx.Tx) error { ... })
func (d *Database) ExecuteInTransaction(ctx context.Context, fn func(tx pgx.Tx) error) error {
	tx, err := d.BeginTx(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback(ctx)
			panic(r)
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback(ctx)
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

// -=-=-=-=-=-=-=-=-=-=-=-=-=-=-
// CONNECTION MANAGEMENT METHODS
// -=-=-=-=-=-=-=-=-=-=-=-=-=-=-

// matiin koneksi database
func (d *Database) Close() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.isClosed {
		return fmt.Errorf("koneksi database nya udh ditutup")
	}

	if !d.isClosed && d.Pool != nil {
		d.Pool.Close()
		d.isClosed = true
	}

	return nil
}

// ngecek status koneksi database
func (d *Database) IsClosed() bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.isClosed
}

// dapetin waktu start koneksinya
func (d *Database) GetStartTime() time.Time {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.startTime
}

// dapetin pool koneksinya (pgxpool.Pool)
func (d *Database) GetConnPool() *pgxpool.Pool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.Pool
}

// dapetin config umum nya
func (d *Database) GetConfig() *Config {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.ConfigDB
}

// dapetin connection string nya
func (d *Database) GetConnString() string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.ConnString
}

// PING!
func (d *Database) Ping(ctx context.Context) error {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.Pool.Ping(ctx)
}

// dapetin stats koneksinya
func (d *Database) GetStats() pgxpool.Stat {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return *d.Pool.Stat()
}

// get kode dc dari config
func (d *Database) GetKodeDc() string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.ConfigDB.KodeDC
}

// / get kode dc dari tabel dc_tabel_dc_t
func (d *Database) GetKodeDcFromDB(ctx context.Context) (string, error) {
	var kodeDc string
	query := "SELECT tbl_dc_kode FROM dc_tabel_dc_t LIMIT 1"
	err := d.SelectOne(ctx, &kodeDc, query)
	if err != nil {
		return "", fmt.Errorf("error getting tbl_dc_kode from database: %w", err)
	}
	return kodeDc, nil
}

// get server db time
func (d *Database) GetServerTime(ctx context.Context) (time.Time, error) {
	var serverTime time.Time
	query := "SELECT NOW()"
	err := d.SelectOne(ctx, &serverTime, query)
	if err != nil {
		return time.Time{}, fmt.Errorf("error getting server time from database: %w", err)
	}
	return serverTime, nil
}

// dapetin uptime koneksinya
func (d *Database) GetUpTime() time.Duration {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return time.Since(d.startTime)
}
