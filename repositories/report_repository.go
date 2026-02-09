package repositories

import (
	"database/sql"
	"kasir-api/models"
	"time"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) GetSalesReport(startDate, endDate time.Time) (*models.SalesReport, error) {
	var report models.SalesReport

	// 1. Calculate Total Revenue
	err := r.db.QueryRow(
		"SELECT COALESCE(SUM(total_amount), 0) FROM transactions WHERE created_at BETWEEN $1 AND $2",
		startDate, endDate,
	).Scan(&report.TotalRevenue)
	if err != nil {
		return nil, err
	}

	// 2. Calculate Total Transactions
	err = r.db.QueryRow(
		"SELECT COUNT(id) FROM transactions WHERE created_at BETWEEN $1 AND $2",
		startDate, endDate,
	).Scan(&report.TotalTransactions)
	if err != nil {
		return nil, err
	}

	// 3. Find Best Seller
	err = r.db.QueryRow(`
		SELECT td.product_id, p.name, SUM(td.quantity) as total_qty
		FROM transaction_details td
		JOIN transactions t ON td.transaction_id = t.id
		JOIN products p ON td.product_id = p.id
		WHERE t.created_at BETWEEN $1 AND $2
		GROUP BY td.product_id, p.name
		ORDER BY total_qty DESC
		LIMIT 1
	`, startDate, endDate).Scan(&report.BestSeller.ProductID, &report.BestSeller.ProductName, &report.BestSeller.Quantity)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	// If no sales, best seller will be empty (zero values), which is fine

	return &report, nil
}
