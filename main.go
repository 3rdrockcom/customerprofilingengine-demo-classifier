package main

import (
	"flag"
	"fmt"

	"github.com/epointpayment/customerprofilingengine-demo-classifier/pkg/classifier"
	"github.com/epointpayment/customerprofilingengine-demo-classifier/pkg/csv"
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
	csvFile := csv.NewCSV(Filename)
	transactions, err := csvFile.Parse()
	if err != nil {
		panic(err)
	}

	// Classify account
	cl, err := classifier.NewClassifier(transactions)
	if err != nil {
		panic(err)
	}
	cl.Debug = Debug
	cl.Process()
	classification := cl.GetClassification()

	fmt.Println("Classification: " + fmt.Sprintf("%s [%.2f]", classification.Name, classification.Value*100.0))
}
