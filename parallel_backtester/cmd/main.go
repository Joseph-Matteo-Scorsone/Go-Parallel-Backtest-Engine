package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"parallel_backtester/models"
	"parallel_backtester/strategies"
	"parallel_backtester/utils"
	"sync"
)

func main() {

	// Define resolutions to iterate over
	resolutions := map[string]struct {
		fileName         string
		stopLossValues   []float64
		takeProfitValues []float64
	}{
		"daily": {
			"daily.csv",
			[]float64{0.03, 0.025, 0.02, 0.015, 0.01, 0.005},
			[]float64{0.07, 0.06, 0.05, 0.04, 0.03, 0.02},
		},
		"hourly": {
			"hourly.csv",
			[]float64{0.03, 0.025, 0.02, 0.015, 0.01, 0.005},
			[]float64{0.06, 0.05, 0.04, 0.03, 0.02, 0.01},
		},
		"5m": {
			"five_minutely.csv",
			[]float64{0.025, 0.02, 0.015, 0.01, 0.005, 0.0025},
			[]float64{0.025, 0.02, 0.015, 0.01, 0.005, 0.0025},
		},
	}

	for resolution, params := range resolutions {
		fmt.Printf("Running backtests for resolution: %s\n", resolution)

		// Load the historical data for the given resolution
		data, err := utils.LoadHistoricalData(params.fileName)
		if err != nil {
			log.Println("Error loading data:", err)
			os.Exit(1)
		}

		// Create a slice to store strategies
		var strategiesList []models.Strategy

		// Create a strategy for each combination of stopLoss and takeProfit values
		for _, stopLoss := range params.stopLossValues {
			for _, takeProfit := range params.takeProfitValues {
				strategy := strategies.NewMovingAverageStrategy(20, 50, stopLoss, takeProfit)
				strategiesList = append(strategiesList, strategy)
			}
		}

		// Run the backtests in parallel
		var wg sync.WaitGroup
		results := make(chan struct {
			AccountValues map[string]float64
			Result        models.BacktestResult
			ID            string
		}, len(strategiesList))

		for i, strategy := range strategiesList {
			wg.Add(1)
			go func(s models.Strategy, id int) {
				defer wg.Done()
				accountValues, result := s.Execute(data, fmt.Sprintf("Test %v", id))
				results <- struct {
					AccountValues map[string]float64
					Result        models.BacktestResult
					ID            string
				}{accountValues, result, fmt.Sprintf("Test %v", id)}
			}(strategy, i)
		}

		wg.Wait()
		close(results)

		// Write each result to a separate CSV for the resolution
		outputFileName := fmt.Sprintf("%s_backtest_results.csv", resolution)
		file, err := os.Create(outputFileName)
		if err != nil {
			log.Fatal("Error creating CSV file:", err)
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		// Write header
		writer.Write([]string{"Date", "Account Balance", "Backtest ID", "Risk"})

		for result := range results {
			for date, balance := range result.AccountValues {
				err := writer.Write([]string{date, fmt.Sprintf("%.2f", balance), result.ID, fmt.Sprintf("TP: %.2f, SL: %.2f", result.Result.TakeProfit, result.Result.StopLoss)})
				if err != nil {
					log.Fatal("Error writing to CSV:", err)
				}
			}
		}

		fmt.Printf("CSV file for resolution %s generated successfully: %s\n", resolution, outputFileName)
	}
}
