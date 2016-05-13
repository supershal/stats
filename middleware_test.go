package stats

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/supershal/stats/metrics"
)

func TestHttpResponseStat(t *testing.T) {
	metrics.Reset()
	w := &statsWriter{
		Writer:    httptest.NewRecorder(),
		StartTime: time.Now(),
	}
	tags := map[string]string{
		"foo": "bar",
	}
	content := []byte("baz")
	w.Write(content)
	httpResponseStat(w, tags)

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
