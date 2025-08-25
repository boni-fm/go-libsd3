package kunci

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/boni-fm/go-libsd3/helper/logging"
)

/*
	TODO:
	- Implementasi ambil sql constring dari settingweb.xml
	- buat supaya bisa baca settingwebgxxx.xml, sekarang masih bisa 1 doang
	- sesuain kalo udh bisa baca gxxx ke dalam dbutil

	TODO: Difficulty >> HARD!!
	- set kunci dari GetVariable mirip dengan SettingLibb punya .net pak edwin
*/

var log = logging.NewLogger()

type SettingWeb[T ConnectionStringPostgre | ConnectionStringSQL] struct {
	XMLName          xml.Name `xml:"SettingConfig"`
	ConnectionString T
}

type ConnectionStringPostgre struct {
	IPPostgres       string `xml:"IPPostgres"`
	PortPostgres     string `xml:"PortPostgres"`
	DatabasePostgres string `xml:"DatabasePostgres"`
	UserPostgres     string `xml:"UserPostgres"`
	PasswordPostgres string `xml:"PasswordPostgres"`
}

type ConnectionStringSQL struct {
	IPSql       string `xml:"IPSql"`
	UserSql     string `xml:"UserSql"`
	PasswordSql string `xml:"PasswordSql"`
	DatabaseSql string `xml:"DatabaseSql"`
}

func GetConnectionInfoPostgre() ConnectionStringPostgre {
	homepath, _ := os.UserHomeDir()
	pathkunci := filepath.Join(homepath, "_docker", "_app", "kunci", "SettingWeb.xml")
	xmlFile, err := os.Open(pathkunci)
	if err != nil {
		log.SayFatalf("Failed to open SettingWeb.xml: %v", err)
	}

	defer xmlFile.Close()

	byteValue, _ := io.ReadAll(xmlFile)
	var connInfo ConnectionStringPostgre
	xml.Unmarshal(byteValue, &connInfo)

	if strings.Contains(connInfo.UserPostgres, "Timeout") {
		connInfo.UserPostgres = strings.Split(connInfo.UserPostgres, ";")[0]
	}

	return connInfo
}

func GetConnectionString(dbtype string) string {
	switch strings.ToUpper(dbtype) {
	case "POSTGRE":
		pgConnInfo := GetConnectionInfoPostgre()
		return fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			pgConnInfo.IPPostgres, pgConnInfo.PortPostgres, pgConnInfo.UserPostgres, pgConnInfo.PasswordPostgres, pgConnInfo.DatabasePostgres)
	}

	return ""
}
