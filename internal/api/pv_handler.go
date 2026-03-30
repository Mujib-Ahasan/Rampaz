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

	pvs, err := s.PVService.ListPVs(ctx)
	if err != nil {
		statusLabel = "error"
		s.Logger.Error("list persistentvolume failed", "err", err)
		return nil, errorHelper(err, "persistentvolume list")

	}

	return &pb.PVListResponse{
		Pvs: pvs,
	}, nil
}
