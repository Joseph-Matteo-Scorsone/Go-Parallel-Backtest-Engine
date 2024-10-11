package utils

import (
	"encoding/csv"
	"os"
	"strconv"
	"time"
)

type PriceData struct {
	Date   time.Time
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume int64
}

func LoadHistoricalData(filename string) ([]PriceData, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var data []PriceData
	for i, record := range records {
		if i == 0 {
			continue // skip header
		}

		date, err := time.Parse("2006-01-02", record[0])
		if err != nil {
			return nil, err
		}

		open, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			return nil, err
		}

		high, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			return nil, err
		}

		low, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			return nil, err
		}

		close, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			return nil, err
		}


		volume, err := strconv.ParseInt(record[5], 10, 64)
		if err != nil {
			return nil, err
		}

		data = append(data, PriceData{date, open, high, low, close, volume})
	}

	return data, nil
}
