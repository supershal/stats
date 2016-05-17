package stats

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/supershal/stats/metrics"
)

// HTTPStats to provide global tags to each metrics and configure custom request and response stats functions.
type HTTPStats struct {
	GlobalTags      map[string]string
	LogRequestStat  HTTPRequestStatFunc
	LogResponseStat HTTPResponseStatFunc
}

// NewHTTPStats provides new instance of HTTPStats
func NewHTTPStats(tags map[string]string) *HTTPStats {
	return &HTTPStats{
		GlobalTags:      tags,
		LogRequestStat:  makeHttpRequestStat(),
		LogResponseStat: makeHttpResponseStat(),
	}
}

// HTTPRequestStatFunc function type to collect HTTP request metrics.
// An application can implement this function to provide custom implementation of request metrics collection.
type HTTPRequestStatFunc func(r http.Request, tags map[string]string)

// makeHttpRequestStat implements a func that returns HTTPRequestStatFunc. it counts number of requests by Method across all URI Paths.
// If app needs additional tags or per URI stats, the app can implement its own HTTPRequestStatFunc function.
func makeHttpRequestStat() HTTPRequestStatFunc {
	return func(r http.Request, tags map[string]string) {
		field := r.Method
		metrics.NewCounter("http_request", tags, field).Add()
	}
}

// HTTPResponseStatFunc function type to collect HTTP response metrics.
// An application can implement this function to provide custom implementation of response metrics collection.
type HTTPResponseStatFunc func(w http.ResponseWriter, tags map[string]string)

// makeHttpResponseStat implements a func that returns HTTPResponseStatFunc. It collects response count by rsponse code, response size and latency.
// If app needs additional tags or per response stats, the app can implement its own HTTPResponseStatFunc function.
func makeHttpResponseStat() HTTPResponseStatFunc {
	// TODO: make min/max latency configurable
	var latencies = make(map[string]*metrics.Histogram)
	var lm sync.Mutex
	// := metrics.NewHistogram("http_response", globalTags, "latency", 0, 10000)

	return func(w http.ResponseWriter, tags map[string]string) {
		var rsc HTTPResponseStatCollector
		var ok bool
		if rsc, ok = w.(HTTPResponseStatCollector); !ok {
			return
		}

		// collect status code counts
		metrics.NewCounter("http_response", tags, strconv.Itoa(rsc.Status())).Add()
		metrics.NewCounter("http_response", tags, "total").Add()

		// collect response size guauge
		metrics.NewGauge("http_response", tags, "size").Set(int64(rsc.Size()))

		// collect response latency histogram
		lat := rsc.Latency().Nanoseconds() / 1000000
		series := metrics.MakeSeries("http_response", tags, "latency")

		lm.Lock()
		if _, ok := latencies[series]; !ok {
			latencies[series] = metrics.NewHistogram("http_response", tags, "latency", 0, 10000)
		}
		lm.Unlock()

		latency := latencies[series]
		latency.RecordValue(lat)
	}
}

// HTTPStatsHandler is a default provided HTTP middleware function to collect global http request and response stats.
func (s *HTTPStats) HTTPStatsHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// custom response writer
		t := time.Now()
		rlw := &StatsWriter{Writer: w, StartTime: t}
		next.ServeHTTP(rlw, r)
		go s.LogRequestStat(*r, s.GlobalTags)
		go s.LogResponseStat(rlw, s.GlobalTags)
	})
}

// ServeHTTP is Negroni compatible interface for httpStats middleware
func (s *HTTPStats) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	rlw := &StatsWriter{Writer: w, StartTime: time.Now()}
	next(rlw, r)
	go s.LogRequestStat(*r, s.GlobalTags)
	go s.LogResponseStat(rlw, s.GlobalTags)
}

// HTTPMetricsSnapshot returns all colleted metrics in metrics Line protocol format.
// https://github.com/influxdata/influxdb/blob/master/tsdb/README.md
func HTTPMetricsSnapshotLines() string {
	return metrics.SnapshotLines()
}
