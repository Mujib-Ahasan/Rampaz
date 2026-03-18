package api

import (
	"context"

	"github.com/Mujib-Ahasan/Rampaz/internal/metrics"
	pb "github.com/Mujib-Ahasan/Rampaz/proto"
	"github.com/prometheus/client_golang/prometheus"
)

func (s *K8SServer) ListNetworkPolicies(ctx context.Context, req *pb.NamespaceRequest) (*pb.NetworkPolicyListResponse, error) {

	endpoint := "list_Network_Policies"
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
	networkPolicies, err := s.NetworkPolicyService.ListNetworkPolicies(ctx, req.Namespace)
	if err != nil {
		status = "error"
		return nil, err
	}

	return &pb.NetworkPolicyListResponse{
		NetworkPolicies: networkPolicies,
	}, nil
}
