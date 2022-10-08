package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	BrandService "github.com/ranggabudipangestu/simple-ecommerce/internal/app/brand/service"
	dto "github.com/ranggabudipangestu/simple-ecommerce/internal/app/product/dto"
	ProductService "github.com/ranggabudipangestu/simple-ecommerce/internal/app/product/service"
	mockBrandRepositores "github.com/ranggabudipangestu/simple-ecommerce/internal/mocks/app/brand/repository"
	mockProductRepositores "github.com/ranggabudipangestu/simple-ecommerce/internal/mocks/app/product/repository"
	"github.com/ranggabudipangestu/simple-ecommerce/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const contextTimeout = 2 * time.Second

var (
	mockBrand             []model.Brand
	mockProduct           []dto.GetProduct
	mockProductRepository = new(mockProductRepositores.ProductRepository)
	mockBrandRepository   = new(mockBrandRepositores.BrandRepository)
	brandService          = BrandService.NewBrandService(mockBrandRepository, contextTimeout)
	productService        = ProductService.NewProductService(mockProductRepository, brandService, contextTimeout)
)

func reset() {
	mockProduct = []dto.GetProduct{}
	mockBrand = []model.Brand{}
	mockProductRepository = new(mockProductRepositores.ProductRepository)
	mockBrandRepository = new(mockBrandRepositores.BrandRepository)

	brandService = BrandService.NewBrandService(mockBrandRepository, contextTimeout)
	productService = ProductService.NewProductService(mockProductRepository, brandService, contextTimeout)

}

func TestCreateProduct(t *testing.T) {
	t.Run("Test Create Product Success", func(t *testing.T) {
		defer reset()

		err := faker.FakeData(&mockBrand)
		assert.NoError(t, err)

		modelProduct := &model.Product{ID: 1}

		payload := dto.InsertProductDto{
			Title:       "Adidas",
			Description: "",
			BrandId:     mockBrand[0].ID,
			Price:       25000000,
		}

		mockBrandRepository.On("GetBrand", mock.Anything, mock.Anything).Return(mockBrand, nil)
		mockProductRepository.On("Create", mock.Anything, payload).Return(modelProduct, nil)

		res, err, state := productService.Create(context.TODO(), payload)

		assert.Equal(t, "SUCCESS", state)
		assert.NotNil(t, res)
		assert.Nil(t, err)
	})

	t.Run("Test Create Product Failed From Database", func(t *testing.T) {
		defer reset()

		err := faker.FakeData(&mockBrand)
		assert.NoError(t, err)

		payload := dto.InsertProductDto{
			Title:       "Adidas",
			Description: "",
			BrandId:     1,
			Price:       25000000,
		}

		mockBrandRepository.On("GetBrand", mock.Anything, mock.Anything).Return(mockBrand, nil)
		mockProductRepository.On("Create", mock.Anything, payload).Return(nil, errors.New("Database Error"))

		res, err, state := productService.Create(context.TODO(), payload)

		assert.Equal(t, "SYSTEM_ERROR", state)
		assert.NotNil(t, err)
		assert.Nil(t, res)
	})

	t.Run("Test Create Product Failed Brand Not Found", func(t *testing.T) {
		defer reset()

		modelProduct := &model.Product{ID: 1}

		payload := dto.InsertProductDto{
			Title:       "Adidas",
			Description: "",
			BrandId:     1,
			Price:       25000000,
		}

		mockBrandRepository.On("GetBrand", mock.Anything, mock.Anything).Return(mockBrand, nil)
		mockProductRepository.On("Create", mock.Anything, payload).Return(modelProduct, nil)

		res, err, state := productService.Create(context.TODO(), payload)

		assert.Equal(t, "NOT_FOUND", state)
		assert.NotNil(t, err)
		assert.Nil(t, res)
	})

}

func TestProductGetById(t *testing.T) {
	filter := dto.FilterProductDto{
		ID: 1, Limit: 1,
	}

	t.Run("Test Produt Get By ID Success", func(t *testing.T) {
		defer reset()

		err := faker.FakeData(&mockProduct)
		assert.NoError(t, err)

		filter := dto.FilterProductDto{
			ID: 1, Limit: 1,
		}

		mockProductRepository.On("GetProduct", mock.Anything, filter).Return(mockProduct, nil)

		res, err, state := productService.GetProductById(context.TODO(), filter.ID)

		assert.Equal(t, "SUCCESS", state)
		assert.NotNil(t, res)
		assert.Nil(t, err)
	})
	t.Run("Test Produt Get By ID Product Not Found", func(t *testing.T) {
		defer reset()

		mockProductRepository.On("GetProduct", mock.Anything, filter).Return(mockProduct, nil)

		res, err, state := productService.GetProductById(context.TODO(), filter.ID)

		assert.Equal(t, "NOT_FOUND", state)
		assert.NotNil(t, err)
		assert.Nil(t, res)
	})
	t.Run("Test Produt Get By ID Product Failed Database Error", func(t *testing.T) {
		defer reset()

		mockProductRepository.On("GetProduct", mock.Anything, filter).Return(mockProduct, errors.New("Database Error"))

		res, err, state := productService.GetProductById(context.TODO(), filter.ID)

		assert.Equal(t, "SYSTEM_ERROR", state)
		assert.NotNil(t, err)
		assert.Nil(t, res)
	})
}

func TestGetProductByBrand(t *testing.T) {
	filter := dto.FilterProductDto{BrandId: 1}

	t.Run("Test Get Product By Brand Success", func(t *testing.T) {
		defer reset()

		err := faker.FakeData(&mockProduct)
		assert.NoError(t, err)

		mockProductRepository.On("GetProduct", mock.Anything, filter).Return(mockProduct, nil)

		res, err, state := productService.GetProductByBrand(context.TODO(), filter.BrandId)

		assert.Equal(t, "SUCCESS", state)
		assert.NotNil(t, res)
		assert.Nil(t, err)
	})

	t.Run("Test Get Product By Brand Not Found", func(t *testing.T) {
		defer reset()

		mockProductRepository.On("GetProduct", mock.Anything, filter).Return(mockProduct, nil)

		res, err, state := productService.GetProductByBrand(context.TODO(), filter.BrandId)

		assert.Equal(t, "NOT_FOUND", state)
		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("Test Get Product By Brand Error Database", func(t *testing.T) {
		defer reset()

		mockProductRepository.On("GetProduct", mock.Anything, filter).Return(nil, errors.New("Database Error"))

		res, err, state := productService.GetProductByBrand(context.TODO(), filter.BrandId)

		assert.Equal(t, "SYSTEM_ERROR", state)
		assert.Nil(t, res)
		assert.NotNil(t, err)
	})
}
