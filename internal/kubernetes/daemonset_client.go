package kubernetes

import (
	"context"
	"fmt"

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

func (c *DaemonSetClient) List(ctx context.Context, namespace, labelSelector string) ([]appsv1.DaemonSet, error) {
	list, err := c.client.AppsV1().DaemonSets(namespace).List(ctx, metav1.ListOptions{LabelSelector: labelSelector})

	if err != nil {
		if namespace == "" {
			return nil, fmt.Errorf("list daemonsets across all namespaces with labelSelector %q: %w", labelSelector, err)
		}
		return nil, fmt.Errorf("list daemonsets in namespace %q with labelSelector %q: %w", namespace, labelSelector, err)

	}

	return list.Items, nil
}
