package dbutil

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/boni-fm/go-libsd3/helper/kunci"
	"github.com/boni-fm/go-libsd3/helper/logging"

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

// buat initialize connection, make singleton
func Connect() (*PostgreDB, error) {
	if pgInstance == nil {
		oncePG.Do(func() {
			connString := kunci.GetConnectionString(POSTGRE_DBTYPE)
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
	err := p.db.Ping()
	if err != nil {
		Connect()
	}

	return p.db.Query(query, args...)
}

func (p *PostgreDB) SelectScalar(query string, args ...interface{}) *sql.Row {
	err := p.db.Ping()
	if err != nil {
		Connect()
	}

	return p.db.QueryRow(query, args...)
}

func (p *PostgreDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	err := p.db.Ping()
	if err != nil {
		Connect()
	}

	return p.db.Exec(query, args...)
}

func (p *PostgreDB) Call(proc string, args ...interface{}) (*sql.Rows, error) {
	err := p.db.Ping()
	if err != nil {
		Connect()
	}

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
