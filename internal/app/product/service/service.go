package service

import (
	"context"
	"errors"
	"time"

	brandService "github.com/ranggabudipangestu/simple-ecommerce/internal/app/brand/service"
	"github.com/ranggabudipangestu/simple-ecommerce/internal/app/product/dto"
	"github.com/ranggabudipangestu/simple-ecommerce/internal/app/product/repository"
	"github.com/ranggabudipangestu/simple-ecommerce/pkg/util"
)

type ProductService interface {
	Create(ctx context.Context, dto dto.InsertProductDto) (interface{}, error, string)
	GetProductById(ctx context.Context, id int) (*dto.GetProduct, error, string)
	GetProductByBrand(ctx context.Context, brandId int) (interface{}, error, string)
}

type Service struct {
	productRepository repository.ProductRepository
	brandService      brandService.BrandService
	contextTimeout    time.Duration
}

func NewProductService(repository repository.ProductRepository, brandService brandService.BrandService, timeout time.Duration) ProductService {
	return &Service{
		productRepository: repository,
		brandService:      brandService,
		contextTimeout:    timeout,
	}
}

func (s *Service) Create(ctx context.Context, payload dto.InsertProductDto) (interface{}, error, string) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	//check brandId exists or not
	_, err, state := s.brandService.CheckBrandById(ctx, payload.BrandId)
	if err != nil {
		return nil, err, state
	}

	result, err := s.productRepository.Create(ctx, payload)
	if err != nil {
		return nil, err, "SYSTEM_ERROR"
	}
	return map[string]interface{}{"id": result.ID}, nil, util.SUCCESS
}

func (s *Service) GetProductById(ctx context.Context, id int) (*dto.GetProduct, error, string) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	filter := dto.FilterProductDto{ID: id, Limit: 1}

	result, err := s.productRepository.GetProduct(ctx, filter)
	if err != nil {
		return nil, err, util.SYSTEM_ERROR
	}

	var data *dto.GetProduct = nil
	if len(result) == 0 {
		return nil, errors.New("Product Not Found"), util.NOT_FOUND

	}

	data = &result[0]
	return data, nil, util.SUCCESS
}

func (s *Service) GetProductByBrand(ctx context.Context, brandId int) (interface{}, error, string) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	filter := dto.FilterProductDto{BrandId: brandId}

	result, err := s.productRepository.GetProduct(ctx, filter)
	if err != nil {
		return nil, err, util.SYSTEM_ERROR
	}

	if len(result) == 0 {
		return nil, errors.New("Product Not Found"), util.NOT_FOUND
	}

	return result, nil, util.SUCCESS
}
