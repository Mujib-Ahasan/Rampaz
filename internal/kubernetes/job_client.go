package kubernetes

import (
	"context"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type JobClient struct {
	client kubernetes.Interface
}

func NewJobClient(client kubernetes.Interface) *JobClient {
	return &JobClient{client: client}
}

func (c *JobClient) List(ctx context.Context, namespace string) ([]batchv1.Job, error) {
	list, err := c.client.BatchV1().Jobs(namespace).List(ctx, metav1.ListOptions{})

	if err != nil {
		return nil, err
	}

	return list.Items, nil
}
