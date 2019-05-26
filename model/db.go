package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rinaldypasya/TokoIjah/inventory"
	"log"
)

var (
	dbEngine = "sqlite3"
	dbName   = "./stock.db"
)

type DB struct {
	*gorm.DB
}

func InitDB() *DB {
	db, err := gorm.Open(dbEngine, dbName)
	if err != nil {
		log.Fatal("failed to initialize database: ", err.Error())
	}

	db.AutoMigrate(&inventory.Stock{}, &inventory.StockIn{}, &inventory.StockOut{}, &inventory.StockValue{}, &inventory.SaleReport{})

	return &DB{db}
}
