package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func InitialiseDb() (*sql.DB, error) {

	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	connectionString := os.Getenv("DB_CONN")

	db, err := sql.Open("mysql", connectionString)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}
	return db, nil
}
