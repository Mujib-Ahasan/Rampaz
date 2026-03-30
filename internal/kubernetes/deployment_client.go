package kubernetes

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type DeploymentClient struct {
	client kubernetes.Interface
}

func NewDeploymentClient(client kubernetes.Interface) *DeploymentClient {
	return &DeploymentClient{client: client}
}

func (c *DeploymentClient) List(ctx context.Context, namespace, labelSelector string) ([]appsv1.Deployment, error) {
	list, err := c.client.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{LabelSelector: labelSelector})

	if err != nil {
		if namespace == "" {
			return nil, fmt.Errorf("list deployments across all namespaces with labelSelector %q: %w", labelSelector, err)
		}
		return nil, fmt.Errorf("list deployments in namespace %q with labelSelector %q: %w", namespace, labelSelector, err)
	}

	return list.Items, nil
}
