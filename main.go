package main

import (
	"database/sql"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	finance "github.com/FlashBoys/go-finance"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
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

func userExists(username string) bool {
	var exists bool
	err := db.QueryRow("SELECT IF(COUNT(*),'true','false') FROM userinfo WHERE username=?", username).Scan(&exists)
	checkErr(err)
	return exists
}

func setupDB() {
	db, err = sql.Open("mysql", "root@/")
	checkFatalErr(err)
	_, err := db.Exec("CREATE DATABASE IF NOT EXISTS stocks")
	checkFatalErr(err)
	db, err = sql.Open("mysql", "root@/stocks")
	checkFatalErr(err)
	cmd := `
    CREATE TABLE IF NOT EXISTS userinfo
    (
      id              	int unsigned NOT NULL auto_increment,
      first_name		varchar(255) NOT NULL,
      last_name			varchar(255) NOT NULL,
      username			varchar(255) NOT NULL UNIQUE,
      email         	varchar(255) NOT NULL UNIQUE,
      password         	varchar(255) NOT NULL,
      added           	varchar(255) NOT NULL,
      PRIMARY KEY    	(id)
    );
    `
	stmt, err := db.Prepare(cmd)
	checkErr(err)
	_, err = stmt.Exec()
	checkErr(err)
	cmd = `
    CREATE TABLE IF NOT EXISTS sessions
    (
      id              	int unsigned NOT NULL auto_increment,
      user_id			int unsigned NOT NULL,
      token				varchar(255),
      added           	varchar(255) NOT NULL,
      expires           varchar(255) NOT NULL,
      PRIMARY KEY    	string(id)
    );
    `
	stmt, err = db.Prepare(cmd)
	checkErr(err)
	_, err = stmt.Exec()
	checkErr(err)
	cmd = `
    CREATE TABLE IF NOT EXISTS portfolio
    (
      id              	int unsigned NOT NULL auto_increment,
      user_id          	int unsigned NOT NULL UNIQUE,
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
	setupDB()
	users := app.Group("/users")
	{
		users.POST("/login", func(c *gin.Context) {
			var id int
			var req loginRequest
			c.BindWith(&req, binding.Form)
			if !userExists(req.Username) {
				c.JSON(401, gin.H{"message": "user does not exist"})
				return
			}
			rows, err := db.Query("SELECT id FROM userinfo WHERE username=?", req.Username)
			checkErr(err)
			rows.Next()
			err = rows.Scan(&id)
			checkErr(err)
			token, err := exec.Command("uuidgen").Output()
			token = token[0 : len(token)-1]
			checkErr(err)
			now := time.Now().Unix()
			then := time.Now().Add(time.Hour * 24).Unix()
			stmt, err := db.Prepare("INSERT sessions SET user_id=?, token=?, added=?, expires=?")
			checkErr(err)
			_, err = stmt.Exec(id, token, now, then)
			checkErr(err)
			fmt.Printf("%s\n", token)
			c.JSON(200, gin.H{
				"userId":  id,
				"token":   fmt.Sprintf("%s", token),
				"expires": then,
			})
		})
		users.POST("/logout", func(c *gin.Context) {
			var req logoutRequest
			c.BindWith(&req, binding.Form)
			var exists bool
			err := db.QueryRow("SELECT IF(COUNT(*),'true','false') FROM sessions WHERE token=?", req.Token).Scan(&exists)
			checkErr(err)
			if !exists {
				c.JSON(401, gin.H{
					"message": "token does not exist",
				})
				return
			}
			_, err = db.Exec("DELETE FROM sessions WHERE token=?", req.Token)
			checkErr(err)
			c.JSON(200, gin.H{
				"message": "successfully logged out",
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
			c.JSON(200, gin.H{"user": "buy"})

		})
		portfolio.POST("/sell", func(c *gin.Context) {
			c.JSON(200, gin.H{"user": "sell"})
		})
		portfolio.GET("/update/:userID", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"page": fmt.Sprintf("/update/%s", c.Param("userID")),
			})
		})
	}

	app.Run(":8080")
}
