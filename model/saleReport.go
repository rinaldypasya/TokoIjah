package model

import (
	"github.com/rinaldypasya/TokoIjah/inventory"
)

func (db *DB) CreateSaleReport(s *inventory.SaleReport) {
	db.Create(s)
}

func (db *DB) GetAllSaleReports() []inventory.SaleReport {
	var allSaleReports []inventory.SaleReport
	db.Find(&allSaleReports)
	return allSaleReports
}

func (db *DB) GetSaleReportsBySKU(sku string) []inventory.SaleReport {
	var saleReports []inventory.SaleReport
	db.Where("sku = ?", sku).Find(&saleReports)
	return saleReports
}

func (db *DB) GetSaleReportsByDate(dateFrom string, dateTo string) []inventory.SaleReport {
	var saleReports []inventory.SaleReport
	db.Where("timestamp BETWEEN ? AND ?", dateFrom, dateTo).Find(&saleReports)
	return saleReports
}
