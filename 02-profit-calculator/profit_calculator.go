package main

import (
	"fmt"
	"os"
)

func main() {
	revenue := getUserInput("Revenue: ")
	expenses := getUserInput("Expenses: ")
	taxRate := getUserInput("Tax Rate: ")

	ebt, profit, ratio := calculateFinancials(revenue, expenses, taxRate)

	fmt.Printf("\nEarnings Before Tax (EBT): %.2f\n", ebt)
	fmt.Printf("Earnings After Tax (Profit): %.2f\n", profit)
	fmt.Printf("Ratio: %.2f\n", ratio)

	storeResults(ebt, profit, ratio)

	fmt.Println("\nPress ENTER to exit...")
	fmt.Scanln()
	fmt.Scanln()
}

func getUserInput(infoText string) float64 {
	var userInput float64

	fmt.Print(infoText)
	fmt.Scan(&userInput)

	if userInput <= 0 {
		panic("The value must be a positive number.")
	}

	return userInput
}

func calculateFinancials(revenue, expenses, taxRate float64) (float64, float64, float64) {
	ebt := revenue - expenses
	profit := ebt * (1 - taxRate/100)
	ratio := ebt / profit

	return ebt, profit, ratio
}

func storeResults(ebt, profit, ratio float64) {
	results := fmt.Sprintf("EBT: %.2f\nProfit: %.2f\nRatio: %.2f\n", ebt, profit, ratio)
	os.WriteFile("results.txt", []byte(results), 0644)
}
