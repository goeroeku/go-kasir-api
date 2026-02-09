package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"kasir-api/models"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(req models.CheckoutRequest) (*models.Transaction, error) {
	ctx := context.Background()
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	// Defer rollback in case of panic or error (if not committed)
	defer tx.Rollback()

	var totalAmount int
	var details []models.TransactionDetail

	// 1. Calculate total and validate stock for all items
	for _, item := range req.Items {
		var price, stock int
		var name string

		// Get product info and lock row for update
		err := tx.QueryRowContext(ctx, "SELECT name, price, stock FROM products WHERE id = $1 FOR UPDATE", item.ProductID).Scan(&name, &price, &stock)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, fmt.Errorf("product with ID %d not found", item.ProductID)
			}
			return nil, err
		}

		if stock < item.Quantity {
			return nil, fmt.Errorf("insufficient stock for product %s (ID: %d)", name, item.ProductID)
		}

		subtotal := price * item.Quantity
		totalAmount += subtotal

		// 2. Decrease stock
		_, err = tx.ExecContext(ctx, "UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}

		details = append(details, models.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: name,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	// 3. Create Transaction record
	var transaction models.Transaction
	err = tx.QueryRowContext(ctx, "INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id, total_amount, created_at", totalAmount).Scan(&transaction.ID, &transaction.TotalAmount, &transaction.CreatedAt)
	if err != nil {
		return nil, err
	}

	// 4. Create Transaction Details records
	for i, detail := range details {
		var detailID int
		err := tx.QueryRowContext(ctx,
			"INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4) RETURNING id",
			transaction.ID, detail.ProductID, detail.Quantity, detail.Subtotal,
		).Scan(&detailID)
		if err != nil {
			return nil, err
		}
		details[i].ID = detailID
		details[i].TransactionID = transaction.ID
	}

	transaction.Details = details

	// 5. Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &transaction, nil
}
