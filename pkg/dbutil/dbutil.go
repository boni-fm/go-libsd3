package dbutil

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/boni-fm/go-libsd3/helper/kunci"
	"github.com/boni-fm/go-libsd3/helper/logging"
	"github.com/boni-fm/go-libsd3/helper/yamlreader"

	_ "github.com/lib/pq"
)

/*
	TODO:
	- buat class nya jadi singleton
*/

type IDatabase interface {
	Connect() (*PostgreDB, error)
	Close() bool
	HealthCheck() string
	Select(query string, args ...interface{}) (*sql.Rows, error)
	SelectScalar(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
	Call(proc string, args ...interface{}) (*sql.Row, error)
}

type PostgreDB struct {
	db *sql.DB
}

var oncePG sync.Once

var log = logging.NewLogger()

const POSTGRE_DBTYPE = "POSTGRE"

var pgInstance *PostgreDB

type DatabaseSetup struct {
	kunciManager kunci.Kunci
}

func SetupConnectionDatabase() (*PostgreDB, error) {

	var err error
	kuncipath, _ := yamlreader.GetKunciConfigFilepath()
	strkunci, _ := yamlreader.ReadConfigDynamicWithKey(kuncipath, "kunci")
	databaseSetup := DatabaseSetup{
		kunciManager: *kunci.NewKunci(strkunci.(string)),
	}

	// Initialize database connection
	if pgInstance, err = databaseSetup.Connect(); err != nil {
		log.SayError("Failed to connect to database: " + err.Error())
		return nil, err
	}
	return pgInstance, nil
}

// buat initialize connection, make singleton
func (d *DatabaseSetup) Connect() (*PostgreDB, error) {
	if pgInstance == nil {
		oncePG.Do(func() {
			connString := d.kunciManager.GetConnectionString(POSTGRE_DBTYPE)
			DB, err := sql.Open("postgres", connString)
			if err != nil {
				log.SayError("Gagal connect ke database: " + err.Error())
				return
			}
			pgInstance = &PostgreDB{db: DB}
		})
		return pgInstance, nil
	}

	return pgInstance, nil
}

func (p *PostgreDB) Close() error {
	if p.db != nil {
		log.Say("Pintu koneksi database ditutup ~")
		pgInstance = nil
		return p.db.Close()
	}

	return nil
}

func (p *PostgreDB) HealthCheck() string {
	if p.db != nil {
		err := p.db.Ping()
		if err != nil {
			msgSakit := "Koneksi DB gk sehat kawan ~ " + err.Error()
			log.SayFatal(msgSakit)
			return msgSakit
		}

		msgSehat := "Koneksi DB sehat walafiat ~ "
		log.Say(msgSehat)
		return msgSehat
	}
	return "Loh kok databasenya kosong?? ~ "
}

func (p *PostgreDB) Select(query string, args ...interface{}) (*sql.Rows, error) {
	return p.db.Query(query, args...)
}

func (p *PostgreDB) SelectScalar(query string, args ...interface{}) *sql.Row {
	return p.db.QueryRow(query, args...)
}

func (p *PostgreDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return p.db.Exec(query, args...)
}

func (p *PostgreDB) Call(proc string, args ...interface{}) (*sql.Rows, error) {
	q := "SELECT * FROM " + proc + "("
	for i := range args {
		if i > 0 {
			q += ","
		}
		q += fmt.Sprintf("$%d", i+1)
	}
	q += ")"

	return p.db.Query(q, args...)
}

func GetDB() *sql.DB {
	return pgInstance.db
}
