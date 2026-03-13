package main

import (
	"fmt"

	"example.com/go-bank/storage"
)

const balanceFileName = "balance.txt"

func main() {
	balance := storage.ReadFloatFromFile(balanceFileName)
	fmt.Println("Welcome to Go Bank!")

	for {
		printMenu()
		userChoice := int(getUserInput("\nYour choice: "))

		switch userChoice {
		case 1:
			fmt.Printf("Your balance: %.2f\n", balance)
		case 2:
			deposit(&balance)
		case 3:
			withdraw(&balance)
		case 4:
			return
		default:
			fmt.Println("Invalid choice!")
		}
	}
}
