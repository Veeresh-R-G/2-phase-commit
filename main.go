package main

import (
	"log"

	"github.com/Veeresh-R-G/2-phase-commit-protocol/database"
)

func main() {

	_, err := database.InitialiseDb()

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database initialised successfully")

}
