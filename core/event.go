package core

import "fmt"

type Event struct {
	Name string
	Data map[string]interface{}
}

func (e *Event) String() string {
	return fmt.Sprintf("Name: %s, Data: %v", e.Name, e.Data)
}
