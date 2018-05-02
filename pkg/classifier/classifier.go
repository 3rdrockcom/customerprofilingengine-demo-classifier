package classifier

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/epointpayment/customerprofilingengine-demo-classifier/pkg/models"
	"github.com/epointpayment/customerprofilingengine-demo-classifier/pkg/ranks"

	"github.com/jinzhu/now"
	"github.com/montanaflynn/stats"
)

type Classifier struct {
	Transactions models.Transactions
	Ranks        ranks.Ranks
	Debug        bool
}

func NewClassifier(t models.Transactions) (*Classifier, error) {
	if len(t) == 0 {
		return nil, errors.New("transactions required")
	}

	sort.Sort(t)

	c := &Classifier{
		Transactions: t,
	}
	return c, nil
}

func (c *Classifier) Process() {
	var rank ranks.Rank

	rank = c.doMonthly()
	rank.Weight = 10
	c.Ranks = append(c.Ranks, rank)

	rank = c.doBiWeekly()
	rank.Weight = 20
	c.Ranks = append(c.Ranks, rank)

	rank = c.doWeekly()
	rank.Weight = 30
	c.Ranks = append(c.Ranks, rank)

	sort.Sort(sort.Reverse(c.Ranks))
}

func (c *Classifier) GetClassification() ranks.Rank {
	classification := c.Ranks[0]
	return classification
}

func (c *Classifier) doMonthly() ranks.Rank {
	t := c.Transactions

	dateMin, dateMax := c.getDateRange()
	dateRangeMin := now.New(dateMin).BeginningOfMonth()
	dateRangeMax := now.New(dateMax).EndOfMonth()

	list := make(Credits)

	//
	for d := dateRangeMin; d.Before(dateRangeMax); d = d.AddDate(0, 1, 0) {
		k, _ := strconv.Atoi(d.Format("20060102"))
		list[k] = []Credit{}

		for i := 0; i < len(t); i++ {
			if (t[i].Date.After(d) || t[i].Date.Equal(d)) && t[i].Date.Before(d.AddDate(0, 1, 0)) {
				list[k] = append(list[k], Credit{
					Amount: c.Transactions[i].Credits,
					Date:   c.Transactions[i].Date,
				})
			}
		}
	}
	rank := ranks.NewRank("Monthly", c.calcRankValue(list), 10)

	if c.Debug {
		fmt.Println(fmt.Sprintf("Class: %s [%.6f]\n", rank.Name, rank.Value))
	}

	return rank
}

func (c *Classifier) doBiWeekly() ranks.Rank {
	t := c.Transactions

	dateMin, dateMax := c.getDateRange()
	dateRangeMin := now.New(dateMin).BeginningOfWeek()
	dateRangeMax := now.New(dateMax).EndOfWeek()

	list := make(Credits)

	//
	for d := dateRangeMin; d.Before(dateRangeMax); d = d.AddDate(0, 0, 14) {
		k, _ := strconv.Atoi(d.Format("20060102"))
		list[k] = []Credit{}

		for i := 0; i < len(t); i++ {
			if (t[i].Date.After(d) || t[i].Date.Equal(d)) && t[i].Date.Before(d.AddDate(0, 0, 14)) {
				list[k] = append(list[k], Credit{
					Amount: c.Transactions[i].Credits,
					Date:   c.Transactions[i].Date,
				})
			}
		}
	}

	rank := ranks.NewRank("BiWeekly", c.calcRankValue(list), 20)

	if c.Debug {
		fmt.Println(fmt.Sprintf("Class: %s [%.6f]\n", rank.Name, rank.Value))
	}

	return rank
}

func (c *Classifier) doWeekly() ranks.Rank {
	t := c.Transactions

	dateMin, dateMax := c.getDateRange()
	dateRangeMin := now.New(dateMin).BeginningOfWeek()
	dateRangeMax := now.New(dateMax).EndOfWeek()

	list := make(Credits)

	//
	for d := dateRangeMin; d.Before(dateRangeMax); d = d.AddDate(0, 0, 7) {
		k, _ := strconv.Atoi(d.Format("20060102"))
		list[k] = []Credit{}

		for i := 0; i < len(c.Transactions); i++ {
			if (t[i].Date.After(d) || t[i].Date.Equal(d)) && t[i].Date.Before(d.AddDate(0, 0, 7)) {
				list[k] = append(list[k], Credit{
					Amount: c.Transactions[i].Credits,
					Date:   c.Transactions[i].Date,
				})
			}
		}
	}

	rank := ranks.NewRank("Weekly", c.calcRankValue(list), 30)

	if c.Debug {
		fmt.Println(fmt.Sprintf("Class: %s [%.6f]\n", rank.Name, rank.Value))
	}

	return rank
}

func (c *Classifier) calcRankValue(list Credits) float64 {
	data := []float64{}

	total := 0.0
	for i := range list {
		sum := 0.0

		for j := range list[i] {
			total += list[i][j].Amount
			sum += list[i][j].Amount
		}

		data = append(data, sum)
	}

	mean, _ := stats.Mean(data)
	sd, _ := stats.StandardDeviation(data)

	if c.Debug {
		fmt.Println(fmt.Sprintf("Statistics: %.2f ± %.2f", mean, sd))
	}

	rankValue := mean / sd
	return rankValue
}

func (c *Classifier) getDateRange() (time.Time, time.Time) {
	dateMin := c.Transactions[0].Date
	dateMax := c.Transactions[len(c.Transactions)-1].Date

	return dateMin, dateMax
}
