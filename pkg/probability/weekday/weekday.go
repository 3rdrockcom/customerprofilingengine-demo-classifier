package weekday

import (
	"github.com/epointpayment/customerprofilingengine-demo-classifier/pkg/models"
)

type Weekday struct {
	t     models.Transactions
	Debug bool
}

func NewWeekday(t models.Transactions) *Weekday {
	w := &Weekday{
		t: t,
	}

	return w
}

func (w *Weekday) Run() Results {
	return w.calc()
}

func (w *Weekday) calc() Results {
	list := make(Results, 31)

	count := 0
	total := 0.0
	for i := range w.t {
		weekday := int(w.t[i].Date.Weekday())

		list[weekday].Weekday = w.t[i].Date.Weekday()
		list[weekday].Total += w.t[i].Credits
		list[weekday].Count++

		total += w.t[i].Credits
		count++
	}

	for i := range list {
		entry := list[i]
		entry.Probability = entry.Total / total

		list[i] = entry
	}

	return list
}
