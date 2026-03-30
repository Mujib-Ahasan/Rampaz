package kubernetes

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type NodeInfoClient struct {
	client *kubernetes.Clientset
}

func NewNodeInfoClient(c *kubernetes.Clientset) *NodeInfoClient {
	return &NodeInfoClient{client: c}
}

func (s *NodeInfoClient) NodeInfo(ctx context.Context, nodeName string) (*corev1.Node, error) {
	nis, err := s.client.CoreV1().Nodes().Get(ctx, nodeName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("get node information for node %q: %w", nodeName, err)
	}

	return nis, nil
}
