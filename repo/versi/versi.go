package versi

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

var msgTidakTerdaftar = "Program .:%s:. belum terdaftar di Master Program DC,\r\n              Segera Hubungi \r\n        --=::>> SUPPORT <<::=-- "

func GetVersiProgramPostgre(Constr, Kodedc, NamaProgram, Versi, IPKomputer string) string {
	text := strings.ReplaceAll(IPKomputer, "'", "")
	Constr = PostgreConstrBuilder(Constr)
	result := ""

	db, err := sql.Open("postgres", Constr)
	if err != nil {
		return "koneksi db gagal!"
	}
	defer db.Close()

	queryVersi :=
		`SELECT
			CASE  WHEN coalesce(APROVE, 'N') = 'Y' AND coalesce(tgl_berlaku,current_date)<=current_date 
			THEN  coalesce(VERSI_BARU, '0')  
			WHEN coalesce(APROVE,'N')='N' THEN coalesce(VERSI_LAMA,'0') 
			ELSE coalesce(VERSI_LAMA ,'0') END AS VERSI  
		FROM dc_program_vbdtl_t 
		WHERE Dc_KODE=$1 and UPPER(Nama_Prog)=$2`
	err = db.QueryRow(queryVersi, Kodedc, strings.ToUpper(NamaProgram)).Scan(&result)
	if err != nil {
		return fmt.Sprintf(msgTidakTerdaftar, strings.ToUpper(NamaProgram))
	}

	versiDB, err := strconv.Atoi(strings.ReplaceAll(strings.TrimSpace(result), ".", ""))
	if err != nil {
		return fmt.Sprintf(msgTidakTerdaftar, strings.ToUpper(NamaProgram))
	}
	versi, err := strconv.Atoi(strings.ReplaceAll(strings.TrimSpace(Versi), ".", ""))
	if err != nil {
		return fmt.Sprintf(msgTidakTerdaftar, strings.ToUpper(NamaProgram))
	}

	if versiDB != versi {
		if versi <= versiDB {
			return fmt.Sprintf("    Program .:%s:. belum update,\r\n      Versi update \r\n--==::>> %s <<::==--", strings.ToUpper(NamaProgram), result)
		} else {
			return fmt.Sprintf("    Program .:%s:. Versi program tidak sama dengan master ,\r\n      Versi Master \r\n--==::>> %s <<::==--", strings.ToUpper(NamaProgram), result)
		}
	} else {
		result = "OKE..."
		queryInsMonitor := `
			INSERT INTO dc_monitoring_program_t (Kode_Dc,Nama_Program,Ip_CLient,versi,Tanggal)
			VALUES ($1, $2, $3, $4, current_timestamp)
		`
		_, err := db.Exec(queryInsMonitor, Kodedc, strings.ToUpper(NamaProgram), text, versi)
		if err != nil {
			return fmt.Sprintf("    Program .:%s:. Gagal mencatat aktivitas monitor.\r\n", strings.ToUpper(NamaProgram))
		}
	}

	return result
}

func PostgreConstrBuilder(constr string) string {
	parts := strings.Split(constr, ";")
	var (
		server, database, userID, password, port string
	)
	for _, part := range parts {
		if strings.HasPrefix(part, "Server=") {
			server = strings.TrimPrefix(part, "Server=")
		} else if strings.HasPrefix(part, "Database=") {
			database = strings.TrimPrefix(part, "Database=")
		} else if strings.HasPrefix(part, "User Id=") {
			userID = strings.TrimPrefix(part, "User Id=")
		} else if strings.HasPrefix(part, "Password=") {
			password = strings.TrimPrefix(part, "Password=")
		} else if strings.HasPrefix(part, "Port=") {
			port = strings.TrimPrefix(part, "Port=")
		}
	}
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", server, port, userID, password, database)
}
