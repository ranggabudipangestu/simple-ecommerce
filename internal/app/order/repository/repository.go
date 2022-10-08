package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/ranggabudipangestu/simple-ecommerce/internal/app/order/dto"
	"github.com/ranggabudipangestu/simple-ecommerce/internal/model"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, dto dto.CreateOrderDto, transactionNumber string) (*model.Transaction, error)
	GetOrderDetails(ctx context.Context, id int) (*dto.GetOrderDto, error)
}

type Repository struct {
	DB *sql.DB
}

func NewOrder(db *sql.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) CreateOrder(ctx context.Context, payload dto.CreateOrderDto, transactionNumber string) (*model.Transaction, error) {

	tx, err := r.DB.BeginTx(ctx, &sql.TxOptions{})

	//PROCESS ORDER
	query := `INSERT into transaction (transactionNumber, deliveryAddress, totalQty, totalTransaction, createdAt) values(?, ?, ?, ?, NOW())`
	result, err := tx.ExecContext(ctx, query, transactionNumber, payload.DeliveryAddress, payload.TotalQty, payload.TotalTransaction)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	var id int64
	id, err = result.LastInsertId()
	if err != nil {
		return nil, err
	}
	//END OF PROCESS ORDER

	//PROCESS ORDER DETAIL
	var (
		placeholders []string
		details      []interface{}
	)

	for _, detail := range payload.Details {
		placeholders = append(placeholders, "(?,?,?,?,?)")
		details = append(details, id, detail.ProductId, detail.Qty, detail.Price, detail.Total)
	}

	query = fmt.Sprintf("INSERT INTO transaction_detail (transactionId, productId, qty, price, total) VALUES %s", strings.Join(placeholders, ","))
	_, err = tx.ExecContext(ctx, query, details...)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	//END OF PROCESS ORDER DETAIL

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &model.Transaction{ID: int(id)}, nil
}

func (r *Repository) GetOrderDetails(ctx context.Context, id int) (*dto.GetOrderDto, error) {

	//PROCESS GET ORDER DATA BY ID
	query := `SELECT id, transactionNumber, deliveryAddress, totalQty, totalTransaction FROM transaction where id = ? LIMIT 1`
	row, err := r.DB.QueryContext(ctx, query, id)

	if err != nil {
		return nil, err
	}

	var data dto.GetOrderDto
	if !row.Next() {
		return nil, nil
	}

	err = row.Scan(
		&data.ID,
		&data.TransactionNumber,
		&data.DeliveryAddress,
		&data.TotalQty,
		&data.TotalTransaction,
	)
	if err != nil {
		return nil, err
	}
	//END OF PROCESS GET ORDER DATA BY ID

	//PROCESS GET ORDER DETAIL BY ORDER ID
	queryDetail := `SELECT
	transaction_detail.id,
	product.title as productName,
	brand.title as brandName,
	transaction_detail.qty,
	transaction_detail.price,
	transaction_detail.total
	FROM transaction_detail
	JOIN product ON product.id = transaction_detail.productId
	JOIN brand ON brand.id = product.brandId
	WHERE transaction_detail.transactionId = ?
	ORDER BY transaction_detail.id`

	detailRows, err := r.DB.QueryContext(ctx, queryDetail, id)
	for detailRows.Next() {
		transactionDetail := dto.GetOrderDetails{}
		err = detailRows.Scan(
			&transactionDetail.ID,
			&transactionDetail.ProductName,
			&transactionDetail.BrandName,
			&transactionDetail.Qty,
			&transactionDetail.Price,
			&transactionDetail.Total,
		)

		data.Details = append(data.Details, transactionDetail)
	}
	if err != nil {
		return nil, err
	}
	//END PROCESS GET ORDER DETAIL BY ORDER ID

	defer func() {
		row.Close()
		detailRows.Close()
	}()

	return &data, nil
}
