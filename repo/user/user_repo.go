package general

import (
	"fmt"

	"github.com/boni-fm/go-libsd3/helper/logging"
	"github.com/boni-fm/go-libsd3/pkg/dbutil"
)

/*
	TODO:
	- kalo lifetime db nya udh aman, bisa di pake db.closenya supaya gk menuhin koneksi
*/

var log = logging.NewLogger()

func GetUserPassword(username string) {
	db, err := dbutil.Connect()
	if err != nil {
		return
	}

	rows, err := db.Select("Select user_password from dc_user_t where user_name = $1", username)
	if err != nil {
		log.SayError(err.Error())
		return
	}
	defer rows.Close()

	var pass string
	if rows.Next() {
		if err := rows.Scan(&pass); err != nil {
			log.SayError(err.Error())
			return
		}
	}
	fmt.Println(pass)
}
