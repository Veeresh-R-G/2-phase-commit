package main

import (
	"log"

	"github.com/Veeresh-R-G/2-phase-commit-protocol/database"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	db, err := database.InitialiseDb()

	if err != nil {
		log.Fatalln("Error connecting to database")
	}
	query := `
    SELECT * from agents;`
	// query := `
	// CREATE TABLE IF NOT EXISTS agents (
	//     id INT AUTO_INCREMENT,
	//     is_reserved BOOLEAN,
	// 	order_id INT NULL,
	//     PRIMARY KEY (id)
	// );`
	txn, _ := db.Begin()

	response, err := txn.Exec(query)
	if err != nil {
		log.Fatalln("Error creating table")
	}
	rowsAffected, err := response.RowsAffected()
	if err != nil {
		log.Fatalf("Failed to retrieve rows affected: %v", err)
	}

	log.Println(rowsAffected)
}
