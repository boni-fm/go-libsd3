package multidbutil

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"

	"github.com/boni-fm/go-libsd3/helper/settinglibgooo"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

// Multiton instance buat db
// src :
// - https://freedium.cfd/https://levelup.gitconnected.com/multiton-design-pattern-in-golang-with-unit-tests-33f194a3fab5
// - https://www.hackingwithgo.nl/2023/10/11/lazy-initialization-and-multiton-a-cheap-way-of-creating-expensive-objects/

var (
	PREFIX_KUNCI = "kunci"
)

type DBConfig struct {
	Client     *settinglibgooo.Kunci
	DB         *sql.DB
	Log        *logrus.Logger
	ConnString string
}

type MultiDB struct {
	KodeDC  string
	Configs map[string]*DBConfig
	mu      sync.Mutex
}

func (m *MultiDB) SetupMultiDB(kodedc string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.Configs == nil {
		m.Configs = make(map[string]*DBConfig)
	}

	if m.Configs[kodedc] == nil {
		client := settinglibgooo.NewSettingLib(PREFIX_KUNCI + strings.ToLower(kodedc))
		connStr := client.GetConnectionString("POSTGRE")
		log := logrus.New()
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			log.Errorf("Error connecting to database for %s: %v", kodedc, err)
			return err
		}

		// di ping buat mastiin koneksinya lagi
		if err = db.Ping(); err != nil {
			log.Errorf("Error pinging database for %s: %v", kodedc, err)
			return err
		}

		m.Configs[kodedc] = &DBConfig{
			Client:     client,
			DB:         db,
			Log:        log,
			ConnString: connStr,
		}

		return nil
	}
	return nil
}

func (m *MultiDB) GetDB(kodedc string) (*sql.DB, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.Configs[kodedc] == nil {
		err := m.SetupMultiDB(kodedc)
		if err != nil {
			return nil, err
		}
	}

	if m.Configs[kodedc] != nil {
		return m.Configs[kodedc].DB, nil
	}

	return nil, fmt.Errorf("terdapat kesalahan saat mencoba mengakses database %s", kodedc)
}

func (m *MultiDB) CloseAllConnection() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for kodedc, config := range m.Configs {
		if config.DB != nil {
			if err := config.DB.Close(); err != nil {
				config.Log.Errorf("Error closing database for %s: %v", kodedc, err)
			} else {
				config.Log.Infof("Database connection for %s closed successfully", kodedc)
			}
		}
	}
}

func (m *MultiDB) SelectScalarByKodedc(kodedc, query string, args ...interface{}) (result interface{}, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.Configs[kodedc] == nil {
		return nil, fmt.Errorf("database configuration for %s not found", kodedc)
	}

	err = m.Configs[kodedc].DB.QueryRow(query, args...).Scan(&result)
	if err != nil {
		m.Configs[kodedc].Log.Errorf("Error executing query for %s: %v", kodedc, err)
		return nil, err
	}

	return result, nil
}
