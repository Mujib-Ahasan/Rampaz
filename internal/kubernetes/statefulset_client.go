package kubernetes

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type StatefulSetClient struct {
	client kubernetes.Interface
}

func NewStatefulSetClient(client kubernetes.Interface) *StatefulSetClient {
	return &StatefulSetClient{client: client}
}

func (c *StatefulSetClient) List(ctx context.Context, namespace, labelSelector string) ([]appsv1.StatefulSet, error) {
	list, err := c.client.AppsV1().StatefulSets(namespace).List(ctx, metav1.ListOptions{LabelSelector: labelSelector})

	if err != nil {
		if namespace == "" {
			return nil, fmt.Errorf("list statefulsets across all namespaces with labelSelector %q: %w", labelSelector, err)
		}
		return nil, fmt.Errorf("list statefulsets in namespace %q with labelSelector %q: %w", namespace, labelSelector, err)
	}

	return list.Items, nil
}
