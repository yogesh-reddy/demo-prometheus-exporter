package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/yogesh-reddy/demo-prometheus-exporter/collector"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	//counter metric
	exCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Namespace: "demo",
			Name:      "simple_counter",
			Help:      "example counter metric",
		})

	//Gauge metric
	exGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "demo",
			Name:      "simple_gauge",
			Help:      "example gauge metric",
		})


	//Histogram metric
	exHistogram = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: "demo",
			Name:      "http_requests_latency",
			Help:      "example gauge metric",
			Buckets: prometheus.LinearBuckets(0 ,0.15,6),
			//prometheus.ExponentialBuckets(0.1,2,5),
		})


)

func main() {

	log.Println("Started exporter")

	prometheus.MustRegister(exCounter,exGauge,exHistogram)

	go AddMetricData()
	prometheus.MustRegister(collector.NewApplicationMetrics())

	//start http server
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8082", nil))
}

func AddMetricData()  {
	go func() {
		for {
			exCounter.Add(rand.Float64() * 5)
			exGauge.Add(rand.Float64()*15 - 5)
			//exSummary.Observe(rand.Float64() * 10)
			time.Sleep(time.Second)
		}
	}()

	go func() {
		for {
			start := time.Now()
			simpleRequest()
			exHistogram.Observe(time.Since(start).Seconds())
			time.Sleep(time.Second)
		}
	}()

}

func simpleRequest() {
	url := "https://jsonplaceholder.typicode.com/todos/1"
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
}
