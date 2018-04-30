package probability

import (
	"sort"

	"github.com/epointpayment/customerprofilingengine-demo-classifier/pkg/models"
	"github.com/epointpayment/customerprofilingengine-demo-classifier/pkg/probability/day"
	"github.com/epointpayment/customerprofilingengine-demo-classifier/pkg/probability/weekday"
)

type Probability struct {
	Transactions models.Transactions
	Debug        bool
}

func New(t models.Transactions) *Probability {
	sort.Sort(t)

	return &Probability{
		Transactions: t,
	}
}

func (p *Probability) RunDay() day.Results {
	d := day.NewDay(p.Transactions)
	d.Debug = p.Debug

	return d.Run()
}

func (p *Probability) RunWeekday() weekday.Results {
	w := weekday.NewWeekday(p.Transactions)
	w.Debug = p.Debug

	return w.Run()
}
