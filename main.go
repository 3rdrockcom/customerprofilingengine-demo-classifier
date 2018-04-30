package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/epointpayment/customerprofilingengine-demo-classifier/pkg/classifier"
	"github.com/epointpayment/customerprofilingengine-demo-classifier/pkg/csv"
	"github.com/epointpayment/customerprofilingengine-demo-classifier/pkg/probability"

	"github.com/fatih/color"
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

	// Probability
	p := probability.New(transactions)
	p.Debug = Debug

	p.RunDay().Display()
	fmt.Println()

	p.RunWeekday().Display()
	fmt.Println()

	// Classify account
	cl, err := classifier.NewClassifier(transactions)
	if err != nil {
		panic(err)
	}
	cl.Debug = Debug
	cl.Process()
	classification := cl.GetClassification()

	o := color.New(color.Bold)
	o.Println(strings.ToUpper(fmt.Sprintf("Classification: %s [%.6f]", classification.Name, classification.Value)))
}
