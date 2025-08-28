package yamlreader

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

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

func GetKunciConfigFilepath() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			dir = filepath.Join(dir, "config.yaml")
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", fmt.Errorf("project root not found")
}
