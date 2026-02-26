package kubernetes

import (
	"context"

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

func (c *StatefulSetClient) List(ctx context.Context, namespace string) ([]appsv1.StatefulSet, error) {
	list, err := c.client.AppsV1().StatefulSets(namespace).List(ctx, metav1.ListOptions{})

	if err != nil {
		return nil, err
	}

	return list.Items, nil
}
