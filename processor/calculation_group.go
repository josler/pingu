package processor

import (
	"fmt"
	"sort"
)

type Pair struct {
	Key   interface{}
	Value Calculation
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value.LessThan(p[j].Value) }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type calculationGroups map[interface{}]Calculation

func (cg calculationGroups) sorted() PairList {
	pl := make(PairList, len(cg))
	i := 0
	for k, v := range cg {
		pl[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}

func (cg calculationGroups) topN(n int) PairList {
	sorted := cg.sorted()
	if len(sorted) <= n {
		return sorted
	}
	return cg.sorted()[:n]
}

func (cg calculationGroups) Report(n int) string {
	if n < 0 {
		n = len(cg)
	}
	str := ""
	for _, pair := range cg.topN(n) {
		str += fmt.Sprintf("%v: %v\n", pair.Key, pair.Value)
	}
	return str
}
