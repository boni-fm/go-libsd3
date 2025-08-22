package dbutil

import (
	"database/sql"
	"fmt"
	"go-libsd3/helper/kunci"
	"go-libsd3/helper/logging"

	_ "github.com/lib/pq"
)

type IDatabase interface {
	Connect() bool
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

var log = logging.NewLogger()

const POSTGRE_DBTYPE = "POSTGRE"

var DBPostgre *sql.DB

// fungsi yang dijalanin paling pertama sebelum exec query
func Connect() (*PostgreDB, error) {
	connString := kunci.GetConnectionString(POSTGRE_DBTYPE)

	DB, err := sql.Open("postgres", connString)
	if err != nil {
		log.SayError(err.Error())
		return nil, err
	}

	return &PostgreDB{db: DB}, nil
}

func (p *PostgreDB) Close() error {
	if p.db != nil {
		log.Say("Closing database connection ~")
		return p.db.Close()
	}
	return nil
}

func (p *PostgreDB) HealthCheck() string {
	if p.db != nil {
		err := p.db.Ping()
		if err != nil {
			log.SayFatal("Koneksi DB gk sehat kawan ~")
			return "Koneksi DB gk sehat kawan ~"
		}

		log.Say("Koneksi DB sehat walafiat ~")
		return "Koneksi DB sehat walafiat ~"
	}
	return "Koneksi DB gk sehat kawan ~"
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
