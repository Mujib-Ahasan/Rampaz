package kubernetes

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type NodeClient struct {
	client *kubernetes.Clientset
}

func NewNodeClient(c *kubernetes.Clientset) *NodeClient {
	return &NodeClient{client: c}
}

func (s *NodeClient) ListNodes(ctx context.Context, nodeName string) (*corev1.Node, error) {
	return s.client.CoreV1().Nodes().Get(ctx, nodeName, metav1.GetOptions{})
}
