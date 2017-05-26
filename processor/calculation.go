package processor

import (
	"fmt"

	"github.com/josler/pingu/core"
)

type Calculation interface {
	Calculate(event *core.Event)
}

type CountCalculation struct {
	value int64
}

func (c *CountCalculation) Calculate(event *core.Event) {
	c.value++
}

func (c *CountCalculation) String() string {
	return fmt.Sprintf("%v", c.value)
}

type SumCalculation struct {
	key     string
	samples int64
	value   int64
}

func NewSumCalculation(key string) *SumCalculation {
	return &SumCalculation{key: key}
}

func (c *SumCalculation) Calculate(event *core.Event) {
	val, found := event.Data[c.key]
	if found {
		switch val := val.(type) {
		case int:
			c.value += int64(val)
		case int64:
			c.value += val
		default:
			return
		}
		c.samples++
	}
}

func (c *SumCalculation) String() string {
	return fmt.Sprintf("%v", c.value)
}

type MeanCalculation struct {
	key     string
	samples float64
	value   float64
}

func NewMeanCalculation(key string) *MeanCalculation {
	return &MeanCalculation{key: key}
}

func (c *MeanCalculation) Calculate(event *core.Event) {
	val, found := event.Data[c.key]
	if found {
		switch val := val.(type) {
		case int:
			c.value += float64(val)
		case int64:
			c.value += float64(val)
		default:
			return
		}
		c.samples++
	}
}

func (c *MeanCalculation) String() string {
	return fmt.Sprintf("%v", (c.value / c.samples))
}
