package repository_test

import (
	"context"
	"errors"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/ranggabudipangestu/simple-ecommerce/internal/app/product/dto"
	"github.com/ranggabudipangestu/simple-ecommerce/internal/app/product/repository"
)

func TestGetProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	query := `SELECT product.id, product.title, product.description, 
		product.brandId, brand.title as brandTitle,
		product.price 
		FROM product
		JOIN brand ON product.brandId = brand.id
		`

	mockProduct := []dto.GetProduct{}
	mockProduct = append(mockProduct, dto.GetProduct{ID: 1, Title: "Nike Airmax", Description: "Sepatu Nike", Brand: dto.BrandDto{ID: 1, Title: "Nike"}, Price: 2000000})
	mockProduct = append(mockProduct, dto.GetProduct{ID: 1, Title: "Adidas Duramo", Description: "Sepatu Adidas", Brand: dto.BrandDto{ID: 2, Title: "Adidas"}, Price: 1500000})

	rows := sqlmock.NewRows([]string{"id", "title", "description", "brandId", "brandTitle", "price"}).
		AddRow(mockProduct[0].ID, mockProduct[0].Title, mockProduct[0].Description, mockProduct[0].Brand.ID, mockProduct[0].Brand.Title, mockProduct[0].Price).
		AddRow(mockProduct[1].ID, mockProduct[1].Title, mockProduct[1].Description, mockProduct[1].Brand.ID, mockProduct[1].Brand.Title, mockProduct[0].Price)

	reset := func() {
		mockProduct = []dto.GetProduct{}
		mockProduct = append(mockProduct, dto.GetProduct{ID: 1, Title: "Nike Airmax", Description: "Sepatu Nike", Brand: dto.BrandDto{ID: 1, Title: "Nike"}, Price: 2000000})
		mockProduct = append(mockProduct, dto.GetProduct{ID: 1, Title: "Adidas Duramo", Description: "Sepatu Adidas", Brand: dto.BrandDto{ID: 2, Title: "Adidas"}, Price: 1500000})

		sqlmock.NewRows([]string{"id", "title", "description", "brandId", "brandTitle", "price"}).
			AddRow(mockProduct[0].ID, mockProduct[0].Title, mockProduct[0].Description, mockProduct[0].Brand.ID, mockProduct[0].Brand.Title, mockProduct[0].Price).
			AddRow(mockProduct[1].ID, mockProduct[1].Title, mockProduct[1].Description, mockProduct[1].Brand.ID, mockProduct[1].Brand.Title, mockProduct[0].Price)
	}

	t.Run("Test Product Without Filter", func(t *testing.T) {

		defer reset()

		mock.ExpectQuery(query).WillReturnRows(rows)
		r := repository.NewProduct(db)
		filter := dto.FilterProductDto{}
		result, err := r.GetProduct(context.TODO(), filter)

		assert.Nil(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Test Product With Filter", func(t *testing.T) {

		defer reset()

		filter := dto.FilterProductDto{
			ID:      1,
			BrandId: 1,
			Title:   "nike",
			Limit:   1,
		}

		mock.ExpectQuery(query).WithArgs(filter.ID, filter.Title, filter.BrandId, filter.Limit).WillReturnRows(rows)
		r := repository.NewProduct(db)

		_, err := r.GetProduct(context.TODO(), filter)

		assert.Nil(t, err)
	})

	t.Run("Test Product Error DB", func(t *testing.T) {

		defer reset()

		mock.ExpectQuery(query).WillReturnError(errors.New("Database Error"))
		r := repository.NewProduct(db)
		filter := dto.FilterProductDto{}
		_, err := r.GetProduct(context.TODO(), filter)

		assert.NotNil(t, err)
	})

}

func TestCreateProduct(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	query := "INSERT INTO product"

	payload := dto.InsertProductDto{
		Title:       "Nike Airmax",
		Description: "Sepatu Nike",
		BrandId:     1,
		Price:       1250000,
	}

	t.Run("Test Create Brand Success", func(t *testing.T) {

		mock.ExpectExec(query).WithArgs(payload.Title, payload.Description, payload.BrandId, payload.Price).WillReturnResult(sqlmock.NewResult(1, 1))
		r := repository.NewProduct(db)
		result, err := r.Create(context.TODO(), payload)

		assert.Nil(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Test Create Brand Error Database", func(t *testing.T) {

		mock.ExpectExec(query).WithArgs(payload.Title, payload.Description, payload.BrandId, payload.Price).WillReturnError(errors.New("Database Error"))
		r := repository.NewProduct(db)
		result, err := r.Create(context.TODO(), payload)

		assert.NotNil(t, err)
		assert.Nil(t, result)
	})
}
