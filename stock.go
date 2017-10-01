package main

import (
	"strings"

	finance "github.com/FlashBoys/go-finance"
)

func getStockPrice(symbol string) (float64, error) {
	sym := strings.ToUpper(symbol)
	q, err := finance.GetQuote(sym)
	checkErr(err)
	if err != nil {
		return -1, err
	}
	val, _ := q.LastTradePrice.Float64()
	return val, nil
}
