package metrics

import (
	"bytes"
	"sort"
	"strconv"
	"sync"

	"github.com/codahale/metrics"
)

// Snapshot lock.
var ss sync.Mutex

// counterSeries is a specialized counter.
// It is synonymous to influxdb series where series is composed of measurement name and tag key values.
type counterSeries struct {
	metrics.Counter
}

// newSeries provides counter name encoded in format of influxdb series: name,tags and field
func newCounterSeries(name string, tags map[string]string, field string) counterSeries {
	s := MakeSeries(name, tags, field)
	return counterSeries{
		metrics.Counter(s),
	}
}

// Counter inherits metrics.Counter and provides counter in influxdb format.
// Use a counter to derive rates (e.g., record total number of requests, derive
// requests per second).
type Counter struct {
	Name          string
	Tags          map[string]string
	Field         string
	counterSeries // counter series.
}

// NewCounter returns new instance of Counter.
func NewCounter(name string, tags map[string]string, field string) *Counter {
	return &Counter{
		Name:          name,
		Tags:          tags,
		Field:         field,
		counterSeries: counterSeries(newCounterSeries(name, tags, field)),
	}
}

// gaugeSeries is a specialized counter.
// It is synonymous to influxdb series where series is composed of measurement name and tag key values.
type gaugeSeries struct {
	metrics.Gauge
}

// newGaugeSeries is new instance of GaugeSeries
func newGaugeSeries(name string, tags map[string]string, field string) gaugeSeries {
	s := MakeSeries(name, tags, field)
	return gaugeSeries{
		metrics.Gauge(s),
	}
}

// Gauge inherits metrics.Gauge and provides gauges in influxdb format.
// A Gauge is an instantaneous measurement of a value.
//
// Use a gauge to track metrics which increase and decrease (e.g., amount of
// free memory).
type Gauge struct {
	Name        string
	Tags        map[string]string
	Field       string
	gaugeSeries // gauge series.
}

// NewCounter returns new instance of Counter
func NewGauge(name string, tags map[string]string, field string) *Gauge {
	return &Gauge{
		Name:        name,
		Tags:        tags,
		Field:       field,
		gaugeSeries: newGaugeSeries(name, tags, field),
	}
}

// A Histogram measures the distribution of a stream of values.
// Use a histogram to track the distribution of a stream of values (e.g., the
// latency associated with HTTP requests).
type Histogram struct {
	*metrics.Histogram
}

func NewHistogram(name string, tags map[string]string, field string, minValue, maxValue int64) *Histogram {
	s := MakeSeries(name, tags, field)
	return &Histogram{
		metrics.NewHistogram(s, minValue, maxValue, 3),
	}
}

// SnapshotLines provies all collected metrics in Line protocol format. https://github.com/influxdata/influxdb/blob/master/tsdb/README.md
func SnapshotLines() string {
	ss.Lock()
	defer ss.Unlock()
	var buffer bytes.Buffer
	counters, gauges := metrics.Snapshot()
	for c, v := range counters {
		buffer.WriteString(c)
		buffer.WriteString("=")
		buffer.Write([]byte(strconv.FormatUint(v, 10)))
		buffer.Write([]byte("\n"))
	}

	for g, v := range gauges {
		buffer.WriteString(g)
		buffer.WriteString("=")
		buffer.Write([]byte(strconv.FormatInt(v, 10)))
		buffer.Write([]byte("\n"))
	}
	return buffer.String()
}

// Snapshot provides all collected metrics.
func Snapshot() (c map[string]uint64, g map[string]int64) {
	ss.Lock()
	defer ss.Unlock()
	return metrics.Snapshot()
}

// TODO: provide snapshot in json format.
// func SnapshotJson() []byte {
// }

// Reset clears all counters and gauges.
func Reset() {
	metrics.Reset()
}

//MakeSeries creates Series in influxdb format: <measurement>,<tag1>=<key1>,<tagN>=<keyN) <field1>=
func MakeSeries(name string, tags map[string]string, field string) string {
	var keys []string
	for k, _ := range tags {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	s := name
	for _, k := range keys {
		s = s + "," + k + "=" + tags[k]
	}
	return s + " " + field
}
