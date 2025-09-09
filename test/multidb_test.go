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

	// Close all connections
	multidb.CloseAllConnection()
}
