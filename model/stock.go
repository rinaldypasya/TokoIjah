package model

import (
	"github.com/rinaldypasya/TokoIjah/inventory"
)

func (db *DB) CreateStock(s *inventory.Stock) {
	db.Create(s)
}

func (db *DB) GetAllStock() []inventory.Stock {
	var allStock []inventory.Stock
	db.Find(&allStock)
	return allStock
}

func (db *DB) GetStockByID(ID int) inventory.Stock {
	var stock inventory.Stock
	db.First(&stock, ID)
	return stock
}

func (db *DB) GetStockBySku(sku string) inventory.Stock {
	var stock inventory.Stock
	db.Where("sku = ?", sku).First(&stock)
	return stock
}

func (db *DB) UpdateStock(s inventory.Stock) inventory.Stock {
	db.Save(s)
	return s
}
