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
	BrandService service.BrandService
}

func NewBrandHandlers(mux *http.ServeMux, service service.BrandService) {
	handler := BrandHandler{BrandService: service}
	mux.HandleFunc("/brand", func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case "POST":
			handler.Create(w, r)
		}
	})

}

func (b *BrandHandler) Create(w http.ResponseWriter, r *http.Request) error {
	var res *util.Response

	var payload dto.InsertBrandDto
	err := json.NewDecoder(r.Body).Decode(&payload)

	if err != nil {
		code := util.GetResCode(util.SYSTEM_ERROR)
		return res.JSON(w, false, code, err.Error(), nil)
	}

	var valid bool
	if valid, err = isRequestValid(&payload); !valid {
		return res.JSON(w, false, util.GetResCode(util.VALIDATION_ERROR), err.Error(), nil)
	}

	result, err, state := b.BrandService.Create(r.Context(), payload)
	if err != nil {
		return res.JSON(w, false, util.GetResCode(state), err.Error(), nil)
	}
	return res.JSON(w, true, http.StatusOK, "success", result)
}

func isRequestValid(payload *dto.InsertBrandDto) (bool, error) {
	validate := validator.New()
	err := validate.Struct(payload)
	if err != nil {
		return false, err
	}
	return true, nil
}
