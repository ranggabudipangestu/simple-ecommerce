package dto

type InsertBrandDto struct {
	Title string `json:"title" validate:"required"`
}

type FilterBrandDto struct {
	ID    int    `json:"int"`
	Title string `json:"title"`
	Limit int    `json:"limit"`
}

type GetBrand struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}
