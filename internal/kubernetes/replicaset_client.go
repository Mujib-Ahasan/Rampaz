package kubernetes

import (
	"context"
	"fmt"

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

func (c *ReplicaSetClient) List(ctx context.Context, namespace, labelSelector string) ([]appsv1.ReplicaSet, error) {
	list, err := c.client.AppsV1().ReplicaSets(namespace).List(ctx, metav1.ListOptions{LabelSelector: labelSelector})

	if err != nil {
		if namespace == "" {
			return nil, fmt.Errorf("list replicasets across all namespaces with labelSelector %q: %w", labelSelector, err)
		}
		return nil, fmt.Errorf("list replicasets in namespace %q with labelSelector %q: %w", namespace, labelSelector, err)
	}

	return list.Items, nil
}
