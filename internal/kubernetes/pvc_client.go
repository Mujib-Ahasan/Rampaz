package kubernetes

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type PVCClient struct {
	client kubernetes.Interface
}

func NewPVCClient(client kubernetes.Interface) *PVCClient {
	return &PVCClient{client: client}
}

func (c *PVCClient) List(ctx context.Context, namespace string) ([]corev1.PersistentVolumeClaim, error) {
	list, err := c.client.CoreV1().PersistentVolumeClaims(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		if namespace == "" {
			return nil, fmt.Errorf("list presistentvolumeclaims across all namespaces: %w", err)
		}
		return nil, fmt.Errorf("list presistentvolumeclaims in namespace %q: %w", namespace, err)
	}
	return list.Items, nil
}
