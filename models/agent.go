package models

type Agent struct {
	id          int
	is_reserved bool
	order_id    int //can be null in table
}
