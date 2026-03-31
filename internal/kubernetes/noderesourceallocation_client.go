package kubernetes

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	metricsv1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metricsclient "k8s.io/metrics/pkg/client/clientset/versioned"
)

type NodeResourceClient struct {
	kubeClient    kubernetes.Interface
	metricsClient metricsclient.Interface
}

func NewNodeResourceClient(client kubernetes.Interface, metclient metricsclient.Interface) *NodeResourceClient {
	return &NodeResourceClient{kubeClient: client, metricsClient: metclient}
}

func (c *NodeResourceClient) GetNode(ctx context.Context, nodeName string) (*corev1.Node, error) {
	return c.kubeClient.CoreV1().Nodes().Get(ctx, nodeName, metav1.GetOptions{})
}

func (c *NodeResourceClient) ListPodsByNode(ctx context.Context, nodeName string) ([]corev1.Pod, error) {
	pods, err := c.kubeClient.CoreV1().Pods("").List(ctx, metav1.ListOptions{
		FieldSelector: fields.OneTermEqualSelector("spec.nodeName", nodeName).String(),
	})
	if err != nil {
		return nil, err
	}
	return pods.Items, nil
}

func (c *NodeResourceClient) GetNodeMetrics(ctx context.Context, nodeName string) (*metricsv1beta1.NodeMetrics, error) {
	return c.metricsClient.MetricsV1beta1().NodeMetricses().Get(ctx, nodeName, metav1.GetOptions{})
}
