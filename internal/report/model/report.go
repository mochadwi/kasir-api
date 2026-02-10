package model

type SalesReport struct {
	TotalRevenue      int    `json:"total_revenue"`
	TotalTransactions int    `json:"total_transactions"`
	ProdukTerlaris    string `json:"produk_terlaris"`
	QtyTerjual        int    `json:"qty_terjual"`
}
