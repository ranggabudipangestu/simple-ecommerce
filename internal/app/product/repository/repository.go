package repository

import (
	"context"
	"database/sql"

	"github.com/ranggabudipangestu/simple-ecommerce/internal/app/product/dto"
	"github.com/ranggabudipangestu/simple-ecommerce/internal/model"
	"github.com/ranggabudipangestu/simple-ecommerce/pkg/util"
)

type ProductRepository interface {
	Create(ctx context.Context, payload dto.InsertProductDto) (*model.Product, error)
	GetProduct(ctx context.Context, filter dto.FilterProductDto) (data []dto.GetProduct, err error)
}

type Repository struct {
	DB *sql.DB
}

func NewProduct(db *sql.DB) *Repository {
	return &Repository{db}
}

func (p *Repository) Create(ctx context.Context, payload dto.InsertProductDto) (*model.Product, error) {

	query := `INSERT INTO product (title, description, brandId, price, createdAt, updatedAt) values(?, ?, ?, ?, NOW(), NOW())`
	result, err := p.DB.ExecContext(ctx, query, payload.Title, payload.Description, payload.BrandId, payload.Price)
	if err != nil {
		return nil, err
	}

	var id int64
	id, err = result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &model.Product{ID: int(id)}, nil
}
func (p *Repository) GetProduct(ctx context.Context, filter dto.FilterProductDto) (data []dto.GetProduct, err error) {
	var filterValues []interface{}
	query := `SELECT product.id, product.title, product.description, 
	product.brandId, brand.title as brandTitle,
	product.price 
	FROM product
	JOIN brand ON product.brandId = brand.id
	`

	if filter.ID > 0 {
		query += util.FilterHandler(filterValues) + ` product.id = ?`
		filterValues = append(filterValues, filter.ID)
	}

	if filter.Title != "" {
		query += util.FilterHandler(filterValues) + ` title = ?`
		filterValues = append(filterValues, filter.Title)
	}

	if filter.BrandId > 0 {
		query += util.FilterHandler(filterValues) + ` brand.id = ?`
		filterValues = append(filterValues, filter.BrandId)
	}

	query += ` ORDER BY id `

	if filter.Limit > 0 {
		filterValues = append(filterValues, filter.Limit)
		query += ` LIMIT ?`
	}

	var rows *sql.Rows
	rows, err = p.DB.QueryContext(ctx, query, filterValues...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		Repository := dto.GetProduct{}
		err = rows.Scan(
			&Repository.ID,
			&Repository.Title,
			&Repository.Description,
			&Repository.Brand.ID,
			&Repository.Brand.Title,
			&Repository.Price,
		)

		data = append(data, Repository)
	}

	return data, nil
}
