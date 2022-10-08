package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/ranggabudipangestu/simple-ecommerce/internal/app/order/dto"
	"github.com/ranggabudipangestu/simple-ecommerce/internal/app/order/helper"
	"github.com/ranggabudipangestu/simple-ecommerce/internal/app/order/repository"
	productService "github.com/ranggabudipangestu/simple-ecommerce/internal/app/product/service"
	"github.com/ranggabudipangestu/simple-ecommerce/internal/model"
	"github.com/ranggabudipangestu/simple-ecommerce/pkg/util"
)

type OrderService interface {
	CreateOrder(ctx context.Context, payload dto.CreateOrderDto) (interface{}, error, string)
	GetOrderDetails(ctx context.Context, id int) (*dto.GetOrderDto, error, string)
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

func (s *Service) CreateOrder(ctx context.Context, payload dto.CreateOrderDto) (interface{}, error, string) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	for i, detail := range payload.Details {
		productResult, err, state := s.productService.GetProductById(ctx, detail.ProductId)
		if err != nil {
			return nil, err, state
		}

		byteData, _ := json.Marshal(productResult)
		receivedProduct := &model.Product{}
		json.Unmarshal(byteData, receivedProduct)

		payload.Details[i].Price = receivedProduct.Price
		payload.Details[i].Total = float32(detail.Qty) * receivedProduct.Price
		payload.TotalTransaction += payload.Details[i].Total
		payload.TotalQty += detail.Qty
	}

	transactionNumber := helper.GenerateTransactionNumber()
	result, err := s.orderRepository.CreateOrder(ctx, payload, transactionNumber)
	if err != nil {
		return nil, err, util.SYSTEM_ERROR
	}

	return map[string]interface{}{"id": result.ID, "transactionNumber": transactionNumber}, nil, util.SUCCESS
}

func (s *Service) GetOrderDetails(ctx context.Context, id int) (*dto.GetOrderDto, error, string) {

	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	result, err := s.orderRepository.GetOrderDetails(ctx, id)

	if err != nil {
		return nil, err, util.SYSTEM_ERROR
	}
	if result == nil {
		return nil, errors.New("Order Not Found"), util.NOT_FOUND
	}
	return result, nil, util.SUCCESS
}
