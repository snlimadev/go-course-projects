package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var scanner = bufio.NewScanner(os.Stdin)

func getUserInput(infoText string) float64 {
	fmt.Print(infoText)
	scanner.Scan()

	input := strings.TrimSpace(scanner.Text())
	value, err := strconv.ParseFloat(input, 64)

	if err != nil {
		return 0
	}

	return value
}
