package dto

type CreateOrderDto struct {
	DeliveryAddress  string               `json:"deliveryAddress" validate:"required"`
	Details          []CreateOrderDetails `json:"details" validate:"required"`
	TotalTransaction float32
	TotalQty         int
}

type CreateOrderDetails struct {
	ProductId int `json:"productId" validate:"required"`
	Price     float32
	Qty       int `json:"qty" validate:"required"`
	Total     float32
}

type GetOrderDto struct {
	ID                int               `json:"id"`
	DeliveryAddress   string            `json:"deliveryAddress"`
	TransactionNumber string            `json:"transactionNumber"`
	TotalTransaction  float32           `json:"totalTransaction"`
	TotalQty          float32           `json:"totalQty"`
	Details           []GetOrderDetails `json:"details"`
}

type GetOrderDetails struct {
	ID          int     `json:"id"`
	ProductName string  `json:"productName"`
	BrandName   string  `json:"brandName"`
	Qty         int     `json:"qty"`
	Price       float32 `json:"price"`
	Total       float32 `json:"total"`
}
