package probability

import (
	"sort"

	"github.com/epointpayment/customerprofilingengine-demo-classifier/pkg/models"
	"github.com/epointpayment/customerprofilingengine-demo-classifier/pkg/probability/day"
	"github.com/epointpayment/customerprofilingengine-demo-classifier/pkg/probability/weekday"
)

var Debug bool

type Probability struct {
	Transactions models.Transactions
}

func New(t models.Transactions) *Probability {
	day.Debug = Debug
	weekday.Debug = Debug

	sort.Sort(t)

	return &Probability{
		Transactions: t,
	}
}

func (p *Probability) RunDay() day.Results {
	d := day.NewDay(p.Transactions)
	return d.Run()
}

func (p *Probability) RunWeekday() weekday.Results {
	w := weekday.NewWeekday(p.Transactions)
	return w.Run()
}
