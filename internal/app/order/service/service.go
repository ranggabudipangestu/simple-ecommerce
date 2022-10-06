package service

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ranggabudipangestu/simple-ecommerce/internal/app/order/dto"
	"github.com/ranggabudipangestu/simple-ecommerce/internal/app/order/helper"
	"github.com/ranggabudipangestu/simple-ecommerce/internal/app/order/repository"
	productService "github.com/ranggabudipangestu/simple-ecommerce/internal/app/product/service"
	"github.com/ranggabudipangestu/simple-ecommerce/internal/model"
	"github.com/ranggabudipangestu/simple-ecommerce/pkg/util"
)

type OrderService interface {
	CreateOrder(ctx context.Context, payload dto.CreateOrderDto) (res *util.Response)
	GetOrderDetails(ctx context.Context, id int) (res *util.Response)
}

type Service struct {
	orderRepository repository.OrderRepository
	productService  productService.ProductService
	contextTimeout  time.Duration
}

func NewOrderService(repository repository.OrderRepository, productService productService.ProductService, timeout time.Duration) OrderService {
	return &Service{
		orderRepository: repository,
		productService:  productService,
		contextTimeout:  timeout,
	}
}

func (s *Service) CreateOrder(ctx context.Context, payload dto.CreateOrderDto) (res *util.Response) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	payload.TransactionNumber = helper.GenerateTransactionNumber()
	for i, detail := range payload.Details {
		productResult := s.productService.GetProductById(ctx, detail.ProductId)
		if !productResult.Success {
			return res.ReturnedData(false, productResult.StatusCode, productResult.Message, nil)
		}

		byteData, _ := json.Marshal(productResult.Data)
		receivedProduct := &model.Product{}
		json.Unmarshal(byteData, receivedProduct)

		payload.Details[i].Price = receivedProduct.Price
		payload.Details[i].Total = float32(detail.Qty) * receivedProduct.Price
		payload.TotalTransaction += payload.Details[i].Total
		payload.TotalQty += detail.Qty
	}

	result, err := s.orderRepository.CreateOrder(ctx, payload)
	if err != nil {
		return res.ReturnedData(false, http.StatusInternalServerError, err.Error(), nil)
	}

	return res.ReturnedData(true, http.StatusOK, "success", map[string]interface{}{"id": result.ID, "transactionNumber": payload.TransactionNumber})
}

func (s *Service) GetOrderDetails(ctx context.Context, id int) (res *util.Response) {

	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	result, err := s.orderRepository.GetOrderDetails(ctx, id)

	if err != nil {
		code := 0
		if err.Error() == "Not Found" {
			code = http.StatusNotFound
		} else {
			code = http.StatusInternalServerError
		}
		return res.ReturnedData(false, code, err.Error(), nil)
	}
	return res.ReturnedData(true, http.StatusOK, "success", result)
}
