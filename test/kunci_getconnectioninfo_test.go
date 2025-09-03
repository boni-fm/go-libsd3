package test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/boni-fm/go-libsd3/helper/settinglibgooo"
)

func TestGetConnectionInfoPostgreE(t *testing.T) {
	// Determine the correct path based on OS
	var xmlPath string
	if runtime.GOOS == "windows" {
		dir := filepath.Join("D:", "_docker", "_app", "kunci")
		os.MkdirAll(dir, 0755)
		xmlPath = filepath.Join(dir, "SettingWeb.xml")
	} else {
		dir := "/_docker/_app/kunci"
		os.MkdirAll(dir, 0755)
		xmlPath = filepath.Join(dir, "SettingWeb.xml")
	}

	xmlContent := `<?xml version="1.0" encoding="UTF-8"?>
<SettingConfig>
	<IPPostgres>127.0.0.1</IPPostgres>
	<PortPostgres>5432</PortPostgres>
	<DatabasePostgres>testdb</DatabasePostgres>
	<UserPostgres>testuser</UserPostgres>
	<PasswordPostgres>testpass</PasswordPostgres>
</SettingConfig>`
	os.WriteFile(xmlPath, []byte(xmlContent), 0644)
	defer os.Remove(xmlPath)

	conn := settinglibgooo.GetConnectionInfoPostgre()
	if conn.IPPostgres != "127.0.0.1" || conn.PortPostgres != "5432" || conn.DatabasePostgres != "testdb" || conn.UserPostgres != "testuser" || conn.PasswordPostgres != "testpass" {
		t.Logf("Unexpected connection info: %+v", conn)
	}
}
