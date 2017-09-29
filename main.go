package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	finance "github.com/FlashBoys/go-finance"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func setupDB() {
	var cmd string
	cmd = `
    CREATE TABLE IF NOT EXISTS users
    (
      id              	int unsigned NOT NULL auto_increment,
      first_name		varchar(255) NOT NULL,
      last_name			varchar(255) NOT NULL,
      email         	varchar(255) NOT NULL,
      password         	varchar(255) NOT NULL,
      added           	datetime NOT NULL,
      PRIMARY KEY    	(id)
    );
    `
	stmt, err := db.Prepare(cmd)
	checkErr(err)
	_, err = stmt.Exec()
	checkErr(err)
	cmd = `
    CREATE TABLE IF NOT EXISTS portfolio
    (
      id              	int unsigned NOT NULL auto_increment,
      user_id          	int unsigned NOT NULL,
      cash				FLOAT(8),
      PRIMARY KEY     	(id)
    );
    `
	stmt, err = db.Prepare(cmd)
	checkErr(err)
	_, err = stmt.Exec()
	checkErr(err)
	cmd = `
    CREATE TABLE IF NOT EXISTS positions
    (
      id              	int unsigned NOT NULL auto_increment,
      user_id          	int unsigned NOT NULL,
      symbol			varchar(32),
      units				int unsigned NOT NULL,
      buy_price			FLOAT(8),
      PRIMARY KEY     	(id)
    );
    `
	stmt, err = db.Prepare(cmd)
	checkErr(err)
	_, err = stmt.Exec()
	checkErr(err)
}

var db *sql.DB
var err error

func main() {
	app := gin.Default()
	db, err = sql.Open("mysql", "root@/stocks")
	if err != nil {
		log.Fatal("Failed to load MySQL Database: %s", err.Error())
	}
	setupDB()
	users := app.Group("/users")
	{
		users.GET("/login", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"page": "login",
			})
		})
		users.GET("/logout", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"page": "logout",
			})
		})
		users.GET("/register", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"page": "register",
			})
		})
	}

	stocks := app.Group("/stock")
	{
		stocks.GET("/:symbol", func(c *gin.Context) {
			sym := strings.ToUpper(c.Param("symbol"))
			q, _ := finance.GetQuote(sym)
			c.JSON(200, q)
		})
	}

	portfolio := app.Group("/portfolio")
	{
		portfolio.POST("/buy", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"page": fmt.Sprintf("/buy"),
			})
		})
		portfolio.POST("/sell", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"page": fmt.Sprintf("/sell"),
			})
		})
		portfolio.GET("/update/:userID", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"page": fmt.Sprintf("/update/%s", c.Param("userID")),
			})
		})
	}

	app.Run(":8080")
}
