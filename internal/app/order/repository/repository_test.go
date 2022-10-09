package repository_test

import (
	"context"
	"errors"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/go-faker/faker/v4"
	"github.com/stretchr/testify/assert"

	"github.com/ranggabudipangestu/simple-ecommerce/internal/app/order/dto"
	"github.com/ranggabudipangestu/simple-ecommerce/internal/app/order/repository"
)

func TestCreateOrder(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	var detailOrder []dto.CreateOrderDetails
	detailOrder = append(detailOrder, dto.CreateOrderDetails{
		ProductId: 1,
		Qty:       1,
		Price:     2000000,
		Total:     2000000,
	})
	transactionNumber := "TRX-4541221212411"
	payload := dto.CreateOrderDto{
		DeliveryAddress:  "Indonesia",
		Details:          detailOrder,
		TotalTransaction: 2000000,
		TotalQty:         1,
	}

	t.Run("Test Create Order Success", func(t *testing.T) {

		mock.ExpectBegin()
		query := "INSERT into transaction"
		mock.ExpectExec(query).WithArgs(transactionNumber, payload.DeliveryAddress, payload.TotalQty, payload.TotalTransaction).WillReturnResult(sqlmock.NewResult(1, 1))

		query = "INSERT INTO transaction_detail"
		mock.ExpectExec(query).WithArgs(1, detailOrder[0].ProductId, detailOrder[0].Qty, detailOrder[0].Price, detailOrder[0].Total).WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()
		r := repository.NewOrder(db)
		result, err := r.CreateOrder(context.TODO(), payload, transactionNumber)

		assert.Nil(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Test Create Order Error Table Transaction", func(t *testing.T) {

		mock.ExpectBegin()
		query := "INSERT into transaction"
		mock.ExpectExec(query).WithArgs(transactionNumber, payload.DeliveryAddress, payload.TotalQty, payload.TotalTransaction).WillReturnError(errors.New("Error Database Transaction"))

		mock.ExpectCommit()
		r := repository.NewOrder(db)
		result, err := r.CreateOrder(context.TODO(), payload, transactionNumber)

		assert.NotNil(t, err)
		assert.Nil(t, result)
	})
}

func TestGetOrderDetail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	var mockOrder dto.GetOrderDto
	var mockDetailOrder []dto.GetOrderDetails

	mockDetailOrder = append(mockDetailOrder, dto.GetOrderDetails{
		ID:          1,
		ProductName: "Nike Airmax",
		BrandName:   "Nike",
		Qty:         1,
		Price:       2000000,
		Total:       2000000,
	})

	query := `SELECT id, transactionNumber, deliveryAddress, totalQty, totalTransaction FROM transaction`
	queryDetail := `SELECT
	transaction_detail.id,
	product.title as productName,
	brand.title as brandName,
	transaction_detail.qty,
	transaction_detail.price,
	transaction_detail.total
	FROM transaction_detail
	JOIN product ON product.id = transaction_detail.productId
	JOIN brand ON brand.id = product.brandId`

	t.Run("Test Get Order Detail Success", func(t *testing.T) {
		err = faker.FakeData(&mockOrder)
		assert.NoError(t, err)
		orderRow := sqlmock.NewRows([]string{"id", "transactionNumber", "deliveryAddres", "totalQty", "totalTransaction"}).
			AddRow(mockOrder.ID, mockOrder.TransactionNumber, mockOrder.DeliveryAddress, mockOrder.TotalQty, mockOrder.TotalTransaction)

		detailRows := sqlmock.NewRows([]string{"id", "productName", "brandName", "qty", "price", "total"}).
			AddRow(mockDetailOrder[0].ID, mockDetailOrder[0].ProductName, mockDetailOrder[0].BrandName, mockDetailOrder[0].Qty, mockDetailOrder[0].Price, mockDetailOrder[0].Total)

		mock.ExpectQuery(query).WithArgs(1).WillReturnRows(orderRow)
		mock.ExpectQuery(queryDetail).WithArgs(1).WillReturnRows(detailRows)
		r := repository.NewOrder(db)
		result, err := r.GetOrderDetails(context.TODO(), 1)

		assert.Nil(t, err)
		assert.NotNil(t, result)
	})

	t.Run("Test Get Order Detail Error Get Order", func(t *testing.T) {
		mock.ExpectQuery(query).WithArgs(1).WillReturnError(errors.New("Database Error"))
		r := repository.NewOrder(db)
		result, err := r.GetOrderDetails(context.TODO(), 1)

		assert.NotNil(t, err)
		assert.Nil(t, result)
	})

	t.Run("Test Get Order Detail not found", func(t *testing.T) {
		orderRow := sqlmock.NewRows([]string{"id", "transactionNumber", "deliveryAddres", "totalQty", "totalTransaction"})
		mock.ExpectQuery(query).WithArgs(1).WillReturnRows(orderRow)
		r := repository.NewOrder(db)
		result, err := r.GetOrderDetails(context.TODO(), 1)

		assert.Nil(t, err)
		assert.Nil(t, result)
	})

	t.Run("Test Get Order Detail Error Get Order Detail", func(t *testing.T) {
		err = faker.FakeData(&mockOrder)
		assert.NoError(t, err)
		orderRow := sqlmock.NewRows([]string{"id", "transactionNumber", "deliveryAddres", "totalQty", "totalTransaction"}).
			AddRow(mockOrder.ID, mockOrder.TransactionNumber, mockOrder.DeliveryAddress, mockOrder.TotalQty, mockOrder.TotalTransaction)

		mock.ExpectQuery(query).WithArgs(1).WillReturnRows(orderRow)
		mock.ExpectQuery(queryDetail).WithArgs(1).WillReturnError(errors.New("Database Error"))

		r := repository.NewOrder(db)
		result, err := r.GetOrderDetails(context.TODO(), 1)

		assert.NotNil(t, err)
		assert.Nil(t, result)
	})
}
