package main

import (
	"fmt"
	"github.com/paulschick/eodhd"
	"log"
	"os"
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

	tickers, res, err := client.TickerService.GetTickers("US", eodhd.GetFormatCsv())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("tickers response code: ", res.StatusCode)
	for i, t := range tickers {
		fmt.Printf("(%d) %s - %s\n", i, t.Name, t.Code)
	}
}
