package domain

import (
	"errors"
	"time"
)

type Exchanges struct {
	PairName  string
	Exchange  string
	Timestamp time.Time
	AvgPrice  float64
	MinPrice  float64
	MaxPrice  float64
}

type Exchange struct {
	ID        string  `json:"id"`
	Exchange  string  `json:"exchange"`
	Symbol    string  `json:"symbol"`
	Price     float64 `json:"price"`
	Timestamp int64   `json:"timestamp"`
}

type PriceSymbol struct {
	Symbol string
	Price  float64
}

type PriceExchangeSymbol struct {
	Exchange string
	Symbol   string
	Price    float64
}

type SystemStatus struct {
	Redis     string `json:"redis"`
	Exchange1 string `json:"exchange1"`
	Exchange2 string `json:"exchange2"`
	Exchange3 string `json:"exchange3"`
}

var ErrorNotFound = errors.New("Not Found")
var ErrorBadRequest = errors.New("Incorrect input")
