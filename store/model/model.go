package model

type Food struct {
	id   int
	name string
}

type Packet struct {
	id          int
	food_id     int
	is_reserved bool
	order_id    int
}
