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

// Metrics for disk space
var (
	diskSpaceTotal = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "disk_space_total",
		Help: "Total disk space in bytes",
	})

	diskSpaceFree = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "disk_space_free",
		Help: "Free disk space in bytes",
	})

	networkTrafficTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "network_traffic_total",
		Help: "Total network traffic in bytes",
	})

	systemStatus = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "system_status",
		Help: "System status (1 = up, 0 = down)",
	})

	systemInfo = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "system_info",
		Help: "System information",
	}, []string{"os", "arch", "version"})
)

func main() {
	// Create a new registry to hold the metrics
	r := prometheus.NewRegistry()

	// Register the metrics with the registry
	r.MustRegister(onlineUsers)
	r.MustRegister(httpRequestsTotal)
	r.MustRegister(httpDuration)
	r.MustRegister(diskSpaceTotal)
	r.MustRegister(diskSpaceFree)
	r.MustRegister(networkTrafficTotal)
	r.MustRegister(systemStatus)
	r.MustRegister(systemInfo)

	// Start a goroutine to update the online users metric periodically
	go func() {
		for {
			onlineUsers.Set(float64(rand.Intn(500)))
			time.Sleep(time.Second)
		}
	}()

	// Simulate disk space metrics
	go func() {
		for {
			// Simulate total disk space and free disk space
			total := rand.Intn(1000) + 500 // Total disk space between 500 and 1500 bytes
			free := rand.Intn(total)       // Free disk space up to total disk space

			diskSpaceTotal.Set(float64(total))
			diskSpaceFree.Set(float64(free))

			time.Sleep(time.Second)
		}
	}()

	// Simulate network traffic metrics
	go func() {
		for {
			// Simulate network traffic
			traffic := rand.Intn(1000) + 500 // Traffic between 500 and 1500 bytes
			networkTrafficTotal.Add(float64(traffic))

			time.Sleep(time.Second)
		}
	}()

	// Simulate system status
	go func() {
		for {
			// Simulate system status (1 = up, 0 = down)
			status := rand.Intn(2)

			systemStatus.Set(float64(status))

			time.Sleep(time.Second)
		}
	}()

	// Simulate system information
	go func() {
		for {
			// Simulate system information
			os := "linux"
			arch := "amd64"
			version := "1.0.0"

			systemInfo.WithLabelValues(os, arch, version).Set(1)

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
