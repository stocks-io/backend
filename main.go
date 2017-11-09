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
var err error

func main() {
	app = gin.Default()
	db = setupDB("stocks")
	if os.Getenv("DB_MOCK") != "" {
		mockDB := setupDB("stocks_mock")
		mockData(mockDB)
		db = mockDB
		log.Println("Using mock database!")
	}
	setupUserRoutes()
	setupPortfolioRoutes()

	app.Run(":8080")
}
