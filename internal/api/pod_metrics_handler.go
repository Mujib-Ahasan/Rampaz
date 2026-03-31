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

func (s *K8SServer) GetPodStats(req *pb.PodRequest, stream grpc.ServerStreamingServer[pb.PodStatsResponse]) error {
	endpoint := "get_pod_stats"
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
		return status.Error(codes.InvalidArgument, "pod stats request cannot be nil")
	}

	if req.Namespace == "" {
		statusLabel = "error"
		return status.Error(codes.InvalidArgument, "namespace is required")
	}

	sendWithMetrics := func(resp *pb.PodStatsResponse) error {
		err := stream.Send(resp)
		if err == nil {
			metrics.StreamMessagesSent.
				WithLabelValues(endpoint).
				Inc()
		}
		return err
	}

	err := s.PodMetService.StreamPodStats(stream.Context(), req.Namespace, sendWithMetrics)
	if err != nil {
		statusLabel = "error"
		s.Logger.Error("pod stats stream failed", "namespace", req.Namespace, "err", err)
		return errorHelper(err, fmt.Sprintf("pod stats stream for namespace %q", req.Namespace))
	}

	return nil
}
