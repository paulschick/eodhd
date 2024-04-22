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

	exchanges, res, err := client.ExchangesService.GetExchanges()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Exchanges request status: %d\n", res.StatusCode)
	fmt.Printf("returned %d exchanges\n", len(exchanges))

	for i, exchange := range exchanges {
		fmt.Printf("(%d) %s, %s, %s\n", i, exchange.Name, exchange.Country, exchange.Code)
	}
}
