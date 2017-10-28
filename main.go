package main

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var app *gin.Engine
var db *sql.DB
var err error

func main() {
	app = gin.Default()
	setupDB()
	mockData()
	//setupUserRoutes()
	//setupPortfolioRoutes()

	//app.Run(":8080")
}
