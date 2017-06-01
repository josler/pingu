package processor

import (
	"fmt"
	"time"

	"github.com/josler/pingu/core"
)

const WAIT_DURATION time.Duration = 1 * time.Second

type Grouping struct {
	*pipelineStage
	key             string
	timestampKey    string
	windowingPeriod time.Duration
	builder         func() Calculation
	buckets         map[time.Time]*BucketCalculation
}

func NewGrouping(key string, builder func() Calculation, in <-chan *core.Event) *Grouping {
	g := Grouping{
		pipelineStage:   &pipelineStage{in: in, out: make(chan *core.Event)},
		key:             key,
		timestampKey:    "timestamp",
		windowingPeriod: 15 * time.Second,
		builder:         builder,
		buckets:         map[time.Time]*BucketCalculation{},
	}
	return &g
}

// SetTimestampKey allows the overriding of the default "timestamp"
func (g *Grouping) SetTimestampKey(timestampKey string) {
	g.timestampKey = timestampKey
}

// SetWindowingPeriod allows the overriding of the default 15 seconds.
func (g *Grouping) SetWindowingPeriod(windowingPeriod time.Duration) {
	g.windowingPeriod = windowingPeriod
}

func (g *Grouping) Start() {
	timer := time.NewTimer(WAIT_DURATION)

	for {
		select {
		case <-timer.C:
			return
		case event := <-g.in:
			timer.Reset(WAIT_DURATION)
			g.process(event)
		}
	}
}

func (g *Grouping) process(event *core.Event) {
	groupKeyVal, found := event.Data[g.key]
	if !found {
		return
	}

	// parse time
	timeVal, found := event.Data[g.timestampKey]
	if !found {
		return
	}
	t, err := time.Parse("2006-01-02T15:04:05.999Z", timeVal.(string))
	if err != nil {
		fmt.Println(err)
		fmt.Println(timeVal)
		return
	}

	// find bucket
	truncated := t.Truncate(g.windowingPeriod)
	bucket, found := g.buckets[truncated]
	if !found {
		bucket = NewBucketCalculation(truncated, g.builder)
		g.buckets[truncated] = bucket
	}

	bucket.add(groupKeyVal, event)
}

func (g *Grouping) Report(n int) string {
	str := fmt.Sprintf("Grouping by %v\n", g.key)
	str += fmt.Sprintf("Gathered over %d buckets, with windowing period: %v\n", len(g.buckets), g.windowingPeriod)

	// Merge buckets and report
	overallBucket := NewBucketCalculation(time.Now(), g.builder)

	for _, bucket := range g.buckets {
		overallBucket.Merge(bucket)
	}

	str += overallBucket.Report(n)
	return str
}

func (g *Grouping) ReportBuckets(n int) string {
	str := fmt.Sprintf("Grouping by %v\n", g.key)
	str += fmt.Sprintf("%d buckets, with windowing period: %v\n", len(g.buckets), g.windowingPeriod)

	for _, bucket := range g.buckets {
		str += bucket.Report(n)
	}

	return str
}
