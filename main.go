package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

func hello(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	resp := make(map[string]interface{})

	resp["message"] = "Hello World"
	resp["status"] = http.StatusOK

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(resp)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	connString := os.Getenv("DB_CONN")
	db, err := sql.Open("mysql", connString)

	if err != nil {
		log.Fatalln(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}

	var resp string
	err = db.QueryRow("SELECT 'HELLO WORLD'").Scan(&resp)
	if err != nil {
		log.Fatalln(err)
	}

	router := httprouter.New()

	router.GET("/resto", hello)

	log.Fatalln(http.ListenAndServe(":8080", router))
}
