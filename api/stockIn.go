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

// StoreProduct make note for new products imported into stock
func StoreProduct(db inventory.InventStockIn, dbStock inventory.InventStock) gin.HandlerFunc {
	return func(gc *gin.Context) {

		var stock inventory.Stock
		var StockIn inventory.StockIn

		if gc.BindJSON(&StockIn) == nil {
			StockIn.Timestamp = time.Now().Format("2006-01-02 15:04:05")
			StockIn.Total = StockIn.OrderAmount * StockIn.BuyingPrice
			db.StoreProduct(&StockIn)

			stock = dbStock.GetStockBySku(StockIn.Sku)
			if stock.Sku != "" { // if product already existed, update the amount
				stock.Amount += StockIn.ReceivedAmount
				updatedStock := dbStock.UpdateStock(stock)
				gc.JSON(http.StatusOK, gin.H{
					"status":  "true",
					"message": "Products stored successfully",
					"id":      StockIn.ID,
					"stock":   updatedStock.Amount,
				})
				return
			} else { // if never exist before, create new stock
				stock.Sku = StockIn.Sku
				stock.Name = StockIn.Name
				stock.Amount = StockIn.ReceivedAmount
				dbStock.CreateStock(&stock)
				gc.JSON(http.StatusOK, gin.H{
					"status":  "true",
					"message": "New Products stored successfully",
					"id":      StockIn.ID,
					"stock":   stock.Amount,
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

// GetAllStoredProducts get all data records from StockIn table
func GetAllStoredProducts(db inventory.InventStockIn) gin.HandlerFunc {
	return func(gc *gin.Context) {

		var StockIn []inventory.StockIn
		StockIn = db.GetAllStoredProducts()
		if len(StockIn) > 0 {
			gc.JSON(http.StatusOK, gin.H{
				"status": "true",
				"data":   StockIn,
			})
			return
		} else {
			gc.JSON(http.StatusOK, gin.H{
				"status":  "true",
				"message": "Store records is empty!",
			})
			return
		}

		return
	}
}

// GetStoredProductsBySku get records data by sku
func GetStoredProductsBySku(db inventory.InventStockIn) gin.HandlerFunc {
	return func(gc *gin.Context) {

		var StockIn []inventory.StockIn
		StockIn = db.GetStoredProductsBySku(gc.Param("sku"))
		if StockIn[0].Sku == gc.Param("sku") {
			gc.JSON(http.StatusOK, gin.H{
				"status": "true",
				"data":   StockIn,
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

// StockInExportToCSV export all records from StockIns table into csv file
func StockInExportToCSV(db inventory.InventStockIn) gin.HandlerFunc {
	return func(gc *gin.Context) {

		var allStockIn []inventory.StockIn
		allStockIn = db.GetAllStoredProducts()

		csvdata := init2dArray(len(allStockIn), 10)

		for i := 0; i < len(allStockIn); i++ {
			csvdata[i][0] = strconv.Itoa(i + 1)
			csvdata[i][1] = allStockIn[i].Timestamp
			csvdata[i][2] = allStockIn[i].Sku
			csvdata[i][3] = allStockIn[i].Name
			csvdata[i][4] = strconv.Itoa(allStockIn[i].OrderAmount)
			csvdata[i][5] = strconv.Itoa(allStockIn[i].ReceivedAmount)
			csvdata[i][6] = strconv.Itoa(allStockIn[i].BuyingPrice)
			csvdata[i][7] = strconv.Itoa(allStockIn[i].Total)
			csvdata[i][8] = allStockIn[i].Receipt
			csvdata[i][9] = allStockIn[i].Note
		}

		fileName := time.Now().Format("2006-01-02") + "-StockIn.csv"
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
			"message":  "StockIn data exported to csv successfully!",
			"filename": fileName,
		})
		return
	}
}

// StockInImportCSV import csv data into StockIns table
func StockInImportCSV(db inventory.InventStockIn) gin.HandlerFunc {
	return func(gc *gin.Context) {

		var StockIn []inventory.StockIn

		file, _ := gc.FormFile("StockInimport")
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

			StockInorderamount, _ := strconv.Atoi(line[4])
			StockInreceivedamount, _ := strconv.Atoi(line[5])
			StockInbuyingprice, _ := strconv.Atoi(line[6])
			StockIntotal, _ := strconv.Atoi(line[7])
			StockIn = append(StockIn, inventory.StockIn{
				Timestamp:      line[1], // start from timestamp column as we ignore id column(assume the csv include the IDs)
				Sku:            line[2],
				Name:           line[3],
				OrderAmount:    StockInorderamount,
				ReceivedAmount: StockInreceivedamount,
				BuyingPrice:    StockInbuyingprice,
				Total:          StockIntotal,
				Receipt:        line[8],
				Note:           line[9],
			})
		}

		if len(StockIn) > 0 {
			for i := 0; i < len(StockIn); i++ {
				db.StoreProduct(&StockIn[i])
			}
			gc.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "data csv migrated successfully to stock table",
				"data":    StockIn,
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
