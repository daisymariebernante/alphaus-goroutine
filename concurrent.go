package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"sync"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run concurrent.go <file_path>")
		return
	}
	inputFilePath := os.Args[1]

	file, err := os.Open(inputFilePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	const maxWorkers = 100
	var wg sync.WaitGroup
	var sumTotal float64
	var countTotal int

	// Buffered channel to limit goroutines
	sem := make(chan struct{}, maxWorkers)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Acquire a slot from the buffered channel
		sem <- struct{}{}

		wg.Add(1)
		go func(line string) {
			defer func() {
				// Release the slot back to the buffered channel
				<-sem
				wg.Done()
			}()

			// Process the line
			lineAverage := calculateLineAverage(line)

			// Update the sumTotal and countTotal safely
			updateTotal(&sumTotal, &countTotal, lineAverage)
		}(line)
	}

	wg.Wait()

	// Calculate and print the overall average
	overallAverage := sumTotal / float64(countTotal)
	fmt.Printf("Overall average: %.2f\n", overallAverage)
}

func calculateLineAverage(line string) float64 {
	sum := 0
	count := 0

	// Scan for all single digits from the input line
	for _, char := range line {
		if char >= '0' && char <= '9' {
			num, _ := strconv.Atoi(string(char))
			sum += num
			count++
		}
	}

	// Calculate the average of all single digit occurrences
	var avg float64
	if count > 0 {
		avg = float64(sum) / float64(count)
	}

	return avg
}

func updateTotal(sumTotal *float64, countTotal *int, lineAverage float64) {
	mutex := sync.Mutex{}
	mutex.Lock()
	defer mutex.Unlock()

	*sumTotal += lineAverage
	*countTotal++
}
