package main

import (
	"encoding/csv"
	"os"
	"strconv"
	"time"
)

type CSV struct {
	Filename string
}

func NewCSV(filename string) CSV {
	return CSV{
		Filename: filename,
	}
}

func (c CSV) Parse() (Transactions, error) {
	var transactions Transactions
	var err error

	// Open CSV file
	f, err := os.Open(c.Filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Read File into a Variable
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, err
	}

	// Loop through lines & turn into object
	for i, line := range lines {
		if i == 0 {
			continue // Skip headers
		}

		// Date
		date, err := time.ParseInLocation(
			"01/02/2006",
			line[0], time.UTC)
		if err != nil {
			return nil, err
		}

		// Credit
		credit, err := strconv.ParseFloat(line[1], 64)
		if err != nil {
			return nil, err
		}

		// Merge transaction
		transactions = append(transactions, Transaction{
			Date:    date,
			Credits: credit,
		})
	}

	return transactions, err
}
