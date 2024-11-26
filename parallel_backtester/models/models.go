package models

import "parallel_backtester/utils"

type Strategy interface {
	Execute(data []utils.PriceData, modelName string) (map[string]float64, BacktestResult)
}

type BacktestResult struct {
	TakeProfit       float64
	StopLoss         float64
	TotalProfit      float64
	WinRate          float64
	Trades           int
	WinningTrades    int
	AccountValues    []float64
	FinalBalance     float64
	PercentageReturn float64
}
