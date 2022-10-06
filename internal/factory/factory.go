package factory

import (
	"database/sql"
	"net/http"
	"time"

	brandHandler "github.com/ranggabudipangestu/simple-ecommerce/internal/app/brand/delivery/http"
	BrandRepository "github.com/ranggabudipangestu/simple-ecommerce/internal/app/brand/repository"
	BrandService "github.com/ranggabudipangestu/simple-ecommerce/internal/app/brand/service"

	productHandler "github.com/ranggabudipangestu/simple-ecommerce/internal/app/product/delivery/http"
	ProductRepository "github.com/ranggabudipangestu/simple-ecommerce/internal/app/product/repository"
	ProductService "github.com/ranggabudipangestu/simple-ecommerce/internal/app/product/service"

	orderHandler "github.com/ranggabudipangestu/simple-ecommerce/internal/app/order/delivery/http"
	OrderRepository "github.com/ranggabudipangestu/simple-ecommerce/internal/app/order/repository"
	OrderService "github.com/ranggabudipangestu/simple-ecommerce/internal/app/order/service"
)

func RegisterHandlers(mux *http.ServeMux, db *sql.DB) {

	const contextTimeout = 2 * time.Second

	brandRepository := BrandRepository.NewBrand(db)
	brandService := BrandService.NewBrandService(brandRepository, contextTimeout)
	brandHandler.NewBrandHandlers(mux, brandService)

	productRepository := ProductRepository.NewProduct(db)
	productService := ProductService.NewProductService(productRepository, brandService, contextTimeout)
	productHandler.NewProductHandler(mux, productService)

	orderRepository := OrderRepository.NewOrder(db)
	orderService := OrderService.NewOrderService(orderRepository, productService, contextTimeout)
	orderHandler.NewOrderHandler(mux, orderService)

}
