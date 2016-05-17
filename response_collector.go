package stats

import (
	"net/http"
	"time"
)

// HTTPResponseStatCollector interface provides methods to collect http response metrics.
type HTTPResponseStatCollector interface {
	http.ResponseWriter
	Status() int
	Size() int
	Latency() time.Duration
}

// statsWriter implements HTTPResponseStatCollector and collects response code, size and latency.
type StatsWriter struct {
	Writer    http.ResponseWriter
	StartTime time.Time
	endTime   time.Time
	status    int
	size      int
}

// Header returns the header map that will be sent by
// WriteHeader. Changing the header after a call to
// WriteHeader (or Write) has no effect unless the modified
// headers were declared as trailers by setting the
// "Trailer" header before the call to WriteHeader (see example).
// To suppress implicit response headers, set their value to nil.
func (l *StatsWriter) Header() http.Header {
	return l.Writer.Header()
}

// Write writes the data to the connection as part of an HTTP reply.
// If WriteHeader has not yet been called, Write calls WriteHeader(http.StatusOK)
// before writing the data.  If the Header does not contain a
// Content-Type line, Write adds a Content-Type set to the result of passing
// the initial 512 bytes of written data to DetectContentType.
func (l *StatsWriter) Write(b []byte) (int, error) {
	if l.status == 0 {
		l.status = http.StatusOK
	}
	size, err := l.Writer.Write(b)
	l.size += size
	// add end time for each chunk of Write operation
	if l.endTime.IsZero() {
		l.endTime = l.StartTime
	}
	l.endTime = l.endTime.Add(time.Now().Sub(l.endTime))
	return size, err
}

// WriteHeader sends an HTTP response header with status code.
// If WriteHeader is not called explicitly, the first call to Write
// will trigger an implicit WriteHeader(http.StatusOK).
// Thus explicit calls to WriteHeader are mainly used to
// send error codes.
func (l *StatsWriter) WriteHeader(s int) {
	l.Writer.WriteHeader(s)
	l.status = s
}

// Size function return current response size in bytes
func (l *StatsWriter) Size() int {
	return l.size
}

// Status function returns current http status code
func (l *StatsWriter) Status() int {
	if l.status == 0 {
		return http.StatusOK
	}
	return l.status
}

// Latency provides response time.
func (l *StatsWriter) Latency() time.Duration {
	return l.endTime.Sub(l.StartTime)
}
