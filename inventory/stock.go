package inventory

type Stock struct {
	ID     int    `gorm:"primary_key;AUTO_INCREMENT" json:"-"`
	Sku    string `gorm:"not null" json:"sku"`
	Name   string `gorm:"not null" json:"name"`
	Amount int    `json:"amount"`
}

// InventStock inventory stock interface
type InventStock interface {
	CreateStock(*Stock)
	GetAllStock() []Stock
	GetStockByID(int) Stock
	GetStockBySku(string) Stock
	UpdateStock(Stock) Stock
}
