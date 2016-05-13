package metrics_test

import (
	"strings"
	"testing"

	"github.com/supershal/stats/metrics"
)

func TestCounter(t *testing.T) {
	metrics.Reset()

	c := metrics.NewCounter(
		"foo",
		map[string]string{
			"bar": "baz"},
		"value")
	c.Add()
	c.AddN(10)

	lines := metrics.SnapshotLines()
	if v, want := lines, "foo,bar=baz value=11\n"; v != want {
		t.Errorf("Counter was %v, but expected %v", v, want)
	}
}

func TestCounterFunc(t *testing.T) {
	metrics.Reset()

	c := metrics.NewCounter(
		"foo",
		map[string]string{
			"bar": "baz"},
		"value")
	c.SetFunc(func() uint64 {
		return 100
	})

	lines := metrics.SnapshotLines()
	if v, want := lines, "foo,bar=baz value=100\n"; v != want {
		t.Errorf("Counter was %v, but expected %v", v, want)
	}
}

func TestCounterBatchFunc(t *testing.T) {
	metrics.Reset()

	c := metrics.NewCounter(
		"foo",
		map[string]string{
			"bar": "baz"},
		"value")

	var a, b uint64

	c.SetBatchFunc(
		"yay",
		func() {
			a, b = 1, 2
		},
		func() uint64 {
			return a
		},
	)

	c1 := metrics.NewCounter(
		"foo1",
		map[string]string{
			"bar1": "baz1"},
		"value1")

	c1.SetBatchFunc(
		"yay",
		func() {
			a, b = 1, 2
		},
		func() uint64 {
			return b
		},
	)

	lines := metrics.SnapshotLines()
	if v, want := lines, "foo,bar=baz value=1\n"; !strings.Contains(v, want) {
		t.Errorf("Counter was %v, but expected %v", v, want)
	}

	if v, want := lines, "foo1,bar1=baz1 value1=2\n"; !strings.Contains(v, want) {
		t.Errorf("Counter was %v, but expected %v", v, want)
	}
}

func TestCounterRemove(t *testing.T) {
	metrics.Reset()

	c := metrics.NewCounter(
		"foo",
		map[string]string{
			"bar": "baz"},
		"value")

	c.Add()
	c.Remove()

	lines := metrics.SnapshotLines()
	if v, want := lines, ""; v != want {
		t.Errorf("Counter was %v, but expected nothing", v)
	}
}

func TestGaugeValue(t *testing.T) {
	metrics.Reset()

	g := metrics.NewGauge(
		"foo",
		map[string]string{
			"bar": "baz"},
		"value")
	g.Set(-100)

	lines := metrics.SnapshotLines()
	if v, want := lines, "foo,bar=baz value=-100\n"; v != want {
		t.Errorf("Gauge was %v, but expected %v", v, want)
	}
}

func TestGaugeFunc(t *testing.T) {
	metrics.Reset()

	g := metrics.NewGauge(
		"foo",
		map[string]string{
			"bar": "baz"},
		"value")

	g.SetFunc(func() int64 {
		return -100
	})

	lines := metrics.SnapshotLines()
	if v, want := lines, "foo,bar=baz value=-100\n"; v != want {
		t.Errorf("Gauge was %v, but expected %v", v, want)
	}
}

func TestGaugeRemove(t *testing.T) {
	metrics.Reset()

	g := metrics.NewGauge(
		"foo",
		map[string]string{
			"bar": "baz"},
		"value")

	g.Set(1)
	g.Remove()

	lines := metrics.SnapshotLines()
	if v, want := lines, ""; v != want {
		t.Errorf("Gauge was %v, but expected %v", v, want)
	}
}

func TestHistogram(t *testing.T) {
	metrics.Reset()

	h := metrics.NewHistogram("foo",
		map[string]string{
			"bar": "baz"},
		"latency",
		1,
		1000)

	for i := 100; i > 0; i-- {
		for j := 0; j < i; j++ {
			h.RecordValue(int64(i))
		}
	}

	lines := metrics.SnapshotLines()

	if v, want := lines, "foo,bar=baz latency.P50=71"; !strings.Contains(v, want) {
		t.Errorf("P50 was %v, but expected %v", v, want)
	}

	if v, want := lines, "foo,bar=baz latency.P75=87"; !strings.Contains(v, want) {
		t.Errorf("P75 was %v, but expected %v", v, want)
	}

	if v, want := lines, "foo,bar=baz latency.P90=95"; !strings.Contains(v, want) {
		t.Errorf("P90 was %v, but expected %v", v, want)
	}

	if v, want := lines, "foo,bar=baz latency.P95=98"; !strings.Contains(v, want) {
		t.Errorf("P95 was %v, but expected %v", v, want)
	}

	if v, want := lines, "foo,bar=baz latency.P99=100"; !strings.Contains(v, want) {
		t.Errorf("P99 was %v, but expected %v", v, want)
	}

	if v, want := lines, "foo,bar=baz latency.P99=100"; !strings.Contains(v, want) {
		t.Errorf("P999 was %v, but expected %v", v, want)
	}
}

func TestHistogramRemove(t *testing.T) {
	metrics.Reset()

	h := metrics.NewHistogram("foo",
		map[string]string{
			"bar": "baz"},
		"latency",
		1,
		1000)
	h.Remove()

	lines := metrics.SnapshotLines()
	if v, want := lines, ""; v != want {
		t.Errorf("Gauge was %v, but expected nothing", v)
	}
}

func BenchmarkCounterAdd(b *testing.B) {
	metrics.Reset()

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			metrics.NewCounter(
				"foo",
				map[string]string{
					"bar": "baz"},
				"value").Add()
		}
	})
}

func BenchmarkCounterAddN(b *testing.B) {
	metrics.Reset()

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			metrics.NewCounter(
				"foo",
				map[string]string{
					"bar": "baz"},
				"value").AddN(100)
		}
	})
}

func BenchmarkGaugeSet(b *testing.B) {
	metrics.Reset()

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			metrics.NewGauge(
				"foo",
				map[string]string{
					"bar": "baz"},
				"value").Set(100)
		}
	})
}

func BenchmarkHistogramRecordValue(b *testing.B) {
	metrics.Reset()
	h := metrics.NewHistogram("foo",
		map[string]string{
			"bar": "baz"},
		"latency",
		1,
		1000)

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			h.RecordValue(100)
		}
	})
}
