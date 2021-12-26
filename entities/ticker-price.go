package entities

import "time"

type TickerPrice struct {
	Ticker string
	Time   time.Time
	Price  string // decimal value. example: "0", "10", "12.2", "13.2345122"
	Err    error
}
