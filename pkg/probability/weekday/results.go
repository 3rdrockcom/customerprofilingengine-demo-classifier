package weekday

import (
	"fmt"
	"sort"
	"time"
)

type Results []Result

func (r Results) Len() int           { return len(r) }
func (r Results) Less(i, j int) bool { return r[i].Probability < r[j].Probability }
func (r Results) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }

type Result struct {
	Weekday     time.Weekday
	Count       int
	Total       float64
	Probability float64
}

func (r Results) Display() {
	sort.Sort(sort.Reverse(r))

	fmt.Println("Probability: Weekday\n---")
	for i := range r {
		if r[i].Probability == 0 {
			break
		}
		fmt.Println(fmt.Sprintf("%-9v: %11.2f %%", r[i].Weekday, r[i].Probability*100))
	}
}
