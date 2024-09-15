package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"

	"github.com/xuri/excelize/v2"
)

// Worker function to process rows and convert them to JSON
func worker(rowsChan chan []string, jsonLinesChan chan map[string]string, wg *sync.WaitGroup) {
	defer wg.Done()
	for row := range rowsChan {
		jsonLinesChan <- processRowToJSON(row)
	}
}

// Function to process a row from Excel to a JSON-like map
func processRowToJSON(row []string) map[string]string {
	rowMap := make(map[string]string)
	for i, colCell := range row {
		rowMap[fmt.Sprintf("column%d", i)] = colCell
	}
	return rowMap
}

// Function to write JSON output to file using streaming
func writeJSONToFile(jsonLinesChan chan map[string]string, doneChan chan bool, outputFile *os.File) {
	encoder := json.NewEncoder(outputFile)
	first := true

	// Begin JSON array
	outputFile.WriteString("[")

	for jsonLine := range jsonLinesChan {
		if !first {
			outputFile.WriteString(",")
		}
		encoder.Encode(jsonLine)
		first = false
	}

	// End JSON array
	outputFile.WriteString("]")
	doneChan <- true
}

func main() {
	// Open the Excel file
	f, err := excelize.OpenFile("large_file.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Open JSON output file
	jsonFile, err := os.Create("output.json")
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	// Create channels for worker communication
	rowsChan := make(chan []string, 100)          // Buffer rows
	jsonLinesChan := make(chan map[string]string) // Buffer JSON lines

	// WaitGroup for goroutines
	var wg sync.WaitGroup

	// Start worker goroutines
	numWorkers := runtime.NumCPU() // Use number of CPU cores available
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(rowsChan, jsonLinesChan, &wg)
	}

	// Start JSON writer goroutine
	doneChan := make(chan bool)
	go writeJSONToFile(jsonLinesChan, doneChan, jsonFile)

	// Read Excel rows and feed them to workers
	rows, err := f.Rows("Sheet1")
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		row, err := rows.Columns()
		if err != nil {
			log.Println(err)
			continue
		}
		rowsChan <- row
	}

	// Close rows channel after all rows are read
	close(rowsChan)

	// Wait for all workers to finish
	wg.Wait()

	// Close the JSON lines channel after all workers are done
	close(jsonLinesChan)

	// Wait for JSON writer to finish
	<-doneChan

	if err := rows.Close(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Conversion complete: Excel to JSON")
}
