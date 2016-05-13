package stats

import (
	"net/http"
	"strconv"
	"time"

	"github.com/supershal/stats/metrics"
)

// HTTPStats to provide global tags to each metrics and configure custom request and response stats functions.
type HTTPStats struct {
	globalTags      map[string]string
	logRequestStat  HTTPRequestStatFunc
	logResponseStat HTTPResponseStatFunc
}

// NewHTTPStats provides new instance of HTTPStats
func NewHTTPStats(tags map[string]string, reqFunc HTTPRequestStatFunc, resFunc HTTPResponseStatFunc) *HTTPStats {
	return &HTTPStats{
		globalTags:      tags,
		logRequestStat:  reqFunc,
		logResponseStat: resFunc,
	}
}

// HTTPRequestStatFunc function type to collect HTTP request metrics.
// An application can implement this function to provide custom implementation of request metrics collection.
type HTTPRequestStatFunc func(r http.Request, tags map[string]string)

// httpRequestStat implements HTTPRequestStatFunc. it counts number of requests by Method across all URI Paths.
// If app needs additional tags or per URI stats, the app can implement its own HTTPRequestStatFunc function.
func httpRequestStat(r http.Request, tags map[string]string) {
	field := r.Method
	metrics.NewCounter("http_request", tags, field).Add()
}

// HTTPResponseStatFunc function type to collect HTTP response metrics.
// An application can implement this function to provide custom implementation of response metrics collection.
type HTTPResponseStatFunc func(w http.ResponseWriter, tags map[string]string)

// httpResponseStat implements HTTPResponseStatFunc. It collects response count by rsponse code, response size and latency.
// If app needs additional tags or per response stats, the app can implement its own HTTPResponseStatFunc function.
func httpResponseStat(w http.ResponseWriter, tags map[string]string) {
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
	// TODO: make min/max latency configurable
	var maxLat int64 = 10000
	h := metrics.NewHistogram("http_response", tags, "latency", 0, maxLat)
	lat := rsc.Latency().Nanoseconds() / 1000000
	h.RecordValue(lat)
}

// HTTPStatsHandler is a default provided HTTP middleware function to collect global http request and response stats.
func (s *HTTPStats) HTTPStatsHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// custom response writer
		rlw := &statsWriter{Writer: w, StartTime: time.Now()}
		next.ServeHTTP(rlw, r)
		go s.logRequestStat(*r, s.globalTags)
		go s.logResponseStat(rlw, s.globalTags)
	})
}

// HTTPMetricsSnapshot returns all colleted metrics in metrics Line protocol format.
// https://github.com/influxdata/influxdb/blob/master/tsdb/README.md
func HTTPMetricsSnapshotLines() string {
	return metrics.SnapshotLines()
}

// // StatsHandler is a genric HTTP middleware function to collect http request and response stats.
// // If the application implements custom HTTPRequestStatFunc and HTTPResponseStatFunc then
// // it can use StatsHandler to collect its own custom metrics.
// func (s *Stats) StatsHandler(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// custom response writer
// 		next.ServeHTTP(w, r)
// 		go s.logRequestStat(*r, s.globalTags)
// 		go s.logResponseStat(w, s.globalTags)
// 	})
// }
