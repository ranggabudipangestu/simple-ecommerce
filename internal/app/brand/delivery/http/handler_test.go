package http_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	brandHttp "github.com/ranggabudipangestu/simple-ecommerce/internal/app/brand/delivery/http"
	"github.com/ranggabudipangestu/simple-ecommerce/internal/app/brand/dto"
	mocks "github.com/ranggabudipangestu/simple-ecommerce/internal/mocks/app/brand/service"
)

func TestCreateBrand(t *testing.T) {
	mux := http.NewServeMux()

	payload := dto.InsertBrandDto{
		Title: "Nike",
	}
	j, err := json.Marshal(payload)
	assert.NoError(t, err)

	mockService := new(mocks.BrandService)

	reset := func() {
		mux = http.NewServeMux()
		mockService = new(mocks.BrandService)
	}

	t.Run("Test Create Brand Success", func(t *testing.T) {
		defer reset()
		mockService.On("Create", context.TODO(), payload).Return(map[string]interface{}{"id": 1}, nil, "SUCCESS")

		brandHttp.NewBrandHandlers(mux, mockService)
		handler := brandHttp.BrandHandler{
			BrandService: mockService,
		}

		req := httptest.NewRequest(http.MethodPost, "/brand", strings.NewReader(string(j)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		err = handler.Create(w, req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, w.Code)

	})

	t.Run("Test Create Brand Duplicate", func(t *testing.T) {
		defer reset()
		mockService.On("Create", context.TODO(), payload).Return(nil, errors.New("DUPLICATE"), "DUPLICATE")

		brandHttp.NewBrandHandlers(mux, mockService)
		handler := brandHttp.BrandHandler{
			BrandService: mockService,
		}

		req := httptest.NewRequest(http.MethodPost, "/brand", strings.NewReader(string(j)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		err = handler.Create(w, req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, w.Code)

	})
	t.Run("Test Create Product No Payload", func(t *testing.T) {
		defer reset()

		brandHttp.NewBrandHandlers(mux, mockService)
		handler := brandHttp.BrandHandler{BrandService: mockService}

		req := httptest.NewRequest(http.MethodPost, "/brand", nil)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		err = handler.Create(w, req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusInternalServerError, w.Code)

	})

	t.Run("Test Create Product Validation Body", func(t *testing.T) {
		defer reset()

		payload := dto.InsertBrandDto{
			Title: "",
		}

		j, err = json.Marshal(payload)
		assert.NoError(t, err)

		brandHttp.NewBrandHandlers(mux, mockService)
		handler := brandHttp.BrandHandler{BrandService: mockService}

		req := httptest.NewRequest(http.MethodPost, "/brand", strings.NewReader(string(j)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		err = handler.Create(w, req)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, w.Code)

	})
}
