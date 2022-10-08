package service

import (
	"context"
	"errors"
	"time"

	"github.com/ranggabudipangestu/simple-ecommerce/internal/app/brand/dto"
	"github.com/ranggabudipangestu/simple-ecommerce/internal/app/brand/repository"
	"github.com/ranggabudipangestu/simple-ecommerce/internal/model"
	"github.com/ranggabudipangestu/simple-ecommerce/pkg/util"
)

type BrandService interface {
	Create(ctx context.Context, payload dto.InsertBrandDto) (interface{}, error, string)
	CheckBrandById(ctx context.Context, id int) (*model.Brand, error, string)
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

func (s *Service) Create(ctx context.Context, payload dto.InsertBrandDto) (interface{}, error, string) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	//Check Brand By Title
	filter := dto.FilterBrandDto{Title: payload.Title, Limit: 1}

	brand, err := s.brandRepository.GetBrand(ctx, filter)
	if err != nil {
		return nil, err, util.SYSTEM_ERROR
	}

	//if brand exists
	if len(brand) > 0 {
		return nil, errors.New("Brand title already Exists"), "DUPLICATE"
	}

	//Create Brand
	result, err := s.brandRepository.Create(ctx, payload)
	if err != nil {
		return nil, err, util.SYSTEM_ERROR
	}

	return map[string]interface{}{"id": result.ID}, nil, util.SUCCESS
}

func (s *Service) CheckBrandById(ctx context.Context, id int) (*model.Brand, error, string) {
	ctx, cancel := context.WithTimeout(ctx, s.contextTimeout)
	defer cancel()

	//Check Brand By Title
	filter := dto.FilterBrandDto{ID: id, Limit: 1}

	brand, err := s.brandRepository.GetBrand(ctx, filter)
	if err != nil {
		return nil, err, util.SYSTEM_ERROR
	}

	//if brand exists
	if len(brand) == 0 {
		return nil, errors.New("Brand Id Doesn't exists"), "NOT_FOUND"
	}

	return &brand[0], nil, util.SUCCESS
}
