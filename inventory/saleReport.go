package inventory

type SaleReport struct {
	ID          int    `gorm:"primary_key;AUTO_INCREMENT" json:"-"`
	OrderID     string `json:"orderid"`
	Timestamp   string `json:"timestamp"`
	Sku         string `json:"sku"`
	Name        string `json:"name"`
	Amount      int    `json:"amount"`
	Saleprice   int    `json:"saleprice"`
	Total       int    `json:"total"`
	Buyingprice int    `json:"buyingprice"`
	Profit      int    `json:"profit"`
}

type InventSaleReport interface {
	CreateSaleReport(*SaleReport)
	GetAllSaleReports() []SaleReport
	GetSaleReportsBySKU(string) []SaleReport
	GetSaleReportsByDate(string, string) []SaleReport
}
