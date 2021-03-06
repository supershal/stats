package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	gmux "github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/supershal/stats"
)

func main() {
	serve()
}

func serve() {
	app := httptest.NewServer(appHandler())
	defer app.Close()
	fmt.Println("Application server started at:", app.URL)

	statserver := httptest.NewServer(metricsHandler())
	defer statserver.Close()
	fmt.Println("Metrics server started at:", statserver.URL)

	fmt.Println("Making 10 successful request to: ", app.URL, "/app")
	for i := 0; i < 10; i++ {
		http.Get(app.URL + "/app")
	}

	fmt.Println("Making 5 unsuccessful request to: ", app.URL, "/error ")
	for i := 0; i < 5; i++ {
		http.Get(app.URL + "/error")
	}

	fmt.Println("Result from: ", statserver.URL, "/metrics: ")
	res, err := http.Get(statserver.URL + "/metrics")
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer res.Body.Close()
	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
		return
	}

	fmt.Println(string(contents))
}

func middlewares() alice.Chain {
	host, _ := os.Hostname()
	tags := map[string]string{
		"host": host,
	}
	s := stats.NewHTTPStats(tags)
	return alice.New(s.HTTPStatsHandler)
}

func appHandler() *gmux.Router {
	m := middlewares()
	g := gmux.NewRouter()

	handler := func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Duration(rand.Int31n(1000)) * time.Millisecond)
		w.Write([]byte("Hello Stats"))
	}

	errorhandler := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Demo error", http.StatusServiceUnavailable)
	}

	g.Handle("/app", m.ThenFunc(handler)).Methods("GET")
	g.Handle("/error", m.ThenFunc(errorhandler)).Methods("GET")
	return g
}

func metricsHandler() *gmux.Router {
	g := gmux.NewRouter()

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(stats.HTTPMetricsSnapshotLines()))
	}

	g.HandleFunc("/metrics", handler).Methods("GET")
	return g
}
