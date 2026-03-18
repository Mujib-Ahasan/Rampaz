package kubernetes

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type NamespaceClient struct {
	client kubernetes.Interface
}

func NewNamespaceClient(client kubernetes.Interface) *NamespaceClient {
	return &NamespaceClient{client: client}
}

func (c *NamespaceClient) ListNamespaces(ctx context.Context) ([]corev1.Namespace, error) {
	list, err := c.client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}
