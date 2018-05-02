package models

import (
	"time"
)

type Transactions []Transaction

func (t Transactions) Len() int           { return len(t) }
func (t Transactions) Less(i, j int) bool { return t[i].Date.Before(t[j].Date) }
func (t Transactions) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }

func (t Transactions) Separator(p float64) []Transactions {
	MaxTransaction := 0.0

	for i := 0; i < len(t); i++ {
		if t[i].Credits > MaxTransaction {
			MaxTransaction = t[i].Credits
		}
	}

	threshold := MaxTransaction * p

	res := make([]Transactions, 2)
	for i := 0; i < len(t); i++ {
		if t[i].Credits >= threshold {
			k := 0
			res[k] = append(res[k], t[i])
		} else {
			k := 1
			res[k] = append(res[k], t[i])
		}
	}

	return res
}

type Transaction struct {
	Date    time.Time
	Credits float64
}
