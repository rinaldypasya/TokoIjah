package inventory

type StockOut struct {
	ID        int    `gorm:"primary_key;AUTO_INCREMENT" json:"-"`
	Timestamp string `gorm:"not null" json:"timestamp"`
	Sku       string `gorm:"not null" json:"sku"`
	Name      string `json:"name"`
	OutAmount int    `json:"outamount"`
	SalePrice int    `json:"saleprice"`
	Total     int    `json:"total"`
	Note      string `json:"note"`
}

type InventStockOut interface {
	RemoveProduct(*StockOut)
	GetAllOutProducts() []StockOut
	GetOutProductsBySku(string) []StockOut
}
