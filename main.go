package main

import (
	"fmt"

	"github.com/josler/pingu/core"
	"github.com/josler/pingu/processor"
)

func main() {
	rule := processor.NewAndRule(
		&processor.NameMatchRule{MatchName: "flarv"},
		&processor.NotFieldMatchRule{FieldName: "bar", Value: "x"},
	)
	eventChan := make(chan *core.Event)

	filter := processor.NewFilter(rule, eventChan)
	grouping := processor.NewGrouping(
		"foo",
		func() processor.Calculation { return processor.NewMeanCalculation("baz") },
		filter.Out(),
	)

	eventChan <- &core.Event{Name: "flarv", Data: map[string]interface{}{"foo": "a", "bar": "x", "baz": 1}}
	eventChan <- &core.Event{Name: "nopey", Data: map[string]interface{}{"foo": "b", "bar": "x"}}
	eventChan <- &core.Event{Name: "flarv", Data: map[string]interface{}{"foo": "a", "baz": 4}}
	eventChan <- &core.Event{Name: "flarv", Data: map[string]interface{}{"foo": "a", "baz": 7}}
	eventChan <- &core.Event{Name: "flarv", Data: map[string]interface{}{"foo": "b", "baz": 2}}

	fmt.Println(grouping)
}
