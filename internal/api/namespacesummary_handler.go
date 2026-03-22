package api

import (
	"context"

	"github.com/Mujib-Ahasan/Rampaz/internal/metrics"
	pb "github.com/Mujib-Ahasan/Rampaz/proto"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *K8SServer) GetNamespaceSummary(ctx context.Context, req *pb.NamespaceRequest) (*pb.NamespaceSummaryResponse, error) {
	endpoint := "get_namespace_summary"
	statusLabel := "success"

	timer := prometheus.NewTimer(
		metrics.RequestLatency.WithLabelValues(endpoint),
	)

	defer func() {
		timer.ObserveDuration()
		metrics.APIRequests.
			WithLabelValues(endpoint, statusLabel).
			Inc()
	}()

	if req.Namespace == "" {
		statusLabel = "error"
		return nil, status.Error(codes.InvalidArgument, "namespace is required")
	}

	resp, err := s.NamespaceSummaryService.GetNamespaceSummary(ctx, req.Namespace)
	if err != nil {
		statusLabel = "error"
		return nil, err
	}

	return resp, nil
}
