package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// OnlineUsers is a gauge metric that represents the number of online users
var onlineUsers = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "goapp_online_users",
	Help: "Online users",
	ConstLabels: map[string]string{
		"course": "funcionameo",
	},
})

// HttpRequestsTotal is a counter metric that counts the total number of HTTP requests
var httpRequestsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "goapp_http_requests_total",
	Help: "Count of all HTTP requests for goapp",
}, []string{})

// HttpDuration is a histogram metric that measures the duration of HTTP requests
var httpDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Name: "goapp_http_request_duration",
	Help: "Duration in seconds of all HTTP requests",
}, []string{"handler"})

func main() {
	// Create a new registry to hold the metrics
	r := prometheus.NewRegistry()

	// Register the metrics with the registry
	r.MustRegister(onlineUsers)
	r.MustRegister(httpRequestsTotal)
	r.MustRegister(httpDuration)

	// Start a goroutine to update the online users metric periodically
	go func() {
		for {
			onlineUsers.Set(float64(rand.Intn(500)))
			time.Sleep(time.Second)
		}
	}()

	// Define two HTTP handlers for the home and contact pages
	home := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate random processing time for the handler
		time.Sleep(time.Duration(rand.Intn(8)) * time.Second)

		// Set the HTTP response status code to 200 OK
		w.WriteHeader(http.StatusOK)

		// Write the response body
		w.Write([]byte("Hello Acacio"))
	})

	contact := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate random processing time for the handler
		time.Sleep(time.Duration(rand.Intn(5)) * time.Second)

		// Set the HTTP response status code to 200 OK
		w.WriteHeader(http.StatusOK)

		// Write the response body
		w.Write([]byte("Contact"))
	})

	// Instrument the home handler with Prometheus metrics
	d := promhttp.InstrumentHandlerDuration(
		httpDuration.MustCurryWith(prometheus.Labels{"handler": "home"}),
		promhttp.InstrumentHandlerCounter(httpRequestsTotal, home),
	)

	// Instrument the contact handler with Prometheus metrics
	d2 := promhttp.InstrumentHandlerDuration(
		httpDuration.MustCurryWith(prometheus.Labels{"handler": "contact"}),
		promhttp.InstrumentHandlerCounter(httpRequestsTotal, contact),
	)

	// Register the handlers with the HTTP server
	http.Handle("/", d)
	http.Handle("/contact", d2)

	// Serve the Prometheus metrics endpoint
	http.Handle("/metrics", promhttp.HandlerFor(r, promhttp.HandlerOpts{}))

	// Start the HTTP server
	log.Fatal(http.ListenAndServe(":8181", nil))
}
