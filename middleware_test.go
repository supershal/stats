package stats

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/supershal/stats/metrics"
)

func TestHTTPResponseStat(t *testing.T) {
	metrics.Reset()
	w := &StatsWriter{
		Writer:    httptest.NewRecorder(),
		StartTime: time.Now(),
	}
	tags := map[string]string{
		"foo": "bar",
	}
	content := []byte("baz")
	w.Write(content)
	f := makeHttpResponseStat()
	f(w, tags)

	c, g := metrics.Snapshot()

	assert.Equal(t, 2, len(c)) // "200" and "total"
	assert.Equal(t, 7, len(g)) // "size" and 6 "latencies"

	assert.Equal(t, uint64(1), c["http_response,foo=bar 200"])
	assert.Equal(t, uint64(1), c["http_response,foo=bar total"])

	assert.Equal(t, int64(len(content)), g["http_response,foo=bar size"])

	assert.Contains(t, g, "http_response,foo=bar latency.P50")
	assert.Contains(t, g, "http_response,foo=bar latency.P75")
	assert.Contains(t, g, "http_response,foo=bar latency.P90")
	assert.Contains(t, g, "http_response,foo=bar latency.P95")
	assert.Contains(t, g, "http_response,foo=bar latency.P99")
	assert.Contains(t, g, "http_response,foo=bar latency.P999")

}

func TestHTTPRequestStat(t *testing.T) {
	metrics.Reset()
	r := http.Request{
		Method: "GET",
	}
	tags := map[string]string{
		"foo": "bar",
	}
	f := makeHttpRequestStat()
	f(r, tags)
	c, _ := metrics.Snapshot()

	assert.Equal(t, 1, len(c))
	assert.Equal(t, uint64(1), c["http_request,foo=bar GET"])
}

func TestHTTPMetricsSnapshotLines(t *testing.T) {
	metrics.Reset()
	w := &StatsWriter{
		Writer:    httptest.NewRecorder(),
		StartTime: time.Now(),
	}
	tags := map[string]string{
		"foo": "bar",
	}
	content := []byte("baz")
	w.Write(content)
	f := makeHttpResponseStat()
	f(w, tags)

	lines := HTTPMetricsSnapshotLines()

	// "200" + "total" + "size" + 6 "latencies"
	assert.Equal(t, 9, len(strings.Split(strings.Trim(lines, "\n"), "\n")))

	assert.Contains(t, lines, "http_response,foo=bar 200=1")
	assert.Contains(t, lines, "http_response,foo=bar total=1")
	assert.Contains(t, lines, "http_response,foo=bar size=3")

	assert.Contains(t, lines, "http_response,foo=bar latency.P50=")
	assert.Contains(t, lines, "http_response,foo=bar latency.P75=")
	assert.Contains(t, lines, "http_response,foo=bar latency.P90=")
	assert.Contains(t, lines, "http_response,foo=bar latency.P95=")
	assert.Contains(t, lines, "http_response,foo=bar latency.P99=")
	assert.Contains(t, lines, "http_response,foo=bar latency.P999=")

}
