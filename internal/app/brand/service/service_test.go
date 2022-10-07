package service_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/ranggabudipangestu/simple-ecommerce/internal/app/brand/dto"
	mockRepositories "github.com/ranggabudipangestu/simple-ecommerce/internal/app/brand/mocks/repository"

	"github.com/go-faker/faker/v4"
	service "github.com/ranggabudipangestu/simple-ecommerce/internal/app/brand/service"
	"github.com/ranggabudipangestu/simple-ecommerce/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const contextTimeout = 2 * time.Second

func TestBrandCheckById(t *testing.T) {
	// var ctx context.Context
	var mockBrand []model.Brand
	mockRepository := new(mockRepositories.BrandRepository)

	reset := func() {
		mockBrand = []model.Brand{}
		mockRepository = new(mockRepositories.BrandRepository)
	}

	t.Run("Test Brand By Id Success", func(t *testing.T) {
		defer reset()

		err := faker.FakeData(&mockBrand)
		assert.NoError(t, err)

		brandService := service.NewBrandService(mockRepository, contextTimeout)
		filter := dto.FilterBrandDto{ID: 1, Limit: 1}
		mockRepository.On("GetBrand", mock.Anything, filter).Return(mockBrand, nil)

		res := brandService.CheckBrandById(context.TODO(), 1)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.True(t, res.Success)
		assert.Equal(t, "success", res.Message)
		assert.NotNil(t, res.Data)
	})

	t.Run("Test Brand By ID Database Error", func(t *testing.T) {
		defer reset()

		brandService := service.NewBrandService(mockRepository, contextTimeout)
		filter := dto.FilterBrandDto{ID: 1, Limit: 1}
		mockRepository.On("GetBrand", mock.Anything, filter).Return(mockBrand, errors.New("Database Error"))
		res := brandService.CheckBrandById(context.TODO(), 1)

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
		assert.False(t, res.Success)
		assert.Equal(t, "Database Error", res.Message)
		assert.Nil(t, res.Data)
	})

	t.Run("Test Brand By ID Not Found", func(t *testing.T) {
		defer reset()

		brandService := service.NewBrandService(mockRepository, contextTimeout)
		filter := dto.FilterBrandDto{ID: 1, Limit: 1}
		mockRepository.On("GetBrand", mock.Anything, filter).Return(mockBrand, nil)
		res := brandService.CheckBrandById(context.TODO(), 1)

		assert.Equal(t, http.StatusNotFound, res.StatusCode)
		assert.False(t, res.Success)
		assert.Equal(t, "Brand Id doesn't exists", res.Message)
		assert.Nil(t, res.Data)
	})
}

func TestBrandCreate(t *testing.T) {
	var mockBrand []model.Brand
	mockRepository := new(mockRepositories.BrandRepository)

	reset := func() {
		mockBrand = []model.Brand{}
		mockRepository = new(mockRepositories.BrandRepository)
	}

	t.Run("Test Brand Create Success", func(t *testing.T) {
		defer reset()

		brandService := service.NewBrandService(mockRepository, contextTimeout)

		payload := dto.InsertBrandDto{}
		err := faker.FakeData(&payload)
		assert.NoError(t, err)

		filter := dto.FilterBrandDto{Title: payload.Title, Limit: 1}
		mockRepository.On("GetBrand", mock.Anything, filter).Return(mockBrand, nil)
		mockRepository.On("Create", mock.Anything, payload).Return(&model.Brand{ID: 1}, nil)

		res := brandService.Create(context.Background(), payload)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.True(t, res.Success)
		assert.Equal(t, "success", res.Message)
		assert.Equal(t, map[string]interface{}{"id": 1}, res.Data)
	})

	t.Run("Test Brand Create Already Exists", func(t *testing.T) {
		defer reset()

		err := faker.FakeData(&mockBrand)
		assert.NoError(t, err)

		payload := dto.InsertBrandDto{}
		err = faker.FakeData(&payload)
		assert.NoError(t, err)

		filter := dto.FilterBrandDto{Title: payload.Title, Limit: 1}
		mockRepository.On("GetBrand", mock.Anything, filter).Return(mockBrand, nil)
		mockRepository.On("Create", mock.Anything, payload).Return(&model.Brand{ID: 1}, nil)

		brandService := service.NewBrandService(mockRepository, contextTimeout)
		res := brandService.Create(context.Background(), payload)

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		assert.False(t, res.Success)
		assert.Equal(t, "Brand Already Exists", res.Message)
		assert.Nil(t, res.Data)
	})

	t.Run("Test Brand Create Error DB Get Brand", func(t *testing.T) {
		defer reset()

		err := faker.FakeData(&mockBrand)
		assert.NoError(t, err)

		payload := dto.InsertBrandDto{}
		err = faker.FakeData(&payload)
		assert.NoError(t, err)

		filter := dto.FilterBrandDto{Title: payload.Title, Limit: 1}
		mockRepository.On("GetBrand", mock.Anything, filter).Return(mockBrand, errors.New("Database Error"))
		mockRepository.On("Create", mock.Anything, payload).Return(&model.Brand{ID: 1}, nil)

		brandService := service.NewBrandService(mockRepository, contextTimeout)
		res := brandService.Create(context.Background(), payload)

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
		assert.False(t, res.Success)
		assert.Equal(t, "Database Error", res.Message)
		assert.Nil(t, res.Data)
	})

	t.Run("Test Brand Create Error DB Create Brand", func(t *testing.T) {
		defer reset()

		payload := dto.InsertBrandDto{}
		err := faker.FakeData(&payload)
		assert.NoError(t, err)

		filter := dto.FilterBrandDto{Title: payload.Title, Limit: 1}
		mockRepository.On("GetBrand", mock.Anything, filter).Return(mockBrand, nil)
		mockRepository.On("Create", mock.Anything, payload).Return(nil, errors.New("Database Error"))

		brandService := service.NewBrandService(mockRepository, contextTimeout)
		res := brandService.Create(context.Background(), payload)

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
		assert.False(t, res.Success)
		assert.Equal(t, "Database Error", res.Message)
		assert.Nil(t, res.Data)
	})

}
