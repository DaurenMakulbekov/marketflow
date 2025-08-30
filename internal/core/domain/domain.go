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

type Price struct {
	Price float64
}

var ErrorNotFound = errors.New("Not Found")
var ErrorBadRequest = errors.New("Incorrect input")
