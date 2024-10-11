package main

import (
	"fmt"
	"parallel_backtester/models"
	"parallel_backtester/strategies"
	"parallel_backtester/utils"
	"sync"
)

func main() {
	data, err := utils.LoadHistoricalData("SPY_10_years.csv")
	if err != nil {
		fmt.Println("Error loading data", err)
	}

	strategies := []models.Stragey{
		strategies.MovingAverageStrategy{ShortPeriod: 10, LongPeriod: 30},
		strategies.MovingAverageStrategy{ShortPeriod: 20, LongPeriod: 50},
	}

	var wg sync.WaitGroup
	results := make(chan models.BacktestResult)

	for _, strategy := range strategies {
		wg.Add(1)
		go func(s models.Stragey){
			defer wg.Done()
			result := s.Execute(data)
			results <- result
		}(strategy)
	}

	go func() {
		wg.Wait()
		close(results)
	} ()

	for result := range results {
		fmt.Printf("Profit: %.2f, Win Rate %.2f, Trades %d\n",
		result.TotalProfit, result.WinRate, result.Trades)
	}

}
