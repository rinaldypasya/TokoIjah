# TokoIjah Inventory Management Service

## Web Service for Handling Simple Inventory Management using Go

## About

This service built with scalability and microservice architecture in mind. The structure of the project was built to ease the development especially in microservices based app. The project consist of layers for abstracting the certain functionalities. In Go projects, these layers are presented as packages. The packages wrap layers functionalities and are accessible to other packages/layers. The packages/layers are:

- Model - Contains the database initialization and configuration in `db.go`. Provide logics for database access(database I/O) and serve the data to/from database.
- Inventory - Contains data structures and interfaces to model the construction of database along with their logics. This package will be become main interface for requests/responds coming in/out to/from database.
- API - Contains the router initialization and server configurations in `router.go`. Provide all APIs to handle requests coming from clients. These APIs will interact with domain as interface for accessing the data which is handled by model. 
- Main - Main package is the root package that act as program executor that starts the program. This package calls API package to initialize the router and starts the server runtime. As program running, while this main package started, it will create server coming from router initialization in API package, which enables APIs package(the router) to be able to handle requests/responds from and to clients.

In short, here's how it works when executing the project:
```
// main calls router in API package to start the server runtime
main -> API.router

// order of accepting requests, for responses work the other way around
requests -> API.router -> API.<endpoints> -> inventory.<interfaces> -> model.<modelname> -> model.db -> database

```

## Requirements and Depedencies

This service at the first time of development used Go1.9.3. Other depedencies/libraries are used to ease the development. Here are the main external depedencies being used:

- [gin](https://github.com/gin-gonic/gin) - Used for routing and handling APIs requests/responses.
- [gin-cors](https://github.com/itsjamie/gin-cors) - Used for configuration and handling middleware.
- [gorm](https://github.com/jinzhu/gorm) - Used for ORM.
- [sqlite3](https://github.com/mattn/go-sqlite3) - SQLite driver for golang.

## Installation Guide

Before installing and running this app, make sure you have install [Go](https://golang.org/doc/install). After Go is installed, Follow these steps:

### Linux/Mac
- Set your [gopath](https://til.codes/how-do-i-set-the-gopath-environment-variable-on-ubuntu/) somewhere in your filesystem, preferably in /home.
- After setting up gopath, create 3 main folder inside gopath. For example, if your gopath is /home/user/go, then create these 3 folder under that go folder:
    - bin/
    - pkg/
    - src/
- Then, open terminal inside your gopath(e.g /home/user/go). Then get this app by using commands:
    >`go get github.com/rinaldypasya/TokoIjah`
- The codebase for the app will be downloaded inside `src/` folder inside your gopath.
- After codebase downloaded, change directory to that app -> e.g 
    >`cd $GOPATH/src/github.com/rinaldypasya/TokoIjah`
- Install/Build the app from source to build the executable and packages object by using this command
    >`make build` (assuming you're already inside the app folder path)

- After finishing installing/building run the app by calling the name of app folder(where main package reside), in this case, if you're inside the app path, just call the name of the app:
    >`./TokoIjah`
    
- If program installed and run successfully, the routing informations will be showed on the terminal showing all API endpoints registered inside router. The port used is :8080, so access it by calling 

    >`localhost:8080/<endpoints>`

### Windows

For installation on Windows, after installing the Go you're required to install **MinGW-64 GCC** because some of depedencies depend and compiled only by gcc 64 bit version, and on Windows you have to use [MinGW64](https://sourceforge.net/projects/mingw-w64/). Remember to download, install, and activate the **64 bit version**.

After installing Go and MinGW64, set [gopath in windows](https://github.com/golang/go/wiki/SettingGOPATH#windows).

After setting up the gopath, make sure the gopath set successfully by calling this command in command prompt:
    >go env

If the gopath set and match with the one you set, the rest of the steps are the same with the Linux/Mac. Follow them one by one and you're all set.

### Rebuilding and Cleaning

Please typing command in project root folder:

- Clean

    >`make clean`

- Rebuild and clean

    >`make rebuild`


## Database and CSV

This app use SQLite 3 database and the file reside inside the app's codebase root folder named `stock.db`. 

>The mechanism of creating the tables is by using the auto-migration feature provided by `gin` which will migrate the database schema into `stock.db`. 

So there's no need to set up the table manually since all the table structures are modeled inside the codebase, specifically inside the domain package.The file inside the domain package, with their correspond data structure will provide the mapping for creating table structure inside database and is executed while first time initializing and auto-migrating the data structure.

This app also provide functionalities to export/import csv data into/from database. This functionalities are available by accessing the API endpoints for export and import csv.

## Business model

Here is the business model for TokoIjah, let say the owner is called Ijah:

- When Ijah wants to record products coming **into** the stock, Ijah uses API endpoint:
    >StoreProduct POST: `/stockin`
- Storing products will affect `stock` table. If the stored products were already existed before, it will increment the amount of that existing product. If not, then create new one.
- When ijah wants to records products coming **out** from the stock, ijah uses API endpoint:
    >RemoveProduct POST: `/stockout`
- Removing products will affects the table stock. If the stored products were already existed before, it will reduce the amount of that existing product. If not, then it's impossible to remove the nonexistent.
- If Ijah wants to generate the values of her stock(stockvalues), Ijah uses API endpoint:
    >GetAllStockValues GET: `/stockvalue`

    If Ijah wants to calculate the pricing and total based on records from products coming into the stock, use:
    >GenerateStockValue GET:`/generatestockvalue`
- If Ijah wants to record the product sales, Ijah used:
    >CreateSaleReport POST:`/salereport`

    If Ijah wants to generate report of sales from certain range of dates, use:
    >GenerateSaleReport POST:`/generatesalereport`
    with the option to **export** it to csv.
- Ijah can also export all the data from stocks, stored products, removed products, stock values, and sale report into CSV formatted data by using the API endpoints providing export csv for each data domain(explore the APIs using postman)
- Ijah can also import existing CSV data into the database(data migration) by using API endpoints providing import from csv for each data domain(explore the APIs using postman)

You can explore all the API endpoints by using postman or insomnia to see all the functionalities provided by the service.

---
That's all. Thank you and if you have any problems or questions, reach me an email at **rinaldypasya@gmail.com**.