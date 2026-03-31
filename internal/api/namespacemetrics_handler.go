package api

import (
	"context"

	"github.com/Mujib-Ahasan/Rampaz/internal/metrics"
	pb "github.com/Mujib-Ahasan/Rampaz/proto"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *K8SServer) GetNamespaceMetrics(ctx context.Context, req *pb.NamespaceMetricsRequest) (*pb.NamespaceMetricsResponse, error) {
	endpoint := "get_namespace_metrics"
	statusLabel := "success"

	timer := prometheus.NewTimer(
		metrics.RequestLatency.WithLabelValues(endpoint),
	)

	defer func() {
		timer.ObserveDuration()
		metrics.APIRequests.WithLabelValues(endpoint, statusLabel).Inc()
	}()

	if req == nil {
		statusLabel = "error"
		return nil, status.Error(codes.InvalidArgument, "namespace metrics request cannot be nil")
	}

	if req.Namespace == "" {
		statusLabel = "error"
		return nil, status.Error(codes.InvalidArgument, "namespace is required")
	}

	resp, err := s.NamespaceMetricsService.GetNamespaceMetrics(ctx, req.Namespace)
	if err != nil {
		statusLabel = "error"
		s.Logger.Error("get namespace metrics failed", "namespace", req.Namespace, "err", err)
		return nil, errorHelper(err, "namespace metrics")

	}

	return resp, nil
}
