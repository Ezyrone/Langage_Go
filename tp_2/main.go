package main

import (
	"errors"
	"fmt"
)

// Exercise 1
func ComputeBasicStats(numbers ...int) (int, int, float64) {
	if len(numbers) == 0 {
		return 0, 0, 0.0
	}

	sum := 0
	for _, n := range numbers {
		sum += n
	}

	count := len(numbers)
	average := float64(sum) / float64(count)

	return sum, count, average
}

// Exercise 2
func ComputeFullStats(numbers ...float64) (float64, float64, float64, float64, int, error) {
	if len(numbers) == 0 {
		return 0, 0, 0, 0, 0, errors.New("no arguments provided")
	}

	min := numbers[0]
	max := numbers[0]
	sum := 0.0

	for _, n := range numbers {
		sum += n
		if n < min {
			min = n
		}
		if n > max {
			max = n
		}
	}

	count := len(numbers)
	average := sum / float64(count)

	return min, max, sum, average, count, nil
}

// Exercise 3
func AnalyzeSensorData(readings ...float64) (float64, float64, float64, int, int, error) {
	var valid []float64
	invalid := 0

	for _, r := range readings {
		if r > 0.0 && r <= 100.0 {
			valid = append(valid, r)
		} else {
			invalid++
		}
	}

	if len(valid) == 0 {
		return 0, 0, 0, 0, invalid, errors.New("no valid readings found")
	}

	min, max, _, average, validCount, _ := ComputeFullStats(valid...)

	return min, max, average, validCount, invalid, nil
}

func main() {

	fmt.Println("Exercise 1")

	sum, count, average := ComputeBasicStats(10, 20, 30, 40)
	fmt.Printf("Sum: %d, Count: %d, Average: %.2f\n", sum, count, average)

	sumEmpty, countEmpty, averageEmpty := ComputeBasicStats()
	fmt.Printf("Sum (empty): %d, Count (empty): %d, Average (empty): %.2f\n", sumEmpty, countEmpty, averageEmpty)

	sumOne, countOne, averageOne := ComputeBasicStats(42)
	fmt.Printf("Sum (single): %d, Count (single): %d, Average (single): %.2f\n", sumOne, countOne, averageOne)

	fmt.Println("Exercise 2")

	min, max, total, avg, cnt, err := ComputeFullStats(1.5, 2.8, 0.7, 3.1)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Printf("Min: %.2f, Max: %.2f, Sum: %.2f, Average: %.2f, Count: %d\n", min, max, total, avg, cnt)
	}

	_, _, _, _, _, errEmpty := ComputeFullStats()
	if errEmpty != nil {
		fmt.Println("Error for empty arguments:", errEmpty)
	}

	fmt.Println("Exercise 3")

	minTemp, maxTemp, avgTemp, validCnt, invalidCnt, errSensor := AnalyzeSensorData(22.5, 23.1, -5.0, 101.0, 21.9, 0.0, 24.0)
	if errSensor != nil {
		fmt.Println("Analysis error:", errSensor)
	} else {
		fmt.Printf("Temp Min: %.2f, Max: %.2f, Average: %.2f, Valid: %d, Invalid: %d\n", minTemp, maxTemp, avgTemp, validCnt, invalidCnt)
	}

	_, _, _, _, _, errAllInvalid := AnalyzeSensorData(-10.0, 105.0, 0.0)
	if errAllInvalid != nil {
		fmt.Println("Error for all invalid data:", errAllInvalid)
	}
}
