package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type symbolsResp []struct {
	Symbol   string `json:"Symbol"`
	Name     string `json:"Name"`
	Industry string `json:"industry"`
}

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

type stockHistory struct {
	Date             string  `json:"date"`
	Open             float64 `json:"open"`
	High             float64 `json:"high"`
	Low              float64 `json:"low"`
	Close            float64 `json:"close"`
	Volume           int     `json:"volume"`
	UnadjustedVolume int     `json:"unadjustedVolume"`
	Change           float64 `json:"change"`
	ChangePercent    float64 `json:"changePercent"`
	Vwap             float64 `json:"vwap"`
	Label            string  `json:"label"`
	ChangeOverTime   float64 `json:"changeOverTime"`
}

func loadSymbols() symbolsResp {
	body, err := ioutil.ReadFile("./companies.min.json")
	checkFatalErr(err)
	resp := symbolsResp{}
	err = json.Unmarshal(body, &resp)
	checkFatalErr(err)
	return resp
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

func getStockHistory(symbol, timeframe string) ([]stockHistory, error) {
	url := fmt.Sprintf("https://api.iextrading.com/1.0/stock/%s/chart/%s", symbol, timeframe)
	body, err := getResponse(url)
	resp := []stockHistory{}
	if err != nil {
		return resp, err
	}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return resp, err
	}
	return resp, nil
}
