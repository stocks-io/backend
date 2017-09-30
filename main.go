package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func checkFatalErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var app *gin.Engine
var db *sql.DB
var err error

func main() {
	app = gin.Default()
	setupDB()
	setupUserRoutes()
	setupStockRoutes()
	setupPortfolioRoutes()

	app.Run(":8080")
}
