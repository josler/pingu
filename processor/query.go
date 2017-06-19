package processor

import (
	"fmt"

	"github.com/josler/pingu/core"
)

type Query struct {
	limit    int
	filter   *Filter
	grouping *Grouping
}

type EventGenerator interface {
	Out() <-chan *core.Event
}

func NewQuery(generator EventGenerator, limit int, config *core.Config) *Query {
	q := Query{limit: limit}
	rule := ParseRules(config.Rules())
	f := ParseCalculation(config.Calculation())
	q.filter = NewFilter(rule, generator.Out())
	q.grouping = NewGrouping(
		config.GroupBy(),
		f,
		q.filter.Out(),
	)

	return &q
}

func (q *Query) Start() {
	q.grouping.Start()
	fmt.Println(q.grouping.Report(q.limit))
	// fmt.Println(grouping.ReportBuckets(limitVar))
}
