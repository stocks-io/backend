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

type buyRequest struct {
	UserId int    `form:"userId" json:"userId" binding:"required"`
	Units  int    `form:"units" json:"units" binding:"required"`
	Symbol string `form:"symbol" json:"symbol" binding:"required"`
}

type loginRequest struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func userExists(username string) bool {
	var exists bool
	err := db.QueryRow("SELECT IF(COUNT(*),'true','false') FROM userinfo WHERE username=?", username).Scan(&exists)
	checkErr(err)
	return exists
}

func setupDB() {
	var cmd string
	cmd = `
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
	db, err = sql.Open("mysql", "root@/stocks")
	if err != nil {
		log.Fatal("Failed to load MySQL Database: %s", err.Error())
	}
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
			checkErr(err)
			now := time.Now().Unix()
			then := time.Now().Add(time.Hour * 24).Unix()
			stmt, err := db.Prepare("INSERT sessions SET user_id=?, token=?, added=?, expires=?")
			checkErr(err)
			_, err = stmt.Exec(id, token, now, then)
			checkErr(err)
			c.JSON(200, gin.H{
				"userId":  id,
				"token":   token,
				"expires": then,
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
