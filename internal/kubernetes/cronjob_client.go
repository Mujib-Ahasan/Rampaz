package kubernetes

import (
	"context"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type CronJobClient struct {
	client kubernetes.Interface
}

func NewCronJobClient(client kubernetes.Interface) *CronJobClient {
	return &CronJobClient{client: client}
}

func (c *CronJobClient) List(ctx context.Context, namespace string) ([]batchv1.CronJob, error) {
	list, err := c.client.BatchV1().CronJobs(namespace).List(ctx, metav1.ListOptions{})

	if err != nil {
		return nil, err
	}

	return list.Items, nil
}
