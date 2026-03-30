package api

import (
	"fmt"

	"github.com/Mujib-Ahasan/Rampaz/internal/metrics"

	pb "github.com/Mujib-Ahasan/Rampaz/proto"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *K8SServer) GetNodeRealTimeStats(req *pb.NodeRequest, stream grpc.ServerStreamingServer[pb.NodeStatsResponse]) error {
	endpoint := "get_node_realtime_stats"
	statusLabel := "success"

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
			WithLabelValues(endpoint, statusLabel).
			Inc()
	}()

	if req == nil {
		statusLabel = "error"
		return status.Error(codes.InvalidArgument, "node stats request cannot be nil")
	}

	if req.NodeName == "" {
		statusLabel = "error"
		return status.Error(codes.InvalidArgument, "node name is required")
	}

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
		statusLabel = "error"
		s.Logger.Error("node realtime stats stream failed", "node", req.NodeName, "err", err)
		return errorHelper(err, fmt.Sprintf("node realtime stats stream for node %q", req.NodeName))

	}

	return nil
}
