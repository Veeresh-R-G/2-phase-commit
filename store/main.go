package main

import (
	"database/sql"
	"errors"
	"log"

	"github.com/Veeresh-R-G/2-phase-commit-protocol/database"
	"github.com/Veeresh-R-G/2-phase-commit-protocol/store/model"
	"github.com/gin-gonic/gin"
)

type FoodRequestBody struct {
	Id          int  `json:"id"`
	Food_id     int  `json:"food_id"`
	Is_reserved bool `json:"is_reserved"`
	Order_id    int  `json:"order_id"`
}

func ReserveFood(food_id int) (*model.Packet, error) {

	var food_packet model.Packet
	db, err := database.InitialiseDb()

	if err != nil {
		return nil, errors.New("error in connecting to database")
	}

	txn, _ := db.Begin()

	//Getting the first food packet that is available
	row := txn.QueryRow(`SELECT * id, food_id, is_reserved, order_id 
	FROM packet
	WHERE
		is_reserved is false and food_id = ? and order_id = -1
	LIMIT 1
	FOR UPDATE
	`, food_id)

	if row.Err() != nil {
		return nil, errors.New("error in getting row")
	}

	err = row.Scan(&food_packet.Id, &food_packet.Food_id, &food_packet.Is_reserved, &food_packet.Order_id)
	if err != nil {
		txn.Rollback()
		return nil, errors.New("no food packet available")
	}

	_, err = txn.Exec(`UPDATE packet 
	SET is_reserved = true
	WHERE id = ?
	`, food_packet.Id)

	if err != nil {
		txn.Rollback()
		return nil, errors.New("error in updating packet")
	}

	err = txn.Commit()
	if err != nil {
		return nil, errors.New("error in committing updating transaction")
	}

	return &food_packet, nil
}

func BookFood(order_id, food_id int) (*model.Packet, error) {

	db, err := database.InitialiseDb()
	if err != nil {
		return nil, errors.New("error in connecting to database")
	}

	txn, _ := db.Begin()

	var packet model.Packet
	row := txn.QueryRow(`
	SELECT * from packet
	WHERE is_reserved is true and order_id = -1 and food_id = ?
	LIMIT 1
	FOR UPDATE`, food_id)

	if row.Err() != nil {
		return nil, errors.New("food packets not available")
	}

	err = row.Scan(&packet.Id, &packet.Food_id, &packet.Is_reserved, &packet.Order_id)
	if err != nil && err == sql.ErrNoRows {
		return nil, errors.New("food packets not available")

	}

	_, err = txn.Exec(`
	UPDATE packet
	SET
	is_reserved = false, order_id = ? 
	WHERE 
	id = ?`, order_id, packet.Id)

	if err != nil {
		txn.Rollback()
		return nil, errors.New("error in updating packet")
	}

	return &packet, nil
}

func main() {

	r := gin.Default()

	r.POST("/store/food/reserve", func(ctx *gin.Context) {

	})

	r.POST("/store/food/book", func(ctx *gin.Context) {
		var reqBody FoodRequestBody

		if err := ctx.BindJSON(&reqBody); err != nil {
			ctx.JSON(400, err)
			return
		}

		_, _ = BookFood(reqBody.Order_id, reqBody.Food_id)
	})
	r.Run(":8081")

	log.Println("Hello, World!")
}
