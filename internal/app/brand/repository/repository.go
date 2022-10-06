package repository

import (
	"context"
	"database/sql"

	"github.com/ranggabudipangestu/simple-ecommerce/internal/app/brand/dto"
	"github.com/ranggabudipangestu/simple-ecommerce/internal/model"
	"github.com/ranggabudipangestu/simple-ecommerce/pkg/util"
)

type BrandRepository interface {
	Create(ctx context.Context, dto dto.InsertBrandDto) (*model.Brand, error)
	GetBrand(ctx context.Context, dto dto.FilterBrandDto) (data []model.Brand, err error)
}

type Repository struct {
	DB *sql.DB
}

func NewBrand(db *sql.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) Create(ctx context.Context, dto dto.InsertBrandDto) (*model.Brand, error) {
	sqlCommand := "INSERT into brand (title, createdAt, updatedAt) values(?, NOW(), NOW())"
	result, err := r.DB.ExecContext(ctx, sqlCommand, &dto.Title)
	if err != nil {
		return nil, err
	}

	var id int64
	id, err = result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &model.Brand{ID: int(id)}, nil
}
func (r *Repository) GetBrand(ctx context.Context, filter dto.FilterBrandDto) (data []model.Brand, err error) {

	var filterValues []interface{}
	query := "SELECT id, title from brand"

	if filter.ID > 0 {
		query += util.FilterHandler(filterValues) + ` id = ?`
		filterValues = append(filterValues, filter.ID)
	}

	if filter.Title != "" {
		query += util.FilterHandler(filterValues) + ` title = ?`
		filterValues = append(filterValues, filter.Title)
	}

	query += ` ORDER BY id `

	if filter.Limit > 0 {
		filterValues = append(filterValues, filter.Limit)
		query += ` LIMIT ?`
	}

	var rows *sql.Rows
	rows, err = r.DB.QueryContext(ctx, query, filterValues...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		brand := model.Brand{}
		err = rows.Scan(
			&brand.ID,
			&brand.Title,
		)

		data = append(data, brand)
	}

	return data, nil
}
