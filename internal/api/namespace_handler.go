package api

import (
	"context"

	"github.com/Mujib-Ahasan/Rampaz/internal/metrics"
	pb "github.com/Mujib-Ahasan/Rampaz/proto"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *K8SServer) ListNamespaces(ctx context.Context, _ *emptypb.Empty) (*pb.NamespaceListResponse, error) {

	endpoint := "list_namespaces"
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
	namespaces, err := s.NamespaceService.ListNamespaces(ctx)
	if err != nil {
		status = "error"
		return nil, err
	}

	return &pb.NamespaceListResponse{
		Namespaces: namespaces,
	}, nil
}
