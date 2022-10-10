package dto

type InsertProductDto struct {
	Title       string  `json:"title" validate:"required"`
	Description string  `json:"description"`
	BrandId     int     `json:"brandId" validate:"required"`
	Price       float32 `json:"price" validate:"required"`
}

type FilterProductDto struct {
	ID          int    `json:"id"`
	BrandId     int    `json:"brandId" validate:"required"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Limit       int    `json:"limit"`
}

type GetProduct struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Brand       BrandDto `json:"brand"`
	Price       float32  `json:"price"`
}

type BrandDto struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}
