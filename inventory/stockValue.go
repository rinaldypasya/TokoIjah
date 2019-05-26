package inventory

type StockValue struct {
	ID          int    `gorm:"primary_key;AUTO_INCREMENT" json:"-"`
	Sku         string `gorm:"not null" json:"sku"`
	Name        string `json:"name"`
	Amount      int    `json:"amount"`
	BuyingPrice int    `json:"buyingprice"`
	Total       int    `json:"total"`
}

type IStockvalue interface {
	CreateStockValue(*StockValue)
	GetAllStockValues() []StockValue
	GetStockValueByID(int) StockValue
	GetStockValuesBySku(string) StockValue
	UpdateStockValue(StockValue) StockValue
}
