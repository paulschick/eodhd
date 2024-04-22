package main

import (
	"fmt"
	"github.com/paulschick/eodhd"
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

	data, res, err := client.OhlcvService.GetOhlcv("AAPL", eodhd.GetPtrString("US"), eodhd.GetFormatJson(), eodhd.GetPtrTime(from), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("status: ", res.StatusCode)
	fmt.Println("number of records returned: ", len(data))
}
