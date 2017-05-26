package processor

import (
	"fmt"

	"github.com/josler/pingu/core"
)

type Calculation interface {
	Calculate(event *core.Event)
	LessThan(other Calculation) bool
}

// ParseCalculation returns a builder function that will build Calculations
func ParseCalculation(calc *core.JsonCalculation) func() Calculation {
	switch calc.Type {
	case "count":
		return func() Calculation { return &CountCalculation{} }
	case "sum":
		return func() Calculation { return NewSumCalculation(calc.Name) }
	case "mean":
		return func() Calculation { return NewMeanCalculation(calc.Name) }
	}
	return nil
}

type CountCalculation struct {
	value int64
}

func (c *CountCalculation) Calculate(event *core.Event) {
	c.value++
}

func (c *CountCalculation) LessThan(other Calculation) bool {
	return c.value < other.(*CountCalculation).value
}

func (c *CountCalculation) String() string {
	return fmt.Sprintf("%v", c.value)
}

type SumCalculation struct {
	key   string
	value float64
}

func NewSumCalculation(key string) *SumCalculation {
	return &SumCalculation{key: key}
}

func (c *SumCalculation) Calculate(event *core.Event) {
	val, found := event.Data[c.key]
	if found {
		switch val := val.(type) {
		case int:
			c.value += float64(val)
		case int64:
			c.value += float64(val)
		case float32:
			c.value += float64(val)
		case float64:
			c.value += val
		default:
			return
		}
	}
}

func (c *SumCalculation) LessThan(other Calculation) bool {
	return c.value < other.(*SumCalculation).value
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
		case float32:
			c.value += float64(val)
		case float64:
			c.value += val
		default:
			return
		}
		c.samples++
	}
}

func (c *MeanCalculation) LessThan(other Calculation) bool {
	return (c.value / c.samples) < (other.(*MeanCalculation).value / other.(*MeanCalculation).samples)
}

func (c *MeanCalculation) String() string {
	return fmt.Sprintf("%v", (c.value / c.samples))
}
