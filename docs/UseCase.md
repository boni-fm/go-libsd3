# Use Case go-libsd3

---
## dbutil
```
db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	row := db.SelectScalar("SELECT $1", "test select")
	var name string
	if err := row.Scan(&name); err != nil {
		log.Fatal(err)
	}
	fmt.Println(name)
```

## logging
```
log := logging.NewLogger()
log.Say("Test log message")
log.Sayf("Test logf %d", 123)
log.SayWithField("Test with field", "foo", "bar")
log.SayWithFields("Test with fields", map[string]interface{}{"a": 1, "b": "c"})
log.Warn("Test warn level")
log.Errorf("Test errorf %s", "err")
```
kesimpen di folder `{home directory}/_docker/_app/logs/logs{tanggal}`