package service

import (
	"context"
	"net/http"
	"time"

	"github.com/ranggabudipangestu/simple-ecommerce/internal/app/brand/dto"
	"github.com/ranggabudipangestu/simple-ecommerce/internal/app/brand/repository"
	"github.com/ranggabudipangestu/simple-ecommerce/pkg/util"
)

type BrandService interface {
	Create(ctx context.Context, payload dto.InsertBrandDto) (res *util.Response)
	CheckBrandById(ctx context.Context, id int) (res *util.Response)
}

type Service struct {
	brandRepository repository.BrandRepository
	contextTimeout  time.Duration
}

func NewBrandService(r repository.BrandRepository, timeout time.Duration) BrandService {
	return &Service{
		brandRepository: r,
		contextTimeout:  timeout,
	}
}

func (s *Service) Create(ctx context.Context, payload dto.InsertBrandDto) (res *util.Response) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	//Check Brand By Title
	filter := dto.FilterBrandDto{Title: payload.Title, Limit: 1}

	brand, err := s.brandRepository.GetBrand(ctx, filter)
	if err != nil {
		return res.ReturnedData(false, http.StatusInternalServerError, err.Error(), nil)
	}

	//if brand exists
	if len(brand) > 0 {
		return res.ReturnedData(false, http.StatusBadRequest, "Brand Already Exists", nil)
	}

	//Create Brand
	result, err := s.brandRepository.Create(ctx, payload)
	if err != nil {
		return res.ReturnedData(false, http.StatusInternalServerError, err.Error(), nil)
	}

	return res.ReturnedData(true, http.StatusOK, "success", map[string]interface{}{"id": result.ID})
}

func (s *Service) CheckBrandById(ctx context.Context, id int) (res *util.Response) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	//Check Brand By Title
	filter := dto.FilterBrandDto{ID: id, Limit: 1}

	brand, err := s.brandRepository.GetBrand(ctx, filter)
	if err != nil {
		return res.ReturnedData(false, http.StatusInternalServerError, err.Error(), nil)
	}

	//if brand exists
	if len(brand) == 0 {
		return res.ReturnedData(false, http.StatusNotFound, "Brand Id doesn't exists", nil)
	}

	return res.ReturnedData(true, http.StatusOK, "success", brand[0])
}
