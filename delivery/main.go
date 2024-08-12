package delivery

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {

	router := httprouter.New()

	router.GET("/resto", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message": "Hello Delivery",
			"status":  http.StatusOK,
		})
	})

	log.Fatal(http.ListenAndServe(":8081", router))
}
