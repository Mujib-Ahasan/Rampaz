package api

import (
	"context"

	"github.com/Mujib-Ahasan/Rampaz/internal/metrics"
	pb "github.com/Mujib-Ahasan/Rampaz/proto"
	"github.com/prometheus/client_golang/prometheus"
)

func (s *K8SServer) ListServices(ctx context.Context, req *pb.NamespaceRequest) (*pb.ServiceListResponse, error) {

	endpoint := "list_service"
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
	services, err := s.SVCService.ListServices(ctx, req.Namespace)
	if err != nil {
		status = "error"
		return nil, err
	}

	return &pb.ServiceListResponse{
		Services: services,
	}, nil
}
