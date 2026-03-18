package api

import (
	"context"

	"github.com/Mujib-Ahasan/Rampaz/internal/metrics"
	pb "github.com/Mujib-Ahasan/Rampaz/proto"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *K8SServer) ListPVs(ctx context.Context, _ *emptypb.Empty) (*pb.PVListResponse, error) {

	endpoint := "list_PVs"
	status := "success"
	timer := prometheus.NewTimer(
		metrics.RequestLatency.WithLabelValues(endpoint),
	)
	defer func() {
		timer.ObserveDuration()
		metrics.APIRequests.
			WithLabelValues(endpoint, status).
			Inc()
	}()
	pvs, err := s.PVService.ListPVs(ctx)
	if err != nil {
		status = "error"
		return nil, err
	}

	return &pb.PVListResponse{
		Pvs: pvs,
	}, nil
}
