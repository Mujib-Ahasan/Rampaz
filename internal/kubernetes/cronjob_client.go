package kubernetes

import (
	"context"
	"fmt"

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

func (c *CronJobClient) List(ctx context.Context, namespace, labelSelector string) ([]batchv1.CronJob, error) {
	list, err := c.client.BatchV1().CronJobs(namespace).List(ctx, metav1.ListOptions{LabelSelector: labelSelector})

	if err != nil {
		if namespace == "" {
			return nil, fmt.Errorf("list cronjob across all namespaces with labelSelector %q: %w", labelSelector, err)
		}
		return nil, fmt.Errorf("list cronjobs in namespace %q with labelSelector %q: %w", namespace, labelSelector, err)

	}

	return list.Items, nil
}
