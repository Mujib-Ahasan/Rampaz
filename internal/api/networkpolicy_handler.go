package api

import (
	"context"

	"github.com/Mujib-Ahasan/Rampaz/internal/metrics"
	pb "github.com/Mujib-Ahasan/Rampaz/proto"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *K8SServer) ListNetworkPolicies(ctx context.Context, req *pb.NamespaceRequest) (*pb.NetworkPolicyListResponse, error) {
	endpoint := "list_Network_Policies"
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

	if req == nil {
		statusLabel = "error"
		return nil, status.Error(codes.InvalidArgument, "networkpolicy request cannot be nil")
	}

	networkPolicies, err := s.NetworkPolicyService.ListNetworkPolicies(ctx, req.Namespace)
	if err != nil {
		statusLabel = "error"

		s.Logger.Error("get networkpolicy failed", "namespace", req.Namespace, "err", err)
		return nil, errorHelper(err, "networkploicy list")
	}

	return &pb.NetworkPolicyListResponse{
		NetworkPolicies: networkPolicies,
	}, nil
}
