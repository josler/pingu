package processor

import "github.com/josler/pingu/core"
import "fmt"

type Grouping struct {
	*pipelineStage
	key     string
	builder func() Calculation
	groups  map[interface{}]Calculation
}

func NewGrouping(key string, builder func() Calculation, in <-chan *core.Event) *Grouping {
	g := Grouping{
		pipelineStage: &pipelineStage{in: in, out: make(chan *core.Event)},
		key:           key,
		builder:       builder,
		groups:        map[interface{}]Calculation{},
	}
	go g.Start()
	return &g
}

func (g *Grouping) Start() {
	for {
		event := <-g.in
		groupVal, found := event.Data[g.key]
		if !found {
			return
		}
		calc, found := g.groups[groupVal]
		if !found {
			calc = g.builder()
			g.groups[groupVal] = calc
		}
		calc.Calculate(event)
	}
}

func (g *Grouping) Flush() {
	// TODO(JO): how to flush to storage every so often
}

func (g *Grouping) String() string {
	str := "Grouping\n"
	for k, v := range g.groups {
		str += fmt.Sprintf("%v: %v\n", k, v)
	}
	return str
}
