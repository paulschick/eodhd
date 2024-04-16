package main

import (
	"eodhd"
	"fmt"
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
	lastExch := exchanges[len(exchanges)-1]
	fmt.Println("name: ", lastExch.Name)
	fmt.Println("currency: ", lastExch.Currency)
}
