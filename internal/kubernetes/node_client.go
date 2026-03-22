package kubernetes

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type NodeClient struct {
	clientset kubernetes.Interface
}

func NewNodeClient(clientset kubernetes.Interface) *NodeClient {
	return &NodeClient{clientset: clientset}
}

func (c *NodeClient) ListNodes(ctx context.Context) ([]corev1.Node, error) {
	nodes, err := c.clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return nodes.Items, nil
}
