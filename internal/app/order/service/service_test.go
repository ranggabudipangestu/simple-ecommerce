package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	BrandService "github.com/ranggabudipangestu/simple-ecommerce/internal/app/brand/service"
	"github.com/ranggabudipangestu/simple-ecommerce/internal/app/order/dto"
	OrderService "github.com/ranggabudipangestu/simple-ecommerce/internal/app/order/service"
	ProductDto "github.com/ranggabudipangestu/simple-ecommerce/internal/app/product/dto"
	ProductService "github.com/ranggabudipangestu/simple-ecommerce/internal/app/product/service"
	mockBrandRepositores "github.com/ranggabudipangestu/simple-ecommerce/internal/mocks/app/brand/repository"
	mockOrderRepositories "github.com/ranggabudipangestu/simple-ecommerce/internal/mocks/app/order/repository"
	mockProductRepositores "github.com/ranggabudipangestu/simple-ecommerce/internal/mocks/app/product/repository"
	"github.com/ranggabudipangestu/simple-ecommerce/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const contextTimeout = 2 * time.Second

var mockProduct []ProductDto.GetProduct
var mockGetOrder *dto.GetOrderDto
var mockOrderRepository = new(mockOrderRepositories.OrderRepository)
var mockProductRepository = new(mockProductRepositores.ProductRepository)
var mockBrandRepository = new(mockBrandRepositores.BrandRepository)
var brandService = BrandService.NewBrandService(mockBrandRepository, contextTimeout)
var productService = ProductService.NewProductService(mockProductRepository, brandService, contextTimeout)
var orderService = OrderService.NewOrderService(mockOrderRepository, productService, contextTimeout)

func reset() {
	mockProduct = []ProductDto.GetProduct{}
	mockGetOrder = &dto.GetOrderDto{}
	mockOrderRepository = new(mockOrderRepositories.OrderRepository)
	mockProductRepository = new(mockProductRepositores.ProductRepository)
	mockBrandRepository = new(mockBrandRepositores.BrandRepository)

	brandService = BrandService.NewBrandService(mockBrandRepository, contextTimeout)
	productService = ProductService.NewProductService(mockProductRepository, brandService, contextTimeout)
	orderService = OrderService.NewOrderService(mockOrderRepository, productService, contextTimeout)

}

func TestCreateOrder(t *testing.T) {
	t.Run("Test Create Order Success", func(t *testing.T) {
		defer reset()

		mockProduct = append(mockProduct, ProductDto.GetProduct{
			ID:          1,
			Title:       "Nike",
			Description: "",
			Brand: ProductDto.BrandDto{
				ID: 1, Title: "Nike",
			},
			Price: 25000000,
		})

		var payloadDetail []dto.CreateOrderDetails
		payloadDetail = append(payloadDetail, dto.CreateOrderDetails{ProductId: 1})
		payload := dto.CreateOrderDto{
			Details:         payloadDetail,
			DeliveryAddress: "Indonesia",
		}

		mockProductRepository.On("GetProduct", mock.Anything, mock.Anything).Return(mockProduct, nil)
		mockOrderRepository.On("CreateOrder", mock.Anything, payload, mock.Anything).Return(&model.Transaction{ID: 1}, nil)

		res, err, state := orderService.CreateOrder(context.TODO(), payload)

		assert.Equal(t, "SUCCESS", state)
		assert.NotNil(t, res)
		assert.Nil(t, err)
	})

	t.Run("Test Create Order Product Id Doesn't Exist", func(t *testing.T) {
		defer reset()

		var payloadDetail []dto.CreateOrderDetails
		payloadDetail = append(payloadDetail, dto.CreateOrderDetails{ProductId: 1})
		payload := dto.CreateOrderDto{
			Details:         payloadDetail,
			DeliveryAddress: "Indonesia",
		}

		mockProductRepository.On("GetProduct", mock.Anything, mock.Anything).Return(mockProduct, nil)
		mockOrderRepository.On("CreateOrder", mock.Anything, payload, mock.Anything).Return(&model.Transaction{ID: 1}, nil)

		res, err, state := orderService.CreateOrder(context.TODO(), payload)

		assert.Equal(t, "NOT_FOUND", state)
		assert.NotNil(t, err)
		assert.Nil(t, res)
	})

	t.Run("Test Create Order Error in Database when Get Product Id", func(t *testing.T) {
		defer reset()

		var payloadDetail []dto.CreateOrderDetails
		payloadDetail = append(payloadDetail, dto.CreateOrderDetails{ProductId: 1})
		payload := dto.CreateOrderDto{
			Details:         payloadDetail,
			DeliveryAddress: "Indonesia",
		}

		mockProductRepository.On("GetProduct", mock.Anything, mock.Anything).Return(mockProduct, errors.New("Database Error"))
		mockOrderRepository.On("CreateOrder", mock.Anything, payload, mock.Anything).Return(&model.Transaction{ID: 1}, nil)

		res, err, state := orderService.CreateOrder(context.TODO(), payload)

		assert.Equal(t, "SYSTEM_ERROR", state)
		assert.NotNil(t, err)
		assert.Nil(t, res)
	})

	t.Run("Test Create Order Error in Database when Create Order", func(t *testing.T) {
		defer reset()

		mockProduct = append(mockProduct, ProductDto.GetProduct{
			ID:          1,
			Title:       "Nike",
			Description: "",
			Brand: ProductDto.BrandDto{
				ID: 1, Title: "Nike",
			},
			Price: 25000000,
		})

		var payloadDetail []dto.CreateOrderDetails
		payloadDetail = append(payloadDetail, dto.CreateOrderDetails{ProductId: 1})
		payload := dto.CreateOrderDto{
			Details:         payloadDetail,
			DeliveryAddress: "Indonesia",
		}

		mockProductRepository.On("GetProduct", mock.Anything, mock.Anything).Return(mockProduct, nil)
		mockOrderRepository.On("CreateOrder", mock.Anything, payload, mock.Anything).Return(nil, errors.New("Database Error"))

		res, err, state := orderService.CreateOrder(context.TODO(), payload)

		assert.Equal(t, "SYSTEM_ERROR", state)
		assert.NotNil(t, err)
		assert.Nil(t, res)
	})
}

func TestGetOrderDetails(t *testing.T) {
	t.Run("Test Get Order Detail Success", func(t *testing.T) {
		defer reset()

		err := faker.FakeData(&mockGetOrder)
		assert.NoError(t, err)

		mockOrderRepository.On("GetOrderDetails", mock.Anything, mockGetOrder.ID).Return(mockGetOrder, nil)

		res, err, state := orderService.GetOrderDetails(context.TODO(), mockGetOrder.ID)

		assert.Equal(t, "SUCCESS", state)
		assert.NotNil(t, res)
		assert.Nil(t, err)
	})
	t.Run("Test Get Order Failed Database Error", func(t *testing.T) {
		defer reset()
		mockOrderRepository.On("GetOrderDetails", mock.Anything, mockGetOrder.ID).Return(nil, errors.New("Database Error"))

		res, err, state := orderService.GetOrderDetails(context.TODO(), mockGetOrder.ID)

		assert.Equal(t, "SYSTEM_ERROR", state)
		assert.NotNil(t, err)
		assert.Nil(t, res)
	})
	t.Run("Test Get Order Data Not Found", func(t *testing.T) {
		defer reset()
		mockOrderRepository.On("GetOrderDetails", mock.Anything, mockGetOrder.ID).Return(nil, nil)

		res, err, state := orderService.GetOrderDetails(context.TODO(), mockGetOrder.ID)

		assert.Equal(t, "NOT_FOUND", state)
		assert.NotNil(t, err)
		assert.Nil(t, res)
	})
}
