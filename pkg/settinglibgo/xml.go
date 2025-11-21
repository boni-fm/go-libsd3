package settinglibgo

/*

	Library settinglib untuk SettingWeb.xml
	ini untuk kebutuhan baca xml setting web secara dinamis
	digunakan dalam api kunci

*/

import (
	"encoding/xml"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/boni-fm/go-libsd3/pkg/config/constant"
	"github.com/sirupsen/logrus"
)

var (
	cachedConnInfo     Node
	cachedConnInfoTime int64
	mu                 sync.Mutex
)

type Node struct {
	XMLName  xml.Name
	Attrs    []xml.Attr `xml:"-"`
	Children []Node     `xml:",any"`
	Text     string     `xml:",chardata"`
}

func DynamicSettingWebXMLReader(key string) string {
	log := logrus.New()
	path := getXMLSettingWebPath()

	info, err := os.Stat(path)
	if err != nil {
		log.Fatalf("Failed to stat SettingWeb.xml: %v", err)
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
		log.Fatalf("Failed to open SettingWeb.xml: %v", err)
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

func getXMLSettingWebPath() string {
	if osName := runtime.GOOS; osName == "windows" {
		return constant.FILEPATH_SETTINGWEB_WINDOWS
	}
	return constant.FILEPATH_SETTINGWEB_LINUX
}
