package day

import (
	"github.com/epointpayment/customerprofilingengine-demo-classifier/pkg/models"
)

var Debug bool

type Day struct {
	t models.Transactions
}

func NewDay(t models.Transactions) *Day {
	d := &Day{
		t: t,
	}

	return d
}

func (d *Day) Run() Results {
	return d.calc()
}

func (d *Day) calc() Results {
	list := make(Results, 31)

	count := 0
	total := 0.0
	for i := range d.t {
		day := d.t[i].Date.Day() - 1

		list[day].Day = d.t[i].Date.Day()
		list[day].Total += d.t[i].Credits
		list[day].Count++

		total += d.t[i].Credits
		count++
	}

	for i := range list {
		entry := list[i]
		entry.Probability = entry.Total / total

		list[i] = entry
	}

	return list
}
