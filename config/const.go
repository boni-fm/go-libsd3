package config

import "time"

const (
	FILEPATH_CONFIG_WINDOWS = "D:/go-libsd3/config.yaml"
	FILEPATH_CONFIG_LINUX   = "/etc/go-libsd3/config.yaml"
	KEY_ENV_KUNCI           = "KunciWeb"

	FILEPATH_LOG_LINUX   = "/var/log/nginx/api"
	FILEPATH_LOG_WINDOWS = "D:/_docker/_app/logs"

	FILEPATH_SETTINGWEB_WINDOWS = "D:\\_docker\\_app\\kunci\\SettingWeb.xml"
	FILEPATH_SETTINGWEB_LINUX   = "SettingWeb.xml"

	DATE_FORMAT     = "2006-01-02"
	DATETIME_FORMAT = "2006-01-02 15:04:05"

	DBTYPE_POSTGRE = "POSTGRE"

	TIME_TWO_MINUTES = 120 * time.Second
	TIME_ONE_MINUTE  = 60 * time.Second
	TIME_TEN_SECONDS = 10 * time.Second
)
