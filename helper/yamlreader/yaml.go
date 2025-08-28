package yamlreader

import (
	"errors"
	"os"

	"gopkg.in/yaml.v2"
)

var ErrorNullConfigValue = errors.New("config.yaml tidak memiliki value data")
var ErrorNullKeyValue = errors.New("value dari key tidak ditemukan dalam config.yaml")

func ReadConfigDynamic(path string) (map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg map[string]interface{}
	_ = yaml.Unmarshal(data, &cfg)
	if cfg == nil {
		return nil, ErrorNullConfigValue
	}
	return cfg, nil
}

func ReadConfigDynamicWithKey(path string, key string) (interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg map[string]interface{}
	_ = yaml.Unmarshal(data, &cfg)
	if cfg[key] == nil {
		return nil, ErrorNullKeyValue
	}
	return cfg[key], nil
}
