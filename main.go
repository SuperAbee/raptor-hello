package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"strconv"
	"time"
)

var (
	opsProcessedVec = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "processed_ops_total",
		Help: "The total number of processed events",
	}, []string{"publisher"})
)

func main () {
	go func() {
		for i := 0; ; i++ {
			opsProcessedVec.With(prometheus.Labels{"publisher": strconv.Itoa(i % 10)}).Inc()
			time.Sleep(3 * time.Second)
		}
	}()

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		err := http.ListenAndServe(":8877", nil)
		if err != nil {
			log.Println(err)
		}
	}()

	go func() {
		http.HandleFunc("/hello", HelloHandler)
		err := http.ListenAndServe(":8878", nil)
		if err != nil {
			log.Println(err)
		}
	}()
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	n, err := strconv.ParseInt(r.Header.Get("sleep"), 10, 64)
	if err != nil {
		n = 0
	}
	time.Sleep(time.Duration(n) * time.Second)
	_, _ = fmt.Fprintf(w, "Hello")
	_, _ = fmt.Println("Hello")
}

