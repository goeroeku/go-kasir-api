package models

// SalesReport represents the sales report data
type SalesReport struct {
	TotalRevenue      int        `json:"total_revenue"`
	TotalTransactions int        `json:"total_transactions"`
	BestSeller        BestSeller `json:"best_seller"`
}

// BestSeller represents the best selling product
type BestSeller struct {
	ProductID   int    `json:"product_id"`
	ProductName string `json:"product_name"`
	Quantity    int    `json:"quantity"`
}
