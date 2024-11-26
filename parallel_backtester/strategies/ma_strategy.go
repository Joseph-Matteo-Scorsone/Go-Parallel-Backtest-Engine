package strategies

import (
	"parallel_backtester/models"
	"parallel_backtester/utils"
)

type MovingAverageStrategy struct {
	ShortPeriod int
	LongPeriod  int
	StopLoss    float64
	TakeProfit  float64
}

func NewMovingAverageStrategy(shortperiod, longperiod int, stopLoss, takeProfit float64) MovingAverageStrategy {
	return MovingAverageStrategy{
		ShortPeriod: shortperiod,
		LongPeriod:  longperiod,
		StopLoss:    stopLoss,
		TakeProfit:  takeProfit,
	}
}

func (s MovingAverageStrategy) Execute(data []utils.PriceData, backtestID string) (map[string]float64, models.BacktestResult) {
	var result models.BacktestResult
	var accountBalance = 10000.0
	var inPosition bool
	var entryPrice float64
	var isLong bool
	var numShares int

	accountValues := make(map[string]float64)

	for i := s.LongPeriod; i < len(data); i++ {

		if i-s.LongPeriod-1 < 0 {
			continue
		}

		shortMa := calculateMA(data[i-s.ShortPeriod:i], s.ShortPeriod)
		longMa := calculateMA(data[i-s.LongPeriod:i], s.LongPeriod)

		prevShortMa := calculateMA(data[i-s.ShortPeriod-1:i-1], s.ShortPeriod)
		prevLongMa := calculateMA(data[i-s.LongPeriod-1:i-1], s.LongPeriod)

		// Check stop loss and take profit if in position
		if inPosition {
			priceChange := data[i].Close - entryPrice
			profitLossPercent := priceChange / entryPrice

			if (isLong && profitLossPercent <= -s.StopLoss) || (!isLong && profitLossPercent >= s.StopLoss) ||
				(isLong && profitLossPercent >= s.TakeProfit) || (!isLong && profitLossPercent <= -s.TakeProfit) || 
				(isLong && prevShortMa <= prevLongMa && shortMa > longMa) || (!isLong && prevShortMa >= prevLongMa && shortMa < longMa) {
				accountBalance += priceChange * float64(numShares)
				inPosition = false
				result.TotalProfit += priceChange * float64(numShares)
				result.Trades++
				if profitLossPercent > 0 {
					result.WinningTrades++
				}
				continue
			}
		} else {
			// Determine buy/sell action
			numShares = int(accountBalance * 0.20 / data[i].Close)
			if prevShortMa <= prevLongMa && shortMa > longMa { // Buy
				entryPrice = data[i].Close
				inPosition = true
				isLong = true
			} else if prevShortMa >= prevLongMa && shortMa < longMa { // Sell
				entryPrice = data[i].Close
				inPosition = true
				isLong = false
			}
		}

		// Store the date and account balance in the map
		date := data[i].Date.Format("2006-01-02")
		accountValues[date] = accountBalance
	}

	// Calculate win rate
	if result.Trades > 0 {
		result.WinRate = float64(result.WinningTrades) / float64(result.Trades)
	}
	result.FinalBalance = accountBalance
	result.PercentageReturn = (accountBalance - 10000) / 10000 * 100

	return accountValues, result
}

func calculateMA(data []utils.PriceData, period int) float64 {
	var sum float64
	for _, price := range data {
		sum += price.Close
	}

	return sum / float64(period)
}
