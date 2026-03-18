package api

import (
	"context"

	"github.com/Mujib-Ahasan/Rampaz/internal/metrics"
	pb "github.com/Mujib-Ahasan/Rampaz/proto"
	"github.com/prometheus/client_golang/prometheus"
)

func (s *K8SServer) ListPVCs(ctx context.Context, req *pb.NamespaceRequest) (*pb.PVCListResponse, error) {

	endpoint := "list_PVCs"
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
	pvcs, err := s.PVCService.ListPVCs(ctx, req.Namespace)
	if err != nil {
		status = "error"
		return nil, err
	}

	return &pb.PVCListResponse{
		Pvcs: pvcs,
	}, nil
}
