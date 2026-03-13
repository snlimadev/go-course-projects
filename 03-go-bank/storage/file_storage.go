package storage

import (
	"fmt"
	"os"
	"strconv"
)

func ReadFloatFromFile(fileName string) float64 {
	data, err := os.ReadFile(fileName)

	if err != nil {
		return 0
	}

	valueText := string(data)
	value, err := strconv.ParseFloat(valueText, 64)

	if err != nil {
		return 0
	}

	return value
}

func WriteFloatToFile(value float64, fileName string) error {
	valueText := fmt.Sprint(value)
	err := os.WriteFile(fileName, []byte(valueText), 0644)

	if err != nil {
		return fmt.Errorf("failed to write value to file: %w", err)
	}

	return nil
}
