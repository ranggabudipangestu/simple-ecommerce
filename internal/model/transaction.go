package model

import "time"

type Transaction struct {
	ID                int
	DeliveryAddress   int
	TransactionNumber string
	TotalTransaction  float32
	CreatedAt         time.Time
}

type TransactionDetails struct {
	ID            int
	TransactionId int
	ProductId     int
	Qty           int
	Price         float32
	Total         float32
}
