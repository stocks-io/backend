package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var app *gin.Engine
var db *sql.DB
var mockDB *sql.DB
var err error

func main() {
	app = gin.Default()
	db = setupDB("stocks")
	mockDB = setupDB("stocks_mock")
	if os.Getenv("MOCK_DB") != "" {
		mockData(mockDB)
	}
	if os.Getenv("USE_MOCK") != "" {
		db = mockDB
		log.Println("Using mock database!")
	}
	setupUserRoutes()
	setupPortfolioRoutes()

	app.Run(":8080")
}
