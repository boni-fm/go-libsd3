package general

import (
	"fmt"

	"github.com/boni-fm/go-libsd3/pkg/dbutil"
)

/*
	TODO:
	- kalo lifetime db nya udh aman, bisa di pake db.closenya supaya gk menuhin koneksi
*/

func GetKodeDC() string {
	db, err := dbutil.SetupConnectionDatabase()
	if err != nil {
		return ""
	}

	// jangan pake ini karena lifetime db nya itu singleton, nanti yang lain kena impact
	// ini dipake kalau bisa connect dari multiple db
	// defer db.Close()

	var kodeDC string
	row := db.SelectScalar("SELECT TBL_DC_KODE FROM DC_TABEL_DC_T")
	if row != nil {
		kodeDC = fmt.Sprintf("%v", row)
	}
	return kodeDC
}

func GetDate() string {
	db, err := dbutil.SetupConnectionDatabase()
	if err != nil {
		return ""
	}

	// jangan pake ini karena lifetime db nya itu singleton, nanti yang lain kena impact
	// ini dipake kalau bisa connect dari multiple db
	// defer db.Close()

	var currentDate string
	row := db.SelectScalar("SELECT CURRENT_DATE")
	if row != nil {
		currentDate = fmt.Sprintf("%v", row)
	}
	return currentDate
}
