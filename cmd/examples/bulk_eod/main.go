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

	data, resp, err := client.BulkEodService.GetBulkEod(eodhd.GetPtrString("US"), eodhd.GetFormatCsv())
	if err != nil {
		log.Fatalf("error getting bulk EOD response: %s", err)
	}

	if resp.StatusCode != 200 {
		log.Fatalf("Received response: %d %s", resp.StatusCode, resp.Status)
	}

	d1 := data[0]
	fmt.Println("first record")
	fmt.Printf("ticker %s - Close %f\n", d1.Code, d1.Close)
}
