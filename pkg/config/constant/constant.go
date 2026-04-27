package constant

import "time"

const (
	// FILEPATH_CONFIG_WINDOWS adalah path file konfigurasi di sistem operasi Windows.
	FILEPATH_CONFIG_WINDOWS = "D:/go-libsd3/config.yaml"
	// FILEPATH_CONFIG_LINUX adalah path file konfigurasi di sistem operasi Linux.
	FILEPATH_CONFIG_LINUX = "/etc/go-libsd3/config.yaml"
	// KEY_ENV_KUNCI adalah nama environment variable untuk menyimpan kunci web.
	KEY_ENV_KUNCI = "KunciWeb"

	// FILEPATH_LOG_LINUX adalah direktori log di sistem operasi Linux.
	FILEPATH_LOG_LINUX = "/var/log/nginx/api"
	// FILEPATH_LOG_WINDOWS adalah direktori log di sistem operasi Windows.
	FILEPATH_LOG_WINDOWS = "D:/go-apps/logs"

	// FILEPATH_SETTINGWEB_WINDOWS adalah path file SettingWeb.xml di sistem operasi Windows.
	FILEPATH_SETTINGWEB_WINDOWS = "D:\\_docker\\_app\\kunci\\SettingWeb.xml"
	// FILEPATH_SETTINGWEB_LINUX adalah path file SettingWeb.xml di sistem operasi Linux.
	FILEPATH_SETTINGWEB_LINUX = "SettingWeb.xml"

	// DATE_FORMAT adalah format tanggal standar yang digunakan di seluruh library.
	DATE_FORMAT = "2006-01-02"
	// DATETIME_FORMAT adalah format tanggal dan waktu standar yang digunakan di seluruh library.
	DATETIME_FORMAT = "2006-01-02 15:04:05"

	// DBTYPE_POSTGRE adalah tipe database PostgreSQL.
	DBTYPE_POSTGRE = "POSTGRE"

	// TIME_FIVE_MINUTES adalah durasi lima menit.
	TIME_FIVE_MINUTES = 300 * time.Second
	// TIME_TWO_MINUTES adalah durasi dua menit.
	TIME_TWO_MINUTES = 120 * time.Second
	// TIME_ONE_MINUTE adalah durasi satu menit.
	TIME_ONE_MINUTE = 60 * time.Second
	// TIME_TEN_SECONDS adalah durasi sepuluh detik.
	TIME_TEN_SECONDS = 10 * time.Second

	// PREFIX_KUNCI adalah prefix yang digunakan untuk mengidentifikasi kunci dalam aplikasi.
	PREFIX_KUNCI = "kunci"
)
