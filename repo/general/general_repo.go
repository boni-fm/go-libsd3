package general

import (
	"fmt"

	"github.com/boni-fm/go-libsd3/helper/logging"
	"github.com/boni-fm/go-libsd3/pkg/dbutil"
)

func GetUser(username string) {
	log := logging.NewLogger()
	db, err := dbutil.Connect()
	if err != nil {
		return
	}
	defer db.Close()

	// Implement your user retrieval logic here
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
