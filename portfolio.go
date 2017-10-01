package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func setupPortfolioRoutes() {
	portfolio := app.Group("/portfolio")
	{
		portfolio.POST("/buy", func(c *gin.Context) {
			var req orderRequest
			c.BindWith(&req, binding.Form)
			if req.Units < 0 {
				c.JSON(401, gin.H{"message": "Cannot buy negative units"})
				return
			}
			userId := tokenToUserId(req.Token)
			if userId == "" {
				c.JSON(401, gin.H{"message": "Unauthorized"})
				return
			}
			cash := getCash(userId)
			currentPrice, err := getStockPrice(req.Symbol)
			if err != nil {
				c.JSON(400, gin.H{"message": err.Error()})
				return
			}

			total := currentPrice * float64(req.Units)
			if cash < total {
				c.JSON(401, gin.H{"message": "Not enough money to buy"})
				return
			}

			cash -= total
			err = setCash(userId, cash)
			checkErr(err)
			updateUnitsOwned(userId, req, true)
			c.JSON(200, gin.H{
				"message":       "Successfully ordered stocks",
				"totalCost":     total,
				"remainingCash": cash,
			})
		})
		portfolio.POST("/sell", func(c *gin.Context) {
			var req orderRequest
			c.BindWith(&req, binding.Form)
			if req.Units < 0 {
				c.JSON(401, gin.H{"message": "Cannot sell negative units"})
				return
			}

			userId := tokenToUserId(req.Token)
			if userId == "" {
				c.JSON(401, gin.H{"message": "Unauthorized"})
				return
			}

			unitsOwned := getUnitsOwned(userId, req.Symbol)
			if req.Units > unitsOwned {
				c.JSON(401, gin.H{"message": "Not enough units to sell"})
				return
			}

			cash := getCash(userId)
			currentPrice, err := getStockPrice(req.Symbol)
			checkErr(err)
			total := currentPrice * float64(req.Units)
			cash += total
			err = setCash(userId, cash)
			updateUnitsOwned(userId, req, false)
			c.JSON(200, gin.H{
				"message":       "Successfully sold stocks",
				"totalCost":     total,
				"remainingCash": cash,
			})
		})

		portfolio.GET("/update/:userID", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"page": fmt.Sprintf("/update/%s", c.Param("userID")),
			})
		})
	}
}
