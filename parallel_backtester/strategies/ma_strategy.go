package strategies

import (
	"parallel_backtester/models"
	"parallel_backtester/utils"
)

type MovingAverageStrategy struct {
	ShortPeriod int
	LongPeriod  int
}

func (s MovingAverageStrategy) Execute(data []utils.PriceData) models.BacktestResult {
	var result models.BacktestResult
	var inPosition bool
	var entryPrice float64

	for i := s.LongPeriod; i < len(data); i++ {
		shortMa := calculateMA(data[i-s.ShortPeriod:i], s.ShortPeriod)
		longMa := calculateMA(data[i-s.LongPeriod:i], s.LongPeriod)

		if shortMa > longMa && !inPosition { //Long signal
			entryPrice = data[i].Close
			inPosition = true
		} else if shortMa < longMa && inPosition { //Exit Long
			result.TotalProfit += data[i].Close - entryPrice
			if data[i].Close > entryPrice {
				result.WinningTrades++
			}
			result.Trades++
			inPosition = false

		} else if shortMa < longMa && !inPosition { //Short signal
			entryPrice = data[i].Close
			inPosition = true
		} else if shortMa > longMa && inPosition { //Exit Short
			result.TotalProfit += entryPrice - data[i].Close
			if data[i].Close < entryPrice {
				result.WinningTrades++
			}
			result.Trades++
			inPosition = false
		}
	}

	if result.Trades > 0 {
		result.WinRate = float64(result.WinningTrades) / float64(result.Trades)
	}

	return result
}

func calculateMA(data []utils.PriceData, period int) float64 {
	var sum float64
	for _, price := range data {
		sum += price.Close
	}

	return sum / float64(period)
}
