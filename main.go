package main

import (
	"fmt"
	"go-stock/marketfall"
	"go-stock/stock"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please specify which task to run: 'stock' or 'marketfall'")
		os.Exit(1)
	}

	task := os.Args[1]
	switch task {
	case "stock":
		fmt.Println("Running stock market analysis...")
		stock.RunStockAnalysis()
	case "marketfall":
		fmt.Println("Running market fall check...")
		marketfall.RunMarketFallCheck()
	default:
		fmt.Println("Invalid task. Please use 'stock' or 'marketfall'")
		os.Exit(1)
	}
}
