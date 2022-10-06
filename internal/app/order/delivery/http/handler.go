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
	orderService service.OrderService
}

func NewOrderHandler(mux *http.ServeMux, service service.OrderService) {
	handler := OrderHandler{orderService: service}

	mux.HandleFunc("/order", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":

			handler.CreateOrder(w, r)
		case "GET":
			handler.GetOrderDetails(w, r)
		}
	})

}

func (b *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var res *util.Response

	var payload dto.CreateOrderDto
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

	for i, detail := range payload.Details {
		if valid, err = isDetailsRequestIsValid(&detail); !valid {
			errMessage := fmt.Sprintf("Error row %s with details %s", strconv.Itoa(i+1), err.Error())
			res.JSON(w, http.StatusInternalServerError, res.ReturnedData(false, http.StatusBadRequest, errMessage, nil))
			return
		}
	}
	result := b.orderService.CreateOrder(r.Context(), payload)

	res.JSON(w, result.StatusCode, result)
	return
}

func (b *OrderHandler) GetOrderDetails(w http.ResponseWriter, r *http.Request) {
	var res *util.Response

	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	result := b.orderService.GetOrderDetails(r.Context(), id)

	res.JSON(w, result.StatusCode, result)
	return
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
