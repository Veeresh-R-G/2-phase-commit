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

	router := httprouter.New()
	/*
		4 APIs
		1. /delivery/agent/reserve -> reserve for any order (Prepare)
		2. /delivery/agent/book -> assign a particular order (commit phase)

		3. /resto/food/reserve
		4. /resto/food/book
	*/
	router.GET("/resto", hello)

	router.POST("/delivery/agent/reserve", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// reserve for any order (Prepare)

	})

	router.POST("/delivery/agent/book", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		// assign a particular order (commit phase)

	})

	log.Fatalln(http.ListenAndServe(":8080", router))
}
