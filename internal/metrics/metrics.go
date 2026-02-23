package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	APIRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rampaz_api_requests_total",
			Help: "Total API requests",
		},
		[]string{"endpoint", "status"},
	)

	KubeOperations = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rampaz_kube_operations_total",
			Help: "Kubernetes client operations",
		},
		[]string{"resource", "operation"},
	)

	RequestLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "rampaz_request_duration_seconds",
			Help: "API request latency",
		},
		[]string{"endpoint"},
	)

	ActiveStreams = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "rampaz_active_streams",
			Help: "Number of active streaming connections",
		},
		[]string{"endpoint"},
	)

	StreamMessagesSent = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rampaz_stream_messages_total",
			Help: "Total number of streaming messages sent",
		},
		[]string{"endpoint"},
	)
)

func Init() {
	prometheus.MustRegister(
		APIRequests,
		KubeOperations,
		RequestLatency,
		ActiveStreams,
		StreamMessagesSent,
	)
}
