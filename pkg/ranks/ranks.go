package ranks

type Ranks []Rank

func (r Ranks) Len() int { return len(r) }
func (r Ranks) Less(i, j int) bool {
	if r[i].Value == r[j].Value {
		return r[i].Weight < r[j].Weight
	} else {
		return r[i].Value < r[j].Value
	}
}
func (r Ranks) Swap(i, j int) { r[i], r[j] = r[j], r[i] }

type Rank struct {
	Name   string
	Value  float64
	Weight int
}

func NewRank(name string, value float64, weight int) Rank {
	return Rank{
		Name:   name,
		Value:  value,
		Weight: weight,
	}
}
