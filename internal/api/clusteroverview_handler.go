package api

import (
	"context"

	"github.com/Mujib-Ahasan/Rampaz/internal/metrics"
	pb "github.com/Mujib-Ahasan/Rampaz/proto"
	"github.com/prometheus/client_golang/prometheus"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *K8SServer) GetClusterOverview(ctx context.Context, _ *emptypb.Empty) (*pb.ClusterOverviewResponse, error) {
	endpoint := "get_cluster_overview"
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

	resp, err := s.ClusterOverviewService.GetClusterOverview(ctx)
	if err != nil {
		statusLabel = "error"
		return nil, err
	}

	return resp, nil
}
