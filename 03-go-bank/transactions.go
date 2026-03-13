package main

import (
	"fmt"

	"example.com/go-bank/storage"
)

func deposit(balance *float64) {
	depositAmount := getUserInput("How much do you want to deposit? ")

	if depositAmount <= 0 {
		fmt.Println("Invalid amount! Must be a number greater than 0.")
		return
	}

	newBalance := *balance + depositAmount
	err := storage.WriteFloatToFile(newBalance, balanceFileName)

	if err != nil {
		fmt.Println("An error occurred while processing the deposit.")
		return
	}

	*balance = newBalance
	fmt.Printf("Money deposited! New balance: %.2f\n", *balance)
}

func withdraw(balance *float64) {
	withdrawAmount := getUserInput("How much do you want to withdraw? ")

	if withdrawAmount <= 0 {
		fmt.Println("Invalid amount! Must be a number greater than 0.")
		return
	}

	if withdrawAmount > *balance {
		fmt.Println("Invalid amount! You can't withdraw more than you have.")
		return
	}

	newBalance := *balance - withdrawAmount
	err := storage.WriteFloatToFile(newBalance, balanceFileName)

	if err != nil {
		fmt.Println("An error occurred while processing the withdrawal.")
		return
	}

	*balance = newBalance
	fmt.Printf("Money withdrawn! New balance: %.2f\n", *balance)
}
