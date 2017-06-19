package tailer

import (
	"bufio"
	"fmt"
	"io"

	"encoding/json"

	"github.com/josler/pingu/core"
)

type Parser struct {
	out   chan *core.Event
	endAt int64
}

func (p *Parser) Parse(reader io.Reader) {
	if p.out == nil {
		p.out = make(chan *core.Event)
	}
	go p.parse(reader)
}

func (p *Parser) ParseTo(reader io.Reader, endAt int64) {
	p.endAt = endAt
	p.Parse(reader)
}

func (p *Parser) parse(reader io.Reader) {
	buf := bufio.NewReader(reader)
	var data []byte
	var err error
	for {
		data, err = buf.ReadBytes('\n')
		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			break
		}
		e := &core.Event{Name: "requests"}
		err := json.Unmarshal(data, &e.Data)
		if err != nil {
			// fmt.Println(err)
			continue
		}
		p.out <- e
	}
}

func (p *Parser) Out() <-chan *core.Event {
	if p.out == nil {
		p.out = make(chan *core.Event)
	}
	return p.out
}
