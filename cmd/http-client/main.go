package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bharath/go-http-client"
)

func main() {
	httpClient := client.NewDefaultClient()
	data, err := client.FetchData(httpClient)
	if err != nil {
		log.Fatalf("Error fetching data: %v", err)
	}

	fmt.Println("Successfully fetched data:")
	fmt.Println(string(data))
	os.Exit(0)
}
