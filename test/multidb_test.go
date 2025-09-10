package test

import (
	"testing"

	"github.com/boni-fm/go-libsd3/pkg/multidbutil"
)

var multidb multidbutil.MultiDB

func TestMultiDB_SetupAndGetDB(t *testing.T) {
	kodedc := "G009sim"

	// Setup DB
	multidb.SetupMultiDB(kodedc)
	db, err := multidb.GetDB(kodedc)
	if err != nil {
		t.Fatalf("GetDB failed: %v", err)
	}
	if db == nil {
		t.Fatalf("DB instance is nil")
	}
	if db.Ping() != nil {
		t.Fatalf("DB ping failed: %v", err)
	}

	row, _ := multidb.SelectScalarByKodedc(kodedc, "select TBL_DC_KODE from dc_tabel_dc_t")
	t.Logf("====> Hasil TBL_DC_KODE: %s", row.(string))

	kodedc2 := "g217SIM"

	// Setup DB
	multidb.SetupMultiDB(kodedc2)
	db2, err := multidb.GetDB(kodedc2)
	if err != nil {
		t.Fatalf("GetDB failed: %v", err)
	}
	if db2 == nil {
		t.Fatalf("DB instance is nil")
	}
	if db2.Ping() != nil {
		t.Fatalf("DB ping failed: %v", err)
	}

	row2, _ := multidb.SelectScalarByKodedc(kodedc2, "select TBL_DC_KODE from dc_tabel_dc_t")
	t.Logf("====> Hasil TBL_DC_KODE: %s", row2.(string))

	for k, v := range multidb.Configs {
		t.Logf("====> Koneksi DB KodeDC: %s, Config: %+v", k, v)
	}

	// Close all connections
	multidb.CloseAllConnection()
}
