package repository_test

import (
	"context"
	"errors"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/ranggabudipangestu/simple-ecommerce/internal/app/brand/dto"
	"github.com/ranggabudipangestu/simple-ecommerce/internal/app/brand/repository"
	"github.com/ranggabudipangestu/simple-ecommerce/internal/model"
)

func TestGetBrand(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mockBrand := []model.Brand{}
	mockBrand = append(mockBrand, model.Brand{ID: 1, Title: "Nike"})
	mockBrand = append(mockBrand, model.Brand{ID: 2, Title: "Adidas"})

	rows := sqlmock.NewRows([]string{"id", "title"}).
		AddRow(mockBrand[0].ID, mockBrand[0].Title).
		AddRow(mockBrand[1].ID, mockBrand[1].Title)

	reset := func() {
		mockBrand := []model.Brand{}
		mockBrand = append(mockBrand, model.Brand{ID: 1, Title: "Nike"})
		mockBrand = append(mockBrand, model.Brand{ID: 2, Title: "Adidas"})

		rows = sqlmock.NewRows([]string{"id", "title"}).
			AddRow(mockBrand[0].ID, mockBrand[0].Title).
			AddRow(mockBrand[1].ID, mockBrand[1].Title)
	}

	t.Run("Test Brand Without Filter", func(t *testing.T) {

		defer reset()
		query := "SELECT id, title from brand ORDER BY id"

		mock.ExpectQuery(query).WillReturnRows(rows)
		r := repository.NewBrand(db)
		filter := dto.FilterBrandDto{}
		result, err := r.GetBrand(context.TODO(), filter)

		assert.Nil(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Test Brand With Filter", func(t *testing.T) {

		query := "SELECT id, title from brand"

		filter := dto.FilterBrandDto{ID: 1, Title: "Nike", Limit: 1}
		mock.ExpectQuery(query).WithArgs(filter.ID, filter.Title, filter.Limit).WillReturnRows(rows)
		r := repository.NewBrand(db)
		result, err := r.GetBrand(context.TODO(), filter)

		assert.Nil(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Test Brand Error", func(t *testing.T) {

		defer reset()
		query := "SELECT id, title from brand ORDER BY id"

		mock.ExpectQuery(query).WillReturnError(errors.New("Database Error"))
		r := repository.NewBrand(db)
		filter := dto.FilterBrandDto{}
		result, err := r.GetBrand(context.TODO(), filter)

		assert.NotNil(t, err)
		assert.Nil(t, result)
	})

}

func TestCreateBrand(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	t.Run("Test Create Brand Success", func(t *testing.T) {

		payload := dto.InsertBrandDto{
			Title: "Puma",
		}
		query := "INSERT into brand"

		mock.ExpectExec(query).WithArgs(payload.Title).WillReturnResult(sqlmock.NewResult(1, 1))
		r := repository.NewBrand(db)
		result, err := r.Create(context.TODO(), payload)

		assert.Nil(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Test Create Brand Error", func(t *testing.T) {

		payload := dto.InsertBrandDto{
			Title: "Puma",
		}
		query := "INSERT into brand"

		mock.ExpectExec(query).WithArgs(payload.Title).WillReturnError(errors.New("Error From Database"))
		r := repository.NewBrand(db)
		result, err := r.Create(context.TODO(), payload)

		assert.NotNil(t, err)
		assert.Nil(t, result)
	})
}
