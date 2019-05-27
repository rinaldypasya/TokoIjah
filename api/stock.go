package api

import (
	"bufio"
	"encoding/csv"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/rinaldypasya/TokoIjah/inventory"
)

// CreateStock create one instance of stock data and save it to stocks table
func CreateStock(db inventory.InventStock) gin.HandlerFunc {
	return func(gc *gin.Context) {

		var stock inventory.Stock
		if gc.BindJSON(&stock) == nil {
			db.CreateStock(&stock)
			gc.JSON(http.StatusOK, gin.H{
				"status":  "true",
				"message": "Stock data created successfully!",
				"id":      stock.ID,
			})
			return
		} else {
			gc.JSON(http.StatusBadRequest, gin.H{
				"status":  "false",
				"message": "Bad Request",
			})
			return
		}
	}
}

// GetAllStock get all rows of stocks from stocks table
func GetAllStock(db inventory.InventStock) gin.HandlerFunc {
	return func(gc *gin.Context) {

		stock := db.GetAllStock()
		if len(stock) > 0 {
			gc.JSON(http.StatusOK, gin.H{
				"status": "true",
				"data":   stock,
			})
			return
		} else {
			gc.JSON(http.StatusOK, gin.H{
				"status":  "true",
				"message": "Stock is empty!",
			})
			return
		}
	}
}

// GetStockByID get stock data by id
func GetStockByID(db inventory.InventStock) gin.HandlerFunc {
	return func(gc *gin.Context) {

		var stock inventory.Stock
		id, err := strconv.Atoi(gc.Param("id"))
		if err != nil {
			gc.JSON(http.StatusBadRequest, gin.H{
				"status":  "false",
				"message": "error!",
			})
			return
		}

		stock = db.GetStockByID(id)
		if stock.ID == id {
			gc.JSON(http.StatusOK, gin.H{
				"status": "true",
				"data":   stock,
			})
			return
		} else {
			gc.JSON(http.StatusNotFound, gin.H{
				"status": "false",
				"data":   "Stock not found!",
			})
			return
		}
	}
}

// GetStockBySku get stock data by sku
func GetStockBySku(db inventory.InventStock) gin.HandlerFunc {
	return func(gc *gin.Context) {

		stock := db.GetStockBySku(gc.Param("sku"))
		if stock.Sku == gc.Param("sku") {
			gc.JSON(http.StatusOK, gin.H{
				"status": "true",
				"data":   stock,
			})
			return
		} else {
			gc.JSON(http.StatusNotFound, gin.H{
				"status": "false",
				"data":   "Stock not found!",
			})
			return
		}
	}
}

// UpdateStock update already existing stock data
func UpdateStock(db inventory.InventStock) gin.HandlerFunc {
	return func(gc *gin.Context) {

		var stock inventory.Stock
		if gc.BindJSON(&stock) == nil {
			updatedata := db.GetStockBySku(stock.Sku)
			updatedata.Name = stock.Name
			updatedata.Amount = stock.Amount
			updated := db.UpdateStock(updatedata)
			gc.JSON(http.StatusOK, gin.H{
				"status":       "true",
				"message":      "Stock updated successfully",
				"Updated Data": updated,
			})
			return
		} else {
			gc.JSON(http.StatusBadRequest, gin.H{
				"status":  "false",
				"message": "Check data carefully",
			})
			return
		}
	}
}

// StockExportToCSV export all stock data from stocks table
func StockExportToCSV(db inventory.InventStock) gin.HandlerFunc {
	return func(gc *gin.Context) {

		allstock := db.GetAllStock()

		csvdata := init2dArray(len(allstock), 4)

		for i := 0; i < len(allstock); i++ {
			csvdata[i][0] = strconv.Itoa(i + 1)
			csvdata[i][1] = allstock[i].Sku
			csvdata[i][2] = allstock[i].Name
			csvdata[i][3] = strconv.Itoa(allstock[i].Amount)
		}

		fileName := time.Now().Format("2006-01-02") + "-Stock.csv"
		file, err := os.Create("./csv/" + fileName)
		if err != nil {
			gc.JSON(http.StatusConflict, gin.H{
				"status":  false,
				"message": "Failed to export file!",
			})
			return
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		for _, value := range csvdata {
			err := writer.Write(value)
			if err != nil {
				gc.JSON(http.StatusConflict, gin.H{
					"status":  false,
					"message": "Failed to export file!",
				})
				return
			}
		}

		gc.JSON(http.StatusOK, gin.H{
			"status":   true,
			"message":  "Stock data exported to csv successfully!",
			"filename": fileName,
		})
	}
}

// StockImportCSV import csv data into stocks table
func StockImportCSV(db inventory.InventStock) gin.HandlerFunc {
	return func(gc *gin.Context) {

		var stock []inventory.Stock

		file, _ := gc.FormFile("stockimport")
		dst := "./csv/" + file.Filename
		gc.SaveUploadedFile(file, dst)
		csvfile, err := os.Open("./csv/" + file.Filename)
		if err != nil {
			gc.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "error opening file, check file again",
			})
		}

		reader := csv.NewReader(bufio.NewReader(csvfile))
		for {
			line, error := reader.Read()
			if error == io.EOF {
				break
			} else if error != nil {
				gc.JSON(http.StatusBadRequest, gin.H{
					"status":  false,
					"message": "something's wrong!",
				})
			}

			stockamount, _ := strconv.Atoi(line[3])
			stock = append(stock, inventory.Stock{
				Sku:    line[1], // start from sku column as we ignore id column(assume the csv include the IDs)
				Name:   line[2],
				Amount: stockamount,
			})
		}

		if len(stock) > 0 {
			for i := 0; i < len(stock); i++ {
				tableStock := db.GetStockBySku(stock[i].Sku)
				if tableStock.Sku != "" { // data already exist in stock table, update the data then
					tableStock.Name = stock[i].Name
					tableStock.Amount += stock[i].Amount
					updatedStock := db.UpdateStock(tableStock)
					_ = updatedStock // LOL
				} else { // data didn't exist in table, create new one then
					db.CreateStock(&stock[i])
				}
			}
			gc.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "data csv migrated successfully to stock table",
				"data":    stock,
			})
			return

		} else {
			gc.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "error reading csv file, check file again for correct format!",
			})
			return

		}

		gc.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "something's wrong!",
		})
	}
}

func init2dArray(y int, x int) [][]string {
	a := make([][]string, y)
	for i := range a {
		a[i] = make([]string, x)
	}
	return a
}
