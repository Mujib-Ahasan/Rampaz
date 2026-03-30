package api

import (
	"context"

	"github.com/Mujib-Ahasan/Rampaz/internal/metrics"
	pb "github.com/Mujib-Ahasan/Rampaz/proto"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *K8SServer) ListNodes(ctx context.Context, _ *emptypb.Empty) (*pb.NodeListResponse, error) {
	endpoint := "list_nodes"
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

	nodes, err := s.NodeService.ListNodes(ctx)
	if err != nil {
		statusLabel = "error"
		s.Logger.Error("list nodes failed", "err", err)
		return nil, errorHelper(err, "node list")

	}

	return &pb.NodeListResponse{
		Nodes: nodes,
	}, nil
}
