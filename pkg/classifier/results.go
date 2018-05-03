package classifier

import (
	"sort"

	"github.com/montanaflynn/stats"
)

type Results []Result

type Result struct {
	Name        string
	Score       float64
	Probability float64
	List        Credits
}

func (r Results) GetClassification(e int) Result {
	classification := r[e]
	return classification
}

func (r Results) GetProbability(e int) float64 {
	data := make([]float64, len(r))

	minScore := 0.0
	for i := range r {
		score := r[i].Score

		if score < minScore {
			minScore = score
		}

		data[i] = score
	}

	sum := 0.0
	for i := range data {
		data[i] = data[i] - minScore
		sum += data[i]
	}

	for i := range data {
		data[i] = (data[i] / sum)
	}

	return data[e]
}

func (r Results) GetAveragePerInterval(e int) float64 {
	data := []float64{}

	classification := r[e]
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

func (r Results) GetAverage(e int) float64 {
	data := []float64{}

	classification := r[e]
	list := classification.List

	for i := range list {
		for j := range list[i] {
			data = append(data, list[i][j].Amount)
		}
	}

	mean, _ := stats.Mean(data)
	return mean
}
