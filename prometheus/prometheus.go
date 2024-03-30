package prometheus

import (
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpRequestsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
	)
	memoryUsageGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "memory_usage_bytes",
			Help: "Current memory usage in bytes.",
		},
	)
	processingTimeHistogram = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "http_processing_time_seconds",
			Help:    "Histogram of processing time for HTTP requests.",
			Buckets: prometheus.DefBuckets,
		},
	)
)

func InitPrometheus() {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(memoryUsageGauge)
	prometheus.MustRegister(processingTimeHistogram)
}

func HandleHTTPRequest() {

	httpRequestsTotal.Inc()

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	memoryUsageGauge.Set(float64(memStats.Alloc))
}
func MeasureRequestDuration(c *fiber.Ctx) error {
	start := time.Now()

	// Pass the request to the next handler
	err := c.Next()

	// Measure processing time
	duration := time.Since(start).Seconds()
	processingTimeHistogram.Observe(duration)

	return err
}
