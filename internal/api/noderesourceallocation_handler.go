package api

import (
	"context"
	"fmt"

	"github.com/Mujib-Ahasan/Rampaz/internal/metrics"
	pb "github.com/Mujib-Ahasan/Rampaz/proto"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *K8SServer) GetNodeResourceAllocation(ctx context.Context, req *pb.NodeRequest) (*pb.NodeResourceAllocationResponse, error) {
	endpoint := "get_node_resource_allocation"
	statusLabel := "success"

	timer := prometheus.NewTimer(
		metrics.RequestLatency.WithLabelValues(endpoint),
	)

	defer func() {
		timer.ObserveDuration()
		metrics.APIRequests.WithLabelValues(endpoint, statusLabel).Inc()
	}()

	if req == nil {
		statusLabel = "error"
		return nil, status.Error(codes.InvalidArgument, "node resource allocation request cannot be nil")
	}

	if req.NodeName == "" {
		statusLabel = "error"
		return nil, status.Error(codes.InvalidArgument, "node name is required")
	}

	resp, err := s.NodeResourceService.GetNodeResourceAllocation(ctx, req.NodeName)
	if err != nil {
		statusLabel = "error"
		s.Logger.Error("node resource allocation failed", "node", req.NodeName, "err", err)
		return nil, errorHelper(err, fmt.Sprintf("node resource allocation for node %q", req.NodeName))

	}

	return resp, nil
}
