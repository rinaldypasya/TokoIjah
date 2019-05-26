package model

import (
	"github.com/rinaldypasya/TokoIjah/inventory"
)

func (db *DB) StoreProduct(s *inventory.StockIn) {
	db.Create(s)
}

func (db *DB) GetAllStoredProducts() []inventory.StockIn {
	var allStoredProducts []inventory.StockIn
	db.Find(&allStoredProducts)
	return allStoredProducts
}

func (db *DB) GetStoredProductsBySku(sku string) []inventory.StockIn {
	var storedProduct []inventory.StockIn
	db.Where("sku = ?", sku).Find(&storedProduct)
	return storedProduct
}
