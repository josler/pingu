package processor

import (
	"fmt"
	"time"

	"github.com/josler/pingu/core"
)

// A BucketCalculation represents a set of calculations for events during a period of time
type BucketCalculation struct {
	startTime  time.Time
	builder    func() Calculation
	calcGroups calculationGroups
	rawEvents  []*core.Event
}

func NewBucketCalculation(startTime time.Time, builder func() Calculation) *BucketCalculation {
	return &BucketCalculation{startTime: startTime, builder: builder, calcGroups: calculationGroups{}, rawEvents: []*core.Event{}}
}

// Merge other into this bucket
func (bucket *BucketCalculation) Merge(other *BucketCalculation) {
	if other.startTime.Before(bucket.startTime) {
		bucket.startTime = other.startTime
	}

	bucket.rawEvents = append(bucket.rawEvents, other.rawEvents...)

	for key, calc := range other.calcGroups {
		existing, found := bucket.calcGroups[key]
		if !found { // if not found just copy
			bucket.calcGroups[key] = calc.Copy()
			continue
		}
		existing.Merge(calc)
	}
}

func (bucket *BucketCalculation) add(key interface{}, event *core.Event) {
	calc, found := bucket.calcGroups[key]
	if !found {
		calc = bucket.builder()
		bucket.calcGroups[key] = calc
	}
	calc.Calculate(event)
	bucket.rawEvents = append(bucket.rawEvents, event)
}

func (bucket *BucketCalculation) Report(n int) string {
	if n < 0 {
		n = len(bucket.calcGroups)
	}

	str := "---------\n"
	str += fmt.Sprintf("%v, %v groups\n", bucket.startTime, len(bucket.calcGroups))
	str += fmt.Sprintf("\n%v\n", bucket.calcGroups.Report(n))

	return str
}
