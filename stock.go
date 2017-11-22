package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type stocksResp struct {
	AvgTotalVolume   int     `json:"avgTotalVolume"`
	CalculationPrice string  `json:"calculationPrice"`
	Change           float64 `json:"change"`
	ChangePercent    float64 `json:"changePercent"`
	Close            float64 `json:"close"`
	CloseTime        int     `json:"closeTime"`
	CompanyName      string  `json:"companyName"`
	DelayedPrice     float64 `json:"delayedPrice"`
	DelayedPriceTime int     `json:"delayedPriceTime"`
	LatestPrice      float64 `json:"latestPrice"`
	LatestSource     string  `json:"latestSource"`
	LatestTime       string  `json:"latestTime"`
	LatestUpdate     int     `json:"latestUpdate"`
	LatestVolume     int     `json:"latestVolume"`
	MarketCap        int     `json:"marketCap"`
	Open             float64 `json:"open"`
	OpenTime         int     `json:"openTime"`
	PeRatio          float64 `json:"peRatio"`
	PreviousClose    float64 `json:"previousClose"`
	PrimaryExchange  string  `json:"primaryExchange"`
	Sector           string  `json:"sector"`
	Symbol           string  `json:"symbol"`
	Week52High       float64 `json:"week52High"`
	Week52Low        float64 `json:"week52Low"`
	YtdChange        float64 `json:"ytdChange"`
}

type symbolsResp []struct {
	Symbol   string `json:"Symbol"`
	Name     string `json:"Name"`
	Industry string `json:"industry"`
}

func getStockPrice(symbol string) (float64, error) {
	url := fmt.Sprintf("https://ws-api.iextrading.com/1.0/stock/%s/quote", symbol)
	body, err := getResponse(url)
	resp := stocksResp{}
	if err != nil {
		return -1, err
	}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return -1, err
	}
	return resp.LatestPrice, nil
}

func loadSymbols() symbolsResp {
	body, err := ioutil.ReadFile("./companies.min.json")
	checkFatalErr(err)
	resp := symbolsResp{}
	err = json.Unmarshal(body, &resp)
	checkFatalErr(err)
	return resp
}
