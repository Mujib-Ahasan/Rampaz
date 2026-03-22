package service

import (
	"context"

	networkingv1 "k8s.io/api/networking/v1"

	"github.com/Mujib-Ahasan/Rampaz/internal/kubernetes"

	pb "github.com/Mujib-Ahasan/Rampaz/proto"
)

type NetworkPolicyService struct {
	client *kubernetes.NetworkPolicyClient
}

func NewNetworkPolicyService(client *kubernetes.NetworkPolicyClient) *NetworkPolicyService {
	return &NetworkPolicyService{client: client}
}

func (s *NetworkPolicyService) ListNetworkPolicies(ctx context.Context, namespace string) ([]*pb.NetworkPolicyInfo, error) {
	policies, err := s.client.List(ctx, namespace)
	if err != nil {
		return nil, err
	}

	var result []*pb.NetworkPolicyInfo

	for _, p := range policies {
		result = append(result, transformNetworkPolicy(&p))
	}

	return result, nil
}

func transformNetworkPolicy(p *networkingv1.NetworkPolicy) *pb.NetworkPolicyInfo {
	var policyTypes []string

	for _, t := range p.Spec.PolicyTypes {
		policyTypes = append(policyTypes, string(t))
	}

	return &pb.NetworkPolicyInfo{
		Name:        p.Name,
		Namespace:   p.Namespace,
		PolicyTypes: policyTypes,
		PodSelector: p.Spec.PodSelector.String(),
		Age:         calculateAge(p.CreationTimestamp.Time),
	}
}
