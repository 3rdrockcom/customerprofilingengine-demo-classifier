package classifier

import "time"

type Credits map[int][]Credit

type Credit struct {
	Date   time.Time
	Amount float64
}
