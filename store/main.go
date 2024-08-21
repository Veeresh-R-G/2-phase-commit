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
	Food_id  int `json:"food_id"`
	Order_id int `json:"order_id"`
}

func ReserveFood(food_id int) (*model.Packet, error) {

	var food_packet model.Packet
	db, err := database.InitialiseDb()

	if err != nil {
		return nil, errors.New("error in connecting to database")
	}

	txn, _ := db.Begin()

	//Getting the first food packet that is available
	row := txn.QueryRow(`SELECT id, food_id, is_reserved, order_id 
	FROM packet
	WHERE
		is_reserved is false and food_id = ? and order_id = -1
	LIMIT 1
	FOR UPDATE
	`, food_id)

	if row.Err() != nil {
		return nil, errors.New("error in getting row / food packet not available")
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
	SELECT id, food_id, is_reserved, order_id from packet
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

	err = txn.Commit()
	if err != nil {
		log.Println("Error in committing transaction : ", err)
		return nil, errors.New("error in committing transaction")
	}

	return &packet, nil
}

func main() {

	r := gin.Default()

	r.POST("/store/food/reserve", func(ctx *gin.Context) {
		var reqBody FoodRequestBody

		if err := ctx.BindJSON(&reqBody); err != nil {
			ctx.JSON(400, err)
			return
		}

		packet, err := ReserveFood(reqBody.Food_id)
		if err != nil {
			log.Println("Error in reserving food : ", err)
			ctx.JSON(429, err)
			return
		}

		ctx.JSON(200, gin.H{"packet reserved": packet})

	})

	r.POST("/store/food/book", func(ctx *gin.Context) {
		var reqBody FoodRequestBody

		if err := ctx.BindJSON(&reqBody); err != nil {
			ctx.JSON(400, err)
			return
		}

		packet, err := BookFood(reqBody.Order_id, reqBody.Food_id)
		if err != nil {
			ctx.JSON(429, err)
			return
		}

		ctx.JSON(200, gin.H{"packet booked": packet})

	})

	log.Println("Store Service started on port 8081")
	r.Run(":8081")
}
