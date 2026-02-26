package kubernetes

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type ReplicaSetClient struct {
	client kubernetes.Interface
}

func NewReplicaSetClient(client kubernetes.Interface) *ReplicaSetClient {
	return &ReplicaSetClient{client: client}
}

func (c *ReplicaSetClient) List(ctx context.Context, namespace string) ([]appsv1.ReplicaSet, error) {
	list, err := c.client.AppsV1().ReplicaSets(namespace).List(ctx, metav1.ListOptions{})

	if err != nil {
		return nil, err
	}

	return list.Items, nil
}
