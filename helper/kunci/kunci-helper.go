package kunci

import (
	"os"

	"github.com/boni-fm/go-libsd3/config"
	"github.com/boni-fm/go-libsd3/helper/yamlreader"
)

func GetConnStringDockerPostgre(tipe string) string {
	var strKunci string
	strKunciDocker := os.Getenv(config.KeyEnvKunci)

	if strKunciDocker != "" {
		strKunci = strKunciDocker
	} else {
		kuncipath, _ := yamlreader.GetKunciConfigFilepath()
		key, _ := yamlreader.ReadConfigDynamicWithKey(kuncipath, "kunci")
		strKunci = key.(string)
	}

	key := NewKunci(strKunci)
	ConStr := key.GetConnectionString(tipe)
	return ConStr
}
