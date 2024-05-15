package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

var onlineUsers = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "goapp_online_users",
	Help: "Online users",
	ConstLabels: map[string]string{
		"course": "funcionameo",
	},
})

var httpRequestsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "goapp_http_requests_total",
	Help: "Count of all HTTP requests for goapp",
}, []string{})

var httpDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Name: "goapp_http_request_duration",
	Help: "Duration in seconds of all HTTP requests",
}, []string{"handler"})

// New metric definitions
var diskSpaceUsage = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "goapp_disk_space_usage",
	Help: "Disk space usage in bytes",
})

var networkTraffic = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "goapp_network_traffic",
	Help: "Network traffic in bytes",
})

var networkTrafficErrors = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "goapp_network_traffic_errors_total",
	Help: "Total number of network traffic errors",
})

var networkTrafficDrops = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "goapp_network_traffic_drops_total",
	Help: "Total number of network traffic drops",
})

var networkSpeed = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "goapp_network_speed",
	Help: "Network speed in Mbps",
})

var filesystemSpaceAvailable = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "goapp_filesystem_space_available",
	Help: "Available filesystem space in bytes",
})

var cpuUsage = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "goapp_cpu_usage",
	Help: "CPU usage in percentage",
})

var fileDescriptor = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "goapp_file_descriptor",
	Help: "Number of file descriptors",
})

var systemStatus = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "goapp_system_status",
	Help: "System status (1 = online, 0 = offline)",
})

var systemInfo = prometheus.NewGaugeVec(prometheus.GaugeOpts{
	Name: "goapp_system_info",
	Help: "Basic system information",
}, []string{"os", "arch", "version"})

var memoryUsage = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "goapp_memory_usage",
	Help: "Memory usage in bytes",
})

func main() {
	r := prometheus.NewRegistry()
	r.MustRegister(onlineUsers)
	r.MustRegister(httpRequestsTotal)
	r.MustRegister(httpDuration)
	r.MustRegister(diskSpaceUsage)
	r.MustRegister(networkTraffic)
	r.MustRegister(networkTrafficErrors)
	r.MustRegister(networkTrafficDrops)
	r.MustRegister(networkSpeed)
	r.MustRegister(filesystemSpaceAvailable)
	r.MustRegister(cpuUsage)
	r.MustRegister(fileDescriptor)
	r.MustRegister(systemStatus)
	r.MustRegister(systemInfo)
	r.MustRegister(memoryUsage)

	go func() {
		for {
			onlineUsers.Set(float64(rand.Intn(500)))
			diskSpaceUsage.Set(randomDiskSpaceUsage()) // Set disk space usage to a random value
			networkTraffic.Set(float64(rand.Intn(1024 * 1024)))        // Simulate network traffic up to 1MB
			networkTrafficErrors.Inc()                                 // Simulate network traffic errors
			networkTrafficDrops.Inc()                                  // Simulate network traffic drops
			networkSpeed.Set(float64(rand.Intn(1000)))                 // Simulate network speed up to 1000 Mbps
			filesystemSpaceAvailable.Set(float64(rand.Intn(1024 * 1024 * 1024))) // Simulate filesystem space available up to 1GB
			cpuUsage.Set(float64(rand.Intn(100)))                      // Simulate CPU usage between 0 and 100%
			fileDescriptor.Set(float64(rand.Intn(500)))                // Simulate number of file descriptors
			systemStatus.Set(float64(rand.Intn(2)))                     // Simulate system status (0 or 1)
			systemInfo.WithLabelValues(runtime.GOOS, runtime.GOARCH, "1.0").Set(1)  // Simulate system information
			memoryUsage.Set(float64(rand.Intn(1024 * 1024 * 1024)))     // Simulate memory usage up to 1GB
			time.Sleep(time.Second)
		}
	}()

	home := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Duration(rand.Intn(8)) * time.Second)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello !!!"))
	})

	contact := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Contact"))
	})

	d := promhttp.InstrumentHandlerDuration(
		httpDuration.MustCurryWith(prometheus.Labels{"handler": "home"}),
		promhttp.InstrumentHandlerCounter(httpRequestsTotal, home),
	)

	d2 := promhttp.InstrumentHandlerDuration(
		httpDuration.MustCurryWith(prometheus.Labels{"handler": "contact"}),
		promhttp.InstrumentHandlerCounter(httpRequestsTotal, contact),
	)

	http.Handle("/", d)
	http.Handle("/contact", d2)
	http.Handle("/metrics", promhttp.HandlerFor(r, promhttp.HandlerOpts{}))
	log.Fatal(http.ListenAndServe(":8181", nil))
}

// randomDiskSpaceUsage simulates the disk space usage with a random value.
func randomDiskSpaceUsage() float64 {
	// Simulate disk space usage between 0 and 500MB
	return float64(rand.Intn(500 * 1024 * 1024))
}
