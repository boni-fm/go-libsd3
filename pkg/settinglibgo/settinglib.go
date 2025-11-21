package settinglibgo

import (
	"os"

	"github.com/boni-fm/go-libsd3/pkg/config/constant"
	"github.com/boni-fm/go-libsd3/pkg/yaml"
)

func GetConnStringPostgre() string {
	const dbtype = "POSTGRE"
	var strKunci string
	strKunciDocker := os.Getenv(constant.KEY_ENV_KUNCI)

	if strKunciDocker != "" {
		strKunci = strKunciDocker
	} else {
		kuncipath, _ := yaml.GetKunciConfigFilepath()
		key, _ := yaml.ReadConfigDynamicWithKey(kuncipath, "kunci")
		strKunci = key.(string)
	}

	key := NewSettingLib(strKunci)
	ConStr := key.GetConnectionString(dbtype)
	return ConStr
}
