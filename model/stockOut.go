package model

import (
	"github.com/rinaldypasya/TokoIjah/inventory"
)

func (db *DB) RemoveProduct(s *inventory.StockOut) {
	db.Create(s)
}

func (db *DB) GetAllOutProducts() []inventory.StockOut {
	var allStockout []inventory.StockOut
	db.Find(&allStockout)
	return allStockout
}

func (db *DB) GetOutProductsBySku(sku string) []inventory.StockOut {
	var stockout []inventory.StockOut
	db.Where("sku = ?", sku).Find(&stockout)
	return stockout
}
