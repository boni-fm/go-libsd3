package config

import "time"

const (
	FILEPATH_CONFIG_WINDOWS = "D:/go-libsd3/config.yaml"
	FILEPATH_CONFIG_LINUX   = "/etc/go-libsd3/config.yaml"
	KEY_ENV_KUNCI           = "KunciWeb"

	FILEPATH_LOG_LINUX   = "/var/log/nginx/api"
	FILEPATH_LOG_WINDOWS = "D:/_docker/_app/logs"

	FILEPATH_SETTINGWEB_WINDOWS = "D:\\_docker\\_app\\kunci\\SettingWeb.xml"
	FILEPATH_SETTINGWEB_LINUX   = "/_docker/_app/kunci/SettingWeb.xml"

	DATE_FORMAT     = "2006-01-02"
	DATETIME_FORMAT = "2006-01-02 15:04:05"

	DBTYPE_POSTGRE = "POSTGRE"

	TIMEOUT_TWO_MINUTES = 120 * time.Second
)
