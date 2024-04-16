package main

import (
	"eodhd"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		panic("no API key provided")
	}

	client, err := eodhd.NewClient(apiKey)
	if err != nil {
		log.Fatalf("error creating client: %s", err)
	}

	from := time.Date(2024, 4, 1, 1, 0, 0, 0, time.Local)

	params := &eodhd.UrlParams{
		ApiToken: apiKey,
		Format:   eodhd.GetFormatCsv(),
		FromTime: &from,
		Symbol:   "AAPL",
	}

	data, resp, err := client.OhlcvService.GetCandles(params)
	if err != nil {
		log.Fatalf("error retrieving candles: %s", err)
	}

	fmt.Printf("Status code: %d\n", resp.StatusCode)
	fmt.Printf("retrieved %d candles for %s\n", len(data), params.Symbol)

	lastCandle := data[len(data)-1]
	_ = lastCandle.ParseDate()
	fmt.Println(lastCandle)
}
