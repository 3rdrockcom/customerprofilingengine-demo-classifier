package models

import "time"

type Transactions []Transaction

func (t Transactions) Len() int           { return len(t) }
func (t Transactions) Less(i, j int) bool { return t[i].Date.Before(t[j].Date) }
func (t Transactions) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }

type Transaction struct {
	Date    time.Time
	Credits float64
}
