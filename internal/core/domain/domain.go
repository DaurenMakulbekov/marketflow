package domain

import (
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
