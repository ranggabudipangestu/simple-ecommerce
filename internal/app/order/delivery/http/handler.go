package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	validator "github.com/go-playground/validator/v10"
	"github.com/ranggabudipangestu/simple-ecommerce/internal/app/order/dto"
	"github.com/ranggabudipangestu/simple-ecommerce/internal/app/order/service"
	"github.com/ranggabudipangestu/simple-ecommerce/pkg/util"
)

type OrderHandler struct {
	OrderService service.OrderService
}

func NewOrderHandler(mux *http.ServeMux, service service.OrderService) {
	handler := OrderHandler{OrderService: service}

	mux.HandleFunc("/order", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":

			handler.CreateOrder(w, r)
		case "GET":
			handler.GetOrderDetails(w, r)
		}
	})

}

func (b *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) error {
	var res *util.Response

	var payload dto.CreateOrderDto
	err := json.NewDecoder(r.Body).Decode(&payload)

	if err != nil {
		return res.JSON(w, false, util.GetResCode(util.SYSTEM_ERROR), err.Error(), nil)
	}

	var valid bool
	if valid, err = isRequestValid(&payload); !valid {
		return res.JSON(w, false, util.GetResCode(util.VALIDATION_ERROR), err.Error(), nil)
	}

	for i, detail := range payload.Details {
		if valid, err = isDetailsRequestIsValid(&detail); !valid {
			errMessage := fmt.Sprintf("Error row %s with details %s", strconv.Itoa(i+1), err.Error())
			return res.JSON(w, false, util.GetResCode(util.VALIDATION_ERROR), errMessage, nil)
		}
	}
	result, err, state := b.OrderService.CreateOrder(r.Context(), payload)

	if err != nil {
		return res.JSON(w, false, util.GetResCode(state), err.Error(), nil)
	}
	return res.JSON(w, true, util.GetResCode(state), "success", result)
}

func (b *OrderHandler) GetOrderDetails(w http.ResponseWriter, r *http.Request) error {
	var res *util.Response

	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	result, err, state := b.OrderService.GetOrderDetails(r.Context(), id)

	if err != nil {
		return res.JSON(w, false, util.GetResCode(state), err.Error(), result)
	}
	return res.JSON(w, true, util.GetResCode(state), "success", result)
}

func isRequestValid(payload *dto.CreateOrderDto) (bool, error) {
	validate := validator.New()
	err := validate.Struct(payload)
	if err != nil {
		return false, err
	}
	return true, nil
}
func isDetailsRequestIsValid(detail *dto.CreateOrderDetails) (bool, error) {
	validate := validator.New()
	err := validate.Struct(detail)
	if err != nil {
		return false, err
	}
	return true, nil
}
