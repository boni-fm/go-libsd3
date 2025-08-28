package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/boni-fm/go-libsd3/helper/kunci"
)

func TestGetConnectionInfoPostgre(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get user home dir: %v", err)
	}
	// Prepare a dummy SettingWeb.xml in a temp dir under home
	dir := filepath.Join(home, "_docker", "kunci")
	os.MkdirAll(dir, 0755)
	filePath := filepath.Join(dir, "SettingWeb.xml")
	xmlContent := `<?xml version="1.0" encoding="UTF-8"?>
<SettingConfig>
	<IPPostgres>127.0.0.1</IPPostgres>
	<PortPostgres>5432</PortPostgres>
	<DatabasePostgres>testdb</DatabasePostgres>
	<UserPostgres>testuser</UserPostgres>
	<PasswordPostgres>testpass</PasswordPostgres>
</SettingConfig>`
	os.WriteFile(filePath, []byte(xmlContent), 0644)
	defer os.Remove(filePath)

	conn := kunci.GetConnectionInfoPostgre()
	if conn.IPPostgres == "" {
		t.Errorf("Unexpected connection info: %+v", conn)
	}
}

func TestGetConnectionString(t *testing.T) {
	//home, _ := os.UserHomeDir()
	//dir := filepath.Join(home, "_docker", "kunci")
	//filePath := filepath.Join(dir, "SettingWeb.xml")
	xmlContent := `<?xml version="1.0" encoding="UTF-8"?>
<SettingConfig>
	<IPPostgres>localhost</IPPostgres>
	<PortPostgres>5433</PortPostgres>
	<DatabasePostgres>db2</DatabasePostgres>
	<UserPostgres>user2</UserPostgres>
	<PasswordPostgres>pass2</PasswordPostgres>
</SettingConfig>`
	os.WriteFile("/_docker/_app/_kunci/SettingWeb.xml", []byte(xmlContent), 0644)
	defer os.Remove("/_docker/_app/_kunci/SettingWeb.xml")

	// cara make kuncinya
	confKunci, _ := kunci.ReadConfig("config.yaml")
	kunciManager := kunci.NewKunci(confKunci)
	connStr := kunciManager.GetConnectionString("POSTGRE")

	expected := "host=localhost port=5433 user=user2 password=pass2 dbname=db2 sslmode=disable"
	t.Logf("Connection String: %s", connStr)
	if connStr == expected {
		t.Errorf("Expected '%s', got '%s'", expected, connStr)
	}
}
