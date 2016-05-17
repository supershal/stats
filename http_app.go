package stats

import (
	"net/http"
	"strconv"

	gmux "github.com/gorilla/mux"
)

//ServeMetrics starts Http server to provide current metrics in influxdb line protocol format.
// It takes port number and path as input.: example  ServeMetrics(8081, "/metrics").
func ServeMetrics(port int, path string) {
	g := gmux.NewRouter()
	g.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(HTTPMetricsSnapshotLines()))
	}).Methods("GET")

	addr := ":" + strconv.Itoa(port)
	err := http.ListenAndServe(addr, g)
	if err != nil {
		panic(err)
	}
}
