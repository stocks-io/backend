package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func setupPortfolioRoutes() {
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
}
