package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ranggabudipangestu/simple-ecommerce/internal/app/brand/dto"
	mockRepositories "github.com/ranggabudipangestu/simple-ecommerce/internal/mocks/app/brand/repository"

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

		res, err, state := brandService.CheckBrandById(context.TODO(), 1)

		assert.Equal(t, "SUCCESS", state)
		assert.NotNil(t, res)
		assert.Nil(t, err)
	})

	t.Run("Test Brand By ID Database Error", func(t *testing.T) {
		defer reset()

		brandService := service.NewBrandService(mockRepository, contextTimeout)
		filter := dto.FilterBrandDto{ID: 1, Limit: 1}
		mockRepository.On("GetBrand", mock.Anything, filter).Return(mockBrand, errors.New("Database Error"))
		res, err, state := brandService.CheckBrandById(context.TODO(), 1)

		assert.Equal(t, "SYSTEM_ERROR", state)
		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Test Brand By ID Not Found", func(t *testing.T) {
		defer reset()

		brandService := service.NewBrandService(mockRepository, contextTimeout)
		filter := dto.FilterBrandDto{ID: 1, Limit: 1}
		mockRepository.On("GetBrand", mock.Anything, filter).Return(mockBrand, nil)
		res, err, state := brandService.CheckBrandById(context.TODO(), 1)

		assert.Equal(t, "NOT_FOUND", state)
		assert.Nil(t, res)
		assert.NotNil(t, err)
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

		res, err, state := brandService.Create(context.Background(), payload)

		assert.Equal(t, map[string]interface{}{"id": 1}, res)
		assert.Equal(t, "SUCCESS", state)
		assert.Nil(t, err)
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
		res, err, state := brandService.Create(context.Background(), payload)

		assert.Nil(t, res)
		assert.NotNil(t, err)
		assert.Equal(t, "DUPLICATE", state)
	})

}
