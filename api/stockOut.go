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

// RemoveProduct make note for products coming out from stock
func RemoveProduct(db inventory.InventStockOut, dbStock inventory.InventStock) gin.HandlerFunc {
	return func(gc *gin.Context) {

		var stock inventory.Stock
		var StockOut inventory.StockOut

		if gc.BindJSON(&StockOut) == nil {
			StockOut.Timestamp = time.Now().Format("2006-01-02 15:04:05")
			StockOut.Total = StockOut.OutAmount * StockOut.SalePrice

			stock = dbStock.GetStockBySku(StockOut.Sku)
			if stock.Sku != "" { // if already existed before, update the stock
				db.RemoveProduct(&StockOut)
				stock.Amount -= StockOut.OutAmount
				updatedStock := dbStock.UpdateStock(stock)
				gc.JSON(http.StatusOK, gin.H{
					"status":  "true",
					"message": "Products removed successfully",
					"id":      StockOut.ID,
					"stock":   updatedStock.Amount,
				})
				return
			} else { // if never exist before, then this product was never in stock before
				gc.JSON(http.StatusBadRequest, gin.H{
					"status":  "false",
					"message": "Products were never in stock before!",
				})
				return
			}

		} else {
			gc.JSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Bad Request",
			})
			return
		}
		return
	}
}

// GetAllOutProducts get all records of products coming out from stock
func GetAllOutProducts(db inventory.InventStockOut) gin.HandlerFunc {
	return func(gc *gin.Context) {

		var StockOut []inventory.StockOut
		StockOut = db.GetAllOutProducts()
		if len(StockOut) > 0 {
			gc.JSON(http.StatusOK, gin.H{
				"status": "true",
				"data":   StockOut,
			})
			return
		} else {
			gc.JSON(http.StatusOK, gin.H{
				"status":  "true",
				"message": "No records of products out from store yet!",
			})
			return
		}

		return
	}
}

// GetOutProductsBySku get products coming out from stock by sku
func GetOutProductsBySku(db inventory.InventStockOut) gin.HandlerFunc {
	return func(gc *gin.Context) {

		var StockOut []inventory.StockOut
		StockOut = db.GetOutProductsBySku(gc.Param("sku"))
		if StockOut[0].Sku == gc.Param("sku") {
			gc.JSON(http.StatusOK, gin.H{
				"status": "true",
				"data":   StockOut,
			})
			return
		} else {
			gc.JSON(http.StatusNotFound, gin.H{
				"status": "false",
				"data":   "No records not found!",
			})
			return
		}
	}
}

// StockOutExportToCSV export all records data into csv file
func StockOutExportToCSV(db inventory.InventStockOut) gin.HandlerFunc {
	return func(gc *gin.Context) {

		var allStockOut []inventory.StockOut
		allStockOut = db.GetAllOutProducts()

		csvdata := init2dArray(len(allStockOut), 10)

		for i := 0; i < len(allStockOut); i++ {
			csvdata[i][0] = strconv.Itoa(i + 1)
			csvdata[i][1] = allStockOut[i].Timestamp
			csvdata[i][2] = allStockOut[i].Sku
			csvdata[i][3] = allStockOut[i].Name
			csvdata[i][4] = strconv.Itoa(allStockOut[i].OutAmount)
			csvdata[i][5] = strconv.Itoa(allStockOut[i].SalePrice)
			csvdata[i][6] = strconv.Itoa(allStockOut[i].Total)
			csvdata[i][7] = allStockOut[i].Note
		}

		fileName := time.Now().Format("2006-01-02") + "-StockOut.csv"
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
			"message":  "StockOut data exported to csv successfully!",
			"filename": fileName,
		})
		return
	}
}

// StockOutImportCSV import csv data into StockOut table
func StockOutImportCSV(db inventory.InventStockOut) gin.HandlerFunc {
	return func(gc *gin.Context) {

		var StockOut []inventory.StockOut

		file, _ := gc.FormFile("StockOutimport")
		dst := "./csv/" + file.Filename
		gc.SaveUploadedFile(file, dst)
		// csvfile, err := os.Open("./csv/import_stock.csv")
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

			StockOutoutamount, _ := strconv.Atoi(line[4])
			StockOutsaleprice, _ := strconv.Atoi(line[5])
			StockOuttotal, _ := strconv.Atoi(line[6])
			StockOut = append(StockOut, inventory.StockOut{
				Timestamp: line[1], // start from timestamp column as we ignore id column(assume the csv include the IDs)
				Sku:       line[2],
				Name:      line[3],
				OutAmount: StockOutoutamount,
				SalePrice: StockOutsaleprice,
				Total:     StockOuttotal,
				Note:      line[7],
			})
		}

		if len(StockOut) > 0 {
			for i := 0; i < len(StockOut); i++ {
				db.RemoveProduct(&StockOut[i])
			}
			gc.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "data csv migrated successfully to StockOut table",
				"data":    StockOut,
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
		return

	}
}
