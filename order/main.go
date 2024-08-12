package order

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Order() {
	fmt.Println("From Order")

	router := httprouter.New()

	router.GET("/resto", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Hello Orders",
			"status":  http.StatusOK,
		})
	})

	log.Fatal(http.ListenAndServe(":8082", router))
}
