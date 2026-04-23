package settinglibgo

import (
	"os"

	"github.com/boni-fm/go-libsd3/pkg/config/constant"
	"github.com/boni-fm/go-libsd3/pkg/yaml"
)

// GetConnStringPostgre mengembalikan connection string PostgreSQL.
// Mengambil kunci dari environment variable atau file konfigurasi YAML.
// Mengembalikan string kosong jika kunci tidak ditemukan atau bukan bertipe string.
func GetConnStringPostgre() string {
	const dbtype = "POSTGRE"
	var strKunci string
	strKunciDocker := os.Getenv(constant.KEY_ENV_KUNCI)

	if strKunciDocker != "" {
		strKunci = strKunciDocker
	} else {
		kuncipath, _ := yaml.GetKunciConfigFilepath()
		key, _ := yaml.ReadConfigDynamicWithKey(kuncipath, "kunci")
		if v, ok := key.(string); ok {
			strKunci = v
		}
	}

	if strKunci == "" {
		return ""
	}

	k := NewSettingLib(strKunci)
	ConStr := k.GetConnectionString(dbtype)
	return ConStr
}
