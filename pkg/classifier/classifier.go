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

var Debug bool

type Classifier struct {
	Transactions models.Transactions
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

func (c *Classifier) Process() Results {
	var listRank ranks.Ranks
	var rank ranks.Rank
	var list = make(map[string]Credits)

	buckets := []string{"monthly", "biweekly", "weekly"}
	for i := range buckets {
		name := buckets[i]
		if Debug {
			fmt.Println(fmt.Sprintf("::: Class: %s :::\n", name))
		}

		switch name {
		case "monthly":
			list[name] = c.processMonthly()
			rank = ranks.NewRank(name, c.calcRankValue(list[name]), 10)
		case "biweekly":
			list[name] = c.processBiWeekly()
			rank = ranks.NewRank(name, c.calcRankValue(list[name]), 20)
		case "weekly":
			list[name] = c.processWeekly()
			rank = ranks.NewRank(name, c.calcRankValue(list[name]), 30)
		}

		if Debug {
			fmt.Println(fmt.Sprintf("Score: %.6f\n", rank.Value))
		}
		listRank = append(listRank, rank)
	}

	sort.Sort(sort.Reverse(listRank))

	res := Results{}
	for i := range listRank {
		entry := Result{
			Name:  listRank[i].Name,
			Score: listRank[i].Value,
			List:  list[listRank[i].Name],
		}
		res = append(res, entry)
	}

	return res
}

func (c *Classifier) processMonthly() Credits {
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

	return list
}

func (c *Classifier) processBiWeekly() Credits {
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

	return list
}

func (c *Classifier) processWeekly() Credits {
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

	return list
}

func (c *Classifier) calcRankValue(list Credits) float64 {
	data := []float64{}

	var keys []int
	for k := range list {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	if Debug {
		fmt.Println("Date Range(s):")
	}

	rankValue := 0.0
	total := 0.0
	for k := range keys {
		i := keys[k]

		sum := 0.0
		for j := range list[i] {
			total += list[i][j].Amount
			sum += list[i][j].Amount
		}

		data = append(data, sum)

		v := 0.0
		if len(list[i]) == 1 {
			v = 1
		}
		if len(list[i]) > 1 {
			v = -(float64(len(list[i])) - 1)
		}
		rankValue += v

		if Debug {
			fmt.Println(fmt.Sprintf("%v: %10v %5v [%v / %v]", i, sum, v, rankValue, k+1))
		}
	}

	rankValue = rankValue / float64(len(list))

	if Debug {
		mean, sd, _ := c.getStatistics(data)

		fmt.Println()
		fmt.Println(fmt.Sprintf("Statistics: %.2f ± %.2f", mean, sd))
	}

	return rankValue
}

func (c *Classifier) getDateRange() (time.Time, time.Time) {
	dateMin := c.Transactions[0].Date
	dateMax := c.Transactions[len(c.Transactions)-1].Date

	return dateMin, dateMax
}

func (c *Classifier) getStatistics(data []float64) (float64, float64, error) {
	var err error

	mean, err := stats.Mean(data)
	if err != nil {
		return 0, 0, err
	}

	sd, err := stats.StandardDeviation(data)
	if err != nil {
		return 0, 0, err
	}

	return mean, sd, nil
}
