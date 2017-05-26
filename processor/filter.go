package processor

import "github.com/josler/pingu/core"

type Filter struct {
	*pipelineStage
	rule Rule
}

func NewFilter(rule Rule, in <-chan *core.Event) *Filter {
	f := Filter{rule: rule, pipelineStage: &pipelineStage{in: in, out: make(chan *core.Event)}}
	go f.Start()
	return &f
}

func (f *Filter) Start() {
	for {
		event := <-f.in
		if f.rule.Match(event) {
			f.out <- event
		}
	}
}
