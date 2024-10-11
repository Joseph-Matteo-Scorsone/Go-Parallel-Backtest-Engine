package models

import "parallel_backtester/utils"

type Stragey interface {
	Execute(data []utils.PriceData) BacktestResult
}

type BacktestResult struct {
	TotalProfit   float64
	WinRate       float64
	Trades        int
	WinningTrades int
}
