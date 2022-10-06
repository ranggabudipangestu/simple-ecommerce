package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/ranggabudipangestu/simple-ecommerce/internal/app/brand/dto"
	"github.com/ranggabudipangestu/simple-ecommerce/internal/app/brand/service"
	"github.com/ranggabudipangestu/simple-ecommerce/pkg/util"
)

type BrandHandler struct {
	brandService service.BrandService
}

func NewBrandHandlers(mux *http.ServeMux, service service.BrandService) {
	handler := BrandHandler{brandService: service}
	mux.HandleFunc("/brand", func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case "POST":
			handler.Create(w, r)
		}

	})

}

func (b *BrandHandler) Create(w http.ResponseWriter, r *http.Request) {
	var res *util.Response

	var payload dto.InsertBrandDto
	err := json.NewDecoder(r.Body).Decode(&payload)

	if err != nil {
		res.JSON(w, http.StatusInternalServerError, res.ReturnedData(false, http.StatusInternalServerError, err.Error(), nil))
		return
	}

	var valid bool
	if valid, err = isRequestValid(&payload); !valid {
		res.JSON(w, http.StatusInternalServerError, res.ReturnedData(false, http.StatusBadRequest, err.Error(), nil))
		return
	}

	result := b.brandService.Create(r.Context(), payload)
	res.JSON(w, result.StatusCode, result)
	return
}

func isRequestValid(payload *dto.InsertBrandDto) (bool, error) {
	validate := validator.New()
	err := validate.Struct(payload)
	if err != nil {
		return false, err
	}
	return true, nil
}
