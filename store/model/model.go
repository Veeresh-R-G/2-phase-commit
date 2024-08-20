package model

type Food struct {
	Id   int
	Name string
}

type Packet struct {
	Id          int
	Food_id     int
	Is_reserved bool
	Order_id    int
}
