package main

import (
	"fmt"
	"math"
)

func main() {
	const inflationRate float64 = 2.5
	var investmentAmount, expectedReturnRate, years float64

	fmt.Print("Investment Amount: ")
	fmt.Scan(&investmentAmount)

	fmt.Print("Expected Return Rate: ")
	fmt.Scan(&expectedReturnRate)

	fmt.Print("Years: ")
	fmt.Scan(&years)

	futureValue := investmentAmount * math.Pow(1+expectedReturnRate/100, years)
	futureRealValue := futureValue / math.Pow(1+inflationRate/100, years)

	fmt.Printf("\nFuture Value: %.2f\n", futureValue)
	fmt.Printf("Future Value (Adjusted for Inflation): %.2f\n", futureRealValue)

	fmt.Println("\nPress ENTER to exit...")
	fmt.Scanln()
	fmt.Scanln()
}
