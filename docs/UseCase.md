# Use Case libsd3

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