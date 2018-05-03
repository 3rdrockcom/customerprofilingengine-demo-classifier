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

	if Debug {
		classifier.Debug = Debug
		probability.Debug = Debug
	}

	// Parse CSV file and extract transactions
	csvFile := csv.NewCSV(Filename)
	t, err := csvFile.Parse()
	if err != nil {
		panic(err)
	}

	tSplit := t.Separator(.5)

	for i := 0; i < len(tSplit); i++ {
		transactions := tSplit[i]

		if len(transactions) == 0 {
			break
		}

		o := color.New(color.Bold).Add(color.FgGreen)
		switch i {
		case 0:
			o.Println(strings.ToUpper("--- Results [Primary] ---"))
			fmt.Println()
		case 1:
			fmt.Println()
			o.Println(strings.ToUpper("--- Results [Secondary] ---"))
			fmt.Println()
		}

		// Probability
		p := probability.New(transactions)
		probDay := p.RunDay()
		probWeekday := p.RunWeekday()

		// Classify account
		cl, err := classifier.NewClassifier(transactions)
		if err != nil {
			panic(err)
		}
		res := cl.Process()

		o = color.New(color.Bold).Add(color.BgBlue).Add(color.FgWhite)

		for i := range res {
			cp := res.GetProbability(i)

			o.Println(strings.ToUpper(fmt.Sprintf("Classification: %s [%.2f %%]", res[i].Name, cp*100)))

			avgPerInterval := res.GetAveragePerInterval(i)
			avg := res.GetAverage(i)

			fmt.Println(strings.ToUpper(fmt.Sprintf("Average Credits Per %s Interval: %.2f", res[i].Name, avgPerInterval)))
			fmt.Println(strings.ToUpper(fmt.Sprintf("Average Credits: %.2f", avg)))
			fmt.Println()

			probDay.Display()

			if res[i].Name == "weekly" {
				fmt.Println()
				probWeekday.Display()
			}

			fmt.Println()
			fmt.Println()
		}
	}
}
