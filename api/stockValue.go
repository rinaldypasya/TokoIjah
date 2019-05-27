package api

import (
	"bufio"
	"encoding/csv"
	"io"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/rinaldypasya/TokoIjah/inventory"
)

// CreateStockValue create one instance of stock value into StockValues table
func CreateStockValue(db inventory.InventStockValue) gin.HandlerFunc {
	return func(gc *gin.Context) {

		var StockValue inventory.StockValue
		if gc.BindJSON(&StockValue) == nil {
			StockValue.Total = StockValue.BuyingPrice * StockValue.Amount
			db.CreateStockValue(&StockValue)
			gc.JSON(http.StatusOK, gin.H{
				"status":  "true",
				"message": "Stock value created successfully",
				"id":      StockValue.ID,
			})
			return
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

// GetAllStockValues get all records from StockValues table
func GetAllStockValues(db inventory.InventStockValue) gin.HandlerFunc {
	return func(gc *gin.Context) {

		var StockValue []inventory.StockValue
		StockValue = db.GetAllStockValues()
		if len(StockValue) > 0 {
			gc.JSON(http.StatusOK, gin.H{
				"status": "true",
				"data":   StockValue,
			})
			return
		} else {
			gc.JSON(http.StatusOK, gin.H{
				"status":  "true",
				"message": "StockValues is empty!",
			})
			return
		}

		return
	}
}

// GetStockValueByID get a StockValue data by id
func GetStockValueByID(db inventory.InventStockValue) gin.HandlerFunc {
	return func(gc *gin.Context) {

		var StockValue inventory.StockValue
		id, err := strconv.Atoi(gc.Param("id"))
		if err != nil {
			gc.JSON(http.StatusBadRequest, gin.H{
				"status":  "false",
				"message": "error!",
			})
			return
		}
		StockValue = db.GetStockValueByID(id)
		if StockValue.ID == id {
			gc.JSON(http.StatusOK, gin.H{
				"status": "true",
				"data":   StockValue,
			})
			return
		} else {
			gc.JSON(http.StatusNotFound, gin.H{
				"status": "false",
				"data":   "StockValue not found!",
			})
			return
		}
	}
}

// GetStockValuesBySku get a stock value data by sku
func GetStockValuesBySku(db inventory.InventStockValue) gin.HandlerFunc {
	return func(gc *gin.Context) {

		var StockValue inventory.StockValue
		StockValue = db.GetStockValuesBySku(gc.Param("sku"))
		if StockValue.Sku == gc.Param("sku") {
			gc.JSON(http.StatusOK, gin.H{
				"status": "true",
				"data":   StockValue,
			})
			return
		} else {
			gc.JSON(http.StatusNotFound, gin.H{
				"status": "false",
				"data":   "StockValue not found!",
			})
			return
		}
	}
}

// UpdateStockValue update an already existing stock value data
func UpdateStockValue(db inventory.InventStockValue) gin.HandlerFunc {
	return func(gc *gin.Context) {

		var StockValue inventory.StockValue
		if gc.BindJSON(&StockValue) == nil {
			updatedata := db.GetStockValuesBySku(StockValue.Sku)
			updatedata.Name = StockValue.Name
			updatedata.Amount = StockValue.Amount
			updatedata.BuyingPrice = StockValue.BuyingPrice
			updatedata.Total = StockValue.BuyingPrice * StockValue.Amount
			updated := db.UpdateStockValue(updatedata)
			gc.JSON(http.StatusOK, gin.H{
				"status":       "true",
				"message":      "StockValue updated successfully",
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

// StockValueExportToCSV export all records from StockValues table
func StockValueExportToCSV(db inventory.InventStockValue) gin.HandlerFunc {
	return func(gc *gin.Context) {

		var allStockValue []inventory.StockValue
		allStockValue = db.GetAllStockValues()

		csvdata := init2dArray(len(allStockValue), 10)

		for i := 0; i < len(allStockValue); i++ {
			csvdata[i][0] = strconv.Itoa(i + 1)
			csvdata[i][1] = allStockValue[i].Sku
			csvdata[i][2] = allStockValue[i].Name
			csvdata[i][3] = strconv.Itoa(allStockValue[i].Amount)
			csvdata[i][4] = strconv.Itoa(allStockValue[i].BuyingPrice)
			csvdata[i][5] = strconv.Itoa(allStockValue[i].Total)
		}

		fileName := time.Now().Format("2006-01-02") + "-StockValue.csv"
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
			"message":  "StockValue data exported to csv successfully!",
			"filename": fileName,
		})
		return
	}
}

// GenerateStockValue generate report from stock value in form of json data or csv
func GenerateStockValue(db inventory.InventStockValue, dbstock inventory.InventStock, dbstockin inventory.InventStockin) gin.HandlerFunc {
	return func(gc *gin.Context) {

		var stock []inventory.Stock
		var stockin []inventory.Stockin
		var tempStockValue inventory.StockValue

		stock = dbstock.GetAllStock() // get all products from actual stock
		if len(stock) > 0 {
			for i := 0; i < len(stock); i++ { // loop through products in stock
				var StockValue inventory.StockValue // if there're products in stock, initialize StockValue
				StockValue.Sku = stock[i].Sku
				StockValue.Name = stock[i].Name
				StockValue.Amount = stock[i].Amount
				stockin = dbstockin.GetStoredProductsBySku(StockValue.Sku) // look stockin records for each products in stock by sku to calculate average price
				if len(stockin) > 0 {                                      // if they exist in stockin records, then calculate average price
					sumTotal := 0          // sum of products total price in stockin records by sku
					sumReceivedAmount := 0 // sum of products receivedamount in stockin records by sku
					for j := 0; j < len(stockin); j++ {
						sumTotal += stockin[j].Total
						sumReceivedAmount += stockin[j].ReceivedAmount
					}
					averageValue := float64(sumTotal) / float64(sumReceivedAmount) // StockValue average price = sum of stockin total / sum of stockin received amount
					StockValue.BuyingPrice = int(Round(averageValue, .5, 0))
					StockValue.Total = StockValue.Amount * StockValue.BuyingPrice

					// check for existing products inside StockValues, if exist update the amount, buying price, and total. If not exist, create new one
					tempStockValue = db.GetStockValuesBySku(StockValue.Sku)
					if tempStockValue.Sku != "" { // already exist in StockValue table, update amount, price and total
						tempStockValue.Amount = StockValue.Amount
						tempStockValue.BuyingPrice = StockValue.BuyingPrice
						tempStockValue.Total = StockValue.Total
						updatedStockValues := db.UpdateStockValue(tempStockValue)
						_ = updatedStockValues
					} else { // not exist and never recorded in StockValue before, then create new one
						db.CreateStockValue(&StockValue)
					}

				} else { // if the products by that sku dont exist in stockin, then the products never recorded into stockin
					gc.JSON(http.StatusBadRequest, gin.H{
						"status":  false,
						"message": "Products in stock were never in Stockin records!",
					})
					return
				} //end if

			} // end loop stock
			// if for loop in stock for stock value finish, calculating value of stock is done, then calculate report
			var allStockValue []inventory.StockValue
			allStockValue = db.GetAllStockValues()
			sumOfStockValueAmounts := 0
			sumOfStockValueTotals := 0
			for k := 0; k < len(allStockValue); k++ {
				sumOfStockValueAmounts += allStockValue[k].Amount
				sumOfStockValueTotals += allStockValue[k].Total
			}
			gc.JSON(http.StatusOK, gin.H{
				"status":            true,
				"message":           "Calculating stock values is done!",
				"Date":              time.Now().Format("2006-01-02"),
				"Total SKU":         len(allStockValue),
				"Total of Products": sumOfStockValueAmounts,
				"Total Value":       sumOfStockValueTotals,
				"Stock Values":      allStockValue,
			})
			return
		} else {
			gc.JSON(http.StatusNoContent, gin.H{
				"status":  false,
				"message": "Stock is empty!",
			})
			return
		}

		return
	}
}

// StockValueImportCSV import csv data into StockValues table
func StockValueImportCSV(db inventory.InventStockValue) gin.HandlerFunc {
	return func(gc *gin.Context) {

		var StockValue []inventory.StockValue

		file, _ := gc.FormFile("StockValueimport")
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

			StockValueamount, _ := strconv.Atoi(line[3])
			StockValuebuyingprice, _ := strconv.Atoi(line[4])
			StockValuetotal, _ := strconv.Atoi(line[5])
			StockValue = append(StockValue, inventory.StockValue{
				Sku:         line[1], // start from sku column as we ignore id column(assume the csv include the IDs)
				Name:        line[2],
				Amount:      StockValueamount,
				BuyingPrice: StockValuebuyingprice,
				Total:       StockValuetotal,
			})
		}

		if len(StockValue) > 0 {
			for i := 0; i < len(StockValue); i++ {
				var tableStockValue inventory.StockValue
				tableStockValue = db.GetStockValuesBySku(StockValue[i].Sku)
				if tableStockValue.Sku != "" { // data already exist in stock table, update the data then
					tableStockValue.Name = StockValue[i].Name
					tableStockValue.Amount += StockValue[i].Amount
					tableStockValue.BuyingPrice = StockValue[i].BuyingPrice
					tableStockValue.Total = StockValue[i].Total
					updatedStockValue := db.UpdateStockValue(tableStockValue)
					_ = updatedStockValue // LOL
				} else { // data didn't exist in table, create new one then
					db.CreateStockValue(&StockValue[i])
				}
			}
			gc.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "data csv migrated successfully to StockValue table",
				"data":    StockValue,
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

// Round Go convert float to int without checking the round up/down value, it will always round down, hence we need this function
func Round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}
