package database

import "log"

func CreateTableAgent() {
	query := `
    CREATE TABLE IF NOT EXISTS agents (
        id INT AUTO_INCREMENT,
        is_reserved BOOLEAN,
		order_id INT NULL,
        PRIMARY KEY (id)
    );`

	db, err := InitialiseDb()

	if err != nil {
		log.Fatalln("Error connecting to database")
	}

	txn, err := db.Begin()

	if err != nil {
		log.Fatalln("Could not begin transaction")
	}

	result, err := txn.Exec(query)

	if err != nil {
		log.Fatalln("Error creating table")
	}

	log.Println(result)

}
