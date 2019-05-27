package api

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/rinaldypasya/TokoIjah/inventory"
)

// SaleReportReqBody data structure for generate SaleReport json req body
type SaleReportReqBody struct {
	From      string `json:"datefrom"`
	To        string `json:"dateto"`
	Csvexport string `json:"exportcsv"`
}

// CreateSaleReport API to create one sale report instance and save it to SaleReports table
func CreateSaleReport(db inventory.InventSaleReport, dbstockvalue inventory.InventStockvalue) gin.HandlerFunc {
	return func(gc *gin.Context) {

		var SaleReport inventory.SaleReport
		var stockvalue inventory.StockValue

		if gc.BindJSON(&SaleReport) == nil {
			SaleReport.Timestamp = time.Now().Format("2006-01-02 15:04:05")
			stockvalue = dbstockvalue.GetStockValuesBySku(SaleReport.Sku)
			if stockvalue.Sku != "" {
				SaleReport.Buyingprice = stockvalue.BuyingPrice
			} else {
				gc.JSON(http.StatusBadRequest, gin.H{
					"status":  false,
					"message": "There is no such product in stockvalue!",
				})
				return
			}
			SaleReport.Total = (SaleReport.Amount * SaleReport.Saleprice)
			SaleReport.Profit = (SaleReport.Amount * SaleReport.Saleprice) - (SaleReport.Amount * SaleReport.Buyingprice)
			db.CreateSaleReport(&SaleReport)
			gc.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "Sale report is created successfully!",
				"id":      SaleReport.ID,
				"sale":    SaleReport,
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

// GetAllSaleReports get all rows of SaleReports from SaleReports table
func GetAllSaleReports(db inventory.InventSaleReport) gin.HandlerFunc {
	return func(gc *gin.Context) {

		var SaleReport []inventory.SaleReport
		SaleReport = db.GetAllSaleReports()
		if len(SaleReport) > 0 {
			gc.JSON(http.StatusOK, gin.H{
				"status": "true",
				"data":   SaleReport,
			})
			return
		} else {
			gc.JSON(http.StatusOK, gin.H{
				"status":  "true",
				"message": "No sale report yet!",
			})
			return
		}

		return
	}
}

// GetSaleReportsBySKU get SaleReports by sku
func GetSaleReportsBySKU(db inventory.InventSaleReport) gin.HandlerFunc {
	return func(gc *gin.Context) {

		var SaleReport []inventory.SaleReport
		SaleReport = db.GetSaleReportsBySKU(gc.Param("sku"))
		if SaleReport[0].Sku == gc.Param("sku") {
			gc.JSON(http.StatusOK, gin.H{
				"status": "true",
				"data":   SaleReport,
			})
			return
		} else {
			gc.JSON(http.StatusNotFound, gin.H{
				"status": "false",
				"data":   "No SaleReports not found!",
			})
			return
		}
	}
}

// SaleReportExportToCSV export SaleReports data into csv file
func SaleReportExportToCSV(db inventory.InventSaleReport) gin.HandlerFunc {
	return func(gc *gin.Context) {

		var allSaleReport []inventory.SaleReport
		allSaleReport = db.GetAllSaleReports()

		csvdata := init2dArray(len(allSaleReport), 10)

		for i := 0; i < len(allSaleReport); i++ {
			csvdata[i][0] = strconv.Itoa(i + 1)
			csvdata[i][1] = allSaleReport[i].OrderID
			csvdata[i][2] = allSaleReport[i].Timestamp
			csvdata[i][3] = allSaleReport[i].Sku
			csvdata[i][4] = allSaleReport[i].Name
			csvdata[i][5] = strconv.Itoa(allSaleReport[i].Amount)
			csvdata[i][6] = strconv.Itoa(allSaleReport[i].Saleprice)
			csvdata[i][7] = strconv.Itoa(allSaleReport[i].Total)
			csvdata[i][8] = strconv.Itoa(allSaleReport[i].Buyingprice)
			csvdata[i][9] = strconv.Itoa(allSaleReport[i].Profit)
		}

		fileName := time.Now().Format("2006-01-02") + "-SaleReport.csv"
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
			"message":  "SaleReport data exported to csv successfully!",
			"filename": fileName,
		})
		return
	}
}

// GetSaleReportsByDate get SaleReports by range of date(timestamp)
func GetSaleReportsByDate(db inventory.InventSaleReport) gin.HandlerFunc {
	return func(gc *gin.Context) {

		var reqBody SaleReportReqBody
		var from string
		var to string
		var csvexport string
		decoder := json.NewDecoder(gc.Request.Body)
		err := decoder.Decode(&reqBody)
		from = reqBody.From
		to = reqBody.To
		csvexport = reqBody.Csvexport
		if err != nil {
			from = "undefined"
			to = "undefined"
		}

		var SaleReports []inventory.SaleReport
		SaleReports = db.GetSaleReportsByDate(from, to)
		if len(SaleReports) > 0 {
			if csvexport == "1" {
				// export to csv
				csvdata := init2dArray(len(SaleReports), 10)
				for i := 0; i < len(SaleReports); i++ {
					csvdata[i][0] = strconv.Itoa(i + 1)
					csvdata[i][1] = SaleReports[i].OrderID
					csvdata[i][2] = SaleReports[i].Timestamp
					csvdata[i][3] = SaleReports[i].Sku
					csvdata[i][4] = SaleReports[i].Name
					csvdata[i][5] = strconv.Itoa(SaleReports[i].Amount)
					csvdata[i][6] = strconv.Itoa(SaleReports[i].Saleprice)
					csvdata[i][7] = strconv.Itoa(SaleReports[i].Total)
					csvdata[i][8] = strconv.Itoa(SaleReports[i].Buyingprice)
					csvdata[i][9] = strconv.Itoa(SaleReports[i].Profit)
				}

				fileName := time.Now().Format("2006-01-02") + "-SaleReport_" + dateOnlyFormat(from) + "_" + dateOnlyFormat(to) + ".csv"
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
					"message":  "SaleReport exported successfully!",
					"filename": fileName,
				})
				return

			} else {
				// just print the json response
				omzet := 0
				grossprofit := 0
				totalsale := 0
				totalproducts := 0
				for i := 0; i < len(SaleReports); i++ {
					omzet += SaleReports[i].Total
					grossprofit += SaleReports[i].Profit
					if SaleReports[i].OrderID != "" {
						totalsale++
					}
					totalproducts += SaleReports[i].Amount
				}

				gc.JSON(http.StatusOK, gin.H{
					"status":            true,
					"datereport":        time.Now().Format("2006-01-02"),
					"daterange":         from + " - " + to,
					"omzet":             omzet,
					"grossprofit":       grossprofit,
					"totalsale":         totalsale,
					"totalsoldproducts": totalproducts,
					"data":              SaleReports,
				})
				return
			}
		} else {
			gc.JSON(http.StatusNotFound, gin.H{
				"status":  false,
				"message": "No sale reports by those dates!",
			})
			return
		}

		// in case nothing to be returned
		gc.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "something wrong!",
			"data":    SaleReports,
		})
		return
	}
}

// SaleReportImportCSV import csv data into SaleReports table
func SaleReportImportCSV(db inventory.InventSaleReport) gin.HandlerFunc {
	return func(gc *gin.Context) {

		var SaleReport []inventory.SaleReport

		file, _ := gc.FormFile("SaleReportimport")
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

			SaleReportamount, _ := strconv.Atoi(line[5])
			SaleReportsaleprice, _ := strconv.Atoi(line[6])
			SaleReporttotal, _ := strconv.Atoi(line[7])
			SaleReportbuyingprice, _ := strconv.Atoi(line[8])
			SaleReportprofit, _ := strconv.Atoi(line[9])
			SaleReport = append(SaleReport, inventory.SaleReport{
				OrderID:     line[1],
				Timestamp:   line[2], // start from timestamp column as we ignore id column(assume the csv include the IDs)
				Sku:         line[3],
				Name:        line[4],
				Amount:      SaleReportamount,
				Saleprice:   SaleReportsaleprice,
				Total:       SaleReporttotal,
				Buyingprice: SaleReportbuyingprice,
				Profit:      SaleReportprofit,
			})
		}

		if len(SaleReport) > 0 {
			for i := 0; i < len(SaleReport); i++ {
				db.CreateSaleReport(&SaleReport[i])
			}
			gc.JSON(http.StatusOK, gin.H{
				"status":  true,
				"message": "data csv migrated successfully to SaleReport table",
				"data":    SaleReport,
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

func dateOnlyFormat(t string) string {
	runes := []rune(t)
	return string(runes[0:10])
}
