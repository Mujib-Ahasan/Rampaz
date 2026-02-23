package api

import (
	"github.com/Mujib-Ahasan/Rampaz/internal/metrics"

	pb "github.com/Mujib-Ahasan/Rampaz/proto"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
)

func (s *K8SServer) GetNodeRealTimeStats(req *pb.NodeRequest, stream grpc.ServerStreamingServer[pb.NodeStatsResponse]) error {
	endpoint := "get_node_realtime_stats"
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
	sendWithMetrics := func(resp *pb.NodeStatsResponse) error {
		err := stream.Send(resp)
		if err == nil {
			metrics.StreamMessagesSent.
				WithLabelValues(endpoint).
				Inc()
		}
		return err
	}

	err := s.NodeMetService.StreamNodeStats(
		stream.Context(),
		req.NodeName,
		sendWithMetrics,
	)

	if err != nil {
		status = "error"
		return err
	}

	return nil
}
