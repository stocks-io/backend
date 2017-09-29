package main

import (
	"fmt"
	"strings"

	finance "github.com/FlashBoys/go-finance"
	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.Default()

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

	stocks := app.Group("/stocks")
	{
		stocks.GET("/:symbol", func(c *gin.Context) {
			sym := strings.ToUpper(c.Param("symbol"))
			q, _ := finance.GetQuote(sym)
			c.JSON(200, gin.H{
				"page":  fmt.Sprintf("/stocks/%s", sym),
				"quote": q,
			})
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

	app.Run() // listen and serve on 0.0.0.0:8080
}
