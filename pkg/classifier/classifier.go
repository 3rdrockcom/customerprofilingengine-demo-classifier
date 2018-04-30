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
	c.Ranks = append(c.Ranks, rank)

	rank = c.doBiWeekly()
	c.Ranks = append(c.Ranks, rank)

	rank = c.doWeekly()
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

	list := make(map[int]int)

	//
	for d := dateRangeMin; d.Before(dateRangeMax); d = d.AddDate(0, 1, 0) {
		k, _ := strconv.Atoi(d.Format("20060102"))
		list[k] = 0

		for i := 0; i < len(t); i++ {
			if (t[i].Date.After(d) || t[i].Date.Equal(d)) && t[i].Date.Before(d.AddDate(0, 1, 0)) {
				list[k]++
			}
		}
	}
	rank := ranks.NewRank("Monthly", c.getScore(list), 10)

	if c.Debug {
		fmt.Println(fmt.Sprintf("Class: %s [%.2f]", rank.Name, rank.Value*100.0))
		fmt.Println(list)
	}

	return rank
}

func (c *Classifier) doBiWeekly() ranks.Rank {
	t := c.Transactions

	dateMin, dateMax := c.getDateRange()
	dateRangeMin := now.New(dateMin).BeginningOfWeek()
	dateRangeMax := now.New(dateMax).EndOfWeek()

	list := make(map[int]int)

	//
	for d := dateRangeMin; d.Before(dateRangeMax); d = d.AddDate(0, 0, 14) {
		k, _ := strconv.Atoi(d.Format("20060102"))
		list[k] = 0

		for i := 0; i < len(t); i++ {
			if (t[i].Date.After(d) || t[i].Date.Equal(d)) && t[i].Date.Before(d.AddDate(0, 0, 14)) {
				list[k]++
			}
		}
	}

	rank := ranks.NewRank("BiWeekly", c.getScore(list), 20)

	if c.Debug {
		fmt.Println(fmt.Sprintf("Class: %s [%.2f]", rank.Name, rank.Value*100.0))
		fmt.Println(list)
	}

	return rank
}

func (c *Classifier) doWeekly() ranks.Rank {
	t := c.Transactions

	dateMin, dateMax := c.getDateRange()
	dateRangeMin := now.New(dateMin).BeginningOfWeek()
	dateRangeMax := now.New(dateMax).EndOfWeek()

	list := make(map[int]int)

	//
	for d := dateRangeMin; d.Before(dateRangeMax); d = d.AddDate(0, 0, 7) {
		k, _ := strconv.Atoi(d.Format("20060102"))
		list[k] = 0

		for i := 0; i < len(c.Transactions); i++ {
			if (t[i].Date.After(d) || t[i].Date.Equal(d)) && t[i].Date.Before(d.AddDate(0, 0, 7)) {
				list[k]++
			}
		}
	}

	rank := ranks.NewRank("Weekly", c.getScore(list), 30)

	if c.Debug {
		fmt.Println(fmt.Sprintf("Class: %s [%.2f]", rank.Name, rank.Value*100.0))
		fmt.Println(list)
	}

	return rank
}

func (c *Classifier) getScore(list map[int]int) float64 {
	l := 0.0
	for i := range list {
		if list[i] > 0 {
			l++
		}
		if list[i] > 1 {
			l -= .5
		}
	}

	return float64(l) / float64(len(list))
}

func (c *Classifier) getDateRange() (time.Time, time.Time) {
	dateMin := c.Transactions[0].Date
	dateMax := c.Transactions[len(c.Transactions)-1].Date

	return dateMin, dateMax
}
