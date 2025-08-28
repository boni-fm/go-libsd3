package kunci

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"runtime"
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

type Config[T PostgreConnectionConfig] struct {
	ConnectionConfig T
}

type PostgreConnectionConfig struct {
	IPPostgres       string `xml:"IPPostgres"`
	PortPostgres     string `xml:"PortPostgres"`
	DatabasePostgres string `xml:"DatabasePostgres"`
	UserPostgres     string `xml:"UserPostgres"`
	PasswordPostgres string `xml:"PasswordPostgres"`
}

type Kunci struct {
	PostgreConfig Config[PostgreConnectionConfig]
	KunciClient   *KunciClient
}

func NewKunci(kuncidc string) *Kunci {
	return &Kunci{
		KunciClient: NewKunciClient(kuncidc),
	}
}

// ini fungsi kalo mau baca langsung dari settingweb.xml
// untuk sekarang tidak digunakan
func GetConnectionInfoPostgre() PostgreConnectionConfig {
	settingWebFile := func() (*os.File, error) {
		if osName := runtime.GOOS; osName == "windows" {
			return os.Open(`D:\_docker\_app\kunci\SettingWeb.xml`)
		} else {
			return os.Open("/_docker/_app/kunci/SettingWeb.xml")
		}
	}

	xmlFile, err := settingWebFile()
	if err != nil {
		log.SayFatalf("Failed to open SettingWeb.xml: %v", err)
	}
	defer xmlFile.Close()

	byteValue, _ := io.ReadAll(xmlFile)
	var connInfo PostgreConnectionConfig
	xml.Unmarshal(byteValue, &connInfo)

	if strings.Contains(connInfo.UserPostgres, "Timeout") {
		connInfo.UserPostgres = strings.Split(connInfo.UserPostgres, ";")[0]
	}

	return connInfo
}

func (k *Kunci) GetConnectionString(dbtype string) string {
	switch strings.ToUpper(dbtype) {
	case "POSTGRE":
		k.SetPGConStringFromWebservice()
		config := k.PostgreConfig
		return fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			config.ConnectionConfig.IPPostgres, config.ConnectionConfig.PortPostgres, config.ConnectionConfig.UserPostgres, config.ConnectionConfig.PasswordPostgres, config.ConnectionConfig.DatabasePostgres)
	}

	return ""
}

func (k *Kunci) SetPGConStringFromWebservice() (*Kunci, error) {
	user, err := k.KunciClient.GetVariable("UserPostgres")
	if err != nil {
		log.Error("Failed to get UserPostgres variable :", err)
		return nil, err
	}

	password, err := k.KunciClient.GetVariable("PasswordPostgres")
	if err != nil {
		log.Error("Failed to get PasswordPostgres variable :", err)
		return nil, err
	}

	port, err := k.KunciClient.GetVariable("PortPostgres")
	if err != nil {
		log.Error("Failed to get PortPostgres variable :", err)
		return nil, err
	}
	ip, err := k.KunciClient.GetVariable("IPPostgres")
	if err != nil {
		log.Error("Failed to get IPPostgres variable :", err)
		return nil, err
	}
	database, err := k.KunciClient.GetVariable("DatabasePostgres")
	if err != nil {
		log.Error("Failed to get DatabasePostgres variable :", err)
		return nil, err
	}

	k.PostgreConfig = Config[PostgreConnectionConfig]{}
	k.PostgreConfig.ConnectionConfig = PostgreConnectionConfig{
		IPPostgres:       ip,
		PortPostgres:     port,
		DatabasePostgres: database,
		UserPostgres:     user,
		PasswordPostgres: password,
	}

	return k, nil
}
