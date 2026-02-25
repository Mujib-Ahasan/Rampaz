package kubernetes

import (
	"context"

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

func (c *DeploymentClient) List(ctx context.Context, namespace string) ([]appsv1.Deployment, error) {
	list, err := c.client.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})

	if err != nil {
		return nil, err
	}

	return list.Items, nil
}
