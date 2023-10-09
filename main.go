package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	count     int
	countLock sync.Mutex

	decreaseCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name:      "down_counter",
			Namespace: "app",
		},
	)
	increaseCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name:      "up_counter",
			Namespace: "app",
		},
	)
)

func main() {
	http.HandleFunc("/up", increaseCount)
	http.HandleFunc("/down", decreaseCount)
	http.HandleFunc("/metrics", promhttp.Handler().ServeHTTP)

	prometheus.MustRegister(increaseCounter)
	prometheus.MustRegister(decreaseCounter)
	http.ListenAndServe(":9999", nil)
}

func increaseCount(w http.ResponseWriter, r *http.Request) {
	countLock.Lock()
	defer countLock.Unlock()
	count = count + 1
	log.Printf("Count increased %d\n", count)
	increaseCounter.Inc()
}

func decreaseCount(w http.ResponseWriter, r *http.Request) {
	countLock.Lock()
	defer countLock.Unlock()
	if count <= 0 {
		http.Error(w, "Cannot be negative", http.StatusBadRequest)
	} else {
		count = count - 1
		log.Printf("Count decreased %d\n", count)
		decreaseCounter.Inc()
	}
}
