package kubernetes

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type DaemonSetClient struct {
	client kubernetes.Interface
}

func NewDaemonSetClient(client kubernetes.Interface) *DaemonSetClient {
	return &DaemonSetClient{client: client}
}

func (c *DaemonSetClient) List(ctx context.Context, namespace string) ([]appsv1.DaemonSet, error) {
	list, err := c.client.AppsV1().DaemonSets(namespace).List(ctx, metav1.ListOptions{})

	if err != nil {
		return nil, err
	}

	return list.Items, nil
}
