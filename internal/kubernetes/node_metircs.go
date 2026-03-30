package kubernetes

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
)

type NodeMetricsClient struct {
	metrics *metricsv.Clientset
}

func NewNodeMetricsClient(metrics *metricsv.Clientset) *NodeMetricsClient {
	return &NodeMetricsClient{metrics: metrics}
}

func (c *NodeMetricsClient) GetNodeMetrics(ctx context.Context, nodeName string) (*v1beta1.NodeMetrics, error) {

	nm, err := c.metrics.MetricsV1beta1().NodeMetricses().Get(ctx, nodeName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("get node metrics for node %q: %w", nodeName, err)
	}
	return nm, nil
}
