package tailer

import (
	"bufio"
	"fmt"
	"io"

	"encoding/json"

	"github.com/josler/pingu/core"
)

type Parser struct {
	out chan *core.Event
}

func (p *Parser) Parse(reader io.Reader) {
	if p.out == nil {
		p.out = make(chan *core.Event)
	}
	go p.parse(reader)
}

func (p *Parser) parse(reader io.Reader) {
	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	for scanner.Scan() {
		e := &core.Event{Name: "requests"}
		err := json.Unmarshal(scanner.Bytes(), &e.Data)
		if err != nil {
			fmt.Println(err)
		}
		p.out <- e
	}
	if scanner.Err() != nil {
		panic(scanner.Err())
	}
}

func (p *Parser) Out() <-chan *core.Event {
	if p.out == nil {
		p.out = make(chan *core.Event)
	}
	return p.out
}
