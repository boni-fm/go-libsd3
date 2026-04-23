package settinglibgo

/*

	Library settinglib untuk SettingWeb.xml
	ini untuk kebutuhan baca xml setting web secara dinamis
	digunakan dalam api kunci

*/

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/boni-fm/go-libsd3/pkg/config/constant"
)

var (
	cachedConnInfo     Node
	cachedConnInfoTime int64
	mu                 sync.Mutex
)

// Node merepresentasikan node dalam struktur XML SettingWeb.xml.
type Node struct {
	XMLName  xml.Name
	Attrs    []xml.Attr `xml:"-"`
	Children []Node     `xml:",any"`
	Text     string     `xml:",chardata"`
}

// DynamicSettingWebXMLReader membaca nilai dari SettingWeb.xml secara dinamis berdasarkan key.
// Menggunakan cache berdasarkan waktu modifikasi file untuk menghindari pembacaan berulang.
// Mengembalikan string kosong jika file tidak dapat dibaca atau key tidak ditemukan.
// Parameter:
//   - key: nama elemen XML yang nilainya akan diambil
func DynamicSettingWebXMLReader(key string) string {
	path := getXMLSettingWebPath()

	info, err := os.Stat(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "DynamicSettingWebXMLReader: gagal stat SettingWeb.xml: %v\n", err)
		return ""
	}

	mu.Lock()
	defer mu.Unlock()

	if cachedConnInfoTime == info.ModTime().Unix() {
		for _, child := range cachedConnInfo.Children {
			if strings.EqualFold(child.XMLName.Local, key) {
				return child.Text
			}
		}
		return ""
	}

	xmlFile, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "DynamicSettingWebXMLReader: gagal membuka SettingWeb.xml: %v\n", err)
		return ""
	}
	defer xmlFile.Close()

	byteValue, err := io.ReadAll(xmlFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "DynamicSettingWebXMLReader: gagal membaca SettingWeb.xml: %v\n", err)
		return ""
	}

	var xmlNode Node
	xml.Unmarshal(byteValue, &xmlNode)

	// Update cache BEFORE searching
	cachedConnInfo = xmlNode
	cachedConnInfoTime = info.ModTime().Unix()

	for _, child := range xmlNode.Children {
		if strings.EqualFold(child.XMLName.Local, key) {
			return child.Text
		}
	}

	return ""
}

func getXMLSettingWebPath() string {
	if osName := runtime.GOOS; osName == "windows" {
		return constant.FILEPATH_SETTINGWEB_WINDOWS
	}
	return constant.FILEPATH_SETTINGWEB_LINUX
}
