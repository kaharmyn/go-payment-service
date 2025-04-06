package main

type Payment struct {
	ID      string  `json:"id"`
	OrderID int     `json:"order_id"`
	Amount  float64 `json:"amount"`
	Status  string  `json:"status"`
}
