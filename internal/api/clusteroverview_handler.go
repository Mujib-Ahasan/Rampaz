package api

import (
	"context"
	"errors"

	"github.com/Mujib-Ahasan/Rampaz/internal/metrics"
	pb "github.com/Mujib-Ahasan/Rampaz/proto"
	"github.com/prometheus/client_golang/prometheus"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		s.Logger.Error("get cluster overview failed", "endpoint", "get_cluster_overview", "err", err)
		return nil, errorHelper(err, "cluster overview")

	}

	return resp, nil
}

func errorHelper(err error, msg string) error {
	if errors.Is(err, context.Canceled) {
		return status.Errorf(codes.Canceled, "%s was canceled", msg)
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return status.Errorf(codes.DeadlineExceeded, "%s timed out", msg)
	}

	return status.Errorf(codes.Internal, "failed to fetch %s", msg)
}
