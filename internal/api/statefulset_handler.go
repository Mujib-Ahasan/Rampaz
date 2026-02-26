package api

import (
	"context"

	"github.com/Mujib-Ahasan/Rampaz/internal/metrics"
	pb "github.com/Mujib-Ahasan/Rampaz/proto"
	"github.com/prometheus/client_golang/prometheus"
)

func (s *K8SServer) ListStatefulSets(ctx context.Context, req *pb.NamespaceRequest) (*pb.WorkloadListResponse, error) {

	endpoint := "list_statefulsets"
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
	workloads, err := s.StatefulSetService.List(ctx, req.Namespace)
	if err != nil {
		status = "error"
		return nil, err
	}

	return &pb.WorkloadListResponse{
		Workloads: workloads,
	}, nil
}
