package main

import (
	"flag"
	"encoding/csv"
	"fmt"
	"encoding/json"
	"os"
)

func main() {
	// Define flags
	filenamePtr := flag.String("filename", ":", "name of CSV file to convert to JSON")
	delimiterPtr := flag.String("delimiter", ",", "delimiter character in the CSV")
	prefixPtr := flag.String("prefix", "", "prefix on each line in result JSON")
	indentPtr := flag.String("indent", "    ", "indentation in result JSON")
	withHeaderPtr := flag.Bool("with-header", false, "does the CSV file has header? (will produce struct for each record if true, list if false")
	_ = flag.String("help", "", "<insert help message here>")

	// Parse flags
	flag.Parse()

	// Validate flags
	if *filenamePtr == ":" {
		panic(fmt.Errorf("filename is required"))
	}

	if len(*delimiterPtr) != 1 {
		panic(fmt.Errorf("invalid delimiter: %s", *delimiterPtr))
	}
	delimiter := ([]rune(*delimiterPtr))[0]

	// Process CSV
	headers, lines, err := readCsv(*filenamePtr, delimiter, *withHeaderPtr)
	if err != nil {
		panic(err)
	}

	var result interface{}
	if *withHeaderPtr {
		result = convertWithHeader(headers, lines)
	} else {
		result = lines
	}

	jsonBytes, err := json.MarshalIndent(result, *prefixPtr, *indentPtr)
	if err != nil {
		panic(fmt.Errorf("Error converting to JSON: %v", err))
	}

	fmt.Println(string(jsonBytes))
}

func readCsv(filename string, delimiter rune, hasHeader bool) ([]string, [][]string, error) {
	csvFile, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}

	reader := csv.NewReader(csvFile)
	reader.Comma = delimiter

	rows, err := reader.ReadAll() 
	if err != nil {
		return nil, rows, err
	}

	if hasHeader {
		return rows[0], rows[1:], nil
	} else {
		return nil, rows, nil
	}
}

func convertWithHeader(headers []string, lines [][]string) []map[string]string {
	var records []map[string]string
	for _, line := range lines {
		record := map[string]string{}
		for idx := 0; idx < len(headers); idx++ {
			header := headers[idx]
			if idx < len(line) {
				record[header] = line[idx]
			} else {
				record[header] = ""
			}
		}
		records = append(records, record)
	}
	return records
}