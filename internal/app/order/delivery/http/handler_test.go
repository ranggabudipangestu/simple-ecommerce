package http_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"

	orderHttp "github.com/ranggabudipangestu/simple-ecommerce/internal/app/order/delivery/http"
	"github.com/ranggabudipangestu/simple-ecommerce/internal/app/order/dto"
	mocks "github.com/ranggabudipangestu/simple-ecommerce/internal/mocks/app/order/service"
)

func TestCreateOrder(t *testing.T) {
	mux := http.NewServeMux()

	var orderDetail []dto.CreateOrderDetails
	orderDetail = append(orderDetail, dto.CreateOrderDetails{
		ProductId: 1,
		Qty:       2,
	})
	payload := dto.CreateOrderDto{
		DeliveryAddress: "Indonesia",
		Details:         orderDetail,
	}

	j, err := json.Marshal(payload)
	assert.NoError(t, err)

	mockService := new(mocks.OrderService)

	reset := func() {
		mux = http.NewServeMux()
		mockService = new(mocks.OrderService)
	}

	t.Run("Test Create Order Success", func(t *testing.T) {
		defer reset()
		mockService.On("CreateOrder", context.TODO(), payload).Return(map[string]interface{}{"id": 1, "transactionNumber": "TRX-21510002451122"}, nil, "SUCCESS")

		orderHttp.NewOrderHandler(mux, mockService)
		handler := orderHttp.OrderHandler{OrderService: mockService}

		req := httptest.NewRequest(http.MethodPost, "/order", strings.NewReader(string(j)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		err = handler.CreateOrder(w, req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, w.Code)

	})

	t.Run("Test Create Order Failed Product Not Found", func(t *testing.T) {
		defer reset()
		mockService.On("CreateOrder", context.TODO(), payload).Return(nil, errors.New("Product Not Found"), "NOT_FOUND")

		orderHttp.NewOrderHandler(mux, mockService)
		handler := orderHttp.OrderHandler{OrderService: mockService}

		req := httptest.NewRequest(http.MethodPost, "/order", strings.NewReader(string(j)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		err = handler.CreateOrder(w, req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, w.Code)

	})

	t.Run("Test Create Order Failed Error In Database", func(t *testing.T) {
		defer reset()
		mockService.On("CreateOrder", context.TODO(), payload).Return(nil, errors.New("DATABASE ERROR"), "SYSTEM_ERROR")

		orderHttp.NewOrderHandler(mux, mockService)
		handler := orderHttp.OrderHandler{OrderService: mockService}

		req := httptest.NewRequest(http.MethodPost, "/order", strings.NewReader(string(j)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		err = handler.CreateOrder(w, req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, w.Code)

	})

	t.Run("Test Create Failed No Payload", func(t *testing.T) {
		defer reset()

		orderHttp.NewOrderHandler(mux, mockService)
		handler := orderHttp.OrderHandler{OrderService: mockService}

		req := httptest.NewRequest(http.MethodPost, "/order", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		err = handler.CreateOrder(w, req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, w.Code)

	})

	t.Run("Test Create Order Failed Validation Body", func(t *testing.T) {
		defer reset()

		var newOrderDetail []dto.CreateOrderDetails
		newOrderDetail = append(newOrderDetail, dto.CreateOrderDetails{
			ProductId: 0,
			Qty:       0,
		})
		payload := dto.CreateOrderDto{
			DeliveryAddress: "",
			Details:         newOrderDetail,
		}

		j, err = json.Marshal(payload)
		assert.NoError(t, err)

		orderHttp.NewOrderHandler(mux, mockService)
		handler := orderHttp.OrderHandler{OrderService: mockService}

		req := httptest.NewRequest(http.MethodPost, "/order", strings.NewReader(string(j)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		err = handler.CreateOrder(w, req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, w.Code)

	})

	t.Run("Test Create Order Failed Validation Details", func(t *testing.T) {
		defer reset()

		var newOrderDetail []dto.CreateOrderDetails
		newOrderDetail = append(newOrderDetail, dto.CreateOrderDetails{
			ProductId: 0,
			Qty:       0,
		})
		payload := dto.CreateOrderDto{
			DeliveryAddress: "Garuda Street",
			Details:         newOrderDetail,
		}

		j, err = json.Marshal(payload)
		assert.NoError(t, err)

		orderHttp.NewOrderHandler(mux, mockService)
		handler := orderHttp.OrderHandler{OrderService: mockService}

		req := httptest.NewRequest(http.MethodPost, "/order", strings.NewReader(string(j)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		err = handler.CreateOrder(w, req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, w.Code)

	})
}

func TestGetOrderDetails(t *testing.T) {
	mux := http.NewServeMux()
	var mockGetOrder *dto.GetOrderDto

	mockService := new(mocks.OrderService)

	reset := func() {
		mux = http.NewServeMux()
		mockService = new(mocks.OrderService)
	}

	t.Run("Test Get Order Detail Success", func(t *testing.T) {
		defer reset()
		err := faker.FakeData(&mockGetOrder)
		mockService.On("GetOrderDetails", context.TODO(), 1).Return(mockGetOrder, nil, "SUCCESS")

		orderHttp.NewOrderHandler(mux, mockService)
		handler := orderHttp.OrderHandler{OrderService: mockService}

		req := httptest.NewRequest(http.MethodGet, "/order?id=1", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		err = handler.GetOrderDetails(w, req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, w.Code)

	})

	t.Run("Test Get Order Detail Failed Error Database", func(t *testing.T) {
		defer reset()
		err := faker.FakeData(&mockGetOrder)
		mockService.On("GetOrderDetails", context.TODO(), 1).Return(nil, errors.New("Databaser Error"), "SYSTEM_ERROR")

		orderHttp.NewOrderHandler(mux, mockService)
		handler := orderHttp.OrderHandler{OrderService: mockService}

		req := httptest.NewRequest(http.MethodGet, "/order?id=1", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		err = handler.GetOrderDetails(w, req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("Test Get Order Detail Data Not Found", func(t *testing.T) {
		defer reset()
		err := faker.FakeData(&mockGetOrder)
		mockService.On("GetOrderDetails", context.TODO(), 1).Return(nil, errors.New("Order Not Found"), "NOT_FOUND")

		orderHttp.NewOrderHandler(mux, mockService)
		handler := orderHttp.OrderHandler{OrderService: mockService}

		req := httptest.NewRequest(http.MethodGet, "/order?id=1", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		err = handler.GetOrderDetails(w, req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
