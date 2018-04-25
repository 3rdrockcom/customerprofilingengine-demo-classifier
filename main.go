package main

import (
	"flag"
	"fmt"
)

var Filename string
var Debug bool

func init() {
	flag.StringVar(&Filename, "file", "sample.csv", "path to sample data")
	flag.BoolVar(&Debug, "debug", false, "show debugging information")

	flag.Parse()
}

func main() {
	// Parse CSV file and extract transactions
	csvFile := NewCSV(Filename)
	transactions, err := csvFile.Parse()
	if err != nil {
		panic(err)
	}

	// Classify account
	classifier, err := NewClassifier(transactions)
	if err != nil {
		panic(err)
	}
	classifier.Process()
	classification := classifier.GetClassification()

	fmt.Println("Classification: " + fmt.Sprintf("%s [%.2f]", classification.Name, classification.Value*100.0))
}
