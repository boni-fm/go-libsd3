package settinglibgooo

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/boni-fm/go-libsd3/config"
	"github.com/boni-fm/go-libsd3/helper/logging"
)

/*
	TODO:
	- sesuain kalo udh bisa baca gxxx ke dalam dbutil
	- buat dokumentasi
*/

var log = logging.NewLoggerWithFilename("db-setup")

var (
	cachedConnInfo     Node
	cachedConnInfoTime int64
	mu                 sync.Mutex
)

type Config[T PostgreConnectionConfig | SqlConnectionConfig] struct {
	ConnectionConfig T
}

type PostgreConnectionConfig struct {
	IPPostgres       string `xml:"IPPostgres"`
	PortPostgres     string `xml:"PortPostgres"`
	DatabasePostgres string `xml:"DatabasePostgres"`
	UserPostgres     string `xml:"UserPostgres"`
	PasswordPostgres string `xml:"PasswordPostgres"`
}

type SqlConnectionConfig struct {
	IPSql       string `xml:"IPSql"`
	UserSql     string `xml:"UserSql"`
	PasswordSql string `xml:"PasswordSql"`
	DatabaseSql string `xml:"DatabaseSql"`
}

type OthersConfig struct {
	Antrian      string `xml:"Antrian"`
	Com          string `xml:"Com"`
	JenisPrinter string `xml:"JenisPrinter"`
}

type DBConfig struct {
	PostgreConnectionConfig
	SqlConnectionConfig
	OthersConfig
}

type Kunci struct {
	PostgreConfig    Config[PostgreConnectionConfig]
	SettingWebClient *SettingLibClient
}

func NewSettingLib(kuncidc string) *Kunci {
	return &Kunci{
		SettingWebClient: NewSettingLibClient(kuncidc),
	}
}

type Node struct {
	XMLName  xml.Name
	Attrs    []xml.Attr `xml:"-"`
	Children []Node     `xml:",any"`
	Text     string     `xml:",chardata"`
}

func DynamicSettingWebXMLReader(key string) string {
	settingWebPath := func() string {
		if osName := runtime.GOOS; osName == "windows" {
			return config.FILEPATH_SETTINGWEB_WINDOWS
		}
		return config.FILEPATH_SETTINGWEB_LINUX
	}

	path := settingWebPath()
	info, err := os.Stat(path)
	if err != nil {
		log.SayFatalf("Failed to stat SettingWeb.xml: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()
	if cachedConnInfoTime == info.ModTime().Unix() {
		for _, child := range cachedConnInfo.Children {
			if strings.EqualFold(child.XMLName.Local, key) {
				return child.Text
			}
		}
	}

	xmlFile, err := os.Open(path)
	if err != nil {
		log.SayFatalf("Failed to open SettingWeb.xml: %v", err)
	}
	defer xmlFile.Close()

	byteValue, _ := io.ReadAll(xmlFile)
	var xmlNode Node
	xml.Unmarshal(byteValue, &xmlNode)

	for _, child := range xmlNode.Children {
		if strings.EqualFold(strings.ToLower(child.XMLName.Local), strings.ToLower(key)) {
			return child.Text
		}
	}

	cachedConnInfo = xmlNode
	cachedConnInfoTime = info.ModTime().Unix()

	return ""
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
	user, err := k.SettingWebClient.GetVariable("UserPostgres")
	if err != nil {
		log.Error("Failed to get UserPostgres variable :", err)
		return nil, err
	}

	password, err := k.SettingWebClient.GetVariable("PasswordPostgres")
	if err != nil {
		log.Error("Failed to get PasswordPostgres variable :", err)
		return nil, err
	}

	port, err := k.SettingWebClient.GetVariable("PortPostgres")
	if err != nil {
		log.Error("Failed to get PortPostgres variable :", err)
		return nil, err
	}
	ip, err := k.SettingWebClient.GetVariable("IPPostgres")
	if err != nil {
		log.Error("Failed to get IPPostgres variable :", err)
		return nil, err
	}
	database, err := k.SettingWebClient.GetVariable("DatabasePostgres")
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
