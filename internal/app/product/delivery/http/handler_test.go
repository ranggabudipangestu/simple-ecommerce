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

	productHttp "github.com/ranggabudipangestu/simple-ecommerce/internal/app/product/delivery/http"
	"github.com/ranggabudipangestu/simple-ecommerce/internal/app/product/dto"
	mocks "github.com/ranggabudipangestu/simple-ecommerce/internal/mocks/app/product/service"
)

func TestCreateProduct(t *testing.T) {
	mux := http.NewServeMux()

	payload := dto.InsertProductDto{
		Title:       "Nike Airmax",
		Description: "The first edition of nike",
		BrandId:     1,
		Price:       1000000,
	}

	j, err := json.Marshal(payload)
	assert.NoError(t, err)

	mockService := new(mocks.ProductService)

	reset := func() {
		mux = http.NewServeMux()
		mockService = new(mocks.ProductService)
	}

	t.Run("Test Create Product success", func(t *testing.T) {
		defer reset()

		mockService.On("Create", context.TODO(), payload).Return(map[string]interface{}{"id": 1}, nil, "SUCCESS")

		productHttp.NewProductHandler(mux, mockService)
		handler := productHttp.ProductHandler{ProductService: mockService}

		req := httptest.NewRequest(http.MethodPost, "/product", strings.NewReader(string(j)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		err = handler.Create(w, req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, w.Code)

	})

	t.Run("Test Create Product Failed Brand Id Not Found", func(t *testing.T) {
		defer reset()

		mockService.On("Create", context.TODO(), payload).Return(nil, errors.New("BrandId Doesn't exist"), "NOT_FOUND")

		productHttp.NewProductHandler(mux, mockService)
		handler := productHttp.ProductHandler{ProductService: mockService}

		req := httptest.NewRequest(http.MethodPost, "/product", strings.NewReader(string(j)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		err = handler.Create(w, req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, w.Code)

	})

	t.Run("Test Create Product Error from Database", func(t *testing.T) {
		defer reset()

		mockService.On("Create", context.TODO(), payload).Return(nil, errors.New("BrandId Doesn't exist"), "SYSTEM_ERROR")

		productHttp.NewProductHandler(mux, mockService)
		handler := productHttp.ProductHandler{ProductService: mockService}

		req := httptest.NewRequest(http.MethodPost, "/product", strings.NewReader(string(j)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		err = handler.Create(w, req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, w.Code)

	})

}

func TestGetProductById(t *testing.T) {
	mux := http.NewServeMux()
	var mockGetProduct *dto.GetProduct

	mockService := new(mocks.ProductService)

	reset := func() {
		mux = http.NewServeMux()
		mockService = new(mocks.ProductService)
	}

	t.Run("Test Get Product By Id Success", func(t *testing.T) {
		defer reset()
		err := faker.FakeData(&mockGetProduct)
		mockService.On("GetProductById", context.TODO(), 1).Return(mockGetProduct, nil, "SUCCESS")

		productHttp.NewProductHandler(mux, mockService)
		handler := productHttp.ProductHandler{ProductService: mockService}

		req := httptest.NewRequest(http.MethodGet, "/product?id=1", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		err = handler.GetProductById(w, req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, w.Code)

	})

	t.Run("Test Get Product By Id Not Found ", func(t *testing.T) {
		defer reset()
		mockService.On("GetProductById", context.TODO(), 1).Return(nil, errors.New("Product Not Found"), "NOT_FOUND")

		productHttp.NewProductHandler(mux, mockService)
		handler := productHttp.ProductHandler{ProductService: mockService}

		req := httptest.NewRequest(http.MethodGet, "/product?id=1", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		err := handler.GetProductById(w, req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, w.Code)

	})

	t.Run("Test Get Product By Id Error System ", func(t *testing.T) {
		defer reset()
		mockService.On("GetProductById", context.TODO(), 1).Return(nil, errors.New("System Error"), "SYSTEM_ERROR")

		productHttp.NewProductHandler(mux, mockService)
		handler := productHttp.ProductHandler{ProductService: mockService}

		req := httptest.NewRequest(http.MethodGet, "/product?id=1", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		err := handler.GetProductById(w, req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, w.Code)

	})

}

func TestGetProductByBrand(t *testing.T) {
	mux := http.NewServeMux()
	var mockGetProduct []dto.GetProduct

	mockService := new(mocks.ProductService)

	reset := func() {
		mux = http.NewServeMux()
		mockService = new(mocks.ProductService)
	}

	t.Run("Test Get Product By Brand Success", func(t *testing.T) {
		defer reset()
		err := faker.FakeData(&mockGetProduct)
		mockService.On("GetProductByBrand", context.TODO(), 1).Return(mockGetProduct, nil, "SUCCESS")

		productHttp.NewProductHandler(mux, mockService)
		handler := productHttp.ProductHandler{ProductService: mockService}

		req := httptest.NewRequest(http.MethodGet, "/product/brand?id=1", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		err = handler.GetProductByBrand(w, req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, w.Code)

	})

	t.Run("Test Get Product By Brand Not Found ", func(t *testing.T) {
		defer reset()
		mockService.On("GetProductByBrand", context.TODO(), 1).Return(nil, errors.New("Product Not Found"), "NOT_FOUND")

		productHttp.NewProductHandler(mux, mockService)
		handler := productHttp.ProductHandler{ProductService: mockService}

		req := httptest.NewRequest(http.MethodGet, "/product/brand?id=1", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		err := handler.GetProductByBrand(w, req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusNotFound, w.Code)

	})

	t.Run("Test Get Product By Brand Error System ", func(t *testing.T) {
		defer reset()
		mockService.On("GetProductByBrand", context.TODO(), 1).Return(nil, errors.New("System Error"), "SYSTEM_ERROR")

		productHttp.NewProductHandler(mux, mockService)
		handler := productHttp.ProductHandler{ProductService: mockService}

		req := httptest.NewRequest(http.MethodGet, "/product/brand?id=1", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		err := handler.GetProductByBrand(w, req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, w.Code)

	})

}
