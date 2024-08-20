package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/Veeresh-R-G/2-phase-commit-protocol/orders/model"
	"github.com/google/uuid"
)

func PlaceOrder(food_id int) (*model.Order, error) {

	body, _ := json.Marshal(map[string]interface{}{
		"food_id": food_id,
	})

	reqBody := bytes.NewBuffer(body)

	//------------------ Preparation Phase ------------------

	//reserve food
	reserve_food_resp, err := http.Post("http://localhost:8081/store/food/reserve", "application/json", reqBody)
	if err != nil || reserve_food_resp.StatusCode != 200 {
		log.Println(`Error in reserving food : `, err.Error())
		return nil, errors.New(`error in reserving food`)
	}

	//reserve agent
	reserve_agent_resp, err := http.Post("http://localhost:8082/delivery/agent/reserve", "application/json", nil)
	if err != nil || reserve_agent_resp.StatusCode != 200 {
		log.Println(`Error in reserving agent : `, err.Error())
		return nil, errors.New(`error in reserving agent`)
	}

	//------------------ Preparation Phase Completed ------------------

	//----------------------- Commit Phase ------------------------
	order_id := uuid.New().String()
	body, _ = json.Marshal(map[string]interface{}{
		"order_id": order_id,
		"food_id":  food_id,
	})
	reqBody = bytes.NewBuffer(body)

	//assign food to order
	book_food_resp, err := http.Post("http://localhost:8081/store/food/book", "application/json", reqBody)
	if err != nil || book_food_resp.StatusCode != 200 {
		log.Println(`Error in booking food : `, err.Error())
		return nil, errors.New(`error in booking food`)
	}

	//assign agent to order
	body, _ = json.Marshal(map[string]interface{}{
		"order_id": order_id,
	})
	reqBody = bytes.NewBuffer(body)
	book_agent_resp, err := http.Post("http://localhost:8082/delivery/agent/book", "application/json", reqBody)
	if err != nil || book_agent_resp.StatusCode != 200 {
		log.Println(`Error in booking agent : `, err.Error())
		return nil, errors.New(`error in booking agent`)
	}

	//----------------------- Commit Phase Completed ------------------------

	return &model.Order{OrderId: order_id}, nil

}

func main() {

	foodId := 1

	var wg sync.WaitGroup
	wg.Add(10)

	for i := 0; i < 5; i++ {
		go func() {
			order, err := PlaceOrder(foodId)
			wg.Done()
			if err != nil {
				fmt.Println(`Order not Placed : `, err.Error())
			} else {
				log.Println(`Order Placed : `, order)
			}

			log.Println(`Order Place : `, order.OrderId)

		}()
	}

	wg.Wait()
}
