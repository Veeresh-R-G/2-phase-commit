- Populated Data in the database

````
stmtsIns, err := db.Prepare("INSERT INTO food VALUES( ?, ?)")
	if err != nil {
		log.Fatalln("Can't Prepare Insert Statements", err)
	}

	defer stmtsIns.Close()

	for i := 0; i < 5; i++ {
		_, err := stmtsIns.Exec(i+1, "food - "+strconv.Itoa(i+1))
		if err != nil {
			log.Fatalln("Error inserting data", err)
		}
	}

	txn, err := db.Begin()
	if err != nil {
		log.Fatalln("Error beginning transaction")
	}

	rows, err := txn.Query("SELECT * FROM food")
	if err != nil {
		log.Fatalln("Error querying data")
	}

	for rows.Next() {
		var id int
		var food string

		err = rows.Scan(&id, &food)
		if err != nil {
			log.Fatalln("Error scanning data", err)
		}
		log.Println(id, food)
		log.Println("--------------------------")
	}
```w
````
