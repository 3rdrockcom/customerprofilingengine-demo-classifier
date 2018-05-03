package classifier

import (
	"sort"

	"github.com/montanaflynn/stats"
)

type Results []Result

type Result struct {
	Name  string
	Score float64
	List  Credits
}

func (r Results) GetClassification() Result {
	classification := r[0]
	return classification
}

func (r Results) GetAverage() float64 {
	data := []float64{}

	classification := r[0]
	list := classification.List

	var keys []int
	for k := range list {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for k := range keys {
		i := keys[k]
		sum := 0.0

		for j := range list[i] {
			sum += list[i][j].Amount
		}

		data = append(data, sum)
	}

	mean, _ := stats.Mean(data)
	return mean
}
