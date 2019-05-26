package model

import (
	"github.com/rinaldypasya/TokoIjah/inventory"
)

func (db *DB) CreateStockValue(s *inventory.StockValue) {
	db.Create(s)
}

func (db *DB) GetAllStockValues() []inventory.StockValue {
	var allStockvalue []inventory.StockValue
	db.Find(&allStockvalue)
	return allStockvalue
}

func (db *DB) GetStockValueByID(ID int) inventory.StockValue {
	var stockvalue inventory.StockValue
	db.First(&stockvalue, ID)
	return stockvalue
}

func (db *DB) GetStockValuesBySku(sku string) inventory.StockValue {
	var stockvalue inventory.StockValue
	db.Where("sku = ?", sku).First(&stockvalue)
	return stockvalue
}

func (db *DB) UpdateStockValue(s inventory.StockValue) inventory.StockValue {
	db.Save(s)
	return s
}
