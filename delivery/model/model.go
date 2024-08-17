package model

type Agent struct {
	Id          int
	Is_reserved bool
	Order_id    int //can be null in table
}
