package day

import (
	"fmt"
	"sort"
)

type Results []Result

func (r Results) Len() int           { return len(r) }
func (r Results) Less(i, j int) bool { return r[i].Probability < r[j].Probability }
func (r Results) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }

type Result struct {
	Day         int
	Count       int
	Total       float64
	Probability float64
}

func (r Results) Display() {
	sort.Sort(sort.Reverse(r))

	fmt.Println("Probability: Day\n---")
	for i := range r {
		if r[i].Probability == 0 {
			break
		}
		fmt.Println(fmt.Sprintf("%-7v: %13.2f %%", fmt.Sprintf("Day %2v", r[i].Day), r[i].Probability*100))
	}
}
