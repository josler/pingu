package processor

import (
	"fmt"
	"sort"
	"time"

	"github.com/josler/pingu/core"
)

const WAIT_DURATION time.Duration = 1 * time.Second

type Grouping struct {
	*pipelineStage
	key         string
	builder     func() Calculation
	calcGroups  map[interface{}]Calculation
	eventGroups map[interface{}][]*core.Event
}

func NewGrouping(key string, builder func() Calculation, in <-chan *core.Event) *Grouping {
	g := Grouping{
		pipelineStage: &pipelineStage{in: in, out: make(chan *core.Event)},
		key:           key,
		builder:       builder,
		calcGroups:    map[interface{}]Calculation{},
		eventGroups:   map[interface{}][]*core.Event{},
	}
	return &g
}

func (g *Grouping) Start() {
	timer := time.NewTimer(WAIT_DURATION)

	for {
		select {
		case <-timer.C:
			return
		case event := <-g.in:
			timer.Reset(WAIT_DURATION)
			g.process(event)
		}
	}
}
func (g *Grouping) process(event *core.Event) {
	groupVal, found := event.Data[g.key]
	if !found {
		return
	}
	calc, found := g.calcGroups[groupVal]
	if !found {
		calc = g.builder()
		g.calcGroups[groupVal] = calc
		g.eventGroups[groupVal] = []*core.Event{}
	}
	calc.Calculate(event)
	g.eventGroups[groupVal] = append(g.eventGroups[groupVal], event)
}

type Pair struct {
	Key   interface{}
	Value Calculation
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value.LessThan(p[j].Value) }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func (g *Grouping) sorted() PairList {
	pl := make(PairList, len(g.calcGroups))
	i := 0
	for k, v := range g.calcGroups {
		pl[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}

func (g *Grouping) topN(n int) PairList {
	sorted := g.sorted()
	if len(sorted) <= n {
		return sorted
	}
	return g.sorted()[:n]
}

func (g *Grouping) Report(n int) string {
	str := fmt.Sprintf("Grouping by %v:\n\n", g.key)

	if n < 0 {
		n = len(g.calcGroups)
	}
	for _, pair := range g.topN(n) {
		str += fmt.Sprintf("%v: %v\n", pair.Key, pair.Value)
	}

	return str
}

func (g *Grouping) String() string {
	str := fmt.Sprintf("Grouping by %v:\n\n", g.key)

	for _, pair := range g.sorted() {
		str += fmt.Sprintf("%v: %v\n", pair.Key, pair.Value)
	}

	// str := fmt.Sprintf("Grouping by %v:\n\n", g.key)
	// for k, v := range g.eventGroups {
	// 	str += fmt.Sprintf("%v: %v\n", k, v)
	// }
	return str
}
