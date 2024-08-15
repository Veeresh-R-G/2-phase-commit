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

	defer db.Close()

}
