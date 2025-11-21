package settinglibgo

/*
	KUNCI!!!
	- untuk dapetin connection string

	Note:
	- untuk manggil constring lewat function ini, perlu buat clientnya dulu
	  > lib = settinglib.NewSettingLib("{kode dc}")
      > constring = lib.GetConnectionString("POSTGRE")
    - kalo mau lgsg manggil bisa make default function nya aja

	TODO:
	- bikin versi cache nya untuk ngurangin resource
	- buat lebih di
*/

import (
	"fmt"
	"strings"

	logger "github.com/boni-fm/go-libsd3/pkg/log"
)

var log = logger.NewLoggerWithFilename("settinglib")

type Config[T PostgreConnectionConfig] struct {
	ConnectionConfig T
}

type PostgreConnectionConfig struct {
	IPPostgres       string
	PortPostgres     string
	DatabasePostgres string
	UserPostgres     string
	PasswordPostgres string
}

type DBConfig struct {
	PostgreConnectionConfig
}

type MainConfiguration struct {
	PostgreConnectionConfig Config[PostgreConnectionConfig]
	SettingLibClient        *SettingLibClient
}

func NewSettingLib(kuncidc string) *MainConfiguration {
	return &MainConfiguration{
		SettingLibClient: NewSettingLibClient(kuncidc),
	}
}

// Get constring nya :D
func (k *MainConfiguration) GetConnectionString(dbtype string) string {
	switch strings.ToUpper(dbtype) {
	case "POSTGRE":
		k.SetPGConStringFromWebservice()
		config := k.PostgreConnectionConfig
		return fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			config.ConnectionConfig.IPPostgres, config.ConnectionConfig.PortPostgres, config.ConnectionConfig.UserPostgres, config.ConnectionConfig.PasswordPostgres, config.ConnectionConfig.DatabasePostgres)
	}

	// return kosong klo gk ada apa2 wkwkwkwk
	return ""
}

// Koleksi dapetin constring dari webservice kunci
func (k *MainConfiguration) SetPGConStringFromWebservice() (*MainConfiguration, error) {
	user, err := k.SettingLibClient.GetVariable("UserPostgres")
	if err != nil {
		log.Error("Failed to get UserPostgres variable :", err)
		return nil, err
	}

	password, err := k.SettingLibClient.GetVariable("PasswordPostgres")
	if err != nil {
		log.Error("Failed to get PasswordPostgres variable :", err)
		return nil, err
	}

	port, err := k.SettingLibClient.GetVariable("PortPostgres")
	if err != nil {
		log.Error("Failed to get PortPostgres variable :", err)
		return nil, err
	}
	ip, err := k.SettingLibClient.GetVariable("IPPostgres")
	if err != nil {
		log.Error("Failed to get IPPostgres variable :", err)
		return nil, err
	}
	database, err := k.SettingLibClient.GetVariable("DatabasePostgres")
	if err != nil {
		log.Error("Failed to get DatabasePostgres variable :", err)
		return nil, err
	}

	k.PostgreConnectionConfig = Config[PostgreConnectionConfig]{}
	k.PostgreConnectionConfig.ConnectionConfig = PostgreConnectionConfig{
		IPPostgres:       ip,
		PortPostgres:     port,
		DatabasePostgres: database,
		UserPostgres:     user,
		PasswordPostgres: password,
	}

	return k, nil
}
