package kubernetes

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
)

type PodMetricsClient struct {
	metricsClient *metricsv.Clientset
}

func NewPodMetricsClient(metrics *metricsv.Clientset) *PodMetricsClient {
	return &PodMetricsClient{metricsClient: metrics}
}

func (c *PodMetricsClient) GetPodMetrics(ctx context.Context, namespace string) (*v1beta1.PodMetricsList, error) {
	pms, err := c.metricsClient.MetricsV1beta1().PodMetricses(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		if namespace == "" {
			return nil, fmt.Errorf("list pod metrics across all namespaces: %w", err)
		}
		return nil, fmt.Errorf("list pod metrics in namespace %q: %w", namespace, err)
	}

	return pms, nil
}
