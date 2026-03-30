package api

import (
	"context"

	"github.com/Mujib-Ahasan/Rampaz/internal/metrics"
	pb "github.com/Mujib-Ahasan/Rampaz/proto"
	"github.com/prometheus/client_golang/prometheus"
)

func (s *K8SServer) ListPVCs(ctx context.Context, req *pb.NamespaceRequest) (*pb.PVCListResponse, error) {
	endpoint := "list_PVCs"
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
	pvcs, err := s.PVCService.ListPVCs(ctx, req.Namespace)
	if err != nil {
		statusLabel = "error"
		s.Logger.Error("list persistentvolumeclaim failed", "namespace", req.Namespace, "err", err)
		return nil, errorHelper(err, "persistentvolumeclaim list")

	}

	return &pb.PVCListResponse{
		Pvcs: pvcs,
	}, nil
}
