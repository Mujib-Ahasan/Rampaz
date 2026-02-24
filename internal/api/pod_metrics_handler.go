package api

import (
	"github.com/Mujib-Ahasan/Rampaz/internal/metrics"
	pb "github.com/Mujib-Ahasan/Rampaz/proto"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
)

func (s *K8SServer) GetPodStats(req *pb.PodRequest, stream grpc.ServerStreamingServer[pb.PodStatsResponse]) error {
	endpoint := "get_pod_stats"
	status := "success"

	timer := prometheus.NewTimer(
		metrics.RequestLatency.WithLabelValues(endpoint),
	)

	metrics.ActiveStreams.
		WithLabelValues(endpoint).
		Inc()
	defer metrics.ActiveStreams.
		WithLabelValues(endpoint).
		Dec()
	defer func() {
		timer.ObserveDuration()
		metrics.APIRequests.
			WithLabelValues(endpoint, status).
			Inc()
	}()

	sendWithMetrics := func(resp *pb.PodStatsResponse) error {
		metrics.StreamMessagesSent.
			WithLabelValues(endpoint).
			Inc()
		return stream.Send(resp)
	}
	err := s.PodMetService.StreamPodStats(
		stream.Context(),
		req.Namespace,
		sendWithMetrics,
	)

	if err != nil {
		status = "error"
		return err
	}

	return nil
}
