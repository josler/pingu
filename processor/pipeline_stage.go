package processor

import "github.com/josler/pingu/core"

type pipelineStage struct {
	in  <-chan *core.Event
	out chan *core.Event
}

func (s *pipelineStage) Out() <-chan *core.Event {
	return s.out
}
