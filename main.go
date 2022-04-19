package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"log"
	"net/http"
	"strconv"
	"time"
)

var (
	taskProcessedDurations = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "task_processed_seconds",
		Help: "任务累计执行的时间总和",
	}, []string{"type"})

	taskProcessing = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "task_processing_total",
		Help: "当前正在执行的任务总数",
	}, []string{"type"})

	taskProcessDurationsHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:        "task_duration_histogram_seconds",
		Help:        "任务执行耗时的柱状图",
		Buckets:     []float64{10, 20, 50, 100},
	}, []string{"code"})

	taskProcessDurationsSummary = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Name:        "task_duration_summary_seconds",
		Help:        "任务执行耗时的分位图",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	}, []string{"code"})

)

func mockData() {
	go func() {
		for i := 0; ; i++ {
			if i % 2 == 0 {
				taskProcessedDurations.With(prometheus.Labels{"type": "GET"}).Add(float64(i % 100))
			} else {
				taskProcessedDurations.With(prometheus.Labels{"type": "POST"}).Add(float64(i % 100))
			}
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for i := 0; ; i++ {
			if i % 2 == 0 {
				taskProcessing.With(prometheus.Labels{"type": "GET"}).Inc()
			} else {
				taskProcessing.With(prometheus.Labels{"type": "POST"}).Inc()
			}
			time.Sleep(1 * time.Second)
			if i % 2 == 0 {
				taskProcessing.With(prometheus.Labels{"type": "GET"}).Dec()
			} else {
				taskProcessing.With(prometheus.Labels{"type": "POST"}).Dec()
			}
		}
	}()

	go func() {
		for i := 0; ; i++ {
			taskProcessDurationsHistogram.With(prometheus.Labels{"code": "200"}).Observe(float64(i % 100))
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for i := 0; ; i++ {
			taskProcessDurationsSummary.With(prometheus.Labels{"code": "200"}).Observe(float64(i % 100))
			time.Sleep(1 * time.Second)
		}
	}()
}

func main () {
	//mockData()

	//go func() {
	//	http.Handle("/metrics", promhttp.Handler())
	//	err := http.ListenAndServe(":8877", nil)
	//	if err != nil {
	//		log.Println(err)
	//	}
	//}()

	http.HandleFunc("/hello", HelloHandler)
	err := http.ListenAndServe(":8878", nil)
	if err != nil {
		log.Println(err)
	}
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	n, err := strconv.ParseInt(r.Header.Get("sleep"), 10, 64)
	if err != nil {
		n = 0
	}
	time.Sleep(time.Duration(n) * time.Second)
	_, _ = fmt.Fprintf(w, "Hello")
	log.Println("Hello")
}

