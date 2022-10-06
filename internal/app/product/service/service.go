package service

import (
	"context"
	"net/http"
	"time"

	brandService "github.com/ranggabudipangestu/simple-ecommerce/internal/app/brand/service"
	"github.com/ranggabudipangestu/simple-ecommerce/internal/app/product/dto"
	"github.com/ranggabudipangestu/simple-ecommerce/internal/app/product/repository"
	"github.com/ranggabudipangestu/simple-ecommerce/pkg/util"
)

type ProductService interface {
	Create(ctx context.Context, dto dto.InsertProductDto) (res *util.Response)
	GetProductById(ctx context.Context, id int) (res *util.Response)
	GetProductByBrand(ctx context.Context, brandId int) (res *util.Response)
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

func (s *Service) Create(ctx context.Context, payload dto.InsertProductDto) (res *util.Response) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	//check brandId exists or not
	brandResult := s.brandService.CheckBrandById(ctx, payload.BrandId)
	if !brandResult.Success {
		return res.ReturnedData(false, brandResult.StatusCode, brandResult.Message, nil)
	}
	result, err := s.productRepository.Create(ctx, payload)
	if err != nil {
		return res.ReturnedData(false, http.StatusBadRequest, err.Error(), nil)
	}

	return res.ReturnedData(true, http.StatusOK, "success", map[string]interface{}{"id": result.ID})
}

func (s *Service) GetProductById(ctx context.Context, id int) (res *util.Response) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	filter := dto.FilterProductDto{ID: id, Limit: 1}

	result, err := s.productRepository.GetProduct(ctx, filter)
	if err != nil {
		return res.ReturnedData(false, http.StatusInternalServerError, err.Error(), nil)
	}

	var data *dto.GetProduct = nil
	if len(result) == 0 {
		return res.ReturnedData(true, http.StatusNotFound, "Product Data Not Found", data)

	}

	data = &result[0]
	return res.ReturnedData(true, http.StatusOK, "success", data)
}

func (s *Service) GetProductByBrand(ctx context.Context, brandId int) (res *util.Response) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	filter := dto.FilterProductDto{BrandId: brandId}

	result, err := s.productRepository.GetProduct(ctx, filter)
	if err != nil {
		return res.ReturnedData(false, http.StatusBadRequest, err.Error(), nil)
	}

	if len(result) == 0 {
		return res.ReturnedData(true, http.StatusNotFound, "Product Data Not Found", result)
	}

	return res.ReturnedData(true, http.StatusOK, "success", result)
}
