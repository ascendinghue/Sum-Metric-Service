package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// Metric is a representation of a submitted metric
type Metric struct {
	Key   string    `json:"key"`
	Value int       `json:"value"`
	Time  time.Time `json:"time"`
}

var metrics []Metric

func removeStaleMetrics() {
	for i, metric := range metrics {
		if metric.Time.Before(time.Now().Add(time.Duration(-1) * time.Hour)) {
			metrics[i] = metrics[len(metrics)-1]
			metrics = metrics[:len(metrics)-1]
		}
	}
}

func postMetric(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var metric Metric
	_ = json.NewDecoder(r.Body).Decode(&metric)
	metric.Key = mux.Vars(r)["key"]
	metric.Time = time.Now()

	metrics = append(metrics, metric)
	fmt.Fprintf(w, "{}")
}

func getMetricSum(w http.ResponseWriter, r *http.Request) {
	removeStaleMetrics()

	var total int
	for _, metric := range metrics {
		if mux.Vars(r)["key"] == metric.Key {
			total = total + metric.Value
		}
	}
	fmt.Fprintf(w, "{\"value\":"+strconv.Itoa(total)+"}")
}

func handleRequests() {
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/metric/{key}/sum", getMetricSum).Methods("GET")
	r.HandleFunc("/metric/{key}", postMetric).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", r))
}

func main() {
	handleRequests()
}
