package repository_test

import (
	"context"
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

	query := "SELECT id, title from brand ORDER BY id"

	mock.ExpectQuery(query).WillReturnRows(rows)
	r := repository.NewBrand(db)
	filter := dto.FilterBrandDto{}
	result, err := r.GetBrand(context.TODO(), filter)

	assert.Nil(t, err)
	assert.NotNil(t, result)
}
