package kubernetes

import (
	"context"

	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type NetworkPolicyClient struct {
	clientset kubernetes.Interface
}

func NewNetworkPolicyClient(clientset kubernetes.Interface) *NetworkPolicyClient {
	return &NetworkPolicyClient{clientset: clientset}
}

func (c *NetworkPolicyClient) ListNetworkPolicies(ctx context.Context, namespace string) ([]networkingv1.NetworkPolicy, error) {
	policies, err := c.clientset.NetworkingV1().NetworkPolicies(namespace).List(ctx, metav1.ListOptions{})

	if err != nil {
		return nil, err
	}

	return policies.Items, nil
}
