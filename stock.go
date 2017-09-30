package main

import (
	"strings"

	finance "github.com/FlashBoys/go-finance"
	"github.com/gin-gonic/gin"
)

func setupStockRoutes() {
	stocks := app.Group("/stock")
	{
		stocks.GET("/:symbol", func(c *gin.Context) {
			sym := strings.ToUpper(c.Param("symbol"))
			q, _ := finance.GetQuote(sym)
			c.JSON(200, q)
		})
	}
}
