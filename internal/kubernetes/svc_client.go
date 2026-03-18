package kubernetes

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type ServiceClient struct {
	client kubernetes.Interface
}

func NewServiceClient(client kubernetes.Interface) *ServiceClient {
	return &ServiceClient{client: client}
}

func (c *ServiceClient) ListServices(ctx context.Context, namespace string) ([]corev1.Service, error) {
	list, err := c.client.CoreV1().Services(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return list.Items, nil
}
