package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	validator "github.com/go-playground/validator/v10"
	"github.com/ranggabudipangestu/simple-ecommerce/internal/app/product/dto"
	"github.com/ranggabudipangestu/simple-ecommerce/internal/app/product/service"
	"github.com/ranggabudipangestu/simple-ecommerce/pkg/util"
)

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(mux *http.ServeMux, service service.ProductService) {
	handler := ProductHandler{productService: service}

	mux.HandleFunc("/product", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			handler.Create(w, r)
		case "GET":
			handler.GetProductById(w, r)
		}
	})

	mux.HandleFunc("/product/brand", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handler.GetProductByBrand(w, r)
		}
	})

}

func (b *ProductHandler) Create(w http.ResponseWriter, r *http.Request) error {
	var res *util.Response

	var payload dto.InsertProductDto
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		return res.JSON(w, false, util.GetResCode("SYSTEM_ERROR"), err.Error(), nil)
	}

	var valid bool
	if valid, err = isRequestValid(&payload); !valid {
		return res.JSON(w, false, util.GetResCode("VALIDATION_ERROR"), err.Error(), nil)

	}

	result, err, state := b.productService.Create(r.Context(), payload)
	if err != nil {
		return res.JSON(w, false, util.GetResCode(state), err.Error(), result)
	}
	return res.JSON(w, true, util.GetResCode(state), "Success", result)
}

func (b *ProductHandler) GetProductById(w http.ResponseWriter, r *http.Request) error {
	var res *util.Response
	ParamId := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(ParamId)
	result, err, state := b.productService.GetProductById(r.Context(), id)

	if err != nil {
		return res.JSON(w, false, util.GetResCode(state), err.Error(), result)
	}
	return res.JSON(w, true, util.GetResCode(state), "Success", result)
}

func (b *ProductHandler) GetProductByBrand(w http.ResponseWriter, r *http.Request) error {
	var res *util.Response
	ParamId := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(ParamId)
	result, err, state := b.productService.GetProductByBrand(r.Context(), id)

	if err != nil {
		return res.JSON(w, false, util.GetResCode(state), err.Error(), result)
	}
	return res.JSON(w, false, util.GetResCode(state), "Success", result)
}

func isRequestValid(dto *dto.InsertProductDto) (bool, error) {
	validate := validator.New()
	err := validate.Struct(dto)
	if err != nil {
		return false, err
	}
	return true, nil
}
