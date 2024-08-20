package main

import (
	"errors"
	"log"

	"github.com/Veeresh-R-G/2-phase-commit-protocol/database"
	"github.com/Veeresh-R-G/2-phase-commit-protocol/delivery/model"
	"github.com/gin-gonic/gin"
)

type BookAgentRequest struct {
	Order_id int `json:"order_id"`
}

func ReserveAgent() (*model.Agent, error) {
	// Reserve an agent for delivery
	db, err := database.InitialiseDb()

	if err != nil {
		return nil, err
	}

	txn, _ := db.Begin()

	row := txn.QueryRow(`SELECT id, is_reserved, order_id from agents
	WHERE is_reserved is false and order_id is -1
	LIMIT 1
	FOR UPDATE`)

	if row.Err() != nil {
		txn.Rollback()
		return nil, row.Err()
	}

	var agent model.Agent
	err = row.Scan(&agent.Id, &agent.Is_reserved, &agent.Order_id)

	if err != nil {
		txn.Rollback()
		return nil, errors.New("no agents available")
	}

	_, err = txn.Exec(`UPDATE agents 
	SET
	is_reserved = true
	WHERE id = ?`, agent.Id)

	if err != nil {
		txn.Rollback()
		return nil, errors.New("error reserving agent")
	}

	return &agent, nil
}

func BookAgent(order_id int) (*model.Agent, error) {

	db, err := database.InitialiseDb()
	if err != nil {
		return nil, errors.New("error connecting to database")
	}

	txn, _ := db.Begin()

	row := txn.QueryRow(`SELECT id, is_reserved, order_id FROM agents
	WHERE is_reserved is true and order_id is -1
	LIMIT 1
	FOR UPDATE`)

	if row.Err() != nil {
		txn.Rollback()
		return nil, errors.New("no agents is free, all are busy delivering the package")
	}

	var agent model.Agent
	err = row.Scan(&agent.Id, &agent.Is_reserved, &agent.Order_id)
	if err != nil {
		txn.Rollback()
		return nil, errors.New("error scanning agent")
	}

	//Booking the agent for a particular order
	_, err = txn.Exec(`UPDATE agents
	SET 
	is_reserved = false,
	order_id = ?
	WHERE id = ?`, order_id, agent.Id)

	if err != nil {
		txn.Rollback()
		return nil, errors.New("error updating agent")
	}

	return &agent, nil
}

func main() {

	r := gin.Default()

	r.POST("/delivery/agent/reserve", func(c *gin.Context) {

		agent, err := ReserveAgent()
		if err != nil {
			c.JSON(429, err)
			return
		}

		c.JSON(200, gin.H{"agent_id": agent.Id})
	})

	r.POST("/delivery/agent/book", func(ctx *gin.Context) {

		var reqBody BookAgentRequest

		if err := ctx.BindJSON(&reqBody); err != nil {
			ctx.JSON(400, err)
			return
		}

		agent, err := BookAgent(reqBody.Order_id)

		if err != nil {
			ctx.JSON(429, err)
			return
		}

		ctx.JSON(200, gin.H{"Booked Agent with agent_id": agent.Id, "Booked for Order-id": reqBody.Order_id})
	})

	log.Printf("Starting delivery service on port 8082")
	r.Run(":8082")

}
