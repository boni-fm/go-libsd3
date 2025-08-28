package kunci

import (
	"os"

	"gopkg.in/yaml.v2"
)

type ConfigKunci struct {
	Kunci string `yaml:"kunci"`
}

func ReadConfig(configPath string) (string, error) {
	//wd, err := os.Getwd()
	// if err != nil {
	// 	return "", err
	// }
	//configPath := filepath.Join(wd, "config", "config.yaml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return "", err
	}
	var cfg ConfigKunci
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return "", err
	}
	return cfg.Kunci, nil
}
