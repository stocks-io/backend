package main

import (
	"fmt"

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
			c.JSON(200, gin.H{
				"page": fmt.Sprintf("/stocks/%s", c.Param("symbol")),
			})
		})
	}

	portfolio := app.Group("/portfolio")
	{
		portfolio.GET("/buy", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"page": fmt.Sprintf("/buy"),
			})
		})
		portfolio.GET("/sell", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"page": fmt.Sprintf("/sell"),
			})
		})
		portfolio.GET("/update", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"page": fmt.Sprintf("/update"),
			})
		})
	}

	app.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	app.Run() // listen and serve on 0.0.0.0:8080
}
