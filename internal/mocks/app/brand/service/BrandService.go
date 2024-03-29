// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	context "context"

	dto "github.com/ranggabudipangestu/simple-ecommerce/internal/app/brand/dto"
	mock "github.com/stretchr/testify/mock"

	model "github.com/ranggabudipangestu/simple-ecommerce/internal/model"
)

// BrandService is an autogenerated mock type for the BrandService type
type BrandService struct {
	mock.Mock
}

// CheckBrandById provides a mock function with given fields: ctx, id
func (_m *BrandService) CheckBrandById(ctx context.Context, id int) (*model.Brand, error, string) {
	ret := _m.Called(ctx, id)

	var r0 *model.Brand
	if rf, ok := ret.Get(0).(func(context.Context, int) *model.Brand); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Brand)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	var r2 string
	if rf, ok := ret.Get(2).(func(context.Context, int) string); ok {
		r2 = rf(ctx, id)
	} else {
		r2 = ret.Get(2).(string)
	}

	return r0, r1, r2
}

// Create provides a mock function with given fields: ctx, payload
func (_m *BrandService) Create(ctx context.Context, payload dto.InsertBrandDto) (interface{}, error, string) {
	ret := _m.Called(ctx, payload)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(context.Context, dto.InsertBrandDto) interface{}); ok {
		r0 = rf(ctx, payload)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, dto.InsertBrandDto) error); ok {
		r1 = rf(ctx, payload)
	} else {
		r1 = ret.Error(1)
	}

	var r2 string
	if rf, ok := ret.Get(2).(func(context.Context, dto.InsertBrandDto) string); ok {
		r2 = rf(ctx, payload)
	} else {
		r2 = ret.Get(2).(string)
	}

	return r0, r1, r2
}

type mockConstructorTestingTNewBrandService interface {
	mock.TestingT
	Cleanup(func())
}

// NewBrandService creates a new instance of BrandService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewBrandService(t mockConstructorTestingTNewBrandService) *BrandService {
	mock := &BrandService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
